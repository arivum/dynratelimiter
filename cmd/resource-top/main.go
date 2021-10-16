/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arivum/resource-ticker/pkg/resources"
	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		ticker       *resources.ResourceTicker
		err          error
		cpuBar       = newProgressBar("CPU ")
		memBar       = newProgressBar("RAM ")
		pool         *pb.Pool
		resourceChan chan resources.Resources
		errChan      chan error
	)

	if ticker, err = resources.NewResourceTicker(resources.WithCPUFloatingAvg(1)); err != nil {
		logrus.Fatal(err)
	}

	clearScreen()
	fmt.Printf("CPU Cores:\t%.2f cores\n", float64(ticker.GetCPUMillicores())/float64(1000))
	fmt.Printf("Total Memory:\t%.2f GiB\n\n", float64(ticker.GetRAMLimitMegabytes())/float64(1024))

	pool, _ = pb.StartPool(cpuBar, memBar)
	resourceChan, errChan = ticker.Run()

	go func() {
		var (
			c = make(chan os.Signal, 1)
		)

		signal.Notify(c, syscall.SIGINT)
		signal.Notify(c, syscall.SIGTERM)
		for range c {
			pool.Stop()
			cpuBar.Finish()
			memBar.Finish()
			clearScreen()
			os.Exit(0)
		}
	}()

	for {
		select {
		case r := <-resourceChan:
			cpuBar.Set64(int64(r.CPU.Usage * 100.0))
			memBar.Set64(int64(r.RAM.Usage * 100.0))
		case err := <-errChan:
			logrus.Error(err)
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func newProgressBar(prefix string) *pb.ProgressBar {
	var bar = pb.New64(100)
	bar.ShowSpeed = false
	bar.ShowElapsedTime = false
	bar.ShowCounters = false
	bar.ShowTimeLeft = false
	bar.AlwaysUpdate = true
	bar.Empty = color.New(color.FgHiWhite).Sprint("■")
	bar.Current = color.New(color.FgHiGreen).Sprint("■")
	bar.CurrentN = color.New(color.FgHiGreen).Sprint("■")
	bar.Prefix(prefix)
	return bar
}
