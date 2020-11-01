package modelaudit

import (
	"database/sql"
	"errors"
	"time"

	"github.com/briams/4g-emailing-api/pkg/models/model"
)

var (
	// ErrIDNotFound error for IDNotFound
	ErrIDNotFound = errors.New("The model does not have an ID")
)

// Model of Model
type Model struct {
	AuditID      uint        `json:"auditId"`
	PrevModel    model.Model `json:"prevModel"`
	Model        model.Model `json:"model"`
	SetUserID    uint        `json:"setUserId"`
	SetDate      time.Time   `json:"setDate,omitempty"`
	SetDateTime  time.Time   `json:"setDateTime,omitempty"`
	SetTimestamp int64       `json:"setTimestamp,omitempty"`
}

// Models slice of Model
type Models []*Model

// Storage interface that must implement a db storage
type Storage interface {
	CreateTx(*sql.Tx, *Model) error
}

// Service of model
type Service struct {
	storage Storage
}

// NewService return a pointer of Service
func NewService(s Storage) *Service {
	return &Service{s}
}
