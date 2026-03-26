package calculator

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// SymbolicCalcCalculatorFixed 修复后的符号计算器
type SymbolicCalcCalculatorFixed struct {
	*BaseCalculator
}

// NewSymbolicCalcCalculatorFixed 创建新的修复后符号计算器
func NewSymbolicCalcCalculatorFixed() *SymbolicCalcCalculatorFixed {
	return &SymbolicCalcCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"symbolic_calc_fixed",
			"修复后的符号计算器，支持表达式解析、符号求导和表达式化简",
		),
	}
}

// ASTNode 抽象语法树节点
type ASTNode struct {
	Type     string     `json:"type"`     // "number", "variable", "operator", "function"
	Value    string     `json:"value"`    // 节点值
	Children []*ASTNode `json:"children"` // 子节点
}

// Token 词法分析令牌
type Token struct {
	Type  string // "number", "variable", "operator", "function", "paren"
	Value string
	Pos   int
}

// 运算符优先级
var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

// 函数列表
var functions = map[string]bool{
	"sin":  true,
	"cos":  true,
	"tan":  true,
	"log":  true,
	"ln":   true,
	"sqrt": true,
	"exp":  true,
}

// Calculate 执行符号计算
func (c *SymbolicCalcCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	symbolicParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(symbolicParams); err != nil {
		return nil, err
	}

	var result *SymbolicResult
	switch symbolicParams.Operation {
	case "parse":
		result, err = c.parseExpression(symbolicParams)
	case "differentiate":
		result, err = c.differentiateExpression(symbolicParams)
	case "simplify":
		result, err = c.simplifyExpression(symbolicParams)
	case "evaluate":
		result, err = c.evaluateExpression(symbolicParams)
	default:
		return nil, fmt.Errorf("不支持的操作类型: %s", symbolicParams.Operation)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// parseParams 解析参数
func (c *SymbolicCalcCalculatorFixed) parseParams(params interface{}) (*SymbolicParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	operation, ok := paramsMap["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation参数必须为字符串")
	}

	expression, ok := paramsMap["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("expression参数必须为字符串")
	}

	paramsObj := &SymbolicParams{
		Operation:  operation,
		Expression: expression,
		Variable:   "x",
	}

	if variable, ok := paramsMap["variable"].(string); ok {
		paramsObj.Variable = variable
	}

	if xValue, ok := paramsMap["x_value"].(float64); ok {
		paramsObj.XValue = xValue
	}

	if yValue, ok := paramsMap["y_value"].(float64); ok {
		paramsObj.YValue = yValue
	}

	if zValue, ok := paramsMap["z_value"].(float64); ok {
		paramsObj.ZValue = zValue
	}

	return paramsObj, nil
}

// validateParams 验证参数
func (c *SymbolicCalcCalculatorFixed) validateParams(params *SymbolicParams) error {
	if params.Expression == "" {
		return fmt.Errorf("表达式不能为空")
	}

	if !c.isValidExpression(params.Expression) {
		return fmt.Errorf("表达式格式不正确: %s", params.Expression)
	}

	return nil
}

// isValidExpression 验证表达式格式
func (c *SymbolicCalcCalculatorFixed) isValidExpression(expr string) bool {
	validPattern := `^[a-zA-Z0-9\+\-\*\/\(\)\^\s\.]+$`
	matched, _ := regexp.MatchString(validPattern, expr)
	return matched
}

// tokenize 词法分析
func (c *SymbolicCalcCalculatorFixed) tokenize(expr string) []Token {
	var tokens []Token
	var current strings.Builder
	pos := 0

	for i, char := range expr {
		switch {
		case unicode.IsSpace(char):
			if current.Len() > 0 {
				tokens = append(tokens, c.classifyToken(current.String(), pos))
				current.Reset()
			}
		case unicode.IsDigit(char) || char == '.':
			if current.Len() == 0 {
				pos = i
			}
			current.WriteRune(char)
		case unicode.IsLetter(char):
			if current.Len() == 0 {
				pos = i
			}
			current.WriteRune(char)
		case strings.ContainsRune("+-*/^()", char):
			if current.Len() > 0 {
				tokens = append(tokens, c.classifyToken(current.String(), pos))
				current.Reset()
			}
			tokenType := "operator"
			if char == '(' || char == ')' {
				tokenType = "paren"
			}
			tokens = append(tokens, Token{Type: tokenType, Value: string(char), Pos: i})
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, c.classifyToken(current.String(), pos))
	}

	// 处理负号（如 -x 或 3*-2）
	return c.processNegativeSigns(tokens)
}

// classifyToken 分类令牌
func (c *SymbolicCalcCalculatorFixed) classifyToken(s string, pos int) Token {
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return Token{Type: "number", Value: s, Pos: pos}
	}
	if functions[s] {
		return Token{Type: "function", Value: s, Pos: pos}
	}
	if len(s) == 1 && unicode.IsLetter(rune(s[0])) {
		return Token{Type: "variable", Value: s, Pos: pos}
	}
	return Token{Type: "function", Value: s, Pos: pos}
}

