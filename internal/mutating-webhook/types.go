/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package mutatingwebhook

import (
	"encoding/json"
)

//[{"op": "add", "path": "/spec/replicas", "value": 3}]

type MutationConfig struct {
	Image string
	Tag   string
}

type jsonPatch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type jsonPatches []jsonPatch

func (j *jsonPatches) ToJSON() []byte {
	var (
		raw []byte
		err error
	)

	if j == nil {
		return []byte("[]")
	}

	if raw, err = json.Marshal(j); err != nil {
		return []byte("[]")
	}
	return raw
}
