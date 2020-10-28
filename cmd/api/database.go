package main

import (
	"database/sql"
	"log"

	"github.com/briams/4g-emailing-api/config"
)

func newConnection() *sql.DB {
	db, err := config.GetDBInstance()
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return db
}
