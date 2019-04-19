package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-logr/glogr"
	"github.com/ory/graceful"
	"github.com/qiniu-ava/pkg/config"
	"github.com/supremind/hodor/handler"
)

func main() {
	address := flag.String("address", "localhost:8080", "listen address")
	flag.Parse()
	log := glogr.New().WithName("hodor")

	cfg := &handler.Config{}
	if e := config.LoadConfigFile(cfg); e != nil {
		log.Error(e, "parse config file failed")
		os.Exit(1)
	}

	done := make(chan struct{})
	go func() {
		server := &http.Server{
			Addr:    *address,
			Handler: handler.Handler(cfg, log),
		}
		log.Info("starting hodor server", "address", *address)
		if e := graceful.Graceful(server.ListenAndServe, server.Shutdown); e != nil {
			log.Error(e, "hodor is shutted down, but not gracefully")
			return
		}
		log.Info("hodor is shutted down gracefully")

		close(done)
	}()

	<-done
}