// processNegativeSigns 处理负号
func (c *SymbolicCalcCalculatorFixed) processNegativeSigns(tokens []Token) []Token {
	var result []Token
	for i, token := range tokens {
		if token.Value == "-" {
			// 负号出现在开头、或前一个是运算符、或前一个是左括号
			if i == 0 ||
				(tokens[i-1].Type == "operator") ||
				(tokens[i-1].Value == "(") {
				result = append(result, Token{Type: "number", Value: "-1", Pos: token.Pos})
				result = append(result, Token{Type: "operator", Value: "*", Pos: token.Pos})
				continue
			}
		}
		result = append(result, token)
	}
	return result
}

// parseToAST 解析为抽象语法树（Shunting-yard算法）
func (c *SymbolicCalcCalculatorFixed) parseToAST(expr string) (*ASTNode, error) {
	tokens := c.tokenize(expr)
	var output []*ASTNode
	var operators []Token

	for _, token := range tokens {
		switch token.Type {
		case "number", "variable":
			output = append(output, &ASTNode{Type: token.Type, Value: token.Value})
		case "function":
			operators = append(operators, token)
		case "operator":
			for len(operators) > 0 {
				top := operators[len(operators)-1]
				if top.Value == "(" {
					break
				}
				if top.Type == "function" || precedence[top.Value] > precedence[token.Value] ||
					(precedence[top.Value] == precedence[token.Value] && token.Value != "^") {
					output = c.applyOperator(&operators, output)
				} else {
					break
				}
			}
			operators = append(operators, token)
		case "paren":
			if token.Value == "(" {
				operators = append(operators, token)
			} else {
				for len(operators) > 0 && operators[len(operators)-1].Value != "(" {
					output = c.applyOperator(&operators, output)
				}
				if len(operators) == 0 {
					return nil, fmt.Errorf("括号不匹配")
				}
				operators = operators[:len(operators)-1] // 弹出 "("
				// 如果前面是函数，应用函数
				if len(operators) > 0 && operators[len(operators)-1].Type == "function" {
					output = c.applyOperator(&operators, output)
				}
			}
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1].Value == "(" {
			return nil, fmt.Errorf("括号不匹配")
		}
		output = c.applyOperator(&operators, output)
	}

	if len(output) != 1 {
		return nil, fmt.Errorf("表达式解析失败")
	}

	return output[0], nil
}

// applyOperator 应用运算符
func (c *SymbolicCalcCalculatorFixed) applyOperator(operators *[]Token, output []*ASTNode) []*ASTNode {
	if len(*operators) == 0 {
		return output
	}

	op := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	if op.Type == "function" {
		if len(output) < 1 {
			return output
		}
		arg := output[len(output)-1]
		output = output[:len(output)-1]
		node := &ASTNode{
			Type:     "function",
			Value:    op.Value,
			Children: []*ASTNode{arg},
		}
		return append(output, node)
	}

	if len(output) < 2 {
		return output
	}

	right := output[len(output)-1]
	left := output[len(output)-2]
	output = output[:len(output)-2]

	node := &ASTNode{
		Type:     "operator",
		Value:    op.Value,
		Children: []*ASTNode{left, right},
	}

	return append(output, node)
}

