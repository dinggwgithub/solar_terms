package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
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

// CalculateFixedResponse 修复版计算响应
type CalculateFixedResponse struct {
	Success     bool        `json:"success"`
	Result      interface{} `json:"result"`
	Warnings    []string    `json:"warnings"`
	Timestamp   string      `json:"timestamp"`
	Calculation string      `json:"calculation"`
	SessionID   string      `json:"session_id,omitempty"`
	Note        string      `json:"note"`
}

// CompareResponse 对比响应
type CompareResponse struct {
	Success        bool                   `json:"success"`
	Timestamp      string                 `json:"timestamp"`
	SessionID      string                 `json:"session_id,omitempty"`
	Comparison     map[string]interface{} `json:"comparison"`
	Differences    []DifferenceDetail     `json:"differences"`
	Summary        string                 `json:"summary"`
	OriginalBug    string                 `json:"original_bug_description"`
	FixExplanation string                 `json:"fix_explanation"`
}

// DifferenceDetail 差异详情
type DifferenceDetail struct {
	Field       string      `json:"field"`
	Original    interface{} `json:"original"`
	Fixed       interface{} `json:"fixed"`
	Description string      `json:"description"`
	IsCritical  bool        `json:"is_critical"`
}

// CalculateFixed 修复后的日出日落计算接口
// @Summary 修复后的日出日落时间计算
// @Description 使用修复后的天文算法计算日出日落时间，修复了原版的符号处理和时区转换错误
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body models.CalculationRequest true "计算请求参数（calculation必须为sunrise_sunset）"
// @Success 200 {object} CalculateFixedResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-fixed [post]
func (h *APIHandler) CalculateFixed(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 只支持日出日落计算
	if req.Calculation != "sunrise_sunset" {
		h.sendError(c, http.StatusBadRequest, "此接口只支持sunrise_sunset计算类型")
		return
	}

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 使用修复后的计算器
	fixedCalc := calculator.NewSunriseSunsetCalculatorFixed()

	// 验证参数
	if err := fixedCalc.Validate(req.GetParams()); err != nil {
		h.sendError(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	// 执行计算
	result, err := fixedCalc.Calculate(req.GetParams())
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	response := CalculateFixedResponse{
		Success:     true,
		Result:      result,
		Timestamp:   time.Now().Format(time.RFC3339),
		Calculation: req.Calculation,
		SessionID:   sessionID,
		Note:        "此结果使用修复后的天文算法计算，修复了原版的赤纬计算、时区转换和符号处理错误",
	}

	c.JSON(http.StatusOK, response)
}

// CalculateCompare 原接口与修复后接口的对比接口
// @Summary 对比原版与修复版的计算结果
// @Description 接收相同请求参数，返回原接口与修复后接口的结果差异，便于验证修复效果
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body models.CalculationRequest true "计算请求参数（calculation必须为sunrise_sunset）"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/calculate-compare [post]
func (h *APIHandler) CalculateCompare(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 只支持日出日落计算
	if req.Calculation != "sunrise_sunset" {
		h.sendError(c, http.StatusBadRequest, "此接口只支持sunrise_sunset计算类型")
		return
	}

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	params := req.GetParams()

	// 执行原版计算
	calcType, _ := calculator.ParseCalculationType(req.Calculation)
	originalResult, _, err := h.calculatorManager.CalculateWithSession(calcType, params, sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原版计算失败: "+err.Error())
		return
	}

	// 执行修复版计算
	fixedCalc := calculator.NewSunriseSunsetCalculatorFixed()
	if err := fixedCalc.Validate(params); err != nil {
		h.sendError(c, http.StatusBadRequest, "修复版参数验证失败: "+err.Error())
		return
	}
	fixedResult, err := fixedCalc.Calculate(params)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复版计算失败: "+err.Error())
		return
	}

	// 对比结果
	comparison, differences := h.compareResults(originalResult, fixedResult)

	// 生成总结
	summary := h.generateComparisonSummary(differences)

	response := CompareResponse{
		Success:        true,
		Timestamp:      time.Now().Format(time.RFC3339),
		SessionID:      sessionID,
		Comparison:     comparison,
		Differences:    differences,
		Summary:        summary,
		OriginalBug:    "原版存在时间逻辑错误：日出显示为18:10（傍晚），日落显示为06:18（早晨），太阳正午显示为00:14（午夜）",
		FixExplanation: "修复内容：1) 使用正确的太阳赤纬计算（而非固定0度）；2) 修正时区转换逻辑；3) 正确处理经度符号和方程时差",
	}

	c.JSON(http.StatusOK, response)
}

// compareResults 对比原版和修复版的结果
func (h *APIHandler) compareResults(original, fixed interface{}) (map[string]interface{}, []DifferenceDetail) {
	comparison := make(map[string]interface{})
	var differences []DifferenceDetail

	// 使用反射获取结果字段
	origVal := reflect.ValueOf(original)
	fixedVal := reflect.ValueOf(fixed)

	// 处理指针
	if origVal.Kind() == reflect.Ptr {
		origVal = origVal.Elem()
	}
	if fixedVal.Kind() == reflect.Ptr {
		fixedVal = fixedVal.Elem()
	}

	// 只处理结构体
	if origVal.Kind() != reflect.Struct || fixedVal.Kind() != reflect.Struct {
		comparison["error"] = "结果类型不支持对比"
		return comparison, differences
	}

	origType := origVal.Type()

	// 遍历所有字段
	for i := 0; i < origVal.NumField(); i++ {
		field := origType.Field(i)
		fieldName := field.Name

		origField := origVal.Field(i)
		fixedField := fixedVal.FieldByName(fieldName)

		if !fixedField.IsValid() {
			continue
		}

		origValue := origField.Interface()
		fixedValue := fixedField.Interface()

		// 添加到对比结果
		comparison[fieldName] = map[string]interface{}{
			"original": origValue,
			"fixed":    fixedValue,
			"same":     reflect.DeepEqual(origValue, fixedValue),
		}

		// 如果不相同，添加到差异列表
		if !reflect.DeepEqual(origValue, fixedValue) {
			desc := h.getDifferenceDescription(fieldName, origValue, fixedValue)
			isCritical := h.isCriticalField(fieldName)
			differences = append(differences, DifferenceDetail{
				Field:       fieldName,
				Original:    origValue,
				Fixed:       fixedValue,
				Description: desc,
				IsCritical:  isCritical,
			})
		}
	}

	return comparison, differences
}

// getDifferenceDescription 获取差异描述
func (h *APIHandler) getDifferenceDescription(field string, original, fixed interface{}) string {
	switch field {
	case "Sunrise":
		return "日出时间存在显著差异，原版可能显示为傍晚时间"
	case "Sunset":
		return "日落时间存在显著差异，原版可能显示为早晨时间"
	case "SolarNoon":
		return "太阳正午时间存在差异，原版可能显示为午夜时间"
	case "DayLength":
		return "白昼长度存在差异"
	default:
		return fmt.Sprintf("%s 字段存在差异", field)
	}
}

// isCriticalField 判断是否为关键字段
func (h *APIHandler) isCriticalField(field string) bool {
	criticalFields := []string{"Sunrise", "Sunset", "SolarNoon", "DayLength"}
	for _, cf := range criticalFields {
		if cf == field {
			return true
		}
	}
	return false
}

// generateComparisonSummary 生成对比总结
func (h *APIHandler) generateComparisonSummary(differences []DifferenceDetail) string {
	if len(differences) == 0 {
		return "原版与修复版结果完全一致，未发现差异"
	}

	criticalCount := 0
	for _, d := range differences {
		if d.IsCritical {
			criticalCount++
		}
	}

	if criticalCount > 0 {
		return fmt.Sprintf("发现 %d 处差异，其中 %d 处为关键字段（日出/日落/正午时间），确认原版存在Bug", len(differences), criticalCount)
	}

	return fmt.Sprintf("发现 %d 处非关键字段差异", len(differences))
}

// RegisterRoutes 注册API路由
func (h *APIHandler) RegisterRoutes(router *gin.Engine) {
	// 系统接口
	router.GET("/api/health", h.HealthCheck)
	router.GET("/api/system-info", h.GetSystemInfo)

	// 科学计算接口
	router.POST("/api/calculate", h.Calculate)

	// 修复版计算接口
	router.POST("/api/calculate-fixed", h.CalculateFixed)

	// 对比接口
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
