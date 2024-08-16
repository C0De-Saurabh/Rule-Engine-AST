package ast

import (
	"strconv"
	"strings"
)

// Token types
const (
	TokenField    = "FIELD"
	TokenOperator = "OPERATOR"
	TokenValue    = "VALUE"
	TokenParenL   = "PAREN_L"
	TokenParenR   = "PAREN_R"
	TokenAnd      = "AND"
	TokenOr       = "OR"
)

// Token represents a lexical token.
type Token struct {
	Type  string
	Value string
}

// Tokenizer splits the input expression into tokens.
func Tokenizer(input string) []Token {
	tokens := []Token{}
	i := 0
	for i < len(input) {
		char := input[i]

		switch {
		case char == '(':
			tokens = append(tokens, Token{Type: TokenParenL, Value: "("})
			i++
		case char == ')':
			tokens = append(tokens, Token{Type: TokenParenR, Value: ")"})
			i++
		case strings.HasPrefix(input[i:], "AND"), strings.HasPrefix(input[i:], "&&"):
			tokens = append(tokens, Token{Type: TokenAnd, Value: "AND"})
			i += 3
		case strings.HasPrefix(input[i:], "OR"), strings.HasPrefix(input[i:], "||"):
			tokens = append(tokens, Token{Type: TokenOr, Value: "OR"})
			i += 2
		case char == ' ':
			i++
		case char == '\'':
			start := i + 1
			i = start
			for i < len(input) && input[i] != '\'' {
				i++
			}
			tokens = append(tokens, Token{Type: TokenValue, Value: input[start:i]})
			i++
		case char >= '0' && char <= '9':
			start := i
			for i < len(input) && (input[i] >= '0' && input[i] <= '9') {
				i++
			}
			tokens = append(tokens, Token{Type: TokenValue, Value: input[start:i]})
		case char >= 'a' && char <= 'z':
			start := i
			for i < len(input) && (input[i] >= 'a' && input[i] <= 'z') {
				i++
			}
			tokens = append(tokens, Token{Type: TokenField, Value: input[start:i]})
		case strings.Contains("><=", string(char)):
			start := i
			if char == '=' && i+1 < len(input) && input[i+1] == '=' {
				i++
			}
			i++
			tokens = append(tokens, Token{Type: TokenOperator, Value: input[start:i]})
		default:
			i++
		}
	}
	return tokens
}

// Parser is responsible for creating the AST from tokens.
type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) currentToken() Token {
	if p.pos >= len(p.tokens) {
		return Token{}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) parsePrimary() *Node {
	token := p.currentToken()

	if token.Type == TokenParenL {
		p.advance()
		node := p.parseExpression()
		if p.currentToken().Type == TokenParenR {
			p.advance() // skip ')'
		}
		return node
	}

	field := token.Value
	p.advance()

	operator := p.currentToken().Value
	p.advance()

	valueToken := p.currentToken()
	var value interface{}
	if valueToken.Type == TokenValue {
		if numValue, err := strconv.Atoi(valueToken.Value); err == nil {
			value = numValue
		} else {
			value = valueToken.Value
		}
	}
	p.advance()

	return &Node{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

func (p *Parser) parseExpression() *Node {
	node := p.parsePrimary()

	for {
		token := p.currentToken()
		if token.Type == TokenAnd || token.Type == TokenOr {
			operator := strings.ToLower(token.Value)
			p.advance()
			right := p.parsePrimary()
			node = &Node{
				Operator: operator,
				Left:     node,
				Right:    right,
			}
		} else {
			break
		}
	}

	return node
}

func ParseAST(input string) *Node {
	tokens := Tokenizer(input)
	parser := &Parser{tokens: tokens}
	return parser.parseExpression()
}
