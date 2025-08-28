package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"render2go/interpreter"
)

func main() {
	// 命令行参数
	var (
		file        = flag.String("file", "", "Script file to execute")
		interactive = flag.Bool("i", false, "Run in interactive mode")
		debug       = flag.Bool("debug", false, "Enable debug mode")
		help        = flag.Bool("help", false, "Show help information")
		version     = flag.Bool("version", false, "Show version information")
	)

	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Println("Render2Go Script Interpreter v1.0.0")
		fmt.Println("A powerful animation scripting language")
		return
	}

	// 显示帮助信息
	if *help {
		printUsage()
		return
	}

	// 创建解释器
	interp := interpreter.NewInterpreter(*debug)

	// 交互式模式
	if *interactive {
		interp.RunInteractive()
		return
	}

	// 执行文件
	if *file != "" {
		if !fileExists(*file) {
			fmt.Printf("❌ Error: File '%s' does not exist\n", *file)
			os.Exit(1)
		}

		fmt.Printf("🎬 Executing script: %s\n", *file)
		err := interp.RunFile(*file)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Script execution completed successfully!")
		return
	}

	// 检查是否有非标志参数（直接的文件名）
	args := flag.Args()
	if len(args) > 0 {
		filename := args[0]
		if !fileExists(filename) {
			fmt.Printf("❌ Error: File '%s' does not exist\n", filename)
			os.Exit(1)
		}

		fmt.Printf("🎬 Executing script: %s\n", filename)
		err := interp.RunFile(filename)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Script execution completed successfully!")
		return
	}

	// 如果没有指定文件且没有交互模式，查找默认脚本
	defaultFiles := []string{"main.r2g", "script.r2g", "animation.r2g"}
	for _, filename := range defaultFiles {
		if fileExists(filename) {
			fmt.Printf("🎬 Found and executing: %s\n", filename)
			err := interp.RunFile(filename)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Script execution completed successfully!")
			return
		}
	}

	// 没有找到脚本文件，启动交互模式
	fmt.Println("No script file specified. Starting interactive mode...")
	fmt.Println("Use 'render2go --help' for usage information")
	fmt.Println()
	interp.RunInteractive()
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println(`🎬 Render2Go Script Interpreter

USAGE:
    render2go [OPTIONS] [FILE]

OPTIONS:
    -file <file>        Execute the specified script file
    -i                  Run in interactive mode
    -debug              Enable debug mode (shows tokens and AST)
    -help               Show this help message
    -version            Show version information

FILE FORMATS:
    .r2g                Render2Go Animation script files

EXAMPLES:
    render2go script.r2g              # Execute script.r2g
    render2go -file animation.r2g     # Execute animation.r2g
    render2go -i                      # Start interactive mode
    render2go -debug script.r2g       # Execute with debug output

SCRIPT LANGUAGE:
    The Render2Go scripting language supports:
    
    Scene Management:
        scene 800 600 "my_project"
    
    Object Creation:
        create circle c1 50 (400, 300)
        create rectangle r1 100 80 (200, 200)
        create line l1 (0, 0) (100, 100)
        create text t1 "Hello World" 24 (400, 100)
    
    Property Setting:
        set c1.color = #576DA2
        set c1.position = (500, 400)
        set c1.opacity = 0.8
    
    Rendering:
        render
        save "my_frame.png"
    
    Control Flow:
        loop 10 {
            render
            save "frame.png"
        }

DEFAULT BEHAVIOR:
    If no file is specified, render2go will look for these files in order:
    - main.r2g
    - script.r2g
    - animation.r2g
    
    If none are found, interactive mode will start.

For more information, visit: https://github.com/render2go/render2go`)
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	return filepath.Ext(filename)
}
