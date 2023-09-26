package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/zk-map-server/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/mapprotocol/atlasclient"
	"github.com/mapprotocol/zk-map-server/dao"
	"github.com/mapprotocol/zk-map-server/entity"
	"github.com/mapprotocol/zk-map-server/resource/log"
	"github.com/mapprotocol/zk-map-server/resp"
)

const statusPending = "pending"

// todo from config
const (
	URLStart  = "http://47.242.33.167:18888/start"
	URLStatus = "http://47.242.33.167:18888/status"
)

const (
	RPCAddressDevNetwork  = "http://43.134.183.62:7445"
	RPCAddressTestNetwork = "http://43.134.183.62:7445"
	RPCAddressMainNetwork = "http://43.134.183.62:7445"
)

var chainID2RPCAddress = map[uint16]string{
	212:   RPCAddressTestNetwork,
	213:   RPCAddressDevNetwork,
	22776: RPCAddressMainNetwork,
}

type response struct {
	Id     string      `json:"id"`
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func GetProof(chainID uint16, height string) (ret *entity.GetProofResponse, code int64) {
	proof, err := dao.NewProofWithHeight(height).Get()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 1. 构建 start api 请求参数
		body, err := generateStartBody(chainID, height)
		if err != nil {
			fmt.Println(err)
			return nil, resp.CodeProofParameterErr
		}
		// 2. 发送 start api 请求 并解析数据
		id, err := RequestStart(URLStart, body)
		if err != nil {
			return nil, resp.CodeExternalServerError
		}

		// 3. 将获取的 id 写入数据库
		pf := &dao.Proof{
			Height:   height,
			UniqueID: id,
		}
		if err = pf.Create(); err != nil {
			log.Logger().WithField("height", height).WithField("uniqueID", id).
				WithField("error", err).Error("failed to create proof")
			return nil, resp.CodeInternalServerError
		}

		// 4. TODO 根据 id 获取 proof (使用异步或者定时任务或者两者结合)
		status, pending, err := RequestStatus(URLStatus, id)
		if err != nil {
			if err := dao.NewProofWithUniqueID(id).SetError(err.Error()); err != nil {
				log.Logger().WithField("uniqueID", id).WithField("error", err).Error("failed to update proof")
				return
			}
		}
		if pending {
			// TODO 稍后重试
		}
		if err := dao.NewProofWithUniqueID(id).SetCompleted(status); err != nil {
			log.Logger().WithField("uniqueID", id).
				WithField("proof", status).WithField("error", err).Error("failed to update proof")
			return
		}
	}
	if err != nil {
		log.Logger().WithField("height", height).WithField("error", err).Error("failed to get proof")
		return nil, resp.CodeInternalServerError
	}

	ret = &entity.GetProofResponse{
		Id:       proof.UniqueID,
		Status:   proof.Status,
		ErrorMsg: proof.ErrorMsg,
	}
	if proof.IsCompleted() {
		result := entity.Result{}
		_ = json.Unmarshal([]byte(proof.Proof), &result)
		ret.Result = result
	}
	return ret, resp.CodeSuccess
}

func generateStartBody(chainID uint16, height string) (string, error) {
	h, err := strconv.ParseUint(height, 10, 64)
	if err != nil {
		return "", errors.New("convert string to uint64 failed:" + height + err.Error())
	}
	c, err := atlasclient.Dial(chainID2RPCAddress[chainID])
	if err != nil {
		return "", err
	}
	block, err := c.MAPBlockByNumber(context.Background(), big.NewInt(int64(h)))
	if err != nil {
		return "", err
	}
	data, err := utils.GetProofParamsForBlock1(c.GetClient(), block)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func RequestStart(url, body string) (string, error) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}
	bs, err := utils.Post(url, headers, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	ret := &response{}
	if err := json.Unmarshal(bs, ret); err != nil {
		return "", err
	}
	return ret.Id, nil
}

func RequestStatus(url, id string) (string, bool, error) {
	url = fmt.Sprintf("%s/%s", url, id)
	bs, err := utils.Get(url, nil, nil)
	if err != nil {
		return "", false, err
	}

	ret := &response{}
	if err := json.Unmarshal(bs, ret); err != nil {
		return "", false, err
	}
	if ret.Status == statusPending {
		return "", true, nil
	}
	return ret.Result.(string), false, nil
}

func IsValidChainID(chainID uint16) bool {
	_, ok := chainID2RPCAddress[chainID]
	return ok
}
