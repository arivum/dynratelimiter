/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package mutatingwebhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	v1admission "k8s.io/api/admission/v1"
	v1core "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *MutationConfig) writeReviewResponse(w http.ResponseWriter, review *v1admission.AdmissionReview) {
	var (
		rawResponse []byte
		err         error
		status      = http.StatusOK
	)
	if rawResponse, err = json.Marshal(review); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(int(review.Response.Result.Code))
		return
	}

	if review.Response.Result != nil && review.Response.Result.Code > 0 {
		status = int(review.Response.Result.Code)
	}

	logrus.Infof("created response %s", string(rawResponse))

	w.WriteHeader(status)
	w.Write(rawResponse)
}

func (m *MutationConfig) HandleMutatingRequest(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		review    = &v1admission.AdmissionReview{}
		pod       *v1core.Pod
		patches   jsonPatches
		patchType = v1admission.PatchTypeJSONPatch
	)

	if review, err = m.parseAdmissionReview(r); err != nil {
		review.Response = &v1admission.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: false,
			Result: &v1meta.Status{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		}
		m.writeReviewResponse(w, review)
		return
	}

	if review.Request == nil {
		review.Response = &v1admission.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: false,
			Result: &v1meta.Status{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		}
		m.writeReviewResponse(w, review)
		return
	}

	if pod, err = m.parsePod(review); err != nil {
		review.Response = &v1admission.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: false,
			Result: &v1meta.Status{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		}
		m.writeReviewResponse(w, review)
		return
	}

	switch review.Request.Operation {
	case v1admission.Create, v1admission.Update:
		if patches, err = m.mutatePod(pod); err != nil {
			review.Response = &v1admission.AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: false,
				Result: &v1meta.Status{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				},
			}
			m.writeReviewResponse(w, review)
			return
		}
		if len(review.Request.UID) == 0 {
			review.Response = &v1admission.AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: false,
				Result: &v1meta.Status{
					Message: "no UID provided in admission request",
					Code:    http.StatusBadRequest,
				},
			}
			m.writeReviewResponse(w, review)
			return
		}

		review.Response = &v1admission.AdmissionResponse{
			UID:       review.Request.UID,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patches.ToJSON(),
		}
		m.writeReviewResponse(w, review)
		return
	default:
		review.Response = &v1admission.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: true,
		}
		m.writeReviewResponse(w, review)
		return
	}
}

func (m *MutationConfig) parseAdmissionReview(r *http.Request) (*v1admission.AdmissionReview, error) {
	var (
		rawBody         []byte
		err             error
		admissionReview = &v1admission.AdmissionReview{}
	)

	defer r.Body.Close()

	if rawBody, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(rawBody, admissionReview); err != nil {
		return nil, err
	}

	return admissionReview, nil
}

func (m *MutationConfig) parsePod(admissionReview *v1admission.AdmissionReview) (*v1core.Pod, error) {
	var (
		pod = &v1core.Pod{}
		err error
	)

	err = json.Unmarshal(admissionReview.Request.Object.Raw, pod)
	return pod, err
}

func (m *MutationConfig) mutatePod(pod *v1core.Pod) (jsonPatches, error) {
	var (
		ok            bool
		injectValue   string
		patches       = make(jsonPatches, 0)
		patchesResult jsonPatches
		patch         *jsonPatch
	)

	if injectValue, ok = pod.ObjectMeta.Annotations[annotationInject]; !ok || injectValue != "enabled" {
		if patch = m.ejectInitContainer(pod.Spec.InitContainers); patch != nil {
			patches = append(patches, *patch)
		}
		if patch = m.ejectVolume(pod.Spec.Volumes); patch != nil {
			patches = append(patches, *patch)
		}
		if patchesResult = m.ejectVolumeMount(pod.Spec.Containers); len(patchesResult) != 0 {
			patches = append(patches, patchesResult...)
		}
		if patchesResult = m.ejectCommand(pod); len(patchesResult) != 0 {
			patches = append(patches, patchesResult...)
		}
	} else {
		if patch = m.injectInitContainer(pod.Spec.InitContainers); patch != nil {
			patches = append(patches, *patch)
		}
		if patch = m.injectVolume(pod.Spec.Volumes); patch != nil {
			patches = append(patches, *patch)
		}
		if patchesResult = m.injectVolumeMount(pod.Spec.Containers); len(patchesResult) != 0 {
			patches = append(patches, patchesResult...)
		}
		if patchesResult = m.injectCommand(pod); len(patchesResult) != 0 {
			patches = append(patches, patchesResult...)
		}
		if patchesResult = m.injectCaps(pod); len(patchesResult) != 0 {
			patches = append(patches, patchesResult...)
		}
	}
	return patches, nil
}

