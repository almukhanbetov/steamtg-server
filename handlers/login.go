package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"context"
)

func DriverLoginHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Phone string `json:"phone"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		var driverID int
		var name string
		err := db.QueryRow(context.Background(),
			`SELECT id, name FROM drivers WHERE phone = $1`,
			req.Phone).Scan(&driverID, &name)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный телефон"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    driverID,
			"name":  name,
			"phone": req.Phone,
		})
	}
}
