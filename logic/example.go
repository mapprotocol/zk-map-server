package logic

import (
	"github.com/mapprotocol/zk-map-server/entity"
	"github.com/mapprotocol/zk-map-server/resp"
)

func Example(msg string) (ret *entity.ExampleResponse, code int64) {

	return &entity.ExampleResponse{
		Msg: msg,
	}, resp.CodeSuccess
}
