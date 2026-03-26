package calculator

import (
	"fmt"
	"time"
)

type SymbolicCalcCalculatorFixed struct {
	*BaseCalculator
}

func NewSymbolicCalcCalculatorFixed() *SymbolicCalcCalculatorFixed {
	return &SymbolicCalcCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"symbolic_calc_fixed",
			"修复版符号计算器，支持表达式解析、符号求导和表达式化简",
		),
	}
}

type SymbolicParamsFixed struct {
	Operation  string  `json:"operation"`
	Expression string  `json:"expression"`
	Variable   string  `json:"variable"`
	XValue     float64 `json:"x_value"`
	YValue     float64 `json:"y_value"`
	ZValue     float64 `json:"z_value"`
}

type SymbolicResultFixed struct {
	OriginalExpression   string             `json:"original_expression"`
	ResultExpression     string             `json:"result_expression"`
	ParsedTree           interface{}        `json:"parsed_tree"`
	Derivative           string             `json:"derivative"`
	Simplified           string             `json:"simplified"`
	NumericValue         float64            `json:"numeric_value"`
	Variables            map[string]float64 `json:"variables"`
	OperationType        string             `json:"operation_type"`
	ExpressionComplexity float64            `json:"expression_complexity,omitempty"`
	VariableCount        int                `json:"variable_count,omitempty"`
	TermCount            float64            `json:"term_count,omitempty"`
	TreeDepth            float64            `json:"tree_depth,omitempty"`
	EvaluationScore      float64            `json:"evaluation_score,omitempty"`
	Error                string             `json:"error,omitempty"`
}

func (c *SymbolicCalcCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	symbolicParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(symbolicParams); err != nil {
		return nil, err
	}

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

func (c *SymbolicCalcCalculatorFixed) parseParams(params interface{}) (*SymbolicParamsFixed, error) {
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

	paramsObj := &SymbolicParamsFixed{
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

func (c *SymbolicCalcCalculatorFixed) validateParams(params *SymbolicParamsFixed) error {
	if params.Expression == "" {
		return fmt.Errorf("表达式不能为空")
	}
	return nil
}

func (c *SymbolicCalcCalculatorFixed) parseExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	cleanExpr := CleanExpression(params.Expression)

	root, err := ParseExpression(cleanExpr)
	if err != nil {
		return nil, fmt.Errorf("表达式解析失败: %v", err)
	}

	parsedTree := NodeToMap(root)
	variables := ExtractVariables(root)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   FormatExpression(root),
		ParsedTree:         parsedTree,
		Variables:          make(map[string]float64),
		OperationType:      "parse",
	}

	for _, v := range variables {
		result.Variables[v] = 0
	}

	c.populateNumericFields(result, params)
	return result, nil
}

func (c *SymbolicCalcCalculatorFixed) differentiateExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	cleanExpr := CleanExpression(params.Expression)

	root, err := ParseExpression(cleanExpr)
	if err != nil {
		return nil, fmt.Errorf("表达式解析失败: %v", err)
	}

	parsedTree := NodeToMap(root)
	variables := ExtractVariables(root)

	derivativeNode := Differentiate(root, params.Variable)
	simplifiedDerivative := Simplify(derivativeNode)
	derivativeStr := FormatExpression(simplifiedDerivative)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   derivativeStr,
		ParsedTree:         parsedTree,
		Derivative:         derivativeStr,
		Simplified:         derivativeStr,
		Variables:          make(map[string]float64),
		OperationType:      "differentiate",
	}

	for _, v := range variables {
		result.Variables[v] = 0
	}

	c.populateNumericFields(result, params)
	return result, nil
}

func (c *SymbolicCalcCalculatorFixed) simplifyExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	cleanExpr := CleanExpression(params.Expression)

	root, err := ParseExpression(cleanExpr)
	if err != nil {
		return nil, fmt.Errorf("表达式解析失败: %v", err)
	}

	parsedTree := NodeToMap(root)
	variables := ExtractVariables(root)

	simplifiedNode := Simplify(root)
	simplifiedStr := FormatExpression(simplifiedNode)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   simplifiedStr,
		ParsedTree:         parsedTree,
		Simplified:         simplifiedStr,
		Variables:          make(map[string]float64),
		OperationType:      "simplify",
	}

	for _, v := range variables {
		result.Variables[v] = 0
	}

	c.populateNumericFields(result, params)
	return result, nil
}

