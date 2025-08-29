package geometry

import (
	"image"
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
	// 设置默认颜色为黑色
	circle.SetColor(color.RGBA{0, 0, 0, 255})
	circle.SetStrokeWidth(2.0)
	circle.SetFillOpacity(0.0) // 默认不填充，只描边
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
	// 文本边界框应该很小，只用于渲染定位，不影响场景边界计算
	// 使用极小的边界框，实际显示由渲染器处理
	size := 0.1 // 固定的小边界框

	// 创建文本的边界框四个角点（以文本位置为中心的小方块）
	points := []gmMath.Vector2{
		{X: t.position.X - size, Y: t.position.Y - size}, // 左下
		{X: t.position.X + size, Y: t.position.Y - size}, // 右下
		{X: t.position.X + size, Y: t.position.Y + size}, // 右上
		{X: t.position.X - size, Y: t.position.Y + size}, // 左上
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

// NewTriangle 创建新的三角形（通过三个顶点）
func NewTriangle(v1, v2, v3 gmMath.Vector2) *Triangle {
	triangle := &Triangle{
		BaseMobject: core.NewBaseMobject(),
		vertices:    [3]gmMath.Vector2{v1, v2, v3},
	}
	// 设置默认颜色为黑色，描边模式
	triangle.SetColor(color.RGBA{0, 0, 0, 255})
	triangle.SetStrokeWidth(2.0)
	triangle.SetFillOpacity(0.0) // 默认不填充，只描边
	triangle.generatePoints()
	return triangle
}

// NewTriangleByCenter 创建等腰直角三角形（通过中心点和大小）
func NewTriangleByCenter(center gmMath.Vector2, size float64) *Triangle {
	// 计算三个顶点 - 直角顶点在上方
	top := gmMath.Vector2{X: center.X, Y: center.Y + size/2}
	// 底边两个顶点
	left := gmMath.Vector2{X: center.X - size/2, Y: center.Y - size/2}
	right := gmMath.Vector2{X: center.X + size/2, Y: center.Y - size/2}

	return NewTriangle(top, left, right)
}

// NewIsoscelesRightTriangle 创建等腰直角三角形（兼容旧版本）
func NewIsoscelesRightTriangle(center gmMath.Vector2, size float64) *Triangle {
	return NewTriangleByCenter(center, size)
}

// NewEquilateralTriangle 创建等边三角形
func NewEquilateralTriangle(center gmMath.Vector2, sideLength float64) *Triangle {
	// 等边三角形的高
	height := sideLength * math.Sqrt(3) / 2

	// 计算三个顶点（顶点朝上）
	top := gmMath.Vector2{X: center.X, Y: center.Y + height/2}
	bottomLeft := gmMath.Vector2{X: center.X - sideLength/2, Y: center.Y - height/2}
	bottomRight := gmMath.Vector2{X: center.X + sideLength/2, Y: center.Y - height/2}

	return NewTriangle(top, bottomLeft, bottomRight)
}

// NewRightTriangle 创建直角三角形（指定两条直角边长度）
func NewRightTriangle(center gmMath.Vector2, width, height float64) *Triangle {
	// 直角顶点在左下角
	bottomLeft := gmMath.Vector2{X: center.X - width/2, Y: center.Y - height/2}
	bottomRight := gmMath.Vector2{X: center.X + width/2, Y: center.Y - height/2}
	topLeft := gmMath.Vector2{X: center.X - width/2, Y: center.Y + height/2}

	return NewTriangle(bottomLeft, bottomRight, topLeft)
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

// SetVertex 设置单个顶点（索引 0, 1, 2）
func (t *Triangle) SetVertex(index int, vertex gmMath.Vector2) *Triangle {
	if index >= 0 && index < 3 {
		t.vertices[index] = vertex
		t.generatePoints()
	}
	return t
}

// GetVertex 获取单个顶点
func (t *Triangle) GetVertex(index int) gmMath.Vector2 {
	if index >= 0 && index < 3 {
		return t.vertices[index]
	}
	return gmMath.Vector2{X: 0, Y: 0}
}

// GetArea 计算三角形面积
func (t *Triangle) GetArea() float64 {
	v1, v2, v3 := t.vertices[0], t.vertices[1], t.vertices[2]

	// 使用向量叉积计算面积
	area := math.Abs((v2.X-v1.X)*(v3.Y-v1.Y)-(v3.X-v1.X)*(v2.Y-v1.Y)) / 2
	return area
}

// GetPerimeter 计算三角形周长
func (t *Triangle) GetPerimeter() float64 {
	v1, v2, v3 := t.vertices[0], t.vertices[1], t.vertices[2]

	side1 := v1.Distance(v2)
	side2 := v2.Distance(v3)
	side3 := v3.Distance(v1)

	return side1 + side2 + side3
}

// IsRightTriangle 判断是否为直角三角形
func (t *Triangle) IsRightTriangle() bool {
	v1, v2, v3 := t.vertices[0], t.vertices[1], t.vertices[2]

	// 计算三边长度的平方
	a2 := v1.Distance(v2) * v1.Distance(v2)
	b2 := v2.Distance(v3) * v2.Distance(v3)
	c2 := v3.Distance(v1) * v3.Distance(v1)

	// 勾股定理检查（允许一定的浮点误差）
	epsilon := 1e-10
	return (math.Abs(a2+b2-c2) < epsilon) ||
		(math.Abs(b2+c2-a2) < epsilon) ||
		(math.Abs(c2+a2-b2) < epsilon)
}

// GetCentroid 获取重心
func (t *Triangle) GetCentroid() gmMath.Vector2 {
	v1, v2, v3 := t.vertices[0], t.vertices[1], t.vertices[2]
	return gmMath.Vector2{
		X: (v1.X + v2.X + v3.X) / 3,
		Y: (v1.Y + v2.Y + v3.Y) / 3,
	}
}

// Image 图像对象
type Image struct {
	*core.BaseMobject
	filename  string
	width     float64
	height    float64
	position  gmMath.Vector2
	imageData image.Image // 添加内存图像数据支持
}

// NewImageFromFile 从文件创建新的图像对象
func NewImageFromFile(filename string, width, height float64) *Image {
	imageObj := &Image{
		BaseMobject: core.NewBaseMobject(),
		filename:    filename,
		width:       width,
		height:      height,
		position:    gmMath.Vector2{X: 0, Y: 0},
		imageData:   nil,
	}

	// 设置默认透明度
	imageObj.SetFillOpacity(1.0)

	// 生成图像边界框
	imageObj.generateBounds()

	return imageObj
}

// NewImageFromData 从内存图像数据创建图像对象
func NewImageFromData(imageData image.Image, x, y, width, height float64) *Image {
	imageObj := &Image{
		BaseMobject: core.NewBaseMobject(),
		filename:    "",
		width:       width,
		height:      height,
		position:    gmMath.Vector2{X: x, Y: y},
		imageData:   imageData,
	}

	// 设置默认透明度
	imageObj.SetFillOpacity(1.0)

	// 生成图像边界框
	imageObj.generateBounds()

	return imageObj
}

// generateBounds 生成图像的边界框点
func (img *Image) generateBounds() {
	// 创建图像的边界框四个角点
	halfWidth := img.width / 2
	halfHeight := img.height / 2

	points := []gmMath.Vector2{
		{X: img.position.X - halfWidth, Y: img.position.Y - halfHeight}, // 左下
		{X: img.position.X + halfWidth, Y: img.position.Y - halfHeight}, // 右下
		{X: img.position.X + halfWidth, Y: img.position.Y + halfHeight}, // 右上
		{X: img.position.X - halfWidth, Y: img.position.Y + halfHeight}, // 左上
	}

	img.SetPoints(points)
}

// GetFilename 获取图像文件名
func (img *Image) GetFilename() string {
	return img.filename
}

// GetDimensions 获取图像尺寸
func (img *Image) GetDimensions() (float64, float64) {
	return img.width, img.height
}

// SetPosition 设置图像位置
func (img *Image) SetPosition(x, y float64) *Image {
	img.position = gmMath.Vector2{X: x, Y: y}
	img.generateBounds()
	return img
}

// GetImageData 获取图像数据
func (img *Image) GetImageData() image.Image {
	return img.imageData
}

// SetImageData 设置图像数据
func (img *Image) SetImageData(imageData image.Image) *Image {
	img.imageData = imageData
	return img
}
