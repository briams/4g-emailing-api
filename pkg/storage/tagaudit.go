package storage

import (
	"database/sql"
	"fmt"

	"git.gdpteam.com/4gen/4g-tags-api/db/mysql"
	"git.gdpteam.com/4gen/4g-tags-api/pkg/models/tagaudit"
	"git.gdpteam.com/4gen/4g-tags-api/pkg/utils"
)

var (
	tableTagAudit       = "TagsAudit"
	mysqlCreateTagAudit = fmt.Sprintf(`
		INSERT INTO %s (
			modelId, mjmlPrev, htmlPrev, variablesPrev,
			mjml, html, variables, setUserId,
			setDate, setDatetime, setTimestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableTagAudit)
)

// MySQLTagAudit used for work with mySQL - paramaudit
type MySQLTagAudit struct {
	db *sql.DB
}

// NewMySQLTagAudit return a new pointer of MySQLTagAudit
func NewMySQLTagAudit(db *sql.DB) *MySQLTagAudit {
	return &MySQLTagAudit{db}
}

// CreateTx registra en la BD
func (m *MySQLTagAudit) CreateTx(tx *sql.Tx, tagAudit *tagaudit.Model) error {
	stmt, err := tx.Prepare(mysqlCreateTagAudit)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := utils.Now()
	return mysql.ExecAffectingOneRow(
		stmt,
		tagAudit.PrevTag.ModelID,
		tagAudit.PrevTag.Mjml,
		tagAudit.PrevTag.Html,
		tagAudit.PrevTag.Variables,
		tagAudit.Tag.Mjml,
		tagAudit.Tag.Html,
		tagAudit.Tag.Variables,
		tagAudit.SetUserID,
		now["date"], now["time"], now["unix"],
	)
}
