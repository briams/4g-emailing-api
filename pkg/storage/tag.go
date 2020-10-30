package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/briams/4g-emailing-api/db/mysql"
	"github.com/briams/4g-emailing-api/pkg/models/tag"
	"github.com/briams/4g-emailing-api/pkg/models/tagaudit"
	"github.com/briams/4g-emailing-api/pkg/utils"
)

var (
	tableTag       = "Models"
	mysqlCreateTag = fmt.Sprintf(`INSERT INTO %s
		(modelId, mjml, html, variables,
		insUserId, insDate, insDatetime, insTimestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tableTag)
	mysqlUpdateTag = fmt.Sprintf(`UPDATE %s SET
	mjml = ?, html = ?, variables = ? WHERE modelId = ?`, tableTag)
	mysqlGetAllTags = fmt.Sprintf(`SELECT
		modelId, mjml, html, variables, active,
		insUserId, insDate, insDatetime, insTimestamp
		FROM %s`, tableTag)
	mysqlGetTagByID   = mysqlGetAllTags + " WHERE modelId = ?"
	mysqlGetTagsByIDs = mysqlGetAllTags + " WHERE modelId IN "
)

// MySQLTag used for work with mySQL - para
type MySQLTag struct {
	db           *sql.DB
	storageAudit tagaudit.Storage
}

// NewMySQLTag return a new pointer of MySQLTag
func NewMySQLTag(db *sql.DB, a tagaudit.Storage) *MySQLTag {
	return &MySQLTag{
		db:           db,
		storageAudit: a,
	}
}

// Create implement the interface tag.Storage
func (t *MySQLTag) Create(m *tag.Model) error {
	stmt, err := t.db.Prepare(mysqlCreateTag)
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

// GetByID implement the interface tag.Storage
func (t *MySQLTag) GetByID(modelID string) (*tag.Model, error) {
	stmt, err := t.db.Prepare(mysqlGetTagByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return t.scanRow(stmt.QueryRow(modelID))
}

// GetByIDs implement the interface tag.Storage
func (t *MySQLTag) GetByIDs(modelIDs ...string) (tag.Models, error) {
	ms := make(tag.Models, 0)

	q := mysqlGetTagsByIDs + "(?" + strings.Repeat(",?", len(modelIDs)-1) + ")"
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

// GetAll implement the interface tag.Storage
func (t *MySQLTag) GetAll() (tag.Models, error) {
	ms := make(tag.Models, 0)

	stmt, err := t.db.Prepare(mysqlGetAllTags)
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

// GetAllWithFields implement the interface tag.Storage
func (t *MySQLTag) GetAllWithFields(fields ...string) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)

	fieldsQuery := strings.Join(fields, ", ")

	mysqlQuery := fmt.Sprintf("SELECT %s FROM %s", fieldsQuery, tableTag)

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

// Update implements the interface tag.Storage
func (t *MySQLTag) Update(m *tag.Model) error {
	prevTag, err := t.GetByID(m.ModelID)
	if err != nil {
		return err
	}

	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(mysqlUpdateTag)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Tag: %w", err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(m.Mjml, m.Html, m.Variables, m.ModelID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Tag: %w", err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return tx.Commit()
	}

	tagAudit := &tagaudit.Model{
		PrevTag:   *prevTag,
		Tag:       *m,
		SetUserID: m.InsUserID,
	}
	if err := t.storageAudit.CreateTx(tx, tagAudit); err != nil {
		tx.Rollback()
		return fmt.Errorf("Audit: %w", err)
	}

	return tx.Commit()
}

func (t *MySQLTag) scanRow(s mysql.RowScanner) (*tag.Model, error) {
	tag := &tag.Model{}

	variablesNull := sql.NullString{}

	if err := s.Scan(
		&tag.ModelID,
		&tag.Mjml,
		&tag.Html,
		&variablesNull,
		&tag.Active,
		&tag.InsUserID,
		&tag.InsDate,
		&tag.InsDateTime,
		&tag.InsTimestamp,
	); err != nil {
		return nil, err
	}

	tag.Variables = variablesNull.String
	return tag, nil
}
