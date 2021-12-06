package service

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"joborder/internal/joborder"
	"joborder/internal/model"
)

type serviceImpl struct {
	logger *zap.Logger
	repo joborder.Repository
}

func NewServiceImpl(logger *zap.Logger,repo joborder.Repository) *serviceImpl {
	return &serviceImpl{
		logger: logger,
		repo: repo,
	}
}

func (service *serviceImpl) CreateUser(ctx context.Context, request *model.UserRequest) error  {
	err := request.PrepareCreate()
	if err != nil {
		return err
	}
	userErr := service.repo.CreateUser(ctx, request)
	if userErr != nil {
		return userErr
	}


	return nil
}

func (service *serviceImpl) CheckUserByPhoneNumber(ctx context.Context,request *model.Login) (*model.User,error)  {
	userByPhoneNumber, err := service.repo.CheckUserByPhoneNumber(ctx,request.PhoneNumber)
	if err != nil {
		return nil, err
	}

	passCodeErr := userByPhoneNumber.CompareHashPassCode([]byte(request.Passcode))
	if passCodeErr != nil {
		return nil, errors.New("Incorrect PassCode")
	}

	userByPhoneNumber.SanitizePassword()

	return userByPhoneNumber,nil
}

func (service *serviceImpl) AddCategories(ctx context.Context,name *[]string) error  {
	if len(*name) == 0 {
		return errors.New("Empty category")
	}
	for _,value := range *name{
		if len(value) == 0 {
			return errors.New("One of your category value is empty")
		}
	}

	err := service.repo.AddCategories(ctx, name)
	if err != nil {
		return err
	}

	return  nil
}

func (service *serviceImpl) GetCategories(ctx context.Context) (*[]model.Category,error)  {
	getCategories, err := service.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return getCategories, nil
}

func (service *serviceImpl) UpdateCategory(ctx context.Context,categoryId uint, name string) error  {
	err := service.repo.UpdateCategory(ctx, categoryId, name)
	if err != nil {
		return  err
	}
	return  nil
}

func (service *serviceImpl) DeleteCategories(ctx context.Context,categoryId *[]uint) error  {
	err := service.repo.DeleteCategories(ctx, categoryId)
	if err != nil {
		return  err
	}

	return nil

}


func (service *serviceImpl) AddProducts(ctx context.Context,request *[]model.ProductRequest) error  {
	var product []model.Product
	for _,value := range *request {
		if value.Custom == 1 {
			product = append(product,model.Product{
				CategoryID: value.CategoryID,
				Custom: true,
				BasePrice: value.BasePrice,
				Stock: &model.Stock{
					Quantity: value.Quantity,
				},
			})
		}else{
			product = append(product,model.Product{
				CategoryID: value.CategoryID,
				Custom: false,
				BasePrice: value.BasePrice,
				Stock: &model.Stock{
					Quantity: value.Quantity,
				},
			})
		}
	}

	err := service.repo.AddProducts(ctx, &product)
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceImpl) GetProducts(ctx context.Context,categoryId uint) (*[]model.Product,error)  {
	products, err := service.repo.GetProducts(ctx,categoryId)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (service *serviceImpl) UpdateProduct(ctx context.Context,categoryId uint,productId uint,product *map[string]interface{}) error {
	for i,value := range *product{
		if value == "" {
			delete(*product,i)
		}
	}
	err := service.repo.UpdateProduct(ctx, categoryId, productId, product)
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceImpl) DeleteProducts(ctx context.Context,productId *[]uint) error  {
	if len(*productId) == 0 {
		return errors.New("Empty Product")
	}
	for _,value := range *productId{
		if value == 0 {
			return errors.New("One of your product value is empty")
		}
	}

	err := service.repo.DeleteProducts(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}



func (service *serviceImpl) SetOrder(ctx context.Context,createdBy uint) error {
	err := service.repo.SetOrder(ctx, createdBy)
	if err != nil {
		return err
	}
	return nil
}

func (service *serviceImpl) AddOrderItems(ctx context.Context, item *[]model.OrderItemRequest) error {
	orderItem := make([]model.OrderItem,0)

	for _, requestValue := range *item {
		orderItem = append(orderItem,model.OrderItem{
			OrderID: requestValue.OrderID,
			ProductID: requestValue.ProductID,
			Measurement: &model.Measurement{
				Width: requestValue.Width,
				Height: requestValue.Height,
			},
		})
	}

	err := service.repo.AddOrderItems(ctx, &orderItem)
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceImpl) GetOrderItems(ctx context.Context,orderId uint) (*[]model.OrderItem,error){

	getOrderItems, err := service.repo.GetOrderItems(ctx,orderId)
	if err != nil {
		return nil, err
	}

	return getOrderItems, nil
}

