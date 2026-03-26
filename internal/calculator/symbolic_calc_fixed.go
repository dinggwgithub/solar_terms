package calculator

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// SymbolicCalcCalculatorFixed 修复版符号计算器
type SymbolicCalcCalculatorFixed struct {
	*BaseCalculator
}

// NewSymbolicCalcCalculatorFixed 创建新的修复版符号计算器
func NewSymbolicCalcCalculatorFixed() *SymbolicCalcCalculatorFixed {
	return &SymbolicCalcCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"symbolic_calc_fixed",
			"修复版符号计算器，支持正确的表达式解析、符号求导和表达式化简",
		),
	}
}

// SymbolicParamsFixed 符号计算参数（与原版兼容）
type SymbolicParamsFixed struct {
	Operation  string  `json:"operation"`  // 操作类型：parse, differentiate, simplify, evaluate
	Expression string  `json:"expression"` // 数学表达式
	Variable   string  `json:"variable"`   // 变量名（求导用）
	XValue     float64 `json:"x_value"`    // x的值（求值用）
	YValue     float64 `json:"y_value"`    // y的值（求值用）
	ZValue     float64 `json:"z_value"`    // z的值（求值用）
}

// SymbolicResultFixed 符号计算结果（与原版格式一致）
type SymbolicResultFixed struct {
	OriginalExpression string             `json:"original_expression"` // 原始表达式
	ResultExpression   string             `json:"result_expression"`   // 结果表达式
	ParsedTree         interface{}        `json:"parsed_tree"`         // 解析树
	Derivative         string             `json:"derivative"`          // 导数结果
	Simplified         string             `json:"simplified"`          // 化简结果
	NumericValue       float64            `json:"numeric_value"`       // 数值结果
	Variables          map[string]float64 `json:"variables"`           // 变量值
	OperationType      string             `json:"operation_type"`      // 操作类型
	// 扩展数据字段
	ExpressionComplexity float64 `json:"expression_complexity,omitempty"` // 表达式复杂度评分
	VariableCount        int     `json:"variable_count,omitempty"`        // 变量数量
	TermCount            float64 `json:"term_count,omitempty"`            // 项数统计
	TreeDepth            float64 `json:"tree_depth,omitempty"`            // 解析树深度
	EvaluationScore      float64 `json:"evaluation_score,omitempty"`      // 综合评估得分
}

// Calculate 执行符号计算
func (c *SymbolicCalcCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	symbolicParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(symbolicParams); err != nil {
		return nil, err
	}

	// 根据操作类型执行计算
	var result *SymbolicResultFixed
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
func (c *SymbolicCalcCalculatorFixed) parseParams(params interface{}) (*SymbolicParamsFixed, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 提取必需参数
	operation, ok := paramsMap["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation参数必须为字符串")
	}

	expression, ok := paramsMap["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("expression参数必须为字符串")
	}

	// 设置默认值
	paramsObj := &SymbolicParamsFixed{
		Operation:  operation,
		Expression: expression,
		Variable:   "x",
	}

	// 提取可选参数
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
func (c *SymbolicCalcCalculatorFixed) validateParams(params *SymbolicParamsFixed) error {
	if params.Expression == "" {
		return fmt.Errorf("表达式不能为空")
	}

	// 验证表达式合法性
	if !c.isValidExpression(params.Expression) {
		return fmt.Errorf("表达式格式不正确: %s", params.Expression)
	}

	return nil
}

// isValidExpression 验证表达式格式
func (c *SymbolicCalcCalculatorFixed) isValidExpression(expr string) bool {
	// 简单的表达式格式验证
	// 允许字母、数字、运算符、括号等
	validPattern := `^[a-zA-Z0-9\+\-\*\/\(\)\^\s\.]+$`
	matched, _ := regexp.MatchString(validPattern, expr)
	return matched
}

// parseExpression 解析数学表达式
func (c *SymbolicCalcCalculatorFixed) parseExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	// 构建解析树
	parsedTree := c.buildParseTree(params.Expression)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   params.Expression,
		ParsedTree:         parsedTree,
		OperationType:      "parse",
	}

	// 填充扩展数值字段
	c.populateNumericFields(result, params)
	return result, nil
}

