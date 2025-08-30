package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

// AST节点接口
type ASTNode interface {
	String() string
}

// 表达式节点
type Expression interface {
	ASTNode
	expressionNode()
}

// 语句节点
type Statement interface {
	ASTNode
	statementNode()
}

// 程序根节点
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// 基础表达式类型

// 标识符
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

// 数字字面量
type NumberLiteral struct {
	Token Token
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return fmt.Sprintf("%.2f", nl.Value) }

// 字符串字面量
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("\"%s\"", sl.Value) }

// 颜色字面量
type ColorLiteral struct {
	Token Token
	Value string
}

func (cl *ColorLiteral) expressionNode() {}
func (cl *ColorLiteral) String() string  { return cl.Value }

// 坐标表达式 (x, y)
type CoordinateExpression struct {
	X Expression
	Y Expression
}

func (ce *CoordinateExpression) expressionNode() {}
func (ce *CoordinateExpression) String() string {
	return fmt.Sprintf("(%s, %s)", ce.X.String(), ce.Y.String())
}

// 数组表达式 [1, 2, 3]
type ArrayExpression struct {
	Token    Token
	Elements []Expression
}

func (ae *ArrayExpression) expressionNode() {}
func (ae *ArrayExpression) String() string {
	var elements []string
	for _, e := range ae.Elements {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// 语句类型

// 场景声明语句
type SceneStatement struct {
	Token  Token
	Width  Expression
	Height Expression
	Name   Expression
}

func (ss *SceneStatement) statementNode() {}
func (ss *SceneStatement) String() string {
	return fmt.Sprintf("scene %s %s %s", ss.Width.String(), ss.Height.String(), ss.Name.String())
}

// 创建对象语句
type CreateStatement struct {
	Token      Token
	ObjectType Token
	Name       *Identifier
	Parameters []Expression
}

func (cs *CreateStatement) statementNode() {}
func (cs *CreateStatement) String() string {
	var params []string
	for _, p := range cs.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("create %s %s(%s)", cs.ObjectType.Literal, cs.Name.String(), strings.Join(params, ", "))
}

// 设置属性语句
type SetStatement struct {
	Token    Token
	Object   *Identifier
	Property Token
	Value    Expression
}

func (ss *SetStatement) statementNode() {}
func (ss *SetStatement) String() string {
	return fmt.Sprintf("set %s.%s = %s", ss.Object.String(), ss.Property.Literal, ss.Value.String())
}

// 动画语句
type AnimateStatement struct {
	Token      Token
	Animation  Token
	Object     *Identifier
	Parameters []Expression
	Duration   Expression
}

func (as *AnimateStatement) statementNode() {}
func (as *AnimateStatement) String() string {
	var params []string
	for _, p := range as.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("animate %s %s(%s) %s", as.Animation.Literal, as.Object.String(), strings.Join(params, ", "), as.Duration.String())
}

// 渲染语句
type RenderStatement struct {
	Token Token
}

func (rs *RenderStatement) statementNode() {}
func (rs *RenderStatement) String() string { return "render" }

// 渲染帧序列语句
type RenderFramesStatement struct {
	Token     Token
	FrameRate Expression // 帧率
	Duration  Expression // 时长
	OutputDir Expression // 输出目录
}

func (rfs *RenderFramesStatement) statementNode() {}
func (rfs *RenderFramesStatement) String() string {
	return fmt.Sprintf("render_frames %s %s %s", rfs.FrameRate.String(), rfs.Duration.String(), rfs.OutputDir.String())
}

// 保存语句
type SaveStatement struct {
	Token    Token
	Filename Expression
}

func (ss *SaveStatement) statementNode() {}
func (ss *SaveStatement) String() string {
	return fmt.Sprintf("save %s", ss.Filename.String())
}

// 导出视频语句
type ExportStatement struct {
	Token    Token
	Filename Expression
	FPS      Expression // 帧率，可选
	Duration Expression // 时长，可选
}

func (es *ExportStatement) statementNode() {}
func (es *ExportStatement) String() string {
	return fmt.Sprintf("export %s", es.Filename.String())
}

// 视频渲染语句
type VideoStatement struct {
	Token    Token
	Filename Expression
	FPS      Expression // 帧率
	Duration Expression // 时长
}

func (vs *VideoStatement) statementNode() {}
func (vs *VideoStatement) String() string {
	return fmt.Sprintf("video %s", vs.Filename.String())
}

// 等待语句
type WaitStatement struct {
	Token    Token
	Duration Expression
}

func (ws *WaitStatement) statementNode() {}
func (ws *WaitStatement) String() string {
	return fmt.Sprintf("wait %s", ws.Duration.String())
}

// 清空语句
type CleanStatement struct {
	Token Token
	Dirs  []Expression // 要清空的目录列表，如果为空则默认清空output和scripts
}

func (cs *CleanStatement) statementNode() {}
func (cs *CleanStatement) String() string {
	if len(cs.Dirs) == 0 {
		return "clean"
	}
	var dirs []string
	for _, d := range cs.Dirs {
		dirs = append(dirs, d.String())
	}
	return fmt.Sprintf("clean %s", strings.Join(dirs, ", "))
}

// 循环语句
type LoopStatement struct {
	Token      Token
	Count      Expression
	Statements []Statement
}

func (ls *LoopStatement) statementNode() {}
func (ls *LoopStatement) String() string {
	var stmts []string
	for _, s := range ls.Statements {
		stmts = append(stmts, s.String())
	}
	return fmt.Sprintf("loop %s {\n%s\n}", ls.Count.String(), strings.Join(stmts, "\n"))
}

// Parser 语法分析器
type Parser struct {
	lexer *Lexer

	curToken  Token
	peekToken Token

	errors []string
}

// NewParser 创建新的语法分析器
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// 读取两个标记，设置 curToken 和 peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken 移动到下一个标记
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseProgram 解析程序
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.curTokenIs(TOKEN_EOF) {
		// 跳过换行符
		if p.curTokenIs(TOKEN_NEWLINE) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement 解析语句
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case TOKEN_SCENE:
		return p.parseSceneStatement()
	case TOKEN_CREATE:
		return p.parseCreateStatement()
	case TOKEN_SET:
		return p.parseSetStatement()
	case TOKEN_ANIMATE:
		return p.parseAnimateStatement()
	case TOKEN_RENDER:
		return p.parseRenderStatement()
	case TOKEN_RENDER_FRAMES:
		return p.parseRenderFramesStatement()
	case TOKEN_SAVE:
		return p.parseSaveStatement()
	case TOKEN_EXPORT:
		return p.parseExportStatement()
	case TOKEN_VIDEO:
		return p.parseVideoStatement()
	case TOKEN_WAIT:
		return p.parseWaitStatement()
	case TOKEN_LOOP:
		return p.parseLoopStatement()
	case TOKEN_CLEAN:
		return p.parseCleanStatement()
	default:
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
}

// parseSceneStatement 解析场景语句
func (p *Parser) parseSceneStatement() *SceneStatement {
	stmt := &SceneStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Width = p.parseNumberLiteral()

	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Height = p.parseNumberLiteral()

	if !p.expectPeek(TOKEN_STRING) {
		return nil
	}
	stmt.Name = p.parseStringLiteral()

	return stmt
}

// parseCreateStatement 解析创建语句
func (p *Parser) parseCreateStatement() *CreateStatement {
	stmt := &CreateStatement{Token: p.curToken}

	if !p.expectPeekObjectType() {
		return nil
	}
	stmt.ObjectType = p.curToken

	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 解析所有参数（括号外的和括号内的）
	var parameters []Expression

	// 先解析括号外的参数
	for !p.peekTokenIs(TOKEN_LPAREN) && !p.peekTokenIs(TOKEN_EOF) && !p.peekTokenIs(TOKEN_NEWLINE) &&
		!p.peekTokenIs(TOKEN_CREATE) && !p.peekTokenIs(TOKEN_SET) && !p.peekTokenIs(TOKEN_ANIMATE) &&
		!p.peekTokenIs(TOKEN_RENDER) && !p.peekTokenIs(TOKEN_SAVE) && !p.peekTokenIs(TOKEN_EXPORT) &&
		!p.peekTokenIs(TOKEN_VIDEO) && !p.peekTokenIs(TOKEN_WAIT) && !p.peekTokenIs(TOKEN_LOOP) {
		p.nextToken()
		expr := p.parseExpression()
		if expr != nil {
			parameters = append(parameters, expr)
		}
	}

	// 然后解析所有的括号参数（作为坐标表达式）
	for p.peekTokenIs(TOKEN_LPAREN) {
		p.nextToken()
		coordExpr := p.parseExpression() // 这会调用 parseCoordinateExpression
		if coordExpr != nil {
			parameters = append(parameters, coordExpr)
		}
	}

	stmt.Parameters = parameters
	return stmt
} // parseSetStatement 解析设置语句
func (p *Parser) parseSetStatement() *SetStatement {
	stmt := &SetStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}
	stmt.Object = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(TOKEN_DOT) {
		return nil
	}

	if !p.expectPeekProperty() {
		return nil
	}
	stmt.Property = p.curToken

	if !p.expectPeek(TOKEN_ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression()

	return stmt
}

// parseAnimateStatement 解析动画语句
func (p *Parser) parseAnimateStatement() *AnimateStatement {
	stmt := &AnimateStatement{Token: p.curToken}

	if !p.expectPeekAnimationType() {
		return nil
	}
	stmt.Animation = p.curToken

	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}
	stmt.Object = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 根据动画类型解析不同的参数
	switch stmt.Animation.Type {
	case TOKEN_COLOR:
		// 颜色动画：一个字符串参数
		if p.peekTokenIs(TOKEN_STRING) {
			p.nextToken()
			expr := p.parseExpression()
			if expr != nil {
				stmt.Parameters = append(stmt.Parameters, expr)
			}
		}
	case TOKEN_PATH:
		// 路径动画：一个数组参数
		if p.peekTokenIs(TOKEN_LBRACKET) {
			p.nextToken()
			expr := p.parseExpression()
			if expr != nil {
				stmt.Parameters = append(stmt.Parameters, expr)
			}
		}
	case TOKEN_ELASTIC:
		// 弹性动画：两个参数（字符串和数字或负数）
		// 解析第一个参数（属性名）
		if p.peekTokenIs(TOKEN_STRING) {
			p.nextToken()
			expr := p.parseExpression()
			if expr != nil {
				stmt.Parameters = append(stmt.Parameters, expr)
			}
		}
		// 解析第二个参数（目标值，可能是负数）
		if p.peekTokenIs(TOKEN_MINUS) || p.peekTokenIs(TOKEN_NUMBER) {
			p.nextToken()
			expr := p.parseExpression()
			if expr != nil {
				stmt.Parameters = append(stmt.Parameters, expr)
			}
		}
	default:
		// 传统动画类型可能有坐标参数
		for !p.peekTokenIs(TOKEN_NUMBER) && !p.peekTokenIs(TOKEN_EOF) && !p.peekTokenIs(TOKEN_NEWLINE) {
			// 如果下一个token是左括号，则解析坐标表达式
			if p.peekTokenIs(TOKEN_LPAREN) {
				p.nextToken()
				expr := p.parseExpression() // 这会调用 parseCoordinateExpression
				if expr != nil {
					stmt.Parameters = append(stmt.Parameters, expr)
				}
			} else {
				break
			}
		}
	}

	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Duration = p.parseNumberLiteral()

	return stmt
}

// parseRenderStatement 解析渲染语句
func (p *Parser) parseRenderStatement() *RenderStatement {
	return &RenderStatement{Token: p.curToken}
}

// parseRenderFramesStatement 解析渲染帧序列语句
func (p *Parser) parseRenderFramesStatement() *RenderFramesStatement {
	stmt := &RenderFramesStatement{Token: p.curToken}

	// 解析帧率
	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.FrameRate = p.parseNumberLiteral()

	// 解析时长
	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Duration = p.parseNumberLiteral()

	// 解析输出目录
	if !p.expectPeek(TOKEN_STRING) {
		return nil
	}
	stmt.OutputDir = p.parseStringLiteral()

	return stmt
}

// parseSaveStatement 解析保存语句
func (p *Parser) parseSaveStatement() *SaveStatement {
	stmt := &SaveStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_STRING) {
		return nil
	}
	stmt.Filename = p.parseStringLiteral()

	return stmt
}

