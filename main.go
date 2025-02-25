package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	// 配置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 创建上传文件的目录
	uploadDir := "uploads"

	// 设置文件大小限制 (默认为32MB)
	r.MaxMultipartMemory = 8 << 20 // 8MB

	// 静态文件服务
	r.Static("/static", uploadDir)
	r.StaticFile("/", "./index.html")

	// 获取文件列表
	r.GET("/files", func(c *gin.Context) {
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取文件列表失败",
			})
			return
		}

		fileList := make([]gin.H, 0)
		for _, file := range files {
			if !file.IsDir() {
				fileList = append(fileList, gin.H{
					"name": file.Name(),
					"url":  fmt.Sprintf("/static/%s", file.Name()),
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"files": fileList,
		})
	})

	// 处理文件上传
	r.POST("/upload", func(c *gin.Context) {
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "获取文件失败",
			})
			return
		}

		// 生成保存路径
		dst := filepath.Join(uploadDir, file.Filename)

		// 保存文件
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "保存文件失败",
			})
			return
		}

		// 返回文件访问URL
		fileURL := fmt.Sprintf("/static/%s", file.Filename)
		c.JSON(http.StatusOK, gin.H{
			"message": "文件上传成功",
			"url":     fileURL,
		})
	})

	// 删除文件
	r.DELETE("/files/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		// 对URL编码的文件名进行解码
		decodedFilename, err := url.QueryUnescape(filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无效的文件名",
			})
			return
		}
		filePath := filepath.Join(uploadDir, decodedFilename)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文件不存在",
			})
			return
		}

		// 删除文件
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "删除文件失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "文件删除成功",
		})
	})

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}