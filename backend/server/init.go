package server

import (
	"net/http"

	"backend/config"
	"backend/controllers"
	"backend/pkg/hs"
	"backend/pkg/log"
)

var cfg config.CfgInfo

func Init() error {
	var err error

	cfg = config.GetCfg()

	// init router
	r, err := hs.New()
	if err != nil {
		log.Errorf("hs.New failed, err=%s", err)
		return err
	}
	r.GET("/get-sticker-rand", controllers.GetStickerRand)

	// start listen
	log.Info("start listen at " + cfg.Server.Addr)
	err = http.ListenAndServe(cfg.Server.Addr, r)
	if err != nil {
		log.Errorf("http.ListenAndServe failed, err=%s", err)
		return err
	}

	return nil
}
