package api

import (
	"fmt"
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

// CalculateFixed 修复后的科学计算接口
// @Summary 执行修复后的科学计算
// @Description 执行修复后的日出日落时间计算
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

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 使用修复后的计算器执行计算
	result, warnings, err := h.calculatorManager.CalculateFixedWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, req.Calculation, sessionID)
}

// CompareResponse 对比响应
type CompareResponse struct {
	Success        bool                   `json:"success"`
	OriginalResult interface{}            `json:"original_result"`
	FixedResult    interface{}            `json:"fixed_result"`
	Differences    map[string]interface{} `json:"differences"`
	Summary        string                 `json:"summary"`
	Timestamp      string                 `json:"timestamp"`
	Calculation    string                 `json:"calculation"`
	SessionID      string                 `json:"session_id,omitempty"`
}

// CalculateCompare 结果对比接口
// @Summary 对比原接口与修复后接口的计算结果
// @Description 接收相同请求参数，返回原接口与修复后接口的结果差异
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID"
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-compare [post]
func (h *APIHandler) CalculateCompare(c *gin.Context) {
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

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 执行原计算
	originalResult, _, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原计算失败: "+err.Error())
		return
	}

	// 执行修复后计算
	fixedResult, _, err := h.calculatorManager.CalculateFixedWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复后计算失败: "+err.Error())
		return
	}

	// 对比结果差异
	differences, summary := compareResults(originalResult, fixedResult)

	response := CompareResponse{
		Success:        true,
		OriginalResult: originalResult,
		FixedResult:    fixedResult,
		Differences:    differences,
		Summary:        summary,
		Timestamp:      time.Now().Format(time.RFC3339),
		Calculation:    req.Calculation,
		SessionID:      sessionID,
	}

	c.JSON(http.StatusOK, response)
}

// compareResults 对比两个结果的差异
func compareResults(original, fixed interface{}) (map[string]interface{}, string) {
	differences := make(map[string]interface{})
	summary := ""

	origMap, ok1 := original.(map[string]interface{})
	fixedMap, ok2 := fixed.(map[string]interface{})

	if !ok1 || !ok2 {
		// 尝试转换为其他类型
		origResult, ok1 := original.(*calculator.SunriseSunsetResult)
		fixedResult, ok2 := fixed.(*calculator.SunriseSunsetResult)

		if ok1 && ok2 {
			return compareSunriseSunsetResults(origResult, fixedResult)
		}

		differences["note"] = "结果类型无法直接对比"
		return differences, "结果结构不同"
	}

	// 简单的字段级对比
	for key, origVal := range origMap {
		if fixedVal, exists := fixedMap[key]; exists {
			if origVal != fixedVal {
				differences[key] = map[string]interface{}{
					"original": origVal,
					"fixed":    fixedVal,
				}
			}
		} else {
			differences[key] = map[string]interface{}{
				"original": origVal,
				"fixed":    "字段不存在",
			}
		}
	}

	for key, fixedVal := range fixedMap {
		if _, exists := origMap[key]; !exists {
			differences[key] = map[string]interface{}{
				"original": "字段不存在",
				"fixed":    fixedVal,
			}
		}
	}

	if len(differences) > 0 {
		summary = fmt.Sprintf("发现 %d 处差异", len(differences))
	} else {
		summary = "结果一致，无差异"
	}

	return differences, summary
}

// compareSunriseSunsetResults 专门对比日出日落结果
func compareSunriseSunsetResults(orig, fixed *calculator.SunriseSunsetResult) (map[string]interface{}, string) {
	differences := make(map[string]interface{})
	issues := []string{}

	if orig.Sunrise != fixed.Sunrise {
		differences["sunrise"] = map[string]string{
			"original": orig.Sunrise,
			"fixed":    fixed.Sunrise,
		}
		issues = append(issues, "日出时间")
	}

	if orig.Sunset != fixed.Sunset {
		differences["sunset"] = map[string]string{
			"original": orig.Sunset,
			"fixed":    fixed.Sunset,
		}
		issues = append(issues, "日落时间")
	}

	if orig.SolarNoon != fixed.SolarNoon {
		differences["solar_noon"] = map[string]string{
			"original": orig.SolarNoon,
			"fixed":    fixed.SolarNoon,
		}
		issues = append(issues, "正午时间")
	}

	// 检查晨昏蒙影
	if orig.CivilTwilight.Morning != fixed.CivilTwilight.Morning {
		differences["civil_twilight.morning"] = map[string]string{
			"original": orig.CivilTwilight.Morning,
			"fixed":    fixed.CivilTwilight.Morning,
		}
		issues = append(issues, "民用晨光")
	}

	if orig.CivilTwilight.Evening != fixed.CivilTwilight.Evening {
		differences["civil_twilight.evening"] = map[string]string{
			"original": orig.CivilTwilight.Evening,
			"fixed":    fixed.CivilTwilight.Evening,
		}
		issues = append(issues, "民用暮光")
	}

	summary := ""
	if len(issues) > 0 {
		summary = fmt.Sprintf("修复了以下问题：%v", issues)
	} else {
		summary = "结果一致，无差异"
	}

	return differences, summary
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
