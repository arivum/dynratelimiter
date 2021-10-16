/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package operator

import (
	"net/http"

	mutatingwebhook "github.com/arivum/dynratelimiter/internal/mutating-webhook"
)

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Server struct {
	httpServer     *http.Server
	configFile     string
	tlsKeyFile     string
	tlsCertFile    string
	Host           string                         `yaml:"host"`
	Port           int                            `yaml:"port"`
	Logging        Logging                        `yaml:"logging"`
	MutationConfig mutatingwebhook.MutationConfig `yaml:"mutationConfig"`
}
