package main

import (
	"adm_bkd/api"
	"adm_bkd/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	api.RegisterHealthAPI(r)
	api.RegisterTeachDocAPI(r, cfg)

	r.Run(cfg.Server.Addr)
}
