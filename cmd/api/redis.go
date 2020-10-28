package main

import (
	"context"
	"log"

	"github.com/briams/4g-emailing-api/config"
	"github.com/briams/4g-emailing-api/db/rds"
	"github.com/briams/4g-emailing-api/pkg/storage"
)

var ctx = context.Background()

func newRedisClient() *rds.Rds {
	redisClient := config.GetRedisClient()
	rds := rds.NewRds(ctx, redisClient)

	res, err := rds.RdsPing()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis:", res)

	return rds
}

func refreshData(rdb *rds.Rds) {
	db := newConnection()
	storageParam := storage.NewRedisTag(db, rdb)
	if err := storageParam.CleanData(); err != nil {
		log.Fatal("Redis:", err)
	}

	if err := storageParam.RefreshData(); err != nil {
		log.Fatal("Redis:", err)
	}
}
