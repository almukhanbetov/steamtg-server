package routes

import (
	"steamtg/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *gin.Engine, db *pgxpool.Pool) {
	// Middleware для передачи db в context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "API is running",
			})
		})
		// 🔹 DRIVERS
		api.POST("/drivers", handlers.CreateDriverHandler(db))
		api.GET("/drivers", handlers.GetDriversHandler(db))
		api.POST("/drivers/nearest", handlers.GeoSearchNearestDriverHandler(db))
		api.POST("/login", handlers.DriverLoginHandler(db))

		// 🔹 CLIENTS
		api.POST("/clients", handlers.CreateClientHandler(db))
		api.POST("/clients-with-order", handlers.CreateClientAndOrderHandler(db))

		// 🔹 CATEGORIES
		api.GET("/categories", handlers.GetCategoriesHandler(db))

		// 🔹 ORDERS
		api.GET("/orders", handlers.GetAllOrdersHandler(db))
		api.GET("/orders/driver/:id", handlers.GetOrdersForDriver(db))
		api.PUT("/orders/:id/status", handlers.UpdateOrderStatusHandler(db))
	}
}
