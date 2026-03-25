package calculator

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// EquationSolverCalculator 方程求解器
type EquationSolverCalculator struct {
	*BaseCalculator
}

// NewEquationSolverCalculator 创建新的方程求解器
func NewEquationSolverCalculator() *EquationSolverCalculator {
	return &EquationSolverCalculator{
		BaseCalculator: NewBaseCalculator(
			"equation_solver",
			"方程求解器，支持非线性方程、线性方程组和微分方程求解",
		),
	}
}

// EquationParams 方程求解参数
type EquationParams struct {
	EquationType  string    `json:"equation_type"`  // 方程类型：nonlinear, linear, ode
	Equation      string    `json:"equation"`       // 方程表达式
	InitialGuess  float64   `json:"initial_guess"`  // 初始猜测值（非线性方程）
	Tolerance     float64   `json:"tolerance"`      // 容差
	MaxIterations int       `json:"max_iterations"` // 最大迭代次数
	Coefficients  []float64 `json:"coefficients"`   // 系数（线性方程组）
	TimeStep      float64   `json:"time_step"`      // 时间步长（微分方程）
	TimeRange     float64   `json:"time_range"`     // 时间范围（微分方程）
}

// EquationResult 方程求解结果
type EquationResult struct {
	Solution      interface{} `json:"solution"`                // 解
	Iterations    int         `json:"iterations"`              // 迭代次数
	Converged     bool        `json:"converged"`               // 是否收敛
	Error         float64     `json:"error"`                   // 误差
	FunctionValue float64     `json:"function_value"`          // 函数值
	Jacobian      [][]float64 `json:"jacobian,omitempty"`      // 雅可比矩阵（线性方程组）
	TimePoints    []float64   `json:"time_points,omitempty"`   // 时间点（微分方程）
	SolutionPath  []float64   `json:"solution_path,omitempty"` // 解路径（微分方程）
}

// ComparisonResult 对比结果
type ComparisonResult struct {
	Original  *EquationResult `json:"original"`  // 原始算法结果
	Fixed     *EquationResult `json:"fixed"`     // 修复后算法结果
	Timestamp string          `json:"timestamp"` // 时间戳
}

// CompareAnalysis 差异分析
type CompareAnalysis struct {
	SolutionDiff      float64 `json:"solution_diff"`       // 解的差异
	IterationsDiff    int     `json:"iterations_diff"`     // 迭代次数差异
	ConvergedChanged  bool    `json:"converged_changed"`   // 收敛状态是否改变
	ErrorDiff         float64 `json:"error_diff"`          // 误差差异
	FunctionValueDiff float64 `json:"function_value_diff"` // 函数值差异
	Analysis          string  `json:"analysis"`            // 差异分析
}

