package animation

import (
	"math"
	gmMath "render2go/math"
)

// InterpolationType 插值类型
type InterpolationType int

const (
	Linear InterpolationType = iota
	Smooth
	EaseIn
	EaseOut
	EaseInOut
	Elastic
	Bounce
)

// Interpolator 插值器接口
type Interpolator interface {
	Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2
	InterpolateFloat(start, end, t float64) float64
}

// LinearInterpolator 线性插值器
type LinearInterpolator struct{}

func NewLinearInterpolator() *LinearInterpolator {
	return &LinearInterpolator{}
}

func (li *LinearInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, t),
		Y: gmMath.Interpolate(start.Y, end.Y, t),
	}
}

func (li *LinearInterpolator) InterpolateFloat(start, end, t float64) float64 {
	return gmMath.Interpolate(start, end, t)
}

// SmoothInterpolator 平滑插值器
type SmoothInterpolator struct{}

func NewSmoothInterpolator() *SmoothInterpolator {
	return &SmoothInterpolator{}
}

func (si *SmoothInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	smoothT := gmMath.SmoothStep(t)
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, smoothT),
		Y: gmMath.Interpolate(start.Y, end.Y, smoothT),
	}
}

func (si *SmoothInterpolator) InterpolateFloat(start, end, t float64) float64 {
	smoothT := gmMath.SmoothStep(t)
	return gmMath.Interpolate(start, end, smoothT)
}

// EaseInInterpolator 缓入插值器
type EaseInInterpolator struct{}

func NewEaseInInterpolator() *EaseInInterpolator {
	return &EaseInInterpolator{}
}

func (ei *EaseInInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	easeT := gmMath.EaseIn(t)
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, easeT),
		Y: gmMath.Interpolate(start.Y, end.Y, easeT),
	}
}

func (ei *EaseInInterpolator) InterpolateFloat(start, end, t float64) float64 {
	easeT := gmMath.EaseIn(t)
	return gmMath.Interpolate(start, end, easeT)
}

// EaseOutInterpolator 缓出插值器
type EaseOutInterpolator struct{}

func NewEaseOutInterpolator() *EaseOutInterpolator {
	return &EaseOutInterpolator{}
}

func (eo *EaseOutInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	easeT := gmMath.EaseOut(t)
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, easeT),
		Y: gmMath.Interpolate(start.Y, end.Y, easeT),
	}
}

func (eo *EaseOutInterpolator) InterpolateFloat(start, end, t float64) float64 {
	easeT := gmMath.EaseOut(t)
	return gmMath.Interpolate(start, end, easeT)
}

// EaseInOutInterpolator 缓入缓出插值器
type EaseInOutInterpolator struct{}

func NewEaseInOutInterpolator() *EaseInOutInterpolator {
	return &EaseInOutInterpolator{}
}

func (eio *EaseInOutInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	easeT := gmMath.EaseInOut(t)
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, easeT),
		Y: gmMath.Interpolate(start.Y, end.Y, easeT),
	}
}

func (eio *EaseInOutInterpolator) InterpolateFloat(start, end, t float64) float64 {
	easeT := gmMath.EaseInOut(t)
	return gmMath.Interpolate(start, end, easeT)
}

// ElasticInterpolator 弹性插值器
type ElasticInterpolator struct {
	amplitude float64
	period    float64
}

func NewElasticInterpolator(amplitude, period float64) *ElasticInterpolator {
	return &ElasticInterpolator{
		amplitude: amplitude,
		period:    period,
	}
}

func (el *ElasticInterpolator) elasticEaseOut(t float64) float64 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}

	p := el.period
	s := p / 4

	return (el.amplitude * math.Pow(2, -10*t) * math.Sin((t*1-s)*(2*math.Pi)/p) + 1)
}

func (el *ElasticInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	easeT := el.elasticEaseOut(t)
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, easeT),
		Y: gmMath.Interpolate(start.Y, end.Y, easeT),
	}
}

