package api

import (
	"net/http"
	"scientific_calc/internal/calculator"
	"time"

	"github.com/gin-gonic/gin"
)

// ODEFixedResponse 修复版ODE求解响应
type ODEFixedResponse struct {
	Success     bool                           `json:"success"`
	Result      *calculator.ODEResultFixed     `json:"result"`
	Warnings    []string                       `json:"warnings"`
	Timestamp   string                         `json:"timestamp"`
	Calculation string                         `json:"calculation"`
	SessionID   string                         `json:"session_id,omitempty"`
	Version     string                         `json:"version"` // 标识修复版本
}

// ODECompareRequest 对比请求
type ODECompareRequest struct {
	OriginalRequest  interface{} `json:"original_request" binding:"required"`
	FixedRequest     interface{} `json:"fixed_request" binding:"required"`
}

// ODECompareResponse 对比响应
type ODECompareResponse struct {
	Success         bool                      `json:"success"`
	OriginalResult  *calculator.ODEResult     `json:"original_result"`
	FixedResult     *calculator.ODEResultFixed `json:"fixed_result"`
	Differences     *ODEDifferences           `json:"differences"`
	Timestamp       string                    `json:"timestamp"`
}

// ODEDifferences 差异分析
type ODEDifferences struct {
	TimePointsMatch     bool      `json:"time_points_match"`
	SolutionPathMatch   bool      `json:"solution_path_match"`
	DerivativePathMatch bool      `json:"derivative_path_match"`
	TimePointsDiff      []float64 `json:"time_points_diff,omitempty"`
	SolutionPathDiff    []float64 `json:"solution_path_diff,omitempty"`
	FinalValueDiff      float64   `json:"final_value_diff"`
	ErrorEstimateDiff   float64   `json:"error_estimate_diff"`
	Issues              []string  `json:"issues"`
}

// CalculateFixed 修复版ODE求解接口
// @Summary 修复版微分方程求解
// @Description 使用修复后的算法求解微分方程，确保时间点均匀、导数路径完整、误差估计准确
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "ODE求解参数"
// @Success 200 {object} ODEFixedResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-fixed [post]
func (h *APIHandler) CalculateFixed(c *gin.Context) {
	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 使用修复后的求解器
	result, err := calculator.SolveODE(params)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	// 验证结果完整性
	warnings := validateODEResult(result)

	response := ODEFixedResponse{
		Success:     true,
		Result:      result,
		Warnings:    warnings,
		Timestamp:   time.Now().Format(time.RFC3339),
		Calculation: "ode_solver_fixed",
		Version:     "2.0-fixed",
	}

	c.JSON(http.StatusOK, response)
}

