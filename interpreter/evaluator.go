package interpreter

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"render2go/animation"
	"render2go/colors"
	"render2go/core"
	"render2go/geometry"
	"render2go/internal/defaults"
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
	case *ExportStatement:
		return e.evalExportStatement(node)
	case *VideoStatement:
		return e.evalVideoStatement(node)
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
	case *ExportStatement:
		return &s.Token
	case *VideoStatement:
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
	// 以下功能已被移除以简化项目
	// case TOKEN_MARKDOWN:
	// 	obj, err = e.createMarkdown(stmt)
	// case TOKEN_TEX:
	// 	obj, err = e.createTex(stmt)
	// case TOKEN_MATHTEX:
	// 	obj, err = e.createMathTex(stmt)
	case TOKEN_COORDINATE_SYSTEM:
		obj, err = e.createCoordinateSystem(stmt)
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

	var radius float64
	var position *gmMath.Vector2

	// 检查第一个参数是坐标还是半径
	if coord, ok := stmt.Parameters[0].(*CoordinateExpression); ok {
		// 第一个参数是坐标，第二个应该是半径
		if len(stmt.Parameters) < 2 {
			return nil, fmt.Errorf("指定位置后还需要指定半径")
		}

		// 解析坐标
		x, err := e.evalExpression(coord.X)
		if err != nil {
			return nil, err
		}
		y, err := e.evalExpression(coord.Y)
		if err != nil {
			return nil, err
		}
		position = &gmMath.Vector2{X: x.(float64), Y: y.(float64)}

		// 解析半径
		radiusVal, err := e.evalExpression(stmt.Parameters[1])
		if err != nil {
			return nil, fmt.Errorf("解析圆形半径参数失败: %v", err)
		}
		radiusFloat, ok := radiusVal.(float64)
		if !ok {
			return nil, fmt.Errorf("圆形半径必须是数字，得到的是: %T", radiusVal)
		}
		radius = radiusFloat
	} else {
		// 第一个参数是半径
		radiusVal, err := e.evalExpression(stmt.Parameters[0])
		if err != nil {
			return nil, fmt.Errorf("解析圆形半径参数失败: %v", err)
		}
		radiusFloat, ok := radiusVal.(float64)
		if !ok {
			return nil, fmt.Errorf("圆形半径必须是数字，得到的是: %T", radiusVal)
		}
		radius = radiusFloat

		// 检查是否有位置参数
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
				position = &gmMath.Vector2{X: x.(float64), Y: y.(float64)}
			}
		}
	}

	if radius <= 0 {
		return nil, fmt.Errorf("圆形半径必须大于0，当前值: %v", radius)
	}

	circle := geometry.NewCircle(radius)

	// 如果有位置，设置位置
	if position != nil {
		circle.MoveTo(*position)
	}

	return circle, nil
}

