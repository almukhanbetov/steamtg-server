package handlers

import (
	"context"
	"net/http"
	"steamtg/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCategoriesHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(), `SELECT id, name, image FROM categories`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при загрузке категорий"})
			return
		}
		defer rows.Close()

		var categories []models.Category
		for rows.Next() {
			var cat models.Category
			if err := rows.Scan(&cat.ID, &cat.Name, &cat.Image); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при парсинге категорий"})
				return
			}
			categories = append(categories, cat)
		}

		c.JSON(http.StatusOK, categories)
	}
}