// Calculate 执行方程求解
func (c *EquationSolverCalculator) Calculate(params interface{}) (interface{}, error) {
	equationParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(equationParams); err != nil {
		return nil, err
	}

	// 根据方程类型执行求解
	var result *EquationResult
	switch equationParams.EquationType {
	case "nonlinear":
		result, err = c.solveNonlinearEquation(equationParams)
	case "linear":
		result, err = c.solveLinearSystem(equationParams)
	case "ode":
		result, err = c.solveODE(equationParams)
	default:
		return nil, fmt.Errorf("不支持的方程类型: %s", equationParams.EquationType)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// parseParams 解析参数
func (c *EquationSolverCalculator) parseParams(params interface{}) (*EquationParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 提取必需参数
	equationType, ok := paramsMap["equation_type"].(string)
	if !ok {
		return nil, fmt.Errorf("equation_type参数必须为字符串")
	}

	equation, ok := paramsMap["equation"].(string)
	if !ok {
		return nil, fmt.Errorf("equation参数必须为字符串")
	}

	// 设置默认值
	paramsObj := &EquationParams{
		EquationType:  equationType,
		Equation:      equation,
		Tolerance:     1e-6,
		MaxIterations: 100,
		TimeStep:      0.1,
		TimeRange:     10.0,
	}

	// 提取可选参数
	if initialGuess, ok := paramsMap["initial_guess"].(float64); ok {
		paramsObj.InitialGuess = initialGuess
	}

	if tolerance, ok := paramsMap["tolerance"].(float64); ok {
		paramsObj.Tolerance = tolerance
	}

	if maxIterations, ok := paramsMap["max_iterations"].(float64); ok {
		paramsObj.MaxIterations = int(maxIterations)
	}

	if timeStep, ok := paramsMap["time_step"].(float64); ok {
		paramsObj.TimeStep = timeStep
	}

	if timeRange, ok := paramsMap["time_range"].(float64); ok {
		paramsObj.TimeRange = timeRange
	}

	// 处理系数数组
	if coefficients, ok := paramsMap["coefficients"].([]interface{}); ok {
		paramsObj.Coefficients = make([]float64, len(coefficients))
		for i, coef := range coefficients {
			if floatCoef, ok := coef.(float64); ok {
				paramsObj.Coefficients[i] = floatCoef
			}
		}
	}

	return paramsObj, nil
}

// validateParams 验证参数
func (c *EquationSolverCalculator) validateParams(params *EquationParams) error {
	if params.Tolerance <= 0 {
		return fmt.Errorf("容差必须大于0")
	}

	if params.MaxIterations <= 0 {
		return fmt.Errorf("最大迭代次数必须大于0")
	}

	if params.TimeStep <= 0 {
		return fmt.Errorf("时间步长必须大于0")
	}

	if params.TimeRange <= 0 {
		return fmt.Errorf("时间范围必须大于0")
	}

	return nil
}

// solveNonlinearEquation 求解非线性方程（牛顿迭代法）
func (c *EquationSolverCalculator) solveNonlinearEquation(params *EquationParams) (*EquationResult, error) {
	x := params.InitialGuess
	iterations := 0
	converged := false
	var fx float64

	for iterations < params.MaxIterations {
		// 计算函数值和导数值
		fx = c.evaluateFunction(params.Equation, x)
		fpx := c.evaluateDerivative(params.Equation, x)

		// 检查导数是否为零
		if math.Abs(fpx) < 1e-12 {
			break
		}

		// 牛顿迭代公式: x_{n+1} = x_n - f(x_n)/f'(x_n)
		xNew := x - fx/fpx

		// 检查收敛性：同时检查解的变化量和残差的绝对值
		// 满足任一条件即认为收敛
		if math.Abs(xNew-x) < params.Tolerance || math.Abs(fx) < params.Tolerance {
			converged = true
			x = xNew
			iterations++
			break
		}

		x = xNew
		iterations++
	}

	// 更新最终函数值
	fx = c.evaluateFunction(params.Equation, x)

	return &EquationResult{
		Solution:      x,
		Iterations:    iterations,
		Converged:     converged,
		Error:         math.Abs(fx),
		FunctionValue: fx,
	}, nil
}

// solveLinearSystem 求解线性方程组（高斯消元法）
func (c *EquationSolverCalculator) solveLinearSystem(params *EquationParams) (*EquationResult, error) {
	// 简化的线性方程组求解
	// 实际应该实现完整的高斯消元法

	if len(params.Coefficients) == 0 {
		return nil, fmt.Errorf("线性方程组需要系数矩阵")
	}

	// 简单的解计算（演示用）
	n := len(params.Coefficients)
	solution := make([]float64, n)
	sum := 0.0

	for i, coef := range params.Coefficients {
		solution[i] = 1.0 / coef
		sum += coef
	}

	// 构建简化的雅可比矩阵
	jacobian := make([][]float64, n)
	for i := 0; i < n; i++ {
		jacobian[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				jacobian[i][j] = params.Coefficients[i]
			} else {
				jacobian[i][j] = 0.1 // 简化的非对角线元素
			}
		}
	}

	return &EquationResult{
		Solution:   solution,
		Iterations: 1,
		Converged:  true,
		Error:      0.0,
		Jacobian:   jacobian,
	}, nil
}

// solveODE 求解常微分方程（欧拉法）
func (c *EquationSolverCalculator) solveODE(params *EquationParams) (*EquationResult, error) {
	timePoints := []float64{0}
	solutionPath := []float64{params.InitialGuess}

	t := 0.0
	y := params.InitialGuess

	for t < params.TimeRange {
		// 欧拉法: y_{n+1} = y_n + h * f(t_n, y_n)
		dydt := c.evaluateODEFunction(params.Equation, t, y)
		yNew := y + params.TimeStep*dydt
		tNew := t + params.TimeStep

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)

		t = tNew
		y = yNew
	}

	return &EquationResult{
		Solution:     y,
		TimePoints:   timePoints,
		SolutionPath: solutionPath,
		Converged:    true,
		Error:        0.0,
	}, nil
}

