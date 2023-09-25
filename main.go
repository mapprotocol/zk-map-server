package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/zk-map-server/resource/db"
	"github.com/mapprotocol/zk-map-server/resource/log"
	"github.com/mapprotocol/zk-map-server/router"
	"github.com/spf13/viper"

	"github.com/mapprotocol/zk-map-server/config"
)

func main() {
	//  init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConf := viper.GetStringMapString("database")
	db.Init(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	engine := gin.Default()
	router.Register(engine)
	_ = endless.ListenAndServe(viper.GetString("address"), engine)
}
