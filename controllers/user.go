package controllers

import (
	"errors"
	"time"
	//"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Permission string `json:"permission"`
}

func (t *User) GetAllUsers() ([]User, error) {
	// Sleep to add some delay in API response
	time.Sleep(time.Millisecond * 1500)

	var records []User

	PgDBConn.Find(&records)

	return records, nil
}

func (t *User) GetUserbySKU() (User, error) {

	query := `SELECT ID, Email, Password, Username, Permission FROM Users
		WHERE ID=t.SKU`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var recoveredUser User
	err = stmt.QueryRow(
		t.ID,
	).Scan(
		&recoveredUser.ID,
		&recoveredUser.Email,
		&recoveredUser.Username,
		&recoveredUser.Password,
		&recoveredUser.Permission,
	if err != nil {
		return User{}, err
	}

	return recoveredUser, nil
}

func (t *User) CreateUser() (User, error) {

	query := `INSERT INTO Users (ID, Email, Password, Username, Permission FROM Users)
		VALUES(?, ?, ?) RETURNING *`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var newUser User
	err = stmt.QueryRow(
		t.ID,
		t.Email,
		t.Username,
		t.Password,
		t.Permission,
	).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Username,
		&newUser.Password,
		&newUser.Permission,
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

	var updatedUser User
	err = stmt.QueryRow(
		t.ID,
		t.Email,
		t.Username,
		t.Password,
		t.Permission,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.Username,
		&updatedUser.Password,
		&updatedUser.Permission,

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
	//	//(	//( ID,Email,Password,Username,Permission ))
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


//****************  Table Load ***********************
func (t *User) UserLoad() error {

	// Load CSV data
	csvFile, err := os.Open("data/csv/data_login.csv") // Replace with your CSV file path
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read all rows from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV test file: %v", err)
		return err
	}

	// Loop through the records and create Category
	for _, record := range records {

		iID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails
		}
		iPermission, err := strconv.Atoi(record[5])
		if err != nil {
			log.Printf("Failed to convert to int - Inventory: %v", err)
			continue // Skip this record if conversion fails
		}

		test := &controllers.Users{
			ID:         iID,
			email:      record[1],
			UserName:   record[3],
			Password:   record[4],
			Permission: iPermission,
		}

		// Save item to the database
		err = db.Create(&test).Error
		if err != nil {
			log.Printf("Failed to insert test record: %v", err)
		}
	}
}