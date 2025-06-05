package router

import (
	db "CRUD_BOOKS/DB"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PagesCount      int    `json:"pages_count"`
	PublicationDate string `json:"publication_date"`
	AuthorID        string `json:"author_id"`
}

type CreateBookRequest struct {
	Name            string `json:"name" binding:"required"`
	PagesCount      int    `json:"pages_count" binding:"required"`
	PublicationDate string `json:"publication_date" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type UpdateBookRequest struct {
	Name            *string `json:"name"`
	PagesCount      *int    `json:"pages_count"`
	PublicationDate *string `json:"publication_date"`
}

const (
	getBookQuery = `
		SELECT books.id, books.name, books.pages_count, books.publication_date::text, books.author_id 
		FROM books 
		WHERE books.id = $1
	`

	getAllBooksQuery = `
		SELECT books.id, books.name, books.pages_count, books.publication_date::text, books.author_id 
		FROM books 
		ORDER BY books.id
	`
)

func CreateBook(c *gin.Context) {
	fmt.Println("Create Book")
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var bookID string
	err := db.DBpool.QueryRow(context.Background(),
		`INSERT INTO books (name, pages_count, publication_date, author_id)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
		req.Name, req.PagesCount, req.PublicationDate, req.AuthorID,
	).Scan(&bookID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"id":      bookID,
	})
}

func UpdateBook(c *gin.Context) {
	bookUuid := c.Param("id")

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	query := "UPDATE books SET "
	params := []interface{}{}
	paramCount := 1

	if req.Name != nil {
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, *req.Name)
		paramCount++
	}
	if req.PagesCount != nil {
		query += fmt.Sprintf("pages_count = $%d, ", paramCount)
		params = append(params, *req.PagesCount)
		paramCount++
	}
	if req.PublicationDate != nil {
		query += fmt.Sprintf("publication_date = $%d, ", paramCount)
		params = append(params, *req.PublicationDate)
		paramCount++
	}

	if len(params) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, bookUuid)

	result, err := db.DBpool.Exec(context.Background(), query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
		"id":      bookUuid,
	})
}

func GetBook(c *gin.Context) {
	bookUuid := c.Param("id")

	var book Book
	err := db.DBpool.QueryRow(context.Background(), getBookQuery, bookUuid).Scan(
		&book.ID,
		&book.Name,
		&book.PagesCount,
		&book.PublicationDate,
		&book.AuthorID,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

func GetAllBooks(c *gin.Context) {
	rows, err := db.DBpool.Query(context.Background(), getAllBooksQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.Name,
			&book.PagesCount,
			&book.PublicationDate,
			&book.AuthorID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error: " + err.Error()})
			return
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rows iteration error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(books),
		"books": books,
	})
}

func DeleteBook(c *gin.Context) {
	bookUuid := c.Param("id")

	result, err := db.DBpool.Exec(context.Background(),
		"DELETE FROM books WHERE id = $1",
		bookUuid,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
		"id":      bookUuid,
	})
}
