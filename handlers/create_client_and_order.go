package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateClientAndOrderHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🔍 Сохраняем и логируем "сырой" JSON на всякий случай
		body, _ := io.ReadAll(c.Request.Body)
		log.Println("📦 Сырой JSON от клиента:", string(body))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Восстанавливаем тело запроса для дальнейшего чтения

		// 👤 Структура запроса
		var req struct {
			Name        string  `json:"name"`
			Phone       string  `json:"phone"`
			Plate       string  `json:"plate"`
			CategoryID  int     `json:"category_id"`
			Lon         float64 `json:"lon"`
			Lat         float64 `json:"lat"`
			AddressAuto string  `json:"address_auto"`
			LonAuto     float64 `json:"lon_auto"`
			LatAuto     float64 `json:"lat_auto"`
		}

		// 📥 Парсим JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("❌ Ошибка разбора JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
			return
		}

		log.Println("📥 Получено от клиента:")
		log.Println("Имя:", req.Name)
		log.Println("Телефон:", req.Phone)
		log.Println("Адрес (авто):", req.AddressAuto)
		log.Println("Координаты (авто):", req.LonAuto, req.LatAuto)

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
			log.Println("❌ Ошибка поиска водителя:", err)
			log.Println("🧭 Используемые координаты:", req.Lon, req.Lat)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске водителя"})
			return
		}
		log.Println("✅ Назначен водитель с ID:", driverID)

		// ✅ 2. Создание клиента и заказа в транзакции
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Println("❌ Ошибка начала транзакции:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}
		defer tx.Rollback(context.Background())

		var clientID int

		log.Println("📦 Вставка клиента в БД...")
		err = tx.QueryRow(context.Background(),
			`INSERT INTO clients (
				name, phone, plate_number, category_id,
				location, address_auto, location_auto
			)
			VALUES (
				$1, $2, $3, $4,
				ST_SetSRID(ST_MakePoint($5, $6), 4326),
				$7,
				ST_SetSRID(ST_MakePoint($8, $9), 4326)
			)
			RETURNING id`,
			req.Name, req.Phone, req.Plate, req.CategoryID,
			req.Lon, req.Lat,
			req.AddressAuto,
			req.LonAuto, req.LatAuto,
		).Scan(&clientID)
		if err != nil {
			log.Println("❌ Ошибка создания клиента:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании клиента"})
			return
		}
		log.Println("✅ Клиент создан с ID:", clientID)

		// ✅ 3. Создание заказа
		_, err = tx.Exec(context.Background(),
			`INSERT INTO orders (client_id, driver_id, status, created_at)
			 VALUES ($1, $2, 'pending', NOW())`,
			clientID, driverID)

		if err != nil {
			log.Println("❌ Ошибка создания заказа:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заказа"})
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Println("❌ Ошибка коммита транзакции:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}

		log.Println("✅ Клиент и заказ успешно сохранены.")
		c.JSON(http.StatusOK, gin.H{"message": "Клиент и заказ созданы", "driver_id": driverID})
	}
}
