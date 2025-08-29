package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"render2go/core"
	"render2go/geometry"
	_ "render2go/interfaces" // 使用 _ 导入接口包
	gmMath "render2go/math"
	"strings"

	"github.com/fogleman/gg"
)

// CanvasRenderer 基于gg库的画布渲染器
type CanvasRenderer struct {
	context             *gg.Context
	width               int
	height              int
	coordinateSystem    *gmMath.CoordinateSystem
	autoSaveProjectName string
	fontLoaded          bool    // 字体是否已加载
	lastFontSize        float64 // 上次加载的字体大小
}

// NewCanvasRenderer 创建新的画布渲染器
func NewCanvasRenderer(width, height int) *CanvasRenderer {
	return &CanvasRenderer{
		context:          gg.NewContext(width, height),
		width:            width,
		height:           height,
		coordinateSystem: gmMath.NewCoordinateSystem(width, height),
		fontLoaded:       false,
		lastFontSize:     0,
	}
}

// GetCoordinateSystem 获取坐标系统
func (r *CanvasRenderer) GetCoordinateSystem() *gmMath.CoordinateSystem {
	return r.coordinateSystem
}

// SetupCoordinateSystem 根据场景内容自动设置坐标系
func (r *CanvasRenderer) SetupCoordinateSystem(objects []core.Mobject) {
	if len(objects) == 0 {
		// 默认设置标准缩放：每单位40像素
		r.coordinateSystem.SetFixedScale(40)
		return
	}

	// 分析对象类型，为文本对象设置特殊处理
	hasTextOnly := true
	realMinX, realMaxX := math.Inf(1), math.Inf(-1)
	realMinY, realMaxY := math.Inf(1), math.Inf(-1)

	for _, obj := range objects {
		switch obj.(type) {
		case *geometry.Text:
			// 文本对象使用其位置而不是边界框
			center := obj.GetCenter()
			if center.X < realMinX {
				realMinX = center.X
			}
			if center.X > realMaxX {
				realMaxX = center.X
			}
			if center.Y < realMinY {
				realMinY = center.Y
			}
			if center.Y > realMaxY {
				realMaxY = center.Y
			}
		case *geometry.CoordinateSystem:
			// 跳过坐标系对象，避免影响布局
			continue
		default:
			hasTextOnly = false
			// 几何对象使用实际边界
			points := obj.GetPoints()
			for _, point := range points {
				if point.X < realMinX {
					realMinX = point.X
				}
				if point.X > realMaxX {
					realMaxX = point.X
				}
				if point.Y < realMinY {
					realMinY = point.Y
				}
				if point.Y > realMaxY {
					realMaxY = point.Y
				}
			}
		}
	}

	// 如果没有有效边界，使用默认缩放
	if realMinX == math.Inf(1) || realMaxX == math.Inf(-1) {
		r.coordinateSystem.SetFixedScale(40)
		return
	}

	// 计算实际内容范围
	rangeX := realMaxX - realMinX
	rangeY := realMaxY - realMinY

	// 如果只有文本对象，使用更大的默认范围
	if hasTextOnly {
		rangeX = math.Max(rangeX, 8) // 最小8单位宽度
		rangeY = math.Max(rangeY, 6) // 最小6单位高度
	}

	// 添加合理边距
	expandedRangeX := rangeX + 4 // 固定4单位边距
	expandedRangeY := rangeY + 3 // 固定3单位边距

	// 设置自动缩放
	r.coordinateSystem.SetAutoScale(expandedRangeX, expandedRangeY)
}

// calculateBounds 计算所有对象的边界
func (r *CanvasRenderer) calculateBounds(objects []core.Mobject) (minX, maxX, minY, maxY float64) {
	minX, maxX = math.Inf(1), math.Inf(-1)
	minY, maxY = math.Inf(1), math.Inf(-1)

	for _, obj := range objects {
		points := obj.GetPoints()
		for _, point := range points {
			if point.X < minX {
				minX = point.X
			}
			if point.X > maxX {
				maxX = point.X
			}
			if point.Y < minY {
				minY = point.Y
			}
			if point.Y > maxY {
				maxY = point.Y
			}
		}
	}

	return minX, maxX, minY, maxY
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
		r.context.SetRGB(0, 0, 0) // 默认黑色而不是白色
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
	case *geometry.CoordinateSystem:
		r.renderCoordinateSystem(obj)
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

	// 获取原始字体大小
	fontSize := text.GetSize()
	if fontSize <= 0 {
		fontSize = 12 // 默认字体大小
	}

	// 简化字体大小计算 - 直接使用原始大小，不进行复杂缩放
	// 只做基本的范围限制
	scaledFontSize := fontSize
	if scaledFontSize < 8 {
		scaledFontSize = 8 // 最小字体大小
	} else if scaledFontSize > 36 {
		scaledFontSize = 36 // 降低最大字体大小
	}

	// 获取文本中心位置并转换到屏幕坐标
	center := text.GetCenter()
	screenPos := r.coordinateSystem.ToScreen(center)

	// 设置文本颜色和透明度
	if c, ok := text.GetColor().(color.RGBA); ok {
		// 对于文本，如果fillOpacity为0，默认设为1.0（完全不透明）
		opacity := text.GetFillOpacity()
		if opacity <= 0 {
			opacity = 1.0
		}
		alpha := float64(c.A) * opacity / 255.0
		r.context.SetRGBA(float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, alpha)
	} else {
		// 如果颜色设置有问题，使用黑色作为默认颜色
		r.context.SetRGBA(0, 0, 0, 1.0)
	}

	// 尝试加载中文字体，使用缓存机制避免重复加载
	if !r.fontLoaded || math.Abs(r.lastFontSize-scaledFontSize) > 1.0 {
		err := r.loadChineseFont(scaledFontSize)
		if err != nil {
			// 如果加载中文字体失败，使用gg的默认字体
			if err := r.context.LoadFontFace("", scaledFontSize); err != nil {
				// 如果连默认字体都加载失败，设置一个基本字体
				r.context.SetFontFace(nil) // 使用内置字体
			}
		}
		r.fontLoaded = true
		r.lastFontSize = scaledFontSize
	}

	// 渲染文本，居中对齐
	r.context.DrawStringAnchored(textContent, screenPos.X, screenPos.Y, 0.5, 0.5)
}

