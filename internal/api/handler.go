package api

import (
	"math/rand"
	"net/http"
	"scientific_calc_bugs/internal/bugs"
	"scientific_calc_bugs/internal/calculator"
	"scientific_calc_bugs/models"
	"time"

	"github.com/gin-gonic/gin"
)

// APIHandler API处理器
type APIHandler struct {
	calculatorManager *calculator.CalculatorManager
	bugManager        *bugs.BugManager
}

// NewAPIHandler 创建新的API处理器
func NewAPIHandler(calculatorManager *calculator.CalculatorManager, bugManager *bugs.BugManager) *APIHandler {
	return &APIHandler{
		calculatorManager: calculatorManager,
		bugManager:        bugManager,
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
	BugType     string      `json:"bug_type,omitempty"`
	Calculation string      `json:"calculation"`
	SessionID   string      `json:"session_id,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

// BugInfoResponse Bug信息响应
type BugInfoResponse struct {
	BugType         string            `json:"bug_type"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Characteristics map[string]string `json:"characteristics"`
}

// CalculatorInfoResponse 计算器信息响应
type CalculatorInfoResponse struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	SupportedBugs []string `json:"supported_bugs"`
}

// SystemInfoResponse 系统信息响应
type SystemInfoResponse struct {
	Version               string   `json:"version"`
	SupportedCalculations []string `json:"supported_calculations"`
	SupportedBugTypes     []string `json:"supported_bug_types"`
	TotalCalculators      int      `json:"total_calculators"`
	TotalBugs             int      `json:"total_bugs"`
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

// CalculateWithBugs 带有bug的计算接口
// @Summary 执行带有bug的计算
// @Description 执行带有bug的科学计算（结果不稳定、约束越界、精度错误）
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param bug_type query string false "bug类型" Enums(instability, constraint, precision) Default(constraint)
// @Param session_id query string false "会话ID，用于保持Bug参数一致性，不传则自动生成"
// @Param mixed_mode query boolean false "是否启用混合Bug模式" Default(false)
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-with-bugs [post]
func (h *APIHandler) CalculateWithBugs(c *gin.Context) {
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

	// 解析Bug类型
	// 防作弊机制：不允许使用 bug_type=none（无Bug模式），未指定时随机分配Bug类型
	bugTypeStr := c.DefaultQuery("bug_type", "")

	// 如果未指定或请求了 none（作弊尝试），随机分配一个Bug类型
	if bugTypeStr == "" || bugTypeStr == "none" {
		availableTypes := []string{"instability", "constraint", "precision"}
		rand.Seed(time.Now().UnixNano())
		bugTypeStr = availableTypes[rand.Intn(len(availableTypes))]
		// 添加Header用于调试，但AI通常不会检查这个
		c.Header("X-Forced-Bug-Type", bugTypeStr)
	}

	bugType, err := h.bugManager.GetBugTypeFromString(bugTypeStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的Bug类型: "+bugTypeStr)
		return
	}

	// 获取或生成会话ID（用于保持Bug参数一致性）
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = bugs.GenerateSessionID()
	}

	// 解析混合模式
	mixedMode := c.DefaultQuery("mixed_mode", "false") == "true"

	// 执行计算
	result, warnings, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), bugType, sessionID, mixedMode)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, bugType.String(), req.Calculation, sessionID)
}

// GetBugInfo 获取Bug信息接口
// @Summary 获取Bug信息
// @Description 获取指定Bug类型的详细信息
// @Tags Bug管理
// @Accept json
// @Produce json
// @Param bug_type query string true "bug类型" Enums(instability, constraint, precision)
// @Success 200 {object} BugInfoResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/bug-info [get]
func (h *APIHandler) GetBugInfo(c *gin.Context) {
	bugTypeStr := c.Query("bug_type")
	if bugTypeStr == "" {
		h.sendError(c, http.StatusBadRequest, "缺少bug_type参数")
		return
	}

	bugType, err := h.bugManager.GetBugTypeFromString(bugTypeStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的Bug类型: "+bugTypeStr)
		return
	}

	bugInfo, err := h.bugManager.GetBugInfo(bugType)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "获取Bug信息失败: "+err.Error())
		return
	}

	response := BugInfoResponse{
		BugType:         bugType.String(),
		Name:            bugInfo["name"],
		Description:     bugInfo["description"],
		Characteristics: make(map[string]string),
	}

	// 提取特征信息
	for key, value := range bugInfo {
		if len(key) > 14 && key[:14] == "characteristic_" {
			response.Characteristics[key[14:]] = value
		}
	}

	c.JSON(http.StatusOK, response)
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
		Name:          calcInfo["name"],
		Description:   calcInfo["description"],
		SupportedBugs: []string{},
	}

	// 解析支持的Bug类型
	if supportedBugs, exists := calcInfo["supported_bugs"]; exists {
		// 这里应该解析字符串为切片
		// 简化处理
		response.SupportedBugs = []string{supportedBugs}
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

	// 获取支持的Bug类型
	supportedBugTypes := h.bugManager.GetSupportedBugTypes()
	bugTypes := make([]string, 0, len(supportedBugTypes))
	for _, bugType := range supportedBugTypes {
		bugTypes = append(bugTypes, bugType.String())
	}

	// 获取统计信息
	bugStats := h.bugManager.GetBugStatistics()
	totalBugs := 0
	if total, exists := bugStats["total_bugs"]; exists {
		totalBugs = total.(int)
	}

	response := SystemInfoResponse{
		Version:               "1.0.0",
		SupportedCalculations: calcTypes,
		SupportedBugTypes:     bugTypes,
		TotalCalculators:      len(supportedCalculations),
		TotalBugs:             totalBugs,
	}

	c.JSON(http.StatusOK, response)
}

// GetBugStatistics 获取Bug统计信息接口
// @Summary 获取Bug统计信息
// @Description 获取Bug系统的统计信息
// @Tags Bug管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/bug-statistics [get]
func (h *APIHandler) GetBugStatistics(c *gin.Context) {
	stats := h.bugManager.GetBugStatistics()
	c.JSON(http.StatusOK, stats)
}

// sendSuccess 发送成功响应
func (h *APIHandler) sendSuccess(c *gin.Context, result interface{}, warnings []string, bugType, calculation string, sessionID ...string) {
	response := CalculationResponse{
		Success:     true,
		Result:      result,
		Warnings:    warnings,
		Timestamp:   time.Now().Format(time.RFC3339),
		BugType:     bugType,
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
	router.POST("/api/calculate-with-bugs", h.CalculateWithBugs)

	// Bug管理接口
	router.GET("/api/bug-info", h.GetBugInfo)
	router.GET("/api/bug-statistics", h.GetBugStatistics)

	// 计算器管理接口
	router.GET("/api/calculator-info", h.GetCalculatorInfo)

	// 实验支持接口
	router.GET("/api/supported-calculations", h.GetSupportedCalculations)
	router.GET("/api/supported-bug-types", h.GetSupportedBugTypes)
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

// GetSupportedBugTypes 获取支持的Bug类型接口
func (h *APIHandler) GetSupportedBugTypes(c *gin.Context) {
	supportedBugTypes := h.bugManager.GetSupportedBugTypes()

	var bugTypes []map[string]interface{}
	for _, bugType := range supportedBugTypes {
		bugInfo, err := h.bugManager.GetBugInfo(bugType)
		if err == nil {
			bugTypes = append(bugTypes, map[string]interface{}{
				"type":        bugType.String(),
				"name":        bugInfo["name"],
				"description": bugInfo["description"],
			})
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"bug_types": bugTypes,
		"total":     len(bugTypes),
	})
}
