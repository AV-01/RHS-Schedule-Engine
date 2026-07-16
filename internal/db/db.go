package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := os.Getenv("SUPABASE_DB_URL")

	if connStr == "" {
		log.Fatal("supabase db url not found!")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed ot open database connection: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed ot ping db: %v", err)
	}
}