// astToString AST转换为字符串
func (c *SymbolicCalcCalculatorFixed) astToString(node *ASTNode) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case "number", "variable":
		return node.Value
	case "function":
		if len(node.Children) > 0 {
			return fmt.Sprintf("%s(%s)", node.Value, c.astToString(node.Children[0]))
		}
		return node.Value + "()"
	case "operator":
		if len(node.Children) == 2 {
			left := c.astToString(node.Children[0])
			right := c.astToString(node.Children[1])

			// 根据优先级决定是否加括号
			if node.Value == "^" || node.Value == "*" || node.Value == "/" {
				if node.Children[0].Type == "operator" &&
					(node.Children[0].Value == "+" || node.Children[0].Value == "-") {
					left = "(" + left + ")"
				}
				if node.Children[1].Type == "operator" &&
					(node.Children[1].Value == "+" || node.Children[1].Value == "-") {
					right = "(" + right + ")"
				}
			}
			if node.Value == "^" {
				if node.Children[1].Type == "operator" {
					right = "(" + right + ")"
				}
			}

			return left + node.Value + right
		}
	}
	return ""
}

// differentiateAST 对AST进行递归求导
func (c *SymbolicCalcCalculatorFixed) differentiateAST(node *ASTNode, variable string) *ASTNode {
	if node == nil {
		return nil
	}

	switch node.Type {
	case "number":
		// d/dx(constant) = 0
		return &ASTNode{Type: "number", Value: "0"}
	case "variable":
		// d/dx(x) = 1, d/dx(y) = 0
		if node.Value == variable {
			return &ASTNode{Type: "number", Value: "1"}
		}
		return &ASTNode{Type: "number", Value: "0"}
	case "operator":
		switch node.Value {
		case "+":
			// d/dx(f + g) = f' + g'
			return &ASTNode{
				Type:  "operator",
				Value: "+",
				Children: []*ASTNode{
					c.differentiateAST(node.Children[0], variable),
					c.differentiateAST(node.Children[1], variable),
				},
			}
		case "-":
			// d/dx(f - g) = f' - g'
			return &ASTNode{
				Type:  "operator",
				Value: "-",
				Children: []*ASTNode{
					c.differentiateAST(node.Children[0], variable),
					c.differentiateAST(node.Children[1], variable),
				},
			}
		case "*":
			// d/dx(f * g) = f'*g + f*g'
			fPrime := c.differentiateAST(node.Children[0], variable)
			gPrime := c.differentiateAST(node.Children[1], variable)
			return &ASTNode{
				Type:  "operator",
				Value: "+",
				Children: []*ASTNode{
					{
						Type:     "operator",
						Value:    "*",
						Children: []*ASTNode{fPrime, node.Children[1]},
					},
					{
						Type:     "operator",
						Value:    "*",
						Children: []*ASTNode{node.Children[0], gPrime},
					},
				},
			}
		case "/":
			// d/dx(f/g) = (f'g - fg')/g²
			fPrime := c.differentiateAST(node.Children[0], variable)
			gPrime := c.differentiateAST(node.Children[1], variable)
			fgPrime := &ASTNode{
				Type:     "operator",
				Value:    "*",
				Children: []*ASTNode{node.Children[0], gPrime},
			}
			fPrimeg := &ASTNode{
				Type:     "operator",
				Value:    "*",
				Children: []*ASTNode{fPrime, node.Children[1]},
			}
			numerator := &ASTNode{
				Type:     "operator",
				Value:    "-",
				Children: []*ASTNode{fPrimeg, fgPrime},
			}
			denominator := &ASTNode{
				Type:     "operator",
				Value:    "^",
				Children: []*ASTNode{node.Children[1], {Type: "number", Value: "2"}},
			}
			return &ASTNode{
				Type:     "operator",
				Value:    "/",
				Children: []*ASTNode{numerator, denominator},
			}
		case "^":
			// d/dx(f^n) = n*f^(n-1)*f'，当n是常数时
			// 或者链式法则: d/dx(u^v) = u^v * (v'*ln(u) + v*u'/u)
			base := node.Children[0]
			exponent := node.Children[1]

			// 检查指数是否为常数
			if exponent.Type == "number" {
				n, _ := strconv.ParseFloat(exponent.Value, 64)
				nMinus1 := strconv.FormatFloat(n-1, 'f', -1, 64)
				basePrime := c.differentiateAST(base, variable)

				// n * f^(n-1)
				term1 := &ASTNode{
					Type:  "operator",
					Value: "*",
					Children: []*ASTNode{
						{Type: "number", Value: exponent.Value},
						{
							Type:     "operator",
							Value:    "^",
							Children: []*ASTNode{base, {Type: "number", Value: nMinus1}},
						},
					},
				}
				// * f'
				return &ASTNode{
					Type:     "operator",
					Value:    "*",
					Children: []*ASTNode{term1, basePrime},
				}
			}

			// 一般情况使用链式法则
			basePrime := c.differentiateAST(base, variable)
			expPrime := c.differentiateAST(exponent, variable)

			lnBase := &ASTNode{
				Type:     "function",
				Value:    "ln",
				Children: []*ASTNode{base},
			}

			term1 := &ASTNode{
				Type:     "operator",
				Value:    "*",
				Children: []*ASTNode{expPrime, lnBase},
			}

			term2 := &ASTNode{
				Type:  "operator",
				Value: "/",
				Children: []*ASTNode{
					{
						Type:     "operator",
						Value:    "*",
						Children: []*ASTNode{exponent, basePrime},
					},
					base,
				},
			}

			sum := &ASTNode{
				Type:     "operator",
				Value:    "+",
				Children: []*ASTNode{term1, term2},
			}

			return &ASTNode{
				Type:     "operator",
				Value:    "*",
				Children: []*ASTNode{node, sum},
			}
		}
	case "function":
		// 函数求导 - 链式法则
		if len(node.Children) > 0 {
			arg := node.Children[0]
			argPrime := c.differentiateAST(arg, variable)

			var derivative *ASTNode
			switch node.Value {
			case "sin":
				derivative = &ASTNode{
					Type:     "function",
					Value:    "cos",
					Children: []*ASTNode{arg},
				}
			case "cos":
				derivative = &ASTNode{
					Type:  "operator",
					Value: "*",
					Children: []*ASTNode{
						{Type: "number", Value: "-1"},
						{
							Type:     "function",
							Value:    "sin",
							Children: []*ASTNode{arg},
						},
					},
				}
			case "tan":
				sec2 := &ASTNode{
					Type:  "operator",
					Value: "^",
					Children: []*ASTNode{
						{
							Type:     "function",
							Value:    "sec",
							Children: []*ASTNode{arg},
						},
						{Type: "number", Value: "2"},
					},
				}
				derivative = sec2
			case "ln":
				derivative = &ASTNode{
					Type:     "operator",
					Value:    "/",
					Children: []*ASTNode{{Type: "number", Value: "1"}, arg},
				}
			case "log":
				derivative = &ASTNode{
					Type:  "operator",
					Value: "/",
					Children: []*ASTNode{
						{Type: "number", Value: "1"},
						{
							Type:     "operator",
							Value:    "*",
							Children: []*ASTNode{arg, {Type: "function", Value: "ln(10)"}},
						},
					},
				}
			case "exp":
				derivative = node
			case "sqrt":
				derivative = &ASTNode{
					Type:  "operator",
					Value: "/",
					Children: []*ASTNode{
						{Type: "number", Value: "1"},
						{
							Type:     "operator",
							Value:    "*",
							Children: []*ASTNode{{Type: "number", Value: "2"}, node},
						},
					},
				}
			default:
				derivative = &ASTNode{Type: "number", Value: "1"}
			}

			return &ASTNode{
				Type:     "operator",
				Value:    "*",
				Children: []*ASTNode{derivative, argPrime},
			}
		}
	}

	return &ASTNode{Type: "number", Value: "0"}
}

