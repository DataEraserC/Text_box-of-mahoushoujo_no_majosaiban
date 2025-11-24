package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

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
}

// CharacterConfig 角色配置
type CharacterConfig struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	EmotionCount int    `json:"emotion_count"`
	FontFile     string `json:"font"`
}

// GenerateImageParams 生成图片的参数
type GenerateImageParams struct {
	CharacterID     string
	Text            string
	EmotionIndex    *int
	BackgroundIndex *int
	TextConfigs     []TextConfig
}

var (
	textBoxConfig    TextBoxConfig
	characterConfigs map[string]CharacterConfig
	appConfig        AppConfig
)

func init() {
	rand.Seed(time.Now().UnixNano())
	
	// 加载应用配置
	loadAppConfig()
	
	// 加载角色配置
	loadCharacterConfigs()
}

// loadAppConfig 加载应用配置
func loadAppConfig() {
	file, err := os.ReadFile("config/app.json")
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		textBoxConfig = TextBoxConfig{
			Position: [2]int{728, 355},
			Over:     [2]int{2339, 800},
		}
		return
	}
	
	if err := json.Unmarshal(file, &appConfig); err != nil {
		// 如果解析失败，使用默认配置
		textBoxConfig = TextBoxConfig{
			Position: [2]int{728, 355},
			Over:     [2]int{2339, 800},
		}
		return
	}
	
	// 设置文本框配置
	textBoxConfig = TextBoxConfig{
		Position: [2]int{appConfig.TextBox.Position[0], appConfig.TextBox.Position[1]},
		Over:     [2]int{appConfig.TextBox.Over[0], appConfig.TextBox.Over[1]},
	}
}

// loadCharacterConfigs 加载角色配置
func loadCharacterConfigs() {
	file, err := os.ReadFile("config/characters.json")
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		setupDefaultCharacterConfigs()
		return
	}
	
	var chars []CharacterConfig
	if err := json.Unmarshal(file, &chars); err != nil {
		// 如果解析失败，使用默认配置
		setupDefaultCharacterConfigs()
		return
	}
	
	characterConfigs = make(map[string]CharacterConfig)
	for _, char := range chars {
		characterConfigs[char.ID] = char
	}
}

// setupDefaultCharacterConfigs 设置默认角色配置
func setupDefaultCharacterConfigs() {
	characterConfigs = map[string]CharacterConfig{
		"ema":    {"ema", "樱羽艾玛", 8, "font3.ttf"},
		"hiro":   {"hiro", "二阶堂希罗", 6, "font3.ttf"},
		"sherri": {"sherri", "橘雪莉", 7, "font3.ttf"},
		"hanna":  {"hanna", "远野汉娜", 5, "font3.ttf"},
		"anan":   {"anan", "夏目安安", 9, "font3.ttf"},
		"yuki":   {"yuki", "月代雪", 18, "font3.ttf"},
		"meruru": {"meruru", "冰上梅露露", 6, "font3.ttf"},
		"noa":    {"noa", "城崎诺亚", 6, "font3.ttf"},
		"reia":   {"reia", "莲见蕾雅", 7, "font3.ttf"},
		"miria":  {"miria", "佐伯米莉亚", 4, "font3.ttf"},
		"nanoka": {"nanoka", "黑部奈叶香", 5, "font3.ttf"},
		"mago":   {"mago", "宝生玛格", 5, "font3.ttf"},
		"alisa":  {"alisa", "紫藤亚里沙", 6, "font3.ttf"},
		"coco":   {"coco", "泽渡可可", 5, "font3.ttf"},
	}
}

// getDefaultCharacter 获取默认角色ID
func getDefaultCharacter() string {
	if appConfig.DefaultCharacter != "" {
		return appConfig.DefaultCharacter
	}
	return "sherri" // 默认角色
}

