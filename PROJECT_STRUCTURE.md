# Render2Go 项目结构

## 📁 核心目录

### 🔧 Go 引擎
- `cmd/` - 命令行工具和主程序入口
- `core/` - 核心数据结构和接口
- `interpreter/` - 脚本解释器
- `renderer/` - 渲染引擎
- `scene/` - 场景管理
- `animation/` - 动画系统
- `geometry/` - 几何图形
- `math/` - 数学工具
- `colors/` - 颜色处理

### 🌐 Web 编辑器
- `web-editor/` - 完整的Web可视化编辑器
  - `index.html` - 主界面
  - `js/` - JavaScript 模块
  - `css/` - 样式文件

### 📝 脚本和输出
- `scripts/` - Render2Go 脚本文件存放目录
- `render2go.exe` - 编译后的可执行文件

## 🚀 快速开始

1. **使用Web编辑器**：
   ```bash
   cd web-editor
   python -m http.server 8000
   # 浏览器打开 http://localhost:8000
   ```

2. **使用命令行**：
   ```bash
   ./render2go.exe script.r2g
   ```

## ✨ 主要功能

- 🎨 **可视化编辑器** - 拖拽式创建动画
- 📝 **脚本语言** - 简洁的动画描述语言
- 🔄 **实时预览** - 所见即所得
- 📤 **多格式导出** - 支持各种输出格式
- 🎭 **变形控制** - 对象形状自由调整
