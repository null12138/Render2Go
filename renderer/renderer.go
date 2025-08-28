package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"render2go/core"
	"render2go/geometry"
	gmMath "render2go/math"
	"strings"

	"github.com/fogleman/gg"
)

// Renderer 渲染器接口
type Renderer interface {
	Clear(r, g, b float64)
	Render(object core.Mobject)
	Present()
	SaveFrame(filename string) error
	GetContext() *gg.Context
	GetCoordinateSystem() *gmMath.CoordinateSystem
	SetAutoSaveProjectName(projectName string)
}

// CanvasRenderer 基于gg库的画布渲染器
type CanvasRenderer struct {
	context             *gg.Context
	width               int
	height              int
	coordinateSystem    *gmMath.CoordinateSystem
	autoSaveProjectName string
}

// NewCanvasRenderer 创建新的画布渲染器
func NewCanvasRenderer(width, height int) *CanvasRenderer {
	return &CanvasRenderer{
		context:          gg.NewContext(width, height),
		width:            width,
		height:           height,
		coordinateSystem: gmMath.NewCoordinateSystem(width, height),
	}
}

// GetCoordinateSystem 获取坐标系统
func (r *CanvasRenderer) GetCoordinateSystem() *gmMath.CoordinateSystem {
	return r.coordinateSystem
}

// SetAutoSaveProjectName 设置自动保存的项目名称
func (r *CanvasRenderer) SetAutoSaveProjectName(projectName string) {
	r.autoSaveProjectName = projectName
}

// Clear 清空画布
func (r *CanvasRenderer) Clear(red, green, blue float64) {
	r.context.SetRGB(red, green, blue)
	r.context.Clear()
}

// GetContext 获取绘图上下文
func (r *CanvasRenderer) GetContext() *gg.Context {
	return r.context
}

// Render 渲染对象
func (r *CanvasRenderer) Render(object core.Mobject) {
	if object == nil {
		return
	}

	points := object.GetPoints()
	if len(points) == 0 {
		return
	}

	// 设置颜色
	if c, ok := object.GetColor().(color.RGBA); ok {
		r.context.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
	} else {
		r.context.SetRGB(1, 1, 1) // 默认白色
	}

	// 设置线宽
	r.context.SetLineWidth(object.GetStrokeWidth())

	// 根据对象类型进行不同的渲染
	switch obj := object.(type) {
	case *geometry.Text:
		r.renderText(obj)
	case *geometry.Circle:
		r.renderCircle(obj)
	case *geometry.Triangle:
		r.renderTriangle(obj)
	case *geometry.Rectangle:
		r.renderRectangle(obj)
	case *geometry.Line:
		r.renderLine(obj)
	case *geometry.Arrow:
		r.renderArrow(obj)
	case *geometry.Polygon:
		r.renderPolygon(obj)
	default:
		r.renderGeneric(object)
	}
}

// renderText 渲染文本
func (r *CanvasRenderer) renderText(text *geometry.Text) {
	// 获取文本内容
	textContent := text.GetText()
	if textContent == "" {
		return // 空文本不渲染
	}

	// 获取字体大小
	fontSize := text.GetSize()
	if fontSize <= 0 {
		fontSize = 24 // 提高默认字体大小到24
	}

	// 获取文本中心位置并转换到屏幕坐标
	center := text.GetCenter()
	screenPos := r.coordinateSystem.ToScreen(center)

	// 设置文本颜色和透明度
	if c, ok := text.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * text.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	} else {
		// 如果颜色设置有问题，使用黑色作为默认颜色
		r.context.SetRGBA(0, 0, 0, 1.0)
	}

	// 尝试加载中文字体，如果失败则使用默认字体
	err := r.loadChineseFont(fontSize)
	if err != nil {
		// 如果加载中文字体失败，使用gg的默认字体但设置合适的大小
		// 这里使用LoadFontFace的空路径，gg会使用内置字体
		r.context.LoadFontFace("", fontSize)
	}

	// 渲染文本，居中对齐
	r.context.DrawStringAnchored(textContent, screenPos.X, screenPos.Y, 0.5, 0.5)
} // renderCircle 渲染圆形
func (r *CanvasRenderer) renderCircle(circle *geometry.Circle) {
	center := circle.GetCenter()
	screenPos := r.coordinateSystem.ToScreen(center)
	radius := circle.GetRadius() * r.coordinateSystem.Scale

	// 应用透明度
	if c, ok := circle.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * circle.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	r.context.DrawCircle(screenPos.X, screenPos.Y, radius)

	if circle.GetFillOpacity() > 0 {
		r.context.Fill()
	} else {
		r.context.Stroke()
	}
}

