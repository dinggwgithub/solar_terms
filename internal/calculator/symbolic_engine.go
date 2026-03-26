package calculator

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type ExprNode interface {
	String() string
	Equals(other ExprNode) bool
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) String() string {
	if n.Value == float64(int(n.Value)) {
		return fmt.Sprintf("%d", int(n.Value))
	}
	return fmt.Sprintf("%.6g", n.Value)
}

func (n *NumberNode) Equals(other ExprNode) bool {
	if o, ok := other.(*NumberNode); ok {
		return n.Value == o.Value
	}
	return false
}

type VariableNode struct {
	Name string
}

func (v *VariableNode) String() string {
	return v.Name
}

func (v *VariableNode) Equals(other ExprNode) bool {
	if o, ok := other.(*VariableNode); ok {
		return v.Name == o.Name
	}
	return false
}

type BinaryOpNode struct {
	Op    string
	Left  ExprNode
	Right ExprNode
}

func (b *BinaryOpNode) String() string {
	leftStr := b.Left.String()
	rightStr := b.Right.String()

	if b.Op == "^" {
		if _, ok := b.Left.(*BinaryOpNode); ok {
			leftStr = "(" + leftStr + ")"
		}
		if _, ok := b.Right.(*BinaryOpNode); ok {
			rightStr = "(" + rightStr + ")"
		}
		return leftStr + "^" + rightStr
	}

	if b.Op == "*" || b.Op == "/" {
		if leftOp, ok := b.Left.(*BinaryOpNode); ok && (leftOp.Op == "+" || leftOp.Op == "-") {
			leftStr = "(" + leftStr + ")"
		}
		if rightOp, ok := b.Right.(*BinaryOpNode); ok && (rightOp.Op == "+" || rightOp.Op == "-") {
			rightStr = "(" + rightStr + ")"
		}
		if b.Op == "/" {
			if _, ok := b.Right.(*BinaryOpNode); ok {
				rightStr = "(" + rightStr + ")"
			}
		}
	}

	return leftStr + b.Op + rightStr
}

func (b *BinaryOpNode) Equals(other ExprNode) bool {
	if o, ok := other.(*BinaryOpNode); ok {
		return b.Op == o.Op && b.Left.Equals(o.Left) && b.Right.Equals(o.Right)
	}
	return false
}

type UnaryOpNode struct {
	Op    string
	Operand ExprNode
}

func (u *UnaryOpNode) String() string {
	if u.Op == "-" {
		return "-" + u.Operand.String()
	}
	return u.Op + "(" + u.Operand.String() + ")"
}

func (u *UnaryOpNode) Equals(other ExprNode) bool {
	if o, ok := other.(*UnaryOpNode); ok {
		return u.Op == o.Op && u.Operand.Equals(o.Operand)
	}
	return false
}

type FunctionNode struct {
	Name string
	Arg  ExprNode
}

func (f *FunctionNode) String() string {
	return f.Name + "(" + f.Arg.String() + ")"
}

func (f *FunctionNode) Equals(other ExprNode) bool {
	if o, ok := other.(*FunctionNode); ok {
		return f.Name == o.Name && f.Arg.Equals(o.Arg)
	}
	return false
}

type Parser struct {
	tokens []string
	pos    int
}

func NewParser(expr string) *Parser {
	tokens := tokenize(expr)
	return &Parser{tokens: tokens, pos: 0}
}

func tokenize(expr string) []string {
	tokens := []string{}
	i := 0
	expr = strings.ReplaceAll(expr, " ", "")

	for i < len(expr) {
		char := expr[i]

		if char >= '0' && char <= '9' || char == '.' {
			j := i
			for j < len(expr) && (expr[j] >= '0' && expr[j] <= '9' || expr[j] == '.') {
				j++
			}
			tokens = append(tokens, expr[i:j])
			i = j
			continue
		}

		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			j := i
			for j < len(expr) && ((expr[j] >= 'a' && expr[j] <= 'z') || (expr[j] >= 'A' && expr[j] <= 'Z') || (expr[j] >= '0' && expr[j] <= '9')) {
				j++
			}
			tokens = append(tokens, expr[i:j])
			i = j
			continue
		}

		if char == '+' || char == '-' || char == '*' || char == '/' || char == '^' || char == '(' || char == ')' {
			tokens = append(tokens, string(char))
			i++
			continue
		}

		i++
	}

	return tokens
}

func (p *Parser) Parse() (ExprNode, error) {
	return p.parseExpression()
}

