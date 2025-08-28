package colors

import "image/color"

// ColorScheme 配色方案
type ColorScheme struct {
	Name   string
	Colors []color.RGBA
}

// 预定义的配色方案
var (
	// ProfessionalBlue 专业蓝配色方案
	// 基于用户提供的配色分析：深蓝、中蓝、紫蓝、青蓝、深色、浅紫
	ProfessionalBlue = ColorScheme{
		Name: "Professional Blue",
		Colors: []color.RGBA{
			HexToRGBA("#051B4A"), // 深蓝色 (27.2%) - 主背景
			HexToRGBA("#274274"), // 中蓝色 (26.5%) - 主要元素
			HexToRGBA("#576DA2"), // 紫蓝色 (25.1%) - 强调色
			HexToRGBA("#2B576E"), // 青蓝色 (7.4%) - 辅助色
			HexToRGBA("#041229"), // 深色 (7.3%) - 深背景
			HexToRGBA("#9BA2C2"), // 浅紫色 (6.6%) - 文字色
		},
	}

	// 快速访问颜色常量 - 新配色方案
	DeepBlue    = HexToRGBA("#051B4A") // 深蓝色 - 主背景
	MidBlue     = HexToRGBA("#274274") // 中蓝色 - 主要元素
	PurpleBlue  = HexToRGBA("#576DA2") // 紫蓝色 - 强调色
	CyanBlue    = HexToRGBA("#2B576E") // 青蓝色 - 辅助色
	DarkColor   = HexToRGBA("#041229") // 深色 - 深背景
	LightPurple = HexToRGBA("#9BA2C2") // 浅紫色 - 文字色

	// 常用的基础颜色
	White       = color.RGBA{255, 255, 255, 255}
	Black       = color.RGBA{0, 0, 0, 255}
	Transparent = color.RGBA{0, 0, 0, 0}
)

// HexToRGBA 将十六进制颜色代码转换为RGBA
func HexToRGBA(hex string) color.RGBA {
	if len(hex) != 7 || hex[0] != '#' {
		return color.RGBA{255, 255, 255, 255} // 默认白色
	}

	r := hexToByte(hex[1:3])
	g := hexToByte(hex[3:5])
	b := hexToByte(hex[5:7])

	return color.RGBA{r, g, b, 255}
}

// hexToByte 将两位十六进制字符串转换为字节
func hexToByte(hex string) uint8 {
	if len(hex) != 2 {
		return 0
	}

	high := hexDigitToValue(hex[0])
	low := hexDigitToValue(hex[1])

	return uint8(high*16 + low)
}

// hexDigitToValue 将十六进制字符转换为数值
func hexDigitToValue(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	default:
		return 0
	}
}

// GetColorByIndex 根据索引获取配色方案中的颜色
func (cs *ColorScheme) GetColorByIndex(index int) color.RGBA {
	if index < 0 || index >= len(cs.Colors) {
		return cs.Colors[0] // 返回第一个颜色作为默认值
	}
	return cs.Colors[index]
}

// GetPrimaryColor 获取主要颜色
func (cs *ColorScheme) GetPrimaryColor() color.RGBA {
	return cs.Colors[1] // 中蓝色
}

// GetSecondaryColor 获取次要颜色
func (cs *ColorScheme) GetSecondaryColor() color.RGBA {
	return cs.Colors[2] // 紫蓝色
}

// GetAccentColor 获取强调色
func (cs *ColorScheme) GetAccentColor() color.RGBA {
	return cs.Colors[5] // 浅紫色
}

// GetBackgroundColor 获取背景色
func (cs *ColorScheme) GetBackgroundColor() color.RGBA {
	return cs.Colors[0] // 深蓝色
}

// GetLightColor 获取浅色
func (cs *ColorScheme) GetLightColor() color.RGBA {
	return cs.Colors[5] // 浅紫色
}

// CreateGradient 创建渐变色
func CreateGradient(start, end color.RGBA, steps int) []color.RGBA {
	if steps <= 1 {
		return []color.RGBA{start}
	}

	gradient := make([]color.RGBA, steps)

	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)

		r := uint8(float64(start.R)*(1-t) + float64(end.R)*t)
		g := uint8(float64(start.G)*(1-t) + float64(end.G)*t)
		b := uint8(float64(start.B)*(1-t) + float64(end.B)*t)
		a := uint8(float64(start.A)*(1-t) + float64(end.A)*t)

		gradient[i] = color.RGBA{r, g, b, a}
	}

	return gradient
}

// RGBAToFloat64 将RGBA颜色转换为0-1范围的浮点数
func RGBAToFloat64(c color.RGBA) (float64, float64, float64, float64) {
	return float64(c.R) / 255.0,
		float64(c.G) / 255.0,
		float64(c.B) / 255.0,
		float64(c.A) / 255.0
}

// Float64ToRGBA 将0-1范围的浮点数转换为RGBA颜色
func Float64ToRGBA(r, g, b, a float64) color.RGBA {
	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: uint8(a * 255),
	}
}

// 默认配色方案设置为专业蓝
var DefaultColorScheme = ProfessionalBlue
