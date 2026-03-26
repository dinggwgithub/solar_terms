package api

import (
	"encoding/json"
	"fmt"
	"math"
	"scientific_calc/internal/calculator"
	"time"

	"github.com/gin-gonic/gin"
)

// CompareRequest 对比请求
type CompareRequest struct {
	Calculation string                 `json:"calculation" example:"ode_solver"`
	Params      map[string]interface{} `json:"params"`
}

type CompareResult struct {
	Original    OriginalResult `json:"original"`
	Fixed       FixedResult    `json:"fixed"`
	Differences Differences    `json:"differences"`
}

type OriginalResult struct {
	Solution       float64   `json:"solution"`
	TimePoints     []float64 `json:"time_points"`
	SolutionPath   []float64 `json:"solution_path"`
	DerivativePath []float64 `json:"derivative_path"`
	MethodUsed     string    `json:"method_used"`
	Stability      string    `json:"stability"`
	ErrorEstimate  float64   `json:"error_estimate"`
}

type FixedResult struct {
	Solution       float64   `json:"solution"`
	TimePoints     []float64 `json:"time_points"`
	SolutionPath   []float64 `json:"solution_path"`
	DerivativePath []float64 `json:"derivative_path"`
	MethodUsed     string    `json:"method_used"`
	Stability      string    `json:"stability"`
	ErrorEstimate  float64   `json:"error_estimate"`
	ExactSolution  float64   `json:"exact_solution"`
	AbsoluteError  float64   `json:"absolute_error"`
}

type Differences struct {
	TimePointsUniform         bool     `json:"time_points_uniform"`
	DerivativePathLength      int      `json:"derivative_path_length_original"`
	DerivativePathLengthFixed int      `json:"derivative_path_length_fixed"`
	SolutionDiff              float64  `json:"solution_diff"`
	TimePointsIssues          []string `json:"time_points_issues"`
	ErrorEstimateImproved     bool     `json:"error_estimate_improved"`
}

