package interpreter

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"render2go/animation"
	"render2go/colors"
	"render2go/core"
	"render2go/geometry"
	gmMath "render2go/math"
	"render2go/renderer"
	"render2go/scene"
	"strings"
	"time"
)

// Evaluator 执行引擎
type Evaluator struct {
	scene       *scene.Scene
	objects     map[string]interface{} // 存储创建的对象
	errors      []string
	projectName string // 项目名称
	currentLine int    // 当前执行行号
	fileName    string // 当前执行的文件名
}

// NewEvaluator 创建新的执行引擎
func NewEvaluator() *Evaluator {
	return &Evaluator{
		objects: make(map[string]interface{}),
		errors:  []string{},
	}
}

// Evaluate 执行程序
func (e *Evaluator) Evaluate(program *Program) error {
	for _, stmt := range program.Statements {
		err := e.evalStatement(stmt)
		if err != nil {
			e.errors = append(e.errors, err.Error())
			return err
		}
	}
	return nil
}

// evalStatement 执行语句
func (e *Evaluator) evalStatement(stmt Statement) error {
	// 更新当前执行的行号，用于错误定位
	if token := getStatementToken(stmt); token != nil {
		e.currentLine = token.Line
	}

	switch node := stmt.(type) {
	case *SceneStatement:
		return e.evalSceneStatement(node)
	case *CreateStatement:
		return e.evalCreateStatement(node)
	case *SetStatement:
		return e.evalSetStatement(node)
	case *AnimateStatement:
		return e.evalAnimateStatement(node)
	case *RenderStatement:
		return e.evalRenderStatement(node)
	case *SaveStatement:
		return e.evalSaveStatement(node)
	case *WaitStatement:
		return e.evalWaitStatement(node)
	case *LoopStatement:
		return e.evalLoopStatement(node)
	case *CleanStatement:
		return e.evalCleanStatement(node)
	default:
		return e.newError("未知语句类型: %T", stmt)
	}
}

// getStatementToken 获取语句的token，用于错误定位
func getStatementToken(stmt Statement) *Token {
	switch s := stmt.(type) {
	case *SceneStatement:
		return &s.Token
	case *CreateStatement:
		return &s.Token
	case *SetStatement:
		return &s.Token
	case *AnimateStatement:
		return &s.Token
	case *RenderStatement:
		return &s.Token
	case *SaveStatement:
		return &s.Token
	case *WaitStatement:
		return &s.Token
	case *LoopStatement:
		return &s.Token
	case *CleanStatement:
		return &s.Token
	default:
		return nil
	}
}

// newError 创建更详细的错误信息
func (e *Evaluator) newError(format string, args ...interface{}) error {
	errorMsg := fmt.Sprintf(format, args...)
	locationInfo := ""

	if e.fileName != "" {
		locationInfo = fmt.Sprintf("文件: %s, 行: %d", e.fileName, e.currentLine)
	} else {
		locationInfo = fmt.Sprintf("行: %d", e.currentLine)
	}

	fullError := fmt.Sprintf("执行错误 (%s): %s", locationInfo, errorMsg)
	fmt.Fprintf(os.Stderr, "❌ %s\n", fullError)
	return fmt.Errorf("%s", fullError)
}

// evalSceneStatement 执行场景语句
func (e *Evaluator) evalSceneStatement(stmt *SceneStatement) error {
	width, err := e.evalExpression(stmt.Width)
	if err != nil {
		return err
	}

	height, err := e.evalExpression(stmt.Height)
	if err != nil {
		return err
	}

	projectName, err := e.evalExpression(stmt.Name)
	if err != nil {
		return err
	}

	w := int(width.(float64))
	h := int(height.(float64))

	// 如果指定为0或负数，使用默认的1920*1080分辨率
	if w <= 0 {
		w = 1920
	}
	if h <= 0 {
		h = 1080
	}

	// 设置项目名称
	e.projectName = projectName.(string)

	// 创建场景
	sc := scene.NewScene(w, h)

	// 设置渲染器
	canvasRenderer := renderer.NewCanvasRenderer(w, h)
	canvasRenderer.SetAutoSaveProjectName(e.projectName) // 设置自动保存项目名称
	sc.SetRenderer(canvasRenderer)

	e.scene = sc
	return nil
}

