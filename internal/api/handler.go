package api

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"scientific_calc/internal/calculator"
	"scientific_calc/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateSessionID 生成会话ID
func GenerateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// APIHandler API处理器
type APIHandler struct {
	calculatorManager *calculator.CalculatorManager
}

// NewAPIHandler 创建新的API处理器
func NewAPIHandler(calculatorManager *calculator.CalculatorManager) *APIHandler {
	return &APIHandler{
		calculatorManager: calculatorManager,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Time      string `json:"time" example:"2024-01-01T12:00:00Z"`
	Version   string `json:"version" example:"1.0.0"`
	BuildTime string `json:"build_time" example:"2024-01-01T12:00:00Z"`
}

// CalculationResponse 计算响应
type CalculationResponse struct {
	Success     bool        `json:"success"`
	Result      interface{} `json:"result"`
	Warnings    []string    `json:"warnings"`
	Timestamp   string      `json:"timestamp"`
	Calculation string      `json:"calculation"`
	SessionID   string      `json:"session_id,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

// SolverComparisonResponse 求解器对比响应
type SolverComparisonResponse struct {
	Success          bool                    `json:"success"`
	OldResult        interface{}             `json:"old_result"`
	NewResult        interface{}             `json:"new_result"`
	ComparisonReport *SolverComparisonReport `json:"comparison_report"`
	Timestamp        string                  `json:"timestamp"`
	Calculation      string                  `json:"calculation"`
	SessionID        string                  `json:"session_id,omitempty"`
}

// SolverComparisonReport 求解器对比报告
type SolverComparisonReport struct {
	TestEquation     string      `json:"test_equation"`
	TimeStep         float64     `json:"time_step"`
	TimeRange        float64     `json:"time_range"`
	InitialValue     float64     `json:"initial_value"`
	FinalValueOld    float64     `json:"final_value_old"`
	FinalValueNew    float64     `json:"final_value_new"`
	FinalValueTheory float64     `json:"final_value_theory"`
	ErrorOldVsTheory float64     `json:"error_old_vs_theory"`
	ErrorNewVsTheory float64     `json:"error_new_vs_theory"`
	ErrorNewVsOld    float64     `json:"error_new_vs_old"`
	ImprovementRatio float64     `json:"improvement_ratio"`
	MaxErrorOld      float64     `json:"max_error_old"`
	MaxErrorNew      float64     `json:"max_error_new"`
	MeanErrorOld     float64     `json:"mean_error_old"`
	MeanErrorNew     float64     `json:"mean_error_new"`
	ConvergedOld     bool        `json:"converged_old"`
	ConvergedNew     bool        `json:"converged_new"`
	IterationsOld    int         `json:"iterations_old"`
	IterationsNew    int         `json:"iterations_new"`
	TimePointsMatch  bool        `json:"time_points_match"`
	IssueSummary     []string    `json:"issue_summary"`
	Recommendations  []string    `json:"recommendations"`
	FieldDifferences []FieldDiff `json:"field_differences,omitempty"`
}

// FieldDiff 字段差异
type FieldDiff struct {
	Field       string      `json:"field"`
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	Difference  float64     `json:"difference,omitempty"`
	Description string      `json:"description"`
}

// CalculatorInfoResponse 计算器信息响应
type CalculatorInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SystemInfoResponse 系统信息响应
type SystemInfoResponse struct {
	Version               string   `json:"version"`
	SupportedCalculations []string `json:"supported_calculations"`
	TotalCalculators      int      `json:"total_calculators"`
}

// CalculateFixed 修复版科学计算接口
// @Summary 执行修复版科学计算
// @Description 执行各种类型的科学计算（修复版）
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID，用于保持计算参数一致性，不传则自动生成"
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-fixed [post]
func (h *APIHandler) CalculateFixed(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 解析计算类型 - 如果是equation_solver则使用修复版
	calcTypeStr := req.Calculation
	if calcTypeStr == "equation_solver" {
		calcTypeStr = "equation_solver_fixed"
	}

	calcType, err := calculator.ParseCalculationType(calcTypeStr)
	if err != nil {
		// 尝试使用原始类型
		calcType, err = calculator.ParseCalculationType(req.Calculation)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
			return
		}
	}

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 执行计算
	result, warnings, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, req.Calculation, sessionID)
}

// SolverCompare 求解器对比接口
// @Summary 对比新旧求解器
// @Description 对比新旧方程求解器的计算结果差异
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} SolverComparisonResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) SolverCompare(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 仅支持equation_solver的对比
	if req.Calculation != "equation_solver" {
		h.sendError(c, http.StatusBadRequest, "仅支持equation_solver类型的对比")
		return
	}

	params := req.GetParams()

	// 转换为map[string]interface{}
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		h.sendError(c, http.StatusBadRequest, "参数格式错误")
		return
	}

	// 使用旧求解器计算
	oldCalcType := calculator.CalculationTypeEquationSolver
	oldResult, _, errOld := h.calculatorManager.Calculate(oldCalcType, paramsMap)
	if errOld != nil {
		h.sendError(c, http.StatusBadRequest, "旧求解器计算失败: "+errOld.Error())
		return
	}

	// 使用新求解器计算
	newCalcType := calculator.CalculationTypeEquationSolverFixed
	newResult, _, errNew := h.calculatorManager.Calculate(newCalcType, paramsMap)
	if errNew != nil {
		h.sendError(c, http.StatusBadRequest, "新求解器计算失败: "+errNew.Error())
		return
	}

	// 生成对比报告
	comparisonReport := h.generateComparisonReport(paramsMap, oldResult, newResult)

	sessionID := GenerateSessionID()

	response := SolverComparisonResponse{
		Success:          true,
		OldResult:        oldResult,
		NewResult:        newResult,
		ComparisonReport: comparisonReport,
		Timestamp:        time.Now().Format(time.RFC3339),
		Calculation:      req.Calculation,
		SessionID:        sessionID,
	}

	c.JSON(http.StatusOK, response)
}

// generateComparisonReport 生成对比报告
func (h *APIHandler) generateComparisonReport(params map[string]interface{}, oldResult, newResult interface{}) *SolverComparisonReport {
	report := &SolverComparisonReport{
		IssueSummary:     make([]string, 0),
		Recommendations:  make([]string, 0),
		FieldDifferences: make([]FieldDiff, 0),
	}

	// 提取参数
	if eq, ok := params["equation"].(string); ok {
		report.TestEquation = eq
	}
	if ts, ok := params["time_step"].(float64); ok {
		report.TimeStep = ts
	}
	if tr, ok := params["time_range"].(float64); ok {
		report.TimeRange = tr
	}
	if iv, ok := params["initial_value"].(float64); ok {
		report.InitialValue = iv
	} else if ig, ok := params["initial_guess"].(float64); ok {
		report.InitialValue = ig
	}

	// 计算理论值（针对dy/dt = -y）
	if report.TestEquation == "dy/dt = -y" {
		report.FinalValueTheory = report.InitialValue * math.Exp(-report.TimeRange)
	}

	// 解析旧结果
	oldRes, ok := oldResult.(*calculator.EquationResult)
	if ok {
		report.FinalValueOld = oldRes.Solution.(float64)
		report.ConvergedOld = oldRes.Converged
		report.IterationsOld = oldRes.Iterations

		// 计算旧版误差
		if report.FinalValueTheory != 0 {
			report.ErrorOldVsTheory = math.Abs(report.FinalValueOld - report.FinalValueTheory)
		}

		// 计算旧版路径误差
		if len(oldRes.SolutionPath) > 0 && len(oldRes.TimePoints) > 0 {
			maxErr := 0.0
			meanErr := 0.0
			for i, t := range oldRes.TimePoints {
				if i < len(oldRes.SolutionPath) {
					theo := report.InitialValue * math.Exp(-t)
					err := math.Abs(oldRes.SolutionPath[i] - theo)
					if err > maxErr {
						maxErr = err
					}
					meanErr += err
				}
			}
			report.MaxErrorOld = maxErr
			report.MeanErrorOld = meanErr / float64(len(oldRes.SolutionPath))
		}
	}

	// 解析新结果
	newRes, ok := newResult.(*calculator.FixedEquationResult)
	if ok {
		report.FinalValueNew = newRes.Solution.(float64)
		report.ConvergedNew = newRes.Converged
		report.IterationsNew = newRes.Iterations
		report.MaxErrorNew = newRes.MaxError
		report.MeanErrorNew = newRes.MeanError

		// 计算新版误差
		if report.FinalValueTheory != 0 {
			report.ErrorNewVsTheory = math.Abs(report.FinalValueNew - report.FinalValueTheory)
		}
	}

	// 计算新旧差异
	report.ErrorNewVsOld = math.Abs(report.FinalValueNew - report.FinalValueOld)

	// 计算改进率
	if report.ErrorOldVsTheory > 0 {
		report.ImprovementRatio = (report.ErrorOldVsTheory - report.ErrorNewVsTheory) / report.ErrorOldVsTheory * 100
	}

	// 检查时间点匹配
	if ok && len(oldRes.TimePoints) > 1 {
		expectedStep := report.TimeStep
		actualStep := oldRes.TimePoints[1] - oldRes.TimePoints[0]
		report.TimePointsMatch = math.Abs(expectedStep-actualStep) < 1e-10
		if !report.TimePointsMatch {
			report.IssueSummary = append(report.IssueSummary,
				fmt.Sprintf("时间步长不准确: 期望%.6f, 实际%.6f", expectedStep, actualStep))
		}
	}

	// 生成问题摘要和建议
	if report.ErrorOldVsTheory > 0.01 {
		report.IssueSummary = append(report.IssueSummary,
			fmt.Sprintf("旧版最终值误差过大: %.6f", report.ErrorOldVsTheory))
	}

	if report.ErrorNewVsTheory < 0.0001 {
		report.Recommendations = append(report.Recommendations,
			"新版求解器精度符合要求，建议使用修复版接口")
	} else {
		report.Recommendations = append(report.Recommendations,
			"建议进一步优化求解算法参数")
	}

	// 生成字段差异
	report.FieldDifferences = append(report.FieldDifferences, FieldDiff{
		Field:       "final_solution",
		OldValue:    report.FinalValueOld,
		NewValue:    report.FinalValueNew,
		Difference:  report.ErrorNewVsOld,
		Description: "最终计算值差异",
	})

	report.FieldDifferences = append(report.FieldDifferences, FieldDiff{
		Field:       "error_vs_theory",
		OldValue:    report.ErrorOldVsTheory,
		NewValue:    report.ErrorNewVsTheory,
		Difference:  report.ErrorOldVsTheory - report.ErrorNewVsTheory,
		Description: "与理论值的误差改进量",
	})

	report.FieldDifferences = append(report.FieldDifferences, FieldDiff{
		Field:       "method_used",
		OldValue:    "Euler(modified)",
		NewValue:    "RK4(standard)",
		Description: "数值方法变更",
	})

	return report
}

// HealthCheck 健康检查接口
// @Summary 健康检查
// @Description 健康检查接口
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/health [get]
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Time:      time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		BuildTime: "2024-03-23T12:00:00Z",
	})
}

// Calculate 科学计算接口
// @Summary 执行科学计算
// @Description 执行各种类型的科学计算
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID，用于保持计算参数一致性，不传则自动生成"
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate [post]
func (h *APIHandler) Calculate(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 解析计算类型
	calcType, err := calculator.ParseCalculationType(req.Calculation)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
		return
	}

	// 获取或生成会话ID（用于保持参数一致性）
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 执行计算
	result, warnings, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, req.Calculation, sessionID)
}

// GetCalculatorInfo 获取计算器信息接口
// @Summary 获取计算器信息
// @Description 获取指定计算器的详细信息
// @Tags 计算器管理
// @Accept json
// @Produce json
// @Param calculation query string true "计算类型: solar_term, ganzhi, astronomy, starting_age, lunar, planet, star, sunrise_sunset, moon_phase"
// @Success 200 {object} CalculatorInfoResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculator-info [get]
func (h *APIHandler) GetCalculatorInfo(c *gin.Context) {
	calcTypeStr := c.Query("calculation")
	if calcTypeStr == "" {
		h.sendError(c, http.StatusBadRequest, "缺少calculation参数")
		return
	}

	calcType, err := calculator.ParseCalculationType(calcTypeStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+calcTypeStr)
		return
	}

	calcInfo, err := h.calculatorManager.GetCalculatorInfo(calcType)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "获取计算器信息失败: "+err.Error())
		return
	}

	response := CalculatorInfoResponse{
		Name:        calcInfo["name"],
		Description: calcInfo["description"],
	}

	c.JSON(http.StatusOK, response)
}

// GetSystemInfo 获取系统信息接口
// @Summary 获取系统信息
// @Description 获取系统支持的完整信息
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} SystemInfoResponse
// @Router /api/system-info [get]
func (h *APIHandler) GetSystemInfo(c *gin.Context) {
	// 获取支持的计算类型
	supportedCalculations := h.calculatorManager.GetSupportedCalculationTypes()
	calcTypes := make([]string, 0, len(supportedCalculations))
	for _, calcType := range supportedCalculations {
		calcTypes = append(calcTypes, calcType.String())
	}

	response := SystemInfoResponse{
		Version:               "1.0.0",
		SupportedCalculations: calcTypes,
		TotalCalculators:      len(supportedCalculations),
	}

	c.JSON(http.StatusOK, response)
}

// sendSuccess 发送成功响应
func (h *APIHandler) sendSuccess(c *gin.Context, result interface{}, warnings []string, calculation string, sessionID ...string) {
	response := CalculationResponse{
		Success:     true,
		Result:      result,
		Warnings:    warnings,
		Timestamp:   time.Now().Format(time.RFC3339),
		Calculation: calculation,
	}

	if len(sessionID) > 0 && sessionID[0] != "" {
		response.SessionID = sessionID[0]
	}

	c.JSON(http.StatusOK, response)
}

// sendError 发送错误响应
func (h *APIHandler) sendError(c *gin.Context, code int, message string) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}

	c.JSON(code, response)
}

// RegisterRoutes 注册API路由
func (h *APIHandler) RegisterRoutes(router *gin.Engine) {
	// 系统接口
	router.GET("/api/health", h.HealthCheck)
	router.GET("/api/system-info", h.GetSystemInfo)

	// 科学计算接口
	router.POST("/api/calculate", h.Calculate)
	router.POST("/api/calculate-fixed", h.CalculateFixed)
	router.POST("/api/solver/compare", h.SolverCompare)

	// 计算器管理接口
	router.GET("/api/calculator-info", h.GetCalculatorInfo)

	// 支持接口
	router.GET("/api/supported-calculations", h.GetSupportedCalculations)
}

// GetSupportedCalculations 获取支持的计算类型接口
func (h *APIHandler) GetSupportedCalculations(c *gin.Context) {
	supportedCalculations := h.calculatorManager.GetSupportedCalculationTypes()

	var calculations []map[string]interface{}
	for _, calcType := range supportedCalculations {
		calcInfo, err := h.calculatorManager.GetCalculatorInfo(calcType)
		if err == nil {
			calculations = append(calculations, map[string]interface{}{
				"type":        calcType.String(),
				"name":        calcInfo["name"],
				"description": calcInfo["description"],
			})
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"calculations": calculations,
		"total":        len(calculations),
	})
}
