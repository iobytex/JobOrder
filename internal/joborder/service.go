package joborder

import (
	"context"
	"joborder/internal/model"
)

type Service interface {

	CreateUser(ctx context.Context, request *model.UserRequest) error
	CheckUserByPhoneNumber(ctx context.Context,request *model.Login) (*model.User,error)

	AddCategories(ctx context.Context,name *[]string) error
	GetCategories(ctx context.Context) (*[]model.Category,error)
	UpdateCategory(ctx context.Context,categoryId uint, name string) error
	DeleteCategories(ctx context.Context,categoryId *[]uint) error

	AddProducts(ctx context.Context,request *[]model.ProductRequest) error
	GetProducts(ctx context.Context,categoryId uint) (*[]model.Product,error)
	UpdateProduct(ctx context.Context,categoryId uint,productId uint,product *map[string]interface{}) error
	DeleteProducts(ctx context.Context,productId *[]uint) error

	SetOrder(ctx context.Context,createdBy uint) error
	AddOrderItems(ctx context.Context, item *[]model.OrderItemRequest) error
	GetOrderItems(ctx context.Context,orderId uint) (*[]model.OrderItem,error)
}