// GenerateImage 生成完整的魔法少女裁判图片
func GenerateImage(params GenerateImageParams) (image.Image, error) {
	// 获取角色配置
	charConfig, exists := characterConfigs[params.CharacterID]
	if !exists {
		return nil, fmt.Errorf("角色 %s 不存在", params.CharacterID)
	}

	// 确定使用的表情索引
	emotionIndex := getRandomEmotionIndex(charConfig.EmotionCount, params.EmotionIndex)

	// 确定使用的背景索引
	backgroundIndex := getRandomBackgroundIndex(params.BackgroundIndex)

	// 构造背景和角色图片路径
	wd, _ := os.Getwd()
	
	// 使用指定或随机的背景图片
	backgroundPath := filepath.Join(wd, "background", fmt.Sprintf("c%d.png", backgroundIndex))
	
	// 构造角色图片路径
	characterImagePath := filepath.Join(wd, params.CharacterID, fmt.Sprintf("%s (%d).png", params.CharacterID, emotionIndex))

	// 打开背景图片
	backgroundImg, err := openImage(backgroundPath)
	if err != nil {
		// 如果背景图片不存在，创建一个默认图片
		backgroundImg = createDefaultImage(1600, 900)
	}

	// 打开角色图片
	characterImg, err := openImage(characterImagePath)
	if err != nil {
		// 如果角色图片不存在，创建一个透明图层
		bounds := backgroundImg.Bounds()
		characterImg = image.NewRGBA(bounds)
	}

	// 创建结果图片
	bounds := backgroundImg.Bounds()
	resultImg := image.NewRGBA(bounds)

	// 绘制背景
	draw.Draw(resultImg, bounds, backgroundImg, image.Point{0, 0}, draw.Src)

	// 绘制角色图片 (在固定位置)
	characterBounds := characterImg.Bounds()
	overlayPosition := image.Point{0, 134} // 与原Python代码保持一致
	draw.Draw(resultImg, 
		image.Rectangle{
			Min: overlayPosition,
			Max: image.Point{
				X: overlayPosition.X + characterBounds.Dx(),
				Y: overlayPosition.Y + characterBounds.Dy(),
			},
		},
		characterImg, 
		characterBounds.Min, 
		draw.Over)

	// 在文本框区域内绘制文本
	err = drawTextOnImage(resultImg, params.Text, charConfig.FontFile)
	if err != nil {
		// 文本绘制失败不中断整个流程
		fmt.Printf("警告: 文本绘制失败: %v\n", err)
	}

	// 绘制角色特定文字水印
	err = drawCharacterTexts(resultImg, params.TextConfigs, charConfig.FontFile)
	if err != nil {
		fmt.Printf("警告: 角色文字水印绘制失败: %v\n", err)
	}

	return resultImg, nil
}

// getRandomEmotionIndex 获取随机或指定的表情索引
func getRandomEmotionIndex(emotionCount int, specifiedIndex *int) int {
	if specifiedIndex != nil && *specifiedIndex >= 1 && *specifiedIndex <= emotionCount {
		return *specifiedIndex
	}

	// 随机选择一个表情
	return rand.Intn(emotionCount) + 1
}

// getRandomBackgroundIndex 获取随机或指定的背景索引
func getRandomBackgroundIndex(specifiedIndex *int) int {
	if specifiedIndex != nil && *specifiedIndex >= 1 && *specifiedIndex <= 16 {
		return *specifiedIndex
	}
	
	// 随机选择一个背景 (1-16)
	return rand.Intn(16) + 1
}

// openImage 打开图片文件
func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// createDefaultImage 创建默认图片
func createDefaultImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// 填充浅灰色背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{200, 200, 200, 255})
		}
	}
	return img
}

