package handlers

import (
	"context"
	"net/http"
	"steamtg/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDriversHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := db.Query(ctx, `
			SELECT id, name, iin, photo, ST_X(location), ST_Y(location), car_id
			FROM drivers
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var drivers []models.Driver
		for rows.Next() {
			var d models.Driver
			if err := rows.Scan(&d.ID, &d.Name, &d.IIN, &d.Photo, &d.Lon, &d.Lat, &d.CarID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			drivers = append(drivers, d)
		}

		c.JSON(http.StatusOK, drivers)
	}
}
