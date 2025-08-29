package geometry

import (
	"fmt"
	"image/color"
	"math"
	"render2go/core"
	gmMath "render2go/math"
)

// CoordinateSystem 坐标系组件
type CoordinateSystem struct {
	*core.BaseMobject
	originX     float64    // 原点X坐标
	originY     float64    // 原点Y坐标
	xAxisLength float64    // X轴长度
	yAxisLength float64    // Y轴长度
	gridSpacing float64    // 网格间距
	showGrid    bool       // 是否显示网格
	showLabels  bool       // 是否显示标签
	showOrigin  bool       // 是否显示原点
	xRange      [2]float64 // X轴范围 [min, max]
	yRange      [2]float64 // Y轴范围 [min, max]

	// 组件
	xAxis     *Arrow  // X轴箭头
	yAxis     *Arrow  // Y轴箭头
	gridLines []*Line // 网格线
	labels    []*Text // 标签
	origin    *Circle // 原点标记
}

// NewCoordinateSystem 创建坐标系
func NewCoordinateSystem(xRange, yRange [2]float64, spacing float64) *CoordinateSystem {
	cs := &CoordinateSystem{
		BaseMobject: core.NewBaseMobject(),
		originX:     0,
		originY:     0,
		xAxisLength: xRange[1] - xRange[0],
		yAxisLength: yRange[1] - yRange[0],
		gridSpacing: spacing,
		showGrid:    true,
		showLabels:  true,
		showOrigin:  true,
		xRange:      xRange,
		yRange:      yRange,
		gridLines:   make([]*Line, 0),
		labels:      make([]*Text, 0),
	}

	// 设置默认样式
	cs.SetColor(color.RGBA{100, 100, 100, 255}) // 灰色
	cs.SetStrokeWidth(1.0)

	cs.generateComponents()
	return cs
}

// NewStandardCoordinateSystem 创建标准坐标系 (-10 到 10)
func NewStandardCoordinateSystem() *CoordinateSystem {
	return NewCoordinateSystem([2]float64{-10, 10}, [2]float64{-10, 10}, 1.0)
}

// NewViewportCoordinateSystem 创建适合视口的坐标系
func NewViewportCoordinateSystem(viewportWidth, viewportHeight float64) *CoordinateSystem {
	// 根据视口大小自动计算合适的范围
	aspectRatio := viewportWidth / viewportHeight
	if aspectRatio > 1 {
		// 宽屏，扩展X轴范围
		halfX := 10.0 * aspectRatio
		return NewCoordinateSystem([2]float64{-halfX, halfX}, [2]float64{-10, 10}, 1.0)
	} else {
		// 高屏，扩展Y轴范围
		halfY := 10.0 / aspectRatio
		return NewCoordinateSystem([2]float64{-10, 10}, [2]float64{-halfY, halfY}, 1.0)
	}
}

// generateComponents 生成坐标系组件
func (cs *CoordinateSystem) generateComponents() {
	cs.generateAxes()
	if cs.showGrid {
		cs.generateGrid()
	}
	if cs.showLabels {
		cs.generateLabels()
	}
	if cs.showOrigin {
		cs.generateOrigin()
	}
	cs.generatePoints()
}

// generateAxes 生成坐标轴
func (cs *CoordinateSystem) generateAxes() {
	// X轴
	xStart := gmMath.Vector2{X: cs.xRange[0], Y: cs.originY}
	xEnd := gmMath.Vector2{X: cs.xRange[1], Y: cs.originY}
	cs.xAxis = NewArrow(xStart, xEnd)
	cs.xAxis.SetColor(cs.GetColor())
	cs.xAxis.SetStrokeWidth(cs.GetStrokeWidth() * 1.5) // X轴稍粗一些

	// Y轴
	yStart := gmMath.Vector2{X: cs.originX, Y: cs.yRange[0]}
	yEnd := gmMath.Vector2{X: cs.originX, Y: cs.yRange[1]}
	cs.yAxis = NewArrow(yStart, yEnd)
	cs.yAxis.SetColor(cs.GetColor())
	cs.yAxis.SetStrokeWidth(cs.GetStrokeWidth() * 1.5) // Y轴稍粗一些
}

