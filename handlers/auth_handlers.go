package handlers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kyp0717/ew-system/controllers"

	"github.com/kyp0717/ew-system/views"
	"github.com/kyp0717/ew-system/views/auth_views"
	"github.com/sujit-baniya/flash"
	"golang.org/x/crypto/bcrypt"
)

/********** Handlers for Auth Views **********/

// Render Home Page
func HandleViewHome(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	hindex := views.HomeIndex(fromProtected)
	home := views.Home("", fromProtected, false, flash.Get(c), hindex)

	handler := adaptor.HTTPHandler(templ.Handler(home))

	return handler(c)
}

// Render Login Page with success/error messages & session management
func HandleViewLogin(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	lindex := auth_views.LoginIndex(fromProtected)
	login := auth_views.Login(
		" | Login", fromProtected, false, flash.Get(c), lindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(login))

	if c.Method() == "POST" {
		tzone := ""
		if len(c.GetReqHeaders()["X-Timezone"]) != 0 {
			tzone = c.GetReqHeaders()["X-Timezone"][0]
		}

		var (
			user controllers.User
			err  error
		)
		fm := fiber.Map{"type": "error"}

		// Fetch user by email
		if user, err = controllers.CheckEmail(c.FormValue("email")); err != nil {
			log.Printf("Failed to fetch user by email: %s, error: %v", c.FormValue("email"), err)
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connection reset by peer") {
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}
			fm["message"] = "There is no user with that email"
			return flash.WithError(c, fm).Redirect("/login")
		}

		log.Printf("Fetched user: %+v", user)

		// Verify password
		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		); err != nil {
			fm["message"] = "Incorrect password"
			return flash.WithError(c, fm).Redirect("/login")
		}

		// Set session
		session, err := store.Get(c)
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)
			return flash.WithError(c, fm).Redirect("/login")
		}

		log.Printf("Setting session: AUTH_KEY=true, USER_ID=%d, TZONE_KEY=%s", user.ID, tzone)
		session.Set(AUTH_KEY, true)
		session.Set(USER_ID, user.ID)
		session.Set(TZONE_KEY, tzone)
		session.Set("username", user.Username) // Store the username in the session

		err = session.Save()
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)
			return flash.WithError(c, fm).Redirect("/login")
		}

		log.Printf("Session before redirect: %+v", session)
		log.Printf("AUTH_KEY: %v, USER_ID: %v", session.Get(AUTH_KEY), session.Get(USER_ID))

		fm = fiber.Map{
			"type":    "success",
			"message": "You have successfully logged in!!",
		}

		//return flash.WithSuccess(c, fm).Redirect("/todo/listpg")
		// Redirect to the homepage instead of /todo/listpg
		return flash.WithSuccess(c, fm).Redirect("/")
	}

	return handler(c)
}

// Render Register Page with success/error messages
func HandleViewRegister(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	rindex := auth_views.RegisterIndex(fromProtected)
	register := auth_views.Register(
		" | Register", fromProtected, false, flash.Get(c), rindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(register))

	if c.Method() == "POST" {
		user := controllers.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			Username: c.FormValue("username"),
		}

		//err := models.CreateUser(user)
		err := controllers.CreateUser(&user)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") ||
				strings.Contains(err.Error(), "connection reset by peer") {
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}
			if strings.Contains(err.Error(), "duplicate key value") {
				err = errors.New("the email is already in use")
			}
			fm := fiber.Map{
				"type":    "error",
				"message": fmt.Sprintf("something went wrong: %s", err),
			}
			return flash.WithError(c, fm).Redirect("/register")
		}

		fm := fiber.Map{
			"type":    "success",
			"message": "You have successfully registered!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/login")
	}

	return handler(c)
}

// Authentication Middleware
func AuthMiddleware(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "error",
	}

	session, err := store.Get(c)
	fmt.Printf("Session: %+v\n", session)
	fmt.Printf("Error: %v\n", err)
	if err != nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	if session.Get(AUTH_KEY) == nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	userId := session.Get(USER_ID)
	if userId == nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	user, err := controllers.GetUserById(fmt.Sprint(userId.(uint64)))
	if err != nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	c.Locals("userId", userId)
	c.Locals("username", user.Username)
	c.Locals(FROM_PROTECTED, true)
	// fromProtected = true

	return c.Next()
}

// Logout Handler
func HandleLogout(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "error",
	}

	session, err := store.Get(c)
	if err != nil {
		fm["message"] = "logged out (no session)"

		return flash.WithError(c, fm).Redirect("/login")
	}

	err = session.Destroy()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/login")
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "You have successfully logged out!!",
	}

	c.Locals(FROM_PROTECTED, false)
	// fromProtected = false

	return flash.WithSuccess(c, fm).Redirect("/login")
}
