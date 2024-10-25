package test

import (
	"github.com/kyp0717/ewxback/webapp/model"
  "github.com/shopspring/decimal"
  "os"
  "log"
  "fmt"
  "encoding/csv"
	"gorm.io/gorm"
  "strconv"
) 

func Testload1(db *gorm.DB ) {

	// Load CSV data
	csvFile, err := os.Open("test/test.csv") // Replace with your CSV file path
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read all rows from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV test file: %v", err)
	}

	// Loop through the records and create Items
	for _, record := range records {
    p, _ := decimal.NewFromString(record[1])
		test := &model.Test1{
			Item:         record[0],
			Price:        p,
		}

		// Save item to the database
		err = db.Create(&test).Error
		if err != nil {
			log.Printf("Failed to insert test record: %v", err)
		}
	}

	fmt.Println("Test Data imported successfully!")
}
// Helper function to parse integers
func parseInt(value string) int {
	if value == "" {
		return 0
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v
}

// Helper function to parse floats
func parseFloat(value string) float64 {
	if value == "" {
		return 0.0
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return v
}
