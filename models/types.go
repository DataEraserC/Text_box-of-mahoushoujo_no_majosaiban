package models

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

// TextBoxConfig 文本框坐标配置
type TextBoxConfig struct {
	Position [2]int // 文本框左上角坐标
	Over     [2]int // 文本框右下角坐标
}

// AppConfig 应用配置
type AppConfig struct {
	TextBox struct {
		Position []int `json:"position"`
		Over     []int `json:"over"`
	} `json:"text_box"`
	DefaultCharacter string `json:"default_character"`
	Port             int    `json:"port"`
}

// GenerateImageParams 生成图片的参数
type GenerateImageParams struct {
	CharacterID     string
	Text            string
	EmotionIndex    *int
	BackgroundIndex *int
	TextConfigs     []TextConfig
}