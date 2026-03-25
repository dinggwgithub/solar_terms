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

	// 方程求解器V2接口（修复版）
	router.POST("/api/solver/v2", h.SolverV2)

	// 方程求解器对比接口
	router.POST("/api/solver/compare", h.SolverCompare)

	// 计算器管理接口
	router.GET("/api/calculator-info", h.GetCalculatorInfo)

	// 支持接口
	router.GET("/api/supported-calculations", h.GetSupportedCalculations)
}

// SolverV2Request V2求解器请求（兼容 /api/calculate 格式）
type SolverV2Request struct {
	Calculation string                 `json:"calculation" example:"equation_solver"`
	Params      map[string]interface{} `json:"params"`
}

// SolverV2 方程求解器V2接口（修复版）
// @Summary 方程求解器V2（修复版）
// @Description 使用修正后的牛顿迭代法求解非线性方程，包含详细迭代信息
// @Tags 方程求解
// @Accept json
// @Produce json
// @Param request body SolverV2Request true "求解请求参数"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/v2 [post]
func (h *APIHandler) SolverV2(c *gin.Context) {
	var req SolverV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if req.Params == nil {
		h.sendError(c, http.StatusBadRequest, "缺少params参数")
		return
	}

	equationType, _ := req.Params["equation_type"].(string)
	if equationType == "" {
		equationType = "nonlinear"
	}

	equation, _ := req.Params["equation"].(string)
	if equation == "" {
		h.sendError(c, http.StatusBadRequest, "缺少equation参数")
		return
	}

	initialGuess := 2.0
	if v, ok := req.Params["initial_guess"].(float64); ok {
		initialGuess = v
	}

	tolerance := 1e-6
	if v, ok := req.Params["tolerance"].(float64); ok {
		tolerance = v
	}

	maxIterations := 100
	if v, ok := req.Params["max_iterations"].(float64); ok {
		maxIterations = int(v)
	}

	calc, exists := h.calculatorManager.GetCalculator(calculator.CalculationTypeEquationSolver)
	if !exists {
		h.sendError(c, http.StatusInternalServerError, "方程求解器未注册")
		return
	}

	solver := calc.(*calculator.EquationSolverCalculator)
	params := &calculator.EquationParams{
		EquationType:  equationType,
		Equation:      equation,
		InitialGuess:  initialGuess,
		Tolerance:     tolerance,
		MaxIterations: maxIterations,
	}

	result, err := solver.SolveNonlinearEquationV2(params)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "求解失败: "+err.Error())
		return
	}

	sessionID := GenerateSessionID()
	h.sendSuccess(c, result, nil, "equation_solver_v2", sessionID)
}

// CompareRequest 对比请求（兼容 /api/calculate 格式）
type CompareRequest struct {
	Calculation string                 `json:"calculation" example:"equation_solver"`
	Params      map[string]interface{} `json:"params"`
}

// CompareResult 对比结果
type CompareResult struct {
	Before   AfterCompare `json:"before"`
	After    AfterCompare `json:"after"`
	Diff     Difference   `json:"diff"`
	Analysis string       `json:"analysis"`
}

// AfterCompare 修复前后对比数据
type AfterCompare struct {
	Solution      float64 `json:"solution"`
	Iterations    int     `json:"iterations"`
	Converged     bool    `json:"converged"`
	Error         float64 `json:"error"`
	FunctionValue float64 `json:"function_value"`
}

// Difference 差异数据
type Difference struct {
	SolutionDiff      float64 `json:"solution_diff"`
	IterationsDiff    int     `json:"iterations_diff"`
	ConvergedChanged  bool    `json:"converged_changed"`
	ErrorDiff         float64 `json:"error_diff"`
	FunctionValueDiff float64 `json:"function_value_diff"`
}

