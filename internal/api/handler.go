package api

import (
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

// CalculationCompareResponse 计算对比响应（用于A/B测试）
type CalculationCompareResponse struct {
	Success        bool                   `json:"success"`
	OriginalResult map[string]interface{} `json:"original_result"`
	FixedResult    map[string]interface{} `json:"fixed_result"`
	Differences    map[string]interface{} `json:"differences"`
	Analysis       *ComparisonAnalysis    `json:"analysis"`
	Timestamp      string                 `json:"timestamp"`
	SessionID      string                 `json:"session_id,omitempty"`
}

// ComparisonAnalysis 对比分析结果
type ComparisonAnalysis struct {
	IsWithinTolerance bool              `json:"is_within_tolerance"`
	Tolerance         float64           `json:"tolerance"`
	IssuesFound       []string          `json:"issues_found"`
	Summary           string            `json:"summary"`
	ParameterStatus   map[string]string `json:"parameter_status"`
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

// CalculateFixed 修复后的科学计算接口
// @Summary 执行修复后的科学计算
// @Description 执行修复后的各种类型的科学计算，特别是天文计算的修复版本
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

	// 如果请求的是astronomy，自动使用修复版本astronomy_fixed
	calcTypeStr := req.Calculation
	if calcTypeStr == "astronomy" {
		calcTypeStr = "astronomy_fixed"
	}

	calcType, err := calculator.ParseCalculationType(calcTypeStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
		return
	}

	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	result, warnings, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, calcTypeStr, sessionID)
}

// CalculateCompare A/B测试对比接口
// @Summary 计算结果对比
// @Description 同时调用原始接口和修复接口，进行A/B测试对比分析
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID，用于保持计算参数一致性，不传则自动生成"
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CalculationCompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-compare [post]
func (h *APIHandler) CalculateCompare(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 调用原始计算器
	originalType, err := calculator.ParseCalculationType("astronomy")
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原始计算器类型错误: "+err.Error())
		return
	}

	originalResult, _, err := h.calculatorManager.CalculateWithSession(originalType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原始计算失败: "+err.Error())
		return
	}

	// 调用修复后的计算器
	fixedType, err := calculator.ParseCalculationType("astronomy_fixed")
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复计算器类型错误: "+err.Error())
		return
	}

	fixedResult, _, err := h.calculatorManager.CalculateWithSession(fixedType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复计算失败: "+err.Error())
		return
	}

	// 对比分析
	comparison := h.compareResults(originalResult, fixedResult)

	response := CalculationCompareResponse{
		Success:        true,
		OriginalResult: h.convertToMap(originalResult),
		FixedResult:    h.convertToMap(fixedResult),
		Differences:    comparison.Differences,
		Analysis:       comparison.Analysis,
		Timestamp:      time.Now().Format(time.RFC3339),
		SessionID:      sessionID,
	}

	c.JSON(http.StatusOK, response)
}

// ComparisonResult 对比结果
type ComparisonResult struct {
	Differences map[string]interface{} `json:"differences"`
	Analysis    *ComparisonAnalysis    `json:"analysis"`
}

// compareResults 对比两个计算结果
func (h *APIHandler) compareResults(original, fixed interface{}) *ComparisonResult {
	originalMap := h.convertToMap(original)
	fixedMap := h.convertToMap(fixed)

	differences := make(map[string]interface{})
	parameterStatus := make(map[string]string)
	issuesFound := make([]string, 0)
	tolerance := 0.5 // 容差±0.5度
	allWithinTolerance := true

	// 定义合理范围
	reasonableRanges := map[string][2]float64{
		"sun_longitude":      {0, 360},
		"apparent_longitude": {0, 360},
		"true_longitude":     {0, 360},
		"mean_longitude":     {0, 360},
		"mean_anomaly":       {0, 360},
		"julian_date":        {2415020.5, 2488070.5},
		"equation_of_center": {-2, 2},
		"nutation":           {-0.01, 0.01},
	}

	for key, origVal := range originalMap {
		if fixedVal, exists := fixedMap[key]; exists {
			origFloat, ok1 := origVal.(float64)
			fixedFloat, ok2 := fixedVal.(float64)

			if ok1 && ok2 {
				diff := fixedFloat - origFloat
				absDiff := diff
				if absDiff < 0 {
					absDiff = -absDiff
				}

				differences[key] = map[string]interface{}{
					"original":   origFloat,
					"fixed":      fixedFloat,
					"difference": diff,
					"abs_diff":   absDiff,
				}

				// 检查是否在合理范围内
				if validRange, hasRange := reasonableRanges[key]; hasRange {
					origInRange := origFloat >= validRange[0] && origFloat <= validRange[1]
					fixedInRange := fixedFloat >= validRange[0] && fixedFloat <= validRange[1]

					if !origInRange && fixedInRange {
						issuesFound = append(issuesFound, key+" 原始值超出合理范围")
						parameterStatus[key] = "已修复"
					} else if !origInRange && !fixedInRange {
						issuesFound = append(issuesFound, key+" 仍超出合理范围")
						parameterStatus[key] = "异常"
						allWithinTolerance = false
					} else if origInRange && fixedInRange {
						if absDiff > tolerance {
							parameterStatus[key] = "差异较大"
							issuesFound = append(issuesFound, key+" 偏差超过容差")
						} else {
							parameterStatus[key] = "正常"
						}
					}
				} else {
					parameterStatus[key] = "无参考范围"
				}

				// 检查差异是否超过容差
				if absDiff > tolerance {
					allWithinTolerance = false
				}
			}
		}
	}

	// 生成摘要
	summary := "修复效果分析："
	if len(issuesFound) > 0 {
		summary += string(len(issuesFound)) + "个问题已解决"
	} else {
		summary += "所有参数在合理范围内"
	}

	return &ComparisonResult{
		Differences: differences,
		Analysis: &ComparisonAnalysis{
			IsWithinTolerance: allWithinTolerance,
			Tolerance:         tolerance,
			IssuesFound:       issuesFound,
			Summary:           summary,
			ParameterStatus:   parameterStatus,
		},
	}
}

// convertToMap 将interface{}转换为map[string]interface{}
func (h *APIHandler) convertToMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]float64); ok {
		result := make(map[string]interface{})
		for k, val := range m {
			result[k] = val
		}
		return result
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return make(map[string]interface{})
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
	router.POST("/api/calculate-compare", h.CalculateCompare)

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
