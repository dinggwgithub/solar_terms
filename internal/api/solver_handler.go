package api

import (
	"fmt"
	"math"
	"net/http"
	"scientific_calc/internal/calculator"
	"time"

	"github.com/gin-gonic/gin"
)

// SolverHandler 方程求解处理器
type SolverHandler struct {
	equationSolver *calculator.EquationSolverCalculator
}

// NewSolverHandler 创建新的方程求解处理器
func NewSolverHandler() *SolverHandler {
	return &SolverHandler{
		equationSolver: calculator.NewEquationSolverCalculator(),
	}
}

// SolverV2Request 求解器V2请求
type SolverV2Request struct {
	EquationType  string      `json:"equation_type"`  // 方程类型：nonlinear, linear, ode
	Equation      string      `json:"equation"`       // 方程表达式
	InitialGuess  float64     `json:"initial_guess"`  // 初始猜测值（非线性方程）
	Tolerance     float64     `json:"tolerance"`      // 容差
	MaxIterations int         `json:"max_iterations"` // 最大迭代次数
	Coefficients  []float64   `json:"coefficients"`   // 系数（线性方程组）
	TimeStep      float64     `json:"time_step"`      // 时间步长（微分方程）
	TimeRange     float64     `json:"time_range"`     // 时间范围（微分方程）
	Params        interface{} `json:"params"`         // 兼容旧格式：嵌套参数
}

// SolverV2DirectRequest 直接格式请求（无嵌套）
type SolverV2DirectRequest struct {
	EquationType  string    `json:"equation_type" binding:"required"`
	Equation      string    `json:"equation" binding:"required"`
	InitialGuess  float64   `json:"initial_guess"`
	Tolerance     float64   `json:"tolerance"`
	MaxIterations int       `json:"max_iterations"`
	Coefficients  []float64 `json:"coefficients"`
	TimeStep      float64   `json:"time_step"`
	TimeRange     float64   `json:"time_range"`
}

// SolverV2Response 求解器V2响应
type SolverV2Response struct {
	Success   bool                      `json:"success"`
	Result    calculator.EquationResult `json:"result"`
	Algorithm string                    `json:"algorithm"` // 使用的算法版本
	Timestamp string                    `json:"timestamp"`
	SessionID string                    `json:"session_id"`
}

// SolverCompareResponse 求解器对比响应
type SolverCompareResponse struct {
	Success    bool                      `json:"success"`
	V1Result   calculator.EquationResult `json:"v1_result"`   // V1版本结果
	V2Result   calculator.EquationResult `json:"v2_result"`   // V2版本结果
	DiffReport SolverDiffReport          `json:"diff_report"` // 差异报告
	Timestamp  string                    `json:"timestamp"`
}

// SolverDiffReport 求解器差异报告
type SolverDiffReport struct {
	SolutionDiff     float64 `json:"solution_diff"`     // 解的差异
	IterationsDiff   int     `json:"iterations_diff"`   // 迭代次数差异
	ConvergedChanged bool    `json:"converged_changed"` // 收敛状态是否改变
	ErrorDiff        float64 `json:"error_diff"`        // 误差差异
	ResidualDiff     float64 `json:"residual_diff"`     // 残差差异
	Analysis         string  `json:"analysis"`          // 差异分析
}

// SolveV2 修复版方程求解接口
// @Summary 修复版方程求解
// @Description 使用修复后的收敛判定逻辑求解方程，支持两种参数格式：直接格式或嵌套params格式
// @Tags 方程求解
// @Accept json
// @Produce json
// @Param request body SolverV2Request true "求解请求参数"
// @Success 200 {object} SolverV2Response
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/v2 [post]
func (h *SolverHandler) SolveV2(c *gin.Context) {
	var req SolverV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// 解析参数（支持嵌套params格式或直接格式）
	params, err := h.parseSolverRequest(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "参数解析错误: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// 执行求解（使用V2版本）
	result, err := h.equationSolver.SolveNonlinearEquationV2(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "求解失败: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, SolverV2Response{
		Success:   true,
		Result:    *result,
		Algorithm: "Newton-Raphson V2 (Fixed)",
		Timestamp: generateTimestamp(),
		SessionID: GenerateSessionID(),
	})
}

