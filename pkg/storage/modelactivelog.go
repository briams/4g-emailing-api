package storage

import (
	"database/sql"
	"fmt"

	"github.com/briams/4g-emailing-api/db/mysql"
	"github.com/briams/4g-emailing-api/pkg/models/modelactivelog"
	"github.com/briams/4g-emailing-api/pkg/utils"
)

var (
	tableModelActiveLog       = "ModelsActiveLog"
	mysqlCreateModelActiveLog = fmt.Sprintf(`
		INSERT INTO %s (
			modelId, active, reason, setUserId, setDate, setDatetime, setTimestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, tableModelActiveLog)
)

// MySQLModelActiveLog used for work with mySQL
type MySQLModelActiveLog struct {
	db *sql.DB
}

// NewMySQLModelActiveLog return a new pointer of MySQLModelActiveLog
func NewMySQLModelActiveLog(db *sql.DB) *MySQLModelActiveLog {
	return &MySQLModelActiveLog{db}
}

// CreateTx registra en la BD
func (m *MySQLModelActiveLog) CreateTx(tx *sql.Tx, modelActiveLog *modelactivelog.Model) error {
	stmt, err := tx.Prepare(mysqlCreateModelActiveLog)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := utils.Now()

	return mysql.ExecAffectingOneRow(
		stmt,
		modelActiveLog.ModelID,
		modelActiveLog.Active,
		modelActiveLog.Reason,
		modelActiveLog.SetUserID,
		now["date"], now["time"], now["unix"],
	)
}
