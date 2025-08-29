# Render2Go 速查手册

一个轻量级的数学可视化引擎，使用Go语言开发，专为数学教育和可视化设计。

## 1. 功能概览

### 核心功能
- **基础图形绘制**：圆形、三角形、矩形、直线等基本图形
- **自动坐标系统**：智能调整坐标系比例和位置
- **图像导出**：支持PNG图像导出
- **视频导出**：基本视频导出功能（需要安装FFmpeg）

### 技术特点
- **高性能渲染**：基于现代图形API的高效渲染管道
- **简洁的脚本语言**：易学易用的r2g脚本语言
- **智能颜色管理**：预定义颜色方案，支持自定义颜色

## 2. 快速入门

### 安装和构建
```bash
# 克隆项目
git clone https://github.com/null12138/Render2Go.git
cd Render2Go

# 构建项目
go build -o render2go.exe cmd/render2go/main.go
```

### 运行示例
```bash
# 基础图形示例
./render2go.exe examples/basic_shapes.r2g

# 简单动画示例
./render2go.exe examples/simple_animation.r2g
```

## 3. r2g脚本语法

### 基本结构
```r2g
# 创建场景
scene <宽度> <高度> "项目名称"

# 创建对象
create <类型> <名称> <参数>...

# 设置属性
set <对象名>.<属性> = <值>

# 渲染和保存
render
save "文件名"

# 导出视频
export "视频文件名.mp4" <帧率> <时长>
```

### 支持的图形类型
| 类型   | 命令格式                                 | 说明                 |
| ------ | ---------------------------------------- | -------------------- |
| 圆形   | `create circle name radius`              | 创建指定半径的圆形   |
| 三角形 | `create triangle name x1 y1 x2 y2 x3 y3` | 创建三点确定的三角形 |
| 矩形   | `create rectangle name width height`     | 创建指定宽高的矩形   |
| 文本   | `create text name "内容" 字号`           | 创建文本对象         |
| 坐标系 | `create coordinate_system name "auto"`   | 创建自动坐标系       |

### 常用属性设置
| 属性   | 格式                        | 示例                       |
| ------ | --------------------------- | -------------------------- |
| 颜色   | `set obj.color = "值"`      | `set c1.color = "#3498DB"` |
| 位置   | `set obj.position = (x, y)` | `set c1.position = (2, 3)` |
| 透明度 | `set obj.opacity = 值`      | `set c1.opacity = 0.5`     |

### 颜色支持
- 十六进制格式：`"#RRGGBB"`
- 预定义颜色名：`"red"`, `"blue"`, `"green"`等

## 4. 脚本示例

### 基础图形
```r2g
# 基础图形示例
scene 800 600 "basic_shapes"

# 创建圆形
create circle c1 1.5
set c1.color = "#3498DB"
set c1.position = (-2, 0)

# 创建三角形
create triangle t1 0 -1 2 -1 1 1
set t1.color = "#E74C3C"

# 创建矩形
create rectangle r1 2 1.5
set r1.color = "#27AE60"
set r1.position = (2, 0)

# 创建标题
create text title "基础图形示例" 24
set title.position = (0, 3)
set title.color = "#2C3E50"

# 渲染和保存
render
save "basic_shapes"
```

### 简单动画
```r2g
# 简单动画示例
scene 800 600 "simple_animation"

# 创建圆形
create circle c1 1.5
set c1.color = "#3498DB"
set c1.position = (-4, 0)

# 创建标题
create text title "简单动画示例" 24
set title.position = (0, 3)
set title.color = "#2C3E50"

# 帧1
render
save "frame_0001"

# 帧2-5（移动圆形）
set c1.position = (-2, 0)
render
save "frame_0002"

set c1.position = (0, 0)
render
save "frame_0003"

set c1.position = (2, 0)
render
save "frame_0004"

set c1.position = (4, 0)
render
save "frame_0005"

# 导出视频
export "simple_animation.mp4" 30 2
```

## 5. 命令行参数

```bash
# 基本用法
render2go.exe [文件名]

# 交互模式
render2go.exe -i

# 帮助信息
render2go.exe -help

# 版本信息
render2go.exe -version

# 清理输出目录
render2go.exe -clean
```