// differentiateExpression 符号求导（修复版）
func (c *SymbolicCalcCalculatorFixed) differentiateExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	// 1. 解析表达式为AST
	ast := c.parseToAST(params.Expression)

	// 2. 对AST进行求导
	derivativeAST := c.diffAST(ast, params.Variable)

	// 3. 将AST转换回字符串
	derivative := c.astToString(derivativeAST)

	// 4. 化简导数结果
	simplified := c.simplifyPolynomial(derivative)

	// 5. 计算数值结果（如果提供了变量值）
	numericValue := 0.0
	if params.XValue != 0 || params.YValue != 0 || params.ZValue != 0 {
		variables := map[string]float64{
			"x": params.XValue,
			"y": params.YValue,
			"z": params.ZValue,
		}
		numericValue = c.evaluateAST(derivativeAST, variables)
	}

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   derivative,
		Derivative:         derivative,
		Simplified:         simplified,
		NumericValue:       numericValue,
		OperationType:      "differentiate",
	}

	// 填充扩展数值字段
	c.populateNumericFields(result, params)
	return result, nil
}

// simplifyExpression 表达式化简（修复版）
func (c *SymbolicCalcCalculatorFixed) simplifyExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	// 解析并化简表达式
	ast := c.parseToAST(params.Expression)
	simplifiedAST := c.simplifyAST(ast)
	simplified := c.astToString(simplifiedAST)

	// 进一步多项式化简
	simplified = c.simplifyPolynomial(simplified)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   simplified,
		Simplified:         simplified,
		OperationType:      "simplify",
	}

	// 填充扩展数值字段
	c.populateNumericFields(result, params)
	return result, nil
}

// evaluateExpression 表达式求值（修复版）
func (c *SymbolicCalcCalculatorFixed) evaluateExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	variables := map[string]float64{
		"x": params.XValue,
		"y": params.YValue,
		"z": params.ZValue,
	}

	// 解析并求值
	ast := c.parseToAST(params.Expression)
	numericValue := c.evaluateAST(ast, variables)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   fmt.Sprintf("%.6f", numericValue),
		NumericValue:       numericValue,
		Variables:          variables,
		OperationType:      "evaluate",
	}

	// 填充扩展数值字段
	c.populateNumericFields(result, params)
	return result, nil
}

// ==================== AST 节点定义 ====================

// ASTNode 抽象语法树节点
type ASTNode struct {
	Type     string     // 节点类型: "number", "variable", "add", "sub", "mul", "div", "pow", "neg"
	Value    float64    // 数值（仅用于number类型）
	Name     string     // 变量名（仅用于variable类型）
	Children []*ASTNode // 子节点
}

// NewNumberNode 创建数值节点
func NewNumberNode(value float64) *ASTNode {
	return &ASTNode{Type: "number", Value: value}
}

// NewVariableNode 创建变量节点
func NewVariableNode(name string) *ASTNode {
	return &ASTNode{Type: "variable", Name: name}
}

// NewBinaryNode 创建二元运算节点
func NewBinaryNode(opType string, left, right *ASTNode) *ASTNode {
	return &ASTNode{Type: opType, Children: []*ASTNode{left, right}}
}

// NewUnaryNode 创建一元运算节点
func NewUnaryNode(opType string, child *ASTNode) *ASTNode {
	return &ASTNode{Type: opType, Children: []*ASTNode{child}}
}

// ==================== 表达式解析 ====================

