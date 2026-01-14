package main

import (
	"AIGO/internal/config"
	"AIGO/internal/router"
)

func main() {
	router.Router.Run(config.Cfg.AppCfg.Port)
}
