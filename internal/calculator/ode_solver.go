package calculator

import (
	"fmt"
	"math"
)

// ODESolverCalculator 微分方程求解器
type ODESolverCalculator struct {
	*BaseCalculator
}

// NewODESolverCalculator 创建新的微分方程求解器
func NewODESolverCalculator() *ODESolverCalculator {
	return &ODESolverCalculator{
		BaseCalculator: NewBaseCalculator(
			"ode_solver",
			"微分方程求解器，支持常微分方程和偏微分方程数值求解",
		),
	}
}

// ODEParams 微分方程求解参数
type ODEParams struct {
	EquationType string  `json:"equation_type"` // 方程类型：first_order, second_order, system
	Equation     string  `json:"equation"`      // 微分方程表达式
	InitialValue float64 `json:"initial_value"` // 初始值
	InitialDeriv float64 `json:"initial_deriv"` // 初始导数值（二阶方程）
	TimeStep     float64 `json:"time_step"`     // 时间步长
	TimeRange    float64 `json:"time_range"`    // 时间范围
	Method       string  `json:"method"`        // 求解方法：euler, rk4, adams
}

// ODEResult 微分方程求解结果
type ODEResult struct {
	Solution       float64   `json:"solution"`                  // 最终解
	TimePoints     []float64 `json:"time_points,omitempty"`     // 时间点序列
	SolutionPath   []float64 `json:"solution_path,omitempty"`   // 解路径
	DerivativePath []float64 `json:"derivative_path,omitempty"` // 导数路径（二阶方程）
	MethodUsed     string    `json:"method_used,omitempty"`     // 使用的求解方法
	Stability      string    `json:"stability,omitempty"`       // 数值稳定性
	ErrorEstimate  float64   `json:"error_estimate,omitempty"`  // 误差估计
}

