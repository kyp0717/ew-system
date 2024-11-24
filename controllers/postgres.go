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

	// Open the database connection with custom naming strategy
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Disable pluralization of table names
		},
	})

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connection successful.")

	// Set the global connection variable
	PgDBConn = db

	// List of models to migrate
	models := []interface{}{&User{}, &Item{}, &Product{}}

	// Iterate over models and run migrations
	for _, model := range models {

		if !db.Migrator().HasTable(model) {
			err = db.AutoMigrate(model)
			if err != nil {
				log.Fatalf("Data migration for %T failed: %v", model, err)
			}

			// Check the model type and load corresponding data
			switch model.(type) {
			case *User:
				LoadUserTable()
			case *Item:
				LoadItemTable()
			case *Product:
				LoadProductTable()
			default:
				fmt.Printf("No data loader defined for %T.\n", model)
			}

			fmt.Printf("Table for %T migrated and data loaded successfully.\n", model)
		} else {
			fmt.Printf("Table for %T already exists. Skipping migration and data load.\n", model)
		}
	}

	fmt.Println("Data migration complete.")
}