// SolverCompare 方程求解器对比接口
// @Summary 方程求解器对比
// @Description 对比修复前后的求解结果，输出结构化差异报告
// @Tags 方程求解
// @Accept json
// @Produce json
// @Param request body CompareRequest true "对比请求参数"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *APIHandler) SolverCompare(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if req.Params == nil {
		h.sendError(c, http.StatusBadRequest, "缺少params参数")
		return
	}

	equationType, _ := req.Params["equation_type"].(string)
	if equationType == "" {
		equationType = "nonlinear"
	}

	equation, _ := req.Params["equation"].(string)
	if equation == "" {
		h.sendError(c, http.StatusBadRequest, "缺少equation参数")
		return
	}

	initialGuess := 2.0
	if v, ok := req.Params["initial_guess"].(float64); ok {
		initialGuess = v
	}

	tolerance := 1e-6
	if v, ok := req.Params["tolerance"].(float64); ok {
		tolerance = v
	}

	maxIterations := 100
	if v, ok := req.Params["max_iterations"].(float64); ok {
		maxIterations = int(v)
	}

	calc, exists := h.calculatorManager.GetCalculator(calculator.CalculationTypeEquationSolver)
	if !exists {
		h.sendError(c, http.StatusInternalServerError, "方程求解器未注册")
		return
	}

	solver := calc.(*calculator.EquationSolverCalculator)
	params := &calculator.EquationParams{
		EquationType:  equationType,
		Equation:      equation,
		InitialGuess:  initialGuess,
		Tolerance:     tolerance,
		MaxIterations: maxIterations,
	}

	oldResult, err := solver.SolveNonlinearEquationOld(params)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原版求解失败: "+err.Error())
		return
	}

	newResult, err := solver.SolveNonlinearEquationV2(params)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "V2求解失败: "+err.Error())
		return
	}

	oldSol, _ := oldResult.Solution.(float64)
	newSol := newResult.Solution

	before := AfterCompare{
		Solution:      oldSol,
		Iterations:    oldResult.Iterations,
		Converged:     oldResult.Converged,
		Error:         oldResult.Error,
		FunctionValue: oldResult.FunctionValue,
	}

	after := AfterCompare{
		Solution:      newSol,
		Iterations:    newResult.Iterations,
		Converged:     newResult.Converged,
		Error:         newResult.Error,
		FunctionValue: newResult.FunctionValue,
	}

	diff := Difference{
		SolutionDiff:      newSol - oldSol,
		IterationsDiff:    newResult.Iterations - oldResult.Iterations,
		ConvergedChanged:  oldResult.Converged != newResult.Converged,
		ErrorDiff:         newResult.Error - oldResult.Error,
		FunctionValueDiff: newResult.FunctionValue - oldResult.FunctionValue,
	}

	analysis := generateAnalysis(before, after, diff)

	result := CompareResult{
		Before:   before,
		After:    after,
		Diff:     diff,
		Analysis: analysis,
	}

	sessionID := GenerateSessionID()
	h.sendSuccess(c, result, nil, "equation_solver_compare", sessionID)
}

func generateAnalysis(before, after AfterCompare, diff Difference) string {
	analysis := "修复前后对比分析:\n"

	analysis += fmt.Sprintf("1. 解的差异: %.10f (相对误差: %.2e)\n", diff.SolutionDiff, math.Abs(diff.SolutionDiff)/math.Abs(after.Solution))

	if diff.IterationsDiff != 0 {
		analysis += fmt.Sprintf("2. 迭代次数变化: %d -> %d (%+d次)\n", before.Iterations, after.Iterations, diff.IterationsDiff)
	} else {
		analysis += fmt.Sprintf("2. 迭代次数: 无变化 (%d次)\n", after.Iterations)
	}

	if diff.ConvergedChanged {
		analysis += fmt.Sprintf("3. 收敛状态: %v -> %v (已修复)\n", before.Converged, after.Converged)
	} else {
		analysis += fmt.Sprintf("3. 收敛状态: 无变化 (%v)\n", after.Converged)
	}

	analysis += fmt.Sprintf("4. 残差改进: %.2e -> %.2e\n", before.FunctionValue, after.FunctionValue)

	analysis += "\n修复内容:\n"
	analysis += "- 修正导数计算错误: 3*x*x + 2 -> 3*x*x - 2\n"
	analysis += "- 完善收敛判定: 双重收敛条件(deltaX & residual)\n"
	analysis += "- 移除硬编码阈值，使用用户指定的tolerance\n"

	return analysis
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