// parseSolverRequest 解析求解请求参数（支持嵌套params格式）
func (h *SolverHandler) parseSolverRequest(req *SolverV2Request) (*calculator.EquationParams, error) {
	// 如果提供了嵌套的params，优先使用
	if req.Params != nil {
		paramsMap, ok := req.Params.(map[string]interface{})
		if ok {
			return h.parseParamsFromMap(paramsMap)
		}
	}

	// 否则使用直接格式的参数
	if req.EquationType == "" || req.Equation == "" {
		return nil, fmt.Errorf("缺少必需的参数: equation_type 和 equation")
	}

	// 设置默认值
	tolerance := req.Tolerance
	if tolerance == 0 {
		tolerance = 1e-6
	}
	maxIterations := req.MaxIterations
	if maxIterations == 0 {
		maxIterations = 100
	}
	timeStep := req.TimeStep
	if timeStep == 0 {
		timeStep = 0.1
	}
	timeRange := req.TimeRange
	if timeRange == 0 {
		timeRange = 10.0
	}

	return &calculator.EquationParams{
		EquationType:  req.EquationType,
		Equation:      req.Equation,
		InitialGuess:  req.InitialGuess,
		Tolerance:     tolerance,
		MaxIterations: maxIterations,
		Coefficients:  req.Coefficients,
		TimeStep:      timeStep,
		TimeRange:     timeRange,
	}, nil
}

// parseParamsFromMap 从map解析参数
func (h *SolverHandler) parseParamsFromMap(paramsMap map[string]interface{}) (*calculator.EquationParams, error) {
	// 提取必需参数
	equationType, ok := paramsMap["equation_type"].(string)
	if !ok || equationType == "" {
		return nil, fmt.Errorf("equation_type参数必须为字符串且不能为空")
	}

	equation, ok := paramsMap["equation"].(string)
	if !ok || equation == "" {
		return nil, fmt.Errorf("equation参数必须为字符串且不能为空")
	}

	// 设置默认值
	params := &calculator.EquationParams{
		EquationType:  equationType,
		Equation:      equation,
		Tolerance:     1e-6,
		MaxIterations: 100,
		TimeStep:      0.1,
		TimeRange:     10.0,
	}

	// 提取可选参数
	if initialGuess, ok := paramsMap["initial_guess"].(float64); ok {
		params.InitialGuess = initialGuess
	}
	if tolerance, ok := paramsMap["tolerance"].(float64); ok && tolerance > 0 {
		params.Tolerance = tolerance
	}
	if maxIterations, ok := paramsMap["max_iterations"].(float64); ok && maxIterations > 0 {
		params.MaxIterations = int(maxIterations)
	}
	if timeStep, ok := paramsMap["time_step"].(float64); ok && timeStep > 0 {
		params.TimeStep = timeStep
	}
	if timeRange, ok := paramsMap["time_range"].(float64); ok && timeRange > 0 {
		params.TimeRange = timeRange
	}

	// 处理系数数组
	if coefficients, ok := paramsMap["coefficients"].([]interface{}); ok {
		params.Coefficients = make([]float64, len(coefficients))
		for i, coef := range coefficients {
			if floatCoef, ok := coef.(float64); ok {
				params.Coefficients[i] = floatCoef
			}
		}
	}

	return params, nil
}

// CompareSolvers 求解器对比接口
// @Summary 对比V1和V2版本的求解器
// @Description 同时运行原始版本和修复版本，输出结构化差异报告，支持两种参数格式
// @Tags 方程求解
// @Accept json
// @Produce json
// @Param request body SolverV2Request true "求解请求参数"
// @Success 200 {object} SolverCompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *SolverHandler) CompareSolvers(c *gin.Context) {
	var req SolverV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// 解析参数（支持嵌套params格式或直接格式）
	params, err := h.parseSolverRequest(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "参数解析错误: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// 运行V1版本
	v1Result, v1Err := h.equationSolver.SolveNonlinearEquationV1(params)
	if v1Err != nil {
		v1Result = &calculator.EquationResult{
			Solution:      req.InitialGuess,
			Iterations:    0,
			Converged:     false,
			Error:         0,
			FunctionValue: 0,
		}
	}

	// 运行V2版本
	v2Result, v2Err := h.equationSolver.SolveNonlinearEquationV2(params)
	if v2Err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "V2求解失败: " + v2Err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// 生成差异报告
	diffReport := generateDiffReport(v1Result, v2Result, params)

	c.JSON(http.StatusOK, SolverCompareResponse{
		Success:    true,
		V1Result:   *v1Result,
		V2Result:   *v2Result,
		DiffReport: diffReport,
		Timestamp:  generateTimestamp(),
	})
}

