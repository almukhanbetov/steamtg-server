package handlers

import (
	"context"
	"net/http"
	"steamtg/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDriverHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Driver
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		_, err := db.Exec(context.Background(),
			`INSERT INTO drivers (name, iin, photo, location, car_id,phone)
			 VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6)`,
			req.Name, req.IIN, req.Photo, req.Lon, req.Lat, req.CarID, req.Phone)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании водителя"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Водитель создан"})
	}
}
