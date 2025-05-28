package main

import (
	"CRUD_BOOKS/DB"
	"CRUD_BOOKS/middleware"
	"CRUD_BOOKS/router"

	"github.com/gin-gonic/gin"
)

func main() {
	dbpool := db.InitDBPool()
	defer dbpool.Close()

	r := gin.Default()
	r.Use(middleware.Logger())
	router.RegisterRoutes(r)

	r.Run(":8080")
}
