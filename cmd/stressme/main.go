/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	highCPU  = false
	highMem  = false
	logLevel string
	host     string
	port     int
)

func parse() {
	var (
		level logrus.Level
		err   error
	)
	flag.BoolVar(&highCPU, "cpu", false, "Stress CPU with each request?")
	flag.BoolVar(&highMem, "mem", false, "Stress Memory with each request?")
	flag.StringVar(&host, "listen-host", "[::]", "Listen host")
	flag.IntVar(&port, "listen-port", 8080, "Listen port")
	flag.StringVar(&logLevel, "loglevel", "info", "Logging level. One of [info, debug, warn, error, fatal]")
	flag.Parse()
	if level, err = logrus.ParseLevel(logLevel); err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(level)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if highMem {
		b := make([]byte, 10, 100<<20)
		if b[0] == 0 {
		}
		logrus.Debugf("allocated %dMiB block", 100<<20)
	}
	if highCPU {
		quit := make(chan bool)
		go func() {
			for {
				select {
				case <-quit:
					return
				default:
				}
			}
		}()
		time.Sleep(10 * time.Millisecond)
		quit <- true
	}
	//time.Sleep(1 * time.Second)
	w.Write([]byte("OK"))
	runtime.GC()
}

func main() {
	var (
		mux    = http.NewServeMux()
		server = http.Server{
			Handler:           mux,
			ReadTimeout:       500 * time.Millisecond,
			WriteTimeout:      2 * time.Second,
			ReadHeaderTimeout: 500 * time.Millisecond,
		}
		output string
	)

	parse()
	runtime.GOMAXPROCS(1)
	mux.HandleFunc("/", handle)
	server.Addr = fmt.Sprintf("%s:%d", host, port)

	output = fmt.Sprintf(`
  ___   _____  ____    ____   ___    ___      _      _  _____ 
/  __ )(_   _)|  _  \ |  __)/  __ )/  __ )   |  \  /  ||  ___)
\  \__   | |  | |_)  )| |__ \  \__ \  \__    |   \/   || |___  
 \__  \  | |  |     / |  __) \__  \ \__  \   | |\  /| ||  ___) 
  _/  /  | |  | |\  \ | |__   _/  /  _/  / _ | | \/ | || |___ 
(____/   |_|  |_| \__)|____)(____/ (____/ |_||_|    |_||_____)

Start listening on %s
`, server.Addr)

	if highMem {
		output += "✔ Memory stress mode active\n"
	}
	if highCPU {
		output += "✔ CPU stress mode active\n"
	}
	fmt.Println(output)
	logrus.Fatal(server.ListenAndServe())

}
