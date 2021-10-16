/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/arivum/dynratelimiter/pkg/xdp"
	"github.com/arivum/resource-ticker/pkg/resources"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApp(version string, commit string) (*App, error) {
	var (
		app          = &App{}
		versionPrint []byte
		err          error
	)

	parseFlags()

	if showVersion {
		versionInformation = Version{
			GitCommit: commit,
			GoVersion: runtime.Version(),
			Version:   version,
			Arch:      runtime.GOARCH,
		}

		versionPrint, _ = json.Marshal(versionInformation)
		fmt.Println(string(versionPrint))
		os.Exit(0)
	}

	if err = app.parse(); err != nil {
		return nil, err
	}

	return app, err
}

func (app *App) Run() {
	app.xdpState.Run()
}

func (app *App) initLogging() {
	var (
		level logrus.Level
		err   error
	)

	if level, err = logrus.ParseLevel(app.Logging.Level); err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	switch app.Logging.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	resources.SetLogger(logrus.StandardLogger())
}

func (app *App) parse() error {
	var (
		absPath string
		err     error
	)

	if len(cpuThreshold) > 0 || len(ramThreshold) > 0 {
		app.Thresholds = &xdp.Thresholds{
			RAM:              ramThreshold,
			CPU:              cpuThreshold,
			InitialRateLimit: int64(initialRateLimit),
		}
		app.Logging = Logging{
			Level:  logLevel,
			Format: "gofmt",
		}
		app.Interfaces = make([]string, 0)
	} else {
		if absPath, err = filepath.Abs(configFile); err != nil {
			return err
		}

		viper.AddConfigPath(filepath.Dir(absPath))
		viper.SetConfigType("yaml")
		viper.SetConfigName(strings.Replace(filepath.Base(absPath), filepath.Ext(absPath), "", 1))

		if err = viper.ReadInConfig(); err != nil {
			return err
		}

		if err = viper.Unmarshal(app); err != nil {
			return err
		}
	}

	app.initLogging()

	if app.xdpState, err = xdp.NewXDPState(app.Interfaces, app.Thresholds); err != nil {
		return err
	}

	if err = app.Thresholds.Init(); err != nil {
		return err
	}

	return nil
}