// evalCreateStatement 执行创建语句
func (e *Evaluator) evalCreateStatement(stmt *CreateStatement) error {
	if e.scene == nil {
		return e.newError("未定义场景，请先使用 'scene' 命令创建场景")
	}

	var obj interface{}
	var err error

	switch stmt.ObjectType.Type {
	case TOKEN_CIRCLE:
		obj, err = e.createCircle(stmt)
	case TOKEN_TRIANGLE:
		obj, err = e.createTriangle(stmt)
	case TOKEN_RECT:
		obj, err = e.createRectangle(stmt)
	case TOKEN_LINE:
		obj, err = e.createLine(stmt)
	case TOKEN_ARROW:
		obj, err = e.createArrow(stmt)
	case TOKEN_POLYGON:
		obj, err = e.createPolygon(stmt)
	case TOKEN_TEXT:
		obj, err = e.createText(stmt)
	default:
		return e.newError("未知对象类型: %s", stmt.ObjectType.Literal)
	}

	if err != nil {
		return e.newError("创建对象 '%s' 失败: %v", stmt.Name.Value, err)
	}

	// 存储对象
	e.objects[stmt.Name.Value] = obj

	// 添加到场景
	if mobject, ok := obj.(core.Mobject); ok {
		e.scene.Add(mobject)
	}

	return nil
}

// createCircle 创建圆形
func (e *Evaluator) createCircle(stmt *CreateStatement) (*geometry.Circle, error) {
	if len(stmt.Parameters) < 1 {
		return nil, fmt.Errorf("创建圆形需要指定半径参数")
	}

	radiusVal, err := e.evalExpression(stmt.Parameters[0])
	if err != nil {
		return nil, fmt.Errorf("解析圆形半径参数失败: %v", err)
	}

	radius, ok := radiusVal.(float64)
	if !ok {
		return nil, fmt.Errorf("圆形半径必须是数字，得到的是: %T", radiusVal)
	}

	if radius <= 0 {
		return nil, fmt.Errorf("圆形半径必须大于0，当前值: %v", radius)
	}

	circle := geometry.NewCircle(radius)

	// 如果有位置参数
	if len(stmt.Parameters) >= 2 {
		if coord, ok := stmt.Parameters[1].(*CoordinateExpression); ok {
			x, err := e.evalExpression(coord.X)
			if err != nil {
				return nil, err
			}
			y, err := e.evalExpression(coord.Y)
			if err != nil {
				return nil, err
			}
			circle.MoveTo(gmMath.Vector2{X: x.(float64), Y: y.(float64)})
		}
	}

	return circle, nil
}

