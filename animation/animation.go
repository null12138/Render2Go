package animation

import (
	"render2go/core"
	gmMath "render2go/math"
	"time"
)

// AnimationType 动画类型
type AnimationType int

const (
	AnimationTypeTransform AnimationType = iota
	AnimationTypeColor
	AnimationTypeOpacity
	AnimationTypeCustom
)

// EasingFunction 缓动函数类型
type EasingFunction func(float64) float64

// Animation 动画接口
type Animation interface {
	Update(progress float64)
	GetDuration() time.Duration
	GetTarget() core.Mobject
	IsFinished() bool
	Reset()
}

// BaseAnimation 基础动画
type BaseAnimation struct {
	target     core.Mobject
	duration   time.Duration
	easingFunc EasingFunction
	progress   float64
	finished   bool
	startTime  time.Time
}

// NewBaseAnimation 创建基础动画
func NewBaseAnimation(target core.Mobject, duration time.Duration) *BaseAnimation {
	return &BaseAnimation{
		target:     target,
		duration:   duration,
		easingFunc: gmMath.SmoothStep,
		progress:   0,
		finished:   false,
	}
}

func (a *BaseAnimation) GetDuration() time.Duration {
	return a.duration
}

func (a *BaseAnimation) GetTarget() core.Mobject {
	return a.target
}

func (a *BaseAnimation) IsFinished() bool {
	return a.finished
}

func (a *BaseAnimation) Reset() {
	a.progress = 0
	a.finished = false
	a.startTime = time.Now()
}

func (a *BaseAnimation) SetEasing(easing EasingFunction) {
	a.easingFunc = easing
}

func (a *BaseAnimation) SetFinished(finished bool) {
	a.finished = finished
}

// MoveToAnimation 移动动画
type MoveToAnimation struct {
	*BaseAnimation
	startPos gmMath.Vector2
	endPos   gmMath.Vector2
}

// NewMoveToAnimation 创建移动动画
func NewMoveToAnimation(target core.Mobject, endPos gmMath.Vector2, duration time.Duration) *MoveToAnimation {
	return &MoveToAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		startPos:      target.GetCenter(),
		endPos:        endPos,
	}
}

func (a *MoveToAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	easedProgress := a.easingFunc(progress)
	currentPos := gmMath.Vector2{
		X: gmMath.Interpolate(a.startPos.X, a.endPos.X, easedProgress),
		Y: gmMath.Interpolate(a.startPos.Y, a.endPos.Y, easedProgress),
	}

	a.target.MoveTo(currentPos)
	a.progress = progress
}

// ScaleAnimation 缩放动画
type ScaleAnimation struct {
	*BaseAnimation
	startScale    float64
	endScale      float64
	initialPoints []gmMath.Vector2
}

// NewScaleAnimation 创建缩放动画
func NewScaleAnimation(target core.Mobject, endScale float64, duration time.Duration) *ScaleAnimation {
	return &ScaleAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		startScale:    1.0,
		endScale:      endScale,
		initialPoints: target.GetPoints(),
	}
}

func (a *ScaleAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	easedProgress := a.easingFunc(progress)
	currentScale := gmMath.Interpolate(a.startScale, a.endScale, easedProgress)

	// 重置到初始状态然后应用缩放
	a.target.SetPoints(a.initialPoints)
	a.target.Scale(currentScale)
	a.progress = progress
}

// RotateAnimation 旋转动画
type RotateAnimation struct {
	*BaseAnimation
	startAngle    float64
	endAngle      float64
	initialPoints []gmMath.Vector2
}

// NewRotateAnimation 创建旋转动画
func NewRotateAnimation(target core.Mobject, angle float64, duration time.Duration) *RotateAnimation {
	return &RotateAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		startAngle:    0,
		endAngle:      angle,
		initialPoints: target.GetPoints(),
	}
}

func (a *RotateAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	easedProgress := a.easingFunc(progress)
	currentAngle := gmMath.Interpolate(a.startAngle, a.endAngle, easedProgress)

	// 重置到初始状态然后应用旋转
	a.target.SetPoints(a.initialPoints)
	a.target.Rotate(currentAngle)
	a.progress = progress
}

// FadeInAnimation 淡入动画
type FadeInAnimation struct {
	*BaseAnimation
	startOpacity float64
	endOpacity   float64
}

// NewFadeInAnimation 创建淡入动画
func NewFadeInAnimation(target core.Mobject, duration time.Duration) *FadeInAnimation {
	return &FadeInAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		startOpacity:  0.0,
		endOpacity:    1.0,
	}
}

func (a *FadeInAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	easedProgress := a.easingFunc(progress)
	currentOpacity := gmMath.Interpolate(a.startOpacity, a.endOpacity, easedProgress)
	a.target.SetFillOpacity(currentOpacity)
	a.progress = progress
}

// FadeOutAnimation 淡出动画
type FadeOutAnimation struct {
	*BaseAnimation
	startOpacity float64
	endOpacity   float64
}

// NewFadeOutAnimation 创建淡出动画
func NewFadeOutAnimation(target core.Mobject, duration time.Duration) *FadeOutAnimation {
	return &FadeOutAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		startOpacity:  target.GetFillOpacity(),
		endOpacity:    0.0,
	}
}

func (a *FadeOutAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	easedProgress := a.easingFunc(progress)
	currentOpacity := gmMath.Interpolate(a.startOpacity, a.endOpacity, easedProgress)
	a.target.SetFillOpacity(currentOpacity)
	a.progress = progress
}

// AnimationGroup 动画组，用于同时播放多个动画
type AnimationGroup struct {
	animations []Animation
	duration   time.Duration
}

// NewAnimationGroup 创建动画组
func NewAnimationGroup(animations ...Animation) *AnimationGroup {
	group := &AnimationGroup{
		animations: animations,
	}

	// 找到最长的动画duration
	for _, anim := range animations {
		if anim.GetDuration() > group.duration {
			group.duration = anim.GetDuration()
		}
	}

	return group
}

func (g *AnimationGroup) Update(progress float64) {
	for _, anim := range g.animations {
		// 计算每个动画的相对进度
		animProgress := progress * float64(g.duration) / float64(anim.GetDuration())
		if animProgress > 1.0 {
			animProgress = 1.0
		}
		anim.Update(animProgress)
	}
}

func (g *AnimationGroup) GetDuration() time.Duration {
	return g.duration
}

func (g *AnimationGroup) IsFinished() bool {
	for _, anim := range g.animations {
		if !anim.IsFinished() {
			return false
		}
	}
	return true
}

func (g *AnimationGroup) Reset() {
	for _, anim := range g.animations {
		anim.Reset()
	}
}

func (g *AnimationGroup) GetTarget() core.Mobject {
	// 对于动画组，返回第一个动画的目标，如果没有动画则返回nil
	if len(g.animations) > 0 {
		return g.animations[0].GetTarget()
	}
	return nil
}
