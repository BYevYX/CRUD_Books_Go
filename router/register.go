package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/book")
	{
		api.GET("/:id", GetBook)
		api.POST("/", CreateBook)
		api.DELETE("/:id", DeleteBook)
		api.PUT("/:id", UpdateBook)

		api.GET("/all", GetAllBooks)
	}
}
