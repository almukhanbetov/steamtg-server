package main

import (
	"context"
	"log"
	"os"
	"steamtg/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Нет .env файла или ошибка загрузки, использую переменные окружения")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	r := gin.Default()
	
	r.Use(cors.Default())
	routes.SetupRoutes(r, db)
	r.Run("0.0.0.0:8989")
}
