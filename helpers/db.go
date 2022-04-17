package helpers

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

var dbConnOnce sync.Once
var db *sqlx.DB

func PrepareTables() {
	schema := `
		CREATE TABLE IF NOT EXISTS "tasmota_devices" (
	"mac_address" VARCHAR(250) NOT NULL,
	"device_ip" VARCHAR(100) NULL,
	"device_name" VARCHAR(250) NULL,
	"tasmota_version" VARCHAR(100) NULL,
	"last_check" DATETIME NULL,
	PRIMARY KEY ("mac_address")
);`
	GetDb().MustExec(schema)
}

func GetDb() *sqlx.DB {
	dbConnOnce.Do(func() {
		var err error
		db, err = sqlx.Connect("sqlite3", "app_data.db")
		if err != nil {
			GetLogger().Fatal(err.Error())
		}
	})
	return db
}
