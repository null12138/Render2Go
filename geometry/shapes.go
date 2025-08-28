package geometry

import (
	"image/color"
	"math"
	"render2go/core"
	gmMath "render2go/math"
)

// Circle 圆形
type Circle struct {
	*core.BaseMobject
	radius float64
	center gmMath.Vector2
}

// NewCircle 创建新的圆形
func NewCircle(radius float64) *Circle {
	circle := &Circle{
		BaseMobject: core.NewBaseMobject(),
		radius:      radius,
		center:      gmMath.Vector2{X: 0, Y: 0},
	}
	circle.generatePoints()
	return circle
}

func (c *Circle) generatePoints() {
	numPoints := 64
	points := make([]gmMath.Vector2, numPoints)

	for i := 0; i < numPoints; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numPoints)
		x := c.center.X + c.radius*math.Cos(angle)
		y := c.center.Y + c.radius*math.Sin(angle)
		points[i] = gmMath.Vector2{X: x, Y: y}
	}

	c.SetPoints(points)
}

// SetRadius 设置半径
func (c *Circle) SetRadius(radius float64) *Circle {
	c.radius = radius
	c.generatePoints()
	return c
}

// GetRadius 获取半径
func (c *Circle) GetRadius() float64 {
	return c.radius
}

// Rectangle 矩形
type Rectangle struct {
	*core.BaseMobject
	width  float64
	height float64
	center gmMath.Vector2
}

// NewRectangle 创建新的矩形
func NewRectangle(width, height float64) *Rectangle {
	rect := &Rectangle{
		BaseMobject: core.NewBaseMobject(),
		width:       width,
		height:      height,
		center:      gmMath.Vector2{X: 0, Y: 0},
	}
	rect.generatePoints()
	return rect
}

func (r *Rectangle) generatePoints() {
	halfWidth := r.width / 2
	halfHeight := r.height / 2

	points := []gmMath.Vector2{
		{X: r.center.X - halfWidth, Y: r.center.Y - halfHeight}, // 左下
		{X: r.center.X + halfWidth, Y: r.center.Y - halfHeight}, // 右下
		{X: r.center.X + halfWidth, Y: r.center.Y + halfHeight}, // 右上
		{X: r.center.X - halfWidth, Y: r.center.Y + halfHeight}, // 左上
		{X: r.center.X - halfWidth, Y: r.center.Y - halfHeight}, // 闭合
	}

	r.SetPoints(points)
}

// Line 直线
type Line struct {
	*core.BaseMobject
	start gmMath.Vector2
	end   gmMath.Vector2
}

// NewLine 创建新的直线
func NewLine(start, end gmMath.Vector2) *Line {
	line := &Line{
		BaseMobject: core.NewBaseMobject(),
		start:       start,
		end:         end,
	}
	line.generatePoints()
	return line
}

func (l *Line) generatePoints() {
	points := []gmMath.Vector2{l.start, l.end}
	l.SetPoints(points)
}

// Arrow 箭头
type Arrow struct {
	*Line
	headSize float64
}

// NewArrow 创建新的箭头
func NewArrow(start, end gmMath.Vector2) *Arrow {
	arrow := &Arrow{
		Line:     NewLine(start, end),
		headSize: 0.2,
	}
	arrow.generateArrowHead()
	return arrow
}

func (a *Arrow) generateArrowHead() {
	direction := a.end.Sub(a.start).Normalize()
	perpendicular := gmMath.Vector2{X: -direction.Y, Y: direction.X}

	headBase := a.end.Sub(direction.Scale(a.headSize))
	headLeft := headBase.Add(perpendicular.Scale(a.headSize / 2))
	headRight := headBase.Sub(perpendicular.Scale(a.headSize / 2))

	points := []gmMath.Vector2{
		a.start,
		a.end,
		headLeft,
		a.end,
		headRight,
	}

	a.SetPoints(points)
}

// Polygon 多边形
type Polygon struct {
	*core.BaseMobject
	vertices []gmMath.Vector2
}

// NewPolygon 创建新的多边形
func NewPolygon(vertices []gmMath.Vector2) *Polygon {
	polygon := &Polygon{
		BaseMobject: core.NewBaseMobject(),
		vertices:    make([]gmMath.Vector2, len(vertices)),
	}
	copy(polygon.vertices, vertices)
	polygon.generatePoints()
	return polygon
}

func (p *Polygon) generatePoints() {
	if len(p.vertices) == 0 {
		return
	}

	// 添加第一个点到最后以闭合多边形
	points := make([]gmMath.Vector2, len(p.vertices)+1)
	copy(points, p.vertices)
	points[len(points)-1] = p.vertices[0]

	p.SetPoints(points)
}