// simplifyAST 化简AST
func (c *SymbolicCalcCalculatorFixed) simplifyAST(node *ASTNode) *ASTNode {
	if node == nil {
		return nil
	}

	// 先递归化简子节点
	for i := range node.Children {
		node.Children[i] = c.simplifyAST(node.Children[i])
	}

	switch node.Type {
	case "operator":
		return c.simplifyOperator(node)
	case "function":
		// 函数化简
		if len(node.Children) > 0 {
			arg := node.Children[0]
			if node.Value == "sqrt" && arg.Type == "number" {
				if n, err := strconv.ParseFloat(arg.Value, 64); err == nil {
					result := math.Sqrt(n)
					if result == math.Trunc(result) {
						return &ASTNode{Type: "number", Value: strconv.FormatFloat(result, 'f', -1, 64)}
					}
				}
			}
		}
	}

	return node
}

// simplifyOperator 化简运算符节点
func (c *SymbolicCalcCalculatorFixed) simplifyOperator(node *ASTNode) *ASTNode {
	if len(node.Children) != 2 {
		return node
	}

	left := node.Children[0]
	right := node.Children[1]

	// 常数计算
	if left.Type == "number" && right.Type == "number" {
		lVal, _ := strconv.ParseFloat(left.Value, 64)
		rVal, _ := strconv.ParseFloat(right.Value, 64)

		var result float64
		switch node.Value {
		case "+":
			result = lVal + rVal
		case "-":
			result = lVal - rVal
		case "*":
			result = lVal * rVal
		case "/":
			if rVal != 0 {
				result = lVal / rVal
			} else {
				return node
			}
		case "^":
			result = math.Pow(lVal, rVal)
		}

		return &ASTNode{Type: "number", Value: strconv.FormatFloat(result, 'f', -1, 64)}
	}

	switch node.Value {
	case "+":
		// x + 0 = x, 0 + x = x
		if left.Type == "number" && left.Value == "0" {
			return right
		}
		if right.Type == "number" && right.Value == "0" {
			return left
		}
		// 合并同类项
		return c.mergeLikeTerms(node)
	case "-":
		// x - 0 = x, 0 - x = -x
		if left.Type == "number" && left.Value == "0" {
			return &ASTNode{
				Type:     "operator",
				Value:    "*",
				Children: []*ASTNode{{Type: "number", Value: "-1"}, right},
			}
		}
		if right.Type == "number" && right.Value == "0" {
			return left
		}
	case "*":
		// x * 0 = 0, 0 * x = 0
		if (left.Type == "number" && left.Value == "0") ||
			(right.Type == "number" && right.Value == "0") {
			return &ASTNode{Type: "number", Value: "0"}
		}
		// x * 1 = x, 1 * x = x
		if left.Type == "number" && left.Value == "1" {
			return right
		}
		if right.Type == "number" && right.Value == "1" {
			return left
		}
		// x * (-1) = -x
		if left.Type == "number" && left.Value == "-1" {
			if right.Type == "operator" && right.Value == "*" && right.Children[0].Type == "number" {
				coeff, _ := strconv.ParseFloat(right.Children[0].Value, 64)
				coeff = -coeff
				return &ASTNode{
					Type:  "operator",
					Value: "*",
					Children: []*ASTNode{
						{Type: "number", Value: strconv.FormatFloat(coeff, 'f', -1, 64)},
						right.Children[1],
					},
				}
			}
		}
	case "^":
		// x^0 = 1
		if right.Type == "number" && right.Value == "0" {
			return &ASTNode{Type: "number", Value: "1"}
		}
		// x^1 = x
		if right.Type == "number" && right.Value == "1" {
			return left
		}
		// 0^x = 0 (x>0)
		if left.Type == "number" && left.Value == "0" {
			return &ASTNode{Type: "number", Value: "0"}
		}
		// 1^x = 1
		if left.Type == "number" && left.Value == "1" {
			return &ASTNode{Type: "number", Value: "1"}
		}
	case "/":
		// 0 / x = 0
		if left.Type == "number" && left.Value == "0" {
			return &ASTNode{Type: "number", Value: "0"}
		}
		// x / 1 = x
		if right.Type == "number" && right.Value == "1" {
			return left
		}
	}

	return node
}

