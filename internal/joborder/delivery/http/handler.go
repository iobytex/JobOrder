package http

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"joborder/internal/joborder"
	"joborder/internal/model"
	"net/http"
)

type jobOrderHandler struct {
	group *gin.RouterGroup
	logger *zap.Logger
	service joborder.Service
	middleware *jwt.GinJWTMiddleware
}

func NewJobOrderHandler( group *gin.RouterGroup,service joborder.Service,logger *zap.Logger,middleware *jwt.GinJWTMiddleware) *jobOrderHandler {
	return &jobOrderHandler{
		group: group,
		logger: logger,
		service: service,
		middleware: middleware,
	}
}

func (joh *jobOrderHandler) CreateUser() gin.HandlerFunc  {
	return func(context *gin.Context) {
		var userRequest model.UserRequest
		if err := context.ShouldBind(&userRequest) ; err != nil {
			joh.logger.Sugar().Error(err.Error())
			joh.logger.Sugar().Info(userRequest)
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		err := joh.service.CreateUser(context, &userRequest)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("%s has been added has an employee",userRequest.Name)})
		return

	}
}

//func (joh *jobOrderHandler) CheckUserByPhoneNumber() gin.HandlerFunc  {
//	return func(context *gin.Context) {
//
//		var login model.Login
//
//		if err := context.ShouldBind(&login) ; err != nil {
//			joh.logger.Sugar().Error(err.Error())
//			joh.logger.Sugar().Info(login)
//			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
//			return
//		}
//
//		user ,err := joh.service.CheckUserByPhoneNumber(context, &login)
//		if err != nil {
//			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
//			return
//		}
//
//		context.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("%s has been added has an employee",userRequest.Name)})
//		return
//
//	}
//}

func (joh *jobOrderHandler) AddCategories() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			Name []string `json:"categories"`
		}

		var req request

		err := context.ShouldBind(&req)
		if err != nil {
			joh.logger.Sugar().Error(err.Error())
			joh.logger.Sugar().Info(req)
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		updateCategoryErr := joh.service.AddCategories(context, &req.Name)
		if updateCategoryErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": updateCategoryErr.Error()})
			return
		}

		context.JSON(http.StatusOK,gin.H{
			"message":"Category's name has been Added",
		})
		return

	}
}

func (joh *jobOrderHandler) GetCategories() gin.HandlerFunc  {
	return func(context *gin.Context) {
		categories, err := joh.service.GetCategories(context)
		if err != nil {
			joh.logger.Sugar().Error(err.Error())
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusNotAcceptable, gin.H{"categories": categories})
		return
	}
}

func (joh *jobOrderHandler) UpdateCategory() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			CategoryId uint `form:"category_id" binding:"required"`
			NewCategoryName string `form:"name" binding:"required"`
		}

		var req request

		err := context.ShouldBind(&req)
		if err != nil {
			joh.logger.Sugar().Error(err.Error())
			joh.logger.Sugar().Info(req)
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		updateCategoryErr := joh.service.UpdateCategory(context, req.CategoryId, req.NewCategoryName)
		if updateCategoryErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": updateCategoryErr.Error()})
			return
		}

		context.JSON(http.StatusOK,gin.H{
			"message":"Category's name has been Updated",
		})
		return
	}
}

func (joh *jobOrderHandler) DeleteCategories() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			CategoryId []uint `form:"category_id" binding:"required"`
		}

		var req request

		err := context.ShouldBind(&req)
		if err != nil {
			joh.logger.Sugar().Error(err.Error())
			joh.logger.Sugar().Info(req)
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		categoriesErr := joh.service.DeleteCategories(context,&req.CategoryId)
		if categoriesErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": categoriesErr.Error()})
			return
		}

		context.JSON(http.StatusOK,gin.H{
			"message":"Category's name has been Updated",
		})
		return

	}
}

func (joh *jobOrderHandler) AddProducts() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			Product []model.ProductRequest `json:"products"`
		}

		var req request

		if err := context.ShouldBindJSON(&req); err != nil {
			joh.logger.Sugar().Error(err.Error())
			joh.logger.Sugar().Info(req)
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		err := joh.service.AddProducts(context,&req.Product)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusCreated,gin.H{
			"message": "All Product has been added to the specified Category",
		})
		return
	}
}

func (joh *jobOrderHandler) GetProducts() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			categoryId uint `form:"category_id" binding:"required"`
		}
		var req request
		err := context.ShouldBind(req)
		if err != nil {
			joh.logger.Sugar().Error(err.Error())
			joh.logger.Sugar().Info(req)
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		products, err := joh.service.GetProducts(context, req.categoryId)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusCreated,gin.H{
			"product": products,
		})
		return
	}
}

func (joh *jobOrderHandler) UpdateProduct() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			CategoryId uint `json:"category_id"`
			ProductId uint `json:"product_id"`
			Product map[string]interface{} `json:"product"`
		}

		var req request

		err := context.ShouldBindJSON(&req)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		updateErr := joh.service.UpdateProduct(context, req.CategoryId, req.ProductId, &req.Product)
		if updateErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": updateErr.Error()})
			return
		}

		context.JSON(http.StatusOK,gin.H{
			"message": "Product has been update",
		})
		return
	}
}

func (joh *jobOrderHandler) DeleteProducts() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			ProductId []uint `json:"product_id"`
		}

		var req request

		err := context.ShouldBindJSON(&req)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		deletedProductErr := joh.service.DeleteProducts(context, &req.ProductId)
		if deletedProductErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": deletedProductErr.Error()})
			return
		}

		context.JSON(http.StatusOK,gin.H{
			"message": "All Selected Product has been deleted",
		})
		return

	}
}


func (joh *jobOrderHandler) SetOrder() gin.HandlerFunc  {
	return func(context *gin.Context) {
		type request struct {
			CreatedBy uint `json:"created_by"`
		}

		var req request

		err := context.ShouldBindJSON(&req)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		setOrderErr := joh.service.SetOrder(context, req.CreatedBy)
		if setOrderErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": setOrderErr.Error()})
			return
		}

		context.JSON(http.StatusCreated,gin.H{
			"message": "Order has been set, Add items to the cart",
		})
		return
	}
}

func (joh *jobOrderHandler) AddOrderItems() gin.HandlerFunc {
	return func(context *gin.Context) {
		type request struct {
			OrderItem []model.OrderItemRequest `json:"order_item"`
		}

		var req request

		err := context.ShouldBindJSON(&req)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		AddOrderItemsErr := joh.service.AddOrderItems(context, &req.OrderItem)
		if AddOrderItemsErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": AddOrderItemsErr.Error()})
			return
		}

		context.JSON(http.StatusCreated,gin.H{
			"message": "Added to cart",
		})
		return
	}
}
func (joh *jobOrderHandler) GetOrderItems() gin.HandlerFunc {
	return func(context *gin.Context) {
		type request struct {
			OrderId uint `json:"order_id"`
		}

		var req request

		err := context.ShouldBindJSON(&req)
		if err != nil {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
			return
		}

		getOrderItems,getOrderItemsErr := joh.service.GetOrderItems(context, req.OrderId)
		if getOrderItemsErr != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": getOrderItemsErr.Error()})
			return
		}

		context.JSON(http.StatusCreated,gin.H{
			"order": getOrderItems,
		})
		return

	}
}
