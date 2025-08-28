# Render2Go 项目概览

## 简介

Render2Go 是一个专为数学教育和动画演示设计的高性能渲染框架。通过简洁的脚本语言，用户可以轻松创建复杂的数学动画、几何图形和教学演示。

## 特性

### 🎯 核心特性
- **简洁的脚本语言**: 类似自然语言的语法，易学易用
- **数学友好**: 专为数学概念可视化设计
- **高质量渲染**: 基于Go语言的高性能2D渲染引擎
- **自动文件管理**: 智能PNG扩展名处理和文件组织
- **调试支持**: 详细的错误信息和执行跟踪

### 🚀 技术特性
- **模块化架构**: 清晰的分层设计，易于扩展
- **零配置**: 开箱即用，无需复杂设置
- **跨平台**: 支持Windows、Linux、macOS
- **命令行友好**: 支持批处理和自动化

## 快速开始

### 安装
```bash
# 克隆项目
git clone <repository-url>
cd Render2Go

# 构建程序
go build -o render2go.exe cmd/render2go/main.go
```

### 第一个脚本
创建文件 `hello.r2g`:
```r2g
scene 800 600 "hello_world"
create circle my_circle 50 (400, 300)
set my_circle.color = "#FF6600"
save "hello"
```

运行脚本:
```bash
./render2go hello.r2g
```

## 项目结构

```
Render2Go/
├── docs/                      # 📚 文档目录
│   ├── SYNTAX_MANUAL.md      # 语法手册
│   ├── ARCHITECTURE.md       # 架构文档
│   └── ANIMATION_GUIDE.md    # 动画指南
├── cmd/                      # 🚀 命令行工具
│   └── render2go/
│       └── main.go          # 程序入口
├── core/                     # 🔧 核心模块
│   └── mobject.go           # 可动画对象基类
├── interpreter/              # 🧠 解释器
│   ├── lexer.go             # 词法分析器
│   ├── parser.go            # 语法分析器
│   ├── evaluator.go         # 执行引擎
│   └── interpreter.go       # 主解释器
├── geometry/                 # 📐 几何图形
│   └── shapes.go            # 基础图形定义
├── animation/                # 🎬 动画系统
│   └── animation.go         # 动画效果
├── renderer/                 # 🎨 渲染引擎
│   └── renderer.go          # 2D图形渲染
├── scene/                    # 🎭 场景管理
│   └── scene.go             # 场景和对象管理
├── math/                     # 🧮 数学库
│   └── vector.go            # 向量和坐标系统
├── colors/                   # 🌈 色彩系统
│   └── colors.go            # 颜色管理和配色方案
├── scripts/                  # 📝 示例脚本
│   ├── tutorials/           # 教程脚本
│   │   ├── basic_shapes.r2g
│   │   └── circle_demo.r2g
│   └── examples/            # 高级示例
│       ├── pythagoras.r2g
│       ├── math_animation.r2g
│       └── circle_circumference.r2g
└── output/                   # 📁 输出目录
    └── [project_name]/
        └── frames/
            └── *.png
```

## 核心概念

### 1. 场景 (Scene)
场景是所有图形对象的容器，定义了画布的尺寸和坐标系统。

### 2. 对象 (Mobject)
所有可见元素都是对象，包括图形、文本等。对象具有位置、颜色、大小等属性。

### 3. 动画 (Animation)
动画是对象属性随时间的变化，如移动、缩放、颜色渐变等。

### 4. 渲染 (Render)
将场景和对象转换为图像文件的过程。

## 语法概览

### 基本语法
```r2g
// 创建场景
scene width height "project_name"

// 创建对象
create shape_type object_name parameters

// 设置属性
set object_name.property = value

// 创建动画
animate object_name property from_value to_value duration

// 保存图像
save "filename"
```

### 支持的图形
- **圆形**: `create circle name radius (x, y)`
- **三角形**: `create triangle name size (x, y)`
- **矩形**: `create rectangle name width height (x, y)`
- **线段**: `create line name (x1, y1) (x2, y2)`
- **文本**: `create text name "content" (x, y)`

### 动画类型
- **位置动画**: `animate obj position (x1,y1) (x2,y2) duration`
- **颜色动画**: `animate obj color "color1" "color2" duration`
- **透明度动画**: `animate obj opacity value1 value2 duration`

## 开发指南

### 添加新图形
1. 在 `geometry/shapes.go` 中定义新图形结构
2. 实现 `Mobject` 接口
3. 在解释器中添加创建逻辑

### 添加新动画
1. 在 `animation/animation.go` 中定义动画类型
2. 实现动画插值逻辑
3. 在执行引擎中注册动画

### 扩展语法
1. 在 `lexer.go` 中添加新Token
2. 在 `parser.go` 中添加语法规则
3. 在 `evaluator.go` 中实现执行逻辑

## 使用案例

### 教育场景
- 数学定理演示（勾股定理、几何变换）
- 物理概念可视化（运动、波动）
- 算法动画演示

### 创意项目
- 艺术图形生成
- 数据可视化
- 交互式演示

## 性能优化

### 最佳实践
- 合理设置场景尺寸
- 避免过度复杂的图形
- 使用适当的动画时长
- 定期清理输出文件

### 性能监控
- 使用 `-debug` 模式查看执行详情
- 监控内存使用情况
- 测试不同场景尺寸的性能影响

## 社区和支持

### 获取帮助
- 查看 `docs/SYNTAX_MANUAL.md` 获取完整语法参考
- 查看 `scripts/tutorials/` 学习基础用法
- 查看 `scripts/examples/` 了解高级特性

### 贡献指南
- 遵循现有的代码风格
- 添加适当的测试用例
- 更新相关文档
- 提交清晰的Pull Request

## 版本历史

### v1.0 (当前)
- ✅ 基础脚本语言支持
- ✅ 核心几何图形
- ✅ 基础动画系统
- ✅ PNG文件输出
- ✅ 自动扩展名修复
- ✅ 调试模式支持

### 计划功能
- 🔄 视频输出支持
- 🔄 交互式编辑器
- 🔄 更多几何图形
- 🔄 高级动画效果
- 🔄 3D渲染支持

## 许可证

本项目采用 MIT 许可证，详见 `LICENSE` 文件。

---

*Render2Go - 让数学动画变得简单*
