package main

import (
	"flag"
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"backend/pkg/log"
	"backend/server"
	"backend/utils"
)

func main() {
	flag.Parse()

	log.Init()

	err := utils.Init()
	if err != nil {
		log.Panicf("utils.Init failed: %s", err)
	}

	cfg, err := config.Init()
	if err != nil {
		log.Panicf("config.Init failed: %s", err)
	}

	err = models.Init(cfg)
	if err != nil {
		log.Panicf("models.Init failed: %s", err)
	}

	s, err := server.New(cfg)
	if err != nil {
		log.Panicf("server.New failed: %s", err)
	}

	hs := http.Server{
		Addr:        cfg.Server.Addr,
		Handler:     s,
		ReadTimeout: 60 * time.Second,
	}

	log.Infof("server start at %s", cfg.Server.Addr)
	err = hs.ListenAndServe()
	if err != nil {
		log.Errorf("server.ListenAndServe failed: %s", err)
	}
}
