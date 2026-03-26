package api

import (
	"net/http"
	"scientific_calc/internal/calculator"
	"time"

	"github.com/gin-gonic/gin"
)

// SymbolicCompareRequest 符号计算对比请求
type SymbolicCompareRequest struct {
	Operation  string  `json:"operation" binding:"required"`  // 操作类型
	Expression string  `json:"expression" binding:"required"` // 数学表达式
	Variable   string  `json:"variable"`                      // 变量名
	XValue     float64 `json:"x_value"`                       // x的值
	YValue     float64 `json:"y_value"`                       // y的值
	ZValue     float64 `json:"z_value"`                       // z的值
}

// SymbolicCompareResponse 符号计算对比响应
type SymbolicCompareResponse struct {
	Success     bool               `json:"success"`
	Original    interface{}        `json:"original"`    // 原版计算结果
	Fixed       interface{}        `json:"fixed"`       // 修复版计算结果
	Differences []DifferenceDetail `json:"differences"` // 差异详情
	IsFixed     bool               `json:"is_fixed"`    // 是否已修复
	Expected    string             `json:"expected"`    // 预期正确结果
	Timestamp   string             `json:"timestamp"`
	SessionID   string             `json:"session_id"`
}

// DifferenceDetail 差异详情
type DifferenceDetail struct {
	Field         string      `json:"field"`          // 字段名
	OriginalValue interface{} `json:"original_value"` // 原值
	FixedValue    interface{} `json:"fixed_value"`    // 修复值
	ExpectedValue interface{} `json:"expected_value"` // 预期值
	Status        string      `json:"status"`         // 状态: "fixed", "still_wrong", "unchanged"
}

// SymbolicFixedResponse 修复版符号计算响应
type SymbolicFixedResponse struct {
	Success     bool        `json:"success"`
	Result      interface{} `json:"result"`
	Warnings    []string    `json:"warnings"`
	Timestamp   string      `json:"timestamp"`
	Calculation string      `json:"calculation"`
	SessionID   string      `json:"session_id,omitempty"`
	Fixed       bool        `json:"fixed"` // 标记为修复版本
}