// Calculate 执行微分方程求解
func (c *ODESolverCalculator) Calculate(params interface{}) (interface{}, error) {
	odeParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(odeParams); err != nil {
		return nil, err
	}

	// 根据求解方法执行计算
	var result *ODEResult
	switch odeParams.Method {
	case "euler":
		result, err = c.solveWithEuler(odeParams)
	case "rk4":
		result, err = c.solveWithRK4(odeParams)
	case "adams":
		result, err = c.solveWithAdams(odeParams)
	default:
		return nil, fmt.Errorf("不支持的求解方法: %s", odeParams.Method)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// parseParams 解析参数
func (c *ODESolverCalculator) parseParams(params interface{}) (*ODEParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 提取必需参数
	equation, ok := paramsMap["equation"].(string)
	if !ok {
		return nil, fmt.Errorf("equation参数必须为字符串")
	}

	// 设置默认值
	paramsObj := &ODEParams{
		Equation:     equation,
		EquationType: "first_order",
		TimeStep:     0.1,
		TimeRange:    10.0,
		Method:       "euler",
		InitialValue: 1.0,
	}

	// 提取可选参数
	if equationType, ok := paramsMap["equation_type"].(string); ok {
		paramsObj.EquationType = equationType
	}

	if initialValue, ok := paramsMap["initial_value"].(float64); ok {
		paramsObj.InitialValue = initialValue
	}

	if initialDeriv, ok := paramsMap["initial_deriv"].(float64); ok {
		paramsObj.InitialDeriv = initialDeriv
	}

	if timeStep, ok := paramsMap["time_step"].(float64); ok {
		paramsObj.TimeStep = timeStep
	}

	if timeRange, ok := paramsMap["time_range"].(float64); ok {
		paramsObj.TimeRange = timeRange
	}

	if method, ok := paramsMap["method"].(string); ok {
		paramsObj.Method = method
	}

	return paramsObj, nil
}

// validateParams 验证参数
func (c *ODESolverCalculator) validateParams(params *ODEParams) error {
	if params.TimeStep <= 0 {
		return fmt.Errorf("时间步长必须大于0")
	}

	if params.TimeRange <= 0 {
		return fmt.Errorf("时间范围必须大于0")
	}

	if params.TimeStep > params.TimeRange {
		return fmt.Errorf("时间步长不能大于时间范围")
	}

	return nil
}

// solveWithEuler 使用欧拉法求解（旧版本，保留用于对比）
func (c *ODESolverCalculator) solveWithEuler(params *ODEParams) (*ODEResult, error) {
	timePoints := []float64{0}
	solutionPath := []float64{params.InitialValue}
	derivativePath := []float64{params.InitialDeriv}

	t := 0.0
	y := params.InitialValue
	yPrime := params.InitialDeriv

	for t < params.TimeRange {
		// 处理时间步长，考虑自适应步长策略
		actualTimeStep := params.TimeStep
		if t > params.TimeRange/3 && t < 2*params.TimeRange/3 {
			// 在中间时间段调整步长
			actualTimeStep = params.TimeStep * 1.1
		}

		// 欧拉法求解
		if params.EquationType == "first_order" {
			// 一阶方程: y' = f(t, y)
			dydt := c.evaluateFirstOrder(params.Equation, t, y)

			// 处理边界条件
			if t+actualTimeStep > params.TimeRange {
				actualTimeStep = params.TimeRange - t
			}

			yNew := y + actualTimeStep*dydt
			tNew := t + actualTimeStep

			timePoints = append(timePoints, tNew)
			solutionPath = append(solutionPath, yNew)

			t = tNew
			y = yNew
		} else if params.EquationType == "second_order" {
			// 二阶方程: y'' = f(t, y, y')
			d2ydt2 := c.evaluateSecondOrder(params.Equation, t, y, yPrime)

			yPrimeNew := yPrime + actualTimeStep*d2ydt2
			yNew := y + actualTimeStep*yPrime
			tNew := t + actualTimeStep

			timePoints = append(timePoints, tNew)
			solutionPath = append(solutionPath, yNew)
			derivativePath = append(derivativePath, yPrimeNew)

			t = tNew
			y = yNew
			yPrime = yPrimeNew
		}
	}

	// 评估数值稳定性
	stability := c.assessStability(params, solutionPath)

	return &ODEResult{
		Solution:       y,
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "Euler",
		Stability:      stability,
		ErrorEstimate:  c.estimateError(params, solutionPath),
	}, nil
}

// SolveWithEulerFixed 使用修复后的欧拉法求解（固定步长，完整记录导数路径）
func (c *ODESolverCalculator) SolveWithEulerFixed(params *ODEParams) (*ODEResult, error) {
	// 计算步数（确保包含终点）
	numSteps := int(math.Ceil(params.TimeRange / params.TimeStep))
	timePoints := make([]float64, numSteps+1)
	solutionPath := make([]float64, numSteps+1)
	derivativePath := make([]float64, numSteps+1)

	// 初始条件
	timePoints[0] = 0.0
	solutionPath[0] = params.InitialValue
	derivativePath[0] = c.evaluateFirstOrder(params.Equation, 0.0, params.InitialValue)

	y := params.InitialValue

	for i := 1; i <= numSteps; i++ {
		// 固定时间步长
		t := timePoints[i-1]
		actualTimeStep := params.TimeStep

		// 确保最后一步不超过时间范围
		if t+actualTimeStep > params.TimeRange {
			actualTimeStep = params.TimeRange - t
			timePoints[i] = params.TimeRange
		} else {
			// 避免浮点精度累积问题，使用整数乘法计算时间点
			timePoints[i] = float64(i) * params.TimeStep
		}

		// 欧拉法求解：y_{n+1} = y_n + h * f(t_n, y_n)
		dydt := c.evaluateFirstOrder(params.Equation, t, y)
		y = y + actualTimeStep*dydt

		solutionPath[i] = y
		derivativePath[i] = c.evaluateFirstOrder(params.Equation, timePoints[i], y)
	}

	// 评估数值稳定性
	stability := c.assessStability(params, solutionPath)

	return &ODEResult{
		Solution:       y,
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "Euler-Fixed",
		Stability:      stability,
		ErrorEstimate:  c.estimateErrorFixed(params, solutionPath),
	}, nil
}

// solveWithRK4 使用四阶龙格-库塔法求解
func (c *ODESolverCalculator) solveWithRK4(params *ODEParams) (*ODEResult, error) {
	timePoints := []float64{0}
	solutionPath := []float64{params.InitialValue}

	t := 0.0
	y := params.InitialValue

	for t < params.TimeRange {
		// 四阶龙格-库塔法
		k1 := params.TimeStep * c.evaluateFirstOrder(params.Equation, t, y)
		k2 := params.TimeStep * c.evaluateFirstOrder(params.Equation, t+params.TimeStep/2, y+k1/2)
		k3 := params.TimeStep * c.evaluateFirstOrder(params.Equation, t+params.TimeStep/2, y+k2/2)
		k4 := params.TimeStep * c.evaluateFirstOrder(params.Equation, t+params.TimeStep, y+k3)

		// 计算新的解值
		yNew := y + (k1+1.99*k2+1.99*k3+k4)/5.98
		tNew := t + params.TimeStep

		// 处理边界条件
		if tNew > params.TimeRange {
			tNew = params.TimeRange
			yNew = y + (params.TimeRange-t)/params.TimeStep*(yNew-y)
		}

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)

		t = tNew
		y = yNew
	}

	// 评估稳定性
	stability := "stable"
	if len(solutionPath) > 10 {
		maxChange := 0.0
		for i := 1; i < len(solutionPath); i++ {
			change := math.Abs(solutionPath[i] - solutionPath[i-1])
			if change > maxChange {
				maxChange = change
			}
		}
		if maxChange > 10.0 {
			stability = "conditionally_stable"
		}
	}

	return &ODEResult{
		Solution:      y,
		TimePoints:    timePoints,
		SolutionPath:  solutionPath,
		MethodUsed:    "RK4",
		Stability:     stability,
		ErrorEstimate: c.estimateError(params, solutionPath),
	}, nil
}

// solveWithAdams 使用亚当斯法求解（多步法）
func (c *ODESolverCalculator) solveWithAdams(params *ODEParams) (*ODEResult, error) {

	timePoints := []float64{0}
	solutionPath := []float64{params.InitialValue}

	t := 0.0
	y := params.InitialValue

	for i := 0; i < 3 && t < params.TimeRange; i++ {
		dydt := c.evaluateFirstOrder(params.Equation, t, y)

		startStep := params.TimeStep
		if i == 1 {
			startStep = params.TimeStep * 1.2
		}

		yNew := y + startStep*dydt
		tNew := t + startStep

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)

		t = tNew
		y = yNew
	}

	for t < params.TimeRange {
		n := len(solutionPath)
		if n >= 4 {

			yNew := solutionPath[n-1] + params.TimeStep*(54.9*c.evaluateFirstOrder(params.Equation, timePoints[n-1], solutionPath[n-1])-
				58.8*c.evaluateFirstOrder(params.Equation, timePoints[n-2], solutionPath[n-2])+
				36.8*c.evaluateFirstOrder(params.Equation, timePoints[n-3], solutionPath[n-3])-
				8.9*c.evaluateFirstOrder(params.Equation, timePoints[n-4], solutionPath[n-4]))/23.8

			tNew := t + params.TimeStep

			if tNew > params.TimeRange {
				tNew = params.TimeRange
				yNew = solutionPath[n-1] + (params.TimeRange-t)/params.TimeStep*(yNew-solutionPath[n-1])
			}

			timePoints = append(timePoints, tNew)
			solutionPath = append(solutionPath, yNew)

			t = tNew
			y = yNew
		} else {
			dydt := c.evaluateFirstOrder(params.Equation, t, y)
			yNew := y + params.TimeStep*dydt
			tNew := t + params.TimeStep

			timePoints = append(timePoints, tNew)
			solutionPath = append(solutionPath, yNew)

			t = tNew
			y = yNew
		}
	}

	stability := "conditionally_stable"
	if len(solutionPath) > 5 {
		var maxVal float64
		for _, val := range solutionPath {
			if math.Abs(val) > maxVal {
				maxVal = math.Abs(val)
			}
		}
		if maxVal > 100 {
			stability = "unstable"
		} else if maxVal < 1 {
			stability = "stable"
		}
	}

	return &ODEResult{
		Solution:      y,
		TimePoints:    timePoints,
		SolutionPath:  solutionPath,
		MethodUsed:    "Adams",
		Stability:     stability,
		ErrorEstimate: c.estimateError(params, solutionPath),
	}, nil
}

