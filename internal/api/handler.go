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
	calculatorManager   *calculator.CalculatorManager
	starCalculatorFixed *calculator.StarCalculatorFixed
}

// NewAPIHandler 创建新的API处理器
func NewAPIHandler(calculatorManager *calculator.CalculatorManager) *APIHandler {
	return &APIHandler{
		calculatorManager:   calculatorManager,
		starCalculatorFixed: calculator.NewStarCalculatorFixed(),
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
	router.POST("/api/solver/compare", h.CompareResults)

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
// @Description 执行修复版科学计算，修正北斗七星、二十八宿方位、干支历法等错误
// @Tags 科学计算
// @Accept json
// @Produce json
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

	sessionID := GenerateSessionID()

	if req.Calculation == "star" {
		result, err := h.starCalculatorFixed.Calculate(req.GetParams())
		if err != nil {
			h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
			return
		}
		h.sendSuccess(c, result, nil, req.Calculation+"_fixed", sessionID)
		return
	}

	calcType, err := calculator.ParseCalculationType(req.Calculation)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
		return
	}

	result, warnings, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, req.Calculation, sessionID)
}

// CompareResult 差异对比结果
type CompareResult struct {
	Field         string      `json:"field"`
	OriginalValue interface{} `json:"original_value"`
	FixedValue    interface{} `json:"fixed_value"`
	Description   string      `json:"description"`
}

// CompareResponse 对比响应
type CompareResponse struct {
	Success        bool            `json:"success"`
	OriginalResult interface{}     `json:"original_result"`
	FixedResult    interface{}     `json:"fixed_result"`
	Differences    []CompareResult `json:"differences"`
	Summary        string          `json:"summary"`
	Timestamp      string          `json:"timestamp"`
}

// CompareResults 对比接口
// @Summary 对比原始计算与修复版计算结果
// @Description 对比原始缺陷响应与修复后响应的结构化对比结果
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) CompareResults(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	sessionID := GenerateSessionID()

	var originalResult, fixedResult interface{}
	var err error

	if req.Calculation == "star" {
		calcType := calculator.CalculationTypeStar
		originalResult, _, err = h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, "原始计算失败: "+err.Error())
			return
		}

		fixedResult, err = h.starCalculatorFixed.Calculate(req.GetParams())
		if err != nil {
			h.sendError(c, http.StatusBadRequest, "修复版计算失败: "+err.Error())
			return
		}
	} else {
		calcType, err := calculator.ParseCalculationType(req.Calculation)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
			return
		}

		originalResult, _, err = h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
		if err != nil {
			h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
			return
		}
		fixedResult = originalResult
	}

	differences := h.compareStarResults(originalResult, fixedResult)

	summary := fmt.Sprintf("共发现 %d 处差异", len(differences))

	c.JSON(http.StatusOK, CompareResponse{
		Success:        true,
		OriginalResult: originalResult,
		FixedResult:    fixedResult,
		Differences:    differences,
		Summary:        summary,
		Timestamp:      time.Now().Format(time.RFC3339),
	})
}

// compareStarResults 对比星曜计算结果
func (h *APIHandler) compareStarResults(original, fixed interface{}) []CompareResult {
	var differences []CompareResult

	origMap, ok1 := original.(map[string]interface{})
	fixedMap, ok2 := fixed.(map[string]interface{})

	if !ok1 || !ok2 {
		origResult, ok1 := original.(*calculator.StarResult)
		fixedResult, ok2 := fixed.(*calculator.StarResultFixed)
		if ok1 && ok2 {
			return h.compareStarResultStructs(origResult, fixedResult)
		}
		return differences
	}

	fieldsToCompare := []struct {
		field       string
		description string
	}{
		{"lunar_date", "农历日期"},
		{"day_ganzhi", "日干支"},
		{"constellation", "二十八宿"},
		{"star_position", "星曜位置"},
		{"auspicious", "是否吉日"},
		{"day_score", "日分值"},
		{"constellation_index", "二十八宿索引"},
		{"auspicious_level", "吉凶程度"},
		{"julian_day", "儒略日"},
	}

	for _, f := range fieldsToCompare {
		origVal, origExists := origMap[f.field]
		fixedVal, fixedExists := fixedMap[f.field]

		if origExists && fixedExists {
			if !h.compareValues(origVal, fixedVal) {
				differences = append(differences, CompareResult{
					Field:         f.field,
					OriginalValue: origVal,
					FixedValue:    fixedVal,
					Description:   f.description + " 存在差异",
				})
			}
		}
	}

	return differences
}

