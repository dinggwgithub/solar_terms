package calculator

import (
	"fmt"
	"math"
	"regexp"
	"scientific_calc_bugs/internal/bugs"
	"strconv"
	"strings"
)

// SymbolicCalcCalculator 符号计算器
type SymbolicCalcCalculator struct {
	*BaseCalculator
}

// NewSymbolicCalcCalculator 创建新的符号计算器
func NewSymbolicCalcCalculator() *SymbolicCalcCalculator {
	return &SymbolicCalcCalculator{
		BaseCalculator: NewBaseCalculator(
			"symbolic_calc",
			"符号计算器，支持表达式解析、符号求导和表达式化简",
		),
	}
}

// SymbolicParams 符号计算参数
type SymbolicParams struct {
	Operation   string  `json:"operation"`   // 操作类型：parse, differentiate, simplify, evaluate
	Expression  string  `json:"expression"`  // 数学表达式
	Variable    string  `json:"variable"`    // 变量名（求导用）
	XValue     float64 `json:"x_value"`     // x的值（求值用）
	YValue     float64 `json:"y_value"`     // y的值（求值用）
	ZValue     float64 `json:"z_value"`     // z的值（求值用）
}

// SymbolicResult 符号计算结果
type SymbolicResult struct {
	OriginalExpression string      `json:"original_expression"` // 原始表达式
	ResultExpression   string      `json:"result_expression"`   // 结果表达式
	ParsedTree         interface{} `json:"parsed_tree"`         // 解析树
	Derivative         string      `json:"derivative"`         // 导数结果
	Simplified         string      `json:"simplified"`          // 化简结果
	NumericValue      float64     `json:"numeric_value"`       // 数值结果
	Variables         map[string]float64 `json:"variables"`     // 变量值
	OperationType     string      `json:"operation_type"`      // 操作类型
}

// Calculate 执行符号计算
func (c *SymbolicCalcCalculator) Calculate(params interface{}) (interface{}, error) {
	symbolicParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(symbolicParams); err != nil {
		return nil, err
	}

	// 根据操作类型执行计算
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
func (c *SymbolicCalcCalculator) parseParams(params interface{}) (*SymbolicParams, error) {
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
	paramsObj := &SymbolicParams{
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
func (c *SymbolicCalcCalculator) validateParams(params *SymbolicParams) error {
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
func (c *SymbolicCalcCalculator) isValidExpression(expr string) bool {
	// 简单的表达式格式验证
	// 允许字母、数字、运算符、括号等
	validPattern := `^[a-zA-Z0-9\+\-\*\/\(\)\^\s\.]+$`
	matched, _ := regexp.MatchString(validPattern, expr)
	return matched
}

// parseExpression 解析数学表达式
func (c *SymbolicCalcCalculator) parseExpression(params *SymbolicParams) (*SymbolicResult, error) {
	// 简化的表达式解析
	// 实际应该实现完整的语法分析
	
	parsedTree := c.buildParseTree(params.Expression)
	
	return &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:  params.Expression,
		ParsedTree:         parsedTree,
		OperationType:      "parse",
	}, nil
}

// differentiateExpression 符号求导
func (c *SymbolicCalcCalculator) differentiateExpression(params *SymbolicParams) (*SymbolicResult, error) {
	// 简化的符号求导
	// 实际应该实现完整的求导规则
	
	derivative := c.calculateDerivative(params.Expression, params.Variable)
	
	return &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:   derivative,
		Derivative:         derivative,
		OperationType:      "differentiate",
	}, nil
}

// simplifyExpression 表达式化简
func (c *SymbolicCalcCalculator) simplifyExpression(params *SymbolicParams) (*SymbolicResult, error) {
	// 简化的表达式化简
	// 实际应该实现完整的化简规则
	
	simplified := c.simplify(params.Expression)
	
	return &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:   simplified,
		Simplified:         simplified,
		OperationType:      "simplify",
	}, nil
}

// evaluateExpression 表达式求值
func (c *SymbolicCalcCalculator) evaluateExpression(params *SymbolicParams) (*SymbolicResult, error) {
	// 表达式数值求值
	
	variables := map[string]float64{
		"x": params.XValue,
		"y": params.YValue,
		"z": params.ZValue,
	}
	
	numericValue := c.evaluate(params.Expression, variables)
	
	return &SymbolicResult{
		OriginalExpression: params.Expression,
		ResultExpression:  fmt.Sprintf("%.6f", numericValue),
		NumericValue:      numericValue,
		Variables:         variables,
		OperationType:     "evaluate",
	}, nil
}

// buildParseTree 构建解析树（简化版）
func (c *SymbolicCalcCalculator) buildParseTree(expr string) map[string]interface{} {
	// 简化的解析树构建
	// 实际应该实现完整的语法分析
	
	tree := map[string]interface{}{
		"type":     "expression",
		"value":    expr,
		"tokens":   c.tokenize(expr),
		"operator": c.findMainOperator(expr),
	}
	
	return tree
}

// tokenize 分词
func (c *SymbolicCalcCalculator) tokenize(expr string) []string {
	// 简单的分词
	tokens := []string{}
	current := ""
	
	for _, char := range expr {
		if char == ' ' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else if strings.Contains("+-*/^()", string(char)) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(char))
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		tokens = append(tokens, current)
	}
	
	return tokens
}

