package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Product struct {
	SKU         string    `json:"SKU" gorm:"primarykey;type:text"`
	Category    string    `json:"Category" gorm:"column:Category;type:text;null"`
	Group       string    `json:"Group" gorm:"column:Group;type:text;null"`
	ProductName string    `json:"ProductName" gorm:"column:ProductName;type:text;null"`
	CreateDate  time.Time `json:"Create_Date" gorm:"column:Create_Date;type:date;null"`
	TotalBoxes  int       `json:"Total_Boxes" gorm:"column:Total_Boxes;null"`

	SKU1   string `json:"SKU_1" gorm:"column:SKU_1;type:text;null"`
	Box1   int    `json:"BOX_1" gorm:"column:BOX_1;null"`
	Piece1 int    `json:"Piece_1" gorm:"column:Piece_1;null"`

	SKU2   string `json:"SKU_2" gorm:"column:SKU_2;type:text;null"`
	Box2   int    `json:"BOX_2" gorm:"column:BOX_2;null"`
	Piece2 int    `json:"Piece_2" gorm:"column:Piece_2;null"`

	SKU3   string `json:"SKU_3" gorm:"column:SKU_3;type:text;null"`
	Box3   int    `json:"BOX_3" gorm:"column:BOX_3;null"`
	Piece3 int    `json:"Piece_3" gorm:"column:Piece_3;null"`

	SKU4   string `json:"SKU_4" gorm:"column:SKU_4;type:text;null"`
	Box4   int    `json:"BOX_4" gorm:"column:BOX_4;null"`
	Piece4 int    `json:"Piece_4" gorm:"column:Piece_4;null"`

	SKU5   string `json:"SKU_5" gorm:"column:SKU_5;type:text;null"`
	Box5   int    `json:"BOX_5" gorm:"column:BOX_5;null"`
	Piece5 int    `json:"Piece_5" gorm:"column:Piece_5;null"`

	SKU6   string `json:"SKU_6" gorm:"column:SKU_6;type:text;null"`
	Box6   int    `json:"BOX_6" gorm:"column:BOX_6;null"`
	Piece6 int    `json:"Piece_6" gorm:"column:Piece_6;null"`

	SKU7   string `json:"SKU_7" gorm:"column:SKU_7;type:text;null"`
	Box7   int    `json:"BOX_7" gorm:"column:BOX_7;null"`
	Piece7 int    `json:"Piece_7" gorm:"column:Piece_7;null"`

	Active    string    `json:"Active" gorm:"column:Active;type:char(1);null"`
	CreatedBy string    `json:"CreatedBy" gorm:"column:CreatedBy;type:text;size:20;null"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName sets the default table name for the Product struct
func (Product) TableName() string {
	return "product" // Replace with your actual table name
}

// ProcessedProduct is used to pass data to the templ view
type ProcessedProduct struct {
	Values []string
}

func LoadProductTable() error {

	// Load CSV data
	csvFile, err := os.Open("./controllers/data/product.csv") // Replace with your CSV file path
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read all rows from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Define the date format
	const dateFormat = "2006-01-02" // Adjust this format if your date format differs

	// Loop through the records and create products
	for _, record := range records {
		// Parse integers
		parseInt := func(value string) int {
			parsed, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("Failed to convert to int: %v", err)
				return 0
			}
			return parsed
		}

		// Parse date
		var createDate time.Time
		if record[4] != "" {
			createDate, err = time.Parse(dateFormat, record[4])
			if err != nil {
				log.Printf("Failed to parse date - CreateDate: %v", err)
				continue
			}
		}

		// Create the product record
		product := &Product{
			SKU:         record[0],
			Category:    record[1],
			Group:       record[2],
			ProductName: record[3],
			CreateDate:  createDate,
			TotalBoxes:  parseInt(record[5]),
			SKU1:        record[6],
			Box1:        parseInt(record[7]),
			Piece1:      parseInt(record[8]),
			SKU2:        record[9],
			Box2:        parseInt(record[10]),
			Piece2:      parseInt(record[11]),
			SKU3:        record[12],
			Box3:        parseInt(record[13]),
			Piece3:      parseInt(record[14]),
			SKU4:        record[15],
			Box4:        parseInt(record[16]),
			Piece4:      parseInt(record[17]),
			SKU5:        record[18],
			Box5:        parseInt(record[19]),
			Piece5:      parseInt(record[20]),
			SKU6:        record[21],
			Box6:        parseInt(record[22]),
			Piece6:      parseInt(record[23]),
			SKU7:        record[24],
			Box7:        parseInt(record[25]),
			Piece7:      parseInt(record[26]),
			Active:      record[27],
			CreatedBy:   record[28],
			CreatedAt:   time.Now(),
		}

		// Save product to the database
		err = PgDBConn.Create(&product).Error
		if err != nil {
			log.Printf("Failed to insert product record: %v", err)
		}
	}

	fmt.Println("....Product Data imported successfully!")

	return nil
}

// GetProductBySKU fetches a Product from the database using its SKU.
func GetProductBySKU(sku string) (*Product, error) {
	var product Product
	result := PgDBConn.Where("SKU = ?", sku).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

// Helper function to retrieve SKU, Box, and Piece fields from a product by field number.
func (p *Product) GetFieldValue(field string, index int) interface{} {
	switch field {
	case "SKU":
		switch index {
		case 1:
			return p.SKU1
		case 2:
			return p.SKU2
		case 3:
			return p.SKU3
		case 4:
			return p.SKU4
		case 5:
			return p.SKU5
		case 6:
			return p.SKU6
		case 7:
			return p.SKU7
		}
	case "Box":
		switch index {
		case 1:
			return p.Box1
		case 2:
			return p.Box2
		case 3:
			return p.Box3
		case 4:
			return p.Box4
		case 5:
			return p.Box5
		case 6:
			return p.Box6
		case 7:
			return p.Box7
		}
	case "Piece":
		switch index {
		case 1:
			return p.Piece1
		case 2:
			return p.Piece2
		case 3:
			return p.Piece3
		case 4:
			return p.Piece4
		case 5:
			return p.Piece5
		case 6:
			return p.Piece6
		case 7:
			return p.Piece7
		}
	}
	return nil
}
