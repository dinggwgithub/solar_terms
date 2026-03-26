package calculator

import (
	"fmt"
	"math"
)

// ODESolverFixed 修复后的微分方程求解器
type ODESolverFixed struct{}

// ODEParamsFixed 微分方程求解参数（与原始版本兼容）
type ODEParamsFixed struct {
	EquationType string  `json:"equation_type"` // 方程类型：first_order, second_order, system
	Equation     string  `json:"equation"`      // 微分方程表达式
	InitialValue float64 `json:"initial_value"` // 初始值
	InitialDeriv float64 `json:"initial_deriv"` // 初始导数值（二阶方程）
	TimeStep     float64 `json:"time_step"`     // 时间步长
	TimeRange    float64 `json:"time_range"`    // 时间范围
	Method       string  `json:"method"`        // 求解方法：euler, rk4, adams
}

// ODEResultFixed 修复后的微分方程求解结果
type ODEResultFixed struct {
	Solution       float64   `json:"solution"`                  // 最终解
	TimePoints     []float64 `json:"time_points,omitempty"`     // 时间点序列
	SolutionPath   []float64 `json:"solution_path,omitempty"`   // 解路径
	DerivativePath []float64 `json:"derivative_path,omitempty"` // 导数路径
	MethodUsed     string    `json:"method_used,omitempty"`     // 使用的求解方法
	Stability      string    `json:"stability,omitempty"`       // 数值稳定性
	ErrorEstimate  float64   `json:"error_estimate,omitempty"`  // 误差估计（与精确解对比）
}

