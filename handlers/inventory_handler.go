package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"
	"github.com/kyp0717/ew-system/handlers/utility"
	"github.com/kyp0717/ew-system/views/item_views"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Inventory Views **********/
func HandleInventoryList(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	// Step 1: Fetch items from the database
	var items []controllers.Item
	err := controllers.PgDBConn.Debug().Find(&items).Error
	if err != nil {
		log.Printf("Database error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("something went wrong: %s", err))
	}
	log.Printf("Fetched items: %+v", items)

	// Step 2: Extract field names dynamically
	fieldNames := utility.GetFieldNames(controllers.Item{})
	log.Printf("Field names: %+v", fieldNames)

	// Step 3: Convert items to ProcessedItem for rendering
	processedItems := utility.MapToProcessedItems(items)
	log.Printf("Processed items: %+v", processedItems)

	// Step 4: Render the inventory list template
	iindex := item_views.ListItemIndex(processedItems, fieldNames)
	ilist := item_views.ListItem(
		" | Inventory List",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		iindex,
	)

	// Step 5: Send the rendered template as the response
	handler := adaptor.HTTPHandler(templ.Handler(ilist))
	return handler(c)
}

// Render Create Item Page with success/error messages
func HandleInventoryCreate(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	if c.Method() == "POST" {
		newItem := new(controllers.Item)
		newItem.CreatedBy = c.Locals("username").(string)
		newItem.ItemName = strings.Trim(c.FormValue("name"), " ")
		newItem.Description = strings.Trim(c.FormValue("description"), " ")

		fm := fiber.Map{
			"type":    "error",
			"message": "Item name empty!!",
		}
		if newItem.ItemName == "" {
			return flash.WithError(c, fm).Redirect("/inventory/list")
		}

		if err := controllers.InsertItem(newItem); err != nil {
			if strings.Contains(err.Error(), "no such table") ||
				strings.Contains(err.Error(), "database is locked") {
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Item successfully created!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/inventory/list")
	}

	cindex := item_views.CreateItemIndex()
	create := item_views.CreateItem(
		" | Create Item",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		cindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(create))

	return handler(c)
}

// Render Edit Item Page with success/error messages
func HandleInventoryEdit(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)
	session, _ := store.Get(c)
	tzone := session.Get(TZONE_KEY).(string)

	idParams, _ := strconv.Atoi(c.Params("id"))
	itemSKU := uint64(idParams)

	item := new(controllers.Item)
	item.SKU = strconv.FormatUint(itemSKU, 10)
	item.CreatedBy = c.Locals("username").(string)

	fm := fiber.Map{"type": "error"}

	recoveredItem, err := controllers.GetItemBySKU(item.SKU)

	if err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}

		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/inventory/list")
	}

	if c.Method() == "POST" {
		item.ItemName = strings.Trim(c.FormValue("name"), " ")
		item.Description = strings.Trim(c.FormValue("description"), " ")

		fm = fiber.Map{
			"type":    "error",
			"message": "Item name empty!!",
		}
		if item.ItemName == "" {
			return flash.WithError(c, fm).Redirect("/inventory/list")
		}

		err := controllers.UpdateItem(item)
		if err != nil {
			if strings.Contains(err.Error(), "no such table") ||
				strings.Contains(err.Error(), "database is locked") {
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}

			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/inventory/list")
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Item successfully updated!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/inventory/list")
	}

	uindex := item_views.UpdateItemIndex(*recoveredItem, tzone)
	update := item_views.UpdateItem(
		fmt.Sprintf(" | Edit Item #%d", itemSKU),
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		uindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(update))

	return handler(c)
}

// Handler Remove Item
func HandleInventoryDelete(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	itemSKU := uint64(idParams)

	item := new(controllers.Item)
	item.SKU = strconv.FormatUint(itemSKU, 10)
	item.CreatedBy = c.Locals("username").(string)

	fm := fiber.Map{"type": "error"}

	if err := controllers.DeleteItem(item.SKU); err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect(
			"/inventory/list",
			fiber.StatusSeeOther,
		)
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "Item successfully deleted!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/inventory/list", fiber.StatusSeeOther)
}

func HandleInventoryDetails(c *fiber.Ctx) error {
	sku := c.Params("sku")

	// Fetch item details
	itemDetails, err := controllers.GetItemDetailsBySKU(sku)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	// Render the page with item details
	return c.Render("itemDetailsPage", fiber.Map{
		"sku":         sku,
		"itemDetails": itemDetails,
	})
}
func HandleSKUItemDetails(c *fiber.Ctx) error {
	// Fetch the item from the database
	sku := c.Params("sku")
	item, err := controllers.GetItemBySKU(sku)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	// Convert the struct to a map
	itemMap := utility.ConvertStructToMap(item)

	// Render the Templ page with the map
	return c.Render("sku_item.templ", fiber.Map{
		"itemMap": itemMap,
	})
}