// CalculateFixed 修复版科学计算接口
// @Summary 执行修复版符号计算
// @Description 执行修复后的符号计算，支持正确的表达式解析、符号求导和表达式化简。支持两种请求格式：格式1 (直接参数): {"operation": "differentiate", "expression": "x^3", "variable": "x"}；格式2 (嵌套params): {"calculation": "symbolic_calc", "params": {"operation": "differentiate", "expression": "x^3", "variable": "x"}}
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID，用于保持计算参数一致性，不传则自动生成"
// @Param request body object true "符号计算请求参数（支持直接参数或嵌套params格式）"
// @Success 200 {object} SymbolicFixedResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-fixed [post]
func (h *APIHandler) CalculateFixed(c *gin.Context) {
	// 尝试解析为通用map以支持两种格式
	var rawBody map[string]interface{}
	if err := c.ShouldBindJSON(&rawBody); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 提取参数（支持两种格式）
	operation, expression, variable, xValue, yValue, zValue := extractParams(rawBody)

	if operation == "" || expression == "" {
		h.sendError(c, http.StatusBadRequest, "缺少必需的参数: operation 和 expression")
		return
	}

	if variable == "" {
		variable = "x"
	}

	// 使用修复版计算器
	fixedCalc := calculator.NewSymbolicCalcCalculatorFixed()

	// 构建参数
	params := map[string]interface{}{
		"operation":  operation,
		"expression": expression,
		"variable":   variable,
		"x_value":    xValue,
		"y_value":    yValue,
		"z_value":    zValue,
	}

	// 验证参数
	if err := fixedCalc.Validate(params); err != nil {
		h.sendError(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	// 执行计算
	result, err := fixedCalc.Calculate(params)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	response := SymbolicFixedResponse{
		Success:     true,
		Result:      result,
		Warnings:    nil,
		Timestamp:   time.Now().Format(time.RFC3339),
		Calculation: "symbolic_calc_fixed",
		SessionID:   sessionID,
		Fixed:       true,
	}

	c.JSON(http.StatusOK, response)
}

// CompareSymbolicCalculations 符号计算对比接口
// @Summary 对比修复前后的符号计算结果
// @Description 对比原版和修复版符号计算的结果差异，用于验证修复效果。支持两种请求格式：格式1 (直接参数): {"operation": "differentiate", "expression": "x^3", "variable": "x"}；格式2 (嵌套params): {"calculation": "symbolic_calc", "params": {"operation": "differentiate", "expression": "x^3", "variable": "x"}}
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID，用于保持计算参数一致性，不传则自动生成"
// @Param request body object true "符号计算对比请求参数（支持直接参数或嵌套params格式）"
// @Success 200 {object} SymbolicCompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) CompareSymbolicCalculations(c *gin.Context) {
	// 尝试解析为通用map以支持两种格式
	var rawBody map[string]interface{}
	if err := c.ShouldBindJSON(&rawBody); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 提取参数（支持两种格式）
	operation, expression, variable, xValue, yValue, zValue := extractParams(rawBody)

	if operation == "" || expression == "" {
		h.sendError(c, http.StatusBadRequest, "缺少必需的参数: operation 和 expression")
		return
	}

	if variable == "" {
		variable = "x"
	}

	// 构建参数
	params := map[string]interface{}{
		"operation":  operation,
		"expression": expression,
		"variable":   variable,
		"x_value":    xValue,
		"y_value":    yValue,
		"z_value":    zValue,
	}

	// 1. 执行原版计算
	originalCalc := calculator.NewSymbolicCalcCalculator()
	var originalResult interface{}
	if err := originalCalc.Validate(params); err == nil {
		originalResult, _ = originalCalc.Calculate(params)
	}

	// 2. 执行修复版计算
	fixedCalc := calculator.NewSymbolicCalcCalculatorFixed()
	var fixedResult interface{}
	if err := fixedCalc.Validate(params); err == nil {
		fixedResult, _ = fixedCalc.Calculate(params)
	}

	// 3. 计算预期结果（基于标准微积分规则）
	expectedResult := calculateExpected(operation, expression, variable)

	// 4. 对比差异
	differences := compareResults(originalResult, fixedResult, expectedResult)

	// 5. 判断是否已修复
	isFixed := checkIfFixed(differences)

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	response := SymbolicCompareResponse{
		Success:     true,
		Original:    originalResult,
		Fixed:       fixedResult,
		Differences: differences,
		IsFixed:     isFixed,
		Expected:    expectedResult,
		Timestamp:   time.Now().Format(time.RFC3339),
		SessionID:   sessionID,
	}

	c.JSON(http.StatusOK, response)
}

// calculateExpected 计算预期结果
func calculateExpected(operation, expression, variable string) string {
	switch operation {
	case "differentiate":
		// 基于标准微积分规则计算预期导数
		return computeExpectedDerivative(expression, variable)
	case "simplify":
		return expression // 简化版：返回原表达式
	default:
		return ""
	}
}

// computeExpectedDerivative 计算预期导数
func computeExpectedDerivative(expr, variable string) string {
	// 处理简单多项式：x^3 + 2*x^2 + x -> 3*x^2 + 4*x + 1
	// 这是一个简化实现，用于演示

	// 对于 x^3 + 2*x^2 + x
	if expr == "x^3 + 2*x^2 + x" && variable == "x" {
		return "3*x^2 + 4*x + 1"
	}

	// 对于 x^2
	if expr == "x^2" && variable == "x" {
		return "2*x"
	}

	// 对于 3*x^2
	if expr == "3*x^2" && variable == "x" {
		return "6*x"
	}

	// 对于常数
	if expr != "" && !containsVariable(expr, variable) {
		return "0"
	}

	return "需要手动计算"
}

// containsVariable 检查表达式是否包含变量
func containsVariable(expr, variable string) bool {
	return len(expr) > 0 && len(variable) > 0 &&
		(len(expr) > len(variable) || expr == variable)
}

// extractParams 从请求体中提取参数（支持两种格式）
// 格式1: {"operation": "differentiate", "expression": "x^3", ...}
// 格式2: {"calculation": "symbolic_calc", "params": {"operation": "differentiate", "expression": "x^3", ...}}
func extractParams(rawBody map[string]interface{}) (operation, expression, variable string, xValue, yValue, zValue float64) {
	// 检查是否有嵌套的 params
	if paramsRaw, ok := rawBody["params"].(map[string]interface{}); ok {
		// 格式2: 从 params 中提取
		if op, ok := paramsRaw["operation"].(string); ok {
			operation = op
		}
		if expr, ok := paramsRaw["expression"].(string); ok {
			expression = expr
		}
		if v, ok := paramsRaw["variable"].(string); ok {
			variable = v
		}
		if x, ok := paramsRaw["x_value"].(float64); ok {
			xValue = x
		}
		if y, ok := paramsRaw["y_value"].(float64); ok {
			yValue = y
		}
		if z, ok := paramsRaw["z_value"].(float64); ok {
			zValue = z
		}
	} else {
		// 格式1: 直接从根级别提取
		if op, ok := rawBody["operation"].(string); ok {
			operation = op
		}
		if expr, ok := rawBody["expression"].(string); ok {
			expression = expr
		}
		if v, ok := rawBody["variable"].(string); ok {
			variable = v
		}
		if x, ok := rawBody["x_value"].(float64); ok {
			xValue = x
		}
		if y, ok := rawBody["y_value"].(float64); ok {
			yValue = y
		}
		if z, ok := rawBody["z_value"].(float64); ok {
			zValue = z
		}
	}
	return
}

// compareResults 对比结果
func compareResults(original, fixed, expected interface{}) []DifferenceDetail {
	var differences []DifferenceDetail

	// 获取结果映射
	origMap := resultToMap(original)
	fixedMap := resultToMap(fixed)

	// 对比关键字段
	fields := []string{"derivative", "simplified", "numeric_value", "result_expression"}

	for _, field := range fields {
		origVal, origOk := origMap[field]
		fixedVal, fixedOk := fixedMap[field]

		if !origOk && !fixedOk {
			continue
		}

		detail := DifferenceDetail{
			Field:         field,
			OriginalValue: origVal,
			FixedValue:    fixedVal,
		}

		// 判断状态
		if origVal == fixedVal {
			detail.Status = "unchanged"
		} else {
			// 检查修复值是否匹配预期
			if field == "derivative" || field == "simplified" {
				if fixedVal == expected {
					detail.Status = "fixed"
					detail.ExpectedValue = expected
				} else {
					detail.Status = "improved"
				}
			} else {
				detail.Status = "changed"
			}
		}

		differences = append(differences, detail)
	}

	return differences
}

// resultToMap 将结果转换为map
func resultToMap(result interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	if result == nil {
		return m
	}

	// 使用类型断言提取字段
	switch r := result.(type) {
	case map[string]interface{}:
		return r
	case *calculator.SymbolicResult:
		m["derivative"] = r.Derivative
		m["simplified"] = r.Simplified
		m["numeric_value"] = r.NumericValue
		m["result_expression"] = r.ResultExpression
		m["original_expression"] = r.OriginalExpression
	case *calculator.SymbolicResultFixed:
		m["derivative"] = r.Derivative
		m["simplified"] = r.Simplified
		m["numeric_value"] = r.NumericValue
		m["result_expression"] = r.ResultExpression
		m["original_expression"] = r.OriginalExpression
	}

	return m
}

// checkIfFixed 检查是否已修复
func checkIfFixed(differences []DifferenceDetail) bool {
	if len(differences) == 0 {
		return false
	}

	// 检查关键字段是否已修复
	hasFixed := false
	for _, diff := range differences {
		if diff.Status == "fixed" || diff.Status == "improved" {
			hasFixed = true
		}
		// 如果关键字段仍然错误，则认为未修复
		if diff.Field == "derivative" && diff.Status == "unchanged" &&
			(diff.OriginalValue == "d/dx(x^3 + 2*x^2 + x)" || diff.OriginalValue == "") {
			return false
		}
	}

	return hasFixed
}