// parseToAST 将表达式字符串解析为AST
func (c *SymbolicCalcCalculatorFixed) parseToAST(expr string) *ASTNode {
	expr = strings.ReplaceAll(expr, " ", "")
	if expr == "" {
		return NewNumberNode(0)
	}
	tokens := c.tokenize(expr)
	return c.parseExpressionTokens(tokens, 0, len(tokens)-1)
}

// tokenize 分词
func (c *SymbolicCalcCalculatorFixed) tokenize(expr string) []string {
	var tokens []string
	i := 0
	for i < len(expr) {
		ch := expr[i]
		if ch == ' ' {
			i++
			continue
		}
		// 数字（包括小数）
		if ch >= '0' && ch <= '9' || ch == '.' {
			j := i
			for j < len(expr) && ((expr[j] >= '0' && expr[j] <= '9') || expr[j] == '.') {
				j++
			}
			tokens = append(tokens, expr[i:j])
			i = j
		} else if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			// 变量名
			j := i
			for j < len(expr) && ((expr[j] >= 'a' && expr[j] <= 'z') || (expr[j] >= 'A' && expr[j] <= 'Z')) {
				j++
			}
			tokens = append(tokens, expr[i:j])
			i = j
		} else {
			// 运算符
			tokens = append(tokens, string(ch))
			i++
		}
	}
	return tokens
}

// parseExpressionTokens 递归解析表达式tokens
func (c *SymbolicCalcCalculatorFixed) parseExpressionTokens(tokens []string, start, end int) *ASTNode {
	if start > end {
		return NewNumberNode(0)
	}

	// 查找最低优先级的运算符
	// 优先级: + - (最低) > * / > ^ (最高)
	parenDepth := 0
	addSubPos := -1
	mulDivPos := -1
	powPos := -1

	for i := start; i <= end; i++ {
		token := tokens[i]
		switch token {
		case "(":
			parenDepth++
		case ")":
			parenDepth--
		case "+", "-":
			if parenDepth == 0 {
				addSubPos = i
			}
		case "*", "/":
			if parenDepth == 0 {
				mulDivPos = i
			}
		case "^":
			if parenDepth == 0 {
				powPos = i
			}
		}
	}

	// 根据优先级选择运算符
	if addSubPos != -1 {
		left := c.parseExpressionTokens(tokens, start, addSubPos-1)
		right := c.parseExpressionTokens(tokens, addSubPos+1, end)
		if tokens[addSubPos] == "+" {
			return NewBinaryNode("add", left, right)
		}
		return NewBinaryNode("sub", left, right)
	}

	if mulDivPos != -1 {
		left := c.parseExpressionTokens(tokens, start, mulDivPos-1)
		right := c.parseExpressionTokens(tokens, mulDivPos+1, end)
		if tokens[mulDivPos] == "*" {
			return NewBinaryNode("mul", left, right)
		}
		return NewBinaryNode("div", left, right)
	}

	if powPos != -1 {
		left := c.parseExpressionTokens(tokens, start, powPos-1)
		right := c.parseExpressionTokens(tokens, powPos+1, end)
		return NewBinaryNode("pow", left, right)
	}

	// 处理括号
	if tokens[start] == "(" && tokens[end] == ")" {
		return c.parseExpressionTokens(tokens, start+1, end-1)
	}

	// 处理单一项
	if start == end {
		token := tokens[start]
		// 尝试解析为数字
		if val, err := strconv.ParseFloat(token, 64); err == nil {
			return NewNumberNode(val)
		}
		// 否则是变量
		return NewVariableNode(token)
	}

	// 处理一元负号
	if tokens[start] == "-" {
		child := c.parseExpressionTokens(tokens, start+1, end)
		return NewUnaryNode("neg", child)
	}

	return NewNumberNode(0)
}

// ==================== 符号求导 ====================

