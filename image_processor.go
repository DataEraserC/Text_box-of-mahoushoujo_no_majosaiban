package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"path/filepath"
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
	backgroundIndex := rand.Intn(16) + 1
	backgroundPath := filepath.Join(wd, "background", fmt.Sprintf("c%d.png", backgroundIndex))
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

	// 计算合适的字体大小 (简化实现)
	fontSize := float64(textBoxHeight) / 6
	if fontSize > 100 {
		fontSize = 100
	}

	// 加载字体
	font, err := loadFont(fontFile, fontSize)
	if err != nil {
		// 如果无法加载指定字体，使用默认字体
		font, err = loadDefaultFont(fontSize)
		if err != nil {
			return err
		}
	}

	// 创建字体绘制上下文
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255})) // 白色文字

	// 简化文本绘制 - 直接在文本框中心绘制
	textX := textBoxConfig.Position[0] + textBoxWidth/2 - len(text)*int(fontSize)/4
	textY := textBoxConfig.Position[1] + textBoxHeight/2

	// 添加阴影效果
	shadowColor := image.NewUniform(color.RGBA{0, 0, 0, 255}) // 黑色阴影
	c.SetSrc(shadowColor)
	_, err = c.DrawString(text, freetype.Pt(textX+2, textY+2))
	if err != nil {
		return err
	}

	// 绘制主文字
	c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255})) // 白色文字
	_, err = c.DrawString(text, freetype.Pt(textX, textY))
	return err
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

		// 绘制阴影
		shadowColor := image.NewUniform(color.RGBA{0, 0, 0, 255})
		c.SetSrc(shadowColor)
		_, err = c.DrawString(config.Text, freetype.Pt(config.Position[0]+2, config.Position[1]+2))
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
		_, err = c.DrawString(config.Text, freetype.Pt(config.Position[0], config.Position[1]))
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