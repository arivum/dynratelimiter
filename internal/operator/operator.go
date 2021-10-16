/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package operator

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func NewServer(configFile string, tlsCertFile string, tlsKeyFile string) (*Server, error) {
	var (
		server = &Server{
			configFile:  configFile,
			tlsKeyFile:  tlsKeyFile,
			tlsCertFile: tlsCertFile,
		}
		err error
	)

	if err = server.init(); err != nil {
		return nil, err
	}
	server.initLogging()

	return server, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServeTLS("", "")
}

func (s *Server) init() error {
	var (
		err       error
		rawConfig []byte
		handler   *http.ServeMux
		cert      tls.Certificate
	)

	if cert, err = tls.LoadX509KeyPair(s.tlsCertFile, s.tlsKeyFile); err != nil {
		return err
	}

	if rawConfig, err = os.ReadFile(s.configFile); err != nil {
		return err
	}

	if err = yaml.Unmarshal(rawConfig, s); err != nil {
		return err
	}

	if len(s.MutationConfig.Image) == 0 {
		return errors.New("please specify the image that should be injected")
	}
	if len(s.MutationConfig.Tag) == 0 {
		s.MutationConfig.Tag = "latest"
	}

	handler = http.NewServeMux()
	handler.HandleFunc("/mutate", s.MutationConfig.HandleMutatingRequest)
	handler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	s.httpServer = &http.Server{
		Addr:      fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler:   handler,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
	}

	return nil
}

func (s *Server) initLogging() {
	var (
		level  logrus.Level
		format logrus.Formatter
		err    error
	)

	if level, err = logrus.ParseLevel(s.Logging.Level); err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	switch s.Logging.Format {
	case "json":
		format = &logrus.JSONFormatter{}
	default:
		format = &logrus.TextFormatter{}
	}
	logrus.SetFormatter(format)
}
