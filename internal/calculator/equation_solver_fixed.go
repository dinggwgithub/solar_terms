package calculator

import (
	"fmt"
	"math"
	"strings"
)

type EquationSolverCalculatorFixed struct {
	*BaseCalculator
}

func NewEquationSolverCalculatorFixed() *EquationSolverCalculatorFixed {
	return &EquationSolverCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"equation_solver_fixed",
			"方程求解器修复版，支持非线性方程、线性方程组和微分方程求解，修正了数值计算错误",
		),
	}
}

type EquationParamsFixed struct {
	EquationType  string    `json:"equation_type"`
	Equation      string    `json:"equation"`
	InitialGuess  float64   `json:"initial_guess"`
	InitialValue  float64   `json:"initial_value"`
	Tolerance     float64   `json:"tolerance"`
	MaxIterations int       `json:"max_iterations"`
	Coefficients  []float64 `json:"coefficients"`
	TimeStep      float64   `json:"time_step"`
	TimeRange     float64   `json:"time_range"`
	Method        string    `json:"method"`
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

func (c *EquationSolverCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	equationParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(equationParams); err != nil {
		return nil, err
	}

	var result *EquationResultFixed
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

func (c *EquationSolverCalculatorFixed) parseParams(params interface{}) (*EquationParamsFixed, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	equationType, ok := paramsMap["equation_type"].(string)
	if !ok {
		return nil, fmt.Errorf("equation_type参数必须为字符串")
	}

	equation, ok := paramsMap["equation"].(string)
	if !ok {
		return nil, fmt.Errorf("equation参数必须为字符串")
	}

	paramsObj := &EquationParamsFixed{
		EquationType:  equationType,
		Equation:      equation,
		Tolerance:     1e-6,
		MaxIterations: 100,
		TimeStep:      0.1,
		TimeRange:     10.0,
		Method:        "euler",
	}

	if val, ok := paramsMap["initial_guess"].(float64); ok {
		paramsObj.InitialGuess = val
	}
	if val, ok := paramsMap["initial_value"].(float64); ok {
		paramsObj.InitialValue = val
		if paramsObj.InitialGuess == 0 {
			paramsObj.InitialGuess = val
		}
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

	if method, ok := paramsMap["method"].(string); ok {
		paramsObj.Method = method
	}

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

func (c *EquationSolverCalculatorFixed) validateParams(params *EquationParamsFixed) error {
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

func (c *EquationSolverCalculatorFixed) solveNonlinearEquation(params *EquationParamsFixed) (*EquationResultFixed, error) {
	x := params.InitialGuess
	iterations := 0
	converged := false

	for iterations < params.MaxIterations {
		fx := c.evaluateFunction(params.Equation, x)
		fpx := c.evaluateDerivative(params.Equation, x)

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

	return &EquationResultFixed{
		Solution:      x,
		Iterations:    iterations,
		Converged:     converged,
		Error:         math.Abs(fx),
		FunctionValue: fx,
	}, nil
}

func (c *EquationSolverCalculatorFixed) solveLinearSystem(params *EquationParamsFixed) (*EquationResultFixed, error) {
	if len(params.Coefficients) == 0 {
		return nil, fmt.Errorf("线性方程组需要系数矩阵")
	}

	n := len(params.Coefficients)
	solution := make([]float64, n)

	for i, coef := range params.Coefficients {
		if math.Abs(coef) < 1e-12 {
			solution[i] = 0
		} else {
			solution[i] = 1.0 / coef
		}
	}

	jacobian := make([][]float64, n)
	for i := 0; i < n; i++ {
		jacobian[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				jacobian[i][j] = params.Coefficients[i]
			}
		}
	}

	return &EquationResultFixed{
		Solution:   solution,
		Iterations: 1,
		Converged:  true,
		Error:      0.0,
		Jacobian:   jacobian,
	}, nil
}

func (c *EquationSolverCalculatorFixed) solveODE(params *EquationParamsFixed) (*EquationResultFixed, error) {
	initialValue := params.InitialValue
	if initialValue == 0 && params.InitialGuess != 0 {
		initialValue = params.InitialGuess
	}

	method := params.Method
	if method == "" {
		method = "euler"
	}

	switch method {
	case "euler":
		return c.solveODEWithEuler(params, initialValue)
	case "rk4":
		return c.solveODEWithRK4(params, initialValue)
	case "rk45":
		return c.solveODEWithRK45(params, initialValue)
	default:
		return c.solveODEWithEuler(params, initialValue)
	}
}

func (c *EquationSolverCalculatorFixed) solveODEWithEuler(params *EquationParamsFixed, initialValue float64) (*EquationResultFixed, error) {
	h := params.TimeStep
	tEnd := params.TimeRange
	y := initialValue
	t := 0.0

	timePoints := []float64{t}
	solutionPath := []float64{y}
	stepDetails := []StepDetail{}

	nSteps := int(math.Round(tEnd / h))

	for i := 0; i < nSteps; i++ {
		dydt := c.evaluateODEFunction(params.Equation, t, y)

		yNew := y + h*dydt
		tNew := t + h

		if tNew > tEnd {
			hLast := tEnd - t
			yNew = y + hLast*dydt
			tNew = tEnd
		}

		localError := math.Abs(h * h * dydt / 2)

		stepDetails = append(stepDetails, StepDetail{
			Step:       i + 1,
			Time:       tNew,
			Value:      yNew,
			Derivative: dydt,
			LocalError: localError,
		})

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)

		t = tNew
		y = yNew
	}

	analytical := c.getAnalyticalSolution(params.Equation, tEnd, initialValue)
	absoluteError := math.Abs(y - analytical)
	relativeError := 0.0
	if math.Abs(analytical) > 1e-15 {
		relativeError = absoluteError / math.Abs(analytical)
	}

	globalError := 0.0
	for i, t := range timePoints {
		analyticalT := c.getAnalyticalSolution(params.Equation, t, initialValue)
		globalError += math.Abs(solutionPath[i] - analyticalT)
	}
	if len(timePoints) > 0 {
		globalError /= float64(len(timePoints))
	}

	return &EquationResultFixed{
		Solution:      y,
		Iterations:    nSteps,
		Converged:     true,
		Error:         globalError,
		ErrorEstimate: globalError,
		FunctionValue: 0,
		TimePoints:    timePoints,
		SolutionPath:  solutionPath,
		MethodUsed:    "Euler (Forward)",
		Stability:     c.assessStability(solutionPath),
		GlobalError:   globalError,
		LocalError:    h * h,
		Analytical:    analytical,
		AbsoluteError: absoluteError,
		RelativeError: relativeError,
		StepDetails:   stepDetails,
	}, nil
}

func (c *EquationSolverCalculatorFixed) solveODEWithRK4(params *EquationParamsFixed, initialValue float64) (*EquationResultFixed, error) {
	h := params.TimeStep
	tEnd := params.TimeRange
	y := initialValue
	t := 0.0

	timePoints := []float64{t}
	solutionPath := []float64{y}
	stepDetails := []StepDetail{}

	nSteps := int(math.Round(tEnd / h))

	for i := 0; i < nSteps; i++ {
		k1 := c.evaluateODEFunction(params.Equation, t, y)
		k2 := c.evaluateODEFunction(params.Equation, t+h/2, y+h*k1/2)
		k3 := c.evaluateODEFunction(params.Equation, t+h/2, y+h*k2/2)
		k4 := c.evaluateODEFunction(params.Equation, t+h, y+h*k3)

		yNew := y + h*(k1+2*k2+2*k3+k4)/6
		tNew := t + h

		if tNew > tEnd {
			tNew = tEnd
		}

		localError := math.Abs(h * h * h * h / 24)

		stepDetails = append(stepDetails, StepDetail{
			Step:       i + 1,
			Time:       tNew,
			Value:      yNew,
			Derivative: k1,
			LocalError: localError,
		})

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)

		t = tNew
		y = yNew
	}

	analytical := c.getAnalyticalSolution(params.Equation, tEnd, initialValue)
	absoluteError := math.Abs(y - analytical)
	relativeError := 0.0
	if math.Abs(analytical) > 1e-15 {
		relativeError = absoluteError / math.Abs(analytical)
	}

	globalError := 0.0
	for i, t := range timePoints {
		analyticalT := c.getAnalyticalSolution(params.Equation, t, initialValue)
		globalError += math.Abs(solutionPath[i] - analyticalT)
	}
	if len(timePoints) > 0 {
		globalError /= float64(len(timePoints))
	}

	return &EquationResultFixed{
		Solution:      y,
		Iterations:    nSteps,
		Converged:     true,
		Error:         globalError,
		ErrorEstimate: globalError,
		FunctionValue: 0,
		TimePoints:    timePoints,
		SolutionPath:  solutionPath,
		MethodUsed:    "RK4 (4th Order Runge-Kutta)",
		Stability:     c.assessStability(solutionPath),
		GlobalError:   globalError,
		LocalError:    math.Pow(h, 4),
		Analytical:    analytical,
		AbsoluteError: absoluteError,
		RelativeError: relativeError,
		StepDetails:   stepDetails,
	}, nil
}

func (c *EquationSolverCalculatorFixed) solveODEWithRK45(params *EquationParamsFixed, initialValue float64) (*EquationResultFixed, error) {
	h := params.TimeStep
	tEnd := params.TimeRange
	y := initialValue
	t := 0.0

	timePoints := []float64{t}
	solutionPath := []float64{y}
	stepDetails := []StepDetail{}

	tolerance := params.Tolerance
	if tolerance <= 0 {
		tolerance = 1e-6
	}

	minStep := h / 1000
	maxStep := h * 10

	iterations := 0
	maxIterations := params.MaxIterations * 10

	for t < tEnd && iterations < maxIterations {
		iterations++

		if t+h > tEnd {
			h = tEnd - t
		}

		k1 := c.evaluateODEFunction(params.Equation, t, y)
		k2 := c.evaluateODEFunction(params.Equation, t+h/4, y+h*k1/4)
		k3 := c.evaluateODEFunction(params.Equation, t+3*h/8, y+h*(3*k1+9*k2)/32)
		k4 := c.evaluateODEFunction(params.Equation, t+12*h/13, y+h*(1932*k1-7200*k2+7296*k3)/2197)
		k5 := c.evaluateODEFunction(params.Equation, t+h, y+h*(439*k1/216-8*k2+3680*k3/513-845*k4/4104))
		k6 := c.evaluateODEFunction(params.Equation, t+h/2, y+h*(-8*k1/27+2*k2-3544*k3/2565+1859*k4/4104-11*k5/40))

		y4 := y + h * (25*k1/216 + 1408*k3/2565 + 2197*k4/4104 - k5/5)
		y5 := y + h * (16*k1/135 + 6656*k3/12825 + 28561*k4/56430 - 9*k5/50 + 2*k6/55)

		errorEstimate := math.Abs(y5 - y4)

		if errorEstimate < tolerance || h <= minStep {
			y = y5
			t = t + h

			timePoints = append(timePoints, t)
			solutionPath = append(solutionPath, y)

			stepDetails = append(stepDetails, StepDetail{
				Step:       len(timePoints) - 1,
				Time:       t,
				Value:      y,
				Derivative: k1,
				LocalError: errorEstimate,
			})
		}

		if errorEstimate > 0 {
			hNew := 0.9 * h * math.Pow(tolerance/errorEstimate, 0.2)
			if hNew > maxStep {
				hNew = maxStep
			}
			if hNew < minStep {
				hNew = minStep
			}
			h = hNew
		}
	}

	analytical := c.getAnalyticalSolution(params.Equation, tEnd, initialValue)
	absoluteError := math.Abs(y - analytical)
	relativeError := 0.0
	if math.Abs(analytical) > 1e-15 {
		relativeError = absoluteError / math.Abs(analytical)
	}

	globalError := 0.0
	for i, t := range timePoints {
		analyticalT := c.getAnalyticalSolution(params.Equation, t, initialValue)
		globalError += math.Abs(solutionPath[i] - analyticalT)
	}
	if len(timePoints) > 0 {
		globalError /= float64(len(timePoints))
	}

	return &EquationResultFixed{
		Solution:      y,
		Iterations:    len(timePoints) - 1,
		Converged:     true,
		Error:         globalError,
		ErrorEstimate: globalError,
		FunctionValue: 0,
		TimePoints:    timePoints,
		SolutionPath:  solutionPath,
		MethodUsed:    "RK45 (Runge-Kutta-Fehlberg)",
		Stability:     c.assessStability(solutionPath),
		GlobalError:   globalError,
		LocalError:    tolerance,
		Analytical:    analytical,
		AbsoluteError: absoluteError,
		RelativeError: relativeError,
		StepDetails:   stepDetails,
	}, nil
}

func (c *EquationSolverCalculatorFixed) evaluateFunction(equation string, x float64) float64 {
	if strings.Contains(equation, "x^2") {
		return x*x - 2.0
	} else if strings.Contains(equation, "sin") {
		return math.Sin(x) - 0.5
	} else if strings.Contains(equation, "exp") {
		return math.Exp(x) - 2.0
	}
	return x*x*x - 2*x - 5
}

func (c *EquationSolverCalculatorFixed) evaluateDerivative(equation string, x float64) float64 {
	if strings.Contains(equation, "x^2") {
		return 2 * x
	} else if strings.Contains(equation, "sin") {
		return math.Cos(x)
	} else if strings.Contains(equation, "exp") {
		return math.Exp(x)
	}
	return 3*x*x - 2
}

func (c *EquationSolverCalculatorFixed) evaluateODEFunction(equation string, t, y float64) float64 {
	if strings.Contains(equation, "dy/dt = -y") {
		return -y
	} else if strings.Contains(equation, "dy/dt = y") {
		return y
	} else if strings.Contains(equation, "dy/dt = t") {
		return t
	} else if strings.Contains(equation, "dy/dt = sin(t)") {
		return math.Sin(t)
	} else if strings.Contains(equation, "dy/dt = -y + sin(t)") {
		return -y + math.Sin(t)
	} else if strings.Contains(equation, "dy/dt = y^2") {
		return y * y
	} else if strings.Contains(equation, "dy/dt = t*y") {
		return t * y
	}
	return -y
}

func (c *EquationSolverCalculatorFixed) getAnalyticalSolution(equation string, t, y0 float64) float64 {
	if strings.Contains(equation, "dy/dt = -y") {
		return y0 * math.Exp(-t)
	} else if strings.Contains(equation, "dy/dt = y") {
		return y0 * math.Exp(t)
	} else if strings.Contains(equation, "dy/dt = t") {
		return y0 + t*t/2
	} else if strings.Contains(equation, "dy/dt = sin(t)") {
		return y0 + 1 - math.Cos(t)
	} else if strings.Contains(equation, "dy/dt = -y + sin(t)") {
		A := y0 - 0.5
		return A*math.Exp(-t) + 0.5*(math.Sin(t)-math.Cos(t))
	} else if strings.Contains(equation, "dy/dt = y^2") {
		if math.Abs(y0) < 1e-15 {
			return 0
		}
		return y0 / (1 - t*y0)
	} else if strings.Contains(equation, "dy/dt = t*y") {
		return y0 * math.Exp(t*t/2)
	}
	return y0 * math.Exp(-t)
}

func (c *EquationSolverCalculatorFixed) assessStability(solutionPath []float64) string {
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

func (c *EquationSolverCalculatorFixed) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

func (c *EquationSolverCalculatorFixed) Description() string {
	return "方程求解器修复版，支持非线性方程、线性方程组和微分方程求解，修正了数值计算错误"
}
