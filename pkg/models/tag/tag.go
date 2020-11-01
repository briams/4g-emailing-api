package tag

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/briams/4g-emailing-api/pkg/mjmlparser"
	"github.com/briams/4g-emailing-api/pkg/utils"
	"github.com/go-redis/redis/v8"
)

var (
	// ErrIDNotFound error for IDNotFound
	ErrIDNotFound = errors.New("The model does not have an ID")

	// ErrModelAlreadyExist error for IDNotFound
	ErrModelAlreadyExist = errors.New("The model already exist")

	// ErrFieldsDoesNotExist error for FieldsDoesNotExist
	ErrFieldsDoesNotExist = fmt.Errorf("Only these existing fields can be accepted: %v", ExistingFields)
)

// ExistingFields for events
var ExistingFields = []string{
	"modelId", "mjml", "html", "variables", "active", "insUserId",
	"insDate", "insDatetime", "insTimestamp",
}

// Model of Param
type Model struct {
	ModelID      string     `json:"modelId,omitempty"`
	Mjml         string     `json:"mjml,omitempty"`
	Html         string     `json:"html,omitempty"`
	Variables    string     `json:"variables"`
	Active       uint8      `json:"active,omitempty"`
	InsUserID    uint       `json:"insUserId,omitempty"`
	InsDate      *time.Time `json:"insDate,omitempty"`
	InsDateTime  *time.Time `json:"insDatetime,omitempty"`
	InsTimestamp int64      `json:"insTimestamp,omitempty"`
}

// Models slice of Model
type Models []*Model

// Storage interface that must implement a db storage
type Storage interface {
	Create(*Model) error
	Update(*Model) error
	GetByID(string) (*Model, error)
	GetByIDs(...string) (Models, error)
	GetAll() (Models, error)
	GetAllWithFields(...string) ([]map[string]interface{}, error)
	Activate(string, string, uint) error
	Deactivate(string, string, uint) error
}

// Service of model
type Service struct {
	storage Storage
}

// NewService return a pointer of Service
func NewService(s Storage) *Service {
	return &Service{s}
}

// Create is used for create a model
func (s *Service) Create(m *Model) error {
	currentModel, err := s.GetByID(m.ModelID)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, redis.Nil) {
		m.Html = m.Mjml
		return s.storage.Create(m)
	}
	if currentModel.ModelID == m.ModelID {
		return ErrModelAlreadyExist
	}

	return err
}

// GetAll is used for get all the models
func (s *Service) GetAll() (Models, error) {
	return s.storage.GetAll()
}

// GetAllWithFields is used for get all the models
func (s *Service) GetAllWithFields(fields ...string) ([]map[string]interface{}, error) {
	for _, field := range fields {
		if !utils.Include(ExistingFields, field) {
			return nil, ErrFieldsDoesNotExist
		}
	}

	return s.storage.GetAllWithFields(fields...)
}

// GetByID is used for get a model
func (s *Service) GetByID(id string) (*Model, error) {
	return s.storage.GetByID(id)
}

// GetByIDs is used for get models by its IDs
func (s *Service) GetByIDs(ModelIDs ...string) (Models, error) {
	return s.storage.GetByIDs(ModelIDs...)
}

// Update is used for update a model
func (s *Service) Update(m *Model) error {
	if m.ModelID == "" {
		return ErrIDNotFound
	}

	mjml, err := mjmlparser.GenerateMJMLWithData(m.Mjml, map[string]interface{}{})

	if err != nil {
		log.Fatal(err)
	}

	res := mjmlparser.ParserMJMLtoHTML(mjml)

	m.Html = res.MJMLReplyResponse.HTML
	// m.Mjml
	// fmt.Println("HTML: ", res.MJMLReplyResponse.HTML)

	return s.storage.Update(m)
}

// Activate is used for active a service
func (s *Service) Activate(id string, reason string, setUserID uint) error {
	return s.storage.Activate(id, reason, setUserID)
}

// Deactivate is used for deactive a service
func (s *Service) Deactivate(id string, reason string, setUserID uint) error {
	return s.storage.Deactivate(id, reason, setUserID)
}
