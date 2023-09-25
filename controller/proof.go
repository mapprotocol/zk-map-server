package controller

import (
	"github.com/gin-gonic/gin"
	"math/big"

	"github.com/mapprotocol/zk-map-server/entity"
	"github.com/mapprotocol/zk-map-server/logic"
	"github.com/mapprotocol/zk-map-server/resp"
	"github.com/mapprotocol/zk-map-server/utils"
)

func GetProof(c *gin.Context) {
	req := &entity.GetProofRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}
	if utils.IsEmpty(req.Height) {
		resp.ParameterErr(c, "missing parameter height")
		return
	}

	_, ok := new(big.Int).SetString(req.Height, 10)
	if !ok {
		resp.ParameterErr(c, "invalid height")
		return
	}

	ret, code := logic.GetProof(req.Height)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}
