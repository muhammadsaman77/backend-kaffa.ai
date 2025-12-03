package main

import (
	"context"
	"fmt"

	"backend-kaffa.ai/configs"
	"backend-kaffa.ai/internal/controllers"
	"backend-kaffa.ai/internal/services"
	"backend-kaffa.ai/internal/sqlc/users"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	ctx := context.Background()
	db := configs.InitDatabase(ctx)
	usersQueries := users.New(db)
	authService := services.NewAuthService(usersQueries)
	authController := controllers.NewAuthController(authService)

	authRouter := r.Group("/api/v1/auth")
	authRouter.POST("/login", authController.LoginUser())       // Placeholder for auth handler
	authRouter.POST("/register", authController.RegisterUser()) // Placeholder for auth handler
	r.Run(":2003")
}