// mergeLikeTerms 合并同类项
func (c *SymbolicCalcCalculatorFixed) mergeLikeTerms(node *ASTNode) *ASTNode {
	if node.Value != "+" && node.Value != "-" {
		return node
	}

	// 收集所有项
	terms := c.collectTerms(node)

	// 按变量部分分组
	termGroups := make(map[string][]*ASTNode)
	constTerms := 0.0

	for _, term := range terms {
		coeff, vars := c.splitTerm(term)
		if vars == "" {
			constTerms += coeff
		} else {
			termGroups[vars] = append(termGroups[vars], term)
		}
	}

	// 合并每组
	var resultTerms []*ASTNode

	// 添加常数项
	if constTerms != 0 || len(termGroups) == 0 {
		resultTerms = append(resultTerms, &ASTNode{
			Type:  "number",
			Value: strconv.FormatFloat(constTerms, 'f', -1, 64),
		})
	}

	// 添加合并后的变量项
	for _, group := range termGroups {
		merged := c.mergeTermGroup(group)
		if merged != nil {
			resultTerms = append(resultTerms, merged)
		}
	}

	// 构建结果树
	if len(resultTerms) == 0 {
		return &ASTNode{Type: "number", Value: "0"}
	}
	if len(resultTerms) == 1 {
		return resultTerms[0]
	}

	result := resultTerms[0]
	for i := 1; i < len(resultTerms); i++ {
		result = &ASTNode{
			Type:     "operator",
			Value:    "+",
			Children: []*ASTNode{result, resultTerms[i]},
		}
	}

	return result
}

