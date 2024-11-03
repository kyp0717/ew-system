package controllers

import (
	"errors"
	"time"

	//"golang.org/x/crypto/bcrypt"

	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/emarifer/gofiber-templ-htmx/controllers"
	"gorm.io/gorm"
) 


type User struct {
	ID       	int `json:"id"`
	Email    	string `json:"email"`
	Password 	string `json:"password"`
	Username 	string `json:"username"`
	Permission	string `json:"permission"`
}

func (t *User) GetAllUsers() ([]User, error) {
	// Sleep to add some delay in API response
	time.Sleep(time.Millisecond * 1500)

	var records []User

	PgDBConn.Find(&records)

	return records, nil
}

func (t *User) GetUserbySKU() (User, error) {

	query := `SELECT SKU, UserName, UPC, Type, Category, Description, Inventory , Cost FROM Users
		WHERE SKU=t.SKU`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var recoveredUser User
	err = stmt.QueryRow(
		t.SKU,
	).Scan(
		&recoveredUser.SKU,
		&recoveredUser.UserName,
		&recoveredUser.UPC,
		&recoveredUser.Category,
		&recoveredUser.Description,
		&recoveredUser.Inventory,
		&recoveredUser.Cost,
	)
	if err != nil {
		return User{}, err
	}

	return recoveredUser, nil
}

func (t *User) CreateUser() (User, error) {

	query := `INSERT INTO Users (SKU, UserName, UPC, Type, Category, Description, Inventory , Cost)
		VALUES(?, ?, ?) RETURNING *`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var newUser User
	err = stmt.QueryRow(
		t.SKU,
		t.UPC,
		t.Description,
	).Scan(
		&newUser.SKU,
		&newUser.UserName,
		&newUser.UPC,
		&newUser.Description,
		&newUser.Inventory,
		&newUser.Cost,
	)
	if err != nil {
		return User{}, err
	}

	/* if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("error: an affected row was expected")
	} */

	return newUser, nil
}
func (t *User) UpdateUser() (User, error) {

	query := `UPDATE Users SET title = ?,  description = ?, status = ?
		WHERE created_by = ? AND id=? RETURNING id, title, description, status`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()
//(SKU, UserName, UPC, Type, Category, Description, Inventory , Cost)
	var updatedUser User
	err = stmt.QueryRow(
		t.SKU,
		t.UserName,
		t.UPC,
		t.Category,
		t.Cost,
	).Scan(
		&updatedUser.SKU,
		&updatedUser.UserName,
		&updatedUser.UPC,
		&updatedUser.Category,
		&updatedUser.Inventory,
		&updatedUser.Cost,
	)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (t *User) DeleteUser() error {

	query := `DELETE FROM Users
		WHERE created_by = ? AND id=?`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
//(SKU, UserName, UPC, Type, Category, Description, Inventory , Cost)
	result, err := stmt.Exec(t.UserName, t.SKU)
	if err != nil {
		return err
	}

	if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}

func ConvertDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return dt.In(loc).Format(time.RFC822Z)
}
