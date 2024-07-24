package database

import (
	"database/sql"
	"log"
)

func OpenDb(driverName, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Databaseni ochishda xatolik - %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Databasega ulanish(Ping)da xatolik - %v", err)
	}

	return db
}
