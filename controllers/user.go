package controllers

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	//"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Permission string `json:"permission"`
}

// ****************  Table Load ***********************
func (t *User) UserLoad(PgDBConn) error {

	// Load CSV data
	csvFile, err := os.Open("data/csv/data_login.csv") // Replace with your CSV file path
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read all rows from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV test file: %v", err)
		return err
	}

	// Loop through the records and create Category
	for _, record := range records {

		iID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails
		}
		iPermission, err := strconv.Atoi(record[5])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails
		}

		test := &User{
			ID:         iID,
			email:      record[1],
			UserName:   record[3],
			Password:   record[4],
			Permission: iPermission,
		}

		// Save item to the database
		err = PgDBConn.Create(&test).Error
		if err != nil {
			log.Printf("Failed to insert test record: %v", err)
		}
	}
}
