package api

import (
	"math"
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

// RegisterRoutes 注册API路由
func (h *APIHandler) RegisterRoutes(router *gin.Engine) {
	// 系统接口
	router.GET("/api/health", h.HealthCheck)
	router.GET("/api/system-info", h.GetSystemInfo)

	// 科学计算接口
	router.POST("/api/calculate", h.Calculate)

	// 修复版科学计算接口
	router.POST("/api/calculate-fixed", h.CalculateFixed)

	// 求解器对比接口
	router.POST("/api/solver/compare", h.CompareSolvers)

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
// @Summary 执行修复后的科学计算
// @Description 执行修复后的方程求解计算，确保数值计算正确
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

	// 只支持 equation_solver 类型使用修复版
	if req.Calculation != "equation_solver" {
		h.sendError(c, http.StatusBadRequest, "修复版接口仅支持 equation_solver 计算类型")
		return
	}

	// 获取或生成会话ID
	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	// 使用修复后的计算器
	fixedCalc := calculator.NewEquationSolverCalculatorFixed()
	if err := fixedCalc.Validate(req.GetParams()); err != nil {
		h.sendError(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	result, err := fixedCalc.Calculate(req.GetParams())
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, nil, req.Calculation, sessionID)
}

// CompareResult 对比结果结构
type CompareResult struct {
	Field        string      `json:"field"`
	OldValue     interface{} `json:"old_value"`
	NewValue     interface{} `json:"new_value"`
	Difference   float64     `json:"difference,omitempty"`
	RelativeDiff float64     `json:"relative_diff,omitempty"`
	Status       string      `json:"status"` // "same", "different", "added", "removed"
}

// CompareResponse 对比响应
type CompareResponse struct {
	Success         bool            `json:"success"`
	OldResult       interface{}     `json:"old_result"`
	NewResult       interface{}     `json:"new_result"`
	Comparisons     []CompareResult `json:"comparisons"`
	TotalFields     int             `json:"total_fields"`
	DifferentFields int             `json:"different_fields"`
	Timestamp       string          `json:"timestamp"`
}

// CompareSolvers 求解器对比接口
// @Summary 对比新旧求解器结果
// @Description 同时运行旧版和新版求解器，返回字段级差异对比
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body models.CalculationRequest true "计算请求参数"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) CompareSolvers(c *gin.Context) {
	var req models.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 只支持 equation_solver 类型
	if req.Calculation != "equation_solver" {
		h.sendError(c, http.StatusBadRequest, "对比接口仅支持 equation_solver 计算类型")
		return
	}

	params := req.GetParams()

	// 运行旧版求解器
	oldCalc := calculator.NewEquationSolverCalculator()
	var oldResult interface{}
	if err := oldCalc.Validate(params); err == nil {
		oldResult, _ = oldCalc.Calculate(params)
	}

	// 运行新版求解器
	newCalc := calculator.NewEquationSolverCalculatorFixed()
	var newResult interface{}
	if err := newCalc.Validate(params); err == nil {
		newResult, _ = newCalc.Calculate(params)
	}

	// 对比结果
	comparisons := h.compareResults(oldResult, newResult)

	// 统计差异
	differentCount := 0
	for _, comp := range comparisons {
		if comp.Status == "different" {
			differentCount++
		}
	}

	c.JSON(http.StatusOK, CompareResponse{
		Success:         true,
		OldResult:       oldResult,
		NewResult:       newResult,
		Comparisons:     comparisons,
		TotalFields:     len(comparisons),
		DifferentFields: differentCount,
		Timestamp:       time.Now().Format(time.RFC3339),
	})
}

// compareResults 对比两个结果结构
func (h *APIHandler) compareResults(oldResult, newResult interface{}) []CompareResult {
	var comparisons []CompareResult

	if oldResult == nil && newResult == nil {
		return comparisons
	}

	// 获取结果的所有字段
	oldFields := h.extractFields(oldResult)
	newFields := h.extractFields(newResult)

	// 收集所有字段名
	allFields := make(map[string]bool)
	for field := range oldFields {
		allFields[field] = true
	}
	for field := range newFields {
		allFields[field] = true
	}

	// 对比每个字段
	for field := range allFields {
		oldVal, oldExists := oldFields[field]
		newVal, newExists := newFields[field]

		comp := CompareResult{
			Field:    field,
			OldValue: oldVal,
			NewValue: newVal,
		}

		switch {
		case !oldExists && newExists:
			comp.Status = "added"
		case oldExists && !newExists:
			comp.Status = "removed"
		case h.valuesEqual(oldVal, newVal):
			comp.Status = "same"
		default:
			comp.Status = "different"
			// 计算数值差异
			if oldFloat, ok1 := h.toFloat64(oldVal); ok1 {
				if newFloat, ok2 := h.toFloat64(newVal); ok2 {
					comp.Difference = newFloat - oldFloat
					if oldFloat != 0 {
						comp.RelativeDiff = (newFloat - oldFloat) / oldFloat
					}
				}
			}
		}

		comparisons = append(comparisons, comp)
	}

	return comparisons
}

// extractFields 提取结构体的字段
func (h *APIHandler) extractFields(result interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	if result == nil {
		return fields
	}

	v := reflect.ValueOf(result)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fields
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		// 使用 json tag 作为字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			// 去掉 omitempty 等选项
			if idx := len(jsonTag); idx > 0 {
				for j, c := range jsonTag {
					if c == ',' {
						idx = j
						break
					}
				}
				fields[jsonTag[:idx]] = value
			}
		} else {
			fields[field.Name] = value
		}
	}

	return fields
}

// valuesEqual 比较两个值是否相等
func (h *APIHandler) valuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// 尝试数值比较
	if aFloat, ok1 := h.toFloat64(a); ok1 {
		if bFloat, ok2 := h.toFloat64(b); ok2 {
			// 使用相对容差比较浮点数
			if aFloat == 0 && bFloat == 0 {
				return true
			}
			diff := math.Abs(aFloat - bFloat)
			maxVal := math.Max(math.Abs(aFloat), math.Abs(bFloat))
			if maxVal == 0 {
				return diff < 1e-10
			}
			return diff/maxVal < 1e-6
		}
	}

	// 使用反射比较
	return reflect.DeepEqual(a, b)
}

// toFloat64 尝试将值转换为 float64
func (h *APIHandler) toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}