func (el *ElasticInterpolator) InterpolateFloat(start, end, t float64) float64 {
	easeT := el.elasticEaseOut(t)
	return gmMath.Interpolate(start, end, easeT)
}

// BounceInterpolator 弹跳插值器
type BounceInterpolator struct{}

func NewBounceInterpolator() *BounceInterpolator {
	return &BounceInterpolator{}
}

func (bi *BounceInterpolator) bounceEaseOut(t float64) float64 {
	if t < 1/2.75 {
		return 7.5625 * t * t
	} else if t < 2/2.75 {
		t -= 1.5 / 2.75
		return 7.5625*t*t + 0.75
	} else if t < 2.5/2.75 {
		t -= 2.25 / 2.75
		return 7.5625*t*t + 0.9375
	} else {
		t -= 2.625 / 2.75
		return 7.5625*t*t + 0.984375
	}
}

func (bi *BounceInterpolator) Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2 {
	easeT := bi.bounceEaseOut(t)
	return gmMath.Vector2{
		X: gmMath.Interpolate(start.X, end.X, easeT),
		Y: gmMath.Interpolate(start.Y, end.Y, easeT),
	}
}

func (bi *BounceInterpolator) InterpolateFloat(start, end, t float64) float64 {
	easeT := bi.bounceEaseOut(t)
	return gmMath.Interpolate(start, end, easeT)
}

// GetInterpolator 获取指定类型的插值器
func GetInterpolator(interpType InterpolationType) Interpolator {
	switch interpType {
	case Linear:
		return NewLinearInterpolator()
	case Smooth:
		return NewSmoothInterpolator()
	case EaseIn:
		return NewEaseInInterpolator()
	case EaseOut:
		return NewEaseOutInterpolator()
	case EaseInOut:
		return NewEaseInOutInterpolator()
	case Elastic:
		return NewElasticInterpolator(1.0, 0.3)
	case Bounce:
		return NewBounceInterpolator()
	default:
		return NewLinearInterpolator()
	}
}

// Keyframe 关键帧
type Keyframe struct {
	Time     float64
	Position gmMath.Vector2
	Value    float64
}

// KeyframeInterpolator 关键帧插值器
type KeyframeInterpolator struct {
	keyframes    []Keyframe
	interpolator Interpolator
}

// NewKeyframeInterpolator 创建关键帧插值器
func NewKeyframeInterpolator(keyframes []Keyframe, interpType InterpolationType) *KeyframeInterpolator {
	return &KeyframeInterpolator{
		keyframes:    keyframes,
		interpolator: GetInterpolator(interpType),
	}
}

// InterpolateAt 在指定时间点进行插值
func (ki *KeyframeInterpolator) InterpolateAt(t float64) (gmMath.Vector2, float64) {
	if len(ki.keyframes) == 0 {
		return gmMath.Vector2{X: 0, Y: 0}, 0
	}

	if len(ki.keyframes) == 1 {
		return ki.keyframes[0].Position, ki.keyframes[0].Value
	}

	// 找到对应的两个关键帧
	for i := 0; i < len(ki.keyframes)-1; i++ {
		if t >= ki.keyframes[i].Time && t <= ki.keyframes[i+1].Time {
			start := ki.keyframes[i]
			end := ki.keyframes[i+1]
			
			// 计算在两个关键帧之间的相对时间
			relativeT := (t - start.Time) / (end.Time - start.Time)
			
			// 使用插值器进行插值
			position := ki.interpolator.Interpolate(start.Position, end.Position, relativeT)
			value := ki.interpolator.InterpolateFloat(start.Value, end.Value, relativeT)
			
			return position, value
		}
	}

	// 如果时间超出范围，返回最后一个关键帧的值
	last := ki.keyframes[len(ki.keyframes)-1]
	return last.Position, last.Value
}