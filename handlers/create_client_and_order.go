package handlers
import (
	"context"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)
func CreateClientAndOrderHandler(db *pgxpool.Pool) gin.HandlerFunc {
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
		// ✅ 1. Поиск ближайшего водителя
		var driverID int
		queryDriver := `
			SELECT id
			FROM drivers
			ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)
			LIMIT 1
		`
		err := db.QueryRow(context.Background(), queryDriver, req.Lon, req.Lat).Scan(&driverID)
		if err != nil {
			log.Println("Ошибка поиска водителя:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске водителя"})
			return
		}

		// ✅ 2. Создание клиента и заказа в транзакции
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("Ошибка начала транзакции:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}
		defer tx.Rollback(context.Background())

		var clientID int
		err = tx.QueryRow(context.Background(),
			`INSERT INTO clients (name, phone, plate_number, category_id, location)
			 VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326))
			 RETURNING id`,
			req.Name, req.Phone, req.Plate, req.CategoryID, req.Lon, req.Lat).Scan(&clientID)

		if err != nil {
			log.Println("Ошибка создания клиента:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании клиента"})
			return
		}

		_, err = tx.Exec(context.Background(),
			`INSERT INTO orders (client_id, driver_id, status, created_at)
			 VALUES ($1, $2, 'pending', NOW())`,
			clientID, driverID)

		if err != nil {
			log.Println("Ошибка создания заказа:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заказа"})
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Println("Ошибка коммита транзакции:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Клиент и заказ созданы", "driver_id": driverID})
	}
}