// renderTriangle 渲染三角形
func (r *CanvasRenderer) renderTriangle(triangle *geometry.Triangle) {
	vertices := triangle.GetVertices()

	// 转换顶点到屏幕坐标
	v1 := r.coordinateSystem.ToScreen(vertices[0])
	v2 := r.coordinateSystem.ToScreen(vertices[1])
	v3 := r.coordinateSystem.ToScreen(vertices[2])

	// 应用透明度
	if c, ok := triangle.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * triangle.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	// 绘制三角形路径
	r.context.MoveTo(v1.X, v1.Y)
	r.context.LineTo(v2.X, v2.Y)
	r.context.LineTo(v3.X, v3.Y)
	r.context.ClosePath()

	if triangle.GetFillOpacity() > 0 {
		r.context.Fill()
	} else {
		r.context.Stroke()
	}
}

// renderRectangle 渲染矩形
func (r *CanvasRenderer) renderRectangle(rect *geometry.Rectangle) {
	points := rect.GetPoints()
	if len(points) < 4 {
		return
	}

	// 应用透明度
	if c, ok := rect.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * rect.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	r.renderPath(points, true, rect.GetFillOpacity() > 0)
}

// renderLine 渲染直线
func (r *CanvasRenderer) renderLine(line *geometry.Line) {
	points := line.GetPoints()
	if len(points) < 2 {
		return
	}

	// 应用透明度
	if c, ok := line.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * line.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	r.renderPath(points, false, false)
}

// renderArrow 渲染箭头
func (r *CanvasRenderer) renderArrow(arrow *geometry.Arrow) {
	points := arrow.GetPoints()
	if len(points) < 2 {
		return
	}

	// 应用透明度
	if c, ok := arrow.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * arrow.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	r.renderPath(points, false, false)
}

// renderPolygon 渲染多边形
func (r *CanvasRenderer) renderPolygon(polygon *geometry.Polygon) {
	points := polygon.GetPoints()
	if len(points) < 3 {
		return
	}

	// 应用透明度
	if c, ok := polygon.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * polygon.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	r.renderPath(points, true, polygon.GetFillOpacity() > 0)
}

// renderGeneric 渲染通用对象
func (r *CanvasRenderer) renderGeneric(object core.Mobject) {
	points := object.GetPoints()
	if len(points) == 0 {
		return
	}

	// 应用透明度
	if c, ok := object.GetColor().(color.RGBA); ok {
		alpha := float64(c.A) * object.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	}

	r.renderPath(points, false, false)
}

// renderPath 渲染路径
func (r *CanvasRenderer) renderPath(points []gmMath.Vector2, closed bool, filled bool) {
	if len(points) == 0 {
		return
	}

	// 使用坐标系统转换逻辑坐标到屏幕坐标
	first := r.coordinateSystem.ToScreen(points[0])
	r.context.MoveTo(first.X, first.Y)

	for i := 1; i < len(points); i++ {
		screenPoint := r.coordinateSystem.ToScreen(points[i])
		r.context.LineTo(screenPoint.X, screenPoint.Y)
	}

	if closed {
		r.context.ClosePath()
	}

	if filled {
		r.context.Fill()
	} else {
		r.context.Stroke()
	}
}

// Present 呈现画面（自动保存到项目目录）
func (r *CanvasRenderer) Present() {
	// 如果设置了项目名称，自动保存当前帧
	if r.autoSaveProjectName != "" {
		outputDir := fmt.Sprintf("output/%s/frames", r.autoSaveProjectName)
		os.MkdirAll(outputDir, 0755)
		filename := filepath.Join(outputDir, r.autoSaveProjectName+".png")
		r.SaveFrame(filename)
	}
}

// SaveFrame 保存当前帧
func (r *CanvasRenderer) SaveFrame(filename string) error {
	// 确保文件名有.png扩展名
	if !strings.HasSuffix(filename, ".png") {
		filename = filename + ".png"
	}

	// 确保目录存在
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建保存目录失败 '%s': %v", dir, err)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建输出文件失败 '%s': %v", filename, err)
	}
	defer file.Close()

	if err := png.Encode(file, r.context.Image()); err != nil {
		return fmt.Errorf("PNG编码失败: %v", err)
	}

	return nil
}

// GetImage 获取当前图像
func (r *CanvasRenderer) GetImage() image.Image {
	return r.context.Image()
}

// loadChineseFont 尝试加载系统中文字体
func (r *CanvasRenderer) loadChineseFont(fontSize float64) error {
	// Windows系统中文字体路径
	var fontPaths []string

	// 检测操作系统并设置相应的字体路径
	switch {
	case strings.Contains(os.Getenv("OS"), "Windows"): // Windows
		fontPaths = []string{
			"C:/Windows/Fonts/msyh.ttc",   // 微软雅黑
			"C:/Windows/Fonts/simhei.ttf", // 黑体
			"C:/Windows/Fonts/simsun.ttc", // 宋体
		}
	default: // Linux/Unix
		fontPaths = []string{
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/System/Library/Fonts/PingFang.ttc", // macOS
		}
	}

	// 尝试加载第一个可用的字体
	for _, fontPath := range fontPaths {
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			continue
		}

		// 使用gg库的LoadFontFace方法加载字体
		err := r.context.LoadFontFace(fontPath, fontSize)
		if err == nil {
			return nil // 成功加载字体
		}
	}

	// 如果没有找到系统字体，使用默认字体
	return r.context.LoadFontFace("", fontSize)
}
