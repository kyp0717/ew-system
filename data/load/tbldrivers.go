package load

import (
	"fmt"
	"log"
	"os"

	// "gorm.io/gorm/logger"

	"github.com/emarifer/gofiber-templ-htmx/data/load"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


func DB_LoadData(db) {

	tablesToCheck := []string{"login", "category", "item", "product", "product_group", "product_detail"}

	for _, table := range tablesToCheck {
		if TableExists(db, table) {
			fmt.Printf("Table %s exists.\n", table)
		} else {
			fmt.Printf("Table %s does not exist.\n", table)

			// Check which table does not exist and generate table respectively
			switch table {
			case "login":
				db.AutoMigrate(&load.Login{})
				load.Login(db)
				if err != nil {
					log.Fatal("AutoMigrate failed to migrate login: ", err)
				}
			
			case "category":
				db.AutoMigrate(&load.Category{})
				load.Category(db)
				if err != nil {
					log.Fatal("AutoMigrate failed to migrate category: ", err)
				}
			case "item":
				db.AutoMigrate(&load.Item{})
				load.Item(db)
				if err != nil {
					log.Fatal("AutoMigrate failed to migrate item: ", err)
				}
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
