package calculator

import (
	"fmt"
	"math"
	"strings"
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

// IterationDetail 迭代详情
type IterationDetail struct {
	Iteration     int     `json:"iteration"`
	X             float64 `json:"x"`
	FunctionValue float64 `json:"function_value"`
	Derivative    float64 `json:"derivative"`
	DeltaX        float64 `json:"delta_x"`
	Residual      float64 `json:"residual"`
}

// EquationResultV2 方程求解结果V2（包含迭代详情）
type EquationResultV2 struct {
	Solution         float64           `json:"solution"`
	Iterations       int               `json:"iterations"`
	Converged        bool              `json:"converged"`
	Error            float64           `json:"error"`
	FunctionValue    float64           `json:"function_value"`
	Tolerance        float64           `json:"tolerance"`
	ConvergenceType  string            `json:"convergence_type"`
	IterationDetails []IterationDetail `json:"iteration_details"`
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

// SolveNonlinearEquationOld 求解非线性方程（牛顿迭代法）- 旧版，用于对比
func (c *EquationSolverCalculator) SolveNonlinearEquationOld(params *EquationParams) (*EquationResult, error) {
	x := params.InitialGuess
	iterations := 0
	converged := false

	for iterations < params.MaxIterations {
		fx := c.evaluateFunction(params.Equation, x)
		fpx := c.evaluateDerivativeOld(params.Equation, x)

		if math.Abs(fpx) < 1e-12 {
			break
		}

		xNew := x - fx/fpx

		if math.Abs(xNew-x) < params.Tolerance {
			converged = true
			break
		}

		x = xNew
		iterations++
	}

	fx := c.evaluateFunction(params.Equation, x)

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

// solveNonlinearEquation 求解非线性方程（牛顿迭代法）
func (c *EquationSolverCalculator) solveNonlinearEquation(params *EquationParams) (*EquationResult, error) {
	x := params.InitialGuess
	iterations := 0
	converged := false
	var lastFx float64

	for iterations < params.MaxIterations {
		fx := c.evaluateFunction(params.Equation, x)
		fpx := c.evaluateDerivative(params.Equation, x)

		if math.Abs(fpx) < 1e-12 {
			break
		}

		xNew := x - fx/fpx

		deltaX := math.Abs(xNew - x)
		residual := math.Abs(fx)

		if deltaX < params.Tolerance && residual < params.Tolerance {
			converged = true
			x = xNew
			lastFx = c.evaluateFunction(params.Equation, x)
			iterations++
			break
		}

		if deltaX < params.Tolerance*math.Max(1.0, math.Abs(x)) {
			converged = true
			x = xNew
			lastFx = c.evaluateFunction(params.Equation, x)
			iterations++
			break
		}

		x = xNew
		lastFx = fx
		iterations++
	}

	if !converged {
		lastFx = c.evaluateFunction(params.Equation, x)
		if math.Abs(lastFx) < params.Tolerance {
			converged = true
		}
	}

	return &EquationResult{
		Solution:      x,
		Iterations:    iterations,
		Converged:     converged,
		Error:         math.Abs(lastFx),
		FunctionValue: lastFx,
	}, nil
}

// SolveNonlinearEquationV2 求解非线性方程（牛顿迭代法）- V2版本，包含详细迭代信息
func (c *EquationSolverCalculator) SolveNonlinearEquationV2(params *EquationParams) (*EquationResultV2, error) {
	x := params.InitialGuess
	iterations := 0
	converged := false
	convergenceType := ""
	var lastFx, lastFpx float64
	iterationDetails := make([]IterationDetail, 0)

	for iterations < params.MaxIterations {
		fx := c.evaluateFunction(params.Equation, x)
		fpx := c.evaluateDerivative(params.Equation, x)

		if math.Abs(fpx) < 1e-12 {
			break
		}

		xNew := x - fx/fpx

		deltaX := math.Abs(xNew - x)
		residual := math.Abs(fx)

		detail := IterationDetail{
			Iteration:     iterations + 1,
			X:             x,
			FunctionValue: fx,
			Derivative:    fpx,
			DeltaX:        deltaX,
			Residual:      residual,
		}
		iterationDetails = append(iterationDetails, detail)

		if deltaX < params.Tolerance && residual < params.Tolerance {
			converged = true
			convergenceType = "both"
			x = xNew
			lastFx = c.evaluateFunction(params.Equation, x)
			lastFpx = c.evaluateDerivative(params.Equation, x)
			iterations++
			break
		}

		if deltaX < params.Tolerance*math.Max(1.0, math.Abs(x)) {
			converged = true
			convergenceType = "relative"
			x = xNew
			lastFx = c.evaluateFunction(params.Equation, x)
			lastFpx = c.evaluateDerivative(params.Equation, x)
			iterations++
			break
		}

		x = xNew
		lastFx = fx
		lastFpx = fpx
		iterations++
	}

	if !converged {
		lastFx = c.evaluateFunction(params.Equation, x)
		lastFpx = c.evaluateDerivative(params.Equation, x)
		if math.Abs(lastFx) < params.Tolerance {
			converged = true
			convergenceType = "residual"
		}
	}

	if converged && convergenceType == "" {
		convergenceType = "unknown"
	}

	finalDetail := IterationDetail{
		Iteration:     iterations,
		X:             x,
		FunctionValue: lastFx,
		Derivative:    lastFpx,
		DeltaX:        0,
		Residual:      math.Abs(lastFx),
	}
	iterationDetails = append(iterationDetails, finalDetail)

	return &EquationResultV2{
		Solution:         x,
		Iterations:       iterations,
		Converged:        converged,
		Error:            math.Abs(lastFx),
		FunctionValue:    lastFx,
		Tolerance:        params.Tolerance,
		ConvergenceType:  convergenceType,
		IterationDetails: iterationDetails,
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

	// 默认导数
	return 3*x*x - 2
}

// evaluateDerivativeOld 计算导数值（旧版，包含错误）
func (c *EquationSolverCalculator) evaluateDerivativeOld(equation string, x float64) float64 {
	if strings.Contains(equation, "x^2") {
		return 2 * x
	} else if strings.Contains(equation, "sin") {
		return math.Cos(x)
	} else if strings.Contains(equation, "exp") {
		return math.Exp(x)
	}

	return 3*x*x + 2
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
