package server

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	handler "joborder/internal/joborder/delivery/http"
	"joborder/internal/joborder/repository"
	"joborder/internal/joborder/service"
	"joborder/internal/middleware"
	"os"
)

type server struct {
	engine *gin.Engine
	db *gorm.DB
	logger *zap.Logger
}

func NewServer(db *gorm.DB,logger *zap.Logger) *server {
	return &server{
		engine: gin.Default(),
		db: db,
		logger: logger,
	}
}

func (s *server) Start()  error {


	repoImpl := repository.NewRepositoryImpl(s.db, s.logger)
	serviceImpl := service.NewServiceImpl(s.logger, repoImpl)

	middleManager := middleware.NewMiddleWareManager(serviceImpl)
	gonicJwt,jwtError := jwt.New(middleManager.MiddleWareHandler())

	v1 := s.engine.Group("/v1")
	{
		jobOrderHandler := handler.NewJobOrderHandler(v1, serviceImpl , s.logger,gonicJwt)
		jobOrderHandler.JobOrderMapRoute()
	}

	s.MapRoutes()

	if jwtError != nil {
		s.logger.Sugar().Infof("%s",jwtError.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ginRunErr := make(chan error)
	go func() {
		err := s.engine.Run(":"+port)
		if err != nil {
			ginRunErr <- errors.Wrap(err,"")
		}
	}()


	if ginRunErr!=nil{
		msg := <- ginRunErr
		s.logger.Error(msg.Error())
		return msg
	}



	return nil
}
