package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/emarifer/gofiber-templ-htmx/controllers"
	"github.com/emarifer/gofiber-templ-htmx/models"
	"github.com/emarifer/gofiber-templ-htmx/views/item_views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for item Views **********/

// Render List Page with success/error messages
func HandleInventoryList(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	item := new(controllers.Item)
	// item.CreatedBy = c.Locals("userId").(uint64)

	// fm := fiber.Map{"type": "error"}

	itemsSlice, err := item.GetAllItems()
	if err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			// "no such table" is the error that SQLite3 produces
			// when some table does not exist, and we have only
			// used it as an example of the errors that can be caught.
			// Here you can add the errors that you are interested
			// in throwing as `500` codes.
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}
		// fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		// return flash.WithError(c, fm).Redirect("/item/list")
	}

	tindex := item_views.itemIndex(itemsSlice)
	tlist := item_views.itemList(
		" | Tasks List",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		tindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(tlist))

	return handler(c)
}

// Render Create item Page with success/error messages
func HandleViewCreatePage(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	if c.Method() == "POST" {
		newitem := new(controllers.Item)
		newitem.CreatedBy = c.Locals("userId").(uint64)
		newitem.SKU = strings.Trim(c.FormValue("SKU"), " ")
		newitem.Description = strings.Trim(c.FormValue("description"), " ")

		fm := fiber.Map{
			"type":    "error",
			"message": "Task title empty!!",
		}
		if newitem.SKU == "" {

			return flash.WithError(c, fm).Redirect("/item/list")
		}

		if _, err := newitem.Createitem(); err != nil {
			if strings.Contains(err.Error(), "no such table") ||
				strings.Contains(err.Error(), "database is locked") {
				// "no such table" is the error that SQLite3 produces
				// when some table does not exist, and we have only
				// used it as an example of the errors that can be caught.
				// Here you can add the errors that you are interested
				// in throwing as `500` codes.
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Task successfully created!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/item/list")
	}

	cindex := item_views.CreateIndex()
	create := item_views.Create(
		" | Create item",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		cindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(create))

	return handler(c)
}

// Render Edit item Page with success/error messages
func HandleViewEditPage(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)
	session, _ := store.Get(c)
	tzone := session.Get(TZONE_KEY).(string)

	idParams, _ := strconv.Atoi(c.Params("id"))
	itemId := uint64(idParams)

	item := new(controllers.Item)
	item.SKU = itemId
	item.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{"type": "error"}

	recovereditem, err := item.GetNoteById()

	if err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			// "no such table" is the error that SQLite3 produces
			// when some table does not exist, and we have only
			// used it as an example of the errors that can be caught.
			// Here you can add the errors that you are interested
			// in throwing as `500` codes.
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}

		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/item/list")
	}

	if c.Method() == "POST" {
		item.Title = strings.Trim(c.FormValue("title"), " ")
		item.Description = strings.Trim(c.FormValue("description"), " ")
		if c.FormValue("status") == "on" {
			item.Status = true
		} else {
			item.Status = false
		}

		fm = fiber.Map{
			"type":    "error",
			"message": "Task title empty!!",
		}
		if item.Title == "" {

			return flash.WithError(c, fm).Redirect("/item/list")
		}

		_, err := item.Updateitem()
		if err != nil {
			if strings.Contains(err.Error(), "no such table") ||
				strings.Contains(err.Error(), "database is locked") {
				// "no such table" is the error that SQLite3 produces
				// when some table does not exist, and we have only
				// used it as an example of the errors that can be caught.
				// Here you can add the errors that you are interested
				// in throwing as `500` codes.
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}

			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/item/list")
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Task successfully updated!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/item/list")
	}

	uindex := item_views.UpdateIndex(recovereditem, tzone)
	update := item_views.Update(
		fmt.Sprintf(" | Edit item #%d", recovereditem.ID),
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		uindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(update))

	return handler(c)
}

// Handler Remove item
func HandleDeleteitem(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	itemId := uint64(idParams)

	item := new(models.item)
	item.ID = itemId
	item.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{"type": "error"}

	if err := item.Deleteitem(); err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			// "no such table" is the error that SQLite3 produces
			// when some table does not exist, and we have only
			// used it as an example of the errors that can be caught.
			// Here you can add the errors that you are interested
			// in throwing as `500` codes.
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect(
			"/item/list",
			fiber.StatusSeeOther,
		)
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "Task successfully deleted!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/item/list", fiber.StatusSeeOther)
}
