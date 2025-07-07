package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	ID          int     `json:"id"`
	ClientName  string  `json:"client_name"`
	ClientPhone string  `json:"client_phone"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Status      string  `json:"status"`
	DriverID    int     `json:"driver_id"`
	DriverName  string  `json:"driver_name"`
	DriverLat   float64 `json:"driver_lat"`
	DriverLon   float64 `json:"driver_lon"`
}

// 🔹 Получение всех заказов без фильтра (Admin)
func GetAllOrdersHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(), `
			SELECT o.id, c.name, c.phone, ST_Y(c.location), ST_X(c.location), o.status, o.driver_id,
			       d.name, ST_Y(d.location), ST_X(d.location)
			FROM orders o
			JOIN clients c ON o.client_id = c.id
			JOIN drivers d ON o.driver_id = d.id
		`)
		if err != nil {
			log.Println("Ошибка db.Query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при загрузке заказов"})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var o Order
			if err := rows.Scan(
				&o.ID,
				&o.ClientName,
				&o.ClientPhone,
				&o.Lat,
				&o.Lon,
				&o.Status,
				&o.DriverID,
				&o.DriverName,
				&o.DriverLat,
				&o.DriverLon,
			); err != nil {
				log.Println("Ошибка rows.Scan:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при парсинге заказов"})
				return
			}
			orders = append(orders, o)
		}

		c.JSON(http.StatusOK, orders)
	}
}

// 🔹 Получение заказов по driver_id + опционально по статусу
func GetOrdersForDriver(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		driverID := c.Param("id")
		status := c.Query("status") // получаем статус из query-параметра

		log.Println("Получение заказов для driverID:", driverID, "со статусом:", status)

		query := `
			SELECT o.id, c.name, c.phone, ST_Y(c.location), ST_X(c.location), o.status,
			       d.id, d.name, ST_Y(d.location), ST_X(d.location)
			FROM orders o
			JOIN clients c ON o.client_id = c.id
			JOIN drivers d ON o.driver_id = d.id
			WHERE o.driver_id = $1
		`

		args := []interface{}{driverID}

		if status != "" {
			query += " AND o.status = $2"
			args = append(args, status)
		}

		rows, err := db.Query(context.Background(), query, args...)
		if err != nil {
			log.Println("Ошибка db.Query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при загрузке заказов"})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var o Order
			if err := rows.Scan(
				&o.ID,
				&o.ClientName,
				&o.ClientPhone,
				&o.Lat,
				&o.Lon,
				&o.Status,
				&o.DriverID,
				&o.DriverName,
				&o.DriverLat,
				&o.DriverLon,
			); err != nil {
				log.Println("Ошибка rows.Scan:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при парсинге заказов"})
				return
			}
			orders = append(orders, o)
		}

		c.JSON(http.StatusOK, orders)
	}
}
