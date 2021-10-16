/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import (
	"github.com/arivum/resource-ticker/pkg/resources"
	"github.com/vishvananda/netlink"
)

type LinkMap map[string]*netlink.Link

type XDPState struct {
	objs           *xdpObjects
	ifaceMap       LinkMap
	rateLimit      uint32
	thresholds     *Thresholds
	resourceTicker *resources.ResourceTicker
}

/*
	Thresholds for resources

	* CPU: integer representing the maximal percentage the CPU should reach. Value has to be inside the range of [0.0, 1.0]

	* RAM: integer representing the maximal percentage of used RAM. Value has to be inside the range of [0.0, 1.0]

	* InitialRateLimit: initial amount of maximal acceptable incoming requests
*/
type Thresholds struct {
	CPU              string `yaml:"cpu"`
	cpu              float64
	RAM              string `yaml:"ram"`
	ram              float64
	InitialRateLimit int64 `yaml:"initialRateLimit"`
	resourceTicker   *resources.ResourceTicker
}