// collectTerms 收集所有加法/减法项
func (c *SymbolicCalcCalculatorFixed) collectTerms(node *ASTNode) []*ASTNode {
	var terms []*ASTNode

	if node.Type != "operator" || (node.Value != "+" && node.Value != "-") {
		return []*ASTNode{node}
	}

	if node.Value == "-" && len(node.Children) == 2 {
		// 转换为 a + (-1)*b
		negB := &ASTNode{
			Type:     "operator",
			Value:    "*",
			Children: []*ASTNode{{Type: "number", Value: "-1"}, node.Children[1]},
		}
		terms = append(terms, c.collectTerms(node.Children[0])...)
		terms = append(terms, c.collectTerms(negB)...)
		return terms
	}

	for _, child := range node.Children {
		terms = append(terms, c.collectTerms(child)...)
	}

	return terms
}

// splitTerm 拆分项为系数和变量部分
func (c *SymbolicCalcCalculatorFixed) splitTerm(node *ASTNode) (float64, string) {
	if node.Type == "number" {
		val, _ := strconv.ParseFloat(node.Value, 64)
		return val, ""
	}
	if node.Type == "variable" {
		return 1.0, node.Value
	}
	if node.Type == "operator" && node.Value == "*" {
		// 系数 * 变量
		if node.Children[0].Type == "number" {
			coeff, _ := strconv.ParseFloat(node.Children[0].Value, 64)
			return coeff, c.astToString(node.Children[1])
		}
		if node.Children[1].Type == "number" {
			coeff, _ := strconv.ParseFloat(node.Children[1].Value, 64)
			return coeff, c.astToString(node.Children[0])
		}
	}
	if node.Type == "operator" && node.Value == "^" {
		return 1.0, c.astToString(node)
	}
	return 1.0, c.astToString(node)
}

// mergeTermGroup 合并同类型组
func (c *SymbolicCalcCalculatorFixed) mergeTermGroup(terms []*ASTNode) *ASTNode {
	if len(terms) == 0 {
		return nil
	}
	if len(terms) == 1 {
		return terms[0]
	}

	totalCoeff := 0.0
	var varPart string

	for _, term := range terms {
		coeff, vars := c.splitTerm(term)
		totalCoeff += coeff
		varPart = vars
	}

	if totalCoeff == 0 {
		return nil
	}
	if totalCoeff == 1 {
		ast, _ := c.parseToAST(varPart)
		return ast
	}

	coeffNode := &ASTNode{
		Type:  "number",
		Value: strconv.FormatFloat(totalCoeff, 'f', -1, 64),
	}

	if varPart == "" {
		return coeffNode
	}

	varNode, _ := c.parseToAST(varPart)
	return &ASTNode{
		Type:     "operator",
		Value:    "*",
		Children: []*ASTNode{coeffNode, varNode},
	}
}