// RegularPolygon 正多边形
func NewRegularPolygon(sides int, radius float64) *Polygon {
	vertices := make([]gmMath.Vector2, sides)

	for i := 0; i < sides; i++ {
		angle := 2 * math.Pi * float64(i) / float64(sides)
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		vertices[i] = gmMath.Vector2{X: x, Y: y}
	}

	return NewPolygon(vertices)
}

// Text 文本对象
type Text struct {
	*core.BaseMobject
	text     string
	size     float64
	position gmMath.Vector2
}

// NewText 创建新的文本对象
func NewText(text string, size float64) *Text {
	textObj := &Text{
		BaseMobject: core.NewBaseMobject(),
		text:        text,
		size:        size,
		position:    gmMath.Vector2{X: 0, Y: 0}, // 默认位置为原点
	}

	// 设置默认文本颜色为黑色，确保在白色背景上可见
	textObj.SetColor(color.RGBA{0, 0, 0, 255})
	textObj.SetFillOpacity(1.0) // 文本默认完全不透明

	// 生成文本的边界框点（用于渲染系统）
	textObj.generateBounds()

	return textObj
}

// generateBounds 生成文本的边界框点
func (t *Text) generateBounds() {
	// 估算文本的大概尺寸（简化计算）
	width := float64(len(t.text)) * t.size * 0.6 // 每个字符大约是字体大小的0.6倍宽
	height := t.size * 1.2                       // 高度略大于字体大小

	halfWidth := width / 2
	halfHeight := height / 2

	// 创建文本的边界框四个角点
	points := []gmMath.Vector2{
		{X: t.position.X - halfWidth, Y: t.position.Y - halfHeight}, // 左下
		{X: t.position.X + halfWidth, Y: t.position.Y - halfHeight}, // 右下
		{X: t.position.X + halfWidth, Y: t.position.Y + halfHeight}, // 右上
		{X: t.position.X - halfWidth, Y: t.position.Y + halfHeight}, // 左上
	}

	t.SetPoints(points)
}

// GetText 获取文本内容
func (t *Text) GetText() string {
	return t.text
}

// SetText 设置文本内容
func (t *Text) SetText(text string) *Text {
	t.text = text
	t.generateBounds() // 重新生成边界框
	return t
}

// GetSize 获取字体大小
func (t *Text) GetSize() float64 {
	return t.size
}

// SetSize 设置字体大小
func (t *Text) SetSize(size float64) *Text {
	t.size = size
	t.generateBounds() // 重新生成边界框
	return t
}

// MoveTo 移动文本到指定位置
func (t *Text) MoveTo(pos gmMath.Vector2) core.Mobject {
	t.position = pos
	t.generateBounds() // 重新生成边界框
	return t
}

// GetCenter 获取文本中心位置
func (t *Text) GetCenter() gmMath.Vector2 {
	return t.position
}

// SetPosition 设置文本位置（别名方法）
func (t *Text) SetPosition(x, y float64) *Text {
	t.MoveTo(gmMath.Vector2{X: x, Y: y})
	return t
}

// Triangle 三角形
type Triangle struct {
	*core.BaseMobject
	vertices [3]gmMath.Vector2
}

// NewTriangle 创建新的三角形
func NewTriangle(v1, v2, v3 gmMath.Vector2) *Triangle {
	triangle := &Triangle{
		BaseMobject: core.NewBaseMobject(),
		vertices:    [3]gmMath.Vector2{v1, v2, v3},
	}
	triangle.generatePoints()
	return triangle
}

// NewIsoscelesRightTriangle 创建等腰直角三角形（底边水平，直角顶点在上方）
func NewIsoscelesRightTriangle(center gmMath.Vector2, size float64) *Triangle {
	// 计算三个顶点
	// 直角顶点在上方
	top := gmMath.Vector2{X: center.X, Y: center.Y + size/2}
	// 底边两个顶点
	left := gmMath.Vector2{X: center.X - size/2, Y: center.Y - size/2}
	right := gmMath.Vector2{X: center.X + size/2, Y: center.Y - size/2}

	return NewTriangle(top, left, right)
}

func (t *Triangle) generatePoints() {
	// 三角形的轮廓点
	points := []gmMath.Vector2{
		t.vertices[0], // 第一个顶点
		t.vertices[1], // 第二个顶点
		t.vertices[2], // 第三个顶点
		t.vertices[0], // 回到第一个顶点闭合
	}
	t.SetPoints(points)
}

// GetVertices 获取顶点
func (t *Triangle) GetVertices() [3]gmMath.Vector2 {
	return t.vertices
}

// SetVertices 设置顶点
func (t *Triangle) SetVertices(v1, v2, v3 gmMath.Vector2) *Triangle {
	t.vertices = [3]gmMath.Vector2{v1, v2, v3}
	t.generatePoints()
	return t
}
