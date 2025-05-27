package main

import (
	"CRUD_BOOKS/DB"
)

func main() {
	dbpool := db.InitDBPool()
	defer dbpool.Close()
}