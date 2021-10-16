/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package main

import (
	"flag"

	"github.com/arivum/dynratelimiter/internal/operator"
	"github.com/sirupsen/logrus"
)

var (
	configFile, tlsCertFile, tlsKeyFile string
)

func parseFlags() {
	flag.StringVar(&configFile, "conf", "/etc/dynratelimiter-operator/config.yaml", "Path to the configuration file")
	flag.StringVar(&tlsCertFile, "tls-cert", "/etc/dynratelimiter-operator/tls.crt", "Path to the TLS certificate PEM")
	flag.StringVar(&tlsKeyFile, "tls-key", "/etc/dynratelimiter-operator/tls.key", "Path to the TLS key PEM")
	flag.Parse()
}

func main() {
	var (
		server *operator.Server
		err    error
	)

	parseFlags()

	if server, err = operator.NewServer(configFile, tlsCertFile, tlsKeyFile); err != nil {
		logrus.Fatal(err)
	}
	logrus.Fatal(server.Run())
}
