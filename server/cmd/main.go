package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/minorcell/pfss/docs" // 导入 swagger docs
	"github.com/minorcell/pfss/internal/handler"
	"github.com/minorcell/pfss/internal/model"
	"github.com/minorcell/pfss/internal/service"
	"github.com/minorcell/pfss/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title PFSS 接口文档
// @version 1.0
// @description Private File Storage System 接口文档
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化数据库连接
	const mysqlDSNFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化数据库连接，从环境变量中获取配置，data base name: DNS
	dsn := fmt.Sprintf(mysqlDSNFormat,
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

	// 初始化数据库模型
	if err := model.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化 Gin 路由器
	router := gin.Default()

	/*
		添加中间件
		1. ErrorHandler() 是自定义的中间件，用于处理错误。
		2. LoggerMiddleware() 是自定义的中间件，用于记录请求日志。
		3. gin.Recovery() 是 Gin 框架提供的一个中间件，用于捕获并处理 panic，防止程序崩溃。其主要功能是在请求处理过程中捕获并处理 panic，防止程序崩溃。当一个请求在处理过程中发生 panic，Gin 会自动调用 Recovery 中间件来恢复程序的正常运行
	*/
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	// 初始化路由，并将数据库连接传递给路由处理函数
	initializeRoutes(router, db)

	// 配置静态文件路由，使得 /upload 目录下的文件可以通过 /upload 访问到
	router.Static("/upload", "upload")

	// 启动服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)

	// 启动服务器，并监听指定端口
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initializeRoutes(router *gin.Engine, db *gorm.DB) {

	// 初始化服务层和处理器
	authService := service.NewAuthService(db)
	userService := service.NewUserService(db)
	bucketService := service.NewBucketService(db)
	fileService := service.NewFileService(db, bucketService)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	bucketHandler := handler.NewBucketHandler(bucketService)
	fileHandler := handler.NewFileHandler(fileService)

	// 配置 Swagger 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查路由，用于检查服务是否正常
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 定义路由组，用于分组管理路由，所有v1版本的路由都以/api/v1开头
	v1 := router.Group("/api/v1")

	// 不需要认证的路由

	// auth路由组，是一个从v1路由组中分出的子路由组，所有auth路由都以/api/v1/auth开头
	auth := v1.Group("/auth")
	{
		// 当路由为 /api/v1/auth/register 时，会调用 authHandler.Register 方法处理请求
		auth.POST("/register", authHandler.Register)
		// 当路由为 /api/v1/auth/login 时，会调用 authHandler.Login 方法处理请求
		auth.POST("/login", authHandler.Login)
	}

	// 受保护的路由，需要认证才能访问
	v1.Use(middleware.AuthMiddleware())
	{
		// 用户管理路由组
		users := v1.Group("/users")
		{
			// 用户管理路由
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)

			// 个人信息路由
			users.POST("/change-password", authHandler.ChangePassword)

			// 管理员权限路由
			users.PUT("/:id", middleware.RootRequired(), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RootRequired(), userHandler.DeleteUser)
			users.PUT("/:id/status", middleware.RootRequired(), userHandler.UpdateUserStatus)
			users.PUT("/:id/permissions", middleware.RootRequired(), userHandler.UpdateUserPermissions)
		}

		// 文件管理路由组
		files := v1.Group("/files")
		{
			// 文件管理路由
			files.POST("", fileHandler.CreateFile)
			files.GET("/bucket/:bucket_id", fileHandler.ListFiles)
			files.GET("/:id", fileHandler.GetFile)
			files.DELETE("/:id", fileHandler.DeleteFile)
		}

		// 桶管理路由组
		buckets := v1.Group("/buckets")
		{
			buckets.POST("", bucketHandler.CreateBucket)
			buckets.GET("", bucketHandler.ListBuckets)
			buckets.GET("/:id", bucketHandler.GetBucket)
			buckets.PUT("/:id", bucketHandler.UpdateBucket)
			buckets.DELETE("/:id", bucketHandler.DeleteBucket)

			buckets.PUT("/:id/permissions", bucketHandler.UpdateBucketPermissions)

			buckets.GET("/:id/stats", bucketHandler.GetBucketStats)
		}

		// 管理员权限路由组
		admin := v1.Group("/admin")
		admin.Use(middleware.RootRequired())
		{
			// TODO: Add admin routes
		}
	}
}
