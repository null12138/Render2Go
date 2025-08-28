package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"render2go/geometry"
	"strings"
)

// Interpreter 主解释器结构
type Interpreter struct {
	evaluator *Evaluator
	debug     bool
}

// NewInterpreter 创建新的解释器实例
func NewInterpreter(debug bool) *Interpreter {
	return &Interpreter{
		evaluator: NewEvaluator(),
		debug:     debug,
	}
}

// RunFile 执行脚本文件
func (i *Interpreter) RunFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	return i.RunReader(file, filename)
}

// RunReader 从Reader执行脚本
func (i *Interpreter) RunReader(reader io.Reader, source string) error {
	// 读取整个输入
	scanner := bufio.NewScanner(reader)
	var content strings.Builder

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read from %s: %w", source, err)
	}

	return i.RunString(content.String(), source)
}

// RunString 直接执行脚本字符串
func (i *Interpreter) RunString(script, source string) error {
	if i.debug {
		fmt.Printf("🔍 Parsing script from %s...\n", source)
	}

	// 词法分析
	lexer := NewLexer(script)

	if i.debug {
		fmt.Println("📝 Tokens:")
		debugLexer := NewLexer(script)
		for {
			token := debugLexer.NextToken()
			if token.Type == TOKEN_EOF {
				break
			}
			if token.Type != TOKEN_NEWLINE {
				fmt.Printf("  %s: %s\n", token.Type, token.Literal)
			}
		}
		fmt.Println()
	}

	// 语法分析
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	// 检查解析错误
	errors := parser.Errors()
	if len(errors) > 0 {
		return fmt.Errorf("parsing errors:\n%s", strings.Join(errors, "\n"))
	}

	if i.debug {
		fmt.Println("🌳 AST:")
		fmt.Println(program.String())
		fmt.Println()
	}

	// 执行程序
	if i.debug {
		fmt.Println("🚀 Executing...")
	}

	err := i.evaluator.Evaluate(program)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	// 检查执行错误
	execErrors := i.evaluator.GetErrors()
	if len(execErrors) > 0 {
		return fmt.Errorf("execution errors:\n%s", strings.Join(execErrors, "\n"))
	}

	if i.debug {
		fmt.Println("✅ Execution completed successfully!")
	}

	// 自动修复PNG文件扩展名
	if i.debug {
		fmt.Println("🔧 Attempting to fix PNG extensions...")
	}
	err = i.fixPNGExtensions()
	if err != nil && i.debug {
		fmt.Printf("⚠️ Warning: Failed to fix PNG extensions: %v\n", err)
	}
	if i.debug {
		fmt.Println("✅ PNG extension fix completed")
	}

	return nil
}

// RunInteractive 运行交互式模式
func (i *Interpreter) RunInteractive() {
	fmt.Println("🎬 Render2Go Script Interpreter")
	fmt.Println("Type your commands or 'exit' to quit")
	fmt.Println("Commands: scene, create, set, animate, render, save, wait, loop")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 1

	for {
		fmt.Printf("[%d]> ", lineNumber)

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		if line == "help" {
			i.printHelp()
			continue
		}

		if line == "debug on" {
			i.debug = true
			fmt.Println("🔍 Debug mode enabled")
			continue
		}

		if line == "debug off" {
			i.debug = false
			fmt.Println("🔍 Debug mode disabled")
			continue
		}

		if line == "clear" {
			i.evaluator = NewEvaluator()
			fmt.Println("🧹 Interpreter state cleared")
			continue
		}

		if line == "objects" {
			i.listObjects()
			continue
		}

		// 执行单行命令
		err := i.RunString(line, fmt.Sprintf("line %d", lineNumber))
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		lineNumber++
	}
}

