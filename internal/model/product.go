package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryID uint
	Category Category `gorm:"foreignKey:CategoryID;references:ID"`
	OrderItem *OrderItem
	Stock *Stock `gorm:"constraint:OnDelete:CASCADE;"`
	Custom bool
	BasePrice decimal.Decimal
}

type ProductRequest struct {
	CategoryID uint `json:"category_id" binding:"required"`
	Custom int `json:"custom" binding:"required"`
	BasePrice decimal.Decimal `json:"base_price" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}

func (Product) TableName() string {
	return "products"
}