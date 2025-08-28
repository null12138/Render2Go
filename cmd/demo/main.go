package main

import (
	"fmt"
	"math"
	"time"

	"render2go/animation"
	"render2go/colors"
	"render2go/geometry"
	gmMath "render2go/math"
)

func demonstrateMath() {
	fmt.Println("=== 数学工具演示 ===")

	// 向量运算
	v1 := gmMath.NewVector2(10, 20)
	v2 := gmMath.NewVector2(5, 15)

	fmt.Printf("向量1: (%.1f, %.1f)\n", v1.X, v1.Y)
	fmt.Printf("向量2: (%.1f, %.1f)\n", v2.X, v2.Y)

	sum := v1.Add(v2)
	fmt.Printf("向量加法: (%.1f, %.1f)\n", sum.X, sum.Y)

	diff := v1.Sub(v2)
	fmt.Printf("向量减法: (%.1f, %.1f)\n", diff.X, diff.Y)

	length := v1.Length()
	fmt.Printf("向量1长度: %.2f\n", length)

	normalized := v1.Normalize()
	fmt.Printf("向量1标准化: (%.3f, %.3f)\n", normalized.X, normalized.Y)

	dot := v1.Dot(v2)
	fmt.Printf("点积: %.1f\n", dot)

	// 插值演示
	fmt.Println("\n插值函数演示:")
	for t := 0.0; t <= 1.0; t += 0.25 {
		linear := gmMath.Interpolate(0, 100, t)
		smooth := gmMath.SmoothStep(t) * 100
		eased := gmMath.EaseInOut(t) * 100

		fmt.Printf("t=%.2f: 线性=%.1f, 平滑=%.1f, 缓动=%.1f\n",
			t, linear, smooth, eased)
	}
}

func demonstrateColors() {
	fmt.Println("\n=== 配色方案演示 ===")

	// 显示主要配色
	fmt.Printf("配色方案: %s\n", colors.ProfessionalBlue.Name)
	fmt.Println("颜色详情:")

	colorNames := []string{"深海蓝", "海洋蓝", "天空蓝", "浅蓝色", "冰蓝色"}
	colorHex := []string{"#0a2639", "#196090", "#3498db", "#8bc4ea", "#d4e9f7"}

	for i, colorName := range colorNames {
		color := colors.ProfessionalBlue.GetColorByIndex(i)
		r, g, b, a := colors.RGBAToFloat64(color)
		fmt.Printf("  %s (%s): RGBA(%.0f, %.0f, %.0f, %.0f)\n",
			colorName, colorHex[i], r*255, g*255, b*255, a*255)
	}

	// 演示渐变
	fmt.Println("\n渐变演示:")
	gradient := colors.CreateGradient(colors.DeepBlue, colors.LightPurple, 5)
	for i, gradColor := range gradient {
		r, g, b, _ := colors.RGBAToFloat64(gradColor)
		fmt.Printf("  渐变步骤%d: RGB(%.0f, %.0f, %.0f)\n", i+1, r*255, g*255, b*255)
	}

	// 演示颜色转换
	fmt.Println("\n颜色转换演示:")
	testColor := colors.CyanBlue
	r, g, b, a := colors.RGBAToFloat64(testColor)
	fmt.Printf("天空蓝 -> 浮点数: (%.3f, %.3f, %.3f, %.3f)\n", r, g, b, a)

	backToRGBA := colors.Float64ToRGBA(r, g, b, a)
	fmt.Printf("浮点数 -> RGBA: (%d, %d, %d, %d)\n",
		backToRGBA.R, backToRGBA.G, backToRGBA.B, backToRGBA.A)
}

func demonstrateGeometry() {
	fmt.Println("\n=== 几何对象演示 ===")
	fmt.Println("使用海洋蓝配色方案:")

	// 显示配色方案
	fmt.Printf("配色方案: %s\n", colors.ProfessionalBlue.Name)
	fmt.Printf("  深海蓝: #0a2639\n")
	fmt.Printf("  海洋蓝: #196090\n")
	fmt.Printf("  天空蓝: #3498db\n")
	fmt.Printf("  浅蓝色: #8bc4ea\n")
	fmt.Printf("  冰蓝色: #d4e9f7\n")

	// 创建圆形
	circle := geometry.NewCircle(50)
	circle.SetColor(colors.MidBlue) // 海洋蓝
	circle.MoveTo(gmMath.Vector2{X: 100, Y: 50})

	center := circle.GetCenter()
	fmt.Printf("圆形: 半径=%.1f, 中心=(%.1f, %.1f), 颜色=海洋蓝\n",
		circle.GetRadius(), center.X, center.Y)

	// 创建矩形
	rect := geometry.NewRectangle(100, 60)
	rect.SetColor(colors.CyanBlue) // 天空蓝
	rect.MoveTo(gmMath.Vector2{X: -50, Y: 25})

	rectCenter := rect.GetCenter()
	fmt.Printf("矩形: 中心=(%.1f, %.1f), 颜色=天空蓝\n", rectCenter.X, rectCenter.Y)

	// 创建文本
	text := geometry.NewText("Hello Render2Go!", 24)
	text.SetColor(colors.LightPurple) // 冰蓝
	text.MoveTo(gmMath.Vector2{X: 0, Y: -100})

	textCenter := text.GetCenter()
	fmt.Printf("文本: '%s', 大小=%.0f, 位置=(%.1f, %.1f), 颜色=冰蓝\n",
		text.GetText(), text.GetSize(), textCenter.X, textCenter.Y)

	// 创建正多边形
	pentagon := geometry.NewRegularPolygon(5, 40)
	pentagon.SetColor(colors.LightPurple) // 浅蓝
	pentagon.MoveTo(gmMath.Vector2{X: 150, Y: -50})

	pentCenter := pentagon.GetCenter()
	fmt.Printf("正五边形: 中心=(%.1f, %.1f), 颜色=浅蓝\n", pentCenter.X, pentCenter.Y)

	// 演示变换
	fmt.Println("\n变换演示:")
	fmt.Printf("圆形缩放前中心: (%.1f, %.1f)\n", circle.GetCenter().X, circle.GetCenter().Y)
	circle.Scale(2.0)
	fmt.Printf("圆形缩放2倍后中心: (%.1f, %.1f)\n", circle.GetCenter().X, circle.GetCenter().Y)

	fmt.Printf("矩形旋转前中心: (%.1f, %.1f)\n", rect.GetCenter().X, rect.GetCenter().Y)
	rect.Rotate(math.Pi / 4) // 45度
	fmt.Printf("矩形旋转45度后中心: (%.1f, %.1f)\n", rect.GetCenter().X, rect.GetCenter().Y)
}