// parseExportStatement 解析导出语句
func (p *Parser) parseExportStatement() *ExportStatement {
	stmt := &ExportStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_STRING) {
		return nil
	}
	stmt.Filename = p.parseStringLiteral()

	// 可选的FPS和Duration参数
	if p.peekTokenIs(TOKEN_NUMBER) {
		p.nextToken()
		stmt.FPS = p.parseNumberLiteral()

		if p.peekTokenIs(TOKEN_NUMBER) {
			p.nextToken()
			stmt.Duration = p.parseNumberLiteral()
		}
	}

	return stmt
}

// parseVideoStatement 解析视频语句
func (p *Parser) parseVideoStatement() *VideoStatement {
	stmt := &VideoStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_STRING) {
		return nil
	}
	stmt.Filename = p.parseStringLiteral()

	// 必须的FPS参数
	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.FPS = p.parseNumberLiteral()

	// 必须的Duration参数
	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Duration = p.parseNumberLiteral()

	return stmt
}

// parseWaitStatement 解析等待语句
func (p *Parser) parseWaitStatement() *WaitStatement {
	stmt := &WaitStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Duration = p.parseNumberLiteral()

	return stmt
}

// parseCleanStatement 解析清空语句
func (p *Parser) parseCleanStatement() *CleanStatement {
	stmt := &CleanStatement{Token: p.curToken}
	stmt.Dirs = []Expression{}

	// 如果后面跟着字符串，则解析为目录列表
	if p.peekTokenIs(TOKEN_STRING) {
		for p.peekTokenIs(TOKEN_STRING) {
			p.nextToken()
			stmt.Dirs = append(stmt.Dirs, p.parseStringLiteral())

			if !p.peekTokenIs(TOKEN_COMMA) {
				break
			}
			p.nextToken() // 跳过逗号
		}
	}

	return stmt
}

