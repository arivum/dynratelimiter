/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package cmd

import (
	"flag"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func parseFlags() {
	flag.StringVar(&configFile, "conf", "dynratelimit.yaml", "Path to config file")
	flag.StringVar(&cpuThreshold, "cpu", "", "Rate limiting threshold based on CPU usage. Notation must be either e.g. 80%% or 300ms")
	flag.StringVar(&ramThreshold, "ram", "", "Rate limiting threshold based on RAM usage. Notation must be either e.g. 80%%, 2048MiB, 2048MB, 2GiB or 2GB")
	flag.StringVar(&logLevel, "loglevel", "info", "Set loglevel to one of [info, debug, warn, error, trace]")
	flag.IntVar(&initialRateLimit, "init-rate-limit", 10, "Initial rate limit. It is recommended to start small.")
	flag.BoolVar(&showVersion, "v", false, "Display version")
	flag.Parse()

	go func() {
		var (
			args    []string
			cmd     *exec.Cmd
			err     error
			ok      bool
			exitErr *exec.ExitError
			status  syscall.WaitStatus
		)

		if args = flag.Args(); len(args) > 0 {
			cmd = exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err = cmd.Start(); err != nil {
				logrus.Fatal(err)
			}

			if err = cmd.Wait(); err != nil {
				if exitErr, ok = err.(*exec.ExitError); ok {
					if status, ok = exitErr.Sys().(syscall.WaitStatus); ok {
						os.Exit(status.ExitStatus())
					}
				} else {
					logrus.Fatalf("cmd.Wait: %v", err)
				}
			}
			os.Exit(0)
		}
	}()
}
