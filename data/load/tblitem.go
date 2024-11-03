package load

import 
(
	"github.com/kyp0717/ewxback/webapp/model"
	"github.com/shopspring/decimal"
//	"github.com/a-h/templ@v0.2.786"
//	"github.com/kyp0717/ewxback/data"
  	"os"
  	"log"
  	"fmt"
  	"encoding/csv"
	"gorm.io/gorm"
  	"strconv"
) 

func Item(db *gorm.DB ) {

	// Load CSV data
	csvFile, err := os.Open("data/csv/data_item.csv") // Replace with your CSV file path
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

		iInventory, err	:= strconv.Atoi(record[6])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails
		}		
		iQtyPerBox, err 	:= strconv.Atoi(record[7])
		if err != nil {
			log.Printf("Failed to convert to int - QtyPerBox: %v", err)
			continue // Skip this record if conversion fails
		}		
		iLength, err 		:= strconv.Atoi(record[16])
		if err != nil {
			log.Printf("Failed to convert to int - Length: %v", err)
			continue // Skip this record if conversion fails
		}		
		iWidth, err 		:= strconv.Atoi(record[17])
		if err != nil {
			log.Printf("Failed to convert to int - Width: %v", err)
			continue // Skip this record if conversion fails
		}
		iHeight, err 		:= strconv.Atoi(record[18])
		if err != nil {
			log.Printf("Failed to convert to int - Height: %v", err)
			continue // Skip this record if conversion fails
		}		
		iBoxLength, err 	:= strconv.Atoi(record[20])
		if err != nil {
			log.Printf("Failed to convert to int - BoxLength: %v", err)
			continue // Skip this record if conversion fails
		}		
		iBoxWidth, err 		:= strconv.Atoi(record[21])
		if err != nil {
			log.Printf("Failed to convert to int - BoxWidth: %v", err)
			continue // Skip this record if conversion fails
		}		
		iBoxHeight, err 	:= strconv.Atoi(record[22])
		if err != nil {
			log.Printf("Failed to convert to int - BoxHeight: %v", err)
			continue // Skip this record if conversion fails
		}		
		iBoxWeight, err 	:= strconv.Atoi(record[23])
		if err != nil {
			log.Printf("Failed to convert to int - BoxWeight: %v", err)
			continue // Skip this record if conversion fails
		}		
		iPiecesContainer, err 	:= strconv.Atoi(record[26])
		if err != nil {
			log.Printf("Failed to convert to int - PiecesContainer: %v", err)
			continue // Skip this record if conversion fails
		}
    	dCost, err 	:= decimal.NewFromString(record[8])	
		if err != nil {
			log.Printf("Failed to convert to decimal - Cost: %v", err)
			continue // Skip this record if conversion fails
		}	
    	dShippingCost, err 	:= decimal.NewFromString(record[28])	
		if err != nil {
			log.Printf("Failed to convert to decimal - dShippingCost: %v", err)
			continue // Skip this record if conversion fails
		}		
    	dPrice, err 	:= decimal.NewFromString(record[9])					
		if err != nil {
			log.Printf("Failed to convert to decimal - Price: %v", err)
			continue // Skip this record if conversion fails
		}
    	dPrice5, err 	:= decimal.NewFromString(record[10])					
		if err != nil {
			log.Printf("Failed to convert to decimal - Price5: %v", err)
			continue // Skip this record if conversion fails
		}
    	dPrice7, err 	:= decimal.NewFromString(record[11])					
		if err != nil {
			log.Printf("Failed to convert to decimal - Price7: %v", err)
			continue // Skip this record if conversion fails
		}
		dPrice10, err 	:= decimal.NewFromString(record[12])					
		if err != nil {
			log.Printf("Failed to convert to decimal - Price10: %v", err)
			continue // Skip this record if conversion fails
		}
		dPrice15, err 	:= decimal.NewFromString(record[13])					
		if err != nil {
			log.Printf("Failed to convert to decimal - Price15: %v", err)
			continue // Skip this record if conversion fails
		}
		dPrice19, err 	:= decimal.NewFromString(record[14])					
		if err != nil {
			log.Printf("Failed to convert to decimal - Price19: %v", err)
			continue // Skip this record if conversion fails
		}


		itemrecord := &model.Item{
			SKU:         			record[0],
			ItemName:         		record[1],	
			UPC:         			record[2],
			Type:         			record[3],
			Category:         		record[4],
			Description:         	record[5],
			Inventory:  			iInventory,
			QtyPerBox:				iQtyPerBox,
			Cost: 					dCost,
			Price:					dPrice,
			Price5:					dPrice5,
			Price7:					dPrice7,
			Price10:				dPrice10,
			Price15:				dPrice15,
			Price19:				dPrice19,
			ItemDimension:  	    record[15],
			Length: 				iLength,
			Width: 					iWidth,
			Height: 				iHeight,
			BoxDimension:         	record[19],	
			BoxLength: 				iBoxLength,
			BoxWidth: 				iBoxWidth,
			BoxHeight: 				iBoxHeight, 
			BoxWeight:				iBoxWeight,
			AvailableDate:         	record[24],	
			ShippingMethod:         record[25],
			PiecesContainer: 		iPiecesContainer,
			Supplier:         		record[27],	
			ShippingCost: 			dShippingCost,
			Active: 				record[29],	
			UserName:         		record[30],	
			UpdateStamp:         	record[31],	
		}

		// Save item to the database
		err = db.Create(&itemrecord).Error
		if err != nil {
			log.Printf("Failed to insert test record: %v", err)
		}
	}

	fmt.Println("....Item Data imported successfully!")
}
