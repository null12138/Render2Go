package renderer

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"render2go/scene"
	"time"
)

// FrameSequenceRenderer 序列帧渲染器
type FrameSequenceRenderer struct {
	outputDir    string
	frameRate    int
	totalFrames  int
	currentFrame int
	width        int
	height       int
}

// NewFrameSequenceRenderer 创建新的序列帧渲染器
func NewFrameSequenceRenderer(outputDir string, frameRate int, duration float64, width, height int) *FrameSequenceRenderer {
	// 使用默认60fps以获得更流畅的动画效果
	if frameRate <= 0 {
		frameRate = 60
	}

	totalFrames := int(duration * float64(frameRate))

	// 创建输出目录
	os.MkdirAll(outputDir, 0755)

	return &FrameSequenceRenderer{
		outputDir:    outputDir,
		frameRate:    frameRate,
		totalFrames:  totalFrames,
		currentFrame: 0,
		width:        width,
		height:       height,
	}
}

// RenderFrame 渲染单帧
func (fsr *FrameSequenceRenderer) RenderFrame(scn *scene.Scene, frameIndex int) error {
	// 设置场景时间
	timePos := float64(frameIndex) / float64(fsr.frameRate)
	scn.SetCurrentTime(timePos)

	// 渲染场景到图像
	img := fsr.renderSceneToImage(scn)

	// 保存帧图像
	filename := fmt.Sprintf("frame_%06d.png", frameIndex)
	filepath := filepath.Join(fsr.outputDir, filename)

	return fsr.saveImage(img, filepath)
}

// RenderSequence 渲染完整序列
func (fsr *FrameSequenceRenderer) RenderSequence(scn *scene.Scene) error {
	fmt.Printf("🎬 开始渲染序列帧...\n")
	fmt.Printf("   输出目录: %s\n", fsr.outputDir)
	fmt.Printf("   帧率: %d fps\n", fsr.frameRate)
	fmt.Printf("   总帧数: %d\n", fsr.totalFrames)

	start := time.Now()

	for i := 0; i < fsr.totalFrames; i++ {
		if err := fsr.RenderFrame(scn, i); err != nil {
			return fmt.Errorf("渲染第 %d 帧失败: %v", i, err)
		}

		// 显示进度
		if i%10 == 0 || i == fsr.totalFrames-1 {
			progress := float64(i+1) / float64(fsr.totalFrames) * 100
			fmt.Printf("   进度: %.1f%% (%d/%d)\n", progress, i+1, fsr.totalFrames)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ 序列帧渲染完成！耗时: %v\n", elapsed)

	// 生成FFmpeg命令提示
	fsr.generateFFmpegCommand()

	return nil
}

// renderSceneToImage 将场景渲染为图像
func (fsr *FrameSequenceRenderer) renderSceneToImage(scn *scene.Scene) image.Image {
	// 创建临时渲染器
	tempRenderer := NewCanvasRenderer(fsr.width, fsr.height)

	// 设置背景色
	backgroundColor := scn.GetBackgroundColor()
	tempRenderer.Clear(backgroundColor[0], backgroundColor[1], backgroundColor[2])

	// 设置坐标系统
	objects := scn.GetObjects()
	tempRenderer.SetupCoordinateSystem(objects)

	// 渲染所有对象
	for _, obj := range objects {
		tempRenderer.Render(obj)
	}

	// 获取渲染结果作为图像
	return tempRenderer.GetImage()
}

// saveImage 保存图像到文件
func (fsr *FrameSequenceRenderer) saveImage(img image.Image, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// generateFFmpegCommand 生成FFmpeg转换命令
func (fsr *FrameSequenceRenderer) generateFFmpegCommand() {
	fmt.Printf("\n🎥 使用FFmpeg生成视频:\n")

	// 生成MP4命令
	mp4Command := fmt.Sprintf(
		"ffmpeg -framerate %d -i \"%s/frame_%%06d.png\" -c:v libx264 -pix_fmt yuv420p output.mp4",
		fsr.frameRate, fsr.outputDir)

	// 生成GIF命令
	gifCommand := fmt.Sprintf(
		"ffmpeg -framerate %d -i \"%s/frame_%%06d.png\" -vf \"palettegen\" palette.png && ffmpeg -framerate %d -i \"%s/frame_%%06d.png\" -i palette.png -lavfi \"paletteuse\" output.gif",
		fsr.frameRate, fsr.outputDir, fsr.frameRate, fsr.outputDir)

	fmt.Printf("\n📹 生成MP4视频:\n%s\n", mp4Command)
	fmt.Printf("\n🎞️ 生成GIF动画:\n%s\n", gifCommand)

	// 保存命令到文件
	cmdFile, err := os.Create(filepath.Join(fsr.outputDir, "generate_video.bat"))
	if err == nil {
		defer cmdFile.Close()
		cmdFile.WriteString("@echo off\n")
		cmdFile.WriteString("echo 正在生成视频...\n")
		cmdFile.WriteString(mp4Command + "\n")
		cmdFile.WriteString("echo 视频生成完成: output.mp4\n")
		cmdFile.WriteString("pause\n")
		fmt.Printf("\n💾 批处理文件已保存: %s/generate_video.bat\n", fsr.outputDir)
	}
}

// GetFrameCount 获取总帧数
func (fsr *FrameSequenceRenderer) GetFrameCount() int {
	return fsr.totalFrames
}

// GetFrameRate 获取帧率
func (fsr *FrameSequenceRenderer) GetFrameRate() int {
	return fsr.frameRate
}