// createTriangle 创建三角形
func (e *Evaluator) createTriangle(stmt *CreateStatement) (*geometry.Triangle, error) {
	if len(stmt.Parameters) < 1 {
		return nil, fmt.Errorf("triangle requires size parameter")
	}

	sizeVal, err := e.evalExpression(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	size := sizeVal.(float64)
	center := gmMath.Vector2{X: 0, Y: 0} // 默认中心

	// 如果有位置参数
	if len(stmt.Parameters) >= 2 {
		if coord, ok := stmt.Parameters[1].(*CoordinateExpression); ok {
			x, err := e.evalExpression(coord.X)
			if err != nil {
				return nil, err
			}
			y, err := e.evalExpression(coord.Y)
			if err != nil {
				return nil, err
			}
			center = gmMath.Vector2{X: x.(float64), Y: y.(float64)}
		}
	}

	triangle := geometry.NewIsoscelesRightTriangle(center, size)
	return triangle, nil
}

// createRectangle 创建矩形
func (e *Evaluator) createRectangle(stmt *CreateStatement) (*geometry.Rectangle, error) {
	if len(stmt.Parameters) < 2 {
		return nil, fmt.Errorf("rectangle requires width and height parameters")
	}

	widthVal, err := e.evalExpression(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	heightVal, err := e.evalExpression(stmt.Parameters[1])
	if err != nil {
		return nil, err
	}

	width := widthVal.(float64)
	height := heightVal.(float64)
	rect := geometry.NewRectangle(width, height)

	// 如果有位置参数
	if len(stmt.Parameters) >= 3 {
		if coord, ok := stmt.Parameters[2].(*CoordinateExpression); ok {
			x, err := e.evalExpression(coord.X)
			if err != nil {
				return nil, err
			}
			y, err := e.evalExpression(coord.Y)
			if err != nil {
				return nil, err
			}
			rect.MoveTo(gmMath.Vector2{X: x.(float64), Y: y.(float64)})
		}
	}

	return rect, nil
}

// createLine 创建线条
func (e *Evaluator) createLine(stmt *CreateStatement) (*geometry.Line, error) {
	if len(stmt.Parameters) < 2 {
		return nil, fmt.Errorf("line requires start and end coordinate parameters")
	}

	start, ok1 := stmt.Parameters[0].(*CoordinateExpression)
	end, ok2 := stmt.Parameters[1].(*CoordinateExpression)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("line requires coordinate expressions")
	}

	startX, err := e.evalExpression(start.X)
	if err != nil {
		return nil, err
	}
	startY, err := e.evalExpression(start.Y)
	if err != nil {
		return nil, err
	}

	endX, err := e.evalExpression(end.X)
	if err != nil {
		return nil, err
	}
	endY, err := e.evalExpression(end.Y)
	if err != nil {
		return nil, err
	}

	startVec := gmMath.Vector2{X: startX.(float64), Y: startY.(float64)}
	endVec := gmMath.Vector2{X: endX.(float64), Y: endY.(float64)}

	return geometry.NewLine(startVec, endVec), nil
}

// createArrow 创建箭头
func (e *Evaluator) createArrow(stmt *CreateStatement) (*geometry.Arrow, error) {
	if len(stmt.Parameters) < 2 {
		return nil, fmt.Errorf("arrow requires start and end coordinate parameters")
	}

	start, ok1 := stmt.Parameters[0].(*CoordinateExpression)
	end, ok2 := stmt.Parameters[1].(*CoordinateExpression)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("arrow requires coordinate expressions")
	}

	startX, err := e.evalExpression(start.X)
	if err != nil {
		return nil, err
	}
	startY, err := e.evalExpression(start.Y)
	if err != nil {
		return nil, err
	}

	endX, err := e.evalExpression(end.X)
	if err != nil {
		return nil, err
	}
	endY, err := e.evalExpression(end.Y)
	if err != nil {
		return nil, err
	}

	startVec := gmMath.Vector2{X: startX.(float64), Y: startY.(float64)}
	endVec := gmMath.Vector2{X: endX.(float64), Y: endY.(float64)}

	return geometry.NewArrow(startVec, endVec), nil
}

// createPolygon 创建多边形
func (e *Evaluator) createPolygon(stmt *CreateStatement) (*geometry.Polygon, error) {
	if len(stmt.Parameters) < 1 {
		return nil, fmt.Errorf("polygon requires points array parameter")
	}

	arrayExpr, ok := stmt.Parameters[0].(*ArrayExpression)
	if !ok {
		return nil, fmt.Errorf("polygon requires array of coordinates")
	}

	var points []gmMath.Vector2
	for _, elem := range arrayExpr.Elements {
		coord, ok := elem.(*CoordinateExpression)
		if !ok {
			return nil, fmt.Errorf("polygon array must contain coordinate expressions")
		}

		x, err := e.evalExpression(coord.X)
		if err != nil {
			return nil, err
		}
		y, err := e.evalExpression(coord.Y)
		if err != nil {
			return nil, err
		}

		points = append(points, gmMath.Vector2{X: x.(float64), Y: y.(float64)})
	}

	return geometry.NewPolygon(points), nil
}

