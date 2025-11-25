package config

import (
	"encoding/json"
	"os"
	"time"
	"math/rand"

	"mahou-textbox/models"
)

var (
	TextBoxConfig    models.TextBoxConfig
	AppConfig        models.AppConfig
	Characters       map[string]models.Character
	TextConfigs      map[string][]models.TextConfig
	Backgrounds      []models.Background
)

func init() {
	rand.Seed(time.Now().UnixNano())
	
	// 加载应用配置
	LoadAppConfig()

	// 加载角色配置
	LoadCharacters()

	// 加载背景配置
	LoadBackgrounds()

	// 初始化文字配置
	InitTextConfigs()
}

// LoadAppConfig 加载应用配置
func LoadAppConfig() {
	file, err := os.ReadFile("config/app.json")
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		TextBoxConfig = models.TextBoxConfig{
			Position: [2]int{728, 355},
			Over:     [2]int{2339, 800},
		}
		return
	}
	
	if err := json.Unmarshal(file, &AppConfig); err != nil {
		// 如果解析失败，使用默认配置
		TextBoxConfig = models.TextBoxConfig{
			Position: [2]int{728, 355},
			Over:     [2]int{2339, 800},
		}
		return
	}
	
	// 设置文本框配置
	TextBoxConfig = models.TextBoxConfig{
		Position: [2]int{AppConfig.TextBox.Position[0], AppConfig.TextBox.Position[1]},
		Over:     [2]int{AppConfig.TextBox.Over[0], AppConfig.TextBox.Over[1]},
	}
}

// GetDefaultCharacter 获取默认角色ID
func GetDefaultCharacter() string {
	if AppConfig.DefaultCharacter != "" {
		return AppConfig.DefaultCharacter
	}
	return "char2" // 橘雪莉作为默认角色
}

// LoadCharacters 加载角色配置
func LoadCharacters() {
	file, err := os.ReadFile("config/characters.json")
	if err != nil {
		panic("无法读取角色配置文件: " + err.Error())
	}

	var chars []models.Character
	if err := json.Unmarshal(file, &chars); err != nil {
		panic("无法解析角色配置文件: " + err.Error())
	}

	Characters = make(map[string]models.Character)
	for _, char := range chars {
		Characters[char.ID] = char
	}
}

// LoadBackgrounds 加载背景配置
func LoadBackgrounds() {
	file, err := os.ReadFile("config/backgrounds.json")
	if err != nil {
		panic("无法读取背景配置文件: " + err.Error())
	}

	if err := json.Unmarshal(file, &Backgrounds); err != nil {
		panic("无法解析背景配置文件: " + err.Error())
	}
}

// InitTextConfigs 初始化文字配置（保留以确保向后兼容）
func InitTextConfigs() {
	TextConfigs = map[string][]models.TextConfig{
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