package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateOrderStatusHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")
		var req struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		_, err := db.Exec(context.Background(),
			`UPDATE orders SET status = $1 WHERE id = $2`,
			req.Status, orderID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Статус заказа обновлён"})
	}
}
