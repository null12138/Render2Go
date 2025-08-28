package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("Render2Go 动画引擎构建工具")
	fmt.Println("========================")

	// 检查Go环境
	fmt.Println("\n1. 检查Go环境...")
	if !checkGoInstallation() {
		fmt.Println("错误: 未找到Go环境，请先安装Go")
		return
	}
	fmt.Println("✓ Go环境检查通过")

	// 整理依赖
	fmt.Println("\n2. 整理项目依赖...")
	if err := runCommand("go", "mod", "tidy"); err != nil {
		fmt.Printf("警告: 依赖整理失败: %v\n", err)
		fmt.Println("这可能是网络问题，可以稍后重试")
	} else {
		fmt.Println("✓ 依赖整理完成")
	}

	// 运行核心演示
	fmt.Println("\n3. 运行核心功能演示...")
	fmt.Println("=====================================")
	if err := runCommand("go", "run", "cmd/demo/main.go"); err != nil {
		fmt.Printf("演示运行失败: %v\n", err)
	}

	// 运行完整示例
	fmt.Println("\n4. 尝试运行完整图形示例...")
	fmt.Println("=====================================")
	fmt.Println("注意：如果出现图形库错误，这是正常的")
	if err := runCommand("go", "run", "examples/basic_example.go"); err != nil {
		fmt.Printf("完整示例运行失败: %v\n", err)
		fmt.Println("这通常是因为缺少图形依赖库")
	}

	// 显示后续步骤
	showNextSteps()
}

func checkGoInstallation() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func showNextSteps() {
	fmt.Println("\n=====================================")
	fmt.Println("🎉 快速开始完成！")
	fmt.Println("\n接下来你可以：")
	fmt.Println("1. 📖 查看 README.md 了解详细文档")
	fmt.Println("2. 🚀 查看 QUICKSTART.md 获取快速指南")
	fmt.Println("3. 💡 查看 examples/ 目录下的示例代码")
	fmt.Println("4. 🎨 安装完整图形依赖来使用渲染功能：")
	fmt.Println("   go get github.com/fogleman/gg")
	fmt.Println("   go get github.com/golang/freetype")
	fmt.Println("   go get golang.org/x/image")
	fmt.Println("\n🎬 开始创建你的动画吧！")

	if runtime.GOOS == "windows" {
		fmt.Println("\n按任意键退出...")
		fmt.Scanln()
	}
}
