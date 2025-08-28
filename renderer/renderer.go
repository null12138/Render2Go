package renderer

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"render2go/core"
	"render2go/geometry"
	gmMath "render2go/math"

	"github.com/fogleman/gg"
)

// Renderer 渲染器接口
type Renderer interface {
	Clear(r, g, b float64)
	Render(object core.Mobject)
	Present()
	SaveFrame(filename string) error
	GetContext() *gg.Context
}

// CanvasRenderer 基于gg库的画布渲染器
type CanvasRenderer struct {
	context *gg.Context
	width   int
	height  int
}

// NewCanvasRenderer 创建新的画布渲染器
func NewCanvasRenderer(width, height int) *CanvasRenderer {
	return &CanvasRenderer{
		context: gg.NewContext(width, height),
		width:   width,
		height:  height,
	}
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
	if err := r.context.LoadFontFace("", text.GetSize()); err != nil {
		// 如果无法加载字体，使用默认字体
		r.context.SetFontFace(nil)
	}

	// 计算文本位置（简化版，假设居中）
	center := text.GetCenter()
	x := center.X + float64(r.width)/2
	y := center.Y + float64(r.height)/2

	r.context.DrawStringAnchored(text.GetText(), x, y, 0.5, 0.5)
}

// renderCircle 渲染圆形
func (r *CanvasRenderer) renderCircle(circle *geometry.Circle) {
	center := circle.GetCenter()
	x := center.X + float64(r.width)/2
	y := center.Y + float64(r.height)/2
	radius := circle.GetRadius()

	r.context.DrawCircle(x, y, radius)

	if circle.GetFillOpacity() > 0 {
		r.context.FillPreserve()
	}
	r.context.Stroke()
}

// renderRectangle 渲染矩形
func (r *CanvasRenderer) renderRectangle(rect *geometry.Rectangle) {
	points := rect.GetPoints()
	if len(points) < 4 {
		return
	}

	r.renderPath(points, true)
}

// renderLine 渲染直线
func (r *CanvasRenderer) renderLine(line *geometry.Line) {
	points := line.GetPoints()
	if len(points) < 2 {
		return
	}

	r.renderPath(points, false)
}

// renderArrow 渲染箭头
func (r *CanvasRenderer) renderArrow(arrow *geometry.Arrow) {
	points := arrow.GetPoints()
	if len(points) < 2 {
		return
	}

	r.renderPath(points, false)
}

// renderPolygon 渲染多边形
func (r *CanvasRenderer) renderPolygon(polygon *geometry.Polygon) {
	points := polygon.GetPoints()
	if len(points) < 3 {
		return
	}

	r.renderPath(points, true)
}

// renderGeneric 渲染通用对象
func (r *CanvasRenderer) renderGeneric(object core.Mobject) {
	points := object.GetPoints()
	if len(points) == 0 {
		return
	}

	r.renderPath(points, false)
}

// renderPath 渲染路径
func (r *CanvasRenderer) renderPath(points []gmMath.Vector2, closed bool) {
	if len(points) == 0 {
		return
	}

	// 转换坐标系（从世界坐标到屏幕坐标）
	first := points[0]
	r.context.MoveTo(first.X+float64(r.width)/2, first.Y+float64(r.height)/2)

	for i := 1; i < len(points); i++ {
		point := points[i]
		r.context.LineTo(point.X+float64(r.width)/2, point.Y+float64(r.height)/2)
	}

	if closed {
		r.context.ClosePath()
	}

	r.context.Stroke()
}

// Present 呈现画面（在这个实现中是空的，因为我们直接绘制）
func (r *CanvasRenderer) Present() {
	// 在实际应用中，这里可能会包含缓冲区交换等操作
}

// SaveFrame 保存当前帧
func (r *CanvasRenderer) SaveFrame(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, r.context.Image())
}

// GetImage 获取当前图像
func (r *CanvasRenderer) GetImage() image.Image {
	return r.context.Image()
}
