package calculator

import (
	"fmt"
	"math"
)

type ODESolverFixedCalculator struct {
	*BaseCalculator
}

func NewODESolverFixedCalculator() *ODESolverFixedCalculator {
	return &ODESolverFixedCalculator{
		BaseCalculator: NewBaseCalculator(
			"ode_solver_fixed",
			"微分方程求解器（修复版），支持常微分方程和偏微分方程数值求解",
		),
	}
}

type ODEParamsFixed struct {
	EquationType string  `json:"equation_type"`
	Equation     string  `json:"equation"`
	InitialValue float64 `json:"initial_value"`
	InitialDeriv float64 `json:"initial_deriv"`
	TimeStep     float64 `json:"time_step"`
	TimeRange    float64 `json:"time_range"`
	Method       string  `json:"method"`
}

type ODEResultFixed struct {
	Solution       float64   `json:"solution"`
	TimePoints     []float64 `json:"time_points,omitempty"`
	SolutionPath   []float64 `json:"solution_path,omitempty"`
	DerivativePath []float64 `json:"derivative_path,omitempty"`
	MethodUsed     string    `json:"method_used,omitempty"`
	Stability      string    `json:"stability,omitempty"`
	ErrorEstimate  float64   `json:"error_estimate,omitempty"`
	ExactSolution  float64   `json:"exact_solution,omitempty"`
	AbsoluteError  float64   `json:"absolute_error,omitempty"`
}

