package routes

import (
	"Backend/clients"
	"Backend/handlers"
	"Backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Dependencies
	soapClient := clients.NewSoapClient("http://servidor-java:8081/ws/image-process")
	batchService := services.NewBatchService(soapClient)
	batchHandler := handlers.NewBatchHandler(batchService)

	productService := services.NewProductService()
	productHandler := handlers.NewProductHandler(productService)

	// Router
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/batch", batchHandler.UploadBatch)

		products := apiV1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}
	
	r.GET("/api", func(c *gin.Context) {
		routes := []struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}{
			{Method: "POST", Path: "/api/v1/batch"},
			{Method: "GET", Path: "/api"},
			{Method: "GET", Path: "/api/v1/products/"},
			{Method: "POST", Path: "/api/v1/products/"},
			{Method: "GET", Path: "/api/v1/products/{id}"},
			{Method: "PUT", Path: "/api/v1/products/{id}"},
			{Method: "DELETE", Path: "/api/v1/products/{id}"},
		}
		c.JSON(http.StatusOK, routes)
	})


	return r
}
