package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"
	"github.com/kyp0717/ew-system/views/todo_views"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Todo Views **********/
func HandleViewListPG(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	// Retrieve user ID from context
	userID := c.Locals("userId").(uint64)

	var todos []controllers.TodoPG

	// Fetch all todos created by the user
	err := controllers.PgDBConn.Where("created_by = ?", userID).Find(&todos).Error
	if err != nil {
		// Handle specific PostgreSQL error cases (such as table or connection issues)
		if strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "could not connect to server") {
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}
		// General error handling
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("something went wrong: %s", err))
	}

	// Process the todos for rendering
	tindex := todo_views.TodoIndexPG(todos)
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

			return flash.WithError(c, fm).Redirect("/todo/listpg")
		}

		if err := controllers.InsertTodoPG(newTodo); err != nil {
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

		return flash.WithSuccess(c, fm).Redirect("/todo/listpg")
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

// Render Edit Todo Page with success/error messages
func HandleViewEditPagePG(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)
	session, _ := store.Get(c)
	tzone := session.Get(TZONE_KEY).(string)

	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(controllers.TodoPG)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{"type": "error"}

	recoveredTodo, err := todo.GetTodoById()

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

		return flash.WithError(c, fm).Redirect("/todo/listpg")
	}

	if c.Method() == "POST" {
		todo.Title = strings.Trim(c.FormValue("title"), " ")
		todo.Description = strings.Trim(c.FormValue("description"), " ")
		if c.FormValue("status") == "on" {
			todo.Status = true
		} else {
			todo.Status = false
		}

		fm = fiber.Map{
			"type":    "error",
			"message": "Task title empty!!",
		}
		if todo.Title == "" {

			return flash.WithError(c, fm).Redirect("/todo/listpg")
		}

		_, err := todo.UpdateTodo()
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

			return flash.WithError(c, fm).Redirect("/todo/listpg")
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Task successfully updated!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/listpg")
	}

	uindex := todo_views.UpdateIndex(recoveredTodo, tzone)
	update := todo_views.Update(
		fmt.Sprintf(" | Edit Todo #%d", recoveredTodo.ID),
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		uindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(update))

	return handler(c)
}

// Handler Remove Todo
func HandleDeleteTodoPG(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(controllers.TodoPG)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{"type": "error"}

	if err := todo.DeleteTodo(); err != nil {
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
			"/todo/listpg",
			fiber.StatusSeeOther,
		)
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "Task successfully deleted!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/todo/listpg", fiber.StatusSeeOther)
}
