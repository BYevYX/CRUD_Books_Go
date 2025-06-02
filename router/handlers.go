package router

import (
	"CRUD_BOOKS/DB"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PagesCount     int    `json:"pages_count"`
	PublicationDate string `json:"publication_date"`
	Author         string `json:"author"`
}

const (
	getBookQuery = `
		SELECT books.id, books.name, books.pages_count, books.publication_date, authors.name as author 
		FROM books 
		JOIN authors ON books.author_id = authors.id 
		WHERE books.id = $1
	`
	
	getAllBooksQuery = `
		SELECT books.id, books.name, books.pages_count, books.publication_date, authors.name as author 
		FROM books 
		JOIN authors ON books.author_id = authors.id
		ORDER BY books.id
	`
)

func GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be positive integer"})
		return
	}

	var book Book
	err = db.DBpool.QueryRow(context.Background(), getBookQuery, id).Scan(
		&book.ID, 
		&book.Name, 
		&book.PagesCount, 
		&book.PublicationDate, 
		&book.Author,
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
			&book.Author,
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
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
        return
    }

    result, err := db.DBpool.Exec(context.Background(),
        "DELETE FROM books WHERE id = $1",
        id,
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
        "id":      id,
    })
}

type CreateBookRequest struct {
    Name           string `json:"name" binding:"required"`
    PagesCount     int    `json:"pages_count" binding:"required"`
    PublicationDate string `json:"publication_date"`
    AuthorID       int    `json:"author_id" binding:"required"`
}

func CreateBook(c *gin.Context) {
    var req CreateBookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    var bookID int
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

type UpdateBookRequest struct {
    Name           *string `json:"name"`
    PagesCount     *int    `json:"pages_count"`
    PublicationDate *string `json:"publication_date"`
    AuthorID       *int    `json:"author_id"`
}

func UpdateBook(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
        return
    }

    var req UpdateBookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    // Динамическое построение запроса
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
    if req.AuthorID != nil {
        query += fmt.Sprintf("author_id = $%d, ", paramCount)
        params = append(params, *req.AuthorID)
        paramCount++
    }

    if len(params) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
        return
    }

    // Удаляем последнюю запятую
    query = query[:len(query)-2]
    query += fmt.Sprintf(" WHERE id = $%d", paramCount)
    params = append(params, id)

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
        "id":      id,
    })
}

