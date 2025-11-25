package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mahou-textbox/config"
)

// GetBackgrounds 获取背景列表
func GetBackgrounds(c *gin.Context) {
	c.JSON(http.StatusOK, config.Backgrounds)
}