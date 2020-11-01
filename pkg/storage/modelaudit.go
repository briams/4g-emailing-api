package storage

import (
	"database/sql"
	"fmt"

	"github.com/briams/4g-emailing-api/db/mysql"
	"github.com/briams/4g-emailing-api/pkg/models/modelaudit"
	"github.com/briams/4g-emailing-api/pkg/utils"
)

var (
	tableModelAudit       = "ModelsAudit"
	mysqlCreateModelAudit = fmt.Sprintf(`
		INSERT INTO %s (
			modelId, mjmlPrev, htmlPrev, variablesPrev,
			mjml, html, variables, setUserId,
			setDate, setDatetime, setTimestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableModelAudit)
)

// MySQLModelAudit used for work with mySQL - modelaudit
type MySQLModelAudit struct {
	db *sql.DB
}

// NewMySQLModelAudit return a new pointer of MySQLModelAudit
func NewMySQLModelAudit(db *sql.DB) *MySQLModelAudit {
	return &MySQLModelAudit{db}
}

// CreateTx registra en la BD
func (m *MySQLModelAudit) CreateTx(tx *sql.Tx, modelAudit *modelaudit.Model) error {
	stmt, err := tx.Prepare(mysqlCreateModelAudit)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := utils.Now()
	return mysql.ExecAffectingOneRow(
		stmt,
		modelAudit.PrevModel.ModelID,
		modelAudit.PrevModel.Mjml,
		modelAudit.PrevModel.Html,
		modelAudit.PrevModel.Variables,
		modelAudit.Model.Mjml,
		modelAudit.Model.Html,
		modelAudit.Model.Variables,
		modelAudit.SetUserID,
		now["date"], now["time"], now["unix"],
	)
}
