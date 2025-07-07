package handlers

import (
	"context"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateClientHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name       string  `json:"name"`
			Phone      string  `json:"phone"`
			Plate      string  `json:"plate"`
			CategoryID int     `json:"category_id"`
			Lon        float64 `json:"lon"`
			Lat        float64 `json:"lat"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		// ✅ UPSERT (insert or update on conflict)
		_, err := db.Exec(context.Background(),
			`INSERT INTO clients (name, phone, plate_number, category_id, location)
			 VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326))
			 ON CONFLICT (phone)
			 DO UPDATE SET
				name = EXCLUDED.name,
				plate_number = EXCLUDED.plate_number,
				category_id = EXCLUDED.category_id,
				location = EXCLUDED.location`,
			req.Name, req.Phone, req.Plate, req.CategoryID, req.Lon, req.Lat)

		if err != nil {
			log.Println("Ошибка UPSERT клиента:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании/обновлении клиента"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Клиент создан или обновлён"})
	}
}