// createText 创建文本对象
func (e *Evaluator) createText(stmt *CreateStatement) (*geometry.Text, error) {
	// 检查参数数量：至少需要文本内容、字体大小和位置
	if len(stmt.Parameters) < 3 {
		return nil, fmt.Errorf("文本对象需要3个参数：文本内容、字体大小和位置坐标")
	}

	// 解析文本内容
	textVal, err := e.evalExpression(stmt.Parameters[0])
	if err != nil {
		return nil, fmt.Errorf("解析文本内容失败: %v", err)
	}
	text, ok := textVal.(string)
	if !ok {
		return nil, fmt.Errorf("文本内容必须是字符串")
	}

	// 解析字体大小
	sizeVal, err := e.evalExpression(stmt.Parameters[1])
	if err != nil {
		return nil, fmt.Errorf("解析字体大小失败: %v", err)
	}
	size, ok := sizeVal.(float64)
	if !ok {
		return nil, fmt.Errorf("字体大小必须是数字")
	}
	if size <= 0 {
		return nil, fmt.Errorf("字体大小必须大于0")
	}

	// 创建文本对象
	textObj := geometry.NewText(text, size)

	// 解析位置坐标
	if coord, ok := stmt.Parameters[2].(*CoordinateExpression); ok {
		x, err := e.evalExpression(coord.X)
		if err != nil {
			return nil, fmt.Errorf("解析X坐标失败: %v", err)
		}
		y, err := e.evalExpression(coord.Y)
		if err != nil {
			return nil, fmt.Errorf("解析Y坐标失败: %v", err)
		}

		xPos, ok1 := x.(float64)
		yPos, ok2 := y.(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("坐标必须是数字")
		}

		// 设置文本位置
		textObj.MoveTo(gmMath.Vector2{X: xPos, Y: yPos})
	} else {
		return nil, fmt.Errorf("第三个参数必须是坐标 (x, y)")
	}

	return textObj, nil
}

// evalSetStatement 执行设置语句
func (e *Evaluator) evalSetStatement(stmt *SetStatement) error {
	obj, exists := e.objects[stmt.Object.Value]
	if !exists {
		return e.newError("对象 '%s' 不存在", stmt.Object.Value)
	}

	value, err := e.evalExpression(stmt.Value)
	if err != nil {
		return e.newError("设置属性 '%s.%s' 时解析值失败: %v",
			stmt.Object.Value, stmt.Property.Literal, err)
	}

	switch stmt.Property.Type {
	case TOKEN_COLOR_PROP:
		return e.setColor(obj, value)
	case TOKEN_SIZE_PROP:
		return e.setSize(obj, value)
	case TOKEN_POSITION_PROP:
		return e.setPosition(obj, value)
	case TOKEN_OPACITY_PROP:
		return e.setOpacity(obj, value)
	case TOKEN_WIDTH_PROP:
		return e.setWidth(obj, value)
	case TOKEN_HEIGHT_PROP:
		return e.setHeight(obj, value)
	default:
		return fmt.Errorf("unknown property: %s", stmt.Property.Literal)
	}
}

// setColor 设置颜色
func (e *Evaluator) setColor(obj interface{}, value interface{}) error {
	var c color.RGBA

	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "#") {
			c = colors.HexToRGBA(v)
		} else {
			// 预定义颜色名称
			switch strings.ToLower(v) {
			case "deepblue":
				c = colors.DeepBlue
			case "midblue":
				c = colors.MidBlue
			case "purpleblue":
				c = colors.PurpleBlue
			case "cyanblue":
				c = colors.CyanBlue
			case "darkcolor":
				c = colors.DarkColor
			case "lightpurple":
				c = colors.LightPurple
			default:
				return e.newError("未知颜色名: %s", v)
			}
		}
	default:
		return e.newError("颜色必须是字符串（如 '#FF0000' 或颜色名），得到的是 %T", value)
	}

	if mobject, ok := obj.(interface{ SetColor(color.Color) }); ok {
		mobject.SetColor(c)
		return nil
	}

	return e.newError("对象不支持颜色属性")
}