// diffAST 对AST进行符号求导
func (c *SymbolicCalcCalculatorFixed) diffAST(node *ASTNode, variable string) *ASTNode {
	if node == nil {
		return NewNumberNode(0)
	}

	switch node.Type {
	case "number":
		// 常数导数为0
		return NewNumberNode(0)

	case "variable":
		// 变量导数：如果是求导变量则为1，否则为0
		if node.Name == variable {
			return NewNumberNode(1)
		}
		return NewNumberNode(0)

	case "add":
		// (u + v)' = u' + v'
		left := c.diffAST(node.Children[0], variable)
		right := c.diffAST(node.Children[1], variable)
		return NewBinaryNode("add", left, right)

	case "sub":
		// (u - v)' = u' - v'
		left := c.diffAST(node.Children[0], variable)
		right := c.diffAST(node.Children[1], variable)
		return NewBinaryNode("sub", left, right)

	case "mul":
		// (u * v)' = u' * v + u * v'
		u := node.Children[0]
		v := node.Children[1]
		uPrime := c.diffAST(u, variable)
		vPrime := c.diffAST(v, variable)
		term1 := NewBinaryNode("mul", uPrime, v)
		term2 := NewBinaryNode("mul", u, vPrime)
		return NewBinaryNode("add", term1, term2)

	case "div":
		// (u / v)' = (u' * v - u * v') / v^2
		u := node.Children[0]
		v := node.Children[1]
		uPrime := c.diffAST(u, variable)
		vPrime := c.diffAST(v, variable)
		term1 := NewBinaryNode("mul", uPrime, v)
		term2 := NewBinaryNode("mul", u, vPrime)
		numerator := NewBinaryNode("sub", term1, term2)
		denominator := NewBinaryNode("pow", v, NewNumberNode(2))
		return NewBinaryNode("div", numerator, denominator)

	case "pow":
		// 处理幂函数
		base := node.Children[0]
		exponent := node.Children[1]

		// 情况1: x^n (变量为底数，常数为指数)
		if base.Type == "variable" && base.Name == variable && exponent.Type == "number" {
			n := exponent.Value
			if n == 0 {
				return NewNumberNode(0)
			}
			if n == 1 {
				return NewNumberNode(1)
			}
			// d/dx(x^n) = n * x^(n-1)
			newExp := NewNumberNode(n - 1)
			newPow := NewBinaryNode("pow", base, newExp)
			return NewBinaryNode("mul", NewNumberNode(n), newPow)
		}

		// 情况2: f(x)^n (复合函数)
		if exponent.Type == "number" {
			n := exponent.Value
			// d/dx(f^n) = n * f^(n-1) * f'
			newExp := NewNumberNode(n - 1)
			newPow := NewBinaryNode("pow", base, newExp)
			part1 := NewBinaryNode("mul", NewNumberNode(n), newPow)
			basePrime := c.diffAST(base, variable)
			return NewBinaryNode("mul", part1, basePrime)
		}

		// 情况3: 一般情况使用对数求导法
		// d/dx(u^v) = u^v * (v' * ln(u) + v * u'/u)
		u := base
		v := exponent
		uPrime := c.diffAST(u, variable)
		vPrime := c.diffAST(v, variable)
		powNode := NewBinaryNode("pow", u, v)
		lnU := &ASTNode{Type: "ln", Children: []*ASTNode{u}}
		termA := NewBinaryNode("mul", vPrime, lnU)
		divNode := NewBinaryNode("div", uPrime, u)
		termB := NewBinaryNode("mul", v, divNode)
		sum := NewBinaryNode("add", termA, termB)
		return NewBinaryNode("mul", powNode, sum)

	case "neg":
		// (-u)' = -u'
		child := c.diffAST(node.Children[0], variable)
		return NewUnaryNode("neg", child)

	default:
		return NewNumberNode(0)
	}
}

// ==================== AST化简 ====================

