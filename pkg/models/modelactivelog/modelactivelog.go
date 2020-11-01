package modelactivelog

import (
	"database/sql"
	"errors"
	"time"
)

var (
	// ErrIDNotFound error for IDNotFound
	ErrIDNotFound = errors.New("The log does not have an ID")
)

// Model of ServiceActiveLog
type Model struct {
	LogID        uint      `json:"logId"`
	ModelID      string    `json:"modelId"`
	Active       uint8     `json:"active"`
	Reason       string    `json:"reason"`
	SetUserID    uint      `json:"setUserId"`
	SetDate      time.Time `json:"setDate,omitempty"`
	SetDateTime  time.Time `json:"setDateTime,omitempty"`
	SetTimestamp int64     `json:"setTimestamp,omitempty"`
}

// Models slice of Model
type Models []*Model

// Storage interface that must implement a db storage
type Storage interface {
	CreateTx(*sql.Tx, *Model) error
}

// ModelActiveLog of Model
type ModelActiveLog struct {
	storage Storage
}

// NewModelActiveLog return a pointer of Model
func NewModelActiveLog(s Storage) *ModelActiveLog {
	return &ModelActiveLog{s}
}