func simulateAnimation(anim animation.Animation, name string) {
	fmt.Printf("\n模拟 %s 动画 (时长: %v):\n", name, anim.GetDuration())

	steps := 5
	for i := 0; i <= steps; i++ {
		progress := float64(i) / float64(steps)
		anim.Update(progress)

		if target := anim.GetTarget(); target != nil {
			center := target.GetCenter()
			fmt.Printf("  进度 %3.0f%%: 位置(%.1f, %.1f)\n",
				progress*100, center.X, center.Y)
		} else {
			fmt.Printf("  进度 %3.0f%%\n", progress*100)
		}

		// 模拟时间延迟
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf("%s 动画完成\n", name)
}

func demonstrateAnimation() {
	fmt.Println("\n=== 动画系统演示 ===")
	fmt.Println("使用海洋蓝配色的动画对象:")

	// 创建测试对象，应用配色方案
	circle := geometry.NewCircle(30)
	circle.SetColor(colors.MidBlue) // 海洋蓝
	circle.MoveTo(gmMath.Vector2{X: -100, Y: 0})
	fmt.Println("创建海洋蓝圆形")

	rect := geometry.NewRectangle(60, 40)
	rect.SetColor(colors.CyanBlue) // 天空蓝
	rect.MoveTo(gmMath.Vector2{X: 100, Y: 0})
	fmt.Println("创建天空蓝矩形")

	// 1. 移动动画
	moveAnim := animation.NewMoveToAnimation(
		circle,
		gmMath.Vector2{X: 150, Y: 100},
		2*time.Second,
	)
	simulateAnimation(moveAnim, "海洋蓝圆形移动")

	// 2. 缩放动画
	scaleAnim := animation.NewScaleAnimation(rect, 2.0, 1500*time.Millisecond)
	simulateAnimation(scaleAnim, "天空蓝矩形缩放")

	// 3. 旋转动画
	rotateAnim := animation.NewRotateAnimation(circle, math.Pi, 1*time.Second)
	simulateAnimation(rotateAnim, "海洋蓝圆形旋转")

	// 4. 淡入动画
	fadeAnim := animation.NewFadeInAnimation(rect, 1*time.Second)
	fmt.Printf("\n模拟天空蓝矩形淡入动画:\n")
	for i := 0; i <= 5; i++ {
		progress := float64(i) / 5.0
		fadeAnim.Update(progress)
		opacity := rect.GetFillOpacity()
		fmt.Printf("  进度 %3.0f%%: 透明度=%.2f\n", progress*100, opacity)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("淡入动画完成")

	// 5. 动画组
	fmt.Println("\n演示动画组合:")
	group := animation.NewAnimationGroup(
		animation.NewMoveToAnimation(circle, gmMath.Vector2{X: 0, Y: 0}, 1*time.Second),
		animation.NewScaleAnimation(rect, 0.5, 1*time.Second),
	)

	fmt.Printf("动画组时长: %v\n", group.GetDuration())
	for i := 0; i <= 3; i++ {
		progress := float64(i) / 3.0
		group.Update(progress)
		fmt.Printf("  组合进度 %3.0f%%\n", progress*100)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("动画组完成")
}

func main() {
	fmt.Println("Render2Go 动画引擎核心功能演示")
	fmt.Println("===============================")
	fmt.Println("这是一个简化的演示版本，展示核心功能而不需要图形库依赖")
	fmt.Println()

	// 演示各个模块
	demonstrateColors() // 新增：配色方案演示
	demonstrateMath()
	demonstrateGeometry()
	demonstrateAnimation()

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("🎨 海洋蓝配色方案应用完成！")
	fmt.Println("核心功能演示已完成！")
	fmt.Println("\n要使用完整的图形渲染功能，请安装以下依赖:")
	fmt.Println("  go get github.com/fogleman/gg")
	fmt.Println("  go get github.com/golang/freetype")
	fmt.Println("  go get golang.org/x/image")
	fmt.Println("\n然后运行: go run examples/basic_example.go")
}
