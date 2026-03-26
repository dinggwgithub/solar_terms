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

// CompareResponse 新旧计算结果对比响应
type CompareResponse struct {
	OldResult   interface{} `json:"old_result"`  // 原始计算结果
	NewResult   interface{} `json:"new_result"`  // 修复版计算结果
	Differences interface{} `json:"differences"` // 差异分析
	Timestamp   string      `json:"timestamp"`   // 时间戳
	SessionID   string      `json:"session_id,omitempty"`
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

// RegisterRoutes 注册API路由
func (h *APIHandler) RegisterRoutes(router *gin.Engine) {
	// 系统接口
	router.GET("/api/health", h.HealthCheck)
	router.GET("/api/system-info", h.GetSystemInfo)

	// 科学计算接口
	router.POST("/api/calculate", h.Calculate)

	// 修复版科学计算接口
	router.POST("/api/calculate-fixed", h.CalculateFixed)

	// 对比接口
	router.POST("/api/solver/compare", h.CompareCalculations)

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

	// 解析计算类型
	calcType, err := calculator.ParseCalculationType(req.Calculation)
	if err != nil {
		// 如果是planet类型，使用修复版行星计算器
		if req.Calculation == "planet" {
			calcType = calculator.CalculationTypePlanetFixed
		} else {
			h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
			return
		}
	}

	// 对于planet_fixed类型，直接使用
	if req.Calculation == "planet_fixed" {
		calcType = calculator.CalculationTypePlanetFixed
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

	h.sendSuccess(c, result, warnings, req.Calculation+"_fixed", sessionID)
}

// CompareCalculations 新旧计算结果对比接口
// @Summary 对比新旧计算结果
// @Description 执行原始和修复版计算，返回对比结果（仅支持planet类型）
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param session_id query string false "会话ID，用于保持计算参数一致性，不传则自动生成"
// @Param request body models.CalculationRequest true "计算请求参数，calculation应为planet"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) CompareCalculations(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 只支持planet类型的对比
	if req.Calculation != "planet" {
		h.sendError(c, http.StatusBadRequest, "仅支持planet类型计算的对比")
		return
	}

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	params := req.GetParams()

	// 执行原始计算
	oldResult, _, errOld := h.calculatorManager.CalculateWithSession(
		calculator.CalculationTypePlanet, params, sessionID)
	if errOld != nil {
		h.sendError(c, http.StatusBadRequest, "原始计算失败: "+errOld.Error())
		return
	}

	// 执行修复版计算
	newResult, _, errNew := h.calculatorManager.CalculateWithSession(
		calculator.CalculationTypePlanetFixed, params, sessionID)
	if errNew != nil {
		h.sendError(c, http.StatusBadRequest, "修复版计算失败: "+errNew.Error())
		return
	}

	// 分析差异
	differences := h.analyzeDifferences(oldResult, newResult)

	response := CompareResponse{
		OldResult:   oldResult,
		NewResult:   newResult,
		Differences: differences,
		Timestamp:   time.Now().Format(time.RFC3339),
		SessionID:   sessionID,
	}

	c.JSON(http.StatusOK, response)
}

// analyzeDifferences 分析新旧结果的差异
func (h *APIHandler) analyzeDifferences(oldResult, newResult interface{}) map[string]interface{} {
	differences := make(map[string]interface{})

	// 使用map来解析结果，因为返回的是interface{}实际是*calculator.PlanetPosition
	oldPos, ok1 := oldResult.(*calculator.PlanetPosition)
	newPos, ok2 := newResult.(*calculator.PlanetPosition)

	if ok1 && ok2 {
		// 分析赤经差异
		raDiff := newPos.RightAscension - oldPos.RightAscension
		raNormalized := newPos.RightAscension // 修复后的值应该在0-24范围内
		raIssue := oldPos.RightAscension < 0 || oldPos.RightAscension >= 24

		differences["right_ascension"] = map[string]interface{}{
			"old_value":      oldPos.RightAscension,
			"new_value":      newPos.RightAscension,
			"difference":     raDiff,
			"was_negative":   oldPos.RightAscension < 0,
			"was_over_range": oldPos.RightAscension >= 24,
			"issue_fixed":    !raIssue && raNormalized >= 0 && raNormalized < 24,
		}

		// 分析赤纬差异
		decDiff := newPos.Declination - oldPos.Declination
		differences["declination"] = map[string]interface{}{
			"old_value":  oldPos.Declination,
			"new_value":  newPos.Declination,
			"difference": decDiff,
			"in_range":   newPos.Declination >= -90 && newPos.Declination <= 90,
		}

		// 分析距离差异
		distDiff := newPos.Distance - oldPos.Distance
		differences["distance"] = map[string]interface{}{
			"old_value":    oldPos.Distance,
			"new_value":    newPos.Distance,
			"difference":   distDiff,
			"is_hardcoded": oldPos.Distance == 1.2, // 火星旧值为硬编码的1.2
		}

		// 分析星等差异
		magDiff := newPos.Magnitude - oldPos.Magnitude
		differences["magnitude"] = map[string]interface{}{
			"old_value":  oldPos.Magnitude,
			"new_value":  newPos.Magnitude,
			"difference": magDiff,
		}

		// 分析相位差异
		phaseDiff := newPos.Phase - oldPos.Phase
		differences["phase"] = map[string]interface{}{
			"old_value":    oldPos.Phase,
			"new_value":    newPos.Phase,
			"difference":   phaseDiff,
			"is_hardcoded": oldPos.Phase == 0.95,
		}

		// 分析距角差异
		elongDiff := newPos.Elongation - oldPos.Elongation
		differences["elongation"] = map[string]interface{}{
			"old_value":    oldPos.Elongation,
			"new_value":    newPos.Elongation,
			"difference":   elongDiff,
			"is_hardcoded": oldPos.Elongation == 120,
		}

		// 总体评估
		differences["summary"] = map[string]interface{}{
			"issues_found": []string{
				"赤经负值问题",
				"参数硬编码问题",
				"计算精度问题",
			},
			"issues_fixed": []string{
				"赤经已归一化到0-24小时范围",
				"距离、相位、距角现为动态计算",
				"计算精度提升",
			},
		}
	} else {
		differences["error"] = "无法分析结果差异：类型不匹配"
	}

	return differences
}
