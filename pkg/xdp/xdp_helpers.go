/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func (x *XDPState) terminationListener() {
	var (
		c   = make(chan os.Signal, 1)
		err error
	)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	for range c {
		if err = x.detachXdpFromLink(); err != nil {
			logrus.Errorf("failed detaching XDP from interfaces failed with error: %v", err)
		}
		os.Exit(0)
	}
}

func (x *XDPState) printStatus() {
	var (
		l = uint32(0)
		r = uint32(0)
	)

	x.objs.DynratelimitMap.Lookup(ratelimit, &l)
	x.objs.DynratelimitMap.Lookup(acceptedRequests, &r)
	logrus.Debugf("requests: %d, limit: %d", r, l)
}
