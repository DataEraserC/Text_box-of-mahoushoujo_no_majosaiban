package handlers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"mahou-textbox/config"
)

// GetCharacters 获取所有角色列表
func GetCharacters(c *gin.Context) {
	// 创建一个有序的角色ID列表
	var characterIds []string
	for id := range config.Characters {
		characterIds = append(characterIds, id)
	}
	
	// 按字符串排序以保证顺序稳定
	sort.Strings(characterIds)

	var chars []map[string]interface{}
	for _, id := range characterIds {
		char := config.Characters[id]
		chars = append(chars, map[string]interface{}{
			"id":   id,
			"name": char.Name,
		})
	}
	c.JSON(http.StatusOK, chars)
}

// GetCurrentCharacter 获取默认角色
func GetCurrentCharacter(c *gin.Context) {
	// 总是返回默认角色，不保存状态
	defaultCharacter := config.GetDefaultCharacter()

	char, exists := config.Characters[defaultCharacter]

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "默认角色不存在"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":   defaultCharacter,
		"name": char.Name,
	})
}

// GetEmotions 获取角色表情列表
func GetEmotions(c *gin.Context) {
	characterId := c.Param("characterId")

	char, exists := config.Characters[characterId]

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "角色不存在"})
		return
	}

	var emotions []map[string]interface{}
	for i, emotion := range char.Emotions {
		emotions = append(emotions, map[string]interface{}{
			"id":   i + 1,
			"name": emotion.Name,
		})
	}

	c.JSON(http.StatusOK, emotions)
}