// simplifyAST 化简AST
func (c *SymbolicCalcCalculatorFixed) simplifyAST(node *ASTNode) *ASTNode {
	if node == nil {
		return nil
	}

	// 递归化简子节点
	for i, child := range node.Children {
		node.Children[i] = c.simplifyAST(child)
	}

	switch node.Type {
	case "add":
		left, right := node.Children[0], node.Children[1]
		// 0 + x = x
		if left.Type == "number" && left.Value == 0 {
			return right
		}
		// x + 0 = x
		if right.Type == "number" && right.Value == 0 {
			return left
		}
		// 常数相加
		if left.Type == "number" && right.Type == "number" {
			return NewNumberNode(left.Value + right.Value)
		}

	case "sub":
		left, right := node.Children[0], node.Children[1]
		// x - 0 = x
		if right.Type == "number" && right.Value == 0 {
			return left
		}
		// 0 - x = -x
		if left.Type == "number" && left.Value == 0 {
			return NewUnaryNode("neg", right)
		}
		// 常数相减
		if left.Type == "number" && right.Type == "number" {
			return NewNumberNode(left.Value - right.Value)
		}
		// x - x = 0
		if c.astEqual(left, right) {
			return NewNumberNode(0)
		}

	case "mul":
		left, right := node.Children[0], node.Children[1]
		// 0 * x = 0
		if (left.Type == "number" && left.Value == 0) || (right.Type == "number" && right.Value == 0) {
			return NewNumberNode(0)
		}
		// 1 * x = x
		if left.Type == "number" && left.Value == 1 {
			return right
		}
		// x * 1 = x
		if right.Type == "number" && right.Value == 1 {
			return left
		}
		// 常数相乘
		if left.Type == "number" && right.Type == "number" {
			return NewNumberNode(left.Value * right.Value)
		}

	case "div":
		left, right := node.Children[0], node.Children[1]
		// 0 / x = 0
		if left.Type == "number" && left.Value == 0 {
			return NewNumberNode(0)
		}
		// x / 1 = x
		if right.Type == "number" && right.Value == 1 {
			return left
		}
		// 常数相除
		if left.Type == "number" && right.Type == "number" && right.Value != 0 {
			return NewNumberNode(left.Value / right.Value)
		}

	case "pow":
		base, exp := node.Children[0], node.Children[1]
		// x^0 = 1
		if exp.Type == "number" && exp.Value == 0 {
			return NewNumberNode(1)
		}
		// x^1 = x
		if exp.Type == "number" && exp.Value == 1 {
			return base
		}
		// 0^x = 0 (x > 0)
		if base.Type == "number" && base.Value == 0 {
			return NewNumberNode(0)
		}
		// 1^x = 1
		if base.Type == "number" && base.Value == 1 {
			return NewNumberNode(1)
		}
		// 常数幂
		if base.Type == "number" && exp.Type == "number" {
			return NewNumberNode(math.Pow(base.Value, exp.Value))
		}

	case "neg":
		child := node.Children[0]
		// -(-x) = x
		if child.Type == "neg" {
			return child.Children[0]
		}
		// -(number)
		if child.Type == "number" {
			return NewNumberNode(-child.Value)
		}
	}

	return node
}

// astEqual 判断两个AST是否相等
func (c *SymbolicCalcCalculatorFixed) astEqual(a, b *ASTNode) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Type != b.Type {
		return false
	}
	if a.Type == "number" {
		return a.Value == b.Value
	}
	if a.Type == "variable" {
		return a.Name == b.Name
	}
	if len(a.Children) != len(b.Children) {
		return false
	}
	for i := range a.Children {
		if !c.astEqual(a.Children[i], b.Children[i]) {
			return false
		}
	}
	return true
}

// ==================== AST求值 ====================