// drawTextOnImage 在图片上绘制文本
func drawTextOnImage(img *image.RGBA, text, fontFile string) error {
	// 获取文本框区域
	textBoxWidth := textBoxConfig.Over[0] - textBoxConfig.Position[0]
	textBoxHeight := textBoxConfig.Over[1] - textBoxConfig.Position[1]

	// 加载字体并搜索最佳字体大小
	bestFontSize := float64(1)
	var bestFont *truetype.Font

	// 搜索最大合适的字体大小
	for fontSize := float64(1); fontSize <= float64(textBoxHeight) && fontSize <= 145; fontSize += 1.0 {
		font, err := loadFont(fontFile, fontSize)
		if err != nil {
			continue
		}

		// 测试当前字体大小是否合适
		if testFontSizeFit(font, text, textBoxWidth, textBoxHeight, fontSize) {
			bestFontSize = fontSize
			bestFont = font
		} else {
			break
		}
	}

	// 如果无法加载指定字体，使用默认字体
	if bestFont == nil {
		var err error
		bestFont, err = loadDefaultFont(bestFontSize)
		if err != nil {
			return err
		}
	}

	// 创建字体绘制上下文
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(bestFont)
	c.SetFontSize(bestFontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	// 文本换行处理
	lines := wrapTextToFit(bestFont, text, textBoxWidth, bestFontSize)

	// 计算行高和总高度 (使用估算值)
	lineHeight := int(bestFontSize * 1.15) // 15% 行间距
	_ = len(lines) * lineHeight

	// 垂直顶部对齐起始位置 (与Python版本一致)
	startY := textBoxConfig.Position[1]

	// 水平左对齐起始位置 (与Python版本一致)
	startX := textBoxConfig.Position[0]

	// 绘制每一行文本
	for i, line := range lines {
		y := startY + i*lineHeight + int(bestFontSize)

		// 添加阴影效果
		shadowColor := image.NewUniform(color.RGBA{0, 0, 0, 255}) // 黑色阴影
		c.SetSrc(shadowColor)
		_, err := c.DrawString(line, freetype.Pt(startX+2, y+2))
		if err != nil {
			return err
		}

		// 绘制主文字 (白色)
		c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255})) // 白色文字
		_, err = c.DrawString(line, freetype.Pt(startX, y))
		if err != nil {
			return err
		}
	}

	return nil
}

// testFontSizeFit 测试指定字体大小是否适合文本框
func testFontSizeFit(font *truetype.Font, text string, maxWidth, maxHeight int, fontSize float64) bool {
	// 文本换行处理
	lines := wrapTextToFit(font, text, maxWidth, fontSize)

	// 计算实际需要的高度 (使用估算值)
	lineHeight := int(fontSize * 1.15) // 15% 行间距
	totalHeight := len(lines) * lineHeight

	// 检查高度是否适合
	if totalHeight > maxHeight {
		return false
	}

	return true
}

// wrapTextToFit 根据字体和文本框宽度自动换行
func wrapTextToFit(font *truetype.Font, text string, maxWidth int, fontSize float64) []string {
	// 创建临时上下文用于测量文本
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(fontSize)

	var lines []string
	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		// 按空格分割单词，如果没有空格则按字符分割
		hasSpace := strings.Contains(paragraph, " ")
		var units []string
		if hasSpace {
			units = strings.Split(paragraph, " ")
		} else {
			// 将字符串转换为字符切片
			for _, r := range paragraph {
				units = append(units, string(r))
			}
		}
		
		line := ""
		
		// 连接单元的辅助函数
		unitJoin := func(a, b string) string {
			if a == "" {
				return b
			}
			if hasSpace {
				return a + " " + b
			}
			return a + b
		}

		for _, unit := range units {
			trial := unitJoin(line, unit)
			// 准确测量文本宽度
			width := getTextWidth(c, trial)
			
			if width <= maxWidth {
				line = trial
			} else {
				if line != "" {
					lines = append(lines, line)
				}
				
				// 如果单元太大，需要进一步拆分（针对无空格情况）
				if hasSpace {
					if getTextWidth(c, unit) <= maxWidth {
						line = unit
					} else {
						// 单词太长也需要拆分
						line = breakLongWord(c, unit, maxWidth)
						if line != "" {
							lines = append(lines, line)
						}
						line = ""
					}
				} else {
					// 字符级别的处理
					if getTextWidth(c, unit) <= maxWidth {
						line = unit
					} else {
						// 特殊情况：单个字符就超过了宽度
						if len(lines) > 0 && lines[len(lines)-1] != "" {
							lines = append(lines, "")
						}
						line = ""
					}
				}
			}
		}

		if line != "" {
			lines = append(lines, line)
		}

		// 如果段落为空，添加空行
		if paragraph == "" && (len(lines) == 0 || lines[len(lines)-1] != "") {
			lines = append(lines, "")
		}
	}

	return lines
}

