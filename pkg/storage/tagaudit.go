package storage

import (
	"database/sql"
	"fmt"

	"github.com/briams/4g-emailing-api/db/mysql"
	"github.com/briams/4g-emailing-api/pkg/models/tagaudit"
	"github.com/briams/4g-emailing-api/pkg/utils"
)

var (
	tableTagAudit       = "ModelsAudit"
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
