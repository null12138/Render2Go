# Render2Go 高级插值功能说明文档

## 概述

Render2Go 现在支持多种高级插值算法，以实现更流畅、更自然的动画效果。通过改进帧间补帧技术，我们能够生成更高质量的动画序列。

## 支持的插值类型

### 1. 线性插值 (Linear Interpolation)
- **特点**: 匀速运动
- **适用场景**: 简单的直线运动
- **视觉效果**: 恒定速度，无加速或减速

### 2. 平滑插值 (Smooth Interpolation)
- **特点**: 缓入缓出效果
- **适用场景**: 大多数动画效果
- **视觉效果**: 开始和结束时较慢，中间较快

### 3. 弹性插值 (Elastic Interpolation)
- **特点**: 带有弹性效果
- **适用场景**: 需要弹性或回弹效果的动画
- **视觉效果**: 超过目标位置后回弹

### 4. 弹跳插值 (Bounce Interpolation)
- **特点**: 带有弹跳效果
- **适用场景**: 物理模拟，如球体落地弹跳
- **视觉效果**: 模拟真实物理弹跳行为

## 技术实现

### 帧率提升
- **默认帧率**: 从30fps提升到60fps
- **效果**: 动画更加流畅，减少卡顿感
- **性能**: 在现代硬件上运行良好

### 插值器实现
所有插值器都实现了以下接口：
```go
type Interpolator interface {
    Interpolate(start, end gmMath.Vector2, t float64) gmMath.Vector2
    InterpolateFloat(start, end, t float64) float64
}
```

### 使用方法
在动画代码中，可以通过以下方式设置插值类型：
```go
anim.SetInterpolation(animation.Smooth) // 设置为平滑插值
```

## 最佳实践

### 1. 选择合适的插值类型
- 对于自然运动，推荐使用平滑插值
- 对于弹性效果，使用弹性插值
- 对于物理模拟，使用弹跳插值
- 对于简单直线运动，使用线性插值

### 2. 帧率设置
- 一般动画推荐使用60fps
- 复杂场景可以使用30fps以节省资源
- 简单动画可以使用更高帧率（如120fps）

### 3. 性能优化
- 避免在动画中进行复杂计算
- 预计算关键值以提高性能
- 合理使用缓动函数

## 示例代码

创建一个使用不同插值类型的动画示例：
```r2g
# 创建场景
scene 1920 1080 "interpolation_demo"

# 创建对象
create circle ball 1
set ball.color = "#3498DB"
set ball.position = (-5, 0)

# 使用不同插值类型的动画
animate move ball (5, 0) 2.0  # 默认平滑插值

# 等待动画完成
wait 2.0

# 渲染并保存
render
save "interpolation_demo_final"
```

## 导出视频

使用export命令导出视频：
```r2g
export "animation.mp4" 60 2.0
```

参数说明：
- `"animation.mp4"`: 输出文件名
- `60`: 帧率（fps）
- `2.0`: 动画时长（秒）

## 注意事项

1. 需要安装FFmpeg才能导出视频文件
2. 帧序列会自动保存到output目录
3. 如果FFmpeg未安装，仍会生成帧序列文件
4. 建议使用60fps以获得最佳视觉效果