func (c *SymbolicCalcCalculatorFixed) evaluateExpression(params *SymbolicParamsFixed) (*SymbolicResultFixed, error) {
	cleanExpr := CleanExpression(params.Expression)

	root, err := ParseExpression(cleanExpr)
	if err != nil {
		return nil, fmt.Errorf("表达式解析失败: %v", err)
	}

	parsedTree := NodeToMap(root)

	varValues := map[string]float64{
		"x": params.XValue,
		"y": params.YValue,
		"z": params.ZValue,
	}

	numericValue, evalErr := Evaluate(root, varValues)

	result := &SymbolicResultFixed{
		OriginalExpression: params.Expression,
		ResultExpression:   fmt.Sprintf("%.6g", numericValue),
		ParsedTree:         parsedTree,
		NumericValue:       numericValue,
		Variables:          varValues,
		OperationType:      "evaluate",
	}

	if evalErr != nil {
		result.Error = evalErr.Error()
	}

	c.populateNumericFields(result, params)
	return result, nil
}

func (c *SymbolicCalcCalculatorFixed) populateNumericFields(result *SymbolicResultFixed, params *SymbolicParamsFixed) {
	if result == nil {
		return
	}

	complexity := c.calculateExpressionComplexity(result.OriginalExpression)
	result.ExpressionComplexity = complexity

	varCount := len(result.Variables)
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

	operatorCount := 0
	for _, char := range expr {
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '^' {
			operatorCount++
		}
	}
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
		baseScore += 10.0
	}

	if result.Derivative != "" && result.Derivative != result.OriginalExpression {
		baseScore += 15.0
	}

	if result.Simplified != "" {
		baseScore += 10.0
	}

	if result.Error != "" {
		baseScore -= 20.0
	}

	return max(0, min(100, baseScore))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (c *SymbolicCalcCalculatorFixed) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

func (c *SymbolicCalcCalculatorFixed) Description() string {
	return "修复版符号计算器，支持表达式解析、符号求导和表达式化简"
}

func (c *SymbolicCalcCalculatorFixed) CalculateDerivativeWithEval(params *SymbolicParamsFixed, evalPoint float64) (*SymbolicResultFixed, error) {
	result, err := c.differentiateExpression(params)
	if err != nil {
		return nil, err
	}

	cleanDerivative := CleanExpression(result.Derivative)
	derivativeRoot, err := ParseExpression(cleanDerivative)
	if err != nil {
		return result, nil
	}

	varValues := map[string]float64{
		params.Variable: evalPoint,
	}

	numericValue, evalErr := Evaluate(derivativeRoot, varValues)
	result.NumericValue = numericValue
	result.Variables = varValues

	if evalErr != nil {
		result.Error = evalErr.Error()
	}

	return result, nil
}

type CompareResult struct {
	OriginalResult interface{}  `json:"original_result"`
	FixedResult    interface{}  `json:"fixed_result"`
	Differences    []Difference `json:"differences"`
	IsFixed        bool         `json:"is_fixed"`
	ComparisonTime string       `json:"comparison_time"`
}

type Difference struct {
	Field    string `json:"field"`
	Original string `json:"original"`
	Fixed    string `json:"fixed"`
	Issue    string `json:"issue"`
}

func CompareResults(original *SymbolicResult, fixed *SymbolicResultFixed) *CompareResult {
	differences := []Difference{}
	isFixed := true

	if original.Derivative != fixed.Derivative {
		differences = append(differences, Difference{
			Field:    "derivative",
			Original: original.Derivative,
			Fixed:    fixed.Derivative,
			Issue:    "导数计算结果不同",
		})
	}

	if original.Simplified != fixed.Simplified {
		differences = append(differences, Difference{
			Field:    "simplified",
			Original: original.Simplified,
			Fixed:    fixed.Simplified,
			Issue:    "化简结果不同",
		})
	}

	if original.NumericValue != fixed.NumericValue {
		differences = append(differences, Difference{
			Field:    "numeric_value",
			Original: fmt.Sprintf("%.6g", original.NumericValue),
			Fixed:    fmt.Sprintf("%.6g", fixed.NumericValue),
			Issue:    "数值结果不同",
		})
	}

	if original.ParsedTree == nil && fixed.ParsedTree != nil {
		differences = append(differences, Difference{
			Field:    "parsed_tree",
			Original: "null",
			Fixed:    "valid tree",
			Issue:    "原版解析树为空",
		})
	}

	if original.Variables == nil && fixed.Variables != nil {
		differences = append(differences, Difference{
			Field:    "variables",
			Original: "null",
			Fixed:    fmt.Sprintf("%v", fixed.Variables),
			Issue:    "原版变量提取失败",
		})
	}

	if len(differences) == 0 {
		isFixed = false
	}

	return &CompareResult{
		OriginalResult: original,
		FixedResult:    fixed,
		Differences:    differences,
		IsFixed:        isFixed,
		ComparisonTime: time.Now().Format(time.RFC3339),
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
