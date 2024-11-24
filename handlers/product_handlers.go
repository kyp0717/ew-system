package handlers

import (
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"
	"github.com/kyp0717/ew-system/handlers/utility"
	"github.com/kyp0717/ew-system/views/product_views"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Product Views **********/
/********** Handlers for Product Views **********/
func HandleProductList(c *fiber.Ctx) error {
	// Step 1: Retrieve "fromProtected" flag safely
	fromProtected, ok := c.Locals(FROM_PROTECTED).(bool)
	if !ok {
		log.Printf("fromProtected not set or invalid. Defaulting to false.")
		fromProtected = false // Default to false if not set
	}

	// Step 2: Retrieve session
	session, err := store.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve session")
	}

	// Step 3: Retrieve username from session
	username, ok := session.Get("username").(string)
	if !ok || username == "" {
		log.Printf("User is not authenticated. Username missing in session.")
		return fiber.NewError(fiber.StatusUnauthorized, "User is not authenticated")
	}

	log.Printf("Authenticated user: %s", username)

	// Step 4: Fetch products from the database
	var products []controllers.Product
	err = controllers.PgDBConn.Debug().Find(&products).Error
	if err != nil {
		log.Printf("Database error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("something went wrong: %s", err))
	}

	// Step 5: Extract field names dynamically
	fieldNames := utility.GetFieldNames(controllers.Product{})
	log.Printf("Field names: %+v", fieldNames)

	// Step 6: Convert products to ProcessedProduct for rendering
	processedProducts := utility.MapToProcessProducts(products)
	log.Printf("utility.MapToProcessProducts: Convert products to ProcessedProduct for rendering")

	// Step 7: Retrieve flash messages safely
	msg := flash.Get(c)
	if msg == nil {
		msg = fiber.Map{}
	}

	// Step 8: Render the product list index
	pindex := product_views.ListProductIndex(processedProducts, fieldNames)
	if pindex == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to render product index")
	}

	// Step 9: Render the full product list view
	plist := product_views.ListProduct(
		" | Product List",
		fromProtected,
		false, // Assuming no error occurred, set isError to false
		msg,
		username, // Pass the username retrieved from the session
		pindex,
	)

	// Step 10: Adapt and send the rendered template as the response
	handler := adaptor.HTTPHandler(templ.Handler(plist))
	return handler(c)
}