// generateGrid 生成网格线
func (cs *CoordinateSystem) generateGrid() {
	cs.gridLines = make([]*Line, 0)

	// 垂直网格线 (平行于Y轴)
	for x := cs.xRange[0]; x <= cs.xRange[1]; x += cs.gridSpacing {
		if math.Abs(x-cs.originX) < 1e-9 { // 跳过坐标轴
			continue
		}
		start := gmMath.Vector2{X: x, Y: cs.yRange[0]}
		end := gmMath.Vector2{X: x, Y: cs.yRange[1]}
		gridLine := NewLine(start, end)
		gridLine.SetColor(color.RGBA{200, 200, 200, 128}) // 淡灰色半透明
		gridLine.SetStrokeWidth(0.5)
		cs.gridLines = append(cs.gridLines, gridLine)
	}

	// 水平网格线 (平行于X轴)
	for y := cs.yRange[0]; y <= cs.yRange[1]; y += cs.gridSpacing {
		if math.Abs(y-cs.originY) < 1e-9 { // 跳过坐标轴
			continue
		}
		start := gmMath.Vector2{X: cs.xRange[0], Y: y}
		end := gmMath.Vector2{X: cs.xRange[1], Y: y}
		gridLine := NewLine(start, end)
		gridLine.SetColor(color.RGBA{200, 200, 200, 128}) // 淡灰色半透明
		gridLine.SetStrokeWidth(0.5)
		cs.gridLines = append(cs.gridLines, gridLine)
	}
}

// generateLabels 生成坐标标签
func (cs *CoordinateSystem) generateLabels() {
	cs.labels = make([]*Text, 0)

	// 计算标签间距 - 根据范围大小自动调整
	labelSpacing := cs.gridSpacing
	rangeSize := cs.xRange[1] - cs.xRange[0]
	if rangeSize > 15 {
		labelSpacing = cs.gridSpacing * 2 // 大范围时减少标签密度
	} else if rangeSize > 30 {
		labelSpacing = cs.gridSpacing * 5 // 超大范围时进一步减少
	}

	// X轴标签
	for x := cs.xRange[0]; x <= cs.xRange[1]; x += labelSpacing {
		if math.Abs(x-cs.originX) < 1e-9 { // 跳过原点
			continue
		}
		label := NewText(formatNumber(x), 12)
		label.SetPosition(x, cs.originY-0.5)     // 稍微偏下
		label.SetColor(color.RGBA{0, 0, 0, 255}) // 黑色
		cs.labels = append(cs.labels, label)
	}

	// Y轴标签
	for y := cs.yRange[0]; y <= cs.yRange[1]; y += labelSpacing {
		if math.Abs(y-cs.originY) < 1e-9 { // 跳过原点
			continue
		}
		label := NewText(formatNumber(y), 12)
		label.SetPosition(cs.originX-0.5, y)     // 稍微偏左
		label.SetColor(color.RGBA{0, 0, 0, 255}) // 黑色
		cs.labels = append(cs.labels, label)
	}

	// 原点标签
	if cs.showOrigin {
		originLabel := NewText("O", 14)
		originLabel.SetPosition(cs.originX-0.3, cs.originY-0.3)
		originLabel.SetColor(color.RGBA{0, 0, 0, 255})
		cs.labels = append(cs.labels, originLabel)
	}
}

// generateOrigin 生成原点标记
func (cs *CoordinateSystem) generateOrigin() {
	cs.origin = NewCircle(0.1)
	cs.origin.MoveTo(gmMath.Vector2{X: cs.originX, Y: cs.originY})
	cs.origin.SetColor(color.RGBA{0, 0, 0, 255}) // 黑色
	cs.origin.SetFillOpacity(1.0)                // 填充
}

