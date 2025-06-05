package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	book_api := r.Group("/api/book")
	{
		book_api.GET("/:id", GetBook)
		book_api.POST("/", CreateBook)
		book_api.DELETE("/:id", DeleteBook)
		book_api.PUT("/:id", UpdateBook)
		book_api.GET("/all", GetAllBooks)
	}
	author_api := r.Group("/api/author")
	{
		author_api.GET("/:id", GetAuthor)
		author_api.POST("/", RegisterAuthor)
		author_api.PUT("/:id", UpdateAuthor)
		author_api.GET("/all", GetAllAuthors)
	}
}
