package controllers

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement:true;type:bigserial" json:"id"`
	Email    string `gorm:"column:email;type:varchar(50);unique;not null" json:"email"`
	Username string `gorm:"column:username;type:varchar(20);unique;not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(100);not null" json:"password"`
}

func (t *User) GetAllUsers() ([]User, error) {
	// Sleep to add some delay in API response
	time.Sleep(time.Millisecond * 1500)
	var records []User

	PgDBConn.Find(&records)

	return records, nil
}

func LoadUserTable() error {

	// Load CSV data
	csvFile, err := os.Open("./controllers/data/users.csv") // Replace with your CSV file path
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

		test := &User{
			Email:    record[0],
			Username: record[1],
			Password: record[2],
		}
		// Save item to the database
		err = PgDBConn.Create(&test).Error
		if err != nil {
			log.Printf("Failed to insert user record: %v", err)
		}
	}
	return nil
}

func CreateUser(user *User) error {
	//Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Use GORM to create the user in the database
	if err := PgDBConn.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserById(id string) (User, error) {
	var user User

	err := PgDBConn.First(&user, "id = ?", id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func CheckEmail(email string) (User, error) {
	var user User

	err := PgDBConn.Where("email = ?", email).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
