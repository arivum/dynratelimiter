/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import (
	"time"

	"github.com/arivum/resource-ticker/pkg/resources"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf xdp ../../ebpf/dynratelimit.c -- -I../../ebpf/include -O2 -Wall

func NewXDPState(ifaces []string, thresholds *Thresholds) (*XDPState, error) {
	var (
		state = &XDPState{
			objs:       &xdpObjects{},
			thresholds: thresholds,
			ifaceMap:   make(LinkMap),
			rateLimit:  100,
		}
		err   error
		iface string
		links []netlink.Link
		link  netlink.Link
	)

	if thresholds.InitialRateLimit > 0 {
		state.rateLimit = uint32(thresholds.InitialRateLimit)
	}

	if state.resourceTicker, err = resources.NewResourceTicker(resources.WithCPUFloatingAvg(1)); err != nil {
		return nil, err
	}

	state.thresholds.resourceTicker = state.resourceTicker

	if len(ifaces) > 0 {
		for _, iface = range ifaces {
			state.ifaceMap[iface] = nil
		}
	} else {
		if links, err = netlink.LinkList(); err != nil {
			return nil, err
		}

		for _, link = range links {
			state.ifaceMap[link.Attrs().Name] = &link
		}
	}

	if err = state.ifaceMap.init(); err != nil {
		return nil, err
	}

	unix.Setrlimit(unix.RLIMIT_MEMLOCK, &unix.Rlimit{
		Cur: unix.RLIM_INFINITY,
		Max: unix.RLIM_INFINITY,
	})

	if err = loadXdpObjects(state.objs, nil); err != nil {
		return nil, err
	}

	if err = state.attachXdpToLink(); err != nil {
		return nil, err
	}

	return state, nil
}

func (x *XDPState) Run() {
	go x.updateRateLimit()
	go x.terminationListener()

	for {
		time.Sleep(1 * time.Second)
	}
}

func (i LinkMap) init() error {
	var (
		err  error
		name string
	)

	for name = range i {
		if i[name], err = getLinkFromName(name); err != nil {
			return err
		}
		logrus.Debugf("found interface %s", name)
	}

	return nil
}