// generatePoints 生成用于渲染的点集
func (cs *CoordinateSystem) generatePoints() {
	points := make([]gmMath.Vector2, 0)

	// 将所有组件的点添加到点集中
	if cs.xAxis != nil {
		points = append(points, cs.xAxis.GetPoints()...)
	}
	if cs.yAxis != nil {
		points = append(points, cs.yAxis.GetPoints()...)
	}

	for _, gridLine := range cs.gridLines {
		points = append(points, gridLine.GetPoints()...)
	}

	if cs.origin != nil {
		points = append(points, cs.origin.GetPoints()...)
	}

	cs.SetPoints(points)
}

// formatNumber 格式化数字显示
func formatNumber(num float64) string {
	if num == math.Trunc(num) {
		return fmt.Sprintf("%.0f", num)
	}
	return fmt.Sprintf("%.1f", num)
}

// 访问器方法

// GetXAxis 获取X轴
func (cs *CoordinateSystem) GetXAxis() *Arrow {
	return cs.xAxis
}

// GetYAxis 获取Y轴
func (cs *CoordinateSystem) GetYAxis() *Arrow {
	return cs.yAxis
}

// GetGridLines 获取网格线
func (cs *CoordinateSystem) GetGridLines() []*Line {
	return cs.gridLines
}

// GetLabels 获取标签
func (cs *CoordinateSystem) GetLabels() []*Text {
	return cs.labels
}

// GetOrigin 获取原点标记
func (cs *CoordinateSystem) GetOrigin() *Circle {
	return cs.origin
}

// 配置方法

// SetShowGrid 设置是否显示网格
func (cs *CoordinateSystem) SetShowGrid(show bool) *CoordinateSystem {
	cs.showGrid = show
	cs.generateComponents()
	return cs
}

// SetShowLabels 设置是否显示标签
func (cs *CoordinateSystem) SetShowLabels(show bool) *CoordinateSystem {
	cs.showLabels = show
	cs.generateComponents()
	return cs
}

// SetShowOrigin 设置是否显示原点
func (cs *CoordinateSystem) SetShowOrigin(show bool) *CoordinateSystem {
	cs.showOrigin = show
	cs.generateComponents()
	return cs
}

// SetGridSpacing 设置网格间距
func (cs *CoordinateSystem) SetGridSpacing(spacing float64) *CoordinateSystem {
	cs.gridSpacing = spacing
	cs.generateComponents()
	return cs
}

// SetOrigin 设置原点位置
func (cs *CoordinateSystem) SetOrigin(x, y float64) *CoordinateSystem {
	cs.originX = x
	cs.originY = y
	cs.generateComponents()
	return cs
}

// SetRange 设置坐标范围
func (cs *CoordinateSystem) SetRange(xMin, xMax, yMin, yMax float64) *CoordinateSystem {
	cs.xRange = [2]float64{xMin, xMax}
	cs.yRange = [2]float64{yMin, yMax}
	cs.xAxisLength = xMax - xMin
	cs.yAxisLength = yMax - yMin
	cs.generateComponents()
	return cs
}

// 实用方法

// PointToCoordinate 将屏幕点转换为坐标系坐标
func (cs *CoordinateSystem) PointToCoordinate(point gmMath.Vector2) gmMath.Vector2 {
	return gmMath.Vector2{
		X: point.X - cs.originX,
		Y: point.Y - cs.originY,
	}
}

// CoordinateToPoint 将坐标系坐标转换为屏幕点
func (cs *CoordinateSystem) CoordinateToPoint(coord gmMath.Vector2) gmMath.Vector2 {
	return gmMath.Vector2{
		X: coord.X + cs.originX,
		Y: coord.Y + cs.originY,
	}
}

// IsInRange 检查坐标是否在范围内
func (cs *CoordinateSystem) IsInRange(x, y float64) bool {
	return x >= cs.xRange[0] && x <= cs.xRange[1] &&
		y >= cs.yRange[0] && y <= cs.yRange[1]
}
