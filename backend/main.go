package main

import (
	"backend/config"
	"backend/models"
	"backend/pkg/log"
	"backend/server"
	"backend/utils"
	"flag"
)

func main() {
	flag.Parse()

	var err error
	if err = utils.Init(); err != nil {
		return
	}
	if err = config.Init(""); err != nil {
		return
	}
	log.Init()

	if err = models.Init(); err != nil {
		return
	}

	if err = server.Init(); err != nil {
		return
	}

	models.Close()
}
