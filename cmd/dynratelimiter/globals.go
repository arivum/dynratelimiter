/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package cmd

var (
	configFile                           string
	showVersion                          bool
	versionInformation                   = Version{}
	cpuThreshold, ramThreshold, logLevel string
	initialRateLimit                     int
)