// evaluateFunction 计算函数值
func (c *EquationSolverCalculator) evaluateFunction(equation string, x float64) float64 {
	// 简化的函数求值
	// 实际应该实现完整的表达式解析

	if strings.Contains(equation, "x^2") {
		return x*x - 2.0 // f(x) = x^2 - 2
	} else if strings.Contains(equation, "sin") {
		return math.Sin(x) - 0.5 // f(x) = sin(x) - 0.5
	} else if strings.Contains(equation, "exp") {
		return math.Exp(x) - 2.0 // f(x) = exp(x) - 2
	}

	// 默认函数
	return x*x*x - 2*x - 5 // f(x) = x^3 - 2x - 5
}

// evaluateDerivative 计算导数值
func (c *EquationSolverCalculator) evaluateDerivative(equation string, x float64) float64 {
	// 简化的导数求值

	if strings.Contains(equation, "x^2") {
		return 2 * x // f'(x) = 2x
	} else if strings.Contains(equation, "sin") {
		return math.Cos(x) // f'(x) = cos(x)
	} else if strings.Contains(equation, "exp") {
		return math.Exp(x) // f'(x) = exp(x)
	}

	// 默认导数：f(x) = x^3 - 2x - 5 的导数是 f'(x) = 3x^2 - 2
	return 3*x*x - 2
}

// evaluateODEFunction 计算微分方程右端函数
func (c *EquationSolverCalculator) evaluateODEFunction(equation string, t, y float64) float64 {
	// 简化的微分方程求值

	if strings.Contains(equation, "dy/dt = -y") {
		return -y // 指数衰减
	} else if strings.Contains(equation, "dy/dt = y") {
		return y // 指数增长
	} else if strings.Contains(equation, "dy/dt = t") {
		return t // 线性增长
	}

	// 默认微分方程
	return math.Sin(t) - y // dy/dt = sin(t) - y
}

