package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/minorcell/pfss/docs" // 导入 swagger docs
	"github.com/minorcell/pfss/internal/handler"
	"github.com/minorcell/pfss/internal/model"
	"github.com/minorcell/pfss/internal/service"
	"github.com/minorcell/pfss/pkg/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title PFSS API
// @version 1.0
// @description Private File Storage System API documentation
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize database schema
	if err := model.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Apply middleware
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())

	// Initialize routes
	initializeRoutes(r, db)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initializeRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize services
	authService := service.NewAuthService(db)
	userService := service.NewUserService(db)
	bucketService := service.NewBucketService(db)
	fileService := service.NewFileService(db, bucketService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	bucketHandler := handler.NewBucketHandler(bucketService)
	fileHandler := handler.NewFileHandler(fileService)
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	
	// Public routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	v1.Use(middleware.AuthMiddleware())
	{
		// Users routes
		users := v1.Group("/users")
		{
			// Public user endpoints
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			
			// Self-service endpoints
			users.POST("/change-password", authHandler.ChangePassword)
			
			// Admin/Root only endpoints
			users.PUT("/:id", middleware.RootRequired(), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RootRequired(), userHandler.DeleteUser)
			users.PUT("/:id/status", middleware.RootRequired(), userHandler.UpdateUserStatus)
			users.PUT("/:id/permissions", middleware.RootRequired(), userHandler.UpdateUserPermissions)
		}

		// Files routes
		files := v1.Group("/files")
		{
			// File management
			files.POST("", fileHandler.CreateFile)
			files.GET("/bucket/:bucket_id", fileHandler.ListFiles)
			files.GET("/:id", fileHandler.GetFile)
			files.PUT("/:id", fileHandler.UpdateFile)
			files.DELETE("/:id", fileHandler.DeleteFile)

			// File upload/download URLs
			files.GET("/:id/upload", fileHandler.GetUploadURL)
			files.GET("/:id/download", fileHandler.GetDownloadURL)
		}

		// Buckets routes
		buckets := v1.Group("/buckets")
		{
			// Bucket management
			buckets.POST("", bucketHandler.CreateBucket)
			buckets.GET("", bucketHandler.ListBuckets)
			buckets.GET("/:id", bucketHandler.GetBucket)
			buckets.PUT("/:id", bucketHandler.UpdateBucket)
			buckets.DELETE("/:id", bucketHandler.DeleteBucket)

			// Bucket permissions
			buckets.PUT("/:id/permissions", bucketHandler.UpdateBucketPermissions)

			// Bucket statistics
			buckets.GET("/:id/stats", bucketHandler.GetBucketStats)
		}

		// Root-only routes
		admin := v1.Group("/admin")
		admin.Use(middleware.RootRequired())
		{
			// TODO: Add admin routes
		}
	}
}
