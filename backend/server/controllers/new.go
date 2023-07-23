package controllers

import "backend/config"

type Controller struct {
	ImgURL string
}

func New(cfg *config.CfgInfo) *Controller {
	return &Controller{
		ImgURL: cfg.Server.ImgURL,
	}
}
