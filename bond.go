/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * bond.go
 */

package bond

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	TOP_N = 23
)

func Run(fullVersion string) {
	mongover := flag.String("mongo", "", "MongoDB version, required if connects to config db")
	port := flag.Int("port", 3618, "web server port number")
	ver := flag.Bool("version", false, "print version number")
	verbose := flag.Bool("v", false, "turn on verbose")
	web := flag.Bool("web", false, "starts a web server")

	flag.Parse()
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	if *ver {
		fmt.Println(fullVersion)
		return
	} else if len(flag.Args()) < 1 {
		log.Fatal("connection string is required")
	}

	uri := flag.Arg(0)
	cfg, err := NewConfigDB(uri, *mongover, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	if err = cfg.GetShardingInfo(); err != nil {
		log.Fatal(err)
	}
	cfg.PrintInfo()
	if err = cfg.CheckLogs(); err != nil {
		log.Fatal(err)
	}
	if err = cfg.CheckWarnings(); err != nil {
		log.Fatal(err)
	}
	if *verbose {
		log.Println(StringifyIndent(cfg))
	}

	if !*web {
		if len(flag.Args()) == 0 {
			flag.PrintDefaults()
		}
		return
	}

	router := httprouter.New()
	router.GET("/", InfoHandler)
	router.GET("/api/bond/v1.0/data/:attr", DataHandler)
	router.GET("/favicon.ico", FaviconHandler)
	router.GET("/bond/info", InfoHandler)

	router.GET("/bond/charts/:attr", ChartsHandler)
	router.GET("/bond/chart/shards/:shard", ShardChartHandler)
	router.GET("/bond/chart/namespaces/:ns", NamespaceChartHandler)

	addr := fmt.Sprintf(":%d", *port)
	if listener, err := net.Listen("tcp", addr); err != nil {
		log.Fatal(err)
	} else {
		listener.Close()
		log.Println("starting web server at", addr)
		log.Fatal(http.ListenAndServe(addr, router))
	}
}

func FaviconHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	r.Close = true
	r.Header.Set("Connection", "close")
	w.Header().Set("Content-Type", "image/x-icon")
	ico, _ := base64.StdEncoding.DecodeString(CHEN_ICO)
	w.Write(ico)
}

// Handler responds to API calls
func Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.Redirect(w, r, "/bond/info", http.StatusSeeOther)
}
