package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// Global DB connection variable
var Con *sql.DB

// InitConnection initializes the MySQL DB connection
func InitConnection() error {
	// Use a consistent and secure way to store credentials in production (env vars!)
	connectionString := "root:123456@tcp(127.0.0.1:3306)/goproject"

	// Open a new connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return fmt.Errorf("unable to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}

	// Set global variable only after success
	Con = db

	fmt.Println("✅ Connected to MySQL successfully!")
	return nil
}
