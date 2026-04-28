package routes

import (
	"Backend/config"
	"Backend/handlers"
	"Backend/infrastructure/soap"
	"Backend/repository"
	"Backend/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func SetupRoutes(cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(handlers.RequestTrace())
	
	// Configuración de CORS oficial para evitar bloqueos en Flutter Web y Ngrok
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authSoapURL := joinURL(cfg.ServerAppSOAPBase, "auth")
	userSoapURL := joinURL(cfg.ServerAppSOAPBase, "user")
	nodeSoapURL := joinURL(cfg.ServerAppSOAPBase, "node")
	bdSoapURL := joinURL(cfg.ServerAppSOAPBase, "bd")

	// 0. Infrastructure
	soapClient := soap.NewClient()

	// 1. Repositories (Adapters)
	authRepo := repository.NewAuthSoapRepository(soapClient, authSoapURL)
	userRepo := repository.NewUserSoapRepository(soapClient, userSoapURL)
	nodeRepo := repository.NewNodeRepository(soapClient, nodeSoapURL)
	bdRepo := repository.NewBdSoapRepository(soapClient, bdSoapURL)
	batchSoapURL := joinURL(cfg.ServerAppSOAPBase, "batches")
	batchRepo := repository.NewBatchSoapRepository(soapClient, batchSoapURL)

	// 2. Services (Logic)
	authService := services.NewAuthService(authRepo)
	userService := services.NewUserService(userRepo)
	nodeService := services.NewNodeService(nodeRepo)
	bdService := services.NewBdService(bdRepo)
	batchService := services.NewBatchService(batchRepo)

	// 3. Handlers (Delivery)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	nodeHandler := handlers.NewNodeHandler(nodeService)
	bdHandler := handlers.NewBdHandler(bdService)
	batchHandler := handlers.NewBatchHandler(batchService)

	registerRESTRoutes(r.Group("/"), authHandler, userHandler, nodeHandler, bdHandler, batchHandler)
	registerRESTRoutes(r.Group("/api/v1"), authHandler, userHandler, nodeHandler, bdHandler, batchHandler)

	// Health check / Info
	r.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "up",
			"version":     "1.1.0",
			"description": "Backend Proxy for Enfok Microservices (SOAP Bridge)",
		})
	})

	return r
}

func joinURL(base, path string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/")
}

func registerRESTRoutes(group *gin.RouterGroup, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, nodeHandler *handlers.NodeHandler, bdHandler *handlers.BdHandler, batchHandler *handlers.BatchHandler) {
	auth := group.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", authHandler.LogOut)
		auth.GET("/validate", authHandler.ValidateToken)
		auth.POST("/validate", authHandler.ValidateTokenPOST)
		auth.POST("/forget-password", authHandler.ForgetPwd)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	user := group.Group("/user")
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.GET("/activity", userHandler.GetActivity)
		user.GET("/search", userHandler.SearchUser)
		user.DELETE("/account", userHandler.DeleteAccount)
		user.GET("/statistics", userHandler.GetStatistics)
	}

	// Rutas añadidas según requerimiento de telemetría
	usersGroup := group.Group("/users/:user_uuid")
	{
		usersGroup.GET("/statistics", userHandler.GetStatistics)
		usersGroup.GET("/activity", userHandler.GetActivity)
	}

	node := group.Group("/node")
	{
		node.POST("/upload", nodeHandler.UploadImages)

		// Async Batch Routes
		batch := node.Group("/batch")
		{
			batch.POST("", nodeHandler.ProcessBatch)
			batch.GET("/:id/status", nodeHandler.GetBatchStatus)
			batch.GET("/:id/results", nodeHandler.GetBatchResults)
		}
	}

	admin := group.Group("/admin")
	{
		admin.GET("/logs/:image_uuid", nodeHandler.GetLogsByImage)
		admin.GET("/metrics/:node_id", nodeHandler.GetMetricsByNode)
	}

	bd := group.Group("/bd")
	{
		bd.GET("/gallery", bdHandler.GetPaginatedImages)
		bd.GET("/batches", bdHandler.GetUserBatchesWithCovers)
		bd.GET("/image/:image_uuid/metrics", bdHandler.GetImageMetrics)
	}

	// New batch route as requested by prompt
	group.GET("/download-batch/:id", batchHandler.DownloadBatch)
}