// renderCircle 渲染圆形
func (r *CanvasRenderer) renderCircle(circle *geometry.Circle) {
	center := circle.GetCenter()
	screenPos := r.coordinateSystem.ToScreen(center)
	radius := circle.GetRadius() * r.coordinateSystem.Scale

	// 设置颜色
	if c, ok := circle.GetColor().(color.RGBA); ok {
		r.context.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
	} else {
		r.context.SetRGB(0, 0, 0) // 默认黑色
	}

	// 设置线宽
	r.context.SetLineWidth(circle.GetStrokeWidth())

	r.context.DrawCircle(screenPos.X, screenPos.Y, radius)

	if circle.GetFillOpacity() > 0 {
		// 如果有填充，先填充再描边
		fillColor := circle.GetColor().(color.RGBA)
		alpha := float64(fillColor.A) * circle.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(fillColor.R)/255.0, float64(fillColor.G)/255.0, float64(fillColor.B)/255.0, alpha)
		r.context.FillPreserve() // 保持路径用于后续描边

		// 重设描边颜色
		r.context.SetRGBA255(int(fillColor.R), int(fillColor.G), int(fillColor.B), int(fillColor.A))
		r.context.Stroke()
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

	// 设置颜色 - 使用通用的颜色设置逻辑
	if c, ok := triangle.GetColor().(color.RGBA); ok {
		r.context.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
	} else {
		r.context.SetRGB(0, 0, 0) // 默认黑色而不是白色
	}

	// 设置线宽
	r.context.SetLineWidth(triangle.GetStrokeWidth())

	// 绘制三角形路径
	r.context.MoveTo(v1.X, v1.Y)
	r.context.LineTo(v2.X, v2.Y)
	r.context.LineTo(v3.X, v3.Y)
	r.context.ClosePath()

	if triangle.GetFillOpacity() > 0 {
		// 如果有填充，先填充再描边
		fillColor := triangle.GetColor().(color.RGBA)
		alpha := float64(fillColor.A) * triangle.GetFillOpacity() / 255.0
		r.context.SetRGBA(float64(fillColor.R)/255.0, float64(fillColor.G)/255.0, float64(fillColor.B)/255.0, alpha)
		r.context.FillPreserve() // 保持路径用于后续描边

		// 重设描边颜色
		r.context.SetRGBA255(int(fillColor.R), int(fillColor.G), int(fillColor.B), int(fillColor.A))
		r.context.Stroke()
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

// renderCoordinateSystem 渲染坐标系
func (r *CanvasRenderer) renderCoordinateSystem(cs *geometry.CoordinateSystem) {
	// 渲染网格线（如果有）
	for _, gridLine := range cs.GetGridLines() {
		r.Render(gridLine)
	}

	// 渲染坐标轴
	if xAxis := cs.GetXAxis(); xAxis != nil {
		r.Render(xAxis)
	}
	if yAxis := cs.GetYAxis(); yAxis != nil {
		r.Render(yAxis)
	}

	// 渲染原点标记（如果有）
	if origin := cs.GetOrigin(); origin != nil {
		r.Render(origin)
	}

	// 渲染标签（如果有）
	for _, label := range cs.GetLabels() {
		r.Render(label)
	}
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
			"C:/Windows/Fonts/msyh.ttc",    // 微软雅黑
			"C:/Windows/Fonts/msyhbd.ttc",  // 微软雅黑 Bold
			"C:/Windows/Fonts/simhei.ttf",  // 黑体
			"C:/Windows/Fonts/simsun.ttc",  // 宋体
			"C:/Windows/Fonts/arial.ttf",   // Arial (英文后备)
			"C:/Windows/Fonts/calibri.ttf", // Calibri (英文后备)
		}
	default: // Linux/Unix
		fontPaths = []string{
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/System/Library/Fonts/PingFang.ttc", // macOS
			"/System/Library/Fonts/Arial.ttf",    // macOS英文后备
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

	// 如果没有找到系统字体，返回错误让调用者处理
	return fmt.Errorf("no suitable font found")
}