func (p *Parser) parseExpression() (ExprNode, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "+" || p.tokens[p.pos] == "-") {
		op := p.tokens[p.pos]
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Op: op, Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parseTerm() (ExprNode, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "*" || p.tokens[p.pos] == "/") {
		op := p.tokens[p.pos]
		p.pos++
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Op: op, Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parsePower() (ExprNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && p.tokens[p.pos] == "^" {
		p.pos++
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Op: "^", Left: left, Right: right}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ExprNode, error) {
	if p.pos < len(p.tokens) && p.tokens[p.pos] == "-" {
		p.pos++
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryOpNode{Op: "-", Operand: operand}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ExprNode, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	token := p.tokens[p.pos]

	if token == "(" {
		p.pos++
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return expr, nil
	}

	if num, err := strconv.ParseFloat(token, 64); err == nil {
		p.pos++
		return &NumberNode{Value: num}, nil
	}

	if isFunction(token) {
		funcName := token
		p.pos++
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != "(" {
			return nil, fmt.Errorf("expected '(' after function %s", funcName)
		}
		p.pos++
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("missing closing parenthesis after function argument")
		}
		p.pos++
		return &FunctionNode{Name: funcName, Arg: arg}, nil
	}

	if isVariable(token) {
		p.pos++
		return &VariableNode{Name: token}, nil
	}

	return nil, fmt.Errorf("unexpected token: %s", token)
}

func isFunction(token string) bool {
	functions := []string{"sin", "cos", "tan", "exp", "log", "ln", "sqrt", "abs"}
	for _, f := range functions {
		if token == f {
			return true
		}
	}
	return false
}

func isVariable(token string) bool {
	if len(token) == 0 {
		return false
	}
	if (token[0] >= 'a' && token[0] <= 'z') || (token[0] >= 'A' && token[0] <= 'Z') {
		return !isFunction(token)
	}
	return false
}

func Differentiate(node ExprNode, variable string) ExprNode {
	switch n := node.(type) {
	case *NumberNode:
		return &NumberNode{Value: 0}

	case *VariableNode:
		if n.Name == variable {
			return &NumberNode{Value: 1}
		}
		return &NumberNode{Value: 0}

	case *BinaryOpNode:
		switch n.Op {
		case "+":
			return &BinaryOpNode{
				Op:    "+",
				Left:  Differentiate(n.Left, variable),
				Right: Differentiate(n.Right, variable),
			}
		case "-":
			return &BinaryOpNode{
				Op:    "-",
				Left:  Differentiate(n.Left, variable),
				Right: Differentiate(n.Right, variable),
			}
		case "*":
			return &BinaryOpNode{
				Op: "+",
				Left: &BinaryOpNode{
					Op:    "*",
					Left:  Differentiate(n.Left, variable),
					Right: n.Right,
				},
				Right: &BinaryOpNode{
					Op:    "*",
					Left:  n.Left,
					Right: Differentiate(n.Right, variable),
				},
			}
		case "/":
			return &BinaryOpNode{
				Op: "/",
				Left: &BinaryOpNode{
					Op: "-",
					Left: &BinaryOpNode{
						Op:    "*",
						Left:  Differentiate(n.Left, variable),
						Right: n.Right,
					},
					Right: &BinaryOpNode{
						Op:    "*",
						Left:  n.Left,
						Right: Differentiate(n.Right, variable),
					},
				},
				Right: &BinaryOpNode{
					Op:    "^",
					Left:  n.Right,
					Right: &NumberNode{Value: 2},
				},
			}
		case "^":
			if isConstant(n.Right, variable) {
				return &BinaryOpNode{
					Op: "*",
					Left: &BinaryOpNode{
						Op:    "*",
						Left:  n.Right,
						Right: Differentiate(n.Left, variable),
					},
					Right: &BinaryOpNode{
						Op:    "^",
						Left:  n.Left,
						Right: &BinaryOpNode{
							Op:    "-",
							Left:  n.Right,
							Right: &NumberNode{Value: 1},
						},
					},
				}
			}
			return &BinaryOpNode{
				Op: "*",
				Left: node,
				Right: &BinaryOpNode{
					Op: "+",
					Left: &BinaryOpNode{
						Op: "*",
						Left: &BinaryOpNode{
							Op:    "*",
							Left:  n.Right,
							Right: &FunctionNode{Name: "ln", Arg: n.Left},
						},
						Right: Differentiate(n.Right, variable),
					},
					Right: &BinaryOpNode{
						Op: "*",
						Left: &BinaryOpNode{
							Op:    "/",
							Left:  n.Right,
							Right: n.Left,
						},
						Right: Differentiate(n.Left, variable),
					},
				},
			}
		}

	case *UnaryOpNode:
		if n.Op == "-" {
			return &UnaryOpNode{
				Op:      "-",
				Operand: Differentiate(n.Operand, variable),
			}
		}

	case *FunctionNode:
		switch n.Name {
		case "sin":
			return &BinaryOpNode{
				Op: "*",
				Left: &FunctionNode{
					Name: "cos",
					Arg:  n.Arg,
				},
				Right: Differentiate(n.Arg, variable),
			}
		case "cos":
			return &UnaryOpNode{
				Op: "-",
				Operand: &BinaryOpNode{
					Op: "*",
					Left: &FunctionNode{
						Name: "sin",
						Arg:  n.Arg,
					},
					Right: Differentiate(n.Arg, variable),
				},
			}
		case "tan":
			return &BinaryOpNode{
				Op: "/",
				Left: Differentiate(n.Arg, variable),
				Right: &BinaryOpNode{
					Op:    "^",
					Left:  &FunctionNode{Name: "cos", Arg: n.Arg},
					Right: &NumberNode{Value: 2},
				},
			}
		case "exp":
			return &BinaryOpNode{
				Op: "*",
				Left: &FunctionNode{
					Name: "exp",
					Arg:  n.Arg,
				},
				Right: Differentiate(n.Arg, variable),
			}
		case "log", "ln":
			return &BinaryOpNode{
				Op: "/",
				Left:  Differentiate(n.Arg, variable),
				Right: n.Arg,
			}
		case "sqrt":
			return &BinaryOpNode{
				Op: "/",
				Left: Differentiate(n.Arg, variable),
				Right: &BinaryOpNode{
					Op: "*",
					Left: &NumberNode{Value: 2},
					Right: &FunctionNode{
						Name: "sqrt",
						Arg:  n.Arg,
					},
				},
			}
		}
	}

	return &NumberNode{Value: 0}
}

func isConstant(node ExprNode, variable string) bool {
	switch n := node.(type) {
	case *NumberNode:
		return true
	case *VariableNode:
		return n.Name != variable
	case *BinaryOpNode:
		return isConstant(n.Left, variable) && isConstant(n.Right, variable)
	case *UnaryOpNode:
		return isConstant(n.Operand, variable)
	case *FunctionNode:
		return isConstant(n.Arg, variable)
	}
	return false
}

func Simplify(node ExprNode) ExprNode {
	node = simplifyRecursive(node)
	return node
}

func simplifyRecursive(node ExprNode) ExprNode {
	switch n := node.(type) {
	case *BinaryOpNode:
		n.Left = simplifyRecursive(n.Left)
		n.Right = simplifyRecursive(n.Right)

		if leftNum, ok := n.Left.(*NumberNode); ok {
			if rightNum, ok := n.Right.(*NumberNode); ok {
				switch n.Op {
				case "+":
					return &NumberNode{Value: leftNum.Value + rightNum.Value}
				case "-":
					return &NumberNode{Value: leftNum.Value - rightNum.Value}
				case "*":
					return &NumberNode{Value: leftNum.Value * rightNum.Value}
				case "/":
					if rightNum.Value != 0 {
						return &NumberNode{Value: leftNum.Value / rightNum.Value}
					}
				case "^":
					return &NumberNode{Value: math.Pow(leftNum.Value, rightNum.Value)}
				}
			}

			switch n.Op {
			case "+":
				if leftNum.Value == 0 {
					return n.Right
				}
			case "-":
				if leftNum.Value == 0 {
					return &UnaryOpNode{Op: "-", Operand: n.Right}
				}
			case "*":
				if leftNum.Value == 0 {
					return &NumberNode{Value: 0}
				}
				if leftNum.Value == 1 {
					return n.Right
				}
			case "/":
				if leftNum.Value == 0 {
					return &NumberNode{Value: 0}
				}
			case "^":
				if leftNum.Value == 0 {
					return &NumberNode{Value: 0}
				}
				if leftNum.Value == 1 {
					return &NumberNode{Value: 1}
				}
			}
		}

		if rightNum, ok := n.Right.(*NumberNode); ok {
			switch n.Op {
			case "+":
				if rightNum.Value == 0 {
					return n.Left
				}
			case "-":
				if rightNum.Value == 0 {
					return n.Left
				}
			case "*":
				if rightNum.Value == 0 {
					return &NumberNode{Value: 0}
				}
				if rightNum.Value == 1 {
					return n.Left
				}
			case "/":
				if rightNum.Value == 1 {
					return n.Left
				}
			case "^":
				if rightNum.Value == 0 {
					return &NumberNode{Value: 1}
				}
				if rightNum.Value == 1 {
					return n.Left
				}
			}
		}

		if n.Op == "+" {
			if leftUnary, ok := n.Left.(*UnaryOpNode); ok && leftUnary.Op == "-" {
				if leftNum, ok := leftUnary.Operand.(*NumberNode); ok {
					return &BinaryOpNode{
						Op:    "-",
						Left:  n.Right,
						Right: leftNum,
					}
				}
			}
		}

		if n.Op == "*" {
			if leftBin, ok := n.Left.(*BinaryOpNode); ok && leftBin.Op == "+" {
				return &BinaryOpNode{
					Op: "+",
					Left: &BinaryOpNode{
						Op:    "*",
						Left:  leftBin.Left,
						Right: n.Right,
					},
					Right: &BinaryOpNode{
						Op:    "*",
						Left:  leftBin.Right,
						Right: n.Right,
					},
				}
			}
		}

		return n

	case *UnaryOpNode:
		n.Operand = simplifyRecursive(n.Operand)
		if num, ok := n.Operand.(*NumberNode); ok {
			if n.Op == "-" {
				return &NumberNode{Value: -num.Value}
			}
		}
		return n

	case *FunctionNode:
		n.Arg = simplifyRecursive(n.Arg)
		if num, ok := n.Arg.(*NumberNode); ok {
			switch n.Name {
			case "sin":
				return &NumberNode{Value: math.Sin(num.Value)}
			case "cos":
				return &NumberNode{Value: math.Cos(num.Value)}
			case "tan":
				return &NumberNode{Value: math.Tan(num.Value)}
			case "exp":
				return &NumberNode{Value: math.Exp(num.Value)}
			case "log", "ln":
				if num.Value > 0 {
					return &NumberNode{Value: math.Log(num.Value)}
				}
			case "sqrt":
				if num.Value >= 0 {
					return &NumberNode{Value: math.Sqrt(num.Value)}
				}
			case "abs":
				return &NumberNode{Value: math.Abs(num.Value)}
			}
		}
		return n
	}

	return node
}

func Evaluate(node ExprNode, variables map[string]float64) (float64, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, nil

	case *VariableNode:
		if val, ok := variables[n.Name]; ok {
			return val, nil
		}
		return 0, fmt.Errorf("undefined variable: %s", n.Name)

	case *BinaryOpNode:
		left, err := Evaluate(n.Left, variables)
		if err != nil {
			return 0, err
		}
		right, err := Evaluate(n.Right, variables)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		case "^":
			return math.Pow(left, right), nil
		}
		return 0, fmt.Errorf("unknown operator: %s", n.Op)

	case *UnaryOpNode:
		val, err := Evaluate(n.Operand, variables)
		if err != nil {
			return 0, err
		}
		if n.Op == "-" {
			return -val, nil
		}
		return 0, fmt.Errorf("unknown unary operator: %s", n.Op)

	case *FunctionNode:
		arg, err := Evaluate(n.Arg, variables)
		if err != nil {
			return 0, err
		}
		switch n.Name {
		case "sin":
			return math.Sin(arg), nil
		case "cos":
			return math.Cos(arg), nil
		case "tan":
			return math.Tan(arg), nil
		case "exp":
			return math.Exp(arg), nil
		case "log", "ln":
			return math.Log(arg), nil
		case "sqrt":
			return math.Sqrt(arg), nil
		case "abs":
			return math.Abs(arg), nil
		}
		return 0, fmt.Errorf("unknown function: %s", n.Name)
	}

	return 0, fmt.Errorf("unknown node type")
}

