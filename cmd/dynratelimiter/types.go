/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package cmd

import "github.com/arivum/dynratelimiter/pkg/xdp"

type Version struct {
	Version   string
	GitCommit string
	GoVersion string
	Arch      string
}

type Logging struct {
	Level  string
	Format string
}

type App struct {
	Thresholds *xdp.Thresholds
	Interfaces []string
	Logging    Logging
	xdpState   *xdp.XDPState
}
