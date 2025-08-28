package core

import (
	"image/color"
	"math"
	gmMath "render2go/math"
)

// Mobject 可移动对象的基类，类似于Manim中的概念
type Mobject interface {
	GetPoints() []gmMath.Vector2
	SetPoints([]gmMath.Vector2)
	GetColor() color.Color
	SetColor(color.Color)
	GetStrokeWidth() float64
	SetStrokeWidth(float64)
	GetFillOpacity() float64
	SetFillOpacity(float64)
	Copy() Mobject
	MoveTo(gmMath.Vector2) Mobject
	Shift(gmMath.Vector2) Mobject
	Scale(float64) Mobject
	Rotate(float64) Mobject
	GetCenter() gmMath.Vector2
}

// BaseMobject 基础可移动对象
type BaseMobject struct {
	points      []gmMath.Vector2
	color       color.Color
	strokeWidth float64
	fillOpacity float64
}

// NewBaseMobject 创建基础可移动对象
func NewBaseMobject() *BaseMobject {
	return &BaseMobject{
		points:      make([]gmMath.Vector2, 0),
		color:       color.RGBA{255, 255, 255, 255},
		strokeWidth: 2.0,
		fillOpacity: 0.0,
	}
}

func (m *BaseMobject) GetPoints() []gmMath.Vector2 {
	return m.points
}

func (m *BaseMobject) SetPoints(points []gmMath.Vector2) {
	m.points = make([]gmMath.Vector2, len(points))
	copy(m.points, points)
}

func (m *BaseMobject) GetColor() color.Color {
	return m.color
}

func (m *BaseMobject) SetColor(c color.Color) {
	m.color = c
}

func (m *BaseMobject) GetStrokeWidth() float64 {
	return m.strokeWidth
}

func (m *BaseMobject) SetStrokeWidth(width float64) {
	m.strokeWidth = width
}

func (m *BaseMobject) GetFillOpacity() float64 {
	return m.fillOpacity
}

func (m *BaseMobject) SetFillOpacity(opacity float64) {
	m.fillOpacity = opacity
}

func (m *BaseMobject) GetCenter() gmMath.Vector2 {
	if len(m.points) == 0 {
		return gmMath.Vector2{X: 0, Y: 0}
	}

	var center gmMath.Vector2
	for _, point := range m.points {
		center = center.Add(point)
	}
	return center.Scale(1.0 / float64(len(m.points)))
}

func (m *BaseMobject) MoveTo(position gmMath.Vector2) Mobject {
	center := m.GetCenter()
	shift := position.Sub(center)
	return m.Shift(shift)
}

func (m *BaseMobject) Shift(offset gmMath.Vector2) Mobject {
	for i := range m.points {
		m.points[i] = m.points[i].Add(offset)
	}
	return m
}

func (m *BaseMobject) Scale(factor float64) Mobject {
	center := m.GetCenter()
	for i := range m.points {
		diff := m.points[i].Sub(center)
		m.points[i] = center.Add(diff.Scale(factor))
	}
	return m
}

func (m *BaseMobject) Rotate(angle float64) Mobject {
	center := m.GetCenter()
	cos, sin := math.Cos(angle), math.Sin(angle)

	for i := range m.points {
		diff := m.points[i].Sub(center)
		rotated := gmMath.Vector2{
			X: diff.X*cos - diff.Y*sin,
			Y: diff.X*sin + diff.Y*cos,
		}
		m.points[i] = center.Add(rotated)
	}
	return m
}

// Copy 创建对象的深拷贝
func (m *BaseMobject) Copy() Mobject {
	newObj := NewBaseMobject()
	newObj.SetPoints(m.points)
	newObj.SetColor(m.color)
	newObj.SetStrokeWidth(m.strokeWidth)
	newObj.SetFillOpacity(m.fillOpacity)
	return newObj
}
