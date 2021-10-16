/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import (
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func (x *XDPState) attachXdpToLink() error {
	var (
		errs, err error
		link      *netlink.Link
	)

	logrus.Debug("attaching XDP to interfaces")
	for _, link = range x.ifaceMap {
		if err = netlink.LinkSetXdpFdWithFlags(*link, x.objs.XdpDynratelimit.FD(), getXDPFlagForInterfaceType((*link).Type())); err != nil {
			errs = multierror.Append(errs, err)
		}
		if err == nil {
			logrus.Infof("successfully attached dynratelimiter to interface %s", (*link).Attrs().Name)
		} else {
			logrus.Errorf("failed attaching dynratelimiter to interface %s with error %v", (*link).Attrs().Name, err)
		}
	}

	return errs
}

func (x *XDPState) detachXdpFromLink() error {
	var (
		err, errs error
		link      *netlink.Link
	)

	logrus.Debug("detaching XDP from interfaces")
	for _, link = range x.ifaceMap {
		if err = netlink.LinkSetXdpFdWithFlags(*link, -1, getXDPFlagForInterfaceType((*link).Type())); err != nil {
			errs = multierror.Append(errs, err)
		}
		if err == nil {
			logrus.Infof("successfully detached dynratelimiter from interface %s", (*link).Attrs().Name)
		} else {
			logrus.Errorf("failed detaching dynratelimiter from interface %s with error %v", (*link).Attrs().Name, err)
		}
	}

	return errs
}
