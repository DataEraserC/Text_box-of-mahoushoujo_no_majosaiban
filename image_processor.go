package main

import (
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

// 文本框坐标配置
type TextBoxConfig struct {
	Position [2]int // 文本框左上角坐标
	Over     [2]int // 文本框右下角坐标
}

// 角色配置
type CharacterConfig struct {
	Name         string
	EmotionCount int
	FontFile     string
}

// 全局配置
var (
	textBoxConfig = TextBoxConfig{
		Position: [2]int{728, 355},
		Over:     [2]int{2339, 800},
	}

	characterConfigs = map[string]CharacterConfig{
		"ema":    {"樱羽艾玛", 8, "font3.ttf"},
		"hiro":   {"二阶堂希罗", 6, "font3.ttf"},
		"sherri": {"橘雪莉", 7, "font3.ttf"},
		"hanna":  {"远野汉娜", 5, "font3.ttf"},
		"anan":   {"夏目安安", 9, "font3.ttf"},
		"yuki":   {"月代雪", 18, "font3.ttf"},
		"meruru": {"冰上梅露露", 6, "font3.ttf"},
		"noa":    {"城崎诺亚", 6, "font3.ttf"},
		"reia":   {"莲见蕾雅", 7, "font3.ttf"},
		"miria":  {"佐伯米莉亚", 4, "font3.ttf"},
		"nanoka": {"黑部奈叶香", 5, "font3.ttf"},
		"mago":   {"宝生玛格", 5, "font3.ttf"},
		"alisa":  {"紫藤亚里沙", 6, "font3.ttf"},
		"coco":   {"泽渡可可", 5, "font3.ttf"},
	}
)

// GenerateImageParams 生成图片的参数
type GenerateImageParams struct {
	CharacterID   string
	Text          string
	EmotionIndex  *int
	TextConfigs   []TextConfig
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

	// 构造背景和角色图片路径
	wd, _ := os.Getwd()
	
	// 随机选择背景图片 (1-16)
	backgroundIndex := rand.Intn(16) + 1
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
	totalHeight := len(lines) * lineHeight

	// 垂直顶部对齐起始位置 (与Python版本一致)
	startY := textBoxConfig.Position[1]
	if totalHeight < textBoxHeight {
		// 可以垂直居中
		startY = textBoxConfig.Position[1] + (textBoxHeight-totalHeight)/2
	}

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
		words := strings.Split(paragraph, " ")
		line := ""
		
		for _, word := range words {
			testLine := line
			if testLine != "" {
				testLine += " "
			}
			testLine += word
			
			// 测试当前行的宽度
			width := int(c.PointToFixed(fontSize).Ceil() * len(testLine) * 96 / 72 / 2) // 粗略估算
			if width <= maxWidth {
				line = testLine
			} else {
				if line != "" {
					lines = append(lines, line)
				}
				line = word
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

		// 使用与Python版本完全一致的位置
		// Python版本中直接使用配置中的位置，没有任何偏移
		positionX := config.Position[0]
		positionY := config.Position[1] + int(float64(config.FontSize)*0.75) // 调整基线对齐

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