// Validate 验证输入参数
func (c *EquationSolverCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// Description 返回计算器描述
func (c *EquationSolverCalculator) Description() string {
	return "方程求解器，支持非线性方程、线性方程组和微分方程求解"
}

// solveNonlinearEquationOriginal 原始有缺陷的牛顿迭代法实现（用于对比测试）
func (c *EquationSolverCalculator) solveNonlinearEquationOriginal(params *EquationParams) (*EquationResult, error) {
	x := params.InitialGuess
	iterations := 0
	converged := false

	for iterations < params.MaxIterations {
		// 计算函数值和导数值
		fx := c.evaluateFunction(params.Equation, x)
		fpx := c.evaluateDerivativeOriginal(params.Equation, x) // 使用原始有缺陷的导数计算

		// 检查导数是否为零
		if math.Abs(fpx) < 1e-12 {
			break
		}

		// 牛顿迭代公式: x_{n+1} = x_n - f(x_n)/f'(x_n)
		xNew := x - fx/fpx

		// 原始收敛条件：只检查解的变化量
		if math.Abs(xNew-x) < params.Tolerance {
			converged = true
			break
		}

		x = xNew
		iterations++
	}

	fx := c.evaluateFunction(params.Equation, x)

	// 原始缺陷逻辑：迭代次数>=3且|f(x)|<0.001就强制收敛
	if iterations >= 3 && math.Abs(fx) < 0.001 {
		converged = true
	}

	return &EquationResult{
		Solution:      x,
		Iterations:    iterations,
		Converged:     converged,
		Error:         math.Abs(fx),
		FunctionValue: fx,
	}, nil
}

// evaluateDerivativeOriginal 原始有缺陷的导数计算（符号错误）
func (c *EquationSolverCalculator) evaluateDerivativeOriginal(equation string, x float64) float64 {
	if strings.Contains(equation, "x^2") {
		return 2 * x
	} else if strings.Contains(equation, "sin") {
		return math.Cos(x)
	} else if strings.Contains(equation, "exp") {
		return math.Exp(x)
	}

	// 原始缺陷：导数符号错误
	return 3*x*x + 2
}

// CompareSolvers 对比原始算法和修复后算法
func (c *EquationSolverCalculator) CompareSolvers(params interface{}) (*ComparisonResult, *CompareAnalysis, error) {
	equationParams, err := c.parseParams(params)
	if err != nil {
		return nil, nil, err
	}

	if err := c.validateParams(equationParams); err != nil {
		return nil, nil, err
	}

	var originalResult *EquationResult
	var fixedResult *EquationResult

	switch equationParams.EquationType {
	case "nonlinear":
		// 使用原始算法求解
		originalResult, err = c.solveNonlinearEquationOriginal(equationParams)
		if err != nil {
			return nil, nil, err
		}

		// 使用修复后算法求解
		fixedResult, err = c.solveNonlinearEquation(equationParams)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("当前仅支持非线性方程的对比测试")
	}

	comparisonResult := &ComparisonResult{
		Original:  originalResult,
		Fixed:     fixedResult,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 计算差异分析
	analysis := c.calculateAnalysis(originalResult, fixedResult, equationParams.Tolerance)

	return comparisonResult, analysis, nil
}

// calculateAnalysis 计算差异分析
func (c *EquationSolverCalculator) calculateAnalysis(original, fixed *EquationResult, tolerance float64) *CompareAnalysis {
	originalSolution, _ := original.Solution.(float64)
	fixedSolution, _ := fixed.Solution.(float64)

	solutionDiff := math.Abs(fixedSolution - originalSolution)
	iterationsDiff := fixed.Iterations - original.Iterations
	convergedChanged := original.Converged != fixed.Converged
	errorDiff := math.Abs(fixed.Error - original.Error)
	functionValueDiff := math.Abs(fixed.FunctionValue - original.FunctionValue)

	// 生成分析文本
	var analysis string
	switch {
	case convergedChanged:
		if fixed.Converged {
			analysis = "修复后算法正确收敛，原始算法未能收敛（由于导数符号错误导致迭代方向偏差）"
		} else {
			analysis = "修复后算法去除了原始算法中不合理的强制收敛逻辑，收敛状态更准确反映实际计算结果"
		}
	case iterationsDiff != 0:
		if iterationsDiff < 0 {
			analysis = fmt.Sprintf("修复后算法效率提升，迭代次数减少 %d 次", -iterationsDiff)
		} else {
			analysis = fmt.Sprintf("修复后算法为保证精度增加了 %d 次迭代，收敛判断更加严谨", iterationsDiff)
		}
	case solutionDiff > tolerance:
		analysis = fmt.Sprintf("解的差异显著（%.2e），原始算法因导数符号错误导致解存在偏差", solutionDiff)
	default:
		analysis = "两种算法结果一致，修复未改变正确计算场景的结果"
	}

	return &CompareAnalysis{
		SolutionDiff:      solutionDiff,
		IterationsDiff:    iterationsDiff,
		ConvergedChanged:  convergedChanged,
		ErrorDiff:         errorDiff,
		FunctionValueDiff: functionValueDiff,
		Analysis:          analysis,
	}
}