// setPosition 设置位置
func (e *Evaluator) setPosition(obj interface{}, value interface{}) error {
	coord, ok := value.(*CoordinateExpression)
	if !ok {
		return e.newError("位置必须是坐标形式 (x, y)，得到的是 %T", value)
	}

	x, err := e.evalExpression(coord.X)
	if err != nil {
		return e.newError("解析X坐标失败: %v", err)
	}
	y, err := e.evalExpression(coord.Y)
	if err != nil {
		return e.newError("解析Y坐标失败: %v", err)
	}

	if mobject, ok := obj.(interface {
		MoveTo(gmMath.Vector2) interface{}
	}); ok {
		mobject.MoveTo(gmMath.Vector2{X: x.(float64), Y: y.(float64)})
		return nil
	}

	return fmt.Errorf("object does not support position property")
}

// setOpacity 设置透明度
func (e *Evaluator) setOpacity(obj interface{}, value interface{}) error {
	opacity, ok := value.(float64)
	if !ok {
		return fmt.Errorf("opacity must be a number")
	}

	if mobject, ok := obj.(interface{ SetFillOpacity(float64) }); ok {
		mobject.SetFillOpacity(opacity)
		return nil
	}

	return fmt.Errorf("object does not support opacity property")
}

// setSize, setWidth, setHeight 等其他属性设置方法...
func (e *Evaluator) setSize(obj interface{}, value interface{}) error {
	size, ok := value.(float64)
	if !ok {
		return fmt.Errorf("size must be a number")
	}

	if circle, ok := obj.(*geometry.Circle); ok {
		circle.SetRadius(size)
		return nil
	}

	return fmt.Errorf("object does not support size property")
}

func (e *Evaluator) setWidth(obj interface{}, value interface{}) error {
	// 实现宽度设置
	return fmt.Errorf("width property not yet implemented")
}

func (e *Evaluator) setHeight(obj interface{}, value interface{}) error {
	// 实现高度设置
	return fmt.Errorf("height property not yet implemented")
}

// evalAnimateStatement 执行动画语句
func (e *Evaluator) evalAnimateStatement(stmt *AnimateStatement) error {
	if e.scene == nil {
		return fmt.Errorf("no scene defined")
	}

	objName := stmt.Object.Value
	obj, ok := e.objects[objName]
	if !ok {
		return fmt.Errorf("object '%s' not found", objName)
	}

	mobj, ok := obj.(core.Mobject)
	if !ok {
		return fmt.Errorf("object '%s' is not animatable", objName)
	}

	durationVal, err := e.evalExpression(stmt.Duration)
	if err != nil {
		return err
	}
	duration := time.Duration(durationVal.(float64) * float64(time.Second))

	var anim animation.Animation
	switch stmt.Animation.Type {
	case TOKEN_MOVE:
		if len(stmt.Parameters) < 1 {
			return fmt.Errorf("move animation requires target position")
		}
		coordExpr, ok := stmt.Parameters[0].(*CoordinateExpression)
		if !ok {
			return fmt.Errorf("move parameter must be coordinate")
		}
		xVal, err := e.evalExpression(coordExpr.X)
		if err != nil {
			return err
		}
		yVal, err := e.evalExpression(coordExpr.Y)
		if err != nil {
			return err
		}
		endPos := gmMath.NewVector2(xVal.(float64), yVal.(float64))
		anim = animation.NewMoveToAnimation(mobj, endPos, duration)
	case TOKEN_SCALE:
		if len(stmt.Parameters) < 1 {
			return fmt.Errorf("scale animation requires scale factor")
		}
		scaleVal, err := e.evalExpression(stmt.Parameters[0])
		if err != nil {
			return err
		}
		anim = animation.NewScaleAnimation(mobj, scaleVal.(float64), duration)
	case TOKEN_ROTATE:
		if len(stmt.Parameters) < 1 {
			return fmt.Errorf("rotate animation requires angle")
		}
		angleVal, err := e.evalExpression(stmt.Parameters[0])
		if err != nil {
			return err
		}
		anim = animation.NewRotateAnimation(mobj, angleVal.(float64), duration)
	case TOKEN_FADE_IN:
		anim = animation.NewFadeInAnimation(mobj, duration)
	case TOKEN_FADE_OUT:
		anim = animation.NewFadeOutAnimation(mobj, duration)
	default:
		return fmt.Errorf("unsupported animation type: %s", stmt.Animation.Literal)
	}

	e.scene.PlayAnimation(anim)
	return nil
}

