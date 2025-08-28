# Render2Go

A high-performance animation rendering framework designed for mathematical education and demonstration. Create complex mathematical animations and geometric graphics through a simple scripting language.

## 🚀 Quick Start

### Installation
```bash
# Clone the repository
git clone [<repository-url>](https://github.com/null12138/Render2Go)
cd Render2Go

# Build the program
go build -o render2go.exe cmd/render2go/main.go
```

### First Script
Create a file named `hello.r2g`:
```r2g
scene 800 600 "hello_world"
create circle my_circle 50 (400, 300)
set my_circle.color = "#FF6600"
save "hello"
```

Run the script:
```bash
./render2go hello.r2g
```

## ✨ Features

- 🎯 **Simple Script Language**: Natural language-like syntax, easy to learn
- 🧮 **Math-Friendly**: Designed specifically for mathematical concept visualization
- 🎨 **High-Quality Rendering**: High-performance 2D rendering engine based on Go
- 📁 **Auto File Management**: Smart PNG extension handling and file organization
- 🐛 **Debug Support**: Detailed error messages and execution tracing

## 📚 Documentation

- 📖 [Complete Syntax Manual](docs/SYNTAX_MANUAL.md) - Full language reference
- 🏗️ [Architecture Guide](docs/ARCHITECTURE.md) - Technical architecture details
- 🎬 [Animation Guide](docs/ANIMATION_GUIDE.md) - Animation system documentation
- 📋 [Project Overview](docs/README.md) - Comprehensive project guide

## 🎓 Learning Path

### 1. Tutorials (start here)
- [`scripts/tutorials/basic_shapes.r2g`](scripts/tutorials/basic_shapes.r2g) - Basic geometric shapes
- [`scripts/tutorials/circle_demo.r2g`](scripts/tutorials/circle_demo.r2g) - Simple animations

### 2. Examples (advanced features)
- [`scripts/examples/pythagoras.r2g`](scripts/examples/pythagoras.r2g) - Pythagorean theorem demonstration
- [`scripts/examples/math_animation.r2g`](scripts/examples/math_animation.r2g) - Mathematical animations
- [`scripts/examples/circle_circumference.r2g`](scripts/examples/circle_circumference.r2g) - Circle animations

## 🔧 Usage

### Basic Commands
```bash
# Execute a script file
./render2go script_file.r2g

# Enable debug mode
./render2go -debug script_file.r2g

# Interactive mode
./render2go -interactive

# Clean output files
./render2go -clean

# Show help
./render2go -help
```

### Example Output
```
🎬 Executing script: hello.r2g
✅ Script execution completed successfully!
```

Generated files: `output/hello_world/frames/hello.png`

## 📝 Language Syntax

### Scene Setup
```r2g
scene width height "project_name"
```

### Object Creation
```r2g
create circle my_circle 50 (400, 300)
create triangle my_triangle 100 (500, 400)
create rectangle my_rect 200 100 (300, 200)
create line my_line (0, 0) (100, 100)
create text my_text "Hello" (400, 300)
```

### Property Settings
```r2g
set my_circle.color = "#FF0000"
set my_circle.opacity = 0.8
set my_circle.position = (200, 300)
```

### Animations
```r2g
animate my_circle position (100, 100) (500, 400) 2.0
animate my_circle color "#FF0000" "#0000FF" 1.5
animate my_circle opacity 1.0 0.0 2.0
```

### File Operations
```r2g
save "frame_name"
wait 1.0
```

## 🏗️ Project Structure

```
Render2Go/
├── docs/                     # 📚 Documentation
├── cmd/render2go/           # 🚀 Command-line tool
├── core/                    # 🔧 Core modules
├── interpreter/             # 🧠 Script interpreter
├── geometry/                # 📐 Geometric shapes
├── animation/               # 🎬 Animation system
├── renderer/                # 🎨 Rendering engine
├── scene/                   # 🎭 Scene management
├── math/                    # 🧮 Math utilities
├── colors/                  # 🌈 Color system
├── scripts/                 # 📝 All example scripts (merged)
│   ├── basic_shapes.r2g    # Basic tutorials
│   ├── math_animation.r2g  # Advanced examples
│   ├── chinese_test.r2g    # Chinese text support
│   └── ...                 # More scripts
└── output/                  # 📁 Generated files (use -clean to clear)
```

## 🎯 Use Cases

### Educational Scenarios
- Mathematical theorem demonstrations (Pythagorean theorem, geometric transformations)
- Physics concept visualization (motion, waves)
- Algorithm animation demonstrations

### Creative Projects
- Artistic graphics generation
- Data visualization
- Interactive demonstrations

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and documentation
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Acknowledgments

- Built with Go and the [gg](https://github.com/fogleman/gg) 2D graphics library
- Inspired by mathematical animation frameworks
- Designed for educational and creative use

---

*Render2Go - Making mathematical animations simple*

create text greeting "Hello Render2Go!" 24 (200, 150)
set greeting.color = lightpurple

render
save "hello.png"
```

### 2
```gma
scene 600 400 "complex_demo"

create text title "Advanced Demo" 28 (300, 40)
set title.color = lightpurple

create circle c1 25 (150, 200)
set c1.color = midblue
set c1.opacity = 0.9

create rectangle r1 50 35 (300, 200)
set r1.color = cyanblue
set r1.opacity = 0.8

create arrow a1 (200, 200) (250, 200)
set a1.color = lightpurple

render
save "complex.png"
```

---

**Render2Go**
