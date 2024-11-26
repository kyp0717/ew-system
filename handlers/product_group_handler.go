package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"
	"github.com/kyp0717/ew-system/handlers/utility"
	"github.com/kyp0717/ew-system/views/product_views"
	"github.com/shopspring/decimal"
	"github.com/sujit-baniya/flash"
	"gorm.io/gorm"
)

func HandleProductDetails(c *fiber.Ctx) error {
	// Step 1: Extract the SKU parameter
	sku := c.Params("sku")
	if sku == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SKU is required"})
	}

	// Step 2: Fetch the product details from the database
	product, err := controllers.GetProductBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An unexpected error occurred"})
	}

	// Step 3: Extract `fromProtected` and `username` from the context
	fromProtected, ok := c.Locals("fromProtected").(bool)
	if !ok {
		fromProtected = false // Default value for unauthenticated users
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		username = "Guest" // Default value for guests
	}

	// Step 4: Convert the product to a map for rendering using utility function
	productMap := utility.ConvertProductGroupToMap(product)

	// Step 5: Render the product details page
	productDetailsComponent := product_views.ProductDetails(productMap) // Create the component for the details view

	productPage := product_views.ListProductBySKU(
		" | Product Details",    // Page title
		fromProtected,           // Whether the user is authenticated
		false,                   // Not an error page
		flash.Get(c),            // Optional flash messages
		username,                // Username for personalization
		productDetailsComponent, // The main content (product details component)
	)

	// Step 6: Adapt the handler to Fiber
	handler := adaptor.HTTPHandler(templ.Handler(productPage))
	return handler(c)
}

func HandleSaveProducts(c *fiber.Ctx) error {
	// Parse the form data into a map
	var itemsData []map[string]string
	if err := c.BodyParser(&itemsData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	for rowIndex, itemRow := range itemsData {
		// Extract fields for each item
		sku := itemRow["SKU"]
		itemName := itemRow["ItemName"]
		upc := itemRow["UPC"]
		itemType := itemRow["Type"]
		category := itemRow["Category"]
		description := itemRow["Description"]
		inventoryStr := itemRow["Inventory"]
		qtyPerBoxStr := itemRow["QtyPerBox"]
		costStr := itemRow["Cost"]
		priceStr := itemRow["Price"]

		// Parse numeric and decimal fields
		inventory, err := strconv.Atoi(inventoryStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid inventory format for row %d", rowIndex),
			})
		}

		qtyPerBox, err := strconv.Atoi(qtyPerBoxStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid quantity per box format for row %d", rowIndex),
			})
		}

		cost, err := decimal.NewFromString(costStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid cost format for row %d", rowIndex),
			})
		}

		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid price format for row %d", rowIndex),
			})
		}

		// Create an updated item object
		updatedItem := controllers.Item{
			SKU:         sku,
			ItemName:    itemName,
			UPC:         upc,
			Type:        itemType,
			Category:    category,
			Description: description,
			Inventory:   inventory,
			QtyPerBox:   qtyPerBox,
			Cost:        cost,
			Price:       price,
		}

		// Update the item in the database
		err = controllers.UpdateItem(&updatedItem) // Pass the pointer to updatedItem
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to update item for row %d: %s", rowIndex, err.Error()),
			})
		}
	}

	// Redirect back to the inventory list or send a success response
	return c.Redirect("/inventory")
}

func HandleUpdateProduct(c *fiber.Ctx) error {
	// Parse the incoming JSON data into the Item struct
	var updatedItem controllers.Item
	if err := c.BodyParser(&updatedItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	// Validate required fields
	if updatedItem.SKU == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SKU is required"})
	}

	// Update the item in the database
	result := controllers.PgDBConn.Model(&controllers.Item{}).Where("SKU = ?", updatedItem.SKU).Updates(updatedItem)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update item"})
	}

	// Send a success response
	return c.JSON(fiber.Map{"message": "Item updated successfully"})
}