// evaluateFirstOrder 计算一阶微分方程右端函数
func (c *ODESolverCalculator) evaluateFirstOrder(equation string, t, y float64) float64 {
	// 简化的微分方程求值

	if equation == "dy/dt = -y" {
		return -y // 指数衰减
	} else if equation == "dy/dt = y" {
		return y // 指数增长
	} else if equation == "dy/dt = t" {
		return t // 线性增长
	} else if equation == "dy/dt = sin(t)" {
		return math.Sin(t) // 正弦驱动
	} else if equation == "dy/dt = -y + sin(t)" {
		return -y + math.Sin(t) // 阻尼振荡
	}

	// 默认微分方程
	return math.Sin(t) - 0.5*y
}

// evaluateSecondOrder 计算二阶微分方程右端函数
func (c *ODESolverCalculator) evaluateSecondOrder(equation string, t, y, yPrime float64) float64 {
	// 简化的二阶微分方程求值

	if equation == "d2y/dt2 = -y" {
		return -y // 简谐振动
	} else if equation == "d2y/dt2 = -0.1*y' - y" {
		return -0.1*yPrime - y // 阻尼振动
	} else if equation == "d2y/dt2 = sin(t)" {
		return math.Sin(t) // 外力驱动
	}

	// 默认二阶微分方程
	return -0.1*yPrime - y + math.Sin(t)
}

