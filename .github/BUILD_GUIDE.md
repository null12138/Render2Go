# GitHub Actions 构建说明

本项目包含了三个GitHub Actions工作流，用于自动化构建和发布Render2Go项目。

## 工作流说明

### 1. 🚀 Release Workflow (`release.yml`)
**触发条件**：
- 推送带有版本标签时 (如: `git tag v1.0.0 && git push origin v1.0.0`)
- 手动触发 (在GitHub仓库的Actions页面)

**功能**：
- 构建所有平台的可执行文件
- 创建GitHub Release
- 自动上传所有平台的二进制文件

**支持平台**：
- Windows (amd64, arm64)
- Linux (amd64, arm64)  
- macOS (amd64, arm64)

### 2. 🔧 Build Test Workflow (`build.yml`)
**触发条件**：
- 推送到main或develop分支
- 创建Pull Request到main分支

**功能**：
- 运行测试
- 验证所有平台构建正常
- 确保代码质量

### 3. 🛠️ Manual Build Workflow (`manual-build.yml`)
**触发条件**：
- 手动触发 (可选择构建特定平台)

**功能**：
- 按需构建特定平台
- 快速测试构建
- 下载构建产物

## 使用方法

### 创建正式发布
1. 确保代码已经推送到main分支
2. 创建并推送版本标签：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions会自动构建并创建Release

### 手动构建测试
1. 进入GitHub仓库页面
2. 点击 "Actions" 选项卡
3. 选择 "Manual Build" 工作流
4. 点击 "Run workflow"
5. 选择要构建的平台 (all/windows/linux/macos)
6. 构建完成后在Artifacts中下载

### 查看构建状态
所有工作流的运行状态都可以在GitHub仓库的Actions页面查看。

## 构建配置

### 编译选项
- `CGO_ENABLED=0`: 禁用CGO，确保静态链接
- `-ldflags="-s -w"`: 减小二进制文件大小
- Go版本: 1.23

### 输出文件命名规则
- Windows: `render2go-windows-{arch}.exe`
- Linux: `render2go-linux-{arch}`
- macOS: `render2go-macos-{arch}`

其中 `{arch}` 为 `amd64` 或 `arm64`

## 注意事项

1. **Release创建**：只有推送标签或手动触发release workflow才会创建GitHub Release
2. **权限要求**：需要仓库的GITHUB_TOKEN权限来创建Release和上传文件
3. **标签格式**：建议使用语义化版本标签，如 `v1.0.0`, `v1.2.3`
4. **构建时间**：全平台构建大约需要5-10分钟

## 故障排除

### 构建失败
- 检查Go模块依赖是否正确
- 确认代码能在本地正常构建
- 查看Actions日志了解具体错误

### Release失败
- 确认有足够的仓库权限
- 检查标签格式是否正确
- 确认没有同名Release存在
