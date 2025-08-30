package animation

import (
	"image/color"
	"render2go/core"
	gmMath "render2go/math"
	"time"
	"math"
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

// BaseAnimation 基动画
type BaseAnimation struct {
	target         core.Mobject
	duration       time.Duration
	easingFunc     EasingFunction
	interpolation  InterpolationType
	progress       float64
	finished       bool
	startTime      time.Time
}

// NewBaseAnimation 创建基础动画
func NewBaseAnimation(target core.Mobject, duration time.Duration) *BaseAnimation {
	return &BaseAnimation{
		target:        target,
		duration:      duration,
		easingFunc:    gmMath.SmoothStep,
		interpolation: Smooth, // 默认使用平滑插值
		progress:      0,
		finished:      false,
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

func (a *BaseAnimation) SetInterpolation(interp InterpolationType) {
	a.interpolation = interp
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

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	currentPos := interpolator.Interpolate(a.startPos, a.endPos, easedProgress)

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

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	currentScale := interpolator.InterpolateFloat(a.startScale, a.endScale, easedProgress)

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

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	currentAngle := interpolator.InterpolateFloat(a.startAngle, a.endAngle, easedProgress)

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

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	currentOpacity := interpolator.InterpolateFloat(a.startOpacity, a.endOpacity, easedProgress)
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

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	currentOpacity := interpolator.InterpolateFloat(a.startOpacity, a.endOpacity, easedProgress)
	a.target.SetFillOpacity(currentOpacity)
	a.progress = progress
}

// ColorAnimation 颜色变换动画
type ColorAnimation struct {
	*BaseAnimation
	startColor color.RGBA
	endColor   color.RGBA
}

// NewColorAnimation 创建颜色变换动画
func NewColorAnimation(target core.Mobject, endColor color.RGBA, duration time.Duration) *ColorAnimation {
	startColor := color.RGBA{255, 255, 255, 255} // 默认白色
	if c, ok := target.GetColor().(color.RGBA); ok {
		startColor = c
	}
	
	return &ColorAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		startColor:    startColor,
		endColor:      endColor,
	}
}

func (a *ColorAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	
	// 使用插值器插值颜色的各个分量
	r := uint8(interpolator.InterpolateFloat(float64(a.startColor.R), float64(a.endColor.R), easedProgress))
	g := uint8(interpolator.InterpolateFloat(float64(a.startColor.G), float64(a.endColor.G), easedProgress))
	b := uint8(interpolator.InterpolateFloat(float64(a.startColor.B), float64(a.endColor.B), easedProgress))
	alpha := uint8(interpolator.InterpolateFloat(float64(a.startColor.A), float64(a.endColor.A), easedProgress))
	
	newColor := color.RGBA{r, g, b, alpha}
	a.target.SetColor(newColor)
	a.progress = progress
}

// PathAnimation 路径动画
type PathAnimation struct {
	*BaseAnimation
	pathPoints []gmMath.Vector2
}

// NewPathAnimation 创建路径动画
func NewPathAnimation(target core.Mobject, pathPoints []gmMath.Vector2, duration time.Duration) *PathAnimation {
	return &PathAnimation{
		BaseAnimation: NewBaseAnimation(target, duration),
		pathPoints:    pathPoints,
	}
}

func (a *PathAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	easedProgress := a.easingFunc(progress)
	
	// 计算路径上的点
	currentPos := a.getPositionOnPath(easedProgress)
	
	// 使用插值器进一步平滑路径点的位置
	pathStart := a.getPositionOnPath(0.0)
	interpolatedPos := interpolator.Interpolate(pathStart, currentPos, 1.0)
	
	a.target.MoveTo(interpolatedPos)
	a.progress = progress
}

// getPositionOnPath 根据进度获取路径上的位置
func (a *PathAnimation) getPositionOnPath(progress float64) gmMath.Vector2 {
	if len(a.pathPoints) == 0 {
		return gmMath.Vector2{X: 0, Y: 0}
	}
	
	if len(a.pathPoints) == 1 {
		return a.pathPoints[0]
	}
	
	// 计算总路径长度
	totalLength := 0.0
	segmentLengths := make([]float64, len(a.pathPoints)-1)
	for i := 0; i < len(a.pathPoints)-1; i++ {
		length := a.pathPoints[i].Distance(a.pathPoints[i+1])
		segmentLengths[i] = length
		totalLength += length
	}
	
	if totalLength == 0 {
		return a.pathPoints[0]
	}
	
	// 根据进度找到对应的线段
	targetDistance := progress * totalLength
	currentDistance := 0.0
	
	for i := 0; i < len(segmentLengths); i++ {
		if currentDistance+segmentLengths[i] >= targetDistance {
			// 在当前线段上插值
			segmentProgress := (targetDistance - currentDistance) / segmentLengths[i]
			startPoint := a.pathPoints[i]
			endPoint := a.pathPoints[i+1]
			
			return gmMath.Vector2{
				X: gmMath.Interpolate(startPoint.X, endPoint.X, segmentProgress),
				Y: gmMath.Interpolate(startPoint.Y, endPoint.Y, segmentProgress),
			}
		}
		currentDistance += segmentLengths[i]
	}
	
	// 如果超出路径，返回最后一个点
	return a.pathPoints[len(a.pathPoints)-1]
}

// ElasticAnimation 弹性动画
type ElasticAnimation struct {
	*BaseAnimation
	startValue float64
	endValue   float64
	amplitude  float64
	period     float64
	property   string // "scale", "opacity", "x", "y"
	target     core.Mobject
}

// NewElasticAnimation 创建弹性动画
func NewElasticAnimation(target core.Mobject, property string, endValue, duration float64) *ElasticAnimation {
	startValue := 0.0
	
	// 根据属性类型获取起始值
	switch property {
	case "scale":
		startValue = 1.0 // 默认缩放为1.0
	case "opacity":
		startValue = target.GetFillOpacity()
	case "x":
		startValue = target.GetCenter().X
	case "y":
		startValue = target.GetCenter().Y
	}
	
	return &ElasticAnimation{
		BaseAnimation: NewBaseAnimation(target, time.Duration(duration*float64(time.Second))),
		startValue:    startValue,
		endValue:      endValue,
		amplitude:     1.0,
		period:        0.3,
		property:      property,
		target:        target,
	}
}

func (a *ElasticAnimation) Update(progress float64) {
	if progress >= 1.0 {
		progress = 1.0
		a.finished = true
	}

	// 使用插值器进行更流畅的插值
	interpolator := GetInterpolator(a.interpolation)
	
	// 弹性缓动函数
	easedProgress := a.elasticEaseOut(progress)
	currentValue := interpolator.InterpolateFloat(a.startValue, a.endValue, easedProgress)
	
	// 根据属性类型应用值
	switch a.property {
	case "scale":
		// 重置到初始状态然后应用缩放
		points := a.target.GetPoints()
		a.target.SetPoints(points)
		a.target.Scale(currentValue)
	case "opacity":
		a.target.SetFillOpacity(currentValue)
	case "x":
		center := a.target.GetCenter()
		a.target.MoveTo(gmMath.Vector2{X: currentValue, Y: center.Y})
	case "y":
		center := a.target.GetCenter()
		a.target.MoveTo(gmMath.Vector2{X: center.X, Y: currentValue})
	}
	
	a.progress = progress
}

// elasticEaseOut 弹性缓出函数
func (a *ElasticAnimation) elasticEaseOut(t float64) float64 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	
	p := a.period
	s := p / 4
	
	return (a.amplitude * math.Pow(2, -10*t) * math.Sin((t*1-s)*(2*math.Pi)/p) + 1)
}

// BouncingBallAnimation 物理弹跳球动画
type BouncingBallAnimation struct {
	*BaseAnimation
	ball        core.Mobject
	initialPos  gmMath.Vector2
	velocity    gmMath.Vector2
	gravity     float64
	elasticity  float64
	groundLevel float64
	lastUpdate  time.Time
}

// NewBouncingBallAnimation 创建物理弹跳球动画
func NewBouncingBallAnimation(ball core.Mobject, duration time.Duration) *BouncingBallAnimation {
	// 初始化时间
	now := time.Now()
	return &BouncingBallAnimation{
		BaseAnimation: NewBaseAnimation(ball, duration),
		ball:          ball,
		initialPos:    ball.GetCenter(),
		velocity:      gmMath.Vector2{X: 0, Y: 0},
		gravity:       -9.8, // 重力加速度 (向下为负)
		elasticity:    0.8,  // 弹性系数
		groundLevel:   -4.5, // 地面高度
		lastUpdate:    now,
	}
}

func (a *BouncingBallAnimation) Update(progress float64) {
	// 计算时间差
	now := time.Now()
	dt := now.Sub(a.lastUpdate).Seconds()
	a.lastUpdate = now
	
	// 更新速度 (v = v0 + a*t)
	a.velocity.Y += a.gravity * dt
	
	// 更新位置 (s = s0 + v*t)
	currentPos := a.ball.GetCenter()
	newPos := gmMath.Vector2{
		X: currentPos.X + a.velocity.X * dt,
		Y: currentPos.Y + a.velocity.Y * dt,
	}
	
	// 检查是否触地
	if newPos.Y <= a.groundLevel {
		// 触地反弹
		newPos.Y = a.groundLevel
		a.velocity.Y = -a.velocity.Y * a.elasticity // 反弹并损失能量
		
		// 如果速度太小，停止弹跳
		if math.Abs(a.velocity.Y) < 0.1 {
			a.velocity.Y = 0
			a.finished = true
		}
	}
	
	// 使用插值器进行更流畅的位置插值
	interpolator := GetInterpolator(a.interpolation)
	interpolatedPos := interpolator.Interpolate(currentPos, newPos, 1.0)
	
	// 移动球到新位置
	a.ball.MoveTo(interpolatedPos)
	a.progress = progress
	
	// 检查是否完成
	if progress >= 1.0 {
		a.finished = true
	}
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
