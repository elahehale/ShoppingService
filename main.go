package main

import (
	"example/web-service-gin/authentication"
	"example/web-service-gin/core"
	"example/web-service-gin/middlewares"
	"example/web-service-gin/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	models.ConnectDataBase()
	defer models.DB.Close()

	authentication_endpoints := router.Group("/user")
	authentication_endpoints.POST("/register", authentication.Register)
	authentication_endpoints.POST("/login", authentication.Login)

	protected := router.Group("/api")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.GET("/user", authentication.CurrentUser)
	protected.GET("/basket/:id", core.GetBasketByID)
	protected.DELETE("/basket/:id", core.DeleteBasketByID)
	protected.PATCH("/basket/:id", core.UpdateBasket)
	protected.POST("/basket", core.AddBasket)
	protected.GET("/basket", core.GetBasketsList)

	router.Run("localhost:8080")

}
