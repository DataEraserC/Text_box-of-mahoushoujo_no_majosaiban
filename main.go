package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Character 角色信息
type Character struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	EmotionCount int    `json:"emotionCount"`
	Font         string `json:"font"`
}

// TextConfig 角色文字配置
type TextConfig struct {
	Text       string `json:"text"`
	Position   []int  `json:"position"`
	FontColor  []int  `json:"fontColor"`
	FontSize   int    `json:"fontSize"`
}

// GenerateRequest 生成图片的请求
type GenerateRequest struct {
	Type          string `json:"type"`
	Content       string `json:"content"`
	TextInput     string `json:"textInput"`
	CharacterId   string `json:"characterId,omitempty"`
	EmotionIndex  *int   `json:"emotionIndex,omitempty"`
	BackgroundIndex *int  `json:"backgroundIndex,omitempty"`
}

// Emotion 表情信息
type Emotion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	characters  map[string]Character
	textConfigs map[string][]TextConfig
	mu          sync.RWMutex
)

func init() {
	// 初始化角色信息
	characters = map[string]Character{
		"ema": {
			ID:           "ema",
			Name:         "樱羽艾玛",
			EmotionCount: 8,
			Font:         "font3.ttf",
		},
		"hiro": {
			ID:           "hiro",
			Name:         "二阶堂希罗",
			EmotionCount: 6,
			Font:         "font3.ttf",
		},
		"sherri": {
			ID:           "sherri",
			Name:         "橘雪莉",
			EmotionCount: 7,
			Font:         "font3.ttf",
		},
		"hanna": {
			ID:           "hanna",
			Name:         "远野汉娜",
			EmotionCount: 5,
			Font:         "font3.ttf",
		},
		"anan": {
			ID:           "anan",
			Name:         "夏目安安",
			EmotionCount: 9,
			Font:         "font3.ttf",
		},
		"yuki": {
			ID:           "yuki",
			Name:         "月代雪",
			EmotionCount: 18,
			Font:         "font3.ttf",
		},
		"meruru": {
			ID:           "meruru",
			Name:         "冰上梅露露",
			EmotionCount: 6,
			Font:         "font3.ttf",
		},
		"noa": {
			ID:           "noa",
			Name:         "城崎诺亚",
			EmotionCount: 6,
			Font:         "font3.ttf",
		},
		"reia": {
			ID:           "reia",
			Name:         "莲见蕾雅",
			EmotionCount: 7,
			Font:         "font3.ttf",
		},
		"miria": {
			ID:           "miria",
			Name:         "佐伯米莉亚",
			EmotionCount: 4,
			Font:         "font3.ttf",
		},
		"nanoka": {
			ID:           "nanoka",
			Name:         "黑部奈叶香",
			EmotionCount: 5,
			Font:         "font3.ttf",
		},
		"mago": {
			ID:           "mago",
			Name:         "宝生玛格",
			EmotionCount: 5,
			Font:         "font3.ttf",
		},
		"alisa": {
			ID:           "alisa",
			Name:         "紫藤亚里沙",
			EmotionCount: 6,
			Font:         "font3.ttf",
		},
		"coco": {
			ID:           "coco",
			Name:         "泽渡可可",
			EmotionCount: 5,
			Font:         "font3.ttf",
		},
	}

	// 初始化文字配置
	textConfigs = map[string][]TextConfig{
		"nanoka": {
			{Text: "黑", Position: []int{759, 63}, FontColor: []int{131, 143, 147}, FontSize: 196},
			{Text: "部", Position: []int{955, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "奈", Position: []int{1053, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "叶香", Position: []int{1197, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"hiro": {
			{Text: "二", Position: []int{759, 63}, FontColor: []int{239, 79, 84}, FontSize: 196},
			{Text: "阶堂", Position: []int{955, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "希", Position: []int{1143, 110}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "罗", Position: []int{1283, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"ema": {
			{Text: "樱", Position: []int{759, 73}, FontColor: []int{253, 145, 175}, FontSize: 186},
			{Text: "羽", Position: []int{949, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "艾", Position: []int{1039, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "玛", Position: []int{1183, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"sherri": {
			{Text: "橘", Position: []int{759, 73}, FontColor: []int{137, 177, 251}, FontSize: 186},
			{Text: "雪", Position: []int{943, 110}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "莉", Position: []int{1093, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "", Position: []int{0, 0}, FontColor: []int{255, 255, 255}, FontSize: 1},
		},
		"anan": {
			{Text: "夏", Position: []int{759, 73}, FontColor: []int{159, 145, 251}, FontSize: 186},
			{Text: "目", Position: []int{949, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "安", Position: []int{1039, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "安", Position: []int{1183, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"noa": {
			{Text: "城", Position: []int{759, 73}, FontColor: []int{104, 223, 231}, FontSize: 186},
			{Text: "崎", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "诺", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "亚", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"coco": {
			{Text: "泽", Position: []int{759, 73}, FontColor: []int{251, 114, 78}, FontSize: 186},
			{Text: "渡", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "可", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "可", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"alisa": {
			{Text: "紫", Position: []int{759, 73}, FontColor: []int{235, 75, 60}, FontSize: 186},
			{Text: "藤", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "亚", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "里沙", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"reia": {
			{Text: "莲", Position: []int{759, 73}, FontColor: []int{253, 177, 88}, FontSize: 186},
			{Text: "见", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "蕾", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "雅", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"mago": {
			{Text: "宝", Position: []int{759, 73}, FontColor: []int{185, 124, 235}, FontSize: 186},
			{Text: "生", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "玛", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "格", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"hanna": {
			{Text: "远", Position: []int{759, 73}, FontColor: []int{169, 199, 30}, FontSize: 186},
			{Text: "野", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "汉", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "娜", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"meruru": {
			{Text: "冰", Position: []int{759, 73}, FontColor: []int{227, 185, 175}, FontSize: 186},
			{Text: "上", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "梅", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "露露", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"miria": {
			{Text: "佐", Position: []int{759, 73}, FontColor: []int{235, 207, 139}, FontSize: 186},
			{Text: "伯", Position: []int{945, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "米", Position: []int{1042, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "莉亚", Position: []int{1186, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
		},
		"yuki": {
			{Text: "月", Position: []int{759, 63}, FontColor: []int{195, 209, 231}, FontSize: 196},
			{Text: "代", Position: []int{948, 175}, FontColor: []int{255, 255, 255}, FontSize: 92},
			{Text: "雪", Position: []int{1053, 117}, FontColor: []int{255, 255, 255}, FontSize: 147},
			{Text: "", Position: []int{0, 0}, FontColor: []int{255, 255, 255}, FontSize: 1},
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

		// 图片生成API
		api.POST("/generate", generateImage)
	}

	fmt.Println("服务器启动在 http://localhost:8080")
	router.Run(":8080")
}

// getCharacters 获取所有角色列表
func getCharacters(c *gin.Context) {
	var chars []Character
	for _, char := range characters {
		chars = append(chars, char)
	}
	c.JSON(http.StatusOK, chars)
}

// getCurrentCharacter 获取默认角色（橘雪莉）
func getCurrentCharacter(c *gin.Context) {
	// 总是返回默认角色橘雪莉，不保存状态
	char := characters["sherri"]
	c.JSON(http.StatusOK, char)
}

// getEmotions 获取角色表情列表
func getEmotions(c *gin.Context) {
	characterId := c.Param("characterId")

	char, exists := characters[characterId]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "角色不存在"})
		return
	}

	var emotions []Emotion
	for i := 1; i <= char.EmotionCount; i++ {
		emotions = append(emotions, Emotion{
			ID:   i,
			Name: fmt.Sprintf("表情%d", i),
		})
	}

	c.JSON(http.StatusOK, emotions)
}

// generateImage 生成图片
func generateImage(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	// 确定使用的角色ID，默认为橘雪莉
	characterId := "sherri"
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

	// 创建输出目录
	outputDir := "./images"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// 生成文件名
	filename := fmt.Sprintf("result_%d_%d.png", os.Getpid(), time.Now().UnixNano())
	filepath := filepath.Join(outputDir, filename)

	// 生成图片
	img, err := createImageWithText(characterId, req.TextInput, req.EmotionIndex, req.BackgroundIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成图片失败: " + err.Error()})
		return
	}

	// 保存图片
	file, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存图片失败: " + err.Error()})
		return
	}
	defer file.Close()

	// 编码并保存PNG图片
	if err := png.Encode(file, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "编码图片失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"imageUrl": "/images/" + filename,
		"character": characterId,  // 添加角色信息用于调试
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
	
	// 如果没有角色，默认返回第一个
	if len(characterIds) > 0 {
		return characterIds[0]
	}
	
	return "sherri" // 默认角色
}

// createImageWithText 创建带文本的图片
func createImageWithText(characterId, text string, emotionIndex *int, backgroundIndex *int) (image.Image, error) {
	// 使用新的图片处理逻辑
	
	// 获取当前角色的文字配置
	mu.RLock()
	textConfigs := textConfigs[characterId]
	mu.RUnlock()
	
	// 构造图片生成参数
	params := GenerateImageParams{
		CharacterID:     characterId,
		Text:            text,
		EmotionIndex:    emotionIndex,
		BackgroundIndex: backgroundIndex,
		TextConfigs:     textConfigs,
	}
	
	// 生成图片
	img, err := GenerateImage(params)
	if err != nil {
		return nil, fmt.Errorf("生成图片失败: %v", err)
	}
	
	return img, nil
}