// ParseODEParams 解析ODE参数
func ParseODEParams(params interface{}) (*ODEParamsFixed, error) {
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
	paramsObj := &ODEParamsFixed{
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

// ValidateODEParams 验证ODE参数
func ValidateODEParams(params *ODEParamsFixed) error {
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

// SolveODEWithEulerFixed 使用修复后的欧拉法求解
func SolveODEWithEulerFixed(params *ODEParamsFixed) (*ODEResultFixed, error) {
	// 计算步数：确保包含终点
	nSteps := int(math.Round(params.TimeRange / params.TimeStep))
	if nSteps < 1 {
		nSteps = 1
	}

	// 初始化数组
	timePoints := make([]float64, nSteps+1)
	solutionPath := make([]float64, nSteps+1)
	derivativePath := make([]float64, nSteps+1)

	// 初始条件
	timePoints[0] = 0.0
	solutionPath[0] = params.InitialValue

	// 计算初始导数
	if params.EquationType == "first_order" {
		derivativePath[0] = evaluateFirstOrderFixed(params.Equation, 0.0, params.InitialValue)
	} else {
		derivativePath[0] = params.InitialDeriv
	}

	// 欧拉法迭代（固定步长）
	for i := 0; i < nSteps; i++ {
		t := timePoints[i]
		y := solutionPath[i]

		// 计算当前点的导数
		var dydt float64
		if params.EquationType == "first_order" {
			dydt = evaluateFirstOrderFixed(params.Equation, t, y)
		} else {
			// 二阶方程：y'' = f(t, y, y')，需要特殊处理
			dydt = evaluateSecondOrderFixed(params.Equation, t, y, derivativePath[i])
		}

		// 计算下一个时间点（确保最后一个点精确到达 TimeRange）
		tNext := float64(i+1) * params.TimeStep
		if i == nSteps-1 {
			tNext = params.TimeRange
		}
		actualStep := tNext - t

		// 欧拉法：y_{n+1} = y_n + h * f(t_n, y_n)
		yNext := y + actualStep*dydt

		// 记录结果
		timePoints[i+1] = tNext
		solutionPath[i+1] = yNext

		// 计算下一个点的导数
		if params.EquationType == "first_order" {
			derivativePath[i+1] = evaluateFirstOrderFixed(params.Equation, tNext, yNext)
		} else {
			// 二阶方程需要更新一阶导数
			if i < nSteps {
				d2ydt2 := evaluateSecondOrderFixed(params.Equation, t, y, derivativePath[i])
				derivativePath[i+1] = derivativePath[i] + actualStep*d2ydt2
			}
		}
	}

	// 评估数值稳定性
	stability := assessStabilityFixed(params, solutionPath)

	// 计算误差估计（与精确解对比）
	errorEstimate := estimateErrorFixed(params, solutionPath, timePoints)

	return &ODEResultFixed{
		Solution:       solutionPath[nSteps],
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "Euler",
		Stability:      stability,
		ErrorEstimate:  errorEstimate,
	}, nil
}

// SolveODEWithRK4Fixed 使用修复后的四阶龙格-库塔法求解
func SolveODEWithRK4Fixed(params *ODEParamsFixed) (*ODEResultFixed, error) {
	// 计算步数
	nSteps := int(math.Round(params.TimeRange / params.TimeStep))
	if nSteps < 1 {
		nSteps = 1
	}

	// 初始化数组
	timePoints := make([]float64, nSteps+1)
	solutionPath := make([]float64, nSteps+1)
	derivativePath := make([]float64, nSteps+1)

	// 初始条件
	timePoints[0] = 0.0
	solutionPath[0] = params.InitialValue
	derivativePath[0] = evaluateFirstOrderFixed(params.Equation, 0.0, params.InitialValue)

	// RK4 迭代
	for i := 0; i < nSteps; i++ {
		t := timePoints[i]
		y := solutionPath[i]
		h := params.TimeStep

		// 计算下一个时间点
		tNext := float64(i+1) * params.TimeStep
		if i == nSteps-1 {
			tNext = params.TimeRange
		}
		h = tNext - t

		// RK4 系数（标准公式）
		k1 := h * evaluateFirstOrderFixed(params.Equation, t, y)
		k2 := h * evaluateFirstOrderFixed(params.Equation, t+h/2, y+k1/2)
		k3 := h * evaluateFirstOrderFixed(params.Equation, t+h/2, y+k2/2)
		k4 := h * evaluateFirstOrderFixed(params.Equation, t+h, y+k3)

		// 标准 RK4 公式：y_{n+1} = y_n + (k1 + 2*k2 + 2*k3 + k4) / 6
		yNext := y + (k1+2*k2+2*k3+k4)/6

		timePoints[i+1] = tNext
		solutionPath[i+1] = yNext
		derivativePath[i+1] = evaluateFirstOrderFixed(params.Equation, tNext, yNext)
	}

	stability := assessStabilityFixed(params, solutionPath)
	errorEstimate := estimateErrorFixed(params, solutionPath, timePoints)

	return &ODEResultFixed{
		Solution:       solutionPath[nSteps],
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "RK4",
		Stability:      stability,
		ErrorEstimate:  errorEstimate,
	}, nil
}

// SolveODEWithAdamsFixed 使用修复后的亚当斯法求解
func SolveODEWithAdamsFixed(params *ODEParamsFixed) (*ODEResultFixed, error) {
	// 计算步数
	nSteps := int(math.Round(params.TimeRange / params.TimeStep))
	if nSteps < 4 {
		// 步数太少时退化为欧拉法
		return SolveODEWithEulerFixed(params)
	}

	// 初始化数组
	timePoints := make([]float64, nSteps+1)
	solutionPath := make([]float64, nSteps+1)
	derivativePath := make([]float64, nSteps+1)

	// 初始条件
	timePoints[0] = 0.0
	solutionPath[0] = params.InitialValue
	derivativePath[0] = evaluateFirstOrderFixed(params.Equation, 0.0, params.InitialValue)

	// 前4步使用RK4启动
	for i := 0; i < 3 && i < nSteps; i++ {
		t := timePoints[i]
		y := solutionPath[i]
		h := params.TimeStep

		tNext := float64(i+1) * params.TimeStep
		if i == nSteps-1 {
			tNext = params.TimeRange
		}
		h = tNext - t

		k1 := h * evaluateFirstOrderFixed(params.Equation, t, y)
		k2 := h * evaluateFirstOrderFixed(params.Equation, t+h/2, y+k1/2)
		k3 := h * evaluateFirstOrderFixed(params.Equation, t+h/2, y+k2/2)
		k4 := h * evaluateFirstOrderFixed(params.Equation, t+h, y+k3)
		yNext := y + (k1+2*k2+2*k3+k4)/6

		timePoints[i+1] = tNext
		solutionPath[i+1] = yNext
		derivativePath[i+1] = evaluateFirstOrderFixed(params.Equation, tNext, yNext)
	}

	// Adams-Bashforth 4步法（标准系数）
	for i := 3; i < nSteps; i++ {
		t := timePoints[i]
		tNext := float64(i+1) * params.TimeStep
		if i == nSteps-1 {
			tNext = params.TimeRange
		}
		h := tNext - t

		// 标准 Adams-Bashforth 4阶公式
		// y_{n+1} = y_n + h/24 * (55*f_n - 59*f_{n-1} + 37*f_{n-2} - 9*f_{n-3})
		f0 := derivativePath[i]
		f1 := derivativePath[i-1]
		f2 := derivativePath[i-2]
		f3 := derivativePath[i-3]

		yNext := solutionPath[i] + h*(55*f0-59*f1+37*f2-9*f3)/24

		timePoints[i+1] = tNext
		solutionPath[i+1] = yNext
		derivativePath[i+1] = evaluateFirstOrderFixed(params.Equation, tNext, yNext)
	}

	stability := "conditionally_stable"
	if len(solutionPath) > 5 {
		maxVal := 0.0
		for _, val := range solutionPath {
			if math.Abs(val) > maxVal {
				maxVal = math.Abs(val)
			}
		}
		if maxVal > 1000 {
			stability = "unstable"
		} else if maxVal < 10 {
			stability = "stable"
		}
	}

	errorEstimate := estimateErrorFixed(params, solutionPath, timePoints)

	return &ODEResultFixed{
		Solution:       solutionPath[nSteps],
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "Adams",
		Stability:      stability,
		ErrorEstimate:  errorEstimate,
	}, nil
}

// evaluateFirstOrderFixed 计算一阶微分方程右端函数
func evaluateFirstOrderFixed(equation string, t, y float64) float64 {
	switch equation {
	case "dy/dt = -y":
		return -y // 指数衰减
	case "dy/dt = y":
		return y // 指数增长
	case "dy/dt = t":
		return t // 线性增长
	case "dy/dt = sin(t)":
		return math.Sin(t) // 正弦驱动
	case "dy/dt = -y + sin(t)":
		return -y + math.Sin(t) // 阻尼振荡
	default:
		return math.Sin(t) - 0.5*y
	}
}

// evaluateSecondOrderFixed 计算二阶微分方程右端函数
func evaluateSecondOrderFixed(equation string, t, y, yPrime float64) float64 {
	switch equation {
	case "d2y/dt2 = -y":
		return -y // 简谐振动
	case "d2y/dt2 = -0.1*y' - y":
		return -0.1*yPrime - y // 阻尼振动
	case "d2y/dt2 = sin(t)":
		return math.Sin(t) // 外力驱动
	default:
		return -0.1*yPrime - y + math.Sin(t)
	}
}

// assessStabilityFixed 评估数值稳定性
func assessStabilityFixed(params *ODEParamsFixed, solutionPath []float64) string {
	if len(solutionPath) < 2 {
		return "unknown"
	}

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
	}
	return "stable"
}

// estimateErrorFixed 估计数值误差（与精确解对比）
func estimateErrorFixed(params *ODEParamsFixed, solutionPath []float64, timePoints []float64) float64 {
	if len(solutionPath) < 2 || len(timePoints) != len(solutionPath) {
		return 0.0
	}

	// 获取精确解（如果可用）
	exactSolution := getExactSolution(params.Equation, params.InitialValue)
	if exactSolution == nil {
		// 无精确解时，使用相邻点变化率估计
		totalChange := 0.0
		for i := 1; i < len(solutionPath); i++ {
			change := math.Abs(solutionPath[i] - solutionPath[i-1])
			totalChange += change
		}
		return totalChange / float64(len(solutionPath)-1)
	}

	// 计算与精确解的均方根误差
	sumSquaredError := 0.0
	for i, t := range timePoints {
		exact := exactSolution(t)
		error := solutionPath[i] - exact
		sumSquaredError += error * error
	}

	rmse := math.Sqrt(sumSquaredError / float64(len(solutionPath)))
	return rmse
}

// exactSolutionFunc 精确解函数类型
type exactSolutionFunc func(t float64) float64

// getExactSolution 获取精确解函数（如果可用）
func getExactSolution(equation string, initialValue float64) exactSolutionFunc {
	switch equation {
	case "dy/dt = -y":
		// y(t) = y0 * exp(-t)
		return func(t float64) float64 {
			return initialValue * math.Exp(-t)
		}
	case "dy/dt = y":
		// y(t) = y0 * exp(t)
		return func(t float64) float64 {
			return initialValue * math.Exp(t)
		}
	case "dy/dt = t":
		// y(t) = y0 + t^2/2
		return func(t float64) float64 {
			return initialValue + t*t/2
		}
	default:
		return nil
	}
}

// SolveODE 统一求解入口
func SolveODE(params interface{}) (*ODEResultFixed, error) {
	odeParams, err := ParseODEParams(params)
	if err != nil {
		return nil, err
	}

	if err := ValidateODEParams(odeParams); err != nil {
		return nil, err
	}

	switch odeParams.Method {
	case "euler":
		return SolveODEWithEulerFixed(odeParams)
	case "rk4":
		return SolveODEWithRK4Fixed(odeParams)
	case "adams":
		return SolveODEWithAdamsFixed(odeParams)
	default:
		return nil, fmt.Errorf("不支持的求解方法: %s", odeParams.Method)
	}
}
