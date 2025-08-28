package scene

import (
	"render2go/animation"
	"render2go/core"
	"render2go/renderer"
	"time"
)

// Scene 场景，类似于Manim中的Scene
type Scene struct {
	objects     []core.Mobject
	animations  []animation.Animation
	renderer    renderer.Renderer
	width       int
	height      int
	background  [3]float64 // RGB background color
	isPlaying   bool
	startTime   time.Time
	currentTime time.Duration
}

// NewScene 创建新场景
func NewScene(width, height int) *Scene {
	return &Scene{
		objects:    make([]core.Mobject, 0),
		animations: make([]animation.Animation, 0),
		width:      width,
		height:     height,
		background: [3]float64{0.0, 0.0, 0.0}, // 黑色背景
		isPlaying:  false,
	}
}

// SetRenderer 设置渲染器
func (s *Scene) SetRenderer(r renderer.Renderer) {
	s.renderer = r
}

// Add 添加对象到场景
func (s *Scene) Add(objects ...core.Mobject) {
	s.objects = append(s.objects, objects...)
}

// Remove 从场景中移除对象
func (s *Scene) Remove(object core.Mobject) {
	for i, obj := range s.objects {
		if obj == object {
			s.objects = append(s.objects[:i], s.objects[i+1:]...)
			break
		}
	}
}

// Clear 清空场景
func (s *Scene) Clear() {
	s.objects = s.objects[:0]
	s.animations = s.animations[:0]
}

// Play 播放动画
func (s *Scene) Play(anims ...animation.Animation) {
	s.animations = append(s.animations, anims...)

	if !s.isPlaying {
		s.isPlaying = true
		s.startTime = time.Now()
		s.runAnimationLoop()
	}
}

// Wait 等待指定时间
func (s *Scene) Wait(duration time.Duration) {
	// 创建一个空动画来实现等待
	waitAnim := &WaitAnimation{
		BaseAnimation: animation.NewBaseAnimation(nil, duration),
	}
	s.Play(waitAnim)
}

// WaitAnimation 等待动画
type WaitAnimation struct {
	*animation.BaseAnimation
}

func (w *WaitAnimation) Update(progress float64) {
	if progress >= 1.0 {
		w.BaseAnimation.SetFinished(true)
	}
}

// runAnimationLoop 运行动画循环
func (s *Scene) runAnimationLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for s.isPlaying {
		select {
		case <-ticker.C:
			s.update()
			s.render()
		}

		// 检查是否所有动画都完成了
		allFinished := true
		for _, anim := range s.animations {
			if !anim.IsFinished() {
				allFinished = false
				break
			}
		}

		if allFinished {
			s.isPlaying = false
			s.animations = s.animations[:0] // 清空完成的动画
		}
	}
}

// update 更新场景状态
func (s *Scene) update() {
	s.currentTime = time.Since(s.startTime)

	for _, anim := range s.animations {
		if !anim.IsFinished() {
			progress := float64(s.currentTime) / float64(anim.GetDuration())
			if progress > 1.0 {
				progress = 1.0
			}
			anim.Update(progress)
		}
	}
}

// render 渲染场景
func (s *Scene) render() {
	if s.renderer != nil {
		s.renderer.Clear(s.background[0], s.background[1], s.background[2])

		for _, obj := range s.objects {
			s.renderer.Render(obj)
		}

		s.renderer.Present()
	}
}

// RenderFrame 公共渲染方法
func (s *Scene) RenderFrame() {
	s.render()
}

// SetBackground 设置背景色
func (s *Scene) SetBackground(r, g, b float64) {
	s.background = [3]float64{r, g, b}
}

// GetObjects 获取场景中的所有对象
func (s *Scene) GetObjects() []core.Mobject {
	return s.objects
}

// GetWidth 获取场景宽度
func (s *Scene) GetWidth() int {
	return s.width
}

// GetHeight 获取场景高度
func (s *Scene) GetHeight() int {
	return s.height
}

// Construct 构造场景内容，子类应该重写此方法
func (s *Scene) Construct() {
	// 默认为空，子类重写
}

// SaveFrame 保存当前帧
func (s *Scene) SaveFrame(filename string) error {
	if s.renderer != nil {
		// 使用简单的文件路径保存
		return s.renderer.SaveFrame(filename)
	}
	return nil
}

// CreateAnimation 创建动画的便捷方法
func (s *Scene) CreateAnimation() *AnimationBuilder {
	return &AnimationBuilder{scene: s}
}

// AnimationBuilder 动画构建器
type AnimationBuilder struct {
	scene      *Scene
	animations []animation.Animation
}

// MoveTo 添加移动动画
func (ab *AnimationBuilder) MoveTo(object core.Mobject, target [2]float64, duration time.Duration) *AnimationBuilder {
	anim := animation.NewMoveToAnimation(object,
		struct{ X, Y float64 }{X: target[0], Y: target[1]}, duration)
	ab.animations = append(ab.animations, anim)
	return ab
}

// Scale 添加缩放动画
func (ab *AnimationBuilder) Scale(object core.Mobject, factor float64, duration time.Duration) *AnimationBuilder {
	anim := animation.NewScaleAnimation(object, factor, duration)
	ab.animations = append(ab.animations, anim)
	return ab
}

// Rotate 添加旋转动画
func (ab *AnimationBuilder) Rotate(object core.Mobject, angle float64, duration time.Duration) *AnimationBuilder {
	anim := animation.NewRotateAnimation(object, angle, duration)
	ab.animations = append(ab.animations, anim)
	return ab
}

// FadeIn 添加淡入动画
func (ab *AnimationBuilder) FadeIn(object core.Mobject, duration time.Duration) *AnimationBuilder {
	anim := animation.NewFadeInAnimation(object, duration)
	ab.animations = append(ab.animations, anim)
	return ab
}

// FadeOut 添加淡出动画
func (ab *AnimationBuilder) FadeOut(object core.Mobject, duration time.Duration) *AnimationBuilder {
	anim := animation.NewFadeOutAnimation(object, duration)
	ab.animations = append(ab.animations, anim)
	return ab
}

// Build 构建动画组
func (ab *AnimationBuilder) Build() animation.Animation {
	if len(ab.animations) == 1 {
		return ab.animations[0]
	}
	return animation.NewAnimationGroup(ab.animations...)
}

// Play 构建并播放动画
func (ab *AnimationBuilder) Play() {
	anim := ab.Build()
	ab.scene.Play(anim)
}

// SaveFrameWithTimestamp 保存当前帧（自动添加时间戳）
func (s *Scene) SaveFrameWithTimestamp(prefix string) error {
	timestamp := time.Now().Format("20060102_150405_000")
	filename := prefix + "_" + timestamp + ".png"
	return s.SaveFrame(filename)
}
