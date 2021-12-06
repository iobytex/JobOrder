package model

import (
	"github.com/go-contrib/uuid"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID     uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy uint
	User User `gorm:"foreignKey:CreatedBy;references:ID"`
	OrderItem []OrderItem
}

type OrderRequest struct {
	CreatedBy uint `form:"created_by" binding:"required"`
}

func (Order) TableName() string {
	return "orders"
}



type OrderItem struct {
	gorm.Model
	OrderID uint
	Order Order `gorm:"foreignKey:OrderID;references:ID"`
	ProductID uint
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
	Measurement *Measurement
	Quantity uint
}

type OrderItemRequest struct {
	OrderID uint `form:"order_id" binding:"required"`
	ProductID uint `form:"product_id" binding:"required"`
	Width uint 	`form:"width" binding:"required"`
	Height uint `form:"height" binding:"required"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type Measurement struct {
	gorm.Model
	Width uint
	Height uint
	OrderItemID uint
	OrderItem OrderItem `gorm:"foreignKey:OrderItemID;references:ID"`
}

func (Measurement) TableName() string{
	return "measurement"
}


