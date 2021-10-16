/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import (
	"math"

	"github.com/arivum/resource-ticker/pkg/resources"
	"github.com/sirupsen/logrus"
)

func (x *XDPState) calculateRatelimit(res resources.Resources) error {
	var (
		deltaCPURateLimit, deltaRAMRateLimit float64
		accRequests                          uint32
		err                                  error
	)

	if err = x.objs.DynratelimitMap.Lookup(acceptedRequests, &accRequests); err != nil {
		return err
	}

	if x.shouldRecalculate(accRequests, res) {
		if x.thresholds.cpu > 0 {
			deltaCPURateLimit = x.calculateCPURateLimitDelta(res.CPU)
		}

		if x.thresholds.ram > 0 {
			deltaRAMRateLimit = x.calculateRAMRateLimitDelta(res.RAM)
		}
		x.rateLimit = x.rateLimit + uint32(math.Min(deltaCPURateLimit, deltaRAMRateLimit))
	}

	return nil
}

func (x *XDPState) updateRateLimit() {
	var (
		err           error
		resourceEvent chan resources.Resources
		resourceTick  resources.Resources
		errChan       chan error
	)

	resourceEvent, errChan = x.resourceTicker.Run()

	for {
		select {
		case resourceTick = <-resourceEvent:
			if logrus.GetLevel() == logrus.DebugLevel {
				x.printStatus()
			}
			logrus.Debugf("%+v", resourceTick.RAM)
			logrus.Debugf("%+v", resourceTick.CPU)
			if err = x.calculateRatelimit(resourceTick); err != nil {
				continue
			}
			if err = x.objs.DynratelimitMap.Put(ratelimit, x.rateLimit); err != nil {
				logrus.Error(err)
			}

			if err = x.objs.DynratelimitMap.Put(acceptedRequests, uint32(0)); err != nil {
				logrus.Error(err)
			}

		case err = <-errChan:
			logrus.Error(err)
		}
	}
}

func (x *XDPState) calculateCPURateLimitDelta(cpu *resources.CPU) float64 {
	return float64(x.rateLimit) * (x.thresholds.cpu - cpu.Usage)
}

func (x *XDPState) calculateRAMRateLimitDelta(ram *resources.RAM) float64 {
	return float64(x.rateLimit) * (x.thresholds.ram - ram.Usage)
}

func (x *XDPState) shouldRecalculate(accRequests uint32, res resources.Resources) bool {
	switch {
	case float64(accRequests) >= 0.8*float64(x.rateLimit):
		return true
	case float64(res.RAM.Used) > x.thresholds.ram*float64(res.RAM.Total):
		return true
	case res.CPU.Usage > x.thresholds.cpu:
		return true
	default:
		return false
	}
}
