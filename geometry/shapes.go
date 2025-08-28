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

// Text 文本（简化版）
type Text struct {
	*core.BaseMobject
	text string
	size float64
}

// NewText 创建新的文本
func NewText(text string, size float64) *Text {
	textObj := &Text{
		BaseMobject: core.NewBaseMobject(),
		text:        text,
		size:        size,
	}
	textObj.SetColor(color.RGBA{255, 255, 255, 255})
	return textObj
}

// GetText 获取文本内容
func (t *Text) GetText() string {
	return t.text
}

// SetText 设置文本内容
func (t *Text) SetText(text string) *Text {
	t.text = text
	return t
}

// GetSize 获取文本大小
func (t *Text) GetSize() float64 {
	return t.size
}

// SetSize 设置文本大小
func (t *Text) SetSize(size float64) *Text {
	t.size = size
	return t
}
