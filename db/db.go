package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection is not available: %v", err)
	}
	log.Println("Database connection successful.")
	createTables()
}

func createTables() {
	createPaymentsTable := `
    CREATE TABLE IF NOT EXISTS payments (
        id INT PRIMARY KEY AUTO_INCREMENT,
        amount INT NOT NULL,
        description TEXT,
        authority VARCHAR(100) NOT NULL UNIQUE,
        ref_id VARCHAR(100),
        status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := DB.Exec(createPaymentsTable)
	if err != nil {
		log.Fatalf("Failed to create payments table: %v", err)
	}

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users(
		id INT PRIMARY KEY AUTO_INCREMENT,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		UNIQUE(email)
		);`

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("failed to create payments table: %v", err)
	}
	log.Println("Payments table checked/created successfully.")
}