## 6. 项目结构

```
render2go/
├── cmd/            # 命令行应用
│   └── render2go/  # 主程序入口
├── core/           # 核心组件
├── geometry/       # 几何图形
├── math/           # 数学工具
├── renderer/       # 渲染器
├── scene/          # 场景管理
├── interpreter/    # 脚本解释器
├── animation/      # 动画系统
├── colors/         # 颜色系统
├── interfaces/     # 接口定义
├── internal/       # 内部工具
└── examples/       # 示例脚本
```

## 7. 高级技巧

### 坐标系统
- 使用`create coordinate_system`创建坐标系
- 坐标系中心点为(0,0)
- 可以使用`"auto"`自动调整坐标系缩放

### 输出管理
- 输出文件保存在`output/<项目名>/frames/`目录中
- 帧序列保存在`output/<项目名>/frames/`目录中
- 视频文件(如果成功)保存在项目根目录

### 视频导出
- 使用`export`命令导出视频
- 需要安装FFmpeg才能成功合成视频
- 如果FFmpeg未安装，仍会生成帧序列

## 8. 常见问题

### Q: 为什么视频无法导出？
A: 需要安装FFmpeg并确保它在系统PATH中。

### Q: 如何调整对象的大小？
A: 对于大多数对象，在创建时指定其大小参数。例如，圆形的半径，矩形的宽高。

### Q: 如何使用自定义字体？
A: 当前版本使用系统默认字体。

### Q: 如何调整坐标系缩放？
A: 使用`create coordinate_system coords "auto"`创建自动调整的坐标系。

---

这是Render2Go的精简版本，专注于数学可视化的核心功能，包括基础图形绘制和视频导出。所有的TEX和Markdown高级功能已被移除以简化项目。

### 基本示例
```r2g
# 创建场景
scene 800 600 "math_demo"

# 数学公式演示
create mathtex formula "x^2 + y^2 = r^2" "large" (0, 2)
set formula.color = "mathblue"

# Markdown文档
create markdown doc "# 勾股定理

对于直角三角形：$a^2 + b^2 = c^2$

## 特殊情况
当 $a = 3, b = 4$ 时：
$$c = sqrt(9 + 16) = 5$$" "normal" (0, 0)

# 添加圆形图示
create circle my_circle 1.5 (-3, 0)
set my_circle.color = "primary"

# 动画效果
animate fadein formula 1.0
animate fadein doc 1.5
animate fadein my_circle 2.0

# 渲染
render
save "math_demo"
```

## ✨ 特性

- 🎨 **多种图形**: 圆形、三角形、矩形、线条、箭头、多边形
- 📝 **数学公式**: LaTeX和Markdown兼容的数学公式渲染
- 📄 **Markdown支持**: 完整的Markdown语法，包括内嵌数学公式  
- 🎬 **动画系统**: 平滑的渐变、移动、缩放、旋转动画
- 📐 **智能坐标系**: 自动缩放，多种坐标系模式
- 🎯 **简洁语法**: 易学易用的 .r2g 脚本语言
- 🌍 **Web编辑器**: 在线编辑和预览

## 📚 文档

- 📖 [完整文档](DOCUMENTATION.md) - 详细语法手册和技术架构
- 🧮 [数学公式和Markdown指南](docs/MATH_MARKDOWN_GUIDE.md) - 数学公式和Markdown渲染完整指南
- 🔧 [快速参考](QUICK_REFERENCE.md) - 常用语法速查
- 🏗️ [项目结构](PROJECT_STRUCTURE.md) - 代码组织说明

## 🛠️ 技术栈

- **渲染引擎**: [fogleman/gg](https://github.com/fogleman/gg)
- **语言**: Go 1.21+
- **图形格式**: PNG
- **字体支持**: 系统字体 + 中文字体

## 📁 项目结构

```
Render2Go/
├── cmd/render2go/      # 主程序
├── core/              # 核心接口
├── geometry/          # 几何图形
├── math/             # 数学计算
├── renderer/         # 渲染引擎
├── interpreter/      # 脚本解释器
└── web-editor/       # Web编辑器
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[MIT License](LICENSE)

---
⭐ 如果这个项目对你有帮助，请给一个 Star！
