# Render2Go 项目结构

```
Render2Go/
├── cmd/
│   └── render2go/
│       └── main.go              # 主程序入口
├── internal/                    # 内部包（仅内部使用）
│   ├── defaults/
│   │   └── defaults.go          # 默认配置系统
│   ├── markdown/
│   │   └── renderer.go          # Markdown文档渲染器
│   └── tex/
│       └── renderer.go          # TeX数学公式渲染器
├── core/
│   └── mobject.go               # 核心对象系统
├── geometry/
│   └── shapes.go                # 几何图形实现
├── math/
│   └── vector.go                # 数学向量系统
├── colors/
│   └── colors.go                # 颜色系统
├── animation/
│   └── animation.go             # 动画系统
├── scene/
│   └── scene.go                 # 场景管理
├── renderer/
│   └── renderer.go              # 渲染引擎
├── interpreter/
│   ├── lexer.go                 # 词法分析器
│   ├── parser.go                # 语法分析器
│   ├── interpreter.go           # 解释器
│   └── evaluator.go             # 求值器
├── examples/
│   ├── math/                    # 数学相关示例
│   │   ├── mathtex_demo.r2g     # MathTeX演示
│   │   ├── markdown_test.r2g    # Markdown测试
│   │   ├── math_animation.r2g   # 数学动画
│   │   ├── comprehensive_demo.r2g # 综合演示
│   │   └── README.md            # 示例说明
│   └── ...                     # 其他示例
├── docs/
│   └── MATH_MARKDOWN_GUIDE.md   # 数学公式和Markdown指南
├── web-editor/                  # Web编辑器
│   ├── index.html               # 主界面
│   ├── css/                     # 样式文件
│   ├── js/                      # JavaScript文件
│   └── assets/                  # 静态资源
├── scripts/                     # 构建和部署脚本
├── output/                      # 渲染输出目录
├── go.mod                       # Go模块定义
├── go.sum                       # Go依赖哈希
├── README.md                    # 项目说明
└── LICENSE                      # 开源许可
```

## 架构说明

### 核心组件

#### 1. 脚本解释器 (`interpreter/`)
- **lexer.go**: 词法分析，将源代码转换为token流
- **parser.go**: 语法分析，构建抽象语法树(AST)
- **interpreter.go**: 解释器主逻辑
- **evaluator.go**: 求值器，执行脚本命令

#### 2. 渲染系统 (`internal/`, `renderer/`)
- **renderer.go**: 核心渲染引擎
- **markdown/renderer.go**: Markdown文档渲染器
- **tex/renderer.go**: TeX数学公式渲染器

#### 3. 对象系统 (`core/`, `geometry/`)
- **mobject.go**: 核心对象接口和基础实现
- **shapes.go**: 几何图形对象实现

#### 4. 动画系统 (`animation/`)
- **animation.go**: 动画效果实现，包括缓动函数

#### 5. 配置系统 (`internal/defaults/`)
- **defaults.go**: 默认配置管理，包括颜色、字体等

### 新增功能模块

#### 数学公式渲染
- **MathTeX支持**: 通用的markdown-tex兼容语法
- **TeX支持**: 完整的LaTeX数学公式语法
- **符号映射**: 支持简化语法和传统LaTeX语法

#### Markdown文档渲染
- **完整语法支持**: 标题、段落、列表、代码块
- **内嵌数学**: 支持`$...$`和`$$...$$`数学公式
- **动画集成**: Markdown内容完全支持动画效果

#### 动画增强
- **数学对象动画**: 公式和文档的淡入淡出、移动、缩放
- **时间控制**: 精确的动画时序控制
- **组合动画**: 多个对象的协调动画

## 开发指南

### 添加新的对象类型
1. 在`geometry/shapes.go`中定义新的结构体
2. 实现`MObject`接口方法
3. 在`interpreter/lexer.go`中添加新的token类型
4. 在`interpreter/parser.go`中更新语法解析
5. 在`interpreter/evaluator.go`中添加创建方法

### 添加新的动画效果
1. 在`animation/animation.go`中定义新的动画类型
2. 实现动画的更新逻辑
3. 在解释器中添加相应的命令支持

### 扩展渲染功能
1. 在`renderer/`目录下添加新的渲染器
2. 实现必要的渲染接口
3. 在主渲染引擎中集成新功能

## 构建和部署

### 开发环境
```bash
# 克隆项目
git clone <repository>
cd Render2Go

# 安装依赖
go mod tidy

# 开发构建
go build -o render2go.exe ./cmd/render2go

# 运行测试
go test ./...
```

### 生产构建
```bash
# 优化构建
go build -ldflags="-s -w" -o render2go.exe ./cmd/render2go

# 交叉编译（示例）
GOOS=linux GOARCH=amd64 go build -o render2go-linux ./cmd/render2go
GOOS=darwin GOARCH=amd64 go build -o render2go-mac ./cmd/render2go
```

## 技术栈

- **语言**: Go 1.21+
- **图形**: 基于OpenGL的2D渲染
- **解析**: 手写词法分析器和递归下降解析器
- **数学**: 自定义向量和矩阵运算
- **Web**: 原生HTML/CSS/JavaScript（无框架依赖）
