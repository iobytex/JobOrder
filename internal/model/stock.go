package model

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	ProductID uint
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
	Quantity uint
}

func (Stock) TableName() string {
	return "stocks"
}

