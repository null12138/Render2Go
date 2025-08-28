package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