// evaluateAST 对AST进行数值求值
func (c *SymbolicCalcCalculatorFixed) evaluateAST(node *ASTNode, variables map[string]float64) float64 {
	if node == nil {
		return 0
	}

	switch node.Type {
	case "number":
		return node.Value

	case "variable":
		if val, ok := variables[node.Name]; ok {
			return val
		}
		return 0

	case "add":
		return c.evaluateAST(node.Children[0], variables) + c.evaluateAST(node.Children[1], variables)

	case "sub":
		return c.evaluateAST(node.Children[0], variables) - c.evaluateAST(node.Children[1], variables)

	case "mul":
		return c.evaluateAST(node.Children[0], variables) * c.evaluateAST(node.Children[1], variables)

	case "div":
		divisor := c.evaluateAST(node.Children[1], variables)
		if divisor != 0 {
			return c.evaluateAST(node.Children[0], variables) / divisor
		}
		return 0

	case "pow":
		base := c.evaluateAST(node.Children[0], variables)
		exp := c.evaluateAST(node.Children[1], variables)
		return math.Pow(base, exp)

	case "neg":
		return -c.evaluateAST(node.Children[0], variables)

	case "ln":
		arg := c.evaluateAST(node.Children[0], variables)
		if arg > 0 {
			return math.Log(arg)
		}
		return 0

	default:
		return 0
	}
}

// ==================== AST转字符串 ====================

// astToString 将AST转换为字符串表达式
func (c *SymbolicCalcCalculatorFixed) astToString(node *ASTNode) string {
	if node == nil {
		return "0"
	}

	switch node.Type {
	case "number":
		// 整数不显示小数点
		if node.Value == float64(int(node.Value)) {
			return fmt.Sprintf("%.0f", node.Value)
		}
		return fmt.Sprintf("%g", node.Value)

	case "variable":
		return node.Name

	case "add":
		return fmt.Sprintf("%s + %s", c.astToString(node.Children[0]), c.astToString(node.Children[1]))

	case "sub":
		return fmt.Sprintf("%s - %s", c.astToString(node.Children[0]), c.astToString(node.Children[1]))

	case "mul":
		left := c.astToString(node.Children[0])
		right := c.astToString(node.Children[1])
		// 系数为1时省略
		if node.Children[0].Type == "number" && node.Children[0].Value == 1 {
			return right
		}
		return fmt.Sprintf("%s*%s", left, right)

	case "div":
		return fmt.Sprintf("(%s)/(%s)", c.astToString(node.Children[0]), c.astToString(node.Children[1]))

	case "pow":
		base := c.astToString(node.Children[0])
		exp := c.astToString(node.Children[1])
		// 简化显示
		if node.Children[0].Type != "variable" && node.Children[0].Type != "number" {
			base = "(" + base + ")"
		}
		return fmt.Sprintf("%s^%s", base, exp)

	case "neg":
		child := c.astToString(node.Children[0])
		if node.Children[0].Type == "add" || node.Children[0].Type == "sub" {
			child = "(" + child + ")"
		}
		return "-" + child

	case "ln":
		return fmt.Sprintf("ln(%s)", c.astToString(node.Children[0]))

	default:
		return "0"
	}
}

// ==================== 多项式化简 ====================

// simplifyPolynomial 化简多项式表达式
func (c *SymbolicCalcCalculatorFixed) simplifyPolynomial(expr string) string {
	// 解析为AST
	ast := c.parseToAST(expr)

	// 收集同类项
	terms := c.collectTerms(ast)

	// 合并同类项
	merged := c.mergeTerms(terms)

	// 构建化简后的表达式
	return c.buildPolynomialString(merged)
}

// Term 表示多项式中的一项
type Term struct {
	Coefficient float64 // 系数
	Power       float64 // 变量的幂次
	Variable    string  // 变量名
}

