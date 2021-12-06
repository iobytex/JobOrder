package joborder

import (
	"context"
	"joborder/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, request *model.UserRequest) error
	CheckUserByPhoneNumber(ctx context.Context,phoneNumber string) (*model.User,error)

	AddCategories(ctx context.Context,name *[]string) error
	GetCategories(ctx context.Context) (*[]model.Category,error)
	UpdateCategory(ctx context.Context,categoryId uint, name string) error
	DeleteCategories(ctx context.Context,categoryId *[]uint) error

	//Stock will be added here
	AddProducts(ctx context.Context,product *[]model.Product) error
	GetProducts(ctx context.Context,categoryId uint) (*[]model.Product,error)
	UpdateProduct(ctx context.Context,categoryId uint,productId uint,product *map[string]interface{}) error
	DeleteProducts(ctx context.Context,productId *[]uint) error

	SetOrder(ctx context.Context,createdBy uint) error
	AddOrderItems(ctx context.Context,item *[]model.OrderItem) error
	GetOrderItems(ctx context.Context,orderId uint) (*[]model.OrderItem,error)
}
