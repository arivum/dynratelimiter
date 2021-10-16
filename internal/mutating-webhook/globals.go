/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package mutatingwebhook

import v1core "k8s.io/api/core/v1"

const (
	annotationNamespace    = "dynratelimiter.arifin.io"
	annotationInject       = annotationNamespace + "/inject"
	annotationLoglevel     = annotationNamespace + "/loglevel"
	annotationThresholdCPU = "thresholds." + annotationNamespace + "/cpu"
	annotationThresholdRAM = "thresholds." + annotationNamespace + "/ram"
	initContainerName      = "dynratelimiter-init"
	volumeName             = "dynratelimiter-volume"
	volumeMount            = "/dynratelimiter"
	containerCommand       = volumeMount + "/dynratelimiter"
)

var (
	defaultCaps          = []v1core.Capability{"NET_ADMIN", "SYS_ADMIN", "SYS_RESOURCE", "IPC_LOCK"}
	initContainerCommand = []string{"cp", "-rf", "/usr/bin/dynratelimiter", volumeMount}
)
