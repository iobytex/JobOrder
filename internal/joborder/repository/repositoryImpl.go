package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"joborder/internal/model"
	"joborder/pkg/gorm_errors"
)

type repositoryImpl struct {
	db *gorm.DB
	logger *zap.Logger
}


func NewRepositoryImpl(db *gorm.DB,logger *zap.Logger) *repositoryImpl{
	return &repositoryImpl{
		db: db,
		logger: logger,
	}
}


func (repository *repositoryImpl) CreateUser(ctx context.Context, request *model.UserRequest) error  {


	user :=  model.User{
		Name: request.Name,
		PhoneNumber: request.PhoneNumber,
		Passcode: request.Passcode,
		Role: request.Role, //Casbin will be user later
	}
	if result := repository.db.WithContext(ctx).Create(&user); result.Error != nil {
		return errors.Wrap(result.Error,"")
	}

	return nil
}

func (repository *repositoryImpl) CheckUserByPhoneNumber(ctx context.Context,phoneNumber string) (*model.User,error)  {

	var user model.User

	if result := repository.db.WithContext(ctx).Where("phone_no = ?",phoneNumber).First(&user) ; result.Error != nil {
		return nil,gorm_errors.GormError(result.Error)
	}

	return &user,nil
}

func (repository *repositoryImpl) AddCategories(ctx context.Context,name *[]string) error  {

	category := make([]model.Category,0)

	for _,value := range  *name{
		category = append(category,model.Category{Name: value})
	}

	if result :=  repository.db.WithContext(ctx).CreateInBatches(category,len(category)); result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return  nil
}

func (repository *repositoryImpl) GetCategories(ctx context.Context) (*[]model.Category,error)  {
	var categories []model.Category

	if result := repository.db.WithContext(ctx).Find(&categories); result.Error != nil {
		return nil,gorm_errors.GormError(result.Error)
	}

	return &categories,nil
}

func (repository *repositoryImpl) UpdateCategory(ctx context.Context,categoryId uint, name string) error  {

	if result := repository.db.WithContext(ctx).Model(&model.Category{}).Where("category_id = ?",categoryId).Update("name",name); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}
	return nil
}

func (repository *repositoryImpl) DeleteCategories(ctx context.Context,categoryId *[]uint) error{

	if result := repository.db.WithContext(ctx).Delete(&[]model.Category{},&categoryId); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}
	return nil
}


func (repository *repositoryImpl) AddProducts(ctx context.Context,product *[]model.Product) error {
	if result := repository.db.WithContext(ctx).Model(&model.Product{}).CreateInBatches(&product,len(*product)); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}
	return nil
}


func (repository *repositoryImpl) GetProducts(ctx context.Context,categoryId uint) (*[]model.Product,error) {
	var products []model.Product
	if result := repository.db.WithContext(ctx).Where("category_id <> ?",categoryId).Find(&products); result.Error != nil {
		return nil,gorm_errors.GormError(result.Error)
	}
	return &products,nil
}

func (repository *repositoryImpl) UpdateProduct(ctx context.Context,categoryId uint,productId uint,product *map[string]interface{}) error {

	if result := repository.db.WithContext(ctx).Where("category_id = ? AND product_id = ?",categoryId,productId).Updates(product); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}

	return nil
}

func (repository *repositoryImpl) DeleteProducts(ctx context.Context,productId *[]uint) error {
	var product model.Product

	if result := repository.db.WithContext(ctx).Delete(&product,*productId); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}

	return nil
}

func (repository *repositoryImpl) SetOrder(ctx context.Context,createdBy uint) error  {

	order :=  model.Order{
		CreatedBy: createdBy,
	}
	if result := repository.db.WithContext(ctx).Select("created_by").Create(&order); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}

	return nil
}

func (repository *repositoryImpl) AddOrderItems(ctx context.Context,item *[]model.OrderItem) error {

	err := repository.db.Transaction(func(tx *gorm.DB) error {

		stockProductId := make([]uint, 0)
		for _, v := range *item {
			stockProductId = append(stockProductId, v.ProductID)
		}
		var stock []model.Stock
		// SELECT * FROM `stock` FOR UPDATE Locking
		if stockResult := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id IN ?", stockProductId).First(&stock); stockResult.Error != nil {
			return gorm_errors.GormError(stockResult.Error)
		}

		for _, data := range stock {
			var value uint
			for _, orderItem := range *item {
				if orderItem.ProductID == data.ProductID {
					value = orderItem.Quantity
					break
				}
			}
			if updateStockQuantityValueResult := tx.WithContext(ctx).Where("product_id = ?", data.ProductID).UpdateColumn("quantity", gorm.Expr("quantity - ?", value)); updateStockQuantityValueResult != nil {
				return gorm_errors.GormError(updateStockQuantityValueResult.Error)
			}
		}

		if result := tx.WithContext(ctx).Create(*item); result.Error != nil {
			return gorm_errors.GormError(result.Error)
		}

		return nil
	},&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	return nil
}

func (repository *repositoryImpl) GetOrderItems(ctx context.Context,orderId uint) (*[]model.OrderItem,error) {

	var orderItems []model.OrderItem

	if result := repository.db.WithContext(ctx).Where("order_id <> ?",orderId).Find(&orderItems); result.Error != nil {
		return nil,gorm_errors.GormError(result.Error)
	}

	return &orderItems,nil

}
