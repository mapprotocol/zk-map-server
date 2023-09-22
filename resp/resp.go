package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mapprotocol/zk-map-server/utils"
)

var EmptyStruct = struct{}{}

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ListData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: data,
	})
}

func SuccessList(c *gin.Context, list interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: ListData{
			List:  list,
			Total: total,
		},
	})
}

func Error(c *gin.Context, code int64) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  code2msg[code],
		Data: EmptyStruct,
	})
}

func ParameterErr(c *gin.Context, msg string) {
	if utils.IsEmpty(msg) {
		msg = code2msg[CodeParameterErr]
	}
	c.JSON(http.StatusOK, Response{
		Code: CodeParameterErr,
		Msg:  msg,
		Data: EmptyStruct,
	})
}
