package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"mahou-textbox/config"
	"mahou-textbox/handlers"
)

var (
	mu sync.RWMutex
)

func main() {
	router := gin.Default()

	// 提供静态文件服务
	router.Static("/frontend", "./frontend")
	router.Static("/images", "./images")

	// API路由
	api := router.Group("/api")
	{
		// 角色相关API
		api.GET("/characters", handlers.GetCharacters)
		api.GET("/characters/current", handlers.GetCurrentCharacter) // 保持这个接口用于获取默认角色
		api.GET("/characters/:characterId/emotions", handlers.GetEmotions)

		// 背景相关API
		api.GET("/backgrounds", handlers.GetBackgrounds)

		// 图片生成API
		api.POST("/generate", handlers.GenerateImage)
	}

	port := 8080
	if config.AppConfig.Port != 0 {
		port = config.AppConfig.Port
	}

	fmt.Printf("服务器启动在 http://localhost:%d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
