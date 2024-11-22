package handlers

import (
	"errors"
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"
	"github.com/kyp0717/ew-system/views/item_views"
	"gorm.io/gorm"
)

func HandleItemDetails(c *fiber.Ctx) error {
	// Extract the SKU parameter
	sku := c.Params("sku")
	if sku == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SKU is required"})
	}

	// Fetch the item
	item, err := controllers.GetItemBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An unexpected error occurred"})
	}

	// Convert item to map
	itemDetails := item_views.ItemDetails(
		controllers.ProcessItemForView(*item),
	)

	// Get fromProtected and username with fallback defaults
	fromProtected, ok := c.Locals("fromProtected").(bool)
	if !ok {
		fromProtected = false // Default value
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		username = "Guest" // Default value
	}

	// Render the page
	itemPage := item_views.ListItemBySKU(
		" | Item Details",
		fromProtected,
		false,
		nil, // Messages (optional)
		username,
		itemDetails,
	)

	handler := adaptor.HTTPHandler(templ.Handler(itemPage))
	return handler(c)
}

func HandleItemSave(c *fiber.Ctx) error {
	var item controllers.Item

	// Step 1: Parse the request body
	if err := c.BodyParser(&item); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	// Log the parsed item for debugging
	log.Printf("Parsed item from request: %+v", item)

	// Step 2: Validate required fields
	if item.SKU == "" {
		log.Println("Validation error: SKU is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SKU is required"})
	}

	// Step 3: Save the item to the database
	if err := controllers.UpdateItem(&item); err != nil {
		log.Printf("Failed to update item with SKU '%s': %v", item.SKU, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save changes"})
	}

	// Log success
	log.Printf("Successfully updated item: %+v", item)

	// Step 4: Return success response
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Item updated successfully",
		"item":    item, // Optionally include the updated item in the response
	})
}
