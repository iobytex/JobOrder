package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string
	Product []Product `gorm:"constraint:OnDelete:CASCADE;"`
}


func (Category) TableName() string {
	return "category"
}
