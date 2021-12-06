package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *server) MapRoutes()  {

	s.engine.HandleMethodNotAllowed = true

	s.engine.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound,gin.H{
			"code": "PAGE_NOT_FOUND",
			"message": "Page not found 🤡",
		})
	})

	s.engine.NoMethod(func(context *gin.Context) {
		context.JSON(http.StatusMethodNotAllowed,gin.H{
			"code": "PAGE_NOT_FOUND",
			"message": "Page not found 🤡",
		})
	})
}
