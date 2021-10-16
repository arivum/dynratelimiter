/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package main

import (
	cmd "github.com/arivum/dynratelimiter/cmd/dynratelimiter"
	"github.com/sirupsen/logrus"
)

var (
	GitCommit = "latest"
	Version   = "0.0.0"
)

func main() {
	var (
		app *cmd.App
		err error
	)

	if app, err = cmd.NewApp(Version, GitCommit); err != nil {
		logrus.Fatal(err)
	}

	app.Run()
}
