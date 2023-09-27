package dao

import (
	"github.com/mapprotocol/zk-map-server/resource/db"
	"time"
)

const (
	ProofStatusPending = iota + 1
	ProofStatusError
	ProofStatusCompleted
)

type Proof struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	ChainID   uint16    `gorm:"column:chain_id" json:"chain_id"`
	Height    string    `gorm:"column:height" json:"height"`
	UniqueID  string    `gorm:"column:unique_id" json:"unique_id"`
	Status    uint8     `gorm:"column:status" json:"status"`
	Proof     string    `gorm:"column:proof" json:"proof"`
	ErrorMsg  string    `gorm:"column:error_msg" json:"error_msg"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" sql:"datetime"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" sql:"datetime"`
}

func NewProofWithHeight(chainID uint16, height string) *Proof {
	return &Proof{
		ChainID: chainID,
		Height:  height,
	}
}

func NewProofWithUniqueID(uniqueID string) *Proof {
	return &Proof{
		UniqueID: uniqueID,
	}
}

func (p *Proof) TableName() string {
	return "proof"
}

func (p *Proof) Create() error {
	return db.GetDB().Create(p).Error
}

func (p *Proof) Get() (proof *Proof, err error) {
	err = db.GetDB().Where(p).First(&proof).Error
	return proof, err
}

func (p *Proof) Updates(np *Proof) error {
	return db.GetDB().Where(p).Updates(np).Error
}

func (p *Proof) SetError(msg string) error {
	np := &Proof{
		ErrorMsg: msg,
		Status:   ProofStatusError,
	}
	return p.Updates(np)
}

func (p *Proof) SetCompleted(proof string) error {
	np := &Proof{
		Proof:  proof,
		Status: ProofStatusCompleted,
	}
	return p.Updates(np)
}

func (p *Proof) IsPending() bool {
	return p.Status == ProofStatusPending
}
func (p *Proof) IsError() bool {
	return p.Status == ProofStatusError
}
func (p *Proof) IsCompleted() bool {
	return p.Status == ProofStatusCompleted
}
