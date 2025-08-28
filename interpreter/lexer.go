// Package interpreter 提供 Render2Go 动画脚本语言的解释器
// 支持独立运行，无需Go环境
package interpreter

import (
	"bufio"
	"io"
	"strings"
)

// TokenType 代表不同的标记类型
type TokenType int

const (
	// 基础标记
	TOKEN_ILLEGAL TokenType = iota
	TOKEN_EOF
	TOKEN_NEWLINE

	// 标识符和字面量
	TOKEN_IDENT  // 变量名、函数名
	TOKEN_NUMBER // 数字
	TOKEN_STRING // 字符串
	TOKEN_COLOR  // 颜色值 #RRGGBB

	// 关键字
	TOKEN_SCENE   // scene
	TOKEN_CREATE  // create
	TOKEN_SET     // set
	TOKEN_ANIMATE // animate
	TOKEN_RENDER  // render
	TOKEN_SAVE    // save
	TOKEN_WAIT    // wait
	TOKEN_LOOP    // loop
	TOKEN_IF      // if
	TOKEN_ELSE    // else
	TOKEN_END     // end

	// 几何类型
	TOKEN_CIRCLE  // circle
	TOKEN_RECT    // rectangle
	TOKEN_LINE    // line
	TOKEN_ARROW   // arrow
	TOKEN_POLYGON // polygon
	TOKEN_TEXT    // text

	// 动画类型
	TOKEN_MOVE     // move
	TOKEN_SCALE    // scale
	TOKEN_ROTATE   // rotate
	TOKEN_FADE_IN  // fadein
	TOKEN_FADE_OUT // fadeout

	// 属性
	TOKEN_COLOR_PROP    // color
	TOKEN_SIZE_PROP     // size
	TOKEN_POSITION_PROP // position
	TOKEN_OPACITY_PROP  // opacity
	TOKEN_WIDTH_PROP    // width
	TOKEN_HEIGHT_PROP   // height

	// 运算符
	TOKEN_ASSIGN   // =
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_MULTIPLY // *
	TOKEN_DIVIDE   // /

	// 分隔符
	TOKEN_COMMA     // ,
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_LBRACE    // {
	TOKEN_RBRACE    // }
	TOKEN_LBRACKET  // [
	TOKEN_RBRACKET  // ]
	TOKEN_DOT       // .
	TOKEN_COLON     // :
	TOKEN_SEMICOLON // ;
)

// Token 表示一个词法标记
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer 词法分析器
type Lexer struct {
	input        string
	position     int  // 当前位置
	readPosition int  // 下一个读取位置
	ch           byte // 当前字符
	line         int  // 当前行号
	column       int  // 当前列号
}

// keywords 关键字映射表
var keywords = map[string]TokenType{
	"scene":     TOKEN_SCENE,
	"create":    TOKEN_CREATE,
	"set":       TOKEN_SET,
	"animate":   TOKEN_ANIMATE,
	"render":    TOKEN_RENDER,
	"save":      TOKEN_SAVE,
	"wait":      TOKEN_WAIT,
	"loop":      TOKEN_LOOP,
	"if":        TOKEN_IF,
	"else":      TOKEN_ELSE,
	"end":       TOKEN_END,
	"circle":    TOKEN_CIRCLE,
	"rectangle": TOKEN_RECT,
	"line":      TOKEN_LINE,
	"arrow":     TOKEN_ARROW,
	"polygon":   TOKEN_POLYGON,
	"text":      TOKEN_TEXT,
	"move":      TOKEN_MOVE,
	"scale":     TOKEN_SCALE,
	"rotate":    TOKEN_ROTATE,
	"fadein":    TOKEN_FADE_IN,
	"fadeout":   TOKEN_FADE_OUT,
	"color":     TOKEN_COLOR_PROP,
	"size":      TOKEN_SIZE_PROP,
	"position":  TOKEN_POSITION_PROP,
	"opacity":   TOKEN_OPACITY_PROP,
	"width":     TOKEN_WIDTH_PROP,
	"height":    TOKEN_HEIGHT_PROP,
}

