package main

import (
	"context"
	"flag"
	"time"

	"backend/app/config"
	"backend/app/pkg/log"
	"backend/app/server"
	"backend/app/utils"

	_ "github.com/lib/pq"
)

func main() {
	time.Local = time.UTC

	ctx := context.Background()

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

	err = server.NewAndRun(ctx, cfg)
	if err != nil {
		log.Panicf("server.New failed: %s", err)
	}
}
