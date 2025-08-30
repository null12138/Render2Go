# Render2Go 语法手册 (R2G Script Language)

⚠️ **重要提示**: 请始终参考此语法手册编写R2G脚本，确保语法正确性。

Render2Go 是一个用于创建动画和可视化的脚本语言，支持序列帧渲染和自动视频生成。

## 目录
1. [场景定义](#1-场景定义)
2. [对象创建](#2-对象创建)
3. [属性设置](#3-属性设置)
4. [动画系统](#4-动画系统)
5. [渲染命令](#5-渲染命令)
6. [坐标系统](#坐标系统)
7. [颜色支持](#颜色支持)
8. [完整示例](#完整示例)
9. [最佳实践](#最佳实践)

---

## 1. 场景定义

### 语法
```r2g
scene <width> <height> "<title>"
```

### 参数说明
- `width`: 场景宽度（像素）
- `height`: 场景高度（像素）  
- `title`: 场景标题（字符串，必须用引号）

### 坐标系统
- **原点**: 屏幕中心 (0, 0)
- **X轴**: 左负(-) 右正(+)
- **Y轴**: 下负(-) 上正(+)
- **坐标范围**: X[-width/2, +width/2], Y[-height/2, +height/2]

### 示例
```r2g
scene 800 600 "我的动画"  # 创建800x600场景，坐标范围X[-400,400], Y[-300,300]
```

---

## 2. 对象创建

### 基本语法
```r2g
create <type> <name> <parameters> (<x>, <y>)
```

### 支持的对象类型

#### 圆形 (circle)
```r2g
create circle <name> <radius> (<x>, <y>)
```
- `radius`: 圆的半径（像素）
- `(x, y)`: 圆心坐标

**示例:**
```r2g
create circle ball 20 (0, 100)    # 创建半径20的球，位于(0,100)
```

#### 矩形 (rectangle)
```r2g
create rectangle <name> <width> <height> (<x>, <y>)
```
- `width`: 矩形宽度（像素）
- `height`: 矩形高度（像素）
- `(x, y)`: 矩形中心坐标

**示例:**
```r2g
create rectangle ground 600 20 (0, -250)    # 创建地面矩形
```

#### 三角形 (triangle)
```r2g
create triangle <name> <size> (<x>, <y>)
```
- `size`: 三角形大小
- `(x, y)`: 三角形中心坐标

#### 线条 (line)
```r2g
create line <name> (<x1>, <y1>) (<x2>, <y2>)
```
- `(x1, y1)`: 起点坐标
- `(x2, y2)`: 终点坐标

#### 文本 (text)
```r2g
create text <name> "<content>" <font_size> (<x>, <y>)
```
- `content`: 文本内容（必须用引号）
- `font_size`: 字体大小（数字或预定义名称）
- `(x, y)`: 文本位置

**支持的字体大小名称:**
- `tiny` - 极小字体
- `small` - 小字体
- `normal` - 正常字体
- `large` - 大字体
- `huge` - 超大字体
- `title` - 标题字体

**示例:**
```r2g
create text title "我的标题" title (0, 200)
create text label "标签文本" normal (100, 100)
create text note "注释" small (0, -50)
```

---

## 3. 属性设置

### 基本语法
```r2g
set <object>.<property> = <value>
```

### 支持的属性

#### 颜色 (color)
```r2g
set <object>.color = <color_name>
```

#### 位置 (position)
```r2g
set <object>.position = (<x>, <y>)
```

#### 大小属性
```r2g
set <circle>.radius = <value>
set <rectangle>.width = <value>
set <rectangle>.height = <value>
```

### 示例
```r2g
set ball.color = red
set ground.color = green
set ball.position = (100, 200)
set circle.radius = 50
```

---

## 4. 动画系统

### 基本语法
```r2g
animate <action> <object> <parameters> <duration>
```

### 支持的动画类型

#### 移动动画 (move)
```r2g
animate move <object> (<target_x>, <target_y>) <duration>
```
- `(target_x, target_y)`: 目标位置
- `duration`: 动画持续时间（秒）

#### 缩放动画 (scale)
```r2g
animate scale <object> <factor> <duration>
```
- `factor`: 缩放倍数
- `duration`: 动画持续时间（秒）

### 动画时长建议
- **短动画**: 0.1 - 0.5秒（快速变化）
- **中等动画**: 0.5 - 1.0秒（正常速度）
- **长动画**: 1.0 - 2.0秒（缓慢变化）

### 示例
```r2g
animate move ball (0, -250) 1.0    # 1秒内移动到(0,-250)
animate move ball (0, 150) 0.8     # 0.8秒内弹起到(0,150)
animate scale circle 1.5 0.5       # 0.5秒内放大1.5倍
```

---

## 5. 渲染命令

### 语法
```r2g
render_frames <fps> <duration> "<output_path>"
```

### 参数说明
- `fps`: 帧率（建议使用60）
- `duration`: 总动画时长（秒）
- `output_path`: 输出路径（必须用引号）

### 输出格式
自动生成：
- PNG序列帧（frame_000001.png, frame_000002.png, ...）
- MP4视频文件（animation.mp4）
- GIF动图文件（animation.gif）

### 示例
```r2g
render_frames 60 4.0 "output/my_animation"    # 60fps，4秒，输出到output/my_animation/
```

---

## 坐标系统

### 坐标映射
- **1:1像素映射**: 逻辑坐标直接对应屏幕像素
- **固定缩放**: 使用SetFixedScale(1.0)，避免自动缩放

### 常用场景尺寸
| 场景尺寸  | X坐标范围   | Y坐标范围   | 用途       |
| --------- | ----------- | ----------- | ---------- |
| 800×600   | [-400, 400] | [-300, 300] | 标清动画   |
| 1200×800  | [-600, 600] | [-400, 400] | 高清动画   |
| 1920×1080 | [-960, 960] | [-540, 540] | 全高清动画 |

### 坐标计算示例
```r2g
# 800×600场景
scene 800 600 "示例"

# 地面位置（底部）: Y = -300 + 地面高度/2
create rectangle ground 600 20 (0, -290)    # 地面顶部在Y=-280

# 天花板位置（顶部）: Y = 300 - 天花板高度/2  
create rectangle ceiling 600 20 (0, 290)    # 天花板底部在Y=280

# 左墙位置: X = -400 + 墙宽度/2
create rectangle leftwall 20 600 (-390, 0)  # 左墙右边缘在X=-380

# 右墙位置: X = 400 - 墙宽度/2
create rectangle rightwall 20 600 (390, 0)  # 右墙左边缘在X=380
```

---

## 颜色支持

### 基础颜色
```
black, white, red, green, blue, yellow, cyan, magenta
```

### 主题颜色
```
primary, secondary, accent, background, surface
```

### 状态颜色
```
error, success, warning, info, muted
```

### 数学颜色
```
mathred, mathblue, mathgreen, mathorange, mathpurple
```

### 专业配色
```
deepblue, midblue, purpleblue, cyanblue
```

### 颜色使用示例
```r2g
set ball.color = red
set ground.color = green
set wall.color = mathblue
set text.color = warning
```

---

## 完整示例

### 物理弹跳球动画
```r2g
# 场景设置
scene 800 600 "物理弹跳球"

# 创建对象
create circle ball 15 (0, 200)
set ball.color = red

create rectangle ground 600 20 (0, -280)
set ground.color = green

create rectangle leftwall 20 350 (-390, -125)
set leftwall.color = mathblue

create rectangle rightwall 20 350 (390, -125) 
set rightwall.color = mathblue

# 物理弹跳动画序列
# 第一次下落
animate move ball (0, 150) 0.15
animate move ball (0, 80) 0.15
animate move ball (0, -20) 0.15
animate move ball (0, -120) 0.15
animate move ball (0, -200) 0.15
animate move ball (0, -260) 0.1    # 撞击地面

# 第一次弹跳（能量70%）
animate move ball (0, -200) 0.1
animate move ball (0, -120) 0.15
animate move ball (0, -20) 0.15
animate move ball (0, 80) 0.15
animate move ball (0, 140) 0.15    # 峰值
animate move ball (0, 80) 0.15
animate move ball (0, -20) 0.15
animate move ball (0, -120) 0.15
animate move ball (0, -200) 0.15
animate move ball (0, -260) 0.1    # 第二次撞击

# 第二次弹跳（能量50%）
animate move ball (0, -200) 0.1
animate move ball (0, -120) 0.12
animate move ball (0, -20) 0.12
animate move ball (0, 40) 0.12     # 较低峰值
animate move ball (0, -20) 0.12
animate move ball (0, -120) 0.12
animate move ball (0, -200) 0.12
animate move ball (0, -260) 0.1    # 第三次撞击

# 最后小弹跳并停止
animate move ball (0, -240) 0.08
animate move ball (0, -220) 0.08
animate move ball (0, -240) 0.08
animate move ball (0, -260) 0.08   # 最终停止
animate move ball (0, -260) 0.5    # 静止

# 渲染输出
render_frames 60 6.0 "output/physics_bounce"
```

---

## 最佳实践

### 1. 脚本组织
```r2g
# 总是以场景定义开始
scene 800 600 "动画标题"

# 创建所有对象
create circle obj1 20 (0, 0)
create rectangle obj2 100 50 (0, -200)

# 设置所有属性
set obj1.color = red
set obj2.color = green

# 动画序列（按时间顺序）
animate move obj1 (0, -150) 1.0
animate move obj1 (0, 100) 0.8
animate move obj1 (0, -150) 0.8

# 最后是渲染命令
render_frames 60 3.0 "output/animation"
```

### 2. 命名规范
- **描述性名称**: `ball`, `ground`, `leftwall`
- **避免数字开头**: 使用 `wall1` 而不是 `1wall`
- **使用下划线**: `my_object` 而不是 `myObject`

### 3. 坐标规划
```r2g
# 先计算场景边界
# 800×600场景: X[-400,400], Y[-300,300]

# 地面位置（留出对象半高度的空间）
create rectangle ground 600 20 (0, -290)    # 顶部在Y=-280

# 确保对象在边界内
create circle ball 15 (0, 250)              # 球心最高在Y=250，球顶在Y=265 < 300
```

### 4. 动画设计原则
- **物理真实性**: 遵循重力、惯性等物理规律
- **能量递减**: 弹跳高度逐渐降低
- **时间合理性**: 避免过快或过慢的动画
- **连续性**: 确保动画序列逻辑连贯

### 5. 性能优化
- **合理帧率**: 30-60fps适合大多数动画
- **适当时长**: 避免过长的动画（> 10秒）
- **对象数量**: 控制同时显示的对象数量
- **动画复杂度**: 避免过于复杂的动画序列

### 6. 颜色搭配
- **对比度**: 确保前景和背景有足够对比
- **一致性**: 使用统一的配色方案
- **功能性**: 不同类型的对象使用不同颜色

---

## 语法要点总结

1. **严格语法**: 所有命令必须严格按照语法格式
2. **坐标精确**: 使用1:1像素映射，坐标直接对应像素
3. **引号规则**: 字符串（标题、路径、文本内容）必须用引号
4. **注释支持**: 使用 `#` 开始注释行
5. **顺序重要**: 先创建对象，再设置属性，然后添加动画，最后渲染
6. **路径格式**: 输出路径使用正斜杠 `/` 或反斜杠 `\`

---

**记住**: 在编写任何R2G脚本之前，请务必参考此语法手册确保正确性！
