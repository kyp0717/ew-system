package controllers

import (
	"fmt"
	"log"
	"os"

	// "gorm.io/gorm/logger"
	//"github.com/emarifer/gofiber-templ-htmx/data/load"

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
	tablesToCheck := []string{"user", "category", "item", "product", "product_group", "product_detail"}

	for _, table := range tablesToCheck {
		if TableExists(PgDBConn, table) {
			fmt.Printf("Table %s exists.\n", table)
		} else {
			fmt.Printf("Table %s does not exist.\n", table)

			// Check which table does not exist and generate table respectively
			switch table {
			case "user":
				PgDBConn.AutoMigrate(&User{})
				err := controller.UserLoad(PgDBConn)
				if err != nil {
					log.Fatal("AutoMigrate failed to migrate login: ", err)
				}

			case "category":
				PgDBConn.AutoMigrate(&Category{})
				err := controllers.Category(PgDBConn)
				if err != nil {
					log.Fatal("AutoMigrate failed to migrate category: ", err)
				}
			case "item":
				PgDBConn.AutoMigrate(&Item{})
				err := controllers.Item(PgDBConn)
				if err != nil {
					log.Fatal("AutoMigrate failed to migrate item: ", err)
				}
				/*
					case "product":
						db.AutoMigrate(	&Product{}, )
						controllers..Product(PgDBConn)
					case "product_group":
						db.AutoMigrate(	&ProductGroup{},	)
						controllers..ProductGroup(PgDBConn)
					case "product_detail":
						db.AutoMigrate(	&ProductDetail{}, )
						controllers..ProductDetail(PgDBConn)
				*/
			}
		}
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