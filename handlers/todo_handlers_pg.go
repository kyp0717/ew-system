package handlers

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"
	"github.com/kyp0717/ew-system/views/todo_views"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Todo Views **********/

// Render List Page with success/error messages
func HandleViewListPG(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	todo := new(controllers.TodoPG)
	todo.CreatedBy = c.Locals("userId").(uint64)

	// fm := fiber.Map{"type": "error"}

	todosSlice, err := todo.GetAllTodos()
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

		// return flash.WithError(c, fm).Redirect("/todo/list")
	}

	tindex := todo_views.TodoIndexPG(todosSlice)
	tlist := todo_views.TodoListPG(
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

// Render Create Todo Page with success/error messages
// Send data to the go backend
func HandleViewCreatePagePG(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	if c.Method() == "POST" {
		newTodo := new(controllers.TodoPG)
		newTodo.CreatedBy = c.Locals("userId").(uint64)
		newTodo.Title = strings.Trim(c.FormValue("title"), " ")
		newTodo.Description = strings.Trim(c.FormValue("description"), " ")

		fm := fiber.Map{
			"type":    "error",
			"message": "Task title empty!!",
		}
		if newTodo.Title == "" {

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		if err := controllers.InsertTodoPG(*newTodo); err != nil {
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

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	cindex := todo_views.CreateIndexPG()
	create := todo_views.CreatePG(
		" | Create Todo",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		cindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(create))

	return handler(c)
}
