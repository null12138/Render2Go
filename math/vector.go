package math

import (
	"math"
)

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
