package tagaudit

import (
	"database/sql"
	"errors"
	"time"

	"github.com/briams/4g-emailing-api/pkg/models/tag"
)

var (
	// ErrIDNotFound error for IDNotFound
	ErrIDNotFound = errors.New("The tag does not have an ID")
)

// Model of Param
type Model struct {
	AuditID      uint      `json:"auditId"`
	PrevTag      tag.Model `json:"prevTag"`
	Tag          tag.Model `json:"tag"`
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

// Service of tag
type Service struct {
	storage Storage
}

// NewService return a pointer of Service
func NewService(s Storage) *Service {
	return &Service{s}
}
