package math

import (
	"math"
)

// CoordinateSystem 坐标系统 - 左下角为原点
type CoordinateSystem struct {
	Width  int
	Height int
	Scale  float64
}

// NewCoordinateSystem 创建新的坐标系统
func NewCoordinateSystem(width, height int) *CoordinateSystem {
	return &CoordinateSystem{
		Width:  width,
		Height: height,
		Scale:  1.0,
	}
}

// ToScreen 将逻辑坐标转换为屏幕坐标（左下角为原点）
func (cs *CoordinateSystem) ToScreen(logical Vector2) Vector2 {
	return Vector2{
		X: logical.X * cs.Scale,
		Y: float64(cs.Height) - logical.Y*cs.Scale, // Y轴翻转，屏幕坐标Y向下，逻辑坐标Y向上
	}
}

// ToLogical 将屏幕坐标转换为逻辑坐标（左下角为原点）
func (cs *CoordinateSystem) ToLogical(screen Vector2) Vector2 {
	return Vector2{
		X: screen.X / cs.Scale,
		Y: (float64(cs.Height) - screen.Y) / cs.Scale, // Y轴翻转
	}
}

// SetScale 设置坐标系缩放
func (cs *CoordinateSystem) SetScale(scale float64) {
	cs.Scale = scale
}

// GetBounds 获取逻辑坐标边界（左下角为原点）
func (cs *CoordinateSystem) GetBounds() (minX, maxX, minY, maxY float64) {
	bottomLeft := cs.ToLogical(Vector2{0, float64(cs.Height)})
	topRight := cs.ToLogical(Vector2{float64(cs.Width), 0})
	return bottomLeft.X, topRight.X, bottomLeft.Y, topRight.Y
}

// Vector2 表示2D向量
type Vector2 struct {
	X, Y float64
}

// Vector3 表示3D向量
type Vector3 struct {
	X, Y, Z float64
}

// NewVector2 创建新的2D向量
func NewVector2(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
}

// NewVector3 创建新的3D向量
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

// Add 向量加法
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub 向量减法
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale 向量缩放
func (v Vector2) Scale(factor float64) Vector2 {
	return Vector2{X: v.X * factor, Y: v.Y * factor}
}

// Length 向量长度
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize 向量标准化
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{0, 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

// Distance 计算两点间距离
func (v Vector2) Distance(other Vector2) float64 {
	return v.Sub(other).Length()
}

// Angle 计算向量角度
func (v Vector2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Rotate 旋转向量
func (v Vector2) Rotate(angle float64) Vector2 {
	cos, sin := math.Cos(angle), math.Sin(angle)
	return Vector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Dot 点积
func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// 3D向量方法
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

func (v Vector3) Scale(factor float64) Vector3 {
	return Vector3{X: v.X * factor, Y: v.Y * factor, Z: v.Z * factor}
}

func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length == 0 {
		return Vector3{0, 0, 0}
	}
	return Vector3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

// Cross 叉积
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Interpolate 线性插值
func Interpolate(a, b, t float64) float64 {
	return a + t*(b-a)
}

// SmoothStep 平滑插值
func SmoothStep(t float64) float64 {
	return t * t * (3 - 2*t)
}

// EaseInOut 缓入缓出
func EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseIn 缓入
func EaseIn(t float64) float64 {
	return t * t
}

// EaseOut 缓出
func EaseOut(t float64) float64 {
	return 1 - (1-t)*(1-t)
}

// Clamp 限制值范围
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// LerpVector2 线性插值向量
func LerpVector2(a, b Vector2, t float64) Vector2 {
	return Vector2{
		X: Interpolate(a.X, b.X, t),
		Y: Interpolate(a.Y, b.Y, t),
	}
}
