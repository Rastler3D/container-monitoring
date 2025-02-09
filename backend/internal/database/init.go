package database

import (
	"database/sql"
	"log"
)

func InitDatabase(db *sql.DB) error {
	log.Println("Initializing database...")

	// Create container_status table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS containermonitor (
			id SERIAL PRIMARY KEY,
			ip_address VARCHAR(15) UNIQUE NOT NULL,
			ping_time FLOAT NOT NULL,
			last_ping TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