// 辅助函数：直接为脚本调用提供动画能力
// AnimateMove 移动动画
func (e *Evaluator) AnimateMove(objName string, x float64, y float64, duration float64) error {
	obj, ok := e.objects[objName]
	if !ok {
		return fmt.Errorf("object '%s' not found", objName)
	}
	mobj, ok := obj.(core.Mobject)
	if !ok {
		return fmt.Errorf("object '%s' is not animatable", objName)
	}
	endPos := gmMath.NewVector2(x, y)
	anim := animation.NewMoveToAnimation(mobj, endPos, time.Duration(duration*float64(time.Second)))
	e.scene.PlayAnimation(anim)
	return nil
}

// AnimateScale 缩放动画
func (e *Evaluator) AnimateScale(objName string, scale float64, duration float64) error {
	obj, ok := e.objects[objName]
	if !ok {
		return fmt.Errorf("object '%s' not found", objName)
	}
	mobj, ok := obj.(core.Mobject)
	if !ok {
		return fmt.Errorf("object '%s' is not animatable", objName)
	}
	anim := animation.NewScaleAnimation(mobj, scale, time.Duration(duration*float64(time.Second)))
	e.scene.PlayAnimation(anim)
	return nil
}

// AnimateRotate 旋转动画
func (e *Evaluator) AnimateRotate(objName string, angle float64, duration float64) error {
	obj, ok := e.objects[objName]
	if !ok {
		return fmt.Errorf("object '%s' not found", objName)
	}
	mobj, ok := obj.(core.Mobject)
	if !ok {
		return fmt.Errorf("object '%s' is not animatable", objName)
	}
	anim := animation.NewRotateAnimation(mobj, angle, time.Duration(duration*float64(time.Second)))
	e.scene.PlayAnimation(anim)
	return nil
}

// AnimateFadeIn 淡入动画
func (e *Evaluator) AnimateFadeIn(objName string, duration float64) error {
	obj, ok := e.objects[objName]
	if !ok {
		return fmt.Errorf("object '%s' not found", objName)
	}
	mobj, ok := obj.(core.Mobject)
	if !ok {
		return fmt.Errorf("object '%s' is not animatable", objName)
	}
	anim := animation.NewFadeInAnimation(mobj, time.Duration(duration*float64(time.Second)))
	e.scene.PlayAnimation(anim)
	return nil
}

// AnimateFadeOut 淡出动画
func (e *Evaluator) AnimateFadeOut(objName string, duration float64) error {
	obj, ok := e.objects[objName]
	if !ok {
		return fmt.Errorf("object '%s' not found", objName)
	}
	mobj, ok := obj.(core.Mobject)
	if !ok {
		return fmt.Errorf("object '%s' is not animatable", objName)
	}
	anim := animation.NewFadeOutAnimation(mobj, time.Duration(duration*float64(time.Second)))
	e.scene.PlayAnimation(anim)
	return nil
}

// evalRenderStatement 执行渲染语句
func (e *Evaluator) evalRenderStatement(stmt *RenderStatement) error {
	if e.scene == nil {
		return fmt.Errorf("no scene defined")
	}

	e.scene.RenderFrame()
	return nil
}

// evalSaveStatement 执行保存语句
func (e *Evaluator) evalSaveStatement(stmt *SaveStatement) error {
	if e.scene == nil {
		return fmt.Errorf("no scene defined")
	}

	filename, err := e.evalExpression(stmt.Filename)
	if err != nil {
		return err
	}

	// 创建输出目录结构
	outputDir := fmt.Sprintf("output/%s/frames", e.projectName)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// 构建完整的文件路径并确保PNG扩展名
	filenameStr := filename.(string)
	if !strings.HasSuffix(filenameStr, ".png") {
		filenameStr = filenameStr + ".png"
	}

	fullPath := filepath.Join(outputDir, filenameStr)

	// 使用统一的保存方法
	return e.saveImageFile(fullPath)
}