// createTriangle 创建三角形
func (e *Evaluator) createTriangle(stmt *CreateStatement) (*geometry.Triangle, error) {
	// 检查参数数量，支持不同的创建方式
	numParams := len(stmt.Parameters)

	if numParams == 0 {
		return nil, fmt.Errorf("triangle requires at least one parameter")
	}

	// 方式1: 通过三个顶点创建 - triangle name (x1,y1) (x2,y2) (x3,y3)
	if numParams == 3 {
		var vertices [3]gmMath.Vector2
		for i, param := range stmt.Parameters {
			if coord, ok := param.(*CoordinateExpression); ok {
				x, err := e.evalExpression(coord.X)
				if err != nil {
					return nil, fmt.Errorf("invalid vertex %d X coordinate: %v", i+1, err)
				}
				y, err := e.evalExpression(coord.Y)
				if err != nil {
					return nil, fmt.Errorf("invalid vertex %d Y coordinate: %v", i+1, err)
				}
				vertices[i] = gmMath.Vector2{X: x.(float64), Y: y.(float64)}
			} else {
				return nil, fmt.Errorf("vertex %d must be a coordinate (x,y)", i+1)
			}
		}
		return geometry.NewTriangle(vertices[0], vertices[1], vertices[2]), nil
	}

	// 方式2: 通过类型和参数创建
	// 第一个参数可能是类型字符串或大小数字
	firstParam := stmt.Parameters[0]

	// 尝试解析第一个参数
	firstVal, err := e.evalExpression(firstParam)
	if err != nil {
		return nil, err
	}

	// 如果第一个参数是字符串，则表示三角形类型
	if typeStr, ok := firstVal.(string); ok {
		return e.createTriangleByType(typeStr, stmt.Parameters[1:])
	}

	// 否则按传统方式处理：第一个参数是大小
	size := firstVal.(float64)
	center := gmMath.Vector2{X: 0, Y: 0} // 默认中心

	// 如果有位置参数
	if numParams >= 2 {
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

	// 默认创建等腰直角三角形
	triangle := geometry.NewTriangleByCenter(center, size)
	return triangle, nil
}

// createTriangleByType 根据类型创建三角形
func (e *Evaluator) createTriangleByType(triangleType string, params []Expression) (*geometry.Triangle, error) {
	switch strings.ToLower(triangleType) {
	case "equilateral":
		// 等边三角形: triangle name "equilateral" sideLength (centerX, centerY)
		return e.createEquilateralTriangle(params)
	case "right":
		// 直角三角形: triangle name "right" width height (centerX, centerY)
		return e.createRightTriangle(params)
	case "isosceles":
		// 等腰直角三角形: triangle name "isosceles" size (centerX, centerY)
		return e.createIsoscelesTriangle(params)
	default:
		return nil, fmt.Errorf("unknown triangle type: %s. Supported types: equilateral, right, isosceles", triangleType)
	}
}

// createEquilateralTriangle 创建等边三角形
func (e *Evaluator) createEquilateralTriangle(params []Expression) (*geometry.Triangle, error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("equilateral triangle requires side length")
	}

	sideLengthVal, err := e.evalExpression(params[0])
	if err != nil {
		return nil, err
	}
	sideLength := sideLengthVal.(float64)

	center := gmMath.Vector2{X: 0, Y: 0}
	if len(params) >= 2 {
		if coord, ok := params[1].(*CoordinateExpression); ok {
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

	return geometry.NewEquilateralTriangle(center, sideLength), nil
}

// createRightTriangle 创建直角三角形
func (e *Evaluator) createRightTriangle(params []Expression) (*geometry.Triangle, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("right triangle requires width and height")
	}

	widthVal, err := e.evalExpression(params[0])
	if err != nil {
		return nil, err
	}
	heightVal, err := e.evalExpression(params[1])
	if err != nil {
		return nil, err
	}

	width := widthVal.(float64)
	height := heightVal.(float64)

	center := gmMath.Vector2{X: 0, Y: 0}
	if len(params) >= 3 {
		if coord, ok := params[2].(*CoordinateExpression); ok {
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

	return geometry.NewRightTriangle(center, width, height), nil
}

// createIsoscelesTriangle 创建等腰直角三角形
func (e *Evaluator) createIsoscelesTriangle(params []Expression) (*geometry.Triangle, error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("isosceles triangle requires size")
	}

	sizeVal, err := e.evalExpression(params[0])
	if err != nil {
		return nil, err
	}
	size := sizeVal.(float64)

	center := gmMath.Vector2{X: 0, Y: 0}
	if len(params) >= 2 {
		if coord, ok := params[1].(*CoordinateExpression); ok {
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

	return geometry.NewTriangleByCenter(center, size), nil
}

// createCoordinateSystem 创建坐标系
func (e *Evaluator) createCoordinateSystem(stmt *CreateStatement) (*geometry.CoordinateSystem, error) {
	numParams := len(stmt.Parameters)

	// 默认创建适应视口的坐标系
	if numParams == 0 {
		// 从场景中获取视口大小
		if e.scene != nil {
			width := float64(e.scene.GetWidth())
			height := float64(e.scene.GetHeight())
			return geometry.NewViewportCoordinateSystem(width, height), nil
		}
		return geometry.NewStandardCoordinateSystem(), nil
	}

	// 如果第一个参数是字符串，可能是预定义类型
	if numParams == 1 {
		if firstVal, err := e.evalExpression(stmt.Parameters[0]); err == nil {
			if typeStr, ok := firstVal.(string); ok {
				switch strings.ToLower(typeStr) {
				case "standard":
					return geometry.NewStandardCoordinateSystem(), nil
				case "small":
					return geometry.NewCoordinateSystem([2]float64{-5, 5}, [2]float64{-5, 5}, 1.0), nil
				case "large":
					return geometry.NewCoordinateSystem([2]float64{-20, 20}, [2]float64{-20, 20}, 2.0), nil
				case "viewport", "auto":
					// 自动适应视口
					if e.scene != nil {
						width := float64(e.scene.GetWidth())
						height := float64(e.scene.GetHeight())
						return geometry.NewViewportCoordinateSystem(width, height), nil
					}
					return geometry.NewStandardCoordinateSystem(), nil
				default:
					return nil, fmt.Errorf("unknown coordinate system type: %s. Supported: standard, small, large, viewport, auto", typeStr)
				}
			}
		}
	}

	// 自定义坐标系: coord_system name xMin xMax yMin yMax spacing
	if numParams >= 5 {
		xMinVal, err := e.evalExpression(stmt.Parameters[0])
		if err != nil {
			return nil, fmt.Errorf("invalid xMin: %v", err)
		}
		xMaxVal, err := e.evalExpression(stmt.Parameters[1])
		if err != nil {
			return nil, fmt.Errorf("invalid xMax: %v", err)
		}
		yMinVal, err := e.evalExpression(stmt.Parameters[2])
		if err != nil {
			return nil, fmt.Errorf("invalid yMin: %v", err)
		}
		yMaxVal, err := e.evalExpression(stmt.Parameters[3])
		if err != nil {
			return nil, fmt.Errorf("invalid yMax: %v", err)
		}
		spacingVal, err := e.evalExpression(stmt.Parameters[4])
		if err != nil {
			return nil, fmt.Errorf("invalid spacing: %v", err)
		}

		xMin := xMinVal.(float64)
		xMax := xMaxVal.(float64)
		yMin := yMinVal.(float64)
		yMax := yMaxVal.(float64)
		spacing := spacingVal.(float64)

		return geometry.NewCoordinateSystem([2]float64{xMin, xMax}, [2]float64{yMin, yMax}, spacing), nil
	}

	// 如果参数数量不匹配，返回错误
	return nil, fmt.Errorf("coordinate system requires 0 (standard), 1 (type), or 5+ (custom) parameters, got %d", numParams)
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
	// 检查参数数量：至少需要文本内容和字体大小
	if len(stmt.Parameters) < 2 {
		return nil, fmt.Errorf("文本对象需要至少2个参数：文本内容和字体大小")
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

	var size float64
	switch s := sizeVal.(type) {
	case float64:
		size = s
	case string:
		// 支持字体大小名称
		sizeName := strings.ToLower(s)
		sizeFound := false

		switch sizeName {
		case "tiny":
			size = defaults.FontSizes.Tiny
			sizeFound = true
		case "small":
			size = defaults.FontSizes.Small
			sizeFound = true
		case "normal":
			size = defaults.FontSizes.Normal
			sizeFound = true
		case "large":
			size = defaults.FontSizes.Large
			sizeFound = true
		case "huge":
			size = defaults.FontSizes.Huge
			sizeFound = true
		case "title":
			size = defaults.FontSizes.Title
			sizeFound = true
		}

		if !sizeFound {
			return nil, fmt.Errorf("未知字体大小名称: %s", s)
		}
	default:
		return nil, fmt.Errorf("字体大小必须是数字或字体大小名称")
	}

	if size <= 0 {
		return nil, fmt.Errorf("字体大小必须大于0")
	}

	// 创建文本对象
	textObj := geometry.NewText(text, size)

	// 如果提供了位置坐标（第3个参数），则设置位置
	if len(stmt.Parameters) >= 3 {
		if coord, ok := stmt.Parameters[2].(*CoordinateExpression); ok {
			x, err := e.evalExpression(coord.X)
			if err != nil {
				return nil, fmt.Errorf("解析X坐标失败: %v", err)
			}
			y, err := e.evalExpression(coord.Y)
			if err != nil {
				return nil, fmt.Errorf("解析Y坐标失败: %v", err)
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
	case TOKEN_VERTEX_PROP:
		return e.setVertex(obj, stmt.Property.Literal, value)
	case TOKEN_VERTICES_PROP:
		return e.setVertices(obj, value)
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
			// 首先检查defaults中的颜色名称
			colorFound := false
			colorName := strings.ToLower(v)

			// 直接匹配颜色名称到defaults.Colors字段
			switch colorName {
			case "black":
				c = defaults.Colors.Black
				colorFound = true
			case "white":
				c = defaults.Colors.White
				colorFound = true
			case "red":
				c = defaults.Colors.Red
				colorFound = true
			case "green":
				c = defaults.Colors.Green
				colorFound = true
			case "blue":
				c = defaults.Colors.Blue
				colorFound = true
			case "yellow":
				c = defaults.Colors.Yellow
				colorFound = true
			case "cyan":
				c = defaults.Colors.Cyan
				colorFound = true
			case "magenta":
				c = defaults.Colors.Magenta
				colorFound = true
			case "primary":
				c = defaults.Colors.Primary
				colorFound = true
			case "secondary":
				c = defaults.Colors.Secondary
				colorFound = true
			case "accent":
				c = defaults.Colors.Accent
				colorFound = true
			case "background":
				c = defaults.Colors.Background
				colorFound = true
			case "surface":
				c = defaults.Colors.Surface
				colorFound = true
			case "error":
				c = defaults.Colors.Error
				colorFound = true
			case "success":
				c = defaults.Colors.Success
				colorFound = true
			case "warning":
				c = defaults.Colors.Warning
				colorFound = true
			case "info":
				c = defaults.Colors.Info
				colorFound = true
			case "muted":
				c = defaults.Colors.Muted
				colorFound = true
			case "mathred":
				c = defaults.Colors.MathRed
				colorFound = true
			case "mathblue":
				c = defaults.Colors.MathBlue
				colorFound = true
			case "mathgreen":
				c = defaults.Colors.MathGreen
				colorFound = true
			case "mathorange":
				c = defaults.Colors.MathOrange
				colorFound = true
			case "mathpurple":
				c = defaults.Colors.MathPurple
				colorFound = true
			}

			if !colorFound {
				// 向后兼容旧的预定义颜色名称
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
		MoveTo(gmMath.Vector2) core.Mobject
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

// setVertex 设置三角形的单个顶点
func (e *Evaluator) setVertex(obj interface{}, property string, value interface{}) error {
	triangle, ok := obj.(*geometry.Triangle)
	if !ok {
		return e.newError("vertex properties are only supported for triangle objects")
	}

	// 确定顶点索引
	var vertexIndex int
	switch property {
	case "vertex1":
		vertexIndex = 0
	case "vertex2":
		vertexIndex = 1
	case "vertex3":
		vertexIndex = 2
	default:
		return e.newError("unknown vertex property: %s. Use vertex1, vertex2, or vertex3", property)
	}

	// 解析坐标值
	coord, ok := value.(*CoordinateExpression)
	if !ok {
		return e.newError("vertex must be a coordinate (x, y), got %T", value)
	}

	x, err := e.evalExpression(coord.X)
	if err != nil {
		return e.newError("invalid vertex X coordinate: %v", err)
	}
	y, err := e.evalExpression(coord.Y)
	if err != nil {
		return e.newError("invalid vertex Y coordinate: %v", err)
	}

	// 设置顶点
	newVertex := gmMath.Vector2{X: x.(float64), Y: y.(float64)}
	triangle.SetVertex(vertexIndex, newVertex)

	return nil
}

// setVertices 设置三角形的所有顶点
func (e *Evaluator) setVertices(obj interface{}, value interface{}) error {
	triangle, ok := obj.(*geometry.Triangle)
	if !ok {
		return e.newError("vertices property is only supported for triangle objects")
	}

	// 解析顶点数组
	array, ok := value.(*ArrayExpression)
	if !ok {
		return e.newError("vertices must be an array of coordinates [(x1,y1), (x2,y2), (x3,y3)], got %T", value)
	}

	if len(array.Elements) != 3 {
		return e.newError("triangle vertices array must contain exactly 3 coordinates, got %d", len(array.Elements))
	}

	var vertices [3]gmMath.Vector2
	for i, element := range array.Elements {
		coord, ok := element.(*CoordinateExpression)
		if !ok {
			return e.newError("vertex %d must be a coordinate (x,y), got %T", i+1, element)
		}

		x, err := e.evalExpression(coord.X)
		if err != nil {
			return e.newError("invalid vertex %d X coordinate: %v", i+1, err)
		}
		y, err := e.evalExpression(coord.Y)
		if err != nil {
			return e.newError("invalid vertex %d Y coordinate: %v", i+1, err)
		}

		vertices[i] = gmMath.Vector2{X: x.(float64), Y: y.(float64)}
	}

	// 设置所有顶点
	triangle.SetVertices(vertices[0], vertices[1], vertices[2])

	return nil
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

	// 获取图像 - 修复接口类型断言
	rendererInterface := e.scene.GetRenderer()
	if canvasRenderer, ok := rendererInterface.(*renderer.CanvasRenderer); ok {
		img := canvasRenderer.GetContext().Image()

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
	} else {
		return fmt.Errorf("不支持的渲染器类型")
	}

	return nil
}

// evalExportStatement 执行导出语句 - 导出序列帧动画
func (e *Evaluator) evalExportStatement(stmt *ExportStatement) error {
	if e.scene == nil {
		return fmt.Errorf("no scene defined")
	}

	filename, err := e.evalExpression(stmt.Filename)
	if err != nil {
		return err
	}

	// 默认参数
	fps := 30.0
	duration := 5.0

	// 解析可选参数
	if stmt.FPS != nil {
		fpsVal, err := e.evalExpression(stmt.FPS)
		if err != nil {
			return err
		}
		fps = fpsVal.(float64)
	}

	if stmt.Duration != nil {
		durationVal, err := e.evalExpression(stmt.Duration)
		if err != nil {
			return err
		}
		duration = durationVal.(float64)
	}

	return e.renderAnimationSequence(filename.(string), float64(fps), duration)
}

// evalVideoStatement 执行视频语句 - 直接生成视频文件
func (e *Evaluator) evalVideoStatement(stmt *VideoStatement) error {
	if e.scene == nil {
		return fmt.Errorf("no scene defined")
	}

	filename, err := e.evalExpression(stmt.Filename)
	if err != nil {
		return err
	}

	fpsVal, err := e.evalExpression(stmt.FPS)
	if err != nil {
		return err
	}

	durationVal, err := e.evalExpression(stmt.Duration)
	if err != nil {
		return err
	}

	fps := int(fpsVal.(float64))
	duration := durationVal.(float64)

	return e.renderVideoDirectly(filename.(string), float64(fps), duration)
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

/*
// createMarkdown 创建Markdown对象
func (e *Evaluator) createMarkdown(stmt *CreateStatement) (interface{}, error) {
	// 该功能已被移除以简化项目
	return nil, fmt.Errorf("markdown功能已被移除")
}

// createTex 创建TeX对象
func (e *Evaluator) createTex(stmt *CreateStatement) (interface{}, error) {
	// 该功能已被移除以简化项目
	return nil, fmt.Errorf("TeX功能已被移除")
}

// createTexWithLatex 创建TeX对象（使用latex库）
func (e *Evaluator) createTexWithLatex(stmt *CreateStatement) (interface{}, error) {
	// 该功能已被移除以简化项目
	return nil, fmt.Errorf("TeX功能已被移除")
}

// createMathTex 创建数学TeX对象
func (e *Evaluator) createMathTex(stmt *CreateStatement) (interface{}, error) {
	// 该功能已被移除以简化项目
	return nil, fmt.Errorf("数学TeX功能已被移除")
}
*/

// renderAnimationSequence 渲染动画序列为帧图片
func (e *Evaluator) renderAnimationSequence(filename string, fps, duration float64) error {
	if e.scene == nil {
		return fmt.Errorf("没有活动的场景")
	}

	// 获取渲染器
	rendererInterface := e.scene.GetRenderer()
	if rendererInterface == nil {
		return fmt.Errorf("没有设置渲染器")
	}

	canvasRenderer, ok := rendererInterface.(*renderer.CanvasRenderer)
	if !ok {
		return fmt.Errorf("渲染器类型不支持")
	}

	// 计算总帧数
	totalFrames := int(fps * duration)
	frameDir := fmt.Sprintf("%s_frames", strings.TrimSuffix(filename, ".mp4"))

	// 创建帧目录
	err := os.MkdirAll(frameDir, 0755)
	if err != nil {
		return fmt.Errorf("创建帧目录失败: %v", err)
	}

	// 准备动画时间轴
	dt := 1.0 / fps

	// 渲染每一帧
	for frame := 0; frame < totalFrames; frame++ {
		currentTime := float64(frame) * dt

		// 清空画布
		canvasRenderer.Clear(1.0, 1.0, 1.0)

		// 更新并渲染所有对象
		for _, obj := range e.scene.GetObjects() {
			// 如果对象支持动画更新
			if mobject, ok := obj.(interface{ UpdateAnimation(float64) }); ok {
				mobject.UpdateAnimation(currentTime)
			}
			// 渲染对象
			canvasRenderer.Render(obj)
		}

		// 保存当前帧
		framePath := fmt.Sprintf("%s/frame_%04d.png", frameDir, frame)
		err := canvasRenderer.SaveFrame(framePath)
		if err != nil {
			return fmt.Errorf("渲染第%d帧失败: %v", frame, err)
		}
	}

	// 使用FFmpeg合成视频
	ffmpegCmd := fmt.Sprintf("ffmpeg -r %.2f -i %s/frame_%%04d.png -c:v libx264 -pix_fmt yuv420p %s", fps, frameDir, filename)

	cmd := exec.Command("cmd", "/C", ffmpegCmd)
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("⚠️ FFmpeg未安装或执行失败，帧图片已保存到: %s\n", frameDir)
		fmt.Printf("您可以手动使用FFmpeg合成视频: %s\n", ffmpegCmd)
		return nil // 不返回错误，只是警告
	}

	// 清理临时帧文件
	os.RemoveAll(frameDir)

	fmt.Printf("动画视频已生成: %s\n", filename)
	return nil
} // renderVideoDirectly 直接渲染视频文件
func (e *Evaluator) renderVideoDirectly(filename string, fps, duration float64) error {
	// 对于直接视频渲染，我们也使用帧序列方法
	// 这确保了与现有渲染系统的兼容性
	return e.renderAnimationSequence(filename, fps, duration)
}