// collectTerms 收集表达式中的各项
func (c *SymbolicCalcCalculatorFixed) collectTerms(node *ASTNode) []Term {
	if node == nil {
		return nil
	}

	var terms []Term

	switch node.Type {
	case "add":
		terms = append(terms, c.collectTerms(node.Children[0])...)
		terms = append(terms, c.collectTerms(node.Children[1])...)

	case "sub":
		leftTerms := c.collectTerms(node.Children[0])
		rightTerms := c.collectTerms(node.Children[1])
		// 右侧取负
		for i := range rightTerms {
			rightTerms[i].Coefficient = -rightTerms[i].Coefficient
		}
		terms = append(terms, leftTerms...)
		terms = append(terms, rightTerms...)

	case "mul":
		// 分析乘法项
		term := c.analyzeTerm(node)
		terms = append(terms, term)

	case "pow":
		term := c.analyzeTerm(node)
		terms = append(terms, term)

	case "number":
		terms = append(terms, Term{Coefficient: node.Value, Power: 0, Variable: ""})

	case "variable":
		terms = append(terms, Term{Coefficient: 1, Power: 1, Variable: node.Name})

	case "neg":
		childTerms := c.collectTerms(node.Children[0])
		for i := range childTerms {
			childTerms[i].Coefficient = -childTerms[i].Coefficient
		}
		terms = append(terms, childTerms...)

	default:
		term := c.analyzeTerm(node)
		terms = append(terms, term)
	}

	return terms
}

// analyzeTerm 分析单项式
func (c *SymbolicCalcCalculatorFixed) analyzeTerm(node *ASTNode) Term {
	term := Term{Coefficient: 1, Power: 0, Variable: ""}

	// 递归分析乘法结构
	c.analyzeTermRecursive(node, &term)

	return term
}

// analyzeTermRecursive 递归分析项的结构
func (c *SymbolicCalcCalculatorFixed) analyzeTermRecursive(node *ASTNode, term *Term) {
	if node == nil {
		return
	}

	switch node.Type {
	case "number":
		term.Coefficient *= node.Value

	case "variable":
		if term.Variable == "" || term.Variable == node.Name {
			term.Variable = node.Name
			term.Power += 1
		}

	case "pow":
		base := node.Children[0]
		exp := node.Children[1]
		if base.Type == "variable" && exp.Type == "number" {
			if term.Variable == "" || term.Variable == base.Name {
				term.Variable = base.Name
				term.Power += exp.Value
			}
		}

	case "mul":
		c.analyzeTermRecursive(node.Children[0], term)
		c.analyzeTermRecursive(node.Children[1], term)

	case "neg":
		term.Coefficient *= -1
		c.analyzeTermRecursive(node.Children[0], term)
	}
}

// mergeTerms 合并同类项
func (c *SymbolicCalcCalculatorFixed) mergeTerms(terms []Term) []Term {
	// 按变量和幂次分组
	groups := make(map[string][]Term)

	for _, term := range terms {
		key := fmt.Sprintf("%s_%.0f", term.Variable, term.Power)
		groups[key] = append(groups[key], term)
	}

	// 合并每组
	var merged []Term
	for _, group := range groups {
		if len(group) == 0 {
			continue
		}
		var sum float64
		for _, t := range group {
			sum += t.Coefficient
		}
		if sum != 0 {
			merged = append(merged, Term{
				Coefficient: sum,
				Power:       group[0].Power,
				Variable:    group[0].Variable,
			})
		}
	}

	// 按幂次降序排序
	for i := 0; i < len(merged); i++ {
		for j := i + 1; j < len(merged); j++ {
			if merged[j].Power > merged[i].Power {
				merged[i], merged[j] = merged[j], merged[i]
			}
		}
	}

	return merged
}

// buildPolynomialString 构建多项式字符串
func (c *SymbolicCalcCalculatorFixed) buildPolynomialString(terms []Term) string {
	if len(terms) == 0 {
		return "0"
	}

	var parts []string
	for i, term := range terms {
		part := c.termToString(term, i == 0)
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return "0"
	}

	// 直接连接，因为 termToString 已经处理了符号
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "-") {
			// 负数，将 "+ -" 替换为 "- "
			result += " - " + strings.TrimPrefix(parts[i], "- ")
		} else {
			result += " + " + parts[i]
		}
	}
	return result
}

