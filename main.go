package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

// Character 角色信息
type Character struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	DisplayName []DisplayNamePart `json:"displayName"`
	Emotions    []Emotion         `json:"emotions"`
}

// DisplayNamePart 角色姓名显示配置
type DisplayNamePart struct {
	Text      string `json:"text"`
	Position  []int  `json:"position"`
	FontColor []int  `json:"fontColor"`
	FontSize  int    `json:"fontSize"`
}

// Emotion 表情信息
type Emotion struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

// Background 背景信息
type Background struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

// TextConfig 角色文字配置 (保留以确保向后兼容)
type TextConfig struct {
	Text      string `json:"text"`
	Position  []int  `json:"position"`
	FontColor []int  `json:"font_color"`
	FontSize  int    `json:"font_size"`
}

// GenerateRequest 生成图片的请求
type GenerateRequest struct {
	Type            string `json:"type"`
	Content         string `json:"content"`
	TextInput       string `json:"textInput"`
	CharacterId     string `json:"characterId,omitempty"`
	EmotionIndex    *int   `json:"emotionIndex,omitempty"`
	BackgroundIndex *int   `json:"backgroundIndex,omitempty"`
}

var (
	characters  map[string]Character
	textConfigs map[string][]TextConfig
	backgrounds []Background
	mu          sync.RWMutex
)

func init() {
	// 加载应用配置
	loadAppConfig()

	// 加载角色配置
	loadCharacters()

	// 加载背景配置
	loadBackgrounds()

	// 初始化文字配置
	initTextConfigs()
}

// loadCharacters 加载角色配置
func loadCharacters() {
	file, err := os.ReadFile("config/characters.json")
	if err != nil {
		panic(fmt.Sprintf("无法读取角色配置文件: %v", err))
	}

	var chars []Character
	if err := json.Unmarshal(file, &chars); err != nil {
		panic(fmt.Sprintf("无法解析角色配置文件: %v", err))
	}

	characters = make(map[string]Character)
	for _, char := range chars {
		characters[char.ID] = char
	}
}

// loadBackgrounds 加载背景配置
func loadBackgrounds() {
	file, err := os.ReadFile("config/backgrounds.json")
	if err != nil {
		panic(fmt.Sprintf("无法读取背景配置文件: %v", err))
	}

	if err := json.Unmarshal(file, &backgrounds); err != nil {
		panic(fmt.Sprintf("无法解析背景配置文件: %v", err))
	}
}

