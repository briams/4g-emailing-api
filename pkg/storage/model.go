package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/briams/4g-emailing-api/db/mysql"
	"github.com/briams/4g-emailing-api/pkg/models/model"
	"github.com/briams/4g-emailing-api/pkg/models/modelactivelog"
	"github.com/briams/4g-emailing-api/pkg/models/modelaudit"
	"github.com/briams/4g-emailing-api/pkg/utils"
)

var (
	tableModel       = "Models"
	mysqlCreateModel = fmt.Sprintf(`INSERT INTO %s
		(modelId, mjml, html, variables,
		insUserId, insDate, insDatetime, insTimestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tableModel)
	mysqlUpdateModel = fmt.Sprintf(`UPDATE %s SET
	mjml = ?, html = ?, variables = ? WHERE modelId = ?`, tableModel)
	mysqlGetAllModels = fmt.Sprintf(`SELECT
		modelId, mjml, html, variables, active,
		insUserId, insDate, insDatetime, insTimestamp
		FROM %s`, tableModel)
	mysqlGetModelByID   = mysqlGetAllModels + " WHERE modelId = ?"
	mysqlGetModelsByIDs = mysqlGetAllModels + " WHERE modelId IN "
	mysqlUpdateActive   = fmt.Sprintf(`UPDATE %s SET
		active = ? WHERE modelId = ?`, tableModel)
)

// MySQLModel used for work with mySQL - para
type MySQLModel struct {
	db               *sql.DB
	storageAudit     modelaudit.Storage
	storageActiveLog modelactivelog.Storage
}

// NewMySQLModel return a new pointer of MySQLModel
func NewMySQLModel(db *sql.DB, a modelaudit.Storage, e modelactivelog.Storage) *MySQLModel {
	return &MySQLModel{
		db:               db,
		storageAudit:     a,
		storageActiveLog: e,
	}
}

// Create implement the interface model.Storage
func (t *MySQLModel) Create(m *model.Model) error {
	stmt, err := t.db.Prepare(mysqlCreateModel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := utils.Now()
	err = mysql.ExecAffectingOneRow(
		stmt,
		m.ModelID, m.Mjml, m.Html, m.Variables,
		m.InsUserID, now["date"], now["time"], now["unix"],
	)
	if err != nil {
		return err
	}
	// m.ModelID = m.ModelID

	return nil
}

// GetByID implement the interface model.Storage
func (t *MySQLModel) GetByID(modelID string) (*model.Model, error) {
	stmt, err := t.db.Prepare(mysqlGetModelByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return t.scanRow(stmt.QueryRow(modelID))
}

// GetByIDs implement the interface model.Storage
func (t *MySQLModel) GetByIDs(modelIDs ...string) (model.Models, error) {
	ms := make(model.Models, 0)

	q := mysqlGetModelsByIDs + "(?" + strings.Repeat(",?", len(modelIDs)-1) + ")"
	stmt, err := t.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := make([]interface{}, len(modelIDs))
	for i, ID := range modelIDs {
		args[i] = ID
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m, err := t.scanRow(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

// GetAll implement the interface model.Storage
func (t *MySQLModel) GetAll() (model.Models, error) {
	ms := make(model.Models, 0)

	stmt, err := t.db.Prepare(mysqlGetAllModels)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m, err := t.scanRow(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

// GetAllWithFields implement the interface model.Storage
func (t *MySQLModel) GetAllWithFields(fields ...string) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)

	fieldsQuery := strings.Join(fields, ", ")

	mysqlQuery := fmt.Sprintf("SELECT %s FROM %s", fieldsQuery, tableModel)

	stmt, err := t.db.Prepare(mysqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		data := make(map[string]interface{})
		for i, colName := range cols {
			data[colName] = columns[i]
		}

		values = append(values, data)
	}

	return values, nil
}

// Update implements the interface model.Storage
func (t *MySQLModel) Update(m *model.Model) error {
	prevModel, err := t.GetByID(m.ModelID)
	if err != nil {
		return err
	}

	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(mysqlUpdateModel)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Model: %w", err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(m.Mjml, m.Html, m.Variables, m.ModelID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Model: %w", err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return tx.Commit()
	}

	modelAudit := &modelaudit.Model{
		PrevModel: *prevModel,
		Model:     *m,
		SetUserID: m.InsUserID,
	}
	if err := t.storageAudit.CreateTx(tx, modelAudit); err != nil {
		tx.Rollback()
		return fmt.Errorf("Audit: %w", err)
	}

	return tx.Commit()
}

// Activate implements the interface service.Storage
func (t *MySQLModel) Activate(modelID string, reason string, setUserID uint) error {
	active := 1
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(mysqlUpdateActive)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Model: %w", err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(active, modelID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Model: %w", err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("Mysql: Could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return tx.Commit()
	}

	modelactivelog := &modelactivelog.Model{
		ModelID:   modelID,
		Reason:    reason,
		Active:    uint8(active),
		SetUserID: setUserID,
	}
	if err := t.storageActiveLog.CreateTx(tx, modelactivelog); err != nil {
		tx.Rollback()
		return fmt.Errorf("Active Log: %w", err)
	}

	return tx.Commit()
}

// Deactivate implements the interface service.Storage
func (t *MySQLModel) Deactivate(modelID string, reason string, setUserID uint) error {
	active := 0
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(mysqlUpdateActive)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Model: %w", err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(active, modelID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Model: %w", err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("Mysql: Could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return tx.Commit()
	}

	modelactivelog := &modelactivelog.Model{
		ModelID:   modelID,
		Reason:    reason,
		Active:    uint8(active),
		SetUserID: setUserID,
	}
	if err := t.storageActiveLog.CreateTx(tx, modelactivelog); err != nil {
		tx.Rollback()
		return fmt.Errorf("Active Log: %w", err)
	}

	return tx.Commit()
}

func (t *MySQLModel) scanRow(s mysql.RowScanner) (*model.Model, error) {
	model := &model.Model{}

	variablesNull := sql.NullString{}

	if err := s.Scan(
		&model.ModelID,
		&model.Mjml,
		&model.Html,
		&variablesNull,
		&model.Active,
		&model.InsUserID,
		&model.InsDate,
		&model.InsDateTime,
		&model.InsTimestamp,
	); err != nil {
		return nil, err
	}

	model.Variables = variablesNull.String
	return model, nil
}
