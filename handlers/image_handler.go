package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"mahou-textbox/config"
	"mahou-textbox/models"
	"mahou-textbox/utils"
)

// GenerateImage 生成图片
func GenerateImage(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	// 确定使用的角色ID，默认为配置文件中指定的默认角色
	characterId := config.GetDefaultCharacter()
	// 如果请求中指定了"random"，则随机选择角色
	if req.CharacterId == "random" {
		characterId = GetRandomCharacter()
	} else if req.CharacterId != "" {
		characterId = req.CharacterId
	}

	_, exists := config.Characters[characterId]

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "角色不存在"})
		return
	}

	// 生成图片
	img, err := CreateImageWithText(characterId, req.TextInput, req.EmotionIndex, req.BackgroundIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成图片失败: " + err.Error()})
		return
	}

	// 将图片编码为base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "编码图片失败: " + err.Error()})
		return
	}

	// 将图片数据转换为base64编码
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"imageData": "data:image/png;base64," + imgBase64,
		"character": characterId, // 添加角色信息用于调试
	})
}

// GetRandomCharacter 随机获取一个角色
func GetRandomCharacter() string {
	// 将map转换为slice以便随机选择
	characterIds := make([]string, 0, len(config.Characters))
	for id := range config.Characters {
		characterIds = append(characterIds, id)
	}

	// 随机选择一个角色
	if len(characterIds) > 0 {
		return characterIds[rand.Intn(len(characterIds))]
	}

	// 如果没有角色，返回默认角色
	return config.GetDefaultCharacter()
}

// CreateImageWithText 创建带文本的图片
func CreateImageWithText(characterId, text string, emotionIndex *int, backgroundIndex *int) (image.Image, error) {
	// 使用新的图片处理逻辑

	// 获取当前角色的文字配置
	character, exists := config.Characters[characterId]
	var configs []models.TextConfig
	
	// 如果角色有displayName配置，则使用它，否则使用旧的textConfigs
	if exists && len(character.DisplayName) > 0 {
		// 将DisplayNamePart转换为TextConfig以保持向后兼容
		for _, part := range character.DisplayName {
			configs = append(configs, models.TextConfig{
				Text:      part.Text,
				Position:  part.Position,
				FontColor: part.FontColor,
				FontSize:  part.FontSize,
			})
		}
	} else {
		// 使用旧的配置
		configs = config.TextConfigs[characterId]
	}

	// 构造图片生成参数
	params := models.GenerateImageParams{
		CharacterID:     characterId,
		Text:            text,
		EmotionIndex:    emotionIndex,
		BackgroundIndex: backgroundIndex,
		TextConfigs:     configs,
	}

	// 生成图片
	img, err := utils.GenerateImage(params)
	if err != nil {
		return nil, fmt.Errorf("生成图片失败: %v", err)
	}

	return img, nil
}