// parseLoopStatement 解析循环语句
func (p *Parser) parseLoopStatement() *LoopStatement {
	stmt := &LoopStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_NUMBER) {
		return nil
	}
	stmt.Count = p.parseNumberLiteral()

	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}

	stmt.Statements = []Statement{}
	p.nextToken()

	for !p.curTokenIs(TOKEN_RBRACE) && !p.curTokenIs(TOKEN_EOF) {
		if p.curTokenIs(TOKEN_NEWLINE) {
			p.nextToken()
			continue
		}

		s := p.parseStatement()
		if s != nil {
			stmt.Statements = append(stmt.Statements, s)
		}
		p.nextToken()
	}

	return stmt
}

// parseExpression 解析表达式
func (p *Parser) parseExpression() Expression {
	switch p.curToken.Type {
	case TOKEN_IDENT:
		return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case TOKEN_NUMBER:
		return p.parseNumberLiteral()
	case TOKEN_STRING:
		return p.parseStringLiteral()
	case TOKEN_COLOR:
		return &ColorLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case TOKEN_LPAREN:
		return p.parseCoordinateExpression()
	case TOKEN_LBRACKET:
		return p.parseArrayExpression()
	default:
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
}

// parseNumberLiteral 解析数字字面量
func (p *Parser) parseNumberLiteral() *NumberLiteral {
	lit := &NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral 解析字符串字面量
func (p *Parser) parseStringLiteral() *StringLiteral {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseCoordinateExpression 解析坐标表达式
func (p *Parser) parseCoordinateExpression() *CoordinateExpression {
	p.nextToken()
	x := p.parseExpression()

	if !p.expectPeek(TOKEN_COMMA) {
		return nil
	}

	p.nextToken()
	y := p.parseExpression()

	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}

	return &CoordinateExpression{X: x, Y: y}
}

// parseArrayExpression 解析数组表达式
func (p *Parser) parseArrayExpression() *ArrayExpression {
	array := &ArrayExpression{Token: p.curToken}
	array.Elements = p.parseExpressionList(TOKEN_RBRACKET)
	return array
}

// parseExpressionList 解析表达式列表
func (p *Parser) parseExpressionList(end TokenType) []Expression {
	args := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression())

	for p.peekTokenIs(TOKEN_COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression())
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

// 辅助方法

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectPeekObjectType() bool {
	types := []TokenType{TOKEN_CIRCLE, TOKEN_TRIANGLE, TOKEN_RECT, TOKEN_LINE, TOKEN_ARROW, TOKEN_POLYGON, TOKEN_TEXT, TOKEN_MARKDOWN, TOKEN_TEX, TOKEN_MATHTEX, TOKEN_COORDINATE_SYSTEM}
	for _, t := range types {
		if p.peekTokenIs(t) {
			p.nextToken()
			return true
		}
	}

	typeNames := []string{"circle", "triangle", "rectangle", "line", "arrow", "polygon", "text", "markdown", "tex", "mathtex", "coordinate_system"}
	p.errors = append(p.errors, fmt.Sprintf("行 %d: 需要对象类型（%s），但得到了 '%s'",
		p.peekToken.Line, strings.Join(typeNames, ", "), p.peekToken.Literal))
	return false
}

func (p *Parser) expectPeekProperty() bool {
	properties := []TokenType{TOKEN_COLOR_PROP, TOKEN_SIZE_PROP, TOKEN_POSITION_PROP, TOKEN_OPACITY_PROP, TOKEN_WIDTH_PROP, TOKEN_HEIGHT_PROP, TOKEN_VERTEX_PROP, TOKEN_VERTICES_PROP}
	for _, t := range properties {
		if p.peekTokenIs(t) {
			p.nextToken()
			return true
		}
	}
	propNames := []string{"color_prop", "size", "position", "opacity", "width", "height", "vertex1", "vertex2", "vertex3", "vertices"}
	p.errors = append(p.errors, fmt.Sprintf("行 %d: 需要属性名（%s），但得到了 '%s'",
		p.peekToken.Line, strings.Join(propNames, ", "), p.peekToken.Literal))
	return false
}

func (p *Parser) expectPeekAnimationType() bool {
	animations := []TokenType{TOKEN_MOVE, TOKEN_SCALE, TOKEN_ROTATE, TOKEN_FADE_IN, TOKEN_FADE_OUT, TOKEN_BOUNCE, TOKEN_COLOR, TOKEN_PATH, TOKEN_ELASTIC}
	for _, t := range animations {
		if p.peekTokenIs(t) {
			p.nextToken()
			return true
		}
	}
	animNames := []string{"move", "scale", "rotate", "fadein", "fadeout", "bounce", "colorchange", "path", "elastic"}
	p.errors = append(p.errors, fmt.Sprintf("行 %d: 需要动画类型（%s），但得到了 '%s'",
		p.peekToken.Line, strings.Join(animNames, ", "), p.peekToken.Literal))
	return false
}

func (p *Parser) peekError(t TokenType) {
	var expected string

	switch t {
	case TOKEN_NUMBER:
		expected = "数字"
	case TOKEN_STRING:
		expected = "字符串"
	case TOKEN_IDENT:
		expected = "标识符"
	case TOKEN_ASSIGN:
		expected = "="
	case TOKEN_LPAREN:
		expected = "("
	case TOKEN_RPAREN:
		expected = ")"
	case TOKEN_LBRACE:
		expected = "{"
	case TOKEN_RBRACE:
		expected = "}"
	case TOKEN_LBRACKET:
		expected = "["
	case TOKEN_RBRACKET:
		expected = "]"
	case TOKEN_COMMA:
		expected = ","
	case TOKEN_DOT:
		expected = "."
	default:
		expected = fmt.Sprintf("%s", t)
	}

	msg := fmt.Sprintf("行 %d: 需要 %s，但得到了 '%s'",
		p.peekToken.Line, expected, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	var suggestion string

	switch t {
	case TOKEN_ILLEGAL:
		suggestion = "请检查输入是否存在非法字符"
	case TOKEN_EOF:
		suggestion = "脚本意外结束"
	default:
		suggestion = fmt.Sprintf("未知语句类型: %s", t)
	}

	msg := fmt.Sprintf("行 %d: %s", p.curToken.Line, suggestion)
	p.errors = append(p.errors, msg)
}

// Errors 返回解析错误
func (p *Parser) Errors() []string {
	return p.errors
}
