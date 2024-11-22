package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Item struct {
	//	SKU             string          `json:"SKU" gorm:"primarykey;type:text;size:20"`
	SKU             string          `json:"SKU" gorm:"column:SKU;type:text;size:20;null"`
	ItemName        string          `json:"ItemName" gorm:"column:ItemName;type:text;null"`
	UPC             string          `json:"UPC" gorm:"column:UPC;type:text;null"`
	Type            string          `json:"Type" gorm:"column:Type;type:text;null"`
	Category        string          `json:"Category" gorm:"column:Category;type:text;null"`
	Description     string          `json:"Description" gorm:"column:Description;type:text;null"`
	Inventory       int             `json:"Inventory" gorm:"column:Inventory;null"`
	QtyPerBox       int             `json:"QtyPerBox" gorm:"column:QtyPerBox;null"`
	Cost            decimal.Decimal `json:"Cost" gorm:"column:Cost;type:decimal(10,2);null"`
	Price           decimal.Decimal `json:"Price" gorm:"column:Price;type:decimal(10,2);null"`
	Price5          decimal.Decimal `json:"Price_5" gorm:"column:Price_5;type:decimal(10,2);null"`
	Price7          decimal.Decimal `json:"Price_7" gorm:"column:Price_7;type:decimal(10,2);null"`
	Price10         decimal.Decimal `json:"Price_10" gorm:"column:Price_10;type:decimal(10,2);null"`
	Price15         decimal.Decimal `json:"Price_15" gorm:"column:Price_15;type:decimal(10,2);null"`
	Price19         decimal.Decimal `json:"Price_19" gorm:"column:Price_19;type:decimal(10,2);null"`
	ItemDimension   string          `json:"ItemDimension" gorm:"column:ItemDimension;type:text;null"`
	Length          int             `json:"Length" gorm:"column:Length;null"`
	Width           int             `json:"Width" gorm:"column:Width;null"`
	Height          int             `json:"Height" gorm:"column:Height;null"`
	BoxDimension    string          `json:"BoxDimension" gorm:"column:BoxDimension;type:text;null"`
	BoxLength       int             `json:"Box_Length" gorm:"column:Box_Length;null"`
	BoxWidth        int             `json:"Box_Width" gorm:"column:Box_Width;null"`
	BoxHeight       int             `json:"Box_Height" gorm:"column:Box_Height;null"`
	BoxWeight       int             `json:"Box_Weight" gorm:"column:Box_Weight;null"`
	AvailableDate   time.Time       `json:"AvailableDate" gorm:"column:AvailableDate;type:date;null"`
	ShippingMethod  string          `json:"Shipping_Method" gorm:"column:Shipping_Method;type:text;null"`
	PiecesContainer int             `json:"PiecesContainer" gorm:"column:PiecesContainer;null"`
	Supplier        string          `json:"Supplier" gorm:"column:Supplier;type:text;null"`
	ShippingCost    decimal.Decimal `json:"ShippingCost" gorm:"column:ShippingCost;type:decimal(10,2);null"`
	Active          string          `json:"Active" gorm:"column:Active;type:char(1);null"`
	CreatedBy       string          `json:"CreatedBy" gorm:"column:CreatedBy;type:text;size:20;null"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// TableName sets the default table name for the Item struct
func (Item) TableName() string {
	return "item" // Replace with your actual table name
}

// ProcessedItem is used to pass data to the templ view
type ProcessedItem struct {
	Values []string
}

// GetFieldNames extracts field names from the Item struct
func GetFieldNames(item interface{}) []string {
	var fieldNames []string
	t := reflect.TypeOf(item)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}
	return fieldNames
}

func LoadItemTable() error {

	// Load CSV data
	csvFile, err := os.Open("./controllers/data/item.csv") // Replace with your CSV file path
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

	// Loop through the records and create items
	for _, record := range records {

		// Parse integers and decimals
		iInventory, err := strconv.Atoi(record[6])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue
		}
		iQtyPerBox, err := strconv.Atoi(record[7])
		if err != nil {
			log.Printf("Failed to convert to int - QtyPerBox: %v", err)
			continue
		}
		iLength, err := strconv.Atoi(record[16])
		if err != nil {
			log.Printf("Failed to convert to int - Length: %v", err)
			continue
		}
		iWidth, err := strconv.Atoi(record[17])
		if err != nil {
			log.Printf("Failed to convert to int - Width: %v", err)
			continue
		}
		iHeight, err := strconv.Atoi(record[18])
		if err != nil {
			log.Printf("Failed to convert to int - Height: %v", err)
			continue
		}
		iBoxLength, err := strconv.Atoi(record[20])
		if err != nil {
			log.Printf("Failed to convert to int - BoxLength: %v", err)
			continue
		}
		iBoxWidth, err := strconv.Atoi(record[21])
		if err != nil {
			log.Printf("Failed to convert to int - BoxWidth: %v", err)
			continue
		}
		iBoxHeight, err := strconv.Atoi(record[22])
		if err != nil {
			log.Printf("Failed to convert to int - BoxHeight: %v", err)
			continue
		}
		iBoxWeight, err := strconv.Atoi(record[23])
		if err != nil {
			log.Printf("Failed to convert to int - BoxWeight: %v", err)
			continue
		}
		iPiecesContainer, err := strconv.Atoi(record[26])
		if err != nil {
			log.Printf("Failed to convert to int - PiecesContainer: %v", err)
			continue
		}
		dCost, err := decimal.NewFromString(record[8])
		if err != nil {
			log.Printf("Failed to convert to decimal - Cost: %v", err)
			continue
		}
		dShippingCost, err := decimal.NewFromString(record[28])
		if err != nil {
			log.Printf("Failed to convert to decimal - ShippingCost: %v", err)
			continue
		}
		dPrice, err := decimal.NewFromString(record[9])
		if err != nil {
			log.Printf("Failed to convert to decimal - Price: %v", err)
			continue
		}
		dPrice5, err := decimal.NewFromString(record[10])
		if err != nil {
			log.Printf("Failed to convert to decimal - Price5: %v", err)
			continue
		}
		dPrice7, err := decimal.NewFromString(record[11])
		if err != nil {
			log.Printf("Failed to convert to decimal - Price7: %v", err)
			continue
		}
		dPrice10, err := decimal.NewFromString(record[12])
		if err != nil {
			log.Printf("Failed to convert to decimal - Price10: %v", err)
			continue
		}
		dPrice15, err := decimal.NewFromString(record[13])
		if err != nil {
			log.Printf("Failed to convert to decimal - Price15: %v", err)
			continue
		}
		dPrice19, err := decimal.NewFromString(record[14])
		if err != nil {
			log.Printf("Failed to convert to decimal - Price19: %v", err)
			continue
		}

		// Parse the AvailableDate
		var availableDate time.Time
		if record[24] != "" {
			availableDate, err = time.Parse(dateFormat, record[24])
			if err != nil {
				log.Printf("Failed to parse date - AvailableDate: %v", err)
				continue
			}
		}

		// Create the item record
		itemrecord := &Item{
			SKU:             record[0],
			ItemName:        record[1],
			UPC:             record[2],
			Type:            record[3],
			Category:        record[4],
			Description:     record[5],
			Inventory:       iInventory,
			QtyPerBox:       iQtyPerBox,
			Cost:            dCost,
			Price:           dPrice,
			Price5:          dPrice5,
			Price7:          dPrice7,
			Price10:         dPrice10,
			Price15:         dPrice15,
			Price19:         dPrice19,
			ItemDimension:   record[15],
			Length:          iLength,
			Width:           iWidth,
			Height:          iHeight,
			BoxDimension:    record[19],
			BoxLength:       iBoxLength,
			BoxWidth:        iBoxWidth,
			BoxHeight:       iBoxHeight,
			BoxWeight:       iBoxWeight,
			AvailableDate:   availableDate, // Assign parsed date
			ShippingMethod:  record[25],
			PiecesContainer: iPiecesContainer,
			Supplier:        record[27],
			ShippingCost:    dShippingCost,
			Active:          record[29],
			CreatedBy:       record[30],
			CreatedAt:       time.Now(),
		}

		// Save item to the database
		err = PgDBConn.Create(&itemrecord).Error
		if err != nil {
			log.Printf("Failed to insert item record: %v", err)
		}
	}

	fmt.Println("....Item Data imported successfully!")

	return nil
}

// CreateItem adds a new item to the database
func CreateItem(item *Item) error {
	err := PgDBConn.Create(item).Error
	if err != nil {
		log.Printf("Failed to create item: %v", err)
		return err
	}
	return nil
}

// GetItemBySKU fetches item details from the database using the SKU
func GetItemBySKU(sku string) (*Item, error) {
	if sku == "" {
		return nil, errors.New("SKU cannot be empty")
	}

	var item Item
	err := PgDBConn.Debug().
		Where(`"SKU" = ?`, sku).
		First(&item).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Item with SKU '%s' not found", sku)
			return nil, fmt.Errorf("item not found for SKU '%s'", sku)
		}
		log.Printf("Failed to fetch item with SKU '%s': %v", sku, err)
		return nil, fmt.Errorf("database error for SKU '%s': %w", sku, err)
	}

	// Log the fetched item for debugging
	log.Printf("Fetched item: %+v", item)

	// Handle Active field safely
	isActive := false // Default to false
	if item.Active != "" {
		isActive = item.Active == "Y"
	}
	log.Printf("Item '%s' active status: %v", sku, isActive)

	// Check for nil or empty Description
	if item.Description == "" {
		log.Printf("Item '%s' has no description", sku)
	}

	return &item, nil
}

// GetItemDetailsBySKU fetches detailed item information from the database
func GetItemDetailsBySKU(sku string) (map[string]interface{}, error) {
	item, err := GetItemBySKU(sku)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":        item.ItemName,
		"price":       item.Price,
		"description": item.Description,
	}, nil
}

// UpdateItem updates an existing item in the database
func UpdateItem(item *Item) error {
	err := PgDBConn.Save(item).Error
	if err != nil {
		log.Printf("Failed to update item: %v", err)
		return err
	}
	return nil
}

// DeleteItem removes an item from the database by SKU
func DeleteItem(sku string) error {
	err := PgDBConn.Where("SKU = ?", sku).Delete(&Item{}).Error
	if err != nil {
		log.Printf("Failed to delete item: %v", err)
		return err
	}
	return nil
}

// ListItems retrieves all items from the database and dynamically generates field names and values
func (t *Item) GetAllItems() ([]Item, error) {
	var items []Item

	// Fetch all todos, ordered by created_at in descending order
	err := PgDBConn.Find(&items).Error

	if err != nil {
		return nil, err
	}
	return items, nil
}

func InsertItem(t *Item) error {

	if err := PgDBConn.Create(t).Error; err != nil {
		log.Fatal("failed to create TodoPG:", err)
		return err
	}
	return nil

}

// ProcessItemForView processes a single Item struct for rendering
func ProcessItemForView(item Item) map[string]interface{} {
	processed := make(map[string]interface{})

	// Format date fields
	var availableDate string
	if !item.AvailableDate.IsZero() {
		availableDate = item.AvailableDate.Format("2006-01-02")
	} else {
		availableDate = "N/A"
	}

	createdAt := item.CreatedAt.Format("2006-01-02")

	// Add fields to the processed map
	processed["SKU"] = item.SKU
	processed["ItemName"] = item.ItemName
	processed["UPC"] = item.UPC
	processed["Type"] = item.Type
	processed["Category"] = item.Category
	processed["Description"] = item.Description
	processed["Inventory"] = fmt.Sprintf("%d", item.Inventory)
	processed["QtyPerBox"] = fmt.Sprintf("%d", item.QtyPerBox)
	processed["Cost"] = item.Cost.String()
	processed["Price"] = item.Price.String()
	processed["Price5"] = item.Price5.String()
	processed["Price7"] = item.Price7.String()
	processed["Price10"] = item.Price10.String()
	processed["Price15"] = item.Price15.String()
	processed["Price19"] = item.Price19.String()
	processed["ItemDimension"] = item.ItemDimension
	processed["Length"] = fmt.Sprintf("%d", item.Length)
	processed["Width"] = fmt.Sprintf("%d", item.Width)
	processed["Height"] = fmt.Sprintf("%d", item.Height)
	processed["BoxDimension"] = item.BoxDimension
	processed["BoxLength"] = fmt.Sprintf("%d", item.BoxLength)
	processed["BoxWidth"] = fmt.Sprintf("%d", item.BoxWidth)
	processed["BoxHeight"] = fmt.Sprintf("%d", item.BoxHeight)
	processed["BoxWeight"] = fmt.Sprintf("%d", item.BoxWeight)
	processed["AvailableDate"] = availableDate
	processed["ShippingMethod"] = item.ShippingMethod
	processed["PiecesContainer"] = fmt.Sprintf("%d", item.PiecesContainer)
	processed["Supplier"] = item.Supplier
	processed["ShippingCost"] = item.ShippingCost.String()
	processed["Active"] = item.Active
	processed["CreatedBy"] = item.CreatedBy
	processed["CreatedAt"] = createdAt

	return processed
}