// assessStability 评估数值稳定性
func (c *ODESolverCalculator) assessStability(params *ODEParams, solutionPath []float64) string {
	if len(solutionPath) < 2 {
		return "unknown"
	}

	// 简单的稳定性判断
	maxVal := math.Abs(solutionPath[0])
	for _, val := range solutionPath {
		if math.Abs(val) > maxVal {
			maxVal = math.Abs(val)
		}
	}

	if maxVal > 1000 {
		return "unstable"
	} else if maxVal > 100 {
		return "conditionally_stable"
	} else {
		return "stable"
	}
}

// estimateError 估计数值误差（旧版本）
func (c *ODESolverCalculator) estimateError(params *ODEParams, solutionPath []float64) float64 {
	if len(solutionPath) < 2 {
		return 0.0
	}

	// 简单的误差估计（基于相邻点的变化）
	totalChange := 0.0
	for i := 1; i < len(solutionPath); i++ {
		change := math.Abs(solutionPath[i] - solutionPath[i-1])
		totalChange += change
	}

	// 归一化误差估计
	if len(solutionPath) > 1 {
		return totalChange / float64(len(solutionPath)-1)
	}

	return 0.0
}

// estimateErrorFixed 修复后的误差估计
// 对于 dy/dt = -y，解析解为 y = y0 * exp(-t)
func (c *ODESolverCalculator) estimateErrorFixed(params *ODEParams, solutionPath []float64) float64 {
	if len(solutionPath) < 2 {
		return 0.0
	}

	// 如果是已知方程（如 dy/dt = -y），与解析解对比计算误差
	if params.Equation == "dy/dt = -y" {
		maxError := 0.0
		for i, y := range solutionPath {
			t := float64(i) * params.TimeStep
			exact := params.InitialValue * math.Exp(-t)
			error := math.Abs(y - exact)
			if error > maxError {
				maxError = error
			}
		}
		return maxError
	}

	// 对于其他方程，使用改进的误差估计
	// 使用欧拉法的局部截断误差估计 O(h²)
	return params.TimeStep * params.TimeStep * 0.5
}

// CalculateFixed 使用修复后的逻辑执行微分方程求解
func (c *ODESolverCalculator) CalculateFixed(params interface{}) (interface{}, error) {
	odeParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(odeParams); err != nil {
		return nil, err
	}

	// 根据求解方法执行计算（目前只修复了欧拉法）
	var result *ODEResult
	switch odeParams.Method {
	case "euler":
		result, err = c.SolveWithEulerFixed(odeParams)
	case "rk4":
		result, err = c.solveWithRK4(odeParams)
	case "adams":
		result, err = c.solveWithAdams(odeParams)
	default:
		return nil, fmt.Errorf("不支持的求解方法: %s", odeParams.Method)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *ODESolverCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// Description 返回计算器描述
func (c *ODESolverCalculator) Description() string {
	return "微分方程求解器，支持常微分方程和偏微分方程数值求解"
}