// evaluateAST 计算AST的数值
func (c *SymbolicCalcCalculatorFixed) evaluateAST(node *ASTNode, variables map[string]float64) (float64, error) {
	if node == nil {
		return 0, fmt.Errorf("空节点")
	}

	switch node.Type {
	case "number":
		return strconv.ParseFloat(node.Value, 64)
	case "variable":
		if val, ok := variables[node.Value]; ok {
			return val, nil
		}
		return 0, fmt.Errorf("未定义变量: %s", node.Value)
	case "operator":
		if len(node.Children) != 2 {
			return 0, fmt.Errorf("运算符需要两个操作数")
		}
		left, err := c.evaluateAST(node.Children[0], variables)
		if err != nil {
			return 0, err
		}
		right, err := c.evaluateAST(node.Children[1], variables)
		if err != nil {
			return 0, err
		}

		switch node.Value {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("除零错误")
			}
			return left / right, nil
		case "^":
			return math.Pow(left, right), nil
		}
	case "function":
		if len(node.Children) != 1 {
			return 0, fmt.Errorf("函数需要一个参数")
		}
		arg, err := c.evaluateAST(node.Children[0], variables)
		if err != nil {
			return 0, err
		}

		switch node.Value {
		case "sin":
			return math.Sin(arg), nil
		case "cos":
			return math.Cos(arg), nil
		case "tan":
			return math.Tan(arg), nil
		case "ln":
			if arg <= 0 {
				return 0, fmt.Errorf("对数参数必须为正")
			}
			return math.Log(arg), nil
		case "log":
			if arg <= 0 {
				return 0, fmt.Errorf("对数参数必须为正")
			}
			return math.Log10(arg), nil
		case "sqrt":
			if arg < 0 {
				return 0, fmt.Errorf("平方根参数不能为负")
			}
			return math.Sqrt(arg), nil
		case "exp":
			return math.Exp(arg), nil
		}
	}

	return 0, fmt.Errorf("未知节点类型: %s", node.Type)
}

// parseExpression 解析数学表达式
func (c *SymbolicCalcCalculatorFixed) parseExpression(params *SymbolicParams) (*SymbolicResult, error) {
	ast, err := c.parseToAST(params.Expression)
	if err != nil {
		return nil, err
	}

	result := &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:   c.astToString(ast),
		ParsedTree:         ast,
		OperationType:      "parse",
	}

	c.populateNumericFields(result, params)
	return result, nil
}

// differentiateExpression 符号求导
func (c *SymbolicCalcCalculatorFixed) differentiateExpression(params *SymbolicParams) (*SymbolicResult, error) {
	// 1. 解析表达式为AST
	ast, err := c.parseToAST(params.Expression)
	if err != nil {
		return nil, err
	}

	// 2. 递归求导
	derivativeAST := c.differentiateAST(ast, params.Variable)

	// 3. 化简导数
	simplifiedAST := c.simplifyAST(derivativeAST)
	simplifiedAST = c.simplifyAST(simplifiedAST) // 二次化简

	// 4. 转换为字符串
	derivativeStr := c.astToString(derivativeAST)
	simplifiedStr := c.astToString(simplifiedAST)

	// 5. 数值计算（默认x=1.0）
	variables := map[string]float64{
		"x": params.XValue,
		"y": params.YValue,
		"z": params.ZValue,
	}
	if params.XValue == 0 && params.YValue == 0 && params.ZValue == 0 {
		variables["x"] = 1.0 // 默认x=1
	}

	numericValue, _ := c.evaluateAST(simplifiedAST, variables)

	result := &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:   simplifiedStr,
		ParsedTree:         ast,
		Derivative:         derivativeStr,
		Simplified:         simplifiedStr,
		NumericValue:       numericValue,
		Variables:          variables,
		OperationType:      "differentiate",
	}

	c.populateNumericFields(result, params)
	return result, nil
}

// simplifyExpression 表达式化简
func (c *SymbolicCalcCalculatorFixed) simplifyExpression(params *SymbolicParams) (*SymbolicResult, error) {
	ast, err := c.parseToAST(params.Expression)
	if err != nil {
		return nil, err
	}

	simplifiedAST := c.simplifyAST(ast)
	simplifiedAST = c.simplifyAST(simplifiedAST)
	simplifiedStr := c.astToString(simplifiedAST)

	result := &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:   simplifiedStr,
		ParsedTree:         ast,
		Simplified:         simplifiedStr,
		OperationType:      "simplify",
	}

	c.populateNumericFields(result, params)
	return result, nil
}