// NewLexer 创建新的词法分析器
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar 读取下一个字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL 表示 EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar 查看下一个字符但不移动位置
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace 跳过空白字符（除换行符）
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment 跳过注释（// 到行尾）
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readIdentifier 读取标识符
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber 读取数字
func (l *Lexer) readNumber() string {
	position := l.position
	hasDot := false

	for isDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString 读取字符串
func (l *Lexer) readString() string {
	position := l.position + 1 // 跳过开始的引号
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// readColor 读取颜色值 #RRGGBB
func (l *Lexer) readColor() string {
	position := l.position
	for isHexDigit(l.ch) || l.ch == '#' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// NextToken 获取下一个标记
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = Token{Type: TOKEN_ASSIGN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '+':
		tok = Token{Type: TOKEN_PLUS, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '-':
		tok = Token{Type: TOKEN_MINUS, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '*':
		tok = Token{Type: TOKEN_MULTIPLY, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken() // 递归获取下一个标记
		}
		tok = Token{Type: TOKEN_DIVIDE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ',':
		tok = Token{Type: TOKEN_COMMA, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '(':
		tok = Token{Type: TOKEN_LPAREN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ')':
		tok = Token{Type: TOKEN_RPAREN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '{':
		tok = Token{Type: TOKEN_LBRACE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '}':
		tok = Token{Type: TOKEN_RBRACE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '[':
		tok = Token{Type: TOKEN_LBRACKET, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ']':
		tok = Token{Type: TOKEN_RBRACKET, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '.':
		tok = Token{Type: TOKEN_DOT, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ':':
		tok = Token{Type: TOKEN_COLON, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ';':
		tok = Token{Type: TOKEN_SEMICOLON, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '\n':
		tok = Token{Type: TOKEN_NEWLINE, Literal: "\\n", Line: l.line, Column: l.column}
	case '"':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.Column = l.column
	case '#':
		tok.Type = TOKEN_COLOR
		tok.Literal = l.readColor()
		tok.Line = l.line
		tok.Column = l.column
		return tok // 不调用 readChar()，因为 readColor 已经处理了
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column
			return tok // 不调用 readChar()，因为 readIdentifier 已经处理了
		} else if isDigit(l.ch) {
			tok.Type = TOKEN_NUMBER
			tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			return tok // 不调用 readChar()，因为 readNumber 已经处理了
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	}

	l.readChar()
	return tok
}

// lookupIdent 查找标识符类型
func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}

// isLetter 检查是否为字母
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isDigit 检查是否为数字
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isHexDigit 检查是否为十六进制数字
func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('A' <= ch && ch <= 'F') || ('a' <= ch && ch <= 'f')
}

// TokenizeFile 从文件中读取并标记化
func TokenizeFile(reader io.Reader) ([]Token, error) {
	scanner := bufio.NewScanner(reader)
	var content strings.Builder

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	lexer := NewLexer(content.String())
	var tokens []Token

	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == TOKEN_EOF {
			break
		}
	}

	return tokens, nil
}

// TokenString 返回标记类型的字符串表示
func (tt TokenType) String() string {
	switch tt {
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_NEWLINE:
		return "NEWLINE"
	case TOKEN_IDENT:
		return "IDENT"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_COLOR:
		return "COLOR"
	case TOKEN_SCENE:
		return "SCENE"
	case TOKEN_CREATE:
		return "CREATE"
	case TOKEN_SET:
		return "SET"
	case TOKEN_ANIMATE:
		return "ANIMATE"
	case TOKEN_RENDER:
		return "RENDER"
	case TOKEN_SAVE:
		return "SAVE"
	case TOKEN_WAIT:
		return "WAIT"
	case TOKEN_LOOP:
		return "LOOP"
	case TOKEN_IF:
		return "IF"
	case TOKEN_ELSE:
		return "ELSE"
	case TOKEN_END:
		return "END"
	case TOKEN_CIRCLE:
		return "CIRCLE"
	case TOKEN_RECT:
		return "RECTANGLE"
	case TOKEN_LINE:
		return "LINE"
	case TOKEN_ARROW:
		return "ARROW"
	case TOKEN_POLYGON:
		return "POLYGON"
	case TOKEN_TEXT:
		return "TEXT"
	case TOKEN_MOVE:
		return "MOVE"
	case TOKEN_SCALE:
		return "SCALE"
	case TOKEN_ROTATE:
		return "ROTATE"
	case TOKEN_FADE_IN:
		return "FADE_IN"
	case TOKEN_FADE_OUT:
		return "FADE_OUT"
	case TOKEN_COLOR_PROP:
		return "COLOR_PROP"
	case TOKEN_SIZE_PROP:
		return "SIZE_PROP"
	case TOKEN_POSITION_PROP:
		return "POSITION_PROP"
	case TOKEN_OPACITY_PROP:
		return "OPACITY_PROP"
	case TOKEN_WIDTH_PROP:
		return "WIDTH_PROP"
	case TOKEN_HEIGHT_PROP:
		return "HEIGHT_PROP"
	case TOKEN_ASSIGN:
		return "ASSIGN"
	case TOKEN_PLUS:
		return "PLUS"
	case TOKEN_MINUS:
		return "MINUS"
	case TOKEN_MULTIPLY:
		return "MULTIPLY"
	case TOKEN_DIVIDE:
		return "DIVIDE"
	case TOKEN_COMMA:
		return "COMMA"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_LBRACKET:
		return "LBRACKET"
	case TOKEN_RBRACKET:
		return "RBRACKET"
	case TOKEN_DOT:
		return "DOT"
	case TOKEN_COLON:
		return "COLON"
	case TOKEN_SEMICOLON:
		return "SEMICOLON"
	default:
		return "UNKNOWN"
	}
}
