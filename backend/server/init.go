package server

import (
	"backend/config"
	"backend/controllers"
	"backend/pkg/hs"
	"backend/pkg/log"
	"net/http"
)

var cfg config.CfgInfo

func Init() error {
	var err error

	cfg = config.GetCfg()

	// init router
	r := hs.New()
	r.POST("/get-sticker-rand", controllers.GetStickerRand)

	// start listen
	log.Info("start listen at " + cfg.Server.Addr)
	err = http.ListenAndServe(cfg.Server.Addr, r)
	if err != nil {
		log.Errorf("http.ListenAndServe failed, err=%s", err)
		return err
	}

	return nil
}
