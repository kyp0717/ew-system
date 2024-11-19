package utility

import (
	"fmt"
	"reflect"

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
		processed = append(processed, controllers.ProcessedItem{
			Values: []string{
				item.SKU,                          // SKU
				item.ItemName,                     // Item Name
				item.UPC,                          // UPC
				item.Type,                         // Type
				item.Category,                     // Category
				item.Description,                  // Description
				fmt.Sprintf("%d", item.Inventory), // Inventory as a string
				fmt.Sprintf("%d", item.QtyPerBox), // Quantity per box
				item.Cost.String(),                // Cost as a string
				item.Price.String(),               // Price as a string
				item.Price5.String(),              // Price for 5 units as a string
				item.Price7.String(),              // Price for 7 units as a string
				item.ItemDimension,                // Item Dimension
				fmt.Sprintf("%d x %d x %d", item.Length, item.Width, item.Height), // Dimensions
				item.BoxDimension,                   // Box Dimension
				fmt.Sprintf("%d", item.BoxWeight),   // Box Weight as a string
				item.AvailableDate,                  // Available Date
				item.ShippingMethod,                 // Shipping Method
				item.Supplier,                       // Supplier
				item.ShippingCost.String(),          // Shipping Cost as a string
				item.Active,                         // Active Status
				item.CreatedBy,                      // Created By
				item.CreatedAt.Format("2006-01-02"), // Created At (formatted date)
			},
		})
	}

	return processed
}
