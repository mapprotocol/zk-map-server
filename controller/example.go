package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/mapprotocol/zk-map-server/entity"
	"github.com/mapprotocol/zk-map-server/logic"
	"github.com/mapprotocol/zk-map-server/resp"
	"github.com/mapprotocol/zk-map-server/utils"
)

func Example(c *gin.Context) {
	req := &entity.ExampleRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}
	if utils.IsEmpty(req.Msg) {
		resp.ParameterErr(c, "missing parameter msg")
		return
	}

	ret, code := logic.Example(req.Msg)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}