// saveImageFile 统一的图像文件保存方法，确保PNG扩展名
func (e *Evaluator) saveImageFile(fullPath string) error {
	// 确保文件路径有.png扩展名
	if !strings.HasSuffix(fullPath, ".png") {
		fullPath = fullPath + ".png"
	}

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败 '%s': %v", dir, err)
	}

	// 获取图像
	img := e.scene.GetRenderer().GetContext().Image()

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败 '%s': %v", fullPath, err)
	}
	defer file.Close()

	// 编码为PNG
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("PNG编码失败: %v", err)
	}

	return nil
}

// evalWaitStatement 执行等待语句
func (e *Evaluator) evalWaitStatement(stmt *WaitStatement) error {
	duration, err := e.evalExpression(stmt.Duration)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(duration.(float64)) * time.Second)
	return nil
}

// evalLoopStatement 执行循环语句
func (e *Evaluator) evalLoopStatement(stmt *LoopStatement) error {
	count, err := e.evalExpression(stmt.Count)
	if err != nil {
		return err
	}

	loopCount := int(count.(float64))
	for i := 0; i < loopCount; i++ {
		for _, s := range stmt.Statements {
			err := e.evalStatement(s)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// evalCleanStatement 执行清空指令
func (e *Evaluator) evalCleanStatement(stmt *CleanStatement) error {
	var dirsToClean []string

	// 如果没有指定目录，则默认清空output和scripts
	if len(stmt.Dirs) == 0 {
		dirsToClean = []string{"output", "scripts"}
	} else {
		// 解析指定的目录
		for _, dirExpr := range stmt.Dirs {
			dirValue, err := e.evalExpression(dirExpr)
			if err != nil {
				return e.newError("解析目录名失败: %v", err)
			}

			dirStr, ok := dirValue.(string)
			if !ok {
				return e.newError("目录名必须是字符串，得到的是: %T", dirValue)
			}
			dirsToClean = append(dirsToClean, dirStr)
		}
	}

	// 执行清空操作
	for _, dir := range dirsToClean {
		// 确保目录名合法，防止安全问题
		if strings.Contains(dir, "..") || strings.Contains(dir, "/") || strings.Contains(dir, "\\") {
			return e.newError("非法目录路径: %s", dir)
		}

		// 清空目录内容，但保留目录本身
		err := cleanDirectory(dir)
		if err != nil {
			return e.newError("清空目录 '%s' 失败: %v", dir, err)
		}

		fmt.Printf("已清空目录: %s\n", dir)
	}

	return nil
}

// cleanDirectory 清空指定目录内的所有文件和子目录
func cleanDirectory(dirPath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建它
		return os.MkdirAll(dirPath, 0755)
	}

	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// 删除所有文件和子目录
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// 递归删除子目录内容
			if err := os.RemoveAll(fullPath); err != nil {
				return err
			}
		} else {
			// 删除文件
			if err := os.Remove(fullPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// evalExpression 计算表达式
func (e *Evaluator) evalExpression(expr Expression) (interface{}, error) {
	switch node := expr.(type) {
	case *Identifier:
		return node.Value, nil
	case *NumberLiteral:
		return node.Value, nil
	case *StringLiteral:
		return node.Value, nil
	case *ColorLiteral:
		return node.Value, nil
	case *CoordinateExpression:
		return node, nil // 返回坐标表达式本身，由调用者处理
	case *ArrayExpression:
		return node, nil // 返回数组表达式本身，由调用者处理
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

// GetErrors 返回执行错误
func (e *Evaluator) GetErrors() []string {
	return e.errors
}

// GetScene 获取当前场景
func (e *Evaluator) GetScene() *scene.Scene {
	return e.scene
}

// GetObjects 获取所有对象
func (e *Evaluator) GetObjects() map[string]interface{} {
	return e.objects
}
