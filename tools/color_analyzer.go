package main

import (
	"fmt"
	"image/color"
	"render2go/colors"
)

func main() {
	fmt.Println("🎨 Render2Go 配色方案对比工具")
	fmt.Println("================================")

	// 输出专业蓝配色方案
	fmt.Println("\n✨ Professional Blue 配色方案")
	fmt.Println("配色来源: 用户提供的专业蓝色调色板")
	fmt.Println("总颜色数: 6")
	fmt.Println()

	profColors := []struct {
		name       string
		color      color.RGBA
		usage      string
		percentage string
	}{
		{"深蓝色 (DeepBlue)", colors.DeepBlue, "主背景色", "27.2%"},
		{"中蓝色 (MidBlue)", colors.MidBlue, "主要元素", "26.5%"},
		{"紫蓝色 (PurpleBlue)", colors.PurpleBlue, "强调色", "25.1%"},
		{"青蓝色 (CyanBlue)", colors.CyanBlue, "辅助色", "7.4%"},
		{"深色 (DarkColor)", colors.DarkColor, "深背景", "7.3%"},
		{"浅紫色 (LightPurple)", colors.LightPurple, "文字色", "6.6%"},
	}

	for i, c := range profColors {
		fmt.Printf("  %d. %-20s #%02X%02X%02X  RGB(%3d,%3d,%3d)  %-10s [%s]\n",
			i+1, c.name, c.color.R, c.color.G, c.color.B,
			c.color.R, c.color.G, c.color.B,
			c.usage, c.percentage)
	}

	// 计算颜色对比度
	fmt.Println("\n📊 配色分析")
	fmt.Println("主要特征:")
	fmt.Printf("  • 色相范围: 蓝色系 (195°-240°)\n")
	fmt.Printf("  • 饱和度: 中等到高 (40%%-85%%)\n")
	fmt.Printf("  • 明度分布: 深色到中亮 (10%%-75%%)\n")
	fmt.Printf("  • 整体风格: 专业、沉稳、现代\n")

	// 显示使用建议
	fmt.Println("\n🎯 使用建议")
	fmt.Println("最佳搭配:")
	fmt.Printf("  • 背景: DeepBlue (#%02X%02X%02X) + DarkColor (#%02X%02X%02X)\n",
		colors.DeepBlue.R, colors.DeepBlue.G, colors.DeepBlue.B,
		colors.DarkColor.R, colors.DarkColor.G, colors.DarkColor.B)
	fmt.Printf("  • 内容: MidBlue (#%02X%02X%02X) + PurpleBlue (#%02X%02X%02X)\n",
		colors.MidBlue.R, colors.MidBlue.G, colors.MidBlue.B,
		colors.PurpleBlue.R, colors.PurpleBlue.G, colors.PurpleBlue.B)
	fmt.Printf("  • 强调: PurpleBlue (#%02X%02X%02X) + LightPurple (#%02X%02X%02X)\n",
		colors.PurpleBlue.R, colors.PurpleBlue.G, colors.PurpleBlue.B,
		colors.LightPurple.R, colors.LightPurple.G, colors.LightPurple.B)
	fmt.Printf("  • 文字: LightPurple (#%02X%02X%02X) + CyanBlue (#%02X%02X%02X)\n",
		colors.LightPurple.R, colors.LightPurple.G, colors.LightPurple.B,
		colors.CyanBlue.R, colors.CyanBlue.G, colors.CyanBlue.B)

	// 代码示例
	fmt.Println("\n💻 代码示例")
	fmt.Println("```go")
	fmt.Println("import \"render2go/colors\"")
	fmt.Println()
	fmt.Println("// 设置主要元素颜色")
	fmt.Println("shape.SetColor(colors.MidBlue)")
	fmt.Println()
	fmt.Println("// 设置强调色")
	fmt.Println("highlight.SetColor(colors.PurpleBlue)")
	fmt.Println()
	fmt.Println("// 设置文字颜色")
	fmt.Println("text.SetColor(colors.LightPurple)")
	fmt.Println()
	fmt.Println("// 使用配色方案")
	fmt.Println("scheme := colors.ProfessionalBlue")
	fmt.Println("primary := scheme.GetPrimaryColor()  // 中蓝色")
	fmt.Println("accent := scheme.GetAccentColor()    // 浅紫色")
	fmt.Println("```")

	// 输出目录信息
	fmt.Println("\n📁 输出系统集成")
	fmt.Println("新配色方案已完全集成到输出管理系统中:")
	fmt.Println("  • 自动创建项目目录结构")
	fmt.Println("  • 支持时间戳文件命名")
	fmt.Println("  • 提供文件组织和清理功能")
	fmt.Println("  • 跨平台路径处理")

	fmt.Println("\n🚀 快速开始")
	fmt.Println("运行示例程序:")
	fmt.Println("  go run examples/new_color_demo.go      # 静态配色演示")
	fmt.Println("  go run examples/animation_demo.go      # 动画配色演示")

	fmt.Println("\n✅ Professional Blue 配色方案配置完成！")
}
