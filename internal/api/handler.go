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

// CalculatorInfoResponse 计算器信息响应
type CalculatorInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CompareRequest 对比请求
type CompareRequest struct {
	Original models.CalculationRequest `json:"original" binding:"required"`
	Fixed    models.CalculationRequest `json:"fixed" binding:"required"`
}

// CompareResponse 对比响应
type CompareResponse struct {
	Success    bool        `json:"success"`
	Original   interface{} `json:"original_result"`
	Fixed      interface{} `json:"fixed_result"`
	Difference interface{} `json:"difference"`
	Timestamp  string      `json:"timestamp"`
}

// SystemInfoResponse 系统信息响应
type SystemInfoResponse struct {
	Version               string   `json:"version"`
	SupportedCalculations []string `json:"supported_calculations"`
	TotalCalculators      int      `json:"total_calculators"`
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

// CalculateFixed 修复版科学计算接口
// @Summary 执行科学计算（修复版）
// @Description 执行各种类型的科学计算（修复版，用于对比测试）
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

	// 获取计算器实例
	calc, exists := h.calculatorManager.GetCalculator(calcType)
	if !exists {
		h.sendError(c, http.StatusBadRequest, "计算器未注册: "+req.Calculation)
		return
	}

	// 验证参数
	if err := calc.Validate(req.GetParams()); err != nil {
		h.sendError(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	// 使用修复版逻辑执行计算
	var result interface{}
	var warnings []string

	if odeCalc, ok := calc.(*calculator.ODESolverCalculator); ok {
		result, err = odeCalc.CalculateFixed(req.GetParams())
	} else {
		// 对于其他计算器，使用原逻辑
		result, err = calc.Calculate(req.GetParams())
	}

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

	// 计算器管理接口
	router.GET("/api/calculator-info", h.GetCalculatorInfo)

	// 支持接口
	router.GET("/api/supported-calculations", h.GetSupportedCalculations)

	// 对比接口
	router.POST("/api/solver/compare", h.SolverCompare)
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

// SolverCompare 对比修复前后求解器结果
// @Summary 对比修复前后求解器结果
// @Description 对比修复前后微分方程求解器的计算结果差异
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body CompareRequest true "对比请求参数（包含原始请求和修复后请求）"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) SolverCompare(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 执行原始计算
	originalResult, _, err := h.executeCalculation(req.Original, false)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原始计算失败: "+err.Error())
		return
	}

	// 执行修复后计算
	fixedResult, _, err := h.executeCalculation(req.Fixed, true)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复后计算失败: "+err.Error())
		return
	}

	// 计算差异
	difference := h.calculateDifference(originalResult, fixedResult)

	c.JSON(http.StatusOK, CompareResponse{
		Success:    true,
		Original:   originalResult,
		Fixed:      fixedResult,
		Difference: difference,
		Timestamp:  time.Now().Format(time.RFC3339),
	})
}

// executeCalculation 执行计算的辅助函数
func (h *APIHandler) executeCalculation(req models.CalculationRequest, useFixed bool) (interface{}, []string, error) {
	calcType, err := calculator.ParseCalculationType(req.Calculation)
	if err != nil {
		return nil, nil, err
	}

	calc, exists := h.calculatorManager.GetCalculator(calcType)
	if !exists {
		return nil, nil, fmt.Errorf("计算器未注册: %s", req.Calculation)
	}

	if err := calc.Validate(req.GetParams()); err != nil {
		return nil, nil, err
	}

	if useFixed {
		if odeCalc, ok := calc.(*calculator.ODESolverCalculator); ok {
			result, err := odeCalc.CalculateFixed(req.GetParams())
			return result, nil, err
		}
	}

	result, err := calc.Calculate(req.GetParams())
	return result, nil, err
}

// calculateDifference 计算结果差异
func (h *APIHandler) calculateDifference(original, fixed interface{}) interface{} {
	// 尝试转换为 ODEResult 格式
	origResult, ok1 := original.(*calculator.ODEResult)
	fixedResult, ok2 := fixed.(*calculator.ODEResult)

	if !ok1 || !ok2 {
		return map[string]string{
			"note": "无法解析为 ODE 求解器结果格式",
		}
	}

	// 计算时间点差异
	timeDiff := map[string]interface{}{}
	if len(origResult.TimePoints) == len(fixedResult.TimePoints) {
		maxDiff := 0.0
		uniformFixed := true
		uniformOrig := true

		for i := 0; i < len(fixedResult.TimePoints)-1; i++ {
			stepFixed := fixedResult.TimePoints[i+1] - fixedResult.TimePoints[i]
			stepOrig := origResult.TimePoints[i+1] - origResult.TimePoints[i]

			if math.Abs(stepFixed-0.1) > 0.0001 {
				uniformFixed = false
			}
			if math.Abs(stepOrig-0.1) > 0.0001 {
				uniformOrig = false
			}

			if math.Abs(stepFixed-stepOrig) > maxDiff {
				maxDiff = math.Abs(stepFixed - stepOrig)
			}
		}

		timeDiff["max_step_difference"] = maxDiff
		timeDiff["original_uniform"] = uniformOrig
		timeDiff["fixed_uniform"] = uniformFixed
	}

	timeDiff["original_count"] = len(origResult.TimePoints)
	timeDiff["fixed_count"] = len(fixedResult.TimePoints)

	// 解路径差异
	solutionDiff := map[string]interface{}{
		"original_count": len(origResult.SolutionPath),
		"fixed_count":    len(fixedResult.SolutionPath),
	}

	// 导数路径差异
	derivativeDiff := map[string]interface{}{
		"original_count": len(origResult.DerivativePath),
		"fixed_count":    len(fixedResult.DerivativePath),
	}

	// 误差估计差异
	errorDiff := map[string]interface{}{
		"original_error": origResult.ErrorEstimate,
		"fixed_error":    fixedResult.ErrorEstimate,
	}

	return map[string]interface{}{
		"time_points":     timeDiff,
		"solution_path":   solutionDiff,
		"derivative_path": derivativeDiff,
		"error_estimate":  errorDiff,
	}
}
