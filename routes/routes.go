package routes

import (
	"Backend/config"
	"Backend/handlers"
	"Backend/repository"
	"Backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func SetupRoutes(cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins: []string{
					"http://localhost:3000",
					"http://127.0.0.1:3000",
					"http://localhost:8080",
					"http://127.0.0.1:8080",
					"http://localhost:9100",
					"http://127.0.0.1:9100",
				},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
				ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)
	r.Use(handlers.RequestTrace())

	authSoapURL := joinURL(cfg.ServerAppSOAPBase, "auth")
	userSoapURL := joinURL(cfg.ServerAppSOAPBase, "user")
	nodeSoapURL := joinURL(cfg.ServerAppSOAPBase, "node")

	// 1. Repositories (Adapters)
	authRepo := repository.NewAuthSoapRepository(authSoapURL)
	userRepo := repository.NewUserSoapRepository(userSoapURL)
	nodeRepo := repository.NewNodeRepository(nodeSoapURL)

	// 2. Services (Logic)
	authService := services.NewAuthService(authRepo)
	userService := services.NewUserService(userRepo)
	nodeService := services.NewNodeService(nodeRepo)

	// 3. Handlers (Delivery)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	nodeHandler := handlers.NewNodeHandler(nodeService)

	registerRESTRoutes(r.Group("/"), authHandler, userHandler, nodeHandler)
	registerRESTRoutes(r.Group("/api/v1"), authHandler, userHandler, nodeHandler)

	// Health check / Info
	r.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "up",
			"version":     "1.0.0",
			"description": "Backend Proxy for SOAP Services",
		})
	})

	return r
}

func joinURL(base, path string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/")
}

func registerRESTRoutes(group *gin.RouterGroup, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, nodeHandler *handlers.NodeHandler) {
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

	node := group.Group("/node")
	{
		node.POST("/upload", nodeHandler.UploadImages)

		// Async Batch Routes
		batch := node.Group("/batch")
		{
			batch.POST("", nodeHandler.ProcessBatch)
			batch.GET("/:id/status", nodeHandler.GetBatchStatus)
			batch.GET("/:id/results", nodeHandler.GetBatchResults)
			batch.GET("/:id/download", nodeHandler.DownloadBatchResult)
		}
	}
}
