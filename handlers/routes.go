package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/kyp0717/ew-system/views"
	"github.com/sujit-baniya/flash"
)

var store *session.Store

const (
	AUTH_KEY       string = "authenticated"
	USER_ID        string = "user_id"
	FROM_PROTECTED string = "from_protected"
	TZONE_KEY      string = "time_zone"
)

func init() {
	store = session.New(session.Config{
		Storage:        memory.New(), // or your preferred storage
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookieHTTPOnly: true,
	})
}

func Setup(app *fiber.App) {
	/* Sessions Config */
	store = session.New(session.Config{
		CookieHTTPOnly: true,
		Expiration:     time.Hour * 1,
	})

	/* Views */
	app.Get("/", flagsMiddleware, HandleViewHome)
	app.Get("/login", flagsMiddleware, HandleViewLogin)
	app.Post("/login", flagsMiddleware, HandleViewLogin)
	app.Get("/register", flagsMiddleware, HandleViewRegister)
	app.Post("/register", flagsMiddleware, HandleViewRegister)

	/* Views protected with session middleware */
	todoApp := app.Group("/todo", AuthMiddleware)
	todoApp.Get("/listpg", HandleViewListPG)
	todoApp.Get("/create", HandleViewCreatePagePG)
	todoApp.Post("/create", HandleViewCreatePagePG)
	todoApp.Get("/edit/:id", HandleViewEditPagePG)
	todoApp.Post("/edit/:id", HandleViewEditPagePG)
	todoApp.Delete("/delete/:id", HandleDeleteTodoPG)
	todoApp.Post("/logout", HandleLogout)

	/* Views protected with session middleware */
	inventoryApp := app.Group("/inventory", AuthMiddleware)
	inventoryApp.Get("/inventorylist", HandleInventoryList)
	inventoryApp.Post("/create", HandleInventoryCreate)
	inventoryApp.Get("/edit/:id", HandleInventoryEdit)
	inventoryApp.Post("/edit/:id", HandleInventoryEdit)
	inventoryApp.Delete("/delete/:id", HandleInventoryDelete)

	/* SKU Item Details Group */
	itemsApp := app.Group("/items")
	itemsApp.Get("/:sku", HandleItemDetails)
	itemsApp.Post("/save-items", HandleSaveItems)
	itemsApp.Post("/update-item", HandleUpdateItem)

	/* ↓ Not Found Management - Fallback Page ↓ */
	app.Get("/*", flagsMiddleware, func(c *fiber.Ctx) error {
		return fiber.NewError(
			fiber.StatusNotFound,
			"error 404: not found",
		)
	})
}

// CustomErrorHandler does centralized error handling.
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Status(code)

	var errorIndex templ.Component

	switch code {
	case 404:
		errorIndex = views.Error404(fromProtected)
	default:
		if c.Route().Path == "/todo/list" {
			// If the path `/todo/list` cannot obtain the to-dos
			// from the database, the error page will only allow the user
			// to return to the home page (fromProtected = false).
			errorIndex = views.Error500(false, code, e.Message)
		} else {
			errorIndex = views.Error500(fromProtected, code, e.Message)
		}
	}

	pageTitle := fmt.Sprintf(" | Error %d", code)

	errorPage := views.Home(
		pageTitle, fromProtected, true, flash.Get(c), errorIndex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(errorPage))

	return handler(c)
}

// flagsMiddleware is a middleware for handling different behaviors
// of non protected pages, specifically not allowing an already
// logged in user to log in or register again.
func flagsMiddleware(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	userId := sess.Get(USER_ID)
	if userId == nil {
		c.Locals(FROM_PROTECTED, false)

		return c.Next()
	}

	c.Locals(FROM_PROTECTED, true)

	return c.Next()
}
