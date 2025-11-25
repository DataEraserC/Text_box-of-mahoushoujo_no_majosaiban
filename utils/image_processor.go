package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"mahou-textbox/config"
	"mahou-textbox/models"
)

// GenerateImage 生成完整的魔法少女裁判图片
func GenerateImage(params models.GenerateImageParams) (image.Image, error) {
	// 获取角色配置
	character, exists := config.Characters[params.CharacterID]
	
	if !exists {
		return nil, fmt.Errorf("角色 %s 不存在", params.CharacterID)
	}

	// 确定使用的表情索引
	emotionIndex := getRandomEmotionIndex(len(character.Emotions), params.EmotionIndex)

	// 确定使用的背景索引
	backgroundIndex := getRandomBackgroundIndex(params.BackgroundIndex)

	// 构造背景和角色图片路径
	wd, _ := os.Getwd()
	
	// 使用指定或随机的背景图片
	var backgroundPath string
	if backgroundIndex > 0 && backgroundIndex <= len(config.Backgrounds) {
		backgroundPath = filepath.Join(wd, config.Backgrounds[backgroundIndex-1].Filename)
	} else {
		backgroundPath = filepath.Join(wd, "backgrounds", fmt.Sprintf("bg%d.png", backgroundIndex))
	}
	
	// 构造角色图片路径
	var characterImagePath string
	if emotionIndex > 0 && emotionIndex <= len(character.Emotions) {
		characterImagePath = filepath.Join(wd, character.Emotions[emotionIndex-1].Filename)
	} else {
		characterImagePath = filepath.Join(wd, "characters", fmt.Sprintf("char_%d.png", emotionIndex))
	}

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
		characterImg, image.Point{0, 0}, draw.Over)

	// 在图片上绘制文本
	if params.Text != "" {
		// 为不同角色使用不同的字体文件
		fontFile := "font3.ttf" // 默认字体
		
		err := drawTextOnImage(resultImg, params.Text, fontFile, params.TextConfigs)
		if err != nil {
			// 如果绘制文本失败，仅记录日志但不中断流程
			fmt.Printf("警告: 绘制文本失败: %v\n", err)
		}
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
	if specifiedIndex != nil && *specifiedIndex >= 1 && *specifiedIndex <= len(config.Backgrounds) {
		return *specifiedIndex
	}
	
	// 随机选择一个背景
	return rand.Intn(len(config.Backgrounds)) + 1
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
func drawTextOnImage(img *image.RGBA, text, fontFile string, textConfigs []models.TextConfig) error {
	// 获取文本框区域
	textBoxWidth := config.TextBoxConfig.Over[0] - config.TextBoxConfig.Position[0]
	textBoxHeight := config.TextBoxConfig.Over[1] - config.TextBoxConfig.Position[1]

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
	startY := config.TextBoxConfig.Position[1]

	// 水平左对齐起始位置 (与Python版本一致)
	startX := config.TextBoxConfig.Position[0]

	// 绘制每一行文本
	for i, line := range lines {
		y := startY + i*lineHeight + int(bestFontSize)

		// 绘制阴影 (偏移2个像素，与Python版本一致)
		shadowColor := image.NewUniform(color.RGBA{0, 0, 0, 255})
		c.SetSrc(shadowColor)
		_, err := c.DrawString(line, freetype.Pt(startX+2, y+2))
		if err != nil {
			continue
		}

		// 绘制主文字
		textColor := image.NewUniform(color.RGBA{255, 255, 255, 255})
		c.SetSrc(textColor)
		_, err = c.DrawString(line, freetype.Pt(startX, y))
		if err != nil {
			continue
		}
	}

	// 绘制角色特定的文本配置（如姓名水印）
	for _, config := range textConfigs {
		// 加载指定字体大小
		font, err := loadFont(fontFile, float64(config.FontSize))
		if err != nil {
			continue
		}

		c.SetFont(font)
		c.SetFontSize(float64(config.FontSize))

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

// testFontSizeFit 测试指定字体大小是否适合文本框
func testFontSizeFit(font *truetype.Font, text string, maxWidth, maxHeight int, fontSize float64) bool {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(fontSize)

	lines := wrapTextToFit(font, text, maxWidth, fontSize)
	lineHeight := int(fontSize * 1.15)
	totalHeight := len(lines) * lineHeight

	return totalHeight <= maxHeight
}

// wrapTextToFit 将文本包装成适合指定宽度的多行
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
						// 单个字符就超过宽度，这种情况理论上不应该出现
						line = unit
					}
				}
			}
		}
		
		// 添加最后一行
		if line != "" {
			lines = append(lines, line)
		}
	}
	
	return lines
}

// getTextWidth 获取文本宽度
func getTextWidth(c *freetype.Context, text string) int {
	width := 0
	
	// 准确测量文本宽度
	w, err := c.DrawString(text, freetype.Pt(0, 0))
	if err == nil {
		width = w.X.Floor()
	}
	
	return width
}

// breakLongWord 拆分长单词
func breakLongWord(c *freetype.Context, word string, maxWidth int) string {
	// 对于超长单词，我们按字符逐步构建直到达到最大宽度
	result := ""
	for _, r := range word {
		trial := result + string(r)
		if getTextWidth(c, trial) > maxWidth {
			// 如果加上这个字符会超出宽度，则停止
			break
		}
		result = trial
	}
	return result
}