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
		log.Fatalf("failed to ping db: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("db connected")
}
