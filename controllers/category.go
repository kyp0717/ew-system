package controllers

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Category struct {
	Category                   string          `json:"Category" gorm:"primaryKey;size:20"`
	Category_Description       string          `json:"Category_Description" gorm:"null;column:Category_Description;size:255"`
	Reserved_Percentage_BySale decimal.Decimal `json:"Reserved_Percentage_BySale" gorm:"column:Reserved_Percentage_BySale;type:decimal(8,3);null"`
	Report_Out_Max             int             `json:"Report_Out_Max" gorm:"null;column:Report_Out_Max"`
}

func (t *Category) GetAllCategory() ([]Category, error) {
	// Sleep to add some delay in API response
	time.Sleep(time.Millisecond * 1500)

	var records []Category

	PgDBConn.Find(&records)

	return records, nil
}

func (t *Category) GetCategorybySKU() (Category, error) {
	
	query := `SELECT Category ,Category_Description , Reserved_Percentage_BySale, Report_Out_Max  FROM Category
		WHERE SKU=t.Category`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return Category{}, err
	}

	defer stmt.Close()


	var recoveredCategory Category
	err = stmt.QueryRow(
		t.Category,
	).Scan(
		&recoveredCategory.Category,
		&recoveredCategory.Category_Description,
		&recoveredCategory.Reserved_Percentage_BySale,
		&recoveredCategory.Report_Out_Max,
	)
	if err != nil {
		return Category{}, err
	}

	return recoveredCategory, nil
}

func (t *Category) CreateCategory() (Category, error) {

	query := `INSERT INTO Category (Category ,Category_Description , Reserved_Percentage_BySale, Report_Out_Max)
		VALUES(?, ?, ?) RETURNING *`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return Category{}, err
	}

	defer stmt.Close()

	var newCategory Category
	err = stmt.QueryRow(
		t.Category,
		t.Category_Description,
		t.Report_Out_Max,
	).Scan(
		&newCategory.Category,
		&newCategory.Category_Description,
		&newCategory.Reserved_Percentage_BySale,
		&newCategory.Report_Out_Max,
	)
	if err != nil {
		return Category{}, err
	}

	/* if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("error: an affected row was expected")
	} */

	return newCategory, nil
}
func (t *Category) UpdateCategory() (Category, error) {

	query := `UPDATE Category SET title = ?,  Category_Description = ?, Reserved_Percentage_BySale = ? , Report_Out_Max = ?
		WHERE Category=? `

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return Category{}, err
	}

	defer stmt.Close()

	var updatedCategory Category
	err = stmt.QueryRow(
		t.Category,
		t.Category_Description,
		t.Reserved_Percentage_BySale,
		t.Report_Out_Max,
	).Scan(
		&updatedCategory.Category,
		&updatedCategory.Category_Description,
		&updatedCategory.Reserved_Percentage_BySale,
		&updatedCategory.Report_Out_Max,
	)
	if err != nil {
		return Category{}, err
	}

	return updatedCategory, nil
}

func (t *Category) DeleteCategory() error {

	query := `DELETE FROM Category
		WHERE Category = ?`

	stmt, err := PgDBConn.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(t.Category, t.Category_Description)
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
