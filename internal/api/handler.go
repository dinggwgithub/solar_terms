package api

import (
	"math"
	"math/rand"
	"net/http"
	"scientific_calc/internal/calculator"
	"scientific_calc/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type APIHandler struct {
	calculatorManager      *calculator.CalculatorManager
	calculatorManagerFixed *calculator.CalculatorManager
}

func NewAPIHandler(calculatorManager *calculator.CalculatorManager) *APIHandler {
	fixedManager := calculator.NewCalculatorManager()
	fixedManager.RegisterCalculator(calculator.CalculationTypeEquationSolver, calculator.NewEquationSolverCalculatorFixed())

	return &APIHandler{
		calculatorManager:      calculatorManager,
		calculatorManagerFixed: fixedManager,
	}
}

type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Time      string `json:"time" example:"2024-01-01T12:00:00Z"`
	Version   string `json:"version" example:"1.0.0"`
	BuildTime string `json:"build_time" example:"2024-01-01T12:00:00Z"`
}

type CalculationResponse struct {
	Success     bool        `json:"success"`
	Result      interface{} `json:"result"`
	Warnings    []string    `json:"warnings"`
	Timestamp   string      `json:"timestamp"`
	Calculation string      `json:"calculation"`
	SessionID   string      `json:"session_id,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

type CalculatorInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SystemInfoResponse struct {
	Version               string   `json:"version"`
	SupportedCalculations []string `json:"supported_calculations"`
	TotalCalculators      int      `json:"total_calculators"`
}

type CompareResponse struct {
	Success      bool                   `json:"success"`
	Calculation  string                 `json:"calculation"`
	Timestamp    string                 `json:"timestamp"`
	SessionID    string                 `json:"session_id"`
	Original     interface{}            `json:"original"`
	Fixed        interface{}            `json:"fixed"`
	Differences  map[string]interface{} `json:"differences"`
	Summary      string                 `json:"summary"`
	Improvements []string               `json:"improvements"`
}

type EquationResultFixed struct {
	Solution       interface{}   `json:"solution"`
	Iterations     int           `json:"iterations"`
	Converged      bool          `json:"converged"`
	Error          float64       `json:"error"`
	ErrorEstimate  float64       `json:"error_estimate"`
	FunctionValue  float64       `json:"function_value"`
	Jacobian       [][]float64   `json:"jacobian,omitempty"`
	TimePoints     []float64     `json:"time_points,omitempty"`
	SolutionPath   []float64     `json:"solution_path,omitempty"`
	MethodUsed     string        `json:"method_used,omitempty"`
	Stability      string        `json:"stability,omitempty"`
	GlobalError    float64       `json:"global_error,omitempty"`
	LocalError     float64       `json:"local_error,omitempty"`
	Analytical     float64       `json:"analytical,omitempty"`
	AbsoluteError  float64       `json:"absolute_error,omitempty"`
	RelativeError  float64       `json:"relative_error,omitempty"`
	StepDetails    []StepDetail  `json:"step_details,omitempty"`
}

type StepDetail struct {
	Step        int     `json:"step"`
	Time        float64 `json:"time"`
	Value       float64 `json:"value"`
	Derivative  float64 `json:"derivative"`
	LocalError  float64 `json:"local_error"`
	Cumulative  float64 `json:"cumulative_error"`
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
// @Description 执行各种类型的科学计算（原始版本）
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

	calcType, err := calculator.ParseCalculationType(req.Calculation)
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

	h.sendSuccess(c, result, warnings, req.Calculation, sessionID)
}

// CalculateFixed 修复版科学计算接口
// @Summary 执行科学计算（修复版）
// @Description 执行各种类型的科学计算（修复版），修正了ODE求解器的数值计算错误，返回严格符合数学要求的计算结果
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

	calcType, err := calculator.ParseCalculationType(req.Calculation)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
		return
	}

	sessionID := c.DefaultQuery("session_id", "")
	if sessionID == "" {
		sessionID = GenerateSessionID()
	}

	result, warnings, err := h.calculatorManagerFixed.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "计算失败: "+err.Error())
		return
	}

	h.sendSuccess(c, result, warnings, req.Calculation, sessionID)
}

// CompareSolvers 求解器对比接口
// @Summary 对比原始版与修复版求解器
// @Description 同时执行原始版和修复版计算，返回字段级的差异对比，用于验证修复效果
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

	calcType, err := calculator.ParseCalculationType(req.Calculation)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "不支持的计算类型: "+req.Calculation)
		return
	}

	sessionID := GenerateSessionID()

	originalResult, _, err := h.calculatorManager.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原始计算失败: "+err.Error())
		return
	}

	fixedResult, _, err := h.calculatorManagerFixed.CalculateWithSession(calcType, req.GetParams(), sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复版计算失败: "+err.Error())
		return
	}

	differences := h.compareResults(originalResult, fixedResult)
	summary := h.generateSummary(differences)
	improvements := h.identifyImprovements(originalResult, fixedResult)

	c.JSON(http.StatusOK, CompareResponse{
		Success:      true,
		Calculation:  req.Calculation,
		Timestamp:    time.Now().Format(time.RFC3339),
		SessionID:    sessionID,
		Original:     originalResult,
		Fixed:        fixedResult,
		Differences:  differences,
		Summary:      summary,
		Improvements: improvements,
	})
}

func (h *APIHandler) compareResults(original, fixed interface{}) map[string]interface{} {
	differences := make(map[string]interface{})

	origMap, origOk := original.(map[string]interface{})
	fixedMap, fixedOk := fixed.(map[string]interface{})

	if !origOk || !fixedOk {
		differences["type_mismatch"] = "结果类型不匹配"
		return differences
	}

	if origSolution, ok := origMap["solution"].(float64); ok {
		if fixedSolution, ok := fixedMap["solution"].(float64); ok {
			diff := origSolution - fixedSolution
			differences["solution_diff"] = diff
			differences["solution_diff_percent"] = diff / fixedSolution * 100
		}
	}

	if origError, ok := origMap["error"].(float64); ok {
		if fixedError, ok := fixedMap["error"].(float64); ok {
			differences["error_improvement"] = origError - fixedError
			if origError > 0 {
				differences["error_improvement_percent"] = (origError - fixedError) / origError * 100
			}
		}
	}

	if origPath, ok := origMap["solution_path"].([]interface{}); ok {
		if fixedPath, ok := fixedMap["solution_path"].([]interface{}); ok {
			if len(origPath) == len(fixedPath) {
				pathDiffs := make([]float64, len(origPath))
				for i := 0; i < len(origPath); i++ {
					if origVal, ok1 := origPath[i].(float64); ok1 {
						if fixedVal, ok2 := fixedPath[i].(float64); ok2 {
							pathDiffs[i] = origVal - fixedVal
						}
					}
				}
				differences["path_differences"] = pathDiffs
			} else {
				differences["path_length_diff"] = len(origPath) - len(fixedPath)
			}
		}
	}

	if origTimePoints, ok := origMap["time_points"].([]interface{}); ok {
		if fixedTimePoints, ok := fixedMap["time_points"].([]interface{}); ok {
			differences["time_points_count_original"] = len(origTimePoints)
			differences["time_points_count_fixed"] = len(fixedTimePoints)
		}
	}

	if analytical, ok := fixedMap["analytical"].(float64); ok {
		if origSolution, ok := origMap["solution"].(float64); ok {
			differences["original_vs_analytical"] = origSolution - analytical
		}
		if fixedSolution, ok := fixedMap["solution"].(float64); ok {
			differences["fixed_vs_analytical"] = fixedSolution - analytical
		}
	}

	return differences
}

func (h *APIHandler) generateSummary(differences map[string]interface{}) string {
	summary := ""

	if diff, ok := differences["solution_diff"].(float64); ok {
		if math.Abs(diff) > 0.001 {
			summary += "修复版与原版解存在显著差异; "
		} else {
			summary += "修复版与原版解基本一致; "
		}
	}

	if improvement, ok := differences["error_improvement"].(float64); ok {
		if improvement > 0 {
			summary += "修复版误差显著降低; "
		} else if improvement < 0 {
			summary += "注意：修复版误差有所增加; "
		}
	}

	if origVsAnal, ok := differences["original_vs_analytical"].(float64); ok {
		if fixedVsAnal, ok := differences["fixed_vs_analytical"].(float64); ok {
			if math.Abs(fixedVsAnal) < math.Abs(origVsAnal) {
				summary += "修复版更接近解析解; "
			}
		}
	}

	if summary == "" {
		summary = "计算结果对比完成"
	}

	return summary
}

func (h *APIHandler) identifyImprovements(original, fixed interface{}) []string {
	improvements := []string{}

	_, origOk := original.(map[string]interface{})
	fixedMap, fixedOk := fixed.(map[string]interface{})

	if !origOk || !fixedOk {
		return improvements
	}

	if _, ok := fixedMap["analytical"]; ok {
		improvements = append(improvements, "修复版提供解析解参考值")
	}

	if _, ok := fixedMap["method_used"]; ok {
		improvements = append(improvements, "修复版标明使用的数值方法")
	}

	if _, ok := fixedMap["global_error"]; ok {
		improvements = append(improvements, "修复版提供全局误差估计")
	}

	if _, ok := fixedMap["relative_error"]; ok {
		improvements = append(improvements, "修复版提供相对误差分析")
	}

	if _, ok := fixedMap["step_details"]; ok {
		improvements = append(improvements, "修复版提供详细迭代过程")
	}

	return improvements
}

// GetCalculatorInfo 获取计算器信息接口
// @Summary 获取计算器信息
// @Description 获取指定计算器的详细信息
// @Tags 计算器管理
// @Accept json
// @Produce json
// @Param calculation query string true "计算类型: solar_term, ganzhi, astronomy, starting_age, lunar, planet, star, sunrise_sunset, moon_phase, equation_solver"
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

func (h *APIHandler) sendError(c *gin.Context, code int, message string) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}

	c.JSON(code, response)
}

func (h *APIHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/health", h.HealthCheck)
	router.GET("/api/system-info", h.GetSystemInfo)

	router.POST("/api/calculate", h.Calculate)
	router.POST("/api/calculate-fixed", h.CalculateFixed)
	router.POST("/api/solver/compare", h.CompareSolvers)

	router.GET("/api/calculator-info", h.GetCalculatorInfo)
	router.GET("/api/supported-calculations", h.GetSupportedCalculations)
}

// GetSupportedCalculations 获取支持的计算类型接口
// @Summary 获取支持的计算类型
// @Description 获取系统支持的所有计算类型列表
// @Tags 计算器管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/supported-calculations [get]
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
