package joborder

import (
	"github.com/gin-gonic/gin"
)

type Delivery interface {

	CreateUser() gin.HandlerFunc

	AddCategories() gin.HandlerFunc
	GetCategories() gin.HandlerFunc
	UpdateCategory() gin.HandlerFunc
	DeleteCategories() gin.HandlerFunc

	AddProducts() gin.HandlerFunc
	GetProducts() gin.HandlerFunc
	UpdateProduct() gin.HandlerFunc
	DeleteProducts() gin.HandlerFunc

	SetOrder() gin.HandlerFunc
	AddOrderItems() gin.HandlerFunc
	GetOrderItems() gin.HandlerFunc
}
