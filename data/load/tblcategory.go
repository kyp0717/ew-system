package load

import 
(
	"github.com/kyp0717/ewxback/webapp/model"
	"github.com/shopspring/decimal"
//	"github.com/kyp0717/ewxback/data"
//	"github.com/a-h/templ@v0.2.786"
  	"os"
  	"log"
  	"fmt"
  	"encoding/csv"
	"gorm.io/gorm"
  	"strconv"
) 

func Category(db *gorm.DB ) {

	// Load CSV data
	csvFile, err := os.Open("data/csv/data_category.csv") // Replace with your CSV file path
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
    	p, _ := decimal.NewFromString(record[2])
		r, err := strconv.Atoi(record[3])
		if err != nil {
			log.Printf("Failed to convert record[3] to int: %v", err)
			continue // Skip this record if conversion fails
		}

		test := &model.Category{
			Category:         		record[0],	
			Category_Description:   record[1],
			Reserved_Percentage_BySale:     p,
			Report_Out_Max:					r,
		}

		// Save item to the database
		err = db.Create(&test).Error
		if err != nil {
			log.Printf("Failed to insert test record: %v", err)
		}
	}

	fmt.Println("....Category Data imported successfully!")
}
// Helper function to parse integers
func ParseInt(value string) int {
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
func ParseFloat(value string) float64 {
	if value == "" {
		return 0.0
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return v
}
