package utility

import (
	"fmt"
	"log"
	"reflect" // Ensure the time package is imported

	"github.com/kyp0717/ew-system/controllers"
)

// GetFieldNames extracts field names from a struct
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

// MapToProcessedItems converts items to ProcessedItem for templ rendering
func MapToProcessedItems(items []controllers.Item) []controllers.ProcessedItem {
	processed := []controllers.ProcessedItem{}

	for _, item := range items {
		// Format AvailableDate as a string (e.g., "YYYY-MM-DD")
		availableDate := ""
		if !item.AvailableDate.IsZero() { // Check for zero-value date
			availableDate = item.AvailableDate.Format("2006-01-02")
		}

		processed = append(processed, controllers.ProcessedItem{
			Values: []string{
				item.SKU,                          // SKU
				item.ItemName,                     // Item Name
				item.UPC,                          // UPC
				item.Type,                         // Type
				item.Category,                     // Category
				item.Description,                  // Description
				fmt.Sprintf("%d", item.Inventory), // Inventory as string
				fmt.Sprintf("%d", item.QtyPerBox), // Quantity per box as string
				item.Cost.String(),                // Cost as string
				item.Price.String(),               // Price as string
				item.Price5.String(),              // Price for 5 units as string
				item.Price7.String(),              // Price for 7 units as string
				item.Price10.String(),             // Price for 10 units as string
				item.Price15.String(),             // Price for 15 units as string
				item.Price19.String(),             // Price for 19 units as string
				item.ItemDimension,                // Item Dimension
				//fmt.Sprintf("%d x %d x %d", item.Length, item.Width, item.Height), // Dimensions
				fmt.Sprintf("%d", item.Length), // Length
				fmt.Sprintf("%d", item.Width),  // Width
				fmt.Sprintf("%d", item.Height), // Height
				item.BoxDimension,              // Box Dimension
				//fmt.Sprintf("%d", item.BoxWeight), // Box Weight as string
				fmt.Sprintf("%d", item.BoxLength), // Box Length
				fmt.Sprintf("%d", item.BoxWidth),  // Box Width
				fmt.Sprintf("%d", item.BoxHeight), // Box Height
				//fmt.Sprintf("%d", item.BoxWeight), // Box Weight as string
				fmt.Sprintf("%d", item.BoxWeight),       // Box Weight
				availableDate,                           // Available Date (formatted as string)
				item.ShippingMethod,                     // Shipping Method
				fmt.Sprintf("%d", item.PiecesContainer), // Picked
				item.Supplier,                           // Supplier
				item.ShippingCost.String(),              // Shipping Cost as string
				item.Active,                             // Active Status
				item.CreatedBy,                          // Created By
				item.CreatedAt.Format("2006-01-02"),     // Created At (formatted as string)
			},
		})
	}

	return processed
}

// MapToProcessProducts converts Product data to ProcessedProduct for templated rendering
func MapToProcessProducts(products []controllers.Product) []controllers.ProcessedProduct {
	processed := []controllers.ProcessedProduct{}

	for _, product := range products {
		// Safely format CreateDate
		createDate := ""
		if !product.CreateDate.IsZero() {
			createDate = product.CreateDate.Format("2006-01-02")
		}

		// Safeguard against zero-value integers
		formatInt := func(value int) string {
			if value == 0 {
				return "0"
			}
			return fmt.Sprintf("%d", value)
		}

		// Append processed product
		processed = append(processed, controllers.ProcessedProduct{
			Values: []string{
				product.SKU,                   // SKU
				product.Category,              // Category
				product.Group,                 // Group
				product.ProductName,           // Product Name
				createDate,                    // Create Date
				formatInt(product.TotalBoxes), // Total Boxes

				product.SKU1,              // SKU_1
				formatInt(product.Box1),   // Box_1
				formatInt(product.Piece1), // Piece_1

				product.SKU2,              // SKU_2
				formatInt(product.Box2),   // Box_2
				formatInt(product.Piece2), // Piece_2

				product.SKU3,              // SKU_3
				formatInt(product.Box3),   // Box_3
				formatInt(product.Piece3), // Piece_3

				product.SKU4,              // SKU_4
				formatInt(product.Box4),   // Box_4
				formatInt(product.Piece4), // Piece_4

				product.SKU5,              // SKU_5
				formatInt(product.Box5),   // Box_5
				formatInt(product.Piece5), // Piece_5

				product.SKU6,              // SKU_6
				formatInt(product.Box6),   // Box_6
				formatInt(product.Piece6), // Piece_6

				product.SKU7,              // SKU_7
				formatInt(product.Box7),   // Box_7
				formatInt(product.Piece7), // Piece_7

				product.Active,                         // Active Status
				product.CreatedBy,                      // Created By
				product.CreatedAt.Format("2006-01-02"), // Created At (formatted as string)
			},
		})
	}
	log.Printf("utility.MapToProcessProducts: return ProcessedProduct")
	return processed
}

// ConvertStructToMap converts a struct to a map of field names and values
func ConvertStructToMap(item interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		result[field.Name] = value
	}
	return result
}
