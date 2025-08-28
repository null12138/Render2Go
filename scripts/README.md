# 🎬 Render2Go Scripts Collection

这个文件夹包含了各种Render2Go脚本示例，现已合并为单一目录结构。

## 📚 脚本分类

### 🔰 基础教程 (Tutorials)
- `basic_shapes.r2g` - 基本几何图形演示
- `circle_demo.r2g` - 圆形动画基础
- `syntax_verification.r2g` - 语法验证测试
- `text_display_test.r2g` - 文本显示测试
- `text_refactor_test.r2g` - 文本重构测试

### 🧮 数学演示 (Mathematical Examples)
- `pythagoras.r2g` - 勾股定理演示
- `pythagoras_complete.r2g` - 完整勾股定理证明
- `pi_derivation.r2g` - π值推导演示
- `circle_circumference.r2g` - 圆周长计算
- `circle_circumference_advanced.r2g` - 高级圆周长动画
- `math_animation.r2g` - 综合数学动画

### 🌏 中文支持测试 (Chinese Support Tests)
- `chinese_test.r2g` - 基础中文字体测试
- `advanced_chinese_test.r2g` - 高级中文字符测试
- `math_chinese_test.r2g` - 中文数学公式演示

## 🚀 快速开始

```bash
# 运行基础图形演示
./render2go scripts/basic_shapes.r2g

# 运行数学动画
./render2go scripts/math_animation.r2g

# 运行中文测试
./render2go scripts/chinese_test.r2g

# 清空输出文件夹
./render2go -clean
```

## 📖 脚本说明

### 教程级别
- **初级** (🔰): 适合初学者，展示基本语法和功能
- **中级** (🔨): 包含复合操作和动画
- **高级** (⚡): 复杂的数学演示和动画序列

### 特色功能
- **几何图形**: 圆形、矩形、三角形、线条
- **文本渲染**: 支持中英文混合、数学公式
- **颜色系统**: RGB、HSV、命名颜色
- **动画效果**: 位置变换、颜色渐变、透明度变化
- **场景管理**: 多对象组合、批量操作

## 💡 使用建议

1. **学习路径**: 建议按照基础 → 数学 → 高级的顺序学习
2. **实验建议**: 尝试修改脚本参数，观察效果变化
3. **创作灵感**: 结合多个示例，创建自己的动画项目

## 📁 输出文件

所有脚本的输出文件都保存在 `output/` 文件夹中，按项目名称分类存储。

## 🧹 维护命令

使用 `./render2go -clean` 命令可以清空所有输出文件，释放磁盘空间。
