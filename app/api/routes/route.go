package route

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"lamoda-test/api/controller"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(db *sql.DB) *gin.Engine {
	// Инициализируем роутер gin
	r := gin.Default()

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Create warehouse
	r.POST("/create-warehouse", func(c *gin.Context) {
		var w controller.Warehouse
		if err := c.BindJSON(&w); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := controller.CreateWarehouse(db, &w); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": w})
	})

	// Create product
	r.POST("/create-product", func(c *gin.Context) {
		var p controller.Product
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		warehouseID, err := strconv.Atoi(c.Param("warehouseID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
			return
		}

		if err := controller.CreateProduct(db, &p, warehouseID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": p})
	})

	// Удаление продукта
	r.DELETE("/delete-product/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		if err := controller.DeleteProduct(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// Резервирование продуктов
	r.POST("/reserve-products", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, controller.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}

		err := controller.ReserveProducts(db, productCodes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, controller.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})

	// Отмена резервирования продуктов
	r.POST("/release-products", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, controller.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}

		err := controller.ReleaseProducts(db, productCodes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, controller.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})

	// Получения оставшегося количества продуктов на складе
	r.GET("/remaining-products/:warehouseID", func(c *gin.Context) {
		warehouseID := c.Param("warehouseID")
		var id int
		if _, err := fmt.Sscan(warehouseID, &id); err != nil {
			c.JSON(http.StatusBadRequest, controller.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid warehouse ID",
			})
			return
		}

		products, err := controller.GetRemainingProducts(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, controller.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, products)
	})

	return r
}