// evaluateExpression 表达式求值
func (c *SymbolicCalcCalculatorFixed) evaluateExpression(params *SymbolicParams) (*SymbolicResult, error) {
	ast, err := c.parseToAST(params.Expression)
	if err != nil {
		return nil, err
	}

	variables := map[string]float64{
		"x": params.XValue,
		"y": params.YValue,
		"z": params.ZValue,
	}

	numericValue, err := c.evaluateAST(ast, variables)
	if err != nil {
		return nil, err
	}

	result := &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:   fmt.Sprintf("%.6f", numericValue),
		ParsedTree:         ast,
		NumericValue:       numericValue,
		Variables:          variables,
		OperationType:      "evaluate",
	}

	c.populateNumericFields(result, params)
	return result, nil
}

// populateNumericFields 填充扩展数值字段
func (c *SymbolicCalcCalculatorFixed) populateNumericFields(result *SymbolicResult, params *SymbolicParams) {
	if result == nil {
		return
	}

	complexity := c.calculateExpressionComplexity(result.OriginalExpression)
	result.ExpressionComplexity = complexity

	varCount := c.countVariables(result.OriginalExpression)
	result.VariableCount = varCount

	termCount := c.countTerms(result.OriginalExpression)
	result.TermCount = float64(termCount)

	treeDepth := c.estimateTreeDepth(result.OriginalExpression)
	result.TreeDepth = float64(treeDepth)

	evalScore := c.calculateEvaluationScore(result, params)
	result.EvaluationScore = evalScore
}

func (c *SymbolicCalcCalculatorFixed) calculateExpressionComplexity(expr string) float64 {
	if expr == "" {
		return 0.0
	}

	baseScore := float64(len(expr)) * 0.5

	operatorCount := strings.Count(expr, "+") + strings.Count(expr, "-") +
		strings.Count(expr, "*") + strings.Count(expr, "/") + strings.Count(expr, "^")
	operatorScore := float64(operatorCount) * 2.0

	nestingLevel := 0
	maxNesting := 0
	for _, char := range expr {
		if char == '(' {
			nestingLevel++
			if nestingLevel > maxNesting {
				maxNesting = nestingLevel
			}
		} else if char == ')' {
			nestingLevel--
		}
	}
	nestingScore := float64(maxNesting) * 3.0

	return baseScore + operatorScore + nestingScore
}

func (c *SymbolicCalcCalculatorFixed) countVariables(expr string) int {
	variables := make(map[rune]bool)
	for _, char := range expr {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			variables[char] = true
		}
	}
	return len(variables)
}

func (c *SymbolicCalcCalculatorFixed) countTerms(expr string) int {
	if expr == "" {
		return 0
	}
	termCount := 1
	for _, char := range expr {
		if char == '+' || char == '-' {
			termCount++
		}
	}
	return termCount
}

func (c *SymbolicCalcCalculatorFixed) estimateTreeDepth(expr string) int {
	if expr == "" {
		return 0
	}

	currentDepth := 1
	maxDepth := 1
	for _, char := range expr {
		switch char {
		case '(':
			currentDepth += 2
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
		case ')':
			currentDepth -= 2
		case '*', '/', '^':
			tempDepth := currentDepth + 1
			if tempDepth > maxDepth {
				maxDepth = tempDepth
			}
		}
	}
	return maxDepth
}

func (c *SymbolicCalcCalculatorFixed) calculateEvaluationScore(result *SymbolicResult, params *SymbolicParams) float64 {
	if result == nil {
		return 0.0
	}

	baseScore := 50.0

	switch result.OperationType {
	case "parse":
		baseScore += 10.0
	case "differentiate":
		baseScore += 25.0
	case "simplify":
		baseScore += 15.0
	case "evaluate":
		baseScore += 20.0
	}

	baseScore += float64(result.VariableCount) * 5.0
	baseScore += result.ExpressionComplexity * 0.1

	if result.NumericValue != 0 {
		baseScore += math.Abs(result.NumericValue) * 0.01
	}

	return math.Max(0, math.Min(100, baseScore))
}

// Validate 验证输入参数
func (c *SymbolicCalcCalculatorFixed) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// Description 返回计算器描述
func (c *SymbolicCalcCalculatorFixed) Description() string {
	return "修复后的符号计算器，支持表达式解析、符号求导和表达式化简"
}