// generateDiffReport 生成差异报告
func generateDiffReport(v1, v2 *calculator.EquationResult, params *calculator.EquationParams) SolverDiffReport {
	var v1Solution, v2Solution float64

	// 处理solution类型
	switch s := v1.Solution.(type) {
	case float64:
		v1Solution = s
	case float32:
		v1Solution = float64(s)
	case int:
		v1Solution = float64(s)
	}

	switch s := v2.Solution.(type) {
	case float64:
		v2Solution = s
	case float32:
		v2Solution = float64(s)
	case int:
		v2Solution = float64(s)
	}

	solutionDiff := math.Abs(v2Solution - v1Solution)
	iterationsDiff := v2.Iterations - v1.Iterations
	convergedChanged := v1.Converged != v2.Converged
	errorDiff := math.Abs(v2.Error - v1.Error)
	residualDiff := math.Abs(v2.FunctionValue - v1.FunctionValue)

	// 生成分析文本
	analysis := generateAnalysis(v1, v2, solutionDiff, iterationsDiff, convergedChanged, params)

	return SolverDiffReport{
		SolutionDiff:     solutionDiff,
		IterationsDiff:   iterationsDiff,
		ConvergedChanged: convergedChanged,
		ErrorDiff:        errorDiff,
		ResidualDiff:     residualDiff,
		Analysis:         analysis,
	}
}

// generateTimestamp 生成时间戳
func generateTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// generateAnalysis 生成差异分析文本
func generateAnalysis(v1, v2 *calculator.EquationResult, solutionDiff float64, iterationsDiff int, convergedChanged bool, params *calculator.EquationParams) string {
	analysis := "差异分析:\n"

	// 收敛状态分析
	if convergedChanged {
		if v2.Converged && !v1.Converged {
			analysis += "- V2版本成功收敛，而V1版本未收敛\n"
		} else if !v2.Converged && v1.Converged {
			analysis += "- V1版本收敛但V2版本未收敛（可能V2更严格）\n"
		}
	} else {
		if v1.Converged && v2.Converged {
			analysis += "- 两个版本都成功收敛\n"
		} else {
			analysis += "- 两个版本都未收敛\n"
		}
	}

	// 迭代次数分析
	if iterationsDiff < 0 {
		analysis += fmt.Sprintf("- V2版本迭代次数减少%d次，收敛更快\n", -iterationsDiff)
	} else if iterationsDiff > 0 {
		analysis += fmt.Sprintf("- V2版本迭代次数增加%d次，可能更精确\n", iterationsDiff)
	} else {
		analysis += "- 两个版本迭代次数相同\n"
	}

	// 解的精度分析
	if solutionDiff < params.Tolerance {
		analysis += "- 两个版本的解在容差范围内一致\n"
	} else {
		analysis += fmt.Sprintf("- 解的差异为%.2e，超出容差范围\n", solutionDiff)
	}

	// 残差分析
	if math.Abs(v2.FunctionValue) < math.Abs(v1.FunctionValue) {
		analysis += "- V2版本的残差更小，精度更高\n"
	} else if math.Abs(v2.FunctionValue) > math.Abs(v1.FunctionValue) {
		analysis += "- V1版本的残差更小\n"
	} else {
		analysis += "- 两个版本的残差相同\n"
	}

	// 总体评价
	if v2.Converged && (v2.Iterations <= v1.Iterations || math.Abs(v2.FunctionValue) <= math.Abs(v1.FunctionValue)) {
		analysis += "结论: V2版本修复有效，收敛判定更准确"
	} else if v2.Converged {
		analysis += "结论: V2版本修复有效，收敛更稳定"
	} else {
		analysis += "结论: 需要进一步检查"
	}

	return analysis
}
