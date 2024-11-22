package controllers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var PgDBConn *gorm.DB

func PgConnectDB() {
	// Access DB credentials from environment
	host := os.Getenv("db_host")
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbname := os.Getenv("db_name")
	dbport := os.Getenv("db_port")

	// Construct a proper DSN string for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, dbport)

	fmt.Println("Starting connection with Postgres Db")

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connection successful.")

	// Set the global connection variable
	PgDBConn = db

	// Run database migrations
	err = db.AutoMigrate(&User{}, &TodoPG{}, &Item{})
	if err != nil {
		log.Fatalf("Data migration failed: %v", err)
	}

	fmt.Println("Data migration complete.")
}
