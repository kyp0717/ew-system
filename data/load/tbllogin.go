package load

import (
	"github.com/kyp0717/ewxback/webapp/model"
  	"os"
  	"log"
  	"fmt"
  	"encoding/csv"
	"gorm.io/gorm"
	"strconv"
) 

func Login(db *gorm.DB ) {
	
	// Load CSV data
	csvFile, err := os.Open("data/csv/data_login.csv") // Replace with your CSV file path
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 			// Allow variable number of fields per record

	// Read all rows from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV test file: %v", err)
	}

	// Loop through the records and create Category
	for _, record := range records{

		iID, err	:= strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails			
		}
		iPermission, err	:= strconv.Atoi(record[5])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails
		}

		test := &model.Login{
			ID:         			iID,	
			email:   				record[1],
			UserName:         		record[3],	
			Password:   			record[4],
			Permission:   			iPermission,
		}

		// Save item to the database
		err = db.Create(&test).Error
		if err != nil {
			log.Printf("Failed to insert test record: %v", err)
		}
	}

	fmt.Println("....Login Data imported successfully!")
}