// termToString 将单项转换为字符串
func (c *SymbolicCalcCalculatorFixed) termToString(term Term, isFirst bool) string {
	if term.Coefficient == 0 {
		return ""
	}

	coef := term.Coefficient
	sign := ""
	if !isFirst {
		if coef < 0 {
			sign = "- "
			coef = -coef
		} else {
			sign = ""
		}
	} else if coef < 0 {
		sign = "-"
		coef = -coef
	}

	// 常数项
	if term.Power == 0 || term.Variable == "" {
		if coef == float64(int(coef)) {
			return fmt.Sprintf("%s%.0f", sign, coef)
		}
		return fmt.Sprintf("%s%g", sign, coef)
	}

	// 变量项
	var coefStr string
	if coef == 1 {
		coefStr = ""
	} else if coef == float64(int(coef)) {
		coefStr = fmt.Sprintf("%.0f*", coef)
	} else {
		coefStr = fmt.Sprintf("%g*", coef)
	}

	var varStr string
	if term.Power == 1 {
		varStr = term.Variable
	} else if term.Power == float64(int(term.Power)) {
		varStr = fmt.Sprintf("%s^%.0f", term.Variable, term.Power)
	} else {
		varStr = fmt.Sprintf("%s^%g", term.Variable, term.Power)
	}

	return sign + coefStr + varStr
}

// ==================== 辅助函数 ====================

// buildParseTree 构建解析树（简化版）
func (c *SymbolicCalcCalculatorFixed) buildParseTree(expr string) map[string]interface{} {
	ast := c.parseToAST(expr)
	return c.astToMap(ast)
}

// astToMap 将AST转换为map表示
func (c *SymbolicCalcCalculatorFixed) astToMap(node *ASTNode) map[string]interface{} {
	if node == nil {
		return nil
	}

	result := map[string]interface{}{
		"type": node.Type,
	}

	switch node.Type {
	case "number":
		result["value"] = node.Value
	case "variable":
		result["name"] = node.Name
	default:
		if len(node.Children) > 0 {
			children := make([]map[string]interface{}, len(node.Children))
			for i, child := range node.Children {
				children[i] = c.astToMap(child)
			}
			result["children"] = children
		}
	}

	return result
}

// populateNumericFields 填充扩展数值字段
func (c *SymbolicCalcCalculatorFixed) populateNumericFields(result *SymbolicResultFixed, params *SymbolicParamsFixed) {
	if result == nil {
		return
	}

	// 1. 表达式复杂度评分
	complexity := c.calculateExpressionComplexity(result.OriginalExpression)
	result.ExpressionComplexity = complexity

	// 2. 变量数量统计
	varCount := c.countVariables(result.OriginalExpression)
	result.VariableCount = varCount

	// 3. 项数统计
	termCount := c.countTerms(result.OriginalExpression)
	result.TermCount = float64(termCount)

	// 4. 解析树深度估算
	treeDepth := c.estimateTreeDepth(result.OriginalExpression)
	result.TreeDepth = float64(treeDepth)

	// 5. 综合评估得分
	evalScore := c.calculateEvaluationScore(result, params)
	result.EvaluationScore = evalScore
}

// calculateExpressionComplexity 计算表达式复杂度
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

// countVariables 统计表达式中的变量数量
func (c *SymbolicCalcCalculatorFixed) countVariables(expr string) int {
	variables := make(map[rune]bool)
	for _, char := range expr {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			variables[char] = true
		}
	}
	return len(variables)
}

// countTerms 统计表达式中的项数
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

// estimateTreeDepth 估算解析树深度
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

// calculateEvaluationScore 计算综合评估得分
func (c *SymbolicCalcCalculatorFixed) calculateEvaluationScore(result *SymbolicResultFixed, params *SymbolicParamsFixed) float64 {
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
	return "修复版符号计算器，支持正确的表达式解析、符号求导和表达式化简"
}
