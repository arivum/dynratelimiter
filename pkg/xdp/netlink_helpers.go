/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import "github.com/vishvananda/netlink"

func getLinkFromName(intf string) (*netlink.Link, error) {
	var (
		err  error
		link netlink.Link
	)

	if link, err = netlink.LinkByName(intf); err != nil {
		return nil, err
	}
	return &link, nil
}

func getXDPFlagForInterfaceType(linkType string) int {
	switch linkType {
	case "veth":
		return 2
	case "tuntap":
		return 2
	default:
		return 0
	}
}