// initTextConfigs 初始化文字配置（保留以确保向后兼容）
func initTextConfigs() {
	textConfigs = map[string][]TextConfig{
		"char0": {
			{Text: "樱", Position: []int{759, 73}, FontColor: []int{253, 145, 175}, FontSize: 186},
			{Text: "羽", Position: []int{949, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "艾", Position: []int{1039, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "玛", Position: []int{1183, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char1": {
			{Text: "二", Position: []int{759, 63}, FontColor: []int{239, 79, 84}, FontSize: 196},
			{Text: "阶堂", Position: []int{955, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "希", Position: []int{1143, 110}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "罗", Position: []int{1283, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char2": {
			{Text: "橘", Position: []int{759, 73}, FontColor: []int{137, 177, 251}, FontSize: 186},
			{Text: "雪", Position: []int{943, 110}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "莉", Position: []int{1093, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "", Position: []int{0, 0}, FontColor: []int{255, 255, 255}, FontSize: 1},
		},
		"char3": {
			{Text: "远", Position: []int{759, 73}, FontColor: []int{169, 199, 30}, FontSize: 186},
			{Text: "野", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "汉", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "娜", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char4": {
			{Text: "夏", Position: []int{759, 73}, FontColor: []int{159, 145, 251}, FontSize: 186},
			{Text: "目", Position: []int{949, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "安", Position: []int{1039, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "安", Position: []int{1183, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char5": {
			{Text: "月", Position: []int{759, 63}, FontColor: []int{195, 209, 231}, FontSize: 196},
			{Text: "代", Position: []int{948, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "雪", Position: []int{1053, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "", Position: []int{0, 0}, FontColor: []int{255, 255, 255}, FontSize: 1},
		},
		"char6": {
			{Text: "冰", Position: []int{759, 73}, FontColor: []int{227, 185, 175}, FontSize: 186},
			{Text: "上", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "梅", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "露露", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char7": {
			{Text: "城", Position: []int{759, 73}, FontColor: []int{104, 223, 231}, FontSize: 186},
			{Text: "崎", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "诺", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "亚", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char8": {
			{Text: "莲", Position: []int{759, 73}, FontColor: []int{253, 177, 88}, FontSize: 186},
			{Text: "见", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "蕾", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "雅", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char9": {
			{Text: "佐", Position: []int{759, 73}, FontColor: []int{235, 207, 139}, FontSize: 186},
			{Text: "伯", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "米", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "莉亚", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char10": {
			{Text: "黑", Position: []int{759, 63}, FontColor: []int{131, 143, 147}, FontSize: 196},
			{Text: "部", Position: []int{955, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "奈", Position: []int{1053, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "叶香", Position: []int{1197, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char11": {
			{Text: "宝", Position: []int{759, 73}, FontColor: []int{185, 124, 235}, FontSize: 186},
			{Text: "生", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "玛", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "格", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char12": {
			{Text: "紫", Position: []int{759, 73}, FontColor: []int{235, 75, 60}, FontSize: 186},
			{Text: "藤", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "亚", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "里沙", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"char13": {
			{Text: "泽", Position: []int{759, 73}, FontColor: []int{251, 114, 78}, FontSize: 186},
			{Text: "渡", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "可", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "可", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
	}
}

func main() {
	router := gin.Default()

	// 提供静态文件服务
	router.Static("/frontend", "./frontend")
	router.Static("/images", "./images")

	// API路由
	api := router.Group("/api")
	{
		// 角色相关API
		api.GET("/characters", getCharacters)
		api.GET("/characters/current", getCurrentCharacter) // 保持这个接口用于获取默认角色
		api.GET("/characters/:characterId/emotions", getEmotions)

		// 背景相关API
		api.GET("/backgrounds", getBackgrounds)

		// 图片生成API
		api.POST("/generate", generateImage)
	}

	fmt.Println("服务器启动在 http://localhost:8080")
	router.Run(":8080")
}

// getCharacters 获取所有角色列表
func getCharacters(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	// 创建一个有序的角色ID列表
	var characterIds []string
	for id := range characters {
		characterIds = append(characterIds, id)
	}
	
	// 按字符串排序以保证顺序稳定
	sort.Strings(characterIds)

	var chars []map[string]interface{}
	for _, id := range characterIds {
		char := characters[id]
		chars = append(chars, map[string]interface{}{
			"id":   id,
			"name": char.Name,
		})
	}
	c.JSON(http.StatusOK, chars)
}

// getCurrentCharacter 获取默认角色
func getCurrentCharacter(c *gin.Context) {
	// 总是返回默认角色，不保存状态
	defaultCharacter := getDefaultCharacter()

	mu.RLock()
	char, exists := characters[defaultCharacter]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "默认角色不存在"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":   defaultCharacter,
		"name": char.Name,
	})
}

// getEmotions 获取角色表情列表
func getEmotions(c *gin.Context) {
	characterId := c.Param("characterId")

	mu.RLock()
	char, exists := characters[characterId]
	mu.RUnlock()

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

// getBackgrounds 获取背景列表
func getBackgrounds(c *gin.Context) {
	c.JSON(http.StatusOK, backgrounds)
}

// generateImage 生成图片
func generateImage(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	// 确定使用的角色ID，默认为配置文件中指定的默认角色
	characterId := getDefaultCharacter()
	// 如果请求中指定了"random"，则随机选择角色
	if req.CharacterId == "random" {
		characterId = getRandomCharacter()
	} else if req.CharacterId != "" {
		characterId = req.CharacterId
	}

	mu.RLock()
	_, exists := characters[characterId]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "角色不存在"})
		return
	}

	// 生成图片
	img, err := createImageWithText(characterId, req.TextInput, req.EmotionIndex, req.BackgroundIndex)
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

// getRandomCharacter 随机获取一个角色
func getRandomCharacter() string {
	mu.RLock()
	defer mu.RUnlock()

	// 将map转换为slice以便随机选择
	characterIds := make([]string, 0, len(characters))
	for id := range characters {
		characterIds = append(characterIds, id)
	}

	// 随机选择一个角色
	if len(characterIds) > 0 {
		return characterIds[rand.Intn(len(characterIds))]
	}

	// 如果没有角色，返回默认角色
	return getDefaultCharacter()
}

// createImageWithText 创建带文本的图片
func createImageWithText(characterId, text string, emotionIndex *int, backgroundIndex *int) (image.Image, error) {
	// 使用新的图片处理逻辑

	// 获取当前角色的文字配置
	mu.RLock()
	character, exists := characters[characterId]
	var configs []TextConfig
	
	// 如果角色有displayName配置，则使用它，否则使用旧的textConfigs
	if exists && len(character.DisplayName) > 0 {
		// 将DisplayNamePart转换为TextConfig以保持向后兼容
		for _, part := range character.DisplayName {
			configs = append(configs, TextConfig{
				Text:      part.Text,
				Position:  part.Position,
				FontColor: part.FontColor,
				FontSize:  part.FontSize,
			})
		}
	} else {
		// 使用旧的配置
		configs = textConfigs[characterId]
	}
	mu.RUnlock()

	// 构造图片生成参数
	params := GenerateImageParams{
		CharacterID:     characterId,
		Text:            text,
		EmotionIndex:    emotionIndex,
		BackgroundIndex: backgroundIndex,
		TextConfigs:     configs,
	}

	// 生成图片
	img, err := GenerateImage(params)
	if err != nil {
		return nil, fmt.Errorf("生成图片失败: %v", err)
	}

	return img, nil
}