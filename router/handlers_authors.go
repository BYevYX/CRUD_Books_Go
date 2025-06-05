package router

import (
	db "CRUD_BOOKS/DB"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Author struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	DeathDate string `json:"death_date"`
}

type RegisterAuthorRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateAuthorRequest struct {
	Name      *string `json:"name"`
	BirthDate *string `json:"birth_date"`
	DeathDate *string `json:"death_date"`
}

const (
	getAuthorQuery = `
		SELECT authors.id, authors.name, authors.birthdate::text, authors.death_date::text 
		FROM authors 
		WHERE authors.id = $1
	`

	getAllAuthorsQuery = `
		SELECT authors.id, authors.name, authors.birthdate::text, authors.death_date::text
		FROM authors 
		ORDER BY authors.id
	`
)

func RegisterAuthor(c *gin.Context) {
	var req RegisterAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var AuthorID string
	err := db.DBpool.QueryRow(context.Background(),
		`INSERT INTO authors (name)
         VALUES ($1)
         RETURNING id
        `,
		req.Name,
	).Scan(&AuthorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register author: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Author registered successfully",
		"id":      AuthorID,
	})
}

func UpdateAuthor(c *gin.Context) {
	authorUuid := c.Param("id")

	var req UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	query := "UPDATE authors SET "
	params := []interface{}{}
	paramCount := 1

	if req.Name != nil {
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, *req.Name)
		paramCount++
	}
	if req.BirthDate != nil {
		query += fmt.Sprintf("birthdate = $%d, ", paramCount)
		params = append(params, *req.BirthDate)
		paramCount++
	}
	if req.DeathDate != nil {
		query += fmt.Sprintf("death_date = $%d, ", paramCount)
		params = append(params, *req.DeathDate)
		paramCount++
	}

	if len(params) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, authorUuid)

	result, err := db.DBpool.Exec(context.Background(), query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "author updated successfully",
		"id":      authorUuid,
	})
}

func GetAuthor(c *gin.Context) {
	authorUuid := c.Param("id")
	var author Author
	err := db.DBpool.QueryRow(context.Background(), getAuthorQuery, authorUuid).Scan(
		&author.ID,
		&author.Name,
		&author.BirthDate,
		&author.DeathDate,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, author)
}

func GetAllAuthors(c *gin.Context) {
	rows, err := db.DBpool.Query(context.Background(), getAllAuthorsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.BirthDate,
			&author.DeathDate,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error: " + err.Error()})
			return
		}
		authors = append(authors, author)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rows iteration error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(authors),
		"authors": authors,
	})
}
