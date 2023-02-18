package main

import (
	"context"
	"lamoda-test/internal/app"
	"lamoda-test/internal/config"
	"lamoda-test/pkg/logging"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.GetLogger(ctx)

	logger.Info("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	ctx = logging.ContextWithLogger(ctx, logger)

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Running Application")
	a.Run(ctx)
}

// 	// Подключение к базе
// 	db, err := sql.Open("postgres", "dbname=lamoda_db user=root password=secret sslmode=require")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	// Инициализируем Роутер
// 	router := gin.Default()

// 	// Роут для swag документации
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// Роут для резервирования продуктов
// 	router.POST("/products/reserve", func(c *gin.Context) {
// 		var productCodes []string
// 		if err := c.ShouldBindJSON(&productCodes); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		if err := ReserveProducts(db, productCodes); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Products reserved successfully"})
// 	})

// 	// Роут для релизов продуктов
// 	router.POST("/products/release", func(c *gin.Context) {
// 		var productCodes []string
// 		if err := c.ShouldBindJSON(&productCodes); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		if err := ReleaseProducts(db, productCodes); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Products released successfully"})
// 	})

// 	// Роут оставшихся продуктов
// 	router.GET("/products/remaining/:warehouseID", func(c *gin.Context) {
// 		warehouseIDStr := c.Param("warehouseID")
// 		warehouseID, err := strconv.Atoi(warehouseIDStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
// 			return
// 		}

// 		products, err := GetRemainingProducts(db, warehouseID)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get remain product"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"products": products})
// 	})

// 	// Запуск сервера
// 	if err := router.Run(":8080"); err != nil {
// 		panic(err)
// 	}
// }