// compareStarResultStructs 对比结构体类型的星曜计算结果
func (h *APIHandler) compareStarResultStructs(orig *calculator.StarResult, fixed *calculator.StarResultFixed) []CompareResult {
	var differences []CompareResult

	if orig.LunarDate != fixed.LunarDate {
		differences = append(differences, CompareResult{
			Field:         "lunar_date",
			OriginalValue: orig.LunarDate,
			FixedValue:    fixed.LunarDate,
			Description:   "农历日期修正：原实现直接使用公历月日作为农历，修复版实现正确的公历转农历算法",
		})
	}

	if orig.DayGanZhi != fixed.DayGanZhi {
		differences = append(differences, CompareResult{
			Field:         "day_ganzhi",
			OriginalValue: orig.DayGanZhi,
			FixedValue:    fixed.DayGanZhi,
			Description:   "日干支修正：原实现使用错误的基准日，修复版使用正确的1900年1月31日（甲子日）作为基准",
		})
	}

	if orig.Constellation != fixed.Constellation {
		differences = append(differences, CompareResult{
			Field:         "constellation",
			OriginalValue: orig.Constellation,
			FixedValue:    fixed.Constellation,
			Description:   "二十八宿修正：调整了二十八宿值日计算算法",
		})
	}

	if orig.StarPosition != fixed.StarPosition {
		differences = append(differences, CompareResult{
			Field:         "star_position",
			OriginalValue: orig.StarPosition,
			FixedValue:    fixed.StarPosition,
			Description:   "星曜位置修正：原实现将南方朱雀标注为西方，修复版正确标注方位",
		})
	}

	if orig.DayScore != fixed.DayScore {
		differences = append(differences, CompareResult{
			Field:         "day_score",
			OriginalValue: orig.DayScore,
			FixedValue:    fixed.DayScore,
			Description:   "日分值修正：原实现使用日期哈希，修复版基于干支和星宿吉凶计算",
		})
	}

	if orig.AuspiciousLevel != fixed.AuspiciousLevel {
		differences = append(differences, CompareResult{
			Field:         "auspicious_level",
			OriginalValue: orig.AuspiciousLevel,
			FixedValue:    fixed.AuspiciousLevel,
			Description:   "吉凶程度修正：修复评分逻辑使其与吉凶判断自洽",
		})
	}

	if fixed.BigDipperInfo != nil {
		differences = append(differences, CompareResult{
			Field:         "big_dipper_info",
			OriginalValue: nil,
			FixedValue:    fixed.BigDipperInfo,
			Description:   "新增北斗七星专属信息：原实现未处理star_name=big_dipper参数，修复版返回北斗七星详细天文数据",
		})
	}

	return differences
}

// compareValues 比较两个值是否相等
func (h *APIHandler) compareValues(a, b interface{}) bool {
	switch a := a.(type) {
	case float64:
		if b, ok := b.(float64); ok {
			return a == b
		}
	case string:
		if b, ok := b.(string); ok {
			return a == b
		}
	case bool:
		if b, ok := b.(bool); ok {
			return a == b
		}
	case int:
		if b, ok := b.(int); ok {
			return a == b
		}
		if b, ok := b.(float64); ok {
			return float64(a) == b
		}
	}
	return false
}