// findMainOperator 查找主运算符
func (c *SymbolicCalcCalculator) findMainOperator(expr string) string {
	// 查找表达式中的主运算符
	operators := []string{"+", "-", "*", "/", "^"}
	
	for _, op := range operators {
		if strings.Contains(expr, op) {
			return op
		}
	}
	
	return "none"
}

// calculateDerivative 计算导数（简化版）
func (c *SymbolicCalcCalculator) calculateDerivative(expr, variable string) string {
	// 简化的导数计算规则
	
	// 常数导数
	if matched, _ := regexp.MatchString(`^\d+$`, expr); matched {
		return "0"
	}
	
	// 变量导数
	if expr == variable {
		return "1"
	}
	
	// 幂函数导数: d/dx(x^n) = n*x^(n-1)
	if match := regexp.MustCompile(`^` + variable + `\^(\d+)$`).FindStringSubmatch(expr); match != nil {
		n, _ := strconv.Atoi(match[1])
		if n == 2 {
			return "2*" + variable
		}
		return fmt.Sprintf("%d*%s^%d", n, variable, n-1)
	}
	
	// 三角函数导数
	if strings.Contains(expr, "sin("+variable+")") {
		return "cos(" + variable + ")"
	}
	if strings.Contains(expr, "cos("+variable+")") {
		return "-sin(" + variable + ")"
	}
	
	// 默认返回原始表达式（表示无法求导）
	return "d/d" + variable + "(" + expr + ")"
}

// simplify 表达式化简（简化版）
func (c *SymbolicCalcCalculator) simplify(expr string) string {
	// 简化的化简规则
	
	// 去除多余空格
	expr = strings.ReplaceAll(expr, " ", "")
	
	// 常数运算化简
	if strings.Contains(expr, "0*") || strings.Contains(expr, "*0") {
		return "0"
	}
	
	if strings.Contains(expr, "1*") {
		return strings.ReplaceAll(expr, "1*", "")
	}
	
	if strings.Contains(expr, "*1") {
		return strings.ReplaceAll(expr, "*1", "")
	}
	
	// 幂运算化简
	if strings.Contains(expr, "^1") {
		return strings.ReplaceAll(expr, "^1", "")
	}
	
	if strings.Contains(expr, "^0") {
		return "1"
	}
	
	return expr
}

// evaluate 表达式求值
func (c *SymbolicCalcCalculator) evaluate(expr string, variables map[string]float64) float64 {
	// 简化的表达式求值
	// 实际应该实现完整的表达式求值器
	
	// 替换变量
	for varName, value := range variables {
		expr = strings.ReplaceAll(expr, varName, fmt.Sprintf("%.6f", value))
	}
	
	// 简单的数值计算
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) == 2 {
			a, _ := strconv.ParseFloat(parts[0], 64)
			b, _ := strconv.ParseFloat(parts[1], 64)
			return a + b
		}
	}
	
	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			a, _ := strconv.ParseFloat(parts[0], 64)
			b, _ := strconv.ParseFloat(parts[1], 64)
			return a * b
		}
	}
	
	if strings.Contains(expr, "^") {
		parts := strings.Split(expr, "^")
		if len(parts) == 2 {
			a, _ := strconv.ParseFloat(parts[0], 64)
			b, _ := strconv.ParseFloat(parts[1], 64)
			return math.Pow(a, b)
		}
	}
	
	// 默认返回0
	return 0.0
}

// Validate 验证输入参数
func (c *SymbolicCalcCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// Description 返回计算器描述
func (c *SymbolicCalcCalculator) Description() string {
	return "符号计算器，支持表达式解析、符号求导和表达式化简"
}

// GetSupportedBugTypes 返回支持的Bug类型
func (c *SymbolicCalcCalculator) GetSupportedBugTypes() []bugs.BugType {
	return []bugs.BugType{
		bugs.BugTypeInstability,
		bugs.BugTypeConstraint,
		bugs.BugTypePrecision,
	}
}