// getTextWidth 准确测量文本宽度
func getTextWidth(c *freetype.Context, text string) int {
	// 使用freetype库准确测量文本宽度
	width, err := c.DrawString(text, freetype.Pt(0, 0))
	if err != nil {
		return 0
	}
	return int(width.X) >> 6 // 将26.6定点数转换为整数
}

// breakLongWord 拆分长单词
func breakLongWord(c *freetype.Context, word string, maxWidth int) string {
	// 对于很长的单词，尝试逐字符添加直到达到最大宽度
	runes := []rune(word)
	result := ""
	
	for i := 0; i < len(runes); i++ {
		trial := result + string(runes[i])
		width := getTextWidth(c, trial)
		
		if width <= maxWidth {
			result = trial
		} else {
			// 达到极限，返回当前结果
			break
		}
	}
	
	return result
}

// drawCharacterTexts 绘制角色特定文字水印
func drawCharacterTexts(img *image.RGBA, textConfigs []TextConfig, fontFile string) error {
	for _, config := range textConfigs {
		if config.Text == "" {
			continue
		}

		// 加载字体
		font, err := loadFont(fontFile, float64(config.FontSize))
		if err != nil {
			// 如果无法加载指定字体，使用默认字体
			font, err = loadDefaultFont(float64(config.FontSize))
			if err != nil {
				continue // 跳过这个文字配置
			}
		}

		// 创建字体绘制上下文
		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(font)
		c.SetFontSize(float64(config.FontSize))
		c.SetClip(img.Bounds())
		c.SetDst(img)

		// 设置文字颜色
		if len(config.FontColor) >= 3 {
			color := image.NewUniform(color.RGBA{
				uint8(config.FontColor[0]),
				uint8(config.FontColor[1]),
				uint8(config.FontColor[2]),
				255,
			})
			c.SetSrc(color)
		} else {
			c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255})) // 默认白色
		}

		// 使用与Python版本一致的位置，并根据用户要求整体向下调整
		positionX := config.Position[0]
		positionY := config.Position[1] + int(float64(config.FontSize))

		// 绘制阴影 (偏移2个像素，与Python版本一致)
		shadowColor := image.NewUniform(color.RGBA{0, 0, 0, 255})
		c.SetSrc(shadowColor)
		_, err = c.DrawString(config.Text, freetype.Pt(positionX+2, positionY+2))
		if err != nil {
			continue
		}

		// 绘制主文字
		if len(config.FontColor) >= 3 {
			color := image.NewUniform(color.RGBA{
				uint8(config.FontColor[0]),
				uint8(config.FontColor[1]),
				uint8(config.FontColor[2]),
				255,
			})
			c.SetSrc(color)
		}
		_, err = c.DrawString(config.Text, freetype.Pt(positionX, positionY))
		if err != nil {
			continue
		}
	}

	return nil
}

// loadFont 加载指定字体文件
func loadFont(fontFile string, size float64) (*truetype.Font, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 构造字体文件路径
	fontPath := filepath.Join(wd, fontFile)

	// 打开字体文件
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	// 解析字体
	return freetype.ParseFont(fontBytes)
}

// loadDefaultFont 加载默认字体
func loadDefaultFont(size float64) (*truetype.Font, error) {
	// 如果无法加载指定字体，返回nil，让调用者处理
	return nil, fmt.Errorf("无法加载字体")
}