// CalculateFixed 修复版微分方程求解接口
// @Summary 修复版微分方程求解
// @Description 使用修复后的微分方程求解器计算，解决时间点不均匀、导数路径缺失、误差估计失真等问题
// @Tags 微分方程求解
// @Accept json
// @Produce json
// @Param request body CompareRequest true "计算请求参数"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-fixed [post]
func (h *APIHandler) CalculateFixed(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.Calculation != "ode_solver" {
		h.sendError(c, 400, "此接口仅支持 ode_solver 计算类型")
		return
	}

	calcType := calculator.CalculationTypeODESolverFixed

	sessionID := GenerateSessionID()

	result, warnings, err := h.calculatorManager.CalculateWithSession(calcType, req.Params, sessionID)
	if err != nil {
		h.sendError(c, 400, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, "ode_solver_fixed", sessionID)
}

// CompareSolvers 对比接口
// @Summary 对比原版与修复版微分方程求解器
// @Description 对比原版和修复版微分方程求解器的计算结果，返回差异分析
// @Tags 微分方程求解
// @Accept json
// @Produce json
// @Param request body CompareRequest true "计算请求参数"
// @Success 200 {object} map[string]interface{} "返回对比结果"
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) CompareSolvers(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.Calculation != "ode_solver" {
		h.sendError(c, 400, "此接口仅支持 ode_solver 计算类型")
		return
	}

	sessionID := GenerateSessionID()

	originalResult, _, err := h.calculatorManager.CalculateWithSession(
		calculator.CalculationTypeODESolver, req.Params, sessionID+"_orig")
	if err != nil {
		h.sendError(c, 400, "原始计算失败: "+err.Error())
		return
	}

	fixedResult, _, err := h.calculatorManager.CalculateWithSession(
		calculator.CalculationTypeODESolverFixed, req.Params, sessionID+"_fixed")
	if err != nil {
		h.sendError(c, 400, "修复版计算失败: "+err.Error())
		return
	}

	origMap, err := structToMap(originalResult)
	if err != nil {
		h.sendError(c, 500, "原始结果格式错误: "+err.Error())
		return
	}

	fixedMap, err := structToMap(fixedResult)
	if err != nil {
		h.sendError(c, 500, "修复版结果格式错误: "+err.Error())
		return
	}

	origResult := extractOriginalResult(origMap)
	fixResult := extractFixedResult(fixedMap)

	differences := compareResults(origResult, fixResult, req.Params)

	response := CompareResult{
		Original:    origResult,
		Fixed:       fixResult,
		Differences: differences,
	}

	c.JSON(200, map[string]interface{}{
		"success":   true,
		"result":    response,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func extractOriginalResult(m map[string]interface{}) OriginalResult {
	result := OriginalResult{}
	if v, ok := m["solution"].(float64); ok {
		result.Solution = v
	}
	if v, ok := m["time_points"].([]interface{}); ok {
		result.TimePoints = make([]float64, len(v))
		for i, val := range v {
			if f, ok := val.(float64); ok {
				result.TimePoints[i] = f
			}
		}
	}
	if v, ok := m["solution_path"].([]interface{}); ok {
		result.SolutionPath = make([]float64, len(v))
		for i, val := range v {
			if f, ok := val.(float64); ok {
				result.SolutionPath[i] = f
			}
		}
	}
	if v, ok := m["derivative_path"].([]interface{}); ok {
		result.DerivativePath = make([]float64, len(v))
		for i, val := range v {
			if f, ok := val.(float64); ok {
				result.DerivativePath[i] = f
			}
		}
	}
	if v, ok := m["method_used"].(string); ok {
		result.MethodUsed = v
	}
	if v, ok := m["stability"].(string); ok {
		result.Stability = v
	}
	if v, ok := m["error_estimate"].(float64); ok {
		result.ErrorEstimate = v
	}
	return result
}

func extractFixedResult(m map[string]interface{}) FixedResult {
	result := FixedResult{}
	if v, ok := m["solution"].(float64); ok {
		result.Solution = v
	}
	if v, ok := m["time_points"].([]interface{}); ok {
		result.TimePoints = make([]float64, len(v))
		for i, val := range v {
			if f, ok := val.(float64); ok {
				result.TimePoints[i] = f
			}
		}
	}
	if v, ok := m["solution_path"].([]interface{}); ok {
		result.SolutionPath = make([]float64, len(v))
		for i, val := range v {
			if f, ok := val.(float64); ok {
				result.SolutionPath[i] = f
			}
		}
	}
	if v, ok := m["derivative_path"].([]interface{}); ok {
		result.DerivativePath = make([]float64, len(v))
		for i, val := range v {
			if f, ok := val.(float64); ok {
				result.DerivativePath[i] = f
			}
		}
	}
	if v, ok := m["method_used"].(string); ok {
		result.MethodUsed = v
	}
	if v, ok := m["stability"].(string); ok {
		result.Stability = v
	}
	if v, ok := m["error_estimate"].(float64); ok {
		result.ErrorEstimate = v
	}
	if v, ok := m["exact_solution"].(float64); ok {
		result.ExactSolution = v
	}
	if v, ok := m["absolute_error"].(float64); ok {
		result.AbsoluteError = v
	}
	return result
}

func compareResults(orig OriginalResult, fix FixedResult, params map[string]interface{}) Differences {
	diff := Differences{}

	timeStep, _ := params["time_step"].(float64)
	timeRange, _ := params["time_range"].(float64)

	diff.TimePointsUniform = checkTimePointsUniform(fix.TimePoints, timeStep)
	diff.DerivativePathLength = len(orig.DerivativePath)
	diff.DerivativePathLengthFixed = len(fix.DerivativePath)
	diff.SolutionDiff = math.Abs(orig.Solution - fix.Solution)

	diff.TimePointsIssues = analyzeTimePointsIssues(orig.TimePoints, timeStep, timeRange)

	diff.ErrorEstimateImproved = fix.ErrorEstimate <= orig.ErrorEstimate

	return diff
}

func checkTimePointsUniform(points []float64, expectedStep float64) bool {
	if len(points) < 2 {
		return true
	}

	for i := 1; i < len(points); i++ {
		actualStep := points[i] - points[i-1]
		if math.Abs(actualStep-expectedStep) > 1e-9 {
			if i != len(points)-1 {
				return false
			}
		}
	}
	return true
}

func analyzeTimePointsIssues(points []float64, expectedStep, timeRange float64) []string {
	var issues []string

	if len(points) < 2 {
		return issues
	}

	for i := 1; i < len(points); i++ {
		actualStep := points[i] - points[i-1]
		deviation := math.Abs(actualStep - expectedStep)
		if deviation > 1e-9 && i != len(points)-1 {
			issues = append(issues,
				fmt.Sprintf("时间点 %.2f 到 %.2f 步长异常: %.6f (预期: %.6f)",
					points[i-1], points[i], actualStep, expectedStep))
		}
	}

	expectedPoints := int(timeRange/expectedStep) + 1
	if len(points) != expectedPoints {
		issues = append(issues,
			fmt.Sprintf("时间点数量错误: %d (预期: %d)", len(points), expectedPoints))
	}

	return issues
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
