package interpreter

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
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
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
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

	// 设置项目名称
	e.projectName = projectName.(string)

	// 创建场景
	sc := scene.NewScene(w, h)

	// 设置渲染器
	canvasRenderer := renderer.NewCanvasRenderer(w, h)
	sc.SetRenderer(canvasRenderer)

	e.scene = sc
	return nil
}

// evalCreateStatement 执行创建语句
func (e *Evaluator) evalCreateStatement(stmt *CreateStatement) error {
	if e.scene == nil {
		return fmt.Errorf("no scene defined, use 'scene' command first")
	}

	var obj interface{}
	var err error

	switch stmt.ObjectType.Type {
	case TOKEN_CIRCLE:
		obj, err = e.createCircle(stmt)
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
		return fmt.Errorf("unknown object type: %s", stmt.ObjectType.Literal)
	}

	if err != nil {
		return err
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
		return nil, fmt.Errorf("circle requires radius parameter")
	}

	radiusVal, err := e.evalExpression(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	radius := radiusVal.(float64)
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

// createText 创建文字
func (e *Evaluator) createText(stmt *CreateStatement) (*geometry.Text, error) {
	if len(stmt.Parameters) < 2 {
		return nil, fmt.Errorf("text requires text and size parameters")
	}

	textVal, err := e.evalExpression(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	sizeVal, err := e.evalExpression(stmt.Parameters[1])
	if err != nil {
		return nil, err
	}

	text := textVal.(string)
	size := sizeVal.(float64)
	textObj := geometry.NewText(text, size)

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
			textObj.MoveTo(gmMath.Vector2{X: x.(float64), Y: y.(float64)})
		}
	}

	return textObj, nil
}

// evalSetStatement 执行设置语句
func (e *Evaluator) evalSetStatement(stmt *SetStatement) error {
	obj, exists := e.objects[stmt.Object.Value]
	if !exists {
		return fmt.Errorf("object '%s' not found", stmt.Object.Value)
	}

	value, err := e.evalExpression(stmt.Value)
	if err != nil {
		return err
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
				return fmt.Errorf("unknown color name: %s", v)
			}
		}
	default:
		return fmt.Errorf("color must be a string")
	}

	if mobject, ok := obj.(interface{ SetColor(color.Color) }); ok {
		mobject.SetColor(c)
		return nil
	}

	return fmt.Errorf("object does not support color property")
}

// setPosition 设置位置
func (e *Evaluator) setPosition(obj interface{}, value interface{}) error {
	coord, ok := value.(*CoordinateExpression)
	if !ok {
		return fmt.Errorf("position must be a coordinate")
	}

	x, err := e.evalExpression(coord.X)
	if err != nil {
		return err
	}
	y, err := e.evalExpression(coord.Y)
	if err != nil {
		return err
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
	// 动画功能的实现
	return fmt.Errorf("animation not yet implemented")
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

	// 构建完整的文件路径
	fullPath := filepath.Join(outputDir, filename.(string))
	return e.scene.SaveFrame(fullPath)
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
