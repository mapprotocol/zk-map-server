package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mapprotocol/zk-map-server/controller"
)

func Register(engine *gin.Engine) {
	engine.GET("/proof", controller.GetProof)
}
