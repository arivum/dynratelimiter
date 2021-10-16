/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package xdp

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"code.cloudfoundry.org/bytefmt"
	"github.com/sirupsen/logrus"
)

var (
	errQuantityParsing = errors.New("cores quantity must be a positive integer with a unit of measurement like m, mcores, µ, mcores, or cores")
)

func (t *Thresholds) Init() error {
	var (
		err error
	)

	if err = t.parseCPU(); err != nil {
		return err
	}

	if err = t.parseRAM(); err != nil {
		return err
	}

	return nil
}

func toMilliCores(s string) (float64, error) {
	var (
		i                    int
		quantifier, multiple string
	)
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if i = strings.IndexFunc(s, unicode.IsLetter); i == -1 {
		return 0, errQuantityParsing
	}
	quantifier, multiple = s[:i], s[i:]
	quantity, err := strconv.ParseFloat(quantifier, 64)
	if err != nil || quantity < 0 {
		return 0, errQuantityParsing
	}
	switch multiple {
	case "m", "mcores":
		return float64(quantity), nil
	case "µ", "µcores":
		return float64(quantity * 1e-3), nil
	case "n", "ncores":
		return float64(quantity * 1e-6), nil
	case "cores":
		return float64(quantity * 1e3), nil
	default:
		return 0, errQuantityParsing
	}
}

func (t *Thresholds) parseCPU() error {
	var (
		err        error
		milliCores float64
		replacer   *strings.Replacer
	)

	if len(t.CPU) > 0 {
		if strings.HasSuffix(t.CPU, "%") {
			replacer = strings.NewReplacer("%", "", " ", "")
			if t.cpu, err = strconv.ParseFloat(replacer.Replace(t.CPU), 64); err != nil {
				return err
			}
			t.cpu /= 100.0
		} else if milliCores, err = toMilliCores(t.CPU); err == nil {
			t.cpu = float64(milliCores) / float64(t.resourceTicker.GetCPUMillicores())
		} else {
			return err
		}
	} else {
		t.cpu = 0.0
	}
	if t.cpu > 1.0 {
		logrus.Warn("setting the CPU resource limit to >100%% bypasses rate limiting")
	}
	return nil
}

func (t *Thresholds) parseRAM() error {
	var (
		err      error
		byteSize uint64
		replacer *strings.Replacer
	)

	if len(t.RAM) > 0 {
		if strings.HasSuffix(t.RAM, "%") {
			replacer = strings.NewReplacer("%", "", " ", "")
			if t.ram, err = strconv.ParseFloat(strings.TrimSpace(replacer.Replace(t.RAM)), 64); err != nil {
				return err
			}
			t.ram /= 100.0
		} else if byteSize, err = bytefmt.ToMegabytes(t.RAM); err == nil {
			t.ram = float64(byteSize) / float64(t.resourceTicker.GetRAMLimitMegabytes())
		} else {
			return err
		}
	} else {
		t.ram = 0.0
	}
	if t.ram > 1.0 {
		logrus.Warn("setting the RAM resource limit to >100%% bypasses rate limiting")
	}
	return nil
}