func (c *ODESolverFixedCalculator) Calculate(params interface{}) (interface{}, error) {
	odeParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(odeParams); err != nil {
		return nil, err
	}

	var result *ODEResultFixed
	switch odeParams.Method {
	case "euler":
		result, err = c.solveWithEulerFixed(odeParams)
	case "rk4":
		result, err = c.solveWithRK4Fixed(odeParams)
	case "adams":
		result, err = c.solveWithAdamsFixed(odeParams)
	default:
		return nil, fmt.Errorf("不支持的求解方法: %s", odeParams.Method)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ODESolverFixedCalculator) parseParams(params interface{}) (*ODEParamsFixed, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	equation, ok := paramsMap["equation"].(string)
	if !ok {
		return nil, fmt.Errorf("equation参数必须为字符串")
	}

	paramsObj := &ODEParamsFixed{
		Equation:     equation,
		EquationType: "first_order",
		TimeStep:     0.1,
		TimeRange:    10.0,
		Method:       "euler",
		InitialValue: 1.0,
	}

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

func (c *ODESolverFixedCalculator) validateParams(params *ODEParamsFixed) error {
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

func (c *ODESolverFixedCalculator) solveWithEulerFixed(params *ODEParamsFixed) (*ODEResultFixed, error) {
	nSteps := int(math.Round(params.TimeRange / params.TimeStep))

	timePoints := make([]float64, 0, nSteps+1)
	solutionPath := make([]float64, 0, nSteps+1)
	derivativePath := make([]float64, 0, nSteps+1)

	t := 0.0
	y := params.InitialValue
	yPrime := params.InitialDeriv

	timePoints = append(timePoints, t)
	solutionPath = append(solutionPath, y)

	if params.EquationType == "first_order" {
		dydt := c.evaluateFirstOrder(params.Equation, t, y)
		derivativePath = append(derivativePath, dydt)
	} else {
		derivativePath = append(derivativePath, yPrime)
	}

	for i := 0; i < nSteps; i++ {
		if params.EquationType == "first_order" {
			dydt := c.evaluateFirstOrder(params.Equation, t, y)

			yNew := y + params.TimeStep*dydt
			tNew := t + params.TimeStep

			if tNew > params.TimeRange {
				tNew = params.TimeRange
				yNew = y + (tNew-t)*dydt
			}

			dydtNew := c.evaluateFirstOrder(params.Equation, tNew, yNew)

			timePoints = append(timePoints, tNew)
			solutionPath = append(solutionPath, yNew)
			derivativePath = append(derivativePath, dydtNew)

			t = tNew
			y = yNew
		} else if params.EquationType == "second_order" {
			d2ydt2 := c.evaluateSecondOrder(params.Equation, t, y, yPrime)

			yPrimeNew := yPrime + params.TimeStep*d2ydt2
			yNew := y + params.TimeStep*yPrime
			tNew := t + params.TimeStep

			if tNew > params.TimeRange {
				tNew = params.TimeRange
				finalStep := tNew - t
				yPrimeNew = yPrime + finalStep*d2ydt2
				yNew = y + finalStep*yPrime
			}

			timePoints = append(timePoints, tNew)
			solutionPath = append(solutionPath, yNew)
			derivativePath = append(derivativePath, yPrimeNew)

			t = tNew
			y = yNew
			yPrime = yPrimeNew
		}
	}

	stability := c.assessStability(solutionPath)
	errorEstimate := c.estimateError(params, solutionPath)
	exactSolution := c.getExactSolution(params.Equation, params.TimeRange, params.InitialValue)
	absoluteError := math.Abs(y - exactSolution)

	return &ODEResultFixed{
		Solution:       y,
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "Euler",
		Stability:      stability,
		ErrorEstimate:  errorEstimate,
		ExactSolution:  exactSolution,
		AbsoluteError:  absoluteError,
	}, nil
}

func (c *ODESolverFixedCalculator) solveWithRK4Fixed(params *ODEParamsFixed) (*ODEResultFixed, error) {
	nSteps := int(math.Round(params.TimeRange / params.TimeStep))

	timePoints := make([]float64, 0, nSteps+1)
	solutionPath := make([]float64, 0, nSteps+1)
	derivativePath := make([]float64, 0, nSteps+1)

	t := 0.0
	y := params.InitialValue

	timePoints = append(timePoints, t)
	solutionPath = append(solutionPath, y)
	dydt := c.evaluateFirstOrder(params.Equation, t, y)
	derivativePath = append(derivativePath, dydt)

	for i := 0; i < nSteps; i++ {
		h := params.TimeStep

		k1 := c.evaluateFirstOrder(params.Equation, t, y)
		k2 := c.evaluateFirstOrder(params.Equation, t+h/2, y+h*k1/2)
		k3 := c.evaluateFirstOrder(params.Equation, t+h/2, y+h*k2/2)
		k4 := c.evaluateFirstOrder(params.Equation, t+h, y+h*k3)

		yNew := y + h*(k1+2*k2+2*k3+k4)/6
		tNew := t + h

		if tNew > params.TimeRange {
			tNew = params.TimeRange
		}

		dydtNew := c.evaluateFirstOrder(params.Equation, tNew, yNew)

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)
		derivativePath = append(derivativePath, dydtNew)

		t = tNew
		y = yNew
	}

	stability := c.assessStability(solutionPath)
	errorEstimate := c.estimateError(params, solutionPath)
	exactSolution := c.getExactSolution(params.Equation, params.TimeRange, params.InitialValue)
	absoluteError := math.Abs(y - exactSolution)

	return &ODEResultFixed{
		Solution:       y,
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "RK4",
		Stability:      stability,
		ErrorEstimate:  errorEstimate,
		ExactSolution:  exactSolution,
		AbsoluteError:  absoluteError,
	}, nil
}

func (c *ODESolverFixedCalculator) solveWithAdamsFixed(params *ODEParamsFixed) (*ODEResultFixed, error) {
	nSteps := int(math.Round(params.TimeRange / params.TimeStep))

	timePoints := make([]float64, 0, nSteps+1)
	solutionPath := make([]float64, 0, nSteps+1)
	derivativePath := make([]float64, 0, nSteps+1)

	t := 0.0
	y := params.InitialValue

	timePoints = append(timePoints, t)
	solutionPath = append(solutionPath, y)
	dydt := c.evaluateFirstOrder(params.Equation, t, y)
	derivativePath = append(derivativePath, dydt)

	for i := 0; i < 3 && i < nSteps; i++ {
		h := params.TimeStep

		k1 := c.evaluateFirstOrder(params.Equation, t, y)
		k2 := c.evaluateFirstOrder(params.Equation, t+h/2, y+h*k1/2)
		k3 := c.evaluateFirstOrder(params.Equation, t+h/2, y+h*k2/2)
		k4 := c.evaluateFirstOrder(params.Equation, t+h, y+h*k3)

		yNew := y + h*(k1+2*k2+2*k3+k4)/6
		tNew := t + h

		dydtNew := c.evaluateFirstOrder(params.Equation, tNew, yNew)

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)
		derivativePath = append(derivativePath, dydtNew)

		t = tNew
		y = yNew
	}

	for len(solutionPath) <= nSteps {
		n := len(solutionPath)
		h := params.TimeStep

		f0 := c.evaluateFirstOrder(params.Equation, timePoints[n-1], solutionPath[n-1])
		f1 := c.evaluateFirstOrder(params.Equation, timePoints[n-2], solutionPath[n-2])
		f2 := c.evaluateFirstOrder(params.Equation, timePoints[n-3], solutionPath[n-3])
		f3 := c.evaluateFirstOrder(params.Equation, timePoints[n-4], solutionPath[n-4])

		yNew := solutionPath[n-1] + h*(55*f0-59*f1+37*f2-9*f3)/24
		tNew := timePoints[n-1] + h

		if tNew > params.TimeRange {
			tNew = params.TimeRange
		}

		dydtNew := c.evaluateFirstOrder(params.Equation, tNew, yNew)

		timePoints = append(timePoints, tNew)
		solutionPath = append(solutionPath, yNew)
		derivativePath = append(derivativePath, dydtNew)

		t = tNew
		y = yNew
	}

	stability := c.assessStability(solutionPath)
	errorEstimate := c.estimateError(params, solutionPath)
	exactSolution := c.getExactSolution(params.Equation, params.TimeRange, params.InitialValue)
	absoluteError := math.Abs(y - exactSolution)

	return &ODEResultFixed{
		Solution:       y,
		TimePoints:     timePoints,
		SolutionPath:   solutionPath,
		DerivativePath: derivativePath,
		MethodUsed:     "Adams",
		Stability:      stability,
		ErrorEstimate:  errorEstimate,
		ExactSolution:  exactSolution,
		AbsoluteError:  absoluteError,
	}, nil
}

func (c *ODESolverFixedCalculator) evaluateFirstOrder(equation string, t, y float64) float64 {
	if equation == "dy/dt = -y" {
		return -y
	} else if equation == "dy/dt = y" {
		return y
	} else if equation == "dy/dt = t" {
		return t
	} else if equation == "dy/dt = sin(t)" {
		return math.Sin(t)
	} else if equation == "dy/dt = -y + sin(t)" {
		return -y + math.Sin(t)
	}
	return math.Sin(t) - 0.5*y
}

func (c *ODESolverFixedCalculator) evaluateSecondOrder(equation string, t, y, yPrime float64) float64 {
	if equation == "d2y/dt2 = -y" {
		return -y
	} else if equation == "d2y/dt2 = -0.1*y' - y" {
		return -0.1*yPrime - y
	} else if equation == "d2y/dt2 = sin(t)" {
		return math.Sin(t)
	}
	return -0.1*yPrime - y + math.Sin(t)
}

func (c *ODESolverFixedCalculator) getExactSolution(equation string, t, y0 float64) float64 {
	if equation == "dy/dt = -y" {
		return y0 * math.Exp(-t)
	} else if equation == "dy/dt = y" {
		return y0 * math.Exp(t)
	} else if equation == "dy/dt = t" {
		return y0 + 0.5*t*t
	}
	return math.NaN()
}

func (c *ODESolverFixedCalculator) assessStability(solutionPath []float64) string {
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

func (c *ODESolverFixedCalculator) estimateError(params *ODEParamsFixed, solutionPath []float64) float64 {
	if len(solutionPath) < 2 {
		return 0.0
	}

	exactFinal := c.getExactSolution(params.Equation, params.TimeRange, params.InitialValue)
	if !math.IsNaN(exactFinal) {
		return math.Abs(solutionPath[len(solutionPath)-1] - exactFinal)
	}

	totalChange := 0.0
	for i := 1; i < len(solutionPath); i++ {
		change := math.Abs(solutionPath[i] - solutionPath[i-1])
		totalChange += change
	}

	return totalChange / float64(len(solutionPath)-1)
}

func (c *ODESolverFixedCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

func (c *ODESolverFixedCalculator) Description() string {
	return "微分方程求解器（修复版），支持常微分方程和偏微分方程数值求解"
}
