package main

import (
	"fmt"
	"time"

	"render2go/animation"
	"render2go/colors"
	"render2go/geometry"
	gmMath "render2go/math"
	"render2go/renderer"
	"render2go/scene"
)

// ColorSchemeScene 配色方案展示场景
type ColorSchemeScene struct {
	*scene.Scene
}

// NewColorSchemeScene 创建配色方案展示场景
func NewColorSchemeScene() *ColorSchemeScene {
	s := &ColorSchemeScene{
		Scene: scene.NewScene(1000, 700),
	}

	// 设置渲染器
	canvasRenderer := renderer.NewCanvasRenderer(1000, 700)
	s.SetRenderer(canvasRenderer)

	return s
}

// Construct 构造场景内容
func (s *ColorSchemeScene) Construct() {
	// 设置深蓝色背景
	r, g, b, _ := colors.RGBAToFloat64(colors.DeepBlue)
	s.SetBackground(r, g, b)

	fmt.Println("构建海洋蓝配色方案展示场景...")

	// 创建标题文本
	title := geometry.NewText("Render2Go 专业蓝配色方案", 36)
	title.SetColor(colors.LightPurple)
	title.MoveTo(gmMath.Vector2{X: 0, Y: 250})

	// 创建配色方案说明
	subtitle := geometry.NewText("Professional Blue Color Scheme", 24)
	subtitle.SetColor(colors.LightPurple)
	subtitle.MoveTo(gmMath.Vector2{X: 0, Y: 200})

	// 创建5个圆形展示配色
	colorNames := []string{"深海蓝", "海洋蓝", "天空蓝", "浅蓝色", "冰蓝色"}
	colorValues := []string{"#0a2639", "#196090", "#3498db", "#8bc4ea", "#d4e9f7"}
	circles := make([]*geometry.Circle, 5)
	labels := make([]*geometry.Text, 5)

	for i := 0; i < 5; i++ {
		// 创建圆形
		circle := geometry.NewCircle(40)
		circle.SetColor(colors.ProfessionalBlue.GetColorByIndex(i))
		x := float64(-200 + i*100)
		circle.MoveTo(gmMath.Vector2{X: x, Y: 50})
		circles[i] = circle

		// 创建标签
		label := geometry.NewText(fmt.Sprintf("%s\n%s", colorNames[i], colorValues[i]), 14)
		label.SetColor(colors.LightPurple)
		label.MoveTo(gmMath.Vector2{X: x, Y: -20})
		labels[i] = label
	}

	// 创建几何形状展示
	rect := geometry.NewRectangle(150, 80)
	rect.SetColor(colors.MidBlue)
	rect.MoveTo(gmMath.Vector2{X: -200, Y: -150})

	triangle := geometry.NewRegularPolygon(3, 50)
	triangle.SetColor(colors.CyanBlue)
	triangle.MoveTo(gmMath.Vector2{X: -50, Y: -150})

	pentagon := geometry.NewRegularPolygon(5, 45)
	pentagon.SetColor(colors.LightPurple)
	pentagon.MoveTo(gmMath.Vector2{X: 100, Y: -150})

	hexagon := geometry.NewRegularPolygon(6, 40)
	hexagon.SetColor(colors.LightPurple)
	hexagon.MoveTo(gmMath.Vector2{X: 250, Y: -150})

	// 添加所有对象到场景
	s.Add(title, subtitle)
	for i := 0; i < 5; i++ {
		s.Add(circles[i], labels[i])
	}
	s.Add(rect, triangle, pentagon, hexagon)

	// 创建动画序列
	fmt.Println("开始播放配色展示动画...")

	// 1. 标题淡入
	titleFade := animation.NewFadeInAnimation(title, 1*time.Second)
	subtitleFade := animation.NewFadeInAnimation(subtitle, 1*time.Second)
	s.Play(titleFade, subtitleFade)
	s.Wait(500 * time.Millisecond)

	// 2. 圆形依次出现
	for i := 0; i < 5; i++ {
		circleAnim := animation.NewFadeInAnimation(circles[i], 500*time.Millisecond)
		labelAnim := animation.NewFadeInAnimation(labels[i], 500*time.Millisecond)
		s.Play(circleAnim, labelAnim)
		s.Wait(200 * time.Millisecond)
	}

	s.Wait(1 * time.Second)

	// 3. 几何形状动画展示
	rectMove := animation.NewMoveToAnimation(rect, gmMath.Vector2{X: -100, Y: -150}, 1*time.Second)
	triRotate := animation.NewRotateAnimation(triangle, 2*3.14159, 2*time.Second)
	pentScale := animation.NewScaleAnimation(pentagon, 1.5, 1500*time.Millisecond)
	hexMove := animation.NewMoveToAnimation(hexagon, gmMath.Vector2{X: 150, Y: -150}, 1*time.Second)

	s.Play(rectMove, triRotate, pentScale, hexMove)
	s.Wait(2 * time.Second)

	fmt.Println("配色方案展示完成！")
}

func main() {
	fmt.Println("Render2Go 海洋蓝配色方案展示")
	fmt.Println("==========================")
	fmt.Println("配色方案: #0a2639, #196090, #3498db, #8bc4ea, #d4e9f7")
	fmt.Println()

	// 创建并运行配色展示场景
	colorScene := NewColorSchemeScene()
	colorScene.Construct()

	// 保存最终帧
	err := colorScene.SaveFrame("color_scheme_demo.png")
	if err != nil {
		fmt.Printf("保存帧时出错: %v\n", err)
	} else {
		fmt.Println("配色方案展示图已保存为 color_scheme_demo.png")
	}

	fmt.Println()
	fmt.Println("配色方案详情:")
	fmt.Println("- 深海蓝 (#0a2639): 主要背景色，营造深邃感")
	fmt.Println("- 海洋蓝 (#196090): 主要元素色，体现专业感")
	fmt.Println("- 天空蓝 (#3498db): 强调色，突出重要内容")
	fmt.Println("- 浅蓝色 (#8bc4ea): 辅助色，用于装饰元素")
	fmt.Println("- 冰蓝色 (#d4e9f7): 文字色，确保良好可读性")
	fmt.Println()
	fmt.Println("🎨 海洋蓝配色方案展示完成！")
}
