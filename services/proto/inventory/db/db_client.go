package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func NewInventoryDbConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("INVENTORY_DB_CONN"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
