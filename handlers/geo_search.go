package handlers

import (
	"context"
	"net/http"
	"log"
	

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NearestDriverRequest struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type NearestDriverResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Lon        float64 `json:"lon"`
	Lat        float64 `json:"lat"`
	Distance_m float64 `json:"distance_m"`
}

func GeoSearchNearestDriverHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req NearestDriverRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		log.Println("Поиск ближайшего водителя к:", req.Lon, req.Lat)

		query := `
			SELECT id, name, ST_X(location), ST_Y(location),
				ST_Distance(location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography)
			FROM drivers
			ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)
			LIMIT 1
		`

		var driver NearestDriverResponse

		err := db.QueryRow(context.Background(), query, req.Lon, req.Lat).
			Scan(&driver.ID, &driver.Name, &driver.Lon, &driver.Lat, &driver.Distance_m)

		if err != nil {
			log.Println("Ошибка поиска водителя:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске водителя"})
			return
		}

		c.JSON(http.StatusOK, driver)
	}
}
