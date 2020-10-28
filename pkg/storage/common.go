package storage

import (
	"database/sql"
	"time"
)

var commonSelectNow = "SELECT NOW() AS now"

// MySQLCommon used for work with mySQL - para
type MySQLCommon struct {
	db *sql.DB
}

// NewMySQLCommon return a new pointer of MySQLCommon
func NewMySQLCommon(db *sql.DB) *MySQLCommon {
	return &MySQLCommon{db}
}

// CheckDBHealth chehks the DB health
func (p *MySQLCommon) CheckDBHealth() (*time.Time, error) {
	now := &time.Time{}
	stmt, err := p.db.Prepare(commonSelectNow)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(now)
	if err != nil {
		return nil, err
	}

	return now, nil
}