func ExtractVariables(node ExprNode) []string {
	vars := make(map[string]bool)
	extractVarsRecursive(node, vars)
	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}
	return result
}

func extractVarsRecursive(node ExprNode, vars map[string]bool) {
	switch n := node.(type) {
	case *VariableNode:
		vars[n.Name] = true
	case *BinaryOpNode:
		extractVarsRecursive(n.Left, vars)
		extractVarsRecursive(n.Right, vars)
	case *UnaryOpNode:
		extractVarsRecursive(n.Operand, vars)
	case *FunctionNode:
		extractVarsRecursive(n.Arg, vars)
	}
}

func NodeToMap(node ExprNode) map[string]interface{} {
	switch n := node.(type) {
	case *NumberNode:
		return map[string]interface{}{
			"type":  "number",
			"value": n.Value,
		}
	case *VariableNode:
		return map[string]interface{}{
			"type": "variable",
			"name": n.Name,
		}
	case *BinaryOpNode:
		return map[string]interface{}{
			"type":  "binary_op",
			"op":    n.Op,
			"left":  NodeToMap(n.Left),
			"right": NodeToMap(n.Right),
		}
	case *UnaryOpNode:
		return map[string]interface{}{
			"type":    "unary_op",
			"op":      n.Op,
			"operand": NodeToMap(n.Operand),
		}
	case *FunctionNode:
		return map[string]interface{}{
			"type": "function",
			"name": n.Name,
			"arg":  NodeToMap(n.Arg),
		}
	}
	return map[string]interface{}{"type": "unknown"}
}

func ParseExpression(expr string) (ExprNode, error) {
	parser := NewParser(expr)
	return parser.Parse()
}

func FormatExpression(node ExprNode) string {
	return node.String()
}

func CleanExpression(expr string) string {
	re := regexp.MustCompile(`\s+`)
	expr = re.ReplaceAllString(expr, "")
	return expr
}
