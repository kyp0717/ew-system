package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/kyp0717/ewxback/data/load"
	"github.com/kyp0717/ewxback/webapp/model"

	// "gorm.io/gorm/logger"

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

	fmt.Println("Starting connection with Postgres Db")
	dsn := user + "://postgres:" + password + "@" + host + ":" + dbport + "/" + dbname + "?sslmode=disable"

	//db, err := gorm.Open(postgres.Open(dsn) , &gorm.Config{})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	if err != nil {
		panic("Database connection failed.")
	}

	log.Println("Connection successful.")

	PgDBConn = db

	// Check if tables exist
	tablesToCheck := []string{"login", "category", "item", "product", "product_group", "product_detail"}

	for _, table := range tablesToCheck {
		if TableExists(db, table) {
			fmt.Printf("Table %s exists.\n", table)
		} else {
			fmt.Printf("Table %s does not exist.\n", table)

			// Check which table does not exist and generate table respectively
			switch table {
			case "login":
				db.AutoMigrate(&model.Login{})
				load.Login(db)
			case "category":
				db.AutoMigrate(&model.Category{})
				load.Category(db)
			case "item":
				db.AutoMigrate(&model.Item{})
				load.Item(db)

				/*
					case "product":
						db.AutoMigrate(	&model.Product{}, )
						load.Product(db)
					case "product_group":
						db.AutoMigrate(	&model.ProductGroup{},	)
						load.ProductGroup(db)
					case "product_detail":
						db.AutoMigrate(	&model.ProductDetail{}, )
						load.ProductDetail(db)
				*/
			}
		}
	}

	// Migrate your models
	err = db.AutoMigrate(
		&model.Login{},
		&model.Item{},
		&model.Product{},
		&model.ProductGroup{},
		&model.ProductDetail{},
		&model.Todo{},
	)
	if err != nil {
		log.Fatal("failed to migrate models: ", err)
	}

	fmt.Println("Data Migration complete.")
}

func TableExists(db *gorm.DB, tableName string) bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s')", tableName)
	err := db.Raw(query).Scan(&exists).Error
	if err != nil {
		log.Printf("Error checking if table %s exists: %v", tableName, err)
		return false
	}
	return exists
}