func (m *MutationConfig) injectCaps(pod *v1core.Pod) jsonPatches {
	var (
		container v1core.Container
		index     int
		patches   = make(jsonPatches, 0)
		caps      []v1core.Capability
		cap       v1core.Capability
		capIndex  int
		found     = false
	)

	for index, container = range pod.Spec.Containers {
		caps = make([]v1core.Capability, 0)

		if container.SecurityContext == nil {
			caps = defaultCaps
		} else if container.SecurityContext.Capabilities == nil {
			caps = defaultCaps
		} else {
			for capIndex = range defaultCaps {
				found = false
				for _, cap = range container.SecurityContext.Capabilities.Add {
					if defaultCaps[capIndex] == cap {
						found = true
					}
				}
				if !found {
					caps = append(caps, defaultCaps[capIndex])
				}
			}
		}

		patches = append(patches, jsonPatch{
			Op:   "replace",
			Path: fmt.Sprintf("/spec/containers/%d/securityContext", index),
			Value: v1core.SecurityContext{
				Capabilities: &v1core.Capabilities{
					Add: caps,
				},
			},
		})
	}

	return patches
}

func (m *MutationConfig) injectVolumeMount(containers []v1core.Container) jsonPatches {
	var (
		container v1core.Container
		mount     v1core.VolumeMount
		index     int
		patches   = make(jsonPatches, 0)
		found     = false
	)

	for index, container = range containers {
		found = false
		for _, mount = range container.VolumeMounts {
			if mount.Name == volumeName {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if len(container.VolumeMounts) == 0 {
			patches = append(patches, jsonPatch{
				Op:   "add",
				Path: fmt.Sprintf("/spec/containers/%d/volumeMounts", index),
				Value: []v1core.VolumeMount{
					{
						Name:      volumeName,
						MountPath: volumeMount,
					},
				},
			})
		} else {
			patches = append(patches, jsonPatch{
				Op:   "add",
				Path: fmt.Sprintf("/spec/containers/%d/volumeMounts/0", index),
				Value: v1core.VolumeMount{
					Name:      volumeName,
					MountPath: volumeMount,
				},
			})
		}
	}

	return patches
}

func (m *MutationConfig) ejectVolumeMount(containers []v1core.Container) jsonPatches {
	var (
		container  v1core.Container
		mount      v1core.VolumeMount
		index      int
		mountIndex int
		patches    = make(jsonPatches, 0)
	)

	for index, container = range containers {
		for mountIndex, mount = range container.VolumeMounts {
			if mount.Name == volumeName {
				patches = append(patches, jsonPatch{
					Op:   "remove",
					Path: fmt.Sprintf("/spec/container/%d/volumeMounts/%d", index, mountIndex),
				})
			}
		}
	}

	return patches
}

func (m *MutationConfig) injectCommand(pod *v1core.Pod) jsonPatches {
	var (
		container      v1core.Container
		index          int
		cmdPrefix      = []string{containerCommand}
		patches        = make(jsonPatches, 0)
		annotation     string
		found          bool
		seperatorIndex int
	)

	if annotation, found = pod.ObjectMeta.Annotations[annotationThresholdCPU]; found && len(annotation) > 0 {
		cmdPrefix = append(cmdPrefix, "-cpu", annotation)
	}
	if annotation, found = pod.ObjectMeta.Annotations[annotationThresholdRAM]; found && len(annotation) > 0 {
		cmdPrefix = append(cmdPrefix, "-ram", annotation)
	}
	if annotation, found = pod.ObjectMeta.Annotations[annotationLoglevel]; found && len(annotation) > 0 {
		cmdPrefix = append(cmdPrefix, "-loglevel", annotation)
	}
	cmdPrefix = append(cmdPrefix, "--")

	for index, container = range pod.Spec.Containers {
		if container.Command[0] != containerCommand {
			patches = append(patches, jsonPatch{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%d/command", index),
				Value: append(cmdPrefix, container.Command...),
			})
		} else {
			for seperatorIndex = 0; seperatorIndex < len(container.Command); seperatorIndex++ {
				if container.Command[seperatorIndex] == "--" {
					break
				}
			}
			if seperatorIndex == len(container.Command) {
				continue
			}
			patches = append(patches, jsonPatch{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%d/command", index),
				Value: append(cmdPrefix, container.Command[seperatorIndex:]...),
			})
		}
	}
	return patches
}

func (m *MutationConfig) ejectCommand(pod *v1core.Pod) jsonPatches {
	var (
		container      v1core.Container
		index          int
		patches        = make(jsonPatches, 0)
		seperatorIndex int
	)

	for index, container = range pod.Spec.Containers {
		if container.Command[0] == containerCommand {
			seperatorIndex = 0
			for seperatorIndex = 0; seperatorIndex < len(container.Command); seperatorIndex++ {
				if container.Command[seperatorIndex] == "--" {
					break
				}
			}
			if seperatorIndex == len(container.Command) {
				continue
			}
			patches = append(patches, jsonPatch{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%d/command", index),
				Value: container.Command[seperatorIndex:],
			})
		}
	}
	return patches
}

func (m *MutationConfig) injectVolume(volumes []v1core.Volume) *jsonPatch {
	var (
		volume v1core.Volume
	)

	for _, volume = range volumes {
		if volume.Name == volumeName {
			return nil
		}
	}
	if len(volumes) == 0 {
		return &jsonPatch{
			Op:   "add",
			Path: "/spec/volumes",
			Value: []v1core.Volume{
				{
					Name: volumeName,
					VolumeSource: v1core.VolumeSource{
						EmptyDir: &v1core.EmptyDirVolumeSource{},
					},
				},
			},
		}
	} else {
		return &jsonPatch{
			Op:   "add",
			Path: "/spec/volumes/0",
			Value: v1core.Volume{
				Name: volumeName,
				VolumeSource: v1core.VolumeSource{
					EmptyDir: &v1core.EmptyDirVolumeSource{},
				},
			},
		}
	}
}

func (m *MutationConfig) ejectVolume(volumes []v1core.Volume) *jsonPatch {
	var (
		volume v1core.Volume
		index  int
	)

	for index, volume = range volumes {
		if volume.Name == volumeName {
			return &jsonPatch{
				Op:   "remove",
				Path: fmt.Sprintf("/spec/volumes/%d", index),
			}
		}
	}

	return nil
}

func (m *MutationConfig) injectInitContainer(initContainers []v1core.Container) *jsonPatch {
	var (
		initContainer v1core.Container
	)

	for _, initContainer = range initContainers {
		if initContainer.Name == initContainerName {
			return nil
		}
	}

	if len(initContainers) == 0 {
		return &jsonPatch{
			Op:   "add",
			Path: "/spec/initContainers",
			Value: []v1core.Container{
				{
					Name:            initContainerName,
					Image:           m.Image + ":" + m.Tag,
					Command:         initContainerCommand,
					ImagePullPolicy: v1core.PullIfNotPresent,
					VolumeMounts: []v1core.VolumeMount{
						{
							Name:      volumeName,
							MountPath: volumeMount,
						},
					},
				},
			},
		}
	} else {
		return &jsonPatch{
			Op:   "add",
			Path: "/spec/initContainers/0",
			Value: v1core.Container{
				Name:            initContainerName,
				Image:           m.Image + ":" + m.Tag,
				Command:         initContainerCommand,
				ImagePullPolicy: v1core.PullIfNotPresent,
				VolumeMounts: []v1core.VolumeMount{
					{
						Name:      volumeName,
						MountPath: volumeMount,
					},
				},
			},
		}
	}
}

func (m *MutationConfig) ejectInitContainer(initContainers []v1core.Container) *jsonPatch {
	var (
		initContainer v1core.Container
		index         int
	)

	for index, initContainer = range initContainers {
		if initContainer.Name == initContainerName {
			return &jsonPatch{
				Op:   "remove",
				Path: fmt.Sprintf("/spec/initContainers/%d", index),
			}
		}
	}
	return nil
}
