package defaults

import (
	"image/color"
)

// FontSizes 预定义字体大小
var FontSizes = struct {
	Tiny   float64
	Small  float64
	Normal float64
	Large  float64
	Huge   float64
	Title  float64
}{
	Tiny:   8,
	Small:  12,
	Normal: 16,
	Large:  20,
	Huge:   24,
	Title:  28,
}

// Colors 预定义颜色
var Colors = struct {
	// 基础颜色
	Black   color.RGBA
	White   color.RGBA
	Red     color.RGBA
	Green   color.RGBA
	Blue    color.RGBA
	Yellow  color.RGBA
	Cyan    color.RGBA
	Magenta color.RGBA

	// 材质设计颜色
	Primary    color.RGBA
	Secondary  color.RGBA
	Accent     color.RGBA
	Background color.RGBA
	Surface    color.RGBA
	Error      color.RGBA

	// 语义颜色
	Success color.RGBA
	Warning color.RGBA
	Info    color.RGBA
	Muted   color.RGBA

	// 数学颜色
	MathRed    color.RGBA
	MathBlue   color.RGBA
	MathGreen  color.RGBA
	MathOrange color.RGBA
	MathPurple color.RGBA
}{
	// 基础颜色
	Black:   color.RGBA{0, 0, 0, 255},
	White:   color.RGBA{255, 255, 255, 255},
	Red:     color.RGBA{255, 0, 0, 255},
	Green:   color.RGBA{0, 255, 0, 255},
	Blue:    color.RGBA{0, 0, 255, 255},
	Yellow:  color.RGBA{255, 255, 0, 255},
	Cyan:    color.RGBA{0, 255, 255, 255},
	Magenta: color.RGBA{255, 0, 255, 255},

	// 材质设计颜色
	Primary:    color.RGBA{33, 150, 243, 255},  // #2196F3
	Secondary:  color.RGBA{255, 152, 0, 255},   // #FF9800
	Accent:     color.RGBA{255, 87, 34, 255},   // #FF5722
	Background: color.RGBA{250, 250, 250, 255}, // #FAFAFA
	Surface:    color.RGBA{255, 255, 255, 255}, // #FFFFFF
	Error:      color.RGBA{244, 67, 54, 255},   // #F44336

	// 语义颜色
	Success: color.RGBA{76, 175, 80, 255},   // #4CAF50
	Warning: color.RGBA{255, 193, 7, 255},   // #FFC107
	Info:    color.RGBA{33, 150, 243, 255},  // #2196F3
	Muted:   color.RGBA{158, 158, 158, 255}, // #9E9E9E

	// 数学颜色（高对比度，清晰区分）
	MathRed:    color.RGBA{231, 76, 60, 255},  // #E74C3C
	MathBlue:   color.RGBA{52, 152, 219, 255}, // #3498DB
	MathGreen:  color.RGBA{39, 174, 96, 255},  // #27AE60
	MathOrange: color.RGBA{243, 156, 18, 255}, // #F39C12
	MathPurple: color.RGBA{142, 68, 173, 255}, // #8E44AD
}

// ColorNames 颜色名称映射
var ColorNames = map[string]color.RGBA{
	"black":   Colors.Black,
	"white":   Colors.White,
	"red":     Colors.Red,
	"green":   Colors.Green,
	"blue":    Colors.Blue,
	"yellow":  Colors.Yellow,
	"cyan":    Colors.Cyan,
	"magenta": Colors.Magenta,

	"primary":    Colors.Primary,
	"secondary":  Colors.Secondary,
	"accent":     Colors.Accent,
	"background": Colors.Background,
	"surface":    Colors.Surface,
	"error":      Colors.Error,

	"success": Colors.Success,
	"warning": Colors.Warning,
	"info":    Colors.Info,
	"muted":   Colors.Muted,

	"math-red":    Colors.MathRed,
	"math-blue":   Colors.MathBlue,
	"math-green":  Colors.MathGreen,
	"math-orange": Colors.MathOrange,
	"math-purple": Colors.MathPurple,
}

// FontSizeNames 字体大小名称映射
var FontSizeNames = map[string]float64{
	"tiny":   FontSizes.Tiny,
	"small":  FontSizes.Small,
	"normal": FontSizes.Normal,
	"large":  FontSizes.Large,
	"huge":   FontSizes.Huge,
	"title":  FontSizes.Title,
}

// GetColorByName 根据名称获取颜色
func GetColorByName(name string) (color.RGBA, bool) {
	c, exists := ColorNames[name]
	return c, exists
}

// GetFontSizeByName 根据名称获取字体大小
func GetFontSizeByName(name string) (float64, bool) {
	size, exists := FontSizeNames[name]
	return size, exists
}

// Opacity 预定义透明度
var Opacity = struct {
	Invisible  float64
	VeryLight  float64
	Light      float64
	Medium     float64
	Strong     float64
	VeryStrong float64
	Solid      float64
}{
	Invisible:  0.0,
	VeryLight:  0.1,
	Light:      0.3,
	Medium:     0.5,
	Strong:     0.7,
	VeryStrong: 0.9,
	Solid:      1.0,
}