// CompareSolvers 求解器对比接口
// @Summary 对比原始和修复后的ODE求解器
// @Description 同时运行原始和修复后的求解器，返回详细对比结果
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body ODECompareRequest true "对比请求参数"
// @Success 200 {object} ODECompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) CompareSolvers(c *gin.Context) {
	var req ODECompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 运行原始求解器
	originalResult, err := h.runOriginalSolver(req.OriginalRequest)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原始求解器运行失败: "+err.Error())
		return
	}

	// 运行修复后求解器
	fixedResult, err := calculator.SolveODE(req.FixedRequest)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复求解器运行失败: "+err.Error())
		return
	}

	// 计算差异
	differences := calculateDifferences(originalResult, fixedResult)

	response := ODECompareResponse{
		Success:        true,
		OriginalResult: originalResult,
		FixedResult:    fixedResult,
		Differences:    differences,
		Timestamp:      time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// runOriginalSolver 运行原始求解器
func (h *APIHandler) runOriginalSolver(params interface{}) (*calculator.ODEResult, error) {
	// 创建原始求解器实例
	calc := calculator.NewODESolverCalculator()
	result, err := calc.Calculate(params)
	if err != nil {
		return nil, err
	}

	// 类型断言
	odeResult, ok := result.(*calculator.ODEResult)
	if !ok {
		return nil, gin.Error{Err: nil, Type: gin.ErrorTypePrivate}
	}

	return odeResult, nil
}

// validateODEResult 验证ODE结果完整性
func validateODEResult(result *calculator.ODEResultFixed) []string {
	var warnings []string

	// 检查时间点数量
	if len(result.TimePoints) == 0 {
		warnings = append(warnings, "时间点序列为空")
	}

	// 检查解路径数量
	if len(result.SolutionPath) == 0 {
		warnings = append(warnings, "解路径为空")
	}

	// 检查导数路径数量
	if len(result.DerivativePath) == 0 {
		warnings = append(warnings, "导数路径为空")
	}

	// 检查数组长度一致性
	if len(result.TimePoints) != len(result.SolutionPath) ||
		len(result.TimePoints) != len(result.DerivativePath) {
		warnings = append(warnings, "时间序列、解路径和导数路径长度不一致")
	}

	// 检查时间点均匀性
	if len(result.TimePoints) >= 2 {
		expectedStep := result.TimePoints[1] - result.TimePoints[0]
		for i := 2; i < len(result.TimePoints); i++ {
			actualStep := result.TimePoints[i] - result.TimePoints[i-1]
			if abs(actualStep-expectedStep) > 1e-10 {
				warnings = append(warnings, "时间点步长不均匀")
				break
			}
		}
	}

	return warnings
}

// calculateDifferences 计算两个结果的差异
func calculateDifferences(original *calculator.ODEResult, fixed *calculator.ODEResultFixed) *ODEDifferences {
	diff := &ODEDifferences{
		Issues: []string{},
	}

	// 检查时间点
	if len(original.TimePoints) == len(fixed.TimePoints) {
		diff.TimePointsMatch = true
		for i := range original.TimePoints {
			if abs(original.TimePoints[i]-fixed.TimePoints[i]) > 1e-10 {
				diff.TimePointsMatch = false
				diff.TimePointsDiff = append(diff.TimePointsDiff, original.TimePoints[i]-fixed.TimePoints[i])
			}
		}
	} else {
		diff.TimePointsMatch = false
		diff.Issues = append(diff.Issues, "时间点数量不一致")
	}

	// 检查解路径
	if len(original.SolutionPath) == len(fixed.SolutionPath) {
		diff.SolutionPathMatch = true
		for i := range original.SolutionPath {
			if abs(original.SolutionPath[i]-fixed.SolutionPath[i]) > 1e-10 {
				diff.SolutionPathMatch = false
				diff.SolutionPathDiff = append(diff.SolutionPathDiff, original.SolutionPath[i]-fixed.SolutionPath[i])
			}
		}
	} else {
		diff.SolutionPathMatch = false
		diff.Issues = append(diff.Issues, "解路径长度不一致")
	}

	// 检查导数路径
	if len(original.DerivativePath) != len(fixed.DerivativePath) {
		diff.DerivativePathMatch = false
		diff.Issues = append(diff.Issues, "导数路径长度不一致")
	} else {
		diff.DerivativePathMatch = true
	}

	// 计算最终值差异
	diff.FinalValueDiff = original.Solution - fixed.Solution

	// 计算误差估计差异
	diff.ErrorEstimateDiff = original.ErrorEstimate - fixed.ErrorEstimate

	// 分析原始结果的问题
	if len(original.DerivativePath) <= 1 {
		diff.Issues = append(diff.Issues, "原始求解器导数路径缺失或过少")
	}

	// 检查原始时间点均匀性
	if len(original.TimePoints) >= 2 {
		expectedStep := 0.1 // 假设标准步长
		for i := 1; i < len(original.TimePoints); i++ {
			actualStep := original.TimePoints[i] - original.TimePoints[i-1]
			if abs(actualStep-expectedStep) > 0.001 {
				diff.Issues = append(diff.Issues, "原始求解器时间点步长不均匀")
				break
			}
		}
	}

	return diff
}

// abs 绝对值辅助函数
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
