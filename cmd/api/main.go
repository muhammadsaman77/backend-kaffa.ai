package main

import (
	"context"

	"backend-kaffa.ai/configs"
	"backend-kaffa.ai/internal/controllers"
	"backend-kaffa.ai/internal/middlewares"
	"backend-kaffa.ai/internal/services"
	"backend-kaffa.ai/internal/sqlc/images"
	"backend-kaffa.ai/internal/sqlc/products"
	"backend-kaffa.ai/internal/sqlc/users"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	configs.InitLogger()
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		configs.Log.Fatal("Error reading config file: " + err.Error())
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
	s3Client := configs.InitStorage(ctx)
	bucketName := viper.GetString("bucket_name")
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		configs.Log.Fatal("Error accessing S3 bucket: " + err.Error())
		return
	}

	usersQueries := users.New(db)
	imagesQueries := images.New(db)
	imageService := services.NewImageService(s3Client, aws.String(bucketName), imagesQueries)
	authService := services.NewAuthService(usersQueries)
	authController := controllers.NewAuthController(authService)

	authRouter := r.Group("/api/v1/auth")
	authRouter.POST("/login", middlewares.LoggerMiddleware, authController.LoginUser)
	authRouter.POST("/register", middlewares.LoggerMiddleware, authController.RegisterUser) // Placeholder for auth handler

	productRouter := r.Group("/api/v1/products")
	productsQueries := products.New(db)
	productService := services.NewProductService(productsQueries, imageService, db)
	productController := controllers.NewProductController(productService)

	productRouter.POST("/", middlewares.LoggerMiddleware, middlewares.AuthMiddleware, productController.CreateProduct)
	productRouter.GET("/", middlewares.LoggerMiddleware, middlewares.AuthMiddleware, productController.ListProducts)
	productRouter.GET("/:id", middlewares.LoggerMiddleware, middlewares.AuthMiddleware, productController.GetProduct)
	productRouter.PUT("/:id", middlewares.LoggerMiddleware, middlewares.AuthMiddleware, productController.UpdateProduct)
	productRouter.DELETE("/:id", middlewares.LoggerMiddleware, middlewares.AuthMiddleware, productController.DeleteProduct)
	r.Run(":2003")
}
