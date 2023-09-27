package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mapprotocol/atlasclient"
	"github.com/mapprotocol/zk-map-server/dao"
	"github.com/mapprotocol/zk-map-server/entity"
	"github.com/mapprotocol/zk-map-server/resource/log"
	"github.com/mapprotocol/zk-map-server/resp"
	"github.com/mapprotocol/zk-map-server/utils"
)

const Interval = 2 * time.Second

const (
	URLStart  = "http://47.242.33.167:18888/start"
	URLStatus = "http://47.242.33.167:18888/status"
)

const (
	RPCAddressDevNetwork  = "http://43.134.183.62:7445"
	RPCAddressTestNetwork = "https://testnet-rpc.maplabs.io"
	RPCAddressMainNetwork = "https://rpc.maplabs.io"
)

const (
	statusPending   = "pending"
	statusFailed    = "failed"
	statusCompleted = "completed"
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
	proof, err := dao.NewProofWithHeight(chainID, height).Get()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Logger().WithField("height", height).WithField("error", err).Error("failed to get proof")
		return nil, resp.CodeInternalServerError
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || proof.IsError() {
		body, err := generateStartBody(chainID, height)
		if err != nil {
			log.Logger().WithField("height", height).WithField("chainID", chainID).
				WithField("error", err).Error("failed to generate start body")
			return nil, resp.CodeProofParameterErr
		}

		id, err := RequestStart(URLStart, body)
		if err != nil {
			log.Logger().WithField("body", body).WithField("error", err).Error("failed to request start")
			return nil, resp.CodeExternalServerError
		}

		if proof.IsError() {
			newProof := &dao.Proof{
				UniqueID: id,
				Status:   dao.ProofStatusPending,
			}
			if err = dao.NewProofWithHeight(chainID, height).Updates(newProof); err != nil {
				log.Logger().WithField("height", height).WithField("uniqueID", id).
					WithField("error", err).Error("failed to create proof")
				return nil, resp.CodeInternalServerError
			}
		} else {
			pf := &dao.Proof{
				ChainID:  chainID,
				Height:   height,
				UniqueID: id,
				Status:   dao.ProofStatusPending,
			}
			err = pf.Create()
			if err != nil && !utils.IsDuplicateError(err.Error()) {
				log.Logger().WithField("chainID", chainID).WithField("height", height).
					WithField("uniqueID", id).WithField("error", err).Error("failed to create proof")
				return nil, resp.CodeInternalServerError
			}
			if err != nil && utils.IsDuplicateError(err.Error()) {
				log.Logger().WithField("chainID", chainID).WithField("height", height).
					WithField("uniqueID", id).WithField("error", err).Warn("duplicated key")
				ret = &entity.GetProofResponse{
					Height: height,
					Status: dao.ProofStatusPending,
				}
				return ret, resp.CodeSuccess
			}
		}
		utils.Go(func() {
			ticker := time.NewTicker(Interval)
			for range ticker.C {
				log.Logger().WithField("uniqueID", id).Info("requesting proof")
				pf, err := dao.NewProofWithUniqueID(id).Get()
				if err != nil {
					log.Logger().WithField("uniqueID", id).WithField("error", err).Error("failed to get proof")
					return
				}
				if !pf.IsPending() {
					return
				}

				result, status, err := RequestStatus(URLStatus, id)
				if err != nil {
					if err := dao.NewProofWithUniqueID(id).SetError(err.Error()); err != nil {
						log.Logger().WithField("uniqueID", id).WithField("error", err).Error("failed to update proof")
						return
					}
				}

				switch status {
				case statusPending:
					log.Logger().WithField("uniqueID", id).Info("proof status is pending")
					continue
				case statusCompleted:
					log.Logger().WithField("uniqueID", id).Info("proof status is completed")
					if err := dao.NewProofWithUniqueID(id).SetCompleted(result.(string)); err != nil {
						log.Logger().WithField("uniqueID", id).
							WithField("proof", result).WithField("error", err).Error("failed to set completed")
						return
					}
				case statusFailed:
					log.Logger().WithField("uniqueID", id).WithField("result", result).Info("proof status is failed")
					if err := dao.NewProofWithUniqueID(id).SetError(statusFailed); err != nil {
						log.Logger().WithField("uniqueID", id).WithField("error", err).Error("failed to set error")
						return
					}
				default:
					log.Logger().WithField("uniqueID", id).WithField("status", status).WithField("result", result).Info("status is unknown")
				}
			}
		})

		ret = &entity.GetProofResponse{
			Height: height,
			Status: dao.ProofStatusPending,
		}
		return ret, resp.CodeSuccess
	}

	ret = &entity.GetProofResponse{
		Height:   height,
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

func RequestStatus(url, id string) (interface{}, string, error) {
	url = fmt.Sprintf("%s/%s", url, id)
	bs, err := utils.Get(url, nil, nil)
	if err != nil {
		return "", "", err
	}

	ret := &response{}
	if err := json.Unmarshal(bs, ret); err != nil {
		return "", "", err
	}
	return ret.Result, ret.Status, nil
}

func IsValidChainID(chainID uint16) bool {
	_, ok := chainID2RPCAddress[chainID]
	return ok
}
