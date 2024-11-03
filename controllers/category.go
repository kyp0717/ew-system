package controllers

import (
	"github.com/shopspring/decimal"
)

type Category struct {
	Category                   string          `json:"Category" gorm:"primaryKey;size:20"`
	Category_Description       string          `json:"Category_Description" gorm:"null;column:Category_Description;size:255"`
	Reserved_Percentage_BySale decimal.Decimal `json:"Reserved_Percentage_BySale" gorm:"column:Reserved_Percentage_BySale;type:decimal(8,3);null"`
	Report_Out_Max             int             `json:"Report_Out_Max" gorm:"null;column:Report_Out_Max"`
}
