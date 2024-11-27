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

	// Step 2: Extract field names dynamically
	fieldNames := utility.GetFieldNames(controllers.Item{})

	// Step 3: Convert items to ProcessedItem for rendering
	processedItems := utility.MapToProcessedItems(items)

	// Step 4: Pagination parameters
	pageNumber, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10")) // Default page size is 10
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Calculate the start and end indices for pagination
	start := (pageNumber - 1) * pageSize
	if start >= len(processedItems) {
		start = len(processedItems)
	}
	end := start + pageSize
	if end > len(processedItems) {
		end = len(processedItems)
	}

	// Slice the processedItems to only include the items for the current page
	paginatedItems := processedItems[start:end]

	// Step 5: Populate the SearchBarArgs with inventory-specific data
	fromMenu := "inventory"
	searchBarArgs := utility.ProcessedSearchBarArgs(fromMenu, items)

	// Debugging: Log the search bar arguments to confirm FromMenu is set
	log.Printf("SearchBarArgs after ProcessedSearchBarArgs: %+v", searchBarArgs)

	// Step 6: Render the inventory list template with paginated items
	iindex := item_views.ListItemIndex(paginatedItems, fieldNames, pageNumber, pageSize, len(processedItems))
	ilist := item_views.ListItem(
		" | Inventory List",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		iindex,
		searchBarArgs,
	)

	// Step 7: Send the rendered template as the response
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
		controllers.SearchBarArgs{},
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
		controllers.SearchBarArgs{},
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
	//item.CreatedBy = c.Locals("username").(string)

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
