package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/zk-map-server/router"
	"github.com/spf13/viper"

	"github.com/mapprotocol/zk-map-server/config"
)

func main() {
	config.Init()

	engine := gin.Default()
	router.Register(engine)
	_ = endless.ListenAndServe(viper.GetString("address"), engine)
}
