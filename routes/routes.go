package routes

import (
	"Backend/handlers"
	"Backend/repository"
	"Backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Configuration (URLs for SOAP services)
	const (
		authSoapURL = "http://localhost:8080/services/auth"
		userSoapURL = "http://localhost:8080/services/user"
		nodeSoapURL = "http://localhost:8080/services/node"
	)

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

	// Router setup
	apiV1 := r.Group("/api/v1")
	{
		// Auth Routes
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authHandler.LogOut)
			auth.GET("/validate", authHandler.ValidateToken)
			auth.POST("/forget-password", authHandler.ForgetPwd)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// User Routes
		user := apiV1.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.GET("/activity", userHandler.GetActivity)
			user.GET("/search", userHandler.SearchUser)
			user.DELETE("/account", userHandler.DeleteAccount)
			user.GET("/statistics", userHandler.GetStatistics)
		}

		// Node Routes
		node := apiV1.Group("/node")
		{
			node.POST("/upload", nodeHandler.UploadImages)
		}
	}

	// Health check / Info
	r.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"version": "1.0.0",
			"description": "Backend Proxy for SOAP Services",
		})
	})

	return r
}