// printHelp 打印帮助信息
func (i *Interpreter) printHelp() {
	fmt.Println(`
📚 Render2Go Script Language Help

Scene Management:
  scene <width> <height> "name"     - Create a new scene

Object Creation:
  create circle <name> <radius> [(<x>, <y>)]
  create rectangle <name> <width> <height> [(<x>, <y>)]
  create line <name> (<x1>, <y1>) (<x2>, <y2>)
  create arrow <name> (<x1>, <y1>) (<x2>, <y2>)
  create text <name> "text" <size> [(<x>, <y>)]
  create polygon <name> [(<x1>, <y1>), (<x2>, <y2>), ...]

Property Setting:
  set <object>.color = #RRGGBB | deepblue | midblue | purpleblue | cyanblue | darkcolor | lightpurple
  set <object>.position = (<x>, <y>)
  set <object>.opacity = <value>
  set <object>.size = <value>

Rendering:
  render                           - Render current frame
  save "filename"                  - Save current frame

Control Flow:
  wait <seconds>                   - Wait for specified time
  loop <count> { ... }            - Repeat commands

Interactive Commands:
  help                            - Show this help
  debug on/off                    - Toggle debug mode
  clear                           - Clear interpreter state
  objects                         - List created objects
  exit/quit                       - Exit interpreter

Color Names:
  deepblue, midblue, purpleblue, cyanblue, darkcolor, lightpurple

Examples:
  scene 800 600 "my_animation"
  create circle c1 50 (400, 300)
  set c1.color = #576DA2
  set c1.opacity = 0.8
  render
  save "frame.png"`)
}

// listObjects 列出已创建的对象
func (i *Interpreter) listObjects() {
	objects := i.evaluator.GetObjects()
	if len(objects) == 0 {
		fmt.Println("📦 No objects created yet")
		return
	}

	fmt.Println("📦 Created Objects:")
	for name, obj := range objects {
		objType := "unknown"
		switch obj.(type) {
		case *geometry.Circle:
			objType = "circle"
		case *geometry.Rectangle:
			objType = "rectangle"
		case *geometry.Line:
			objType = "line"
		case *geometry.Arrow:
			objType = "arrow"
		case *geometry.Polygon:
			objType = "polygon"
		case *geometry.Text:
			objType = "text"
		}
		fmt.Printf("  %s: %s\n", name, objType)
	}
}

// GetEvaluator 获取评估器（用于测试）
func (i *Interpreter) GetEvaluator() *Evaluator {
	return i.evaluator
}

// fixPNGExtensions 自动修复输出目录中的PNG文件扩展名
func (i *Interpreter) fixPNGExtensions() error {
	outputPath := "output"

	// 检查输出目录是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return nil // 输出目录不存在，无需处理
	}

	// 遍历输出目录中的所有文件
	return filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查没有扩展名且大于4字节的文件
		if filepath.Ext(path) == "" && info.Size() > 4 {
			// 读取文件头部检查是否为PNG
			func() {
				file, err := os.Open(path)
				if err != nil {
					return // 跳过无法读取的文件
				}
				defer file.Close()

				header := make([]byte, 4)
				_, err = file.Read(header)
				if err != nil {
					return
				}

				// PNG文件头部：89 50 4E 47
				if header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 {
					// 确保文件关闭后再重命名
					file.Close()

					// 重命名文件添加.png扩展名
					newPath := path + ".png"
					if i.debug {
						fmt.Printf("🔧 Attempting to rename: %s -> %s\n", path, newPath)
					}
					err = os.Rename(path, newPath)
					if err != nil {
						if i.debug {
							fmt.Printf("❌ Rename failed: %v\n", err)
						}
						// 如果重命名失败，尝试复制+删除
						err = i.copyAndDelete(path, newPath)
						if err == nil && i.debug {
							fmt.Printf("🔧 Fixed PNG extension via copy+delete: %s -> %s\n", filepath.Base(path), filepath.Base(newPath))
						}
					} else if i.debug {
						fmt.Printf("🔧 Fixed PNG extension: %s -> %s\n", filepath.Base(path), filepath.Base(newPath))
					}
				}
			}()
		}

		return nil
	})
}

// copyAndDelete 复制文件到新位置并删除原文件
func (i *Interpreter) copyAndDelete(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 删除原文件
	return os.Remove(src)
}
