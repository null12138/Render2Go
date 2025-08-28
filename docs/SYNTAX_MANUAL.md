# Render2Go 语法手册

Render2Go 是一个基于Go语言开发的动画渲染框架，提供简洁易用的脚本语言来创建数学动画和几何图形。

## 🚀 快速参考

### 基本脚本结构
```r2g
// 1. 设置场景
scene 1200 800 "my_project"

// 2. 创建对象
create circle my_circle 50 (400, 300)
create text title "Hello World" 32 (600, 500)

// 3. 设置属性
set my_circle.color = "#3366CC"
set title.color = "#2C3E50"

// 4. 保存和渲染
save "my_frame"
render
```

### 关键语法要点
- **字符串必须用引号**: `"project_name"`, `"Hello World"`
- **坐标使用圆括号**: `(x, y)`, `(400, 300)`
- **属性设置格式**: `set object.property = value`
- **注释支持**: `//` 或 `#` 开头的单行注释
- **颜色格式**: `"#FF0000"` 或 `"red"`

---

## 目录

1. [基础语法](#基础语法)
2. [场景设置](#场景设置)
3. [对象创建](#对象创建)
4. [属性设置](#属性设置)
5. [动画控制](#动画控制)
6. [文件操作](#文件操作)
7. [循环和控制](#循环和控制)
8. [语法检查和错误排除](#语法检查和错误排除)
9. [示例脚本](#示例脚本)
10. [命令行工具](#命令行工具)

---

## 基础语法

### 注释
```r2g
// 这是单行注释（支持 # 开头的注释）
# 这也是单行注释
```

**注意**: 
- 使用 `//` 或 `#` 开头的单行注释
- 颜色值使用 `#` 时需要紧跟数字或字母，如 `#FF0000`

### 基本数据类型
- **数字**: `123`, `45.67`, `-89.01`
- **字符串**: `"hello"`, `"world"`
- **坐标**: `(x, y)` 如 `(100, 200)`
- **颜色**: 十六进制 `"#FF0000"` 或颜色名 `"red"`

---

## 场景设置

### scene - 创建场景
```r2g
scene width height "project_name"
```
- `width`: 场景宽度（像素）
- `height`: 场景高度（像素）  
- `project_name`: 项目名称（字符串）

**示例:**
```r2g
scene 1920 1080 "my_animation"
scene 800 600 "simple_demo"
```

---

## 对象创建

### create - 创建几何对象

#### 圆形 (circle)
```r2g
create circle object_name radius (center_x, center_y)
```
**参数说明:**
- `object_name`: 对象名称（标识符）
- `radius`: 半径（数字）
- `(center_x, center_y)`: 圆心坐标

**示例:**
```r2g
create circle my_circle 50 (400, 300)
create circle small_dot 10 (100, 100)
create circle reference_circle 150 (400, 400)
```

#### 三角形 (triangle) 
```r2g
create triangle object_name size (center_x, center_y)
```
**参数说明:**
- `object_name`: 对象名称（标识符）
- `size`: 三角形大小（数字）
- `(center_x, center_y)`: 中心坐标

**示例:**
```r2g
create triangle my_triangle 100 (500, 400)
create triangle red_triangle 150 (960, 540)
```

#### 矩形 (rectangle)
```r2g
create rectangle object_name width height (center_x, center_y)
```
**示例:**
```r2g
create rectangle my_rect 200 100 (300, 200)
```

#### 线段 (line)
```r2g
create line object_name (start_x, start_y) (end_x, end_y)
```
**参数说明:**
- `object_name`: 对象名称（标识符）
- `(start_x, start_y)`: 起始点坐标
- `(end_x, end_y)`: 结束点坐标

**示例:**
```r2g
create line my_line (0, 0) (100, 100)
create line triangle_side1 (400, 250) (530, 475)
create line hex_side1 (550, 325) (550, 475)
```

#### 文本 (text)
```r2g
create text object_name "content" size (x, y)
```
**参数说明:**
- `object_name`: 对象名称（标识符）
- `"content"`: 文本内容（字符串，必须用引号）
- `size`: 字体大小（数字）
- `(x, y)`: 文本位置坐标

**默认属性:**
- 默认颜色：黑色 (避免在白色背景上不可见)
- 默认透明度：1.0 (完全不透明)
- 文本居中对齐

**示例:**
```r2g
create text title "Hello World" 32 (400, 300)
create text subtitle "正多边形逼近圆的方法" 24 (600, 700)
create text step1 "第1步: 正三角形 (3边)" 20 (600, 550)
```

---

## 属性设置

### set - 设置对象属性

**基本语法:**
```r2g
set object_name.property = value
```

#### 颜色设置
```r2g
set object_name.color = "color_value"
```
**支持的颜色格式:**
- 十六进制: `"#FF0000"` (红色), `"#3366CC"` (蓝色)
- 颜色名: `"red"`, `"blue"`, `"green"`, `"yellow"`, `"purple"`, `"orange"`

**示例:**
```r2g
set my_circle.color = "#3366CC"
set triangle_side1.color = "#E74C3C"
set title.color = "#2C3E50"
```

#### 透明度设置
```r2g
set object_name.opacity = value
```
- `value`: 0.0 (完全透明) 到 1.0 (完全不透明)

**示例:**
```r2g
set reference_circle.opacity = 0.3
set my_triangle.opacity = 1.0
```

#### 线条宽度
```r2g
set object_name.stroke_width = value
```
**示例:**
```r2g
set my_line.stroke_width = 3.0
```

#### 位置设置
```r2g
set object_name.position = (x, y)
```
**示例:**
```r2g
set my_circle.position = (200, 300)
```

---

## 动画控制

### animate - 创建动画
```r2g
animate object_name property from_value to_value duration
```

#### 位置动画
```r2g
animate object_name position (start_x, start_y) (end_x, end_y) duration
```
**示例:**
```r2g
animate my_circle position (100, 100) (500, 400) 2.0
```

#### 颜色动画
```r2g
animate object_name color "start_color" "end_color" duration
```
**示例:**
```r2g
animate my_circle color "#FF0000" "#0000FF" 1.5
```

#### 透明度动画
```r2g
animate object_name opacity start_value end_value duration
```
**示例:**
```r2g
animate my_circle opacity 1.0 0.0 2.0
```

### wait - 等待
```r2g
wait duration
```
**示例:**
```r2g
wait 1.0    // 等待1秒
wait 0.5    // 等待0.5秒
```

---

## 文件操作

### save - 保存当前帧
```r2g
save "filename"
```
**参数说明:**
- `"filename"`: 文件名（字符串，必须用引号）
- 文件会自动添加 `.png` 扩展名
- 保存到 `output/项目名/frames/` 目录

**示例:**
```r2g
save "pi_derivation_start"
save "pi_derivation_triangle"
save "pi_derivation_complete"
```

### render - 渲染当前帧
```r2g
render
```
**说明:**
- 渲染当前场景中的所有对象
- 通常在 `save` 命令后调用
- 自动保存到项目目录

**示例:**
```r2g
save "my_frame"
render
```

### render - 渲染动画
```r2g
render fps duration "output_name"
```
- `fps`: 帧率 (如 30, 60)
- `duration`: 持续时间（秒）
- `output_name`: 输出文件名

**示例:**
```r2g
render 30 5.0 "my_animation"
```

---

## 循环和控制

### loop - 循环执行
```r2g
loop count {
    // 循环体
}
```
**示例:**
```r2g
loop 10 {
    animate my_circle position (100, 100) (500, 100) 0.5
    wait 0.1
}
```

---

## 示例脚本

### 基础图形演示
```r2g
// 基础图形演示
scene 1920 1080 "basic_shapes"

// 创建圆形
create circle blue_circle 80 (400, 540)
set blue_circle.color = "#3366CC"
set blue_circle.opacity = 1.0

// 创建三角形
create triangle red_triangle 150 (960, 540)
set red_triangle.color = "#CC3366"
set red_triangle.opacity = 0.8

// 创建矩形
create rectangle green_rect 200 100 (1520, 540)
set green_rect.color = "#33CC66"
set green_rect.opacity = 0.9

// 保存图像
save "basic_shapes"
```

### 简单动画
```r2g
// 圆形移动动画
scene 800 600 "circle_animation"

create circle moving_circle 30 (50, 300)
set moving_circle.color = "#FF6600"

// 动画：从左移动到右
animate moving_circle position (50, 300) (750, 300) 3.0

// 渐变透明
animate moving_circle opacity 1.0 0.0 1.0

save "final_frame"
```

### 数学演示 - π推导
```r2g
// 圆周率π推导演示 - 使用正多边形逼近圆的方法
scene 1200 800 "pi_derivation"

// 创建参考圆
create circle reference_circle 150 (400, 400)
set reference_circle.color = "#3366CC"
set reference_circle.opacity = 0.3

// 标题文本
create text title "圆周率π的推导演示" 32 (600, 750)
set title.color = "#2C3E50"

create text subtitle "正多边形逼近圆的方法" 24 (600, 700)
set subtitle.color = "#34495E"

// 保存初始状态
save "pi_derivation_start"
render

// 第一步：正三角形
create text step1 "第1步: 正三角形 (3边)" 20 (600, 550)
set step1.color = "#C0392B"

// 创建正三角形的三条边
create line triangle_side1 (400, 250) (530, 475)
set triangle_side1.color = "#E74C3C"

create line triangle_side2 (530, 475) (270, 475)
set triangle_side2.color = "#E74C3C"

create line triangle_side3 (270, 475) (400, 250)
set triangle_side3.color = "#E74C3C"

create text triangle_result "π ≈ 2.598 (误差很大)" 18 (600, 500)
set triangle_result.color = "#E74C3C"

save "pi_derivation_triangle"
render

// 第二步：正六边形
create text step2 "第2步: 正六边形 (6边)" 20 (600, 450)
set step2.color = "#D35400"

// 创建六边形边
create line hex_side1 (550, 325) (550, 475)
set hex_side1.color = "#D35400"

create line hex_side2 (550, 475) (400, 550)
set hex_side2.color = "#D35400"

create line hex_side3 (400, 550) (250, 475)
set hex_side3.color = "#D35400"

create text hex_result "π ≈ 3.000 (更接近了)" 18 (600, 400)
set hex_result.color = "#D35400"

save "pi_derivation_complete"
render
```

### 基础图形演示
```r2g
// 基础图形演示
scene 1920 1080 "basic_shapes"

// 创建圆形
create circle blue_circle 80 (400, 540)
set blue_circle.color = "#3366CC"
set blue_circle.opacity = 1.0

// 创建三角形
create triangle red_triangle 150 (960, 540)
set red_triangle.color = "#CC3366"
set red_triangle.opacity = 0.8

// 保存图像
save "basic_shapes"
render
```

---

## 命令行工具

### 基本用法
```bash
# 执行脚本文件
./render2go script_file.r2g

# 启用调试模式
./render2go -debug script_file.r2g

# 交互式模式
./render2go -interactive

# 清理输出文件
./render2go -clean

# 显示帮助
./render2go -help
```

### 命令行选项
- `-debug`: 启用调试模式，显示详细的解析和执行信息
- `-interactive`: 启动交互式命令行模式
- `-clean`: 清理输出目录中的所有文件
- `-help`: 显示帮助信息
- `-version`: 显示版本信息

### 调试模式输出
启用调试模式时，会显示：
- 🔍 词法分析结果（Token列表）
- 🌳 语法分析结果（AST抽象语法树）
- 🚀 执行过程信息
- 🔧 PNG文件扩展名自动修复过程
- ✅ 执行完成确认

---

## 文件组织

### 项目结构
```
project_name/
├── script.r2g          # 脚本文件
└── output/
    └── project_name/
        └── frames/
            ├── frame1.png
            ├── frame2.png
            └── ...
```

### 输出文件
- 所有生成的图像文件保存在 `output/项目名/frames/` 目录下
- 文件格式为PNG，会自动添加 `.png` 扩展名
- 支持自动PNG文件扩展名修复功能

---

## 错误处理

### 常见错误类型
1. **语法错误**: 脚本语法不正确
2. **对象未找到**: 引用了不存在的对象
3. **类型错误**: 参数类型不匹配
4. **文件错误**: 无法创建或写入输出文件

### 错误信息格式
```
❌ Error: [错误类型] 
详细错误描述 (文件名:行号)
```

---

## 语法检查和错误排除

### 常见语法错误

#### 1. 字符串引号错误
```r2g
// ❌ 错误：缺少引号
scene 800 600 project_name
create text title Hello World 32 (400, 300)

// ✅ 正确：字符串必须用引号
scene 800 600 "project_name"
create text title "Hello World" 32 (400, 300)
```

#### 2. 注释与颜色冲突
```r2g
// ❌ 错误：# 后直接跟空格会被当作注释
set circle.color = # FF0000

// ✅ 正确：颜色值需要引号包围
set circle.color = "#FF0000"
```

#### 3. 属性设置语法错误
```r2g
// ❌ 错误：缺少等号和点号
set circle color "#FF0000"
set circle opacity 0.5

// ✅ 正确：必须使用点号和等号
set circle.color = "#FF0000"
set circle.opacity = 0.5
```

#### 4. 坐标格式错误
```r2g
// ❌ 错误：坐标格式不正确
create circle test 50 [400, 300]
create line test1 400,300 500,400

// ✅ 正确：坐标使用圆括号
create circle test 50 (400, 300)
create line test1 (400, 300) (500, 400)
```

### 调试技巧

#### 1. 检查解析错误
- 运行脚本时注意错误信息中的行号
- 检查该行及前后行的语法
- 确认所有字符串都有引号

#### 2. 验证对象创建
- 先创建简单对象测试
- 逐步添加复杂属性
- 使用 `render` 命令验证渲染结果

#### 3. 文件输出检查
- 确认输出目录存在
- 检查生成的PNG文件大小（空白图片通常很小）
- 查看终端输出的成功/错误信息

---

## 最佳实践

### 1. 命名规范
- 对象名使用下划线分隔: `reference_circle`, `triangle_side1`
- 项目名使用描述性名称: `"pythagoras_theorem"`, `"pi_derivation"`

### 2. 代码组织
- 先设置场景
- 再创建对象
- 然后设置属性
- 最后执行动画和保存

### 3. 性能考虑
- 合理设置场景尺寸，避免过大的画布
- 适当使用等待时间控制动画节奏
- 定期保存重要帧以便调试

### 4. 调试技巧
- 使用 `-debug` 模式查看执行详情
- 先测试简单图形，再添加复杂动画
- 分步保存帧来验证每个阶段的结果

---

## 版本信息

**当前版本**: Render2Go v1.0  
**语法版本**: R2G v1.0  
**支持平台**: Windows, Linux, macOS  
**依赖**: Go 1.19+

---

*本手册持续更新中，如有问题请查看项目文档或提交Issue。*
