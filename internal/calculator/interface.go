package calculator

import (
	"fmt"
	"math"
	"math/rand"
	"scientific_calc_bugs/internal/bugs"
	"strings"
	"time"
)

// CalculationType 计算类型枚举
type CalculationType int

const (
	CalculationTypeSolarTerm CalculationType = iota
	CalculationTypeGanZhi
	CalculationTypeAstronomy
	CalculationTypeStartingAge
	CalculationTypeLunar
	CalculationTypePlanet
	CalculationTypeStar
	CalculationTypeSunriseSunset
	CalculationTypeMoonPhase
	CalculationTypeEquationSolver
	CalculationTypeSymbolicCalc
	CalculationTypeODESolver
)

// String 返回计算类型的字符串表示
func (ct CalculationType) String() string {
	switch ct {
	case CalculationTypeSolarTerm:
		return "solar_term"
	case CalculationTypeGanZhi:
		return "ganzhi"
	case CalculationTypeAstronomy:
		return "astronomy"
	case CalculationTypeStartingAge:
		return "starting_age"
	case CalculationTypeLunar:
		return "lunar"
	case CalculationTypePlanet:
		return "planet"
	case CalculationTypeStar:
		return "star"
	case CalculationTypeSunriseSunset:
		return "sunrise_sunset"
	case CalculationTypeMoonPhase:
		return "moon_phase"
	case CalculationTypeEquationSolver:
		return "equation_solver"
	case CalculationTypeSymbolicCalc:
		return "symbolic_calc"
	case CalculationTypeODESolver:
		return "ode_solver"
	default:
		return "unknown"
	}
}

// ParseCalculationType 从字符串解析计算类型
func ParseCalculationType(calcTypeStr string) (CalculationType, error) {
	switch calcTypeStr {
	case "solar_term":
		return CalculationTypeSolarTerm, nil
	case "ganzhi":
		return CalculationTypeGanZhi, nil
	case "astronomy":
		return CalculationTypeAstronomy, nil
	case "starting_age":
		return CalculationTypeStartingAge, nil
	case "lunar":
		return CalculationTypeLunar, nil
	case "planet":
		return CalculationTypePlanet, nil
	case "star":
		return CalculationTypeStar, nil
	case "sunrise_sunset":
		return CalculationTypeSunriseSunset, nil
	case "moon_phase":
		return CalculationTypeMoonPhase, nil
	case "equation_solver":
		return CalculationTypeEquationSolver, nil
	case "symbolic_calc":
		return CalculationTypeSymbolicCalc, nil
	case "ode_solver":
		return CalculationTypeODESolver, nil
	default:
		return CalculationTypeSolarTerm, fmt.Errorf("不支持的计算类型: %s", calcTypeStr)
	}
}

// Calculator 计算器接口定义
type Calculator interface {
	// Calculate 执行计算
	Calculate(params interface{}) (interface{}, error)
	// Validate 验证输入参数
	Validate(params interface{}) error
	// Description 返回计算器描述
	Description() string
	// GetSupportedBugTypes 返回支持的Bug类型
	GetSupportedBugTypes() []bugs.BugType
}

// CalculatorManager 计算器管理器
type CalculatorManager struct {
	calculators map[CalculationType]Calculator
	bugManager  *bugs.BugManager
}

// NewCalculatorManager 创建新的计算器管理器
func NewCalculatorManager(bugManager *bugs.BugManager) *CalculatorManager {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	return &CalculatorManager{
		calculators: make(map[CalculationType]Calculator),
		bugManager:  bugManager,
	}
}

// RegisterCalculator 注册计算器
func (m *CalculatorManager) RegisterCalculator(calcType CalculationType, calculator Calculator) {
	m.calculators[calcType] = calculator
}

// Calculate 执行计算（支持Bug应用）
func (m *CalculatorManager) Calculate(calcType CalculationType, params interface{}, bugType bugs.BugType) (interface{}, []string, error) {
	calculator, exists := m.calculators[calcType]
	if !exists {
		return nil, nil, fmt.Errorf("计算器未注册: %s", calcType.String())
	}

	// 验证输入参数
	if err := calculator.Validate(params); err != nil {
		return nil, nil, fmt.Errorf("参数验证失败: %v", err)
	}

	// 如果没有Bug，直接计算
	if bugType == bugs.BugTypeNone {
		result, err := calculator.Calculate(params)
		return result, []string{"正常计算模式"}, err
	}

	// 检查是否支持该Bug类型
	supportedBugs := calculator.GetSupportedBugTypes()
	bugSupported := false
	for _, supportedBug := range supportedBugs {
		if supportedBug == bugType {
			bugSupported = true
			break
		}
	}

	if !bugSupported {
		return nil, nil, fmt.Errorf("计算器不支持该Bug类型: %s", bugType.String())
	}

	// 应用Bug
	result, warnings := m.bugManager.ApplyBug(bugType, calcType.String(), params)

	// 【关键修复】如果Bug应用不支持该计算类型（返回nil）
	// 则先正常计算，然后对结果应用通用的Bug效果
	if result == nil && len(warnings) > 0 {
		// 检查警告信息是否包含"不支持的计算类型"
		if strings.Contains(warnings[0], "不支持的计算类型") {
			// 先执行正常计算
			normalResult, err := calculator.Calculate(params)
			if err != nil {
				return nil, nil, err
			}

			// 根据Bug类型应用通用效果
			result = m.applyGenericBug(bugType, normalResult)
			warnings = []string{
				fmt.Sprintf("%s模式（通用Bug效果）", bugType.String()),
				"该计算类型使用通用Bug注入",
			}
		}
	}

	return result, warnings, nil
}

// CalculateWithSession 执行计算（支持会话级Bug参数一致性）
func (m *CalculatorManager) CalculateWithSession(calcType CalculationType, params interface{}, bugType bugs.BugType, sessionID string, mixedMode bool) (interface{}, []string, error) {
	calculator, exists := m.calculators[calcType]
	if !exists {
		return nil, nil, fmt.Errorf("计算器未注册: %s", calcType.String())
	}

	// 验证输入参数
	if err := calculator.Validate(params); err != nil {
		return nil, nil, fmt.Errorf("参数验证失败: %v", err)
	}

	// 如果没有Bug，直接计算
	if bugType == bugs.BugTypeNone {
		result, err := calculator.Calculate(params)
		return result, []string{"正常计算模式"}, err
	}

	// 检查是否支持该Bug类型
	supportedBugs := calculator.GetSupportedBugTypes()
	bugSupported := false
	for _, supportedBug := range supportedBugs {
		if supportedBug == bugType {
			bugSupported = true
			break
		}
	}

	if !bugSupported {
		return nil, nil, fmt.Errorf("计算器不支持该Bug类型: %s", bugType.String())
	}

	// 首先尝试使用BugManager应用Bug
	result, warnings := m.bugManager.ApplyBug(bugType, calcType.String(), params)

	// 如果Bug应用不支持该计算类型（返回nil）
	// 则先正常计算，然后使用动态配置应用Bug效果
	if result == nil && len(warnings) > 0 {
		if strings.Contains(warnings[0], "不支持的计算类型") {
			// 先执行正常计算
			normalResult, err := calculator.Calculate(params)
			if err != nil {
				return nil, nil, err
			}

			// 获取或创建会话配置
			configManager := bugs.GetGlobalConfigManager()
			config := configManager.GetOrCreateSessionConfig(sessionID)

			// 设置混合模式 - 如果需要启用混合模式但配置中没有Bug类型，重新随机化
			config.EnableMixedMode = mixedMode
			if mixedMode && len(config.MixedBugTypes) == 0 {
				config.Randomize() // 重新随机化以生成混合Bug类型
			}

			// 使用动态配置应用Bug效果
			if mixedMode {
				// 混合Bug模式
				result = m.applyMixedBugsWithConfig(normalResult, config)
				warnings = []string{
					"混合Bug模式（动态参数）",
					fmt.Sprintf("会话ID: %s", sessionID),
					fmt.Sprintf("应用Bug类型: %v", config.MixedBugTypes),
				}
			} else {
				// 单一Bug模式
				result = m.applyGenericBugWithConfig(bugType, normalResult, config)
				warnings = []string{
					fmt.Sprintf("%s模式（动态参数）", bugType.String()),
					fmt.Sprintf("会话ID: %s", sessionID),
					"该计算类型使用动态Bug注入",
				}
			}
		}
	}

	return result, warnings, nil
}

// applyGenericBugWithConfig 使用动态配置应用Bug
func (m *CalculatorManager) applyGenericBugWithConfig(bugType bugs.BugType, result interface{}, config *bugs.BugDynamicConfig) interface{} {
	switch bugType {
	case bugs.BugTypeInstability:
		return m.addInstabilityToResultWithConfig(result, config)
	case bugs.BugTypeConstraint:
		return m.addConstraintViolationToResultWithConfig(result, config)
	case bugs.BugTypePrecision:
		return m.addPrecisionLossToResultWithConfig(result, config)
	default:
		return result
	}
}

// applyMixedBugsWithConfig 应用混合Bug
func (m *CalculatorManager) applyMixedBugsWithConfig(result interface{}, config *bugs.BugDynamicConfig) interface{} {
	intermediate := result
	for _, bt := range config.MixedApplyOrder {
		intermediate = m.applyGenericBugWithConfig(bt, intermediate, config)
	}
	return intermediate
}

// applyGenericBug 对计算结果应用通用的Bug效果（保持原有硬编码版本用于兼容）
func (m *CalculatorManager) applyGenericBug(bugType bugs.BugType, result interface{}) interface{} {
	switch bugType {
	case bugs.BugTypeInstability:
		// 不稳定性：对数值结果添加微小的随机变化
		return m.addInstabilityToResult(result)
	case bugs.BugTypeConstraint:
		// 约束越界：将结果调整到合理范围之外
		return m.addConstraintViolationToResult(result)
	case bugs.BugTypePrecision:
		// 精度损失：降低数值结果的精度
		return m.addPrecisionLossToResult(result)
	default:
		return result
	}
}

// addInstabilityToResult 对结果添加不稳定性
func (m *CalculatorManager) addInstabilityToResult(result interface{}) interface{} {
	// 递归或直接修改数值结果
	switch v := result.(type) {
	case float64:
		// 添加±5%的随机变化
		change := 1.0 + (rand.Float64()-0.5)*0.1
		return v * change
	case int:
		// 随机±1
		return v + rand.Intn(3) - 1
	case []float64:
		newSlice := make([]float64, len(v))
		for i, val := range v {
			change := 1.0 + (rand.Float64()-0.5)*0.1
			newSlice[i] = val * change
		}
		return newSlice
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, val := range v {
			newMap[key] = m.addInstabilityToResult(val)
		}
		return newMap
	case *EquationResult:
		// 专门处理方程求解器的结果
		return &EquationResult{
			Solution:      m.addInstabilityToResult(v.Solution),
			Iterations:    v.Iterations + rand.Intn(5) - 2, // ±2次迭代变化
			Converged:     v.Converged,
			Error:         v.Error * (1.0 + (rand.Float64()-0.5)*0.2), // ±10%误差变化
			FunctionValue: v.FunctionValue * (1.0 + (rand.Float64()-0.5)*0.1),
			Jacobian:      v.Jacobian,
			TimePoints:    m.addInstabilityToResult(v.TimePoints).([]float64),
			SolutionPath:  m.addInstabilityToResult(v.SolutionPath).([]float64),
		}
	default:
		// 对于其他类型，尝试通过反射处理
		if m, ok := result.(map[string]float64); ok {
			newMap := make(map[string]float64)
			for key, val := range m {
				change := 1.0 + (rand.Float64()-0.5)*0.1
				newMap[key] = val * change
			}
			return newMap
		}
		return result
	}
}

// addConstraintViolationToResult 对结果添加约束越界
func (m *CalculatorManager) addConstraintViolationToResult(result interface{}) interface{} {
	switch v := result.(type) {
	case float64:
		// 将数值放大或缩小到不合理范围
		if v > 0 {
			return v * 1000 // 放大1000倍
		}
		return v / 1000 // 缩小1000倍
	case int:
		if v > 0 {
			return v * 1000
		}
		return v / 1000
	case []float64:
		newSlice := make([]float64, len(v))
		for i, val := range v {
			if val > 0 {
				newSlice[i] = val * 1000
			} else {
				newSlice[i] = val / 1000
			}
		}
		return newSlice
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, val := range v {
			newMap[key] = m.addConstraintViolationToResult(val)
		}
		return newMap
	case *EquationResult:
		// 专门处理方程求解器的结果
		return &EquationResult{
			Solution:      m.addConstraintViolationToResult(v.Solution),
			Iterations:    v.Iterations * 1000, // 不合理的迭代次数
			Converged:     v.Converged,
			Error:         v.Error * 1000, // 不合理的误差
			FunctionValue: v.FunctionValue * 1000,
			Jacobian:      v.Jacobian,
			TimePoints:    m.addConstraintViolationToResult(v.TimePoints).([]float64),
			SolutionPath:  m.addConstraintViolationToResult(v.SolutionPath).([]float64),
		}
	default:
		if m, ok := result.(map[string]float64); ok {
			newMap := make(map[string]float64)
			for key, val := range m {
				if val > 0 {
					newMap[key] = val * 1000
				} else {
					newMap[key] = val / 1000
				}
			}
			return newMap
		}
		return result
	}
}

// addPrecisionLossToResult 对结果添加精度损失
func (m *CalculatorManager) addPrecisionLossToResult(result interface{}) interface{} {
	switch v := result.(type) {
	case float64:
		// 只保留两位小数
		return math.Round(v*100) / 100
	case []float64:
		newSlice := make([]float64, len(v))
		for i, val := range v {
			newSlice[i] = math.Round(val*100) / 100
		}
		return newSlice
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, val := range v {
			newMap[key] = m.addPrecisionLossToResult(val)
		}
		return newMap
	case *EquationResult:
		// 专门处理方程求解器的结果
		return &EquationResult{
			Solution:      m.addPrecisionLossToResult(v.Solution),
			Iterations:    v.Iterations,
			Converged:     v.Converged,
			Error:         math.Round(v.Error*100) / 100,
			FunctionValue: math.Round(v.FunctionValue*100) / 100,
			Jacobian:      v.Jacobian,
			TimePoints:    m.addPrecisionLossToResult(v.TimePoints).([]float64),
			SolutionPath:  m.addPrecisionLossToResult(v.SolutionPath).([]float64),
		}
	default:
		if m, ok := result.(map[string]float64); ok {
			newMap := make(map[string]float64)
			for key, val := range m {
				newMap[key] = math.Round(val*100) / 100
			}
			return newMap
		}
		return result
	}
}

// addConstraintViolationToResultWithConfig 使用动态配置添加约束越界Bug
func (m *CalculatorManager) addConstraintViolationToResultWithConfig(result interface{}, config *bugs.BugDynamicConfig) interface{} {
	switch v := result.(type) {
	case float64:
		return config.ApplyConstraint(v)
	case int:
		converted := float64(v)
		return int(config.ApplyConstraint(converted))
	case []float64:
		newSlice := make([]float64, len(v))
		for i, val := range v {
			newSlice[i] = config.ApplyConstraint(val)
		}
		return newSlice
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, val := range v {
			newMap[key] = m.addConstraintViolationToResultWithConfig(val, config)
		}
		return newMap
	case *EquationResult:
		// 专门处理方程求解器的结果 - Solution可能是float64或[]float64
		var newSolution interface{}
		if sol, ok := v.Solution.(float64); ok {
			newSolution = config.ApplyConstraint(sol)
		} else if sol, ok := v.Solution.([]float64); ok {
			newSolution = m.addConstraintViolationToResultWithConfig(sol, config).([]float64)
		} else {
			newSolution = v.Solution // 保持不变
		}
		return &EquationResult{
			Solution:      newSolution,
			Iterations:    int(config.ApplyConstraint(float64(v.Iterations))),
			Converged:     v.Converged,
			Error:         config.ApplyConstraint(v.Error),
			FunctionValue: config.ApplyConstraint(v.FunctionValue),
			Jacobian:      v.Jacobian,
			TimePoints:    m.addConstraintViolationToResultWithConfig(v.TimePoints, config).([]float64),
			SolutionPath:  m.addConstraintViolationToResultWithConfig(v.SolutionPath, config).([]float64),
		}
	case *PlanetPosition:
		// 专门处理行星位置计算结果
		return &PlanetPosition{
			RightAscension: config.ApplyConstraint(v.RightAscension),
			Declination:    config.ApplyConstraint(v.Declination),
			Distance:       config.ApplyConstraint(v.Distance),
			Magnitude:      config.ApplyConstraint(v.Magnitude),
			Phase:          config.ApplyConstraint(v.Phase),
			Elongation:     config.ApplyConstraint(v.Elongation),
		}
	case *StarResult:
		// 专门处理星曜计算结果
		return &StarResult{
			LunarDate:        v.LunarDate,      // 字符串字段保持不变
			DayGanZhi:        v.DayGanZhi,      // 字符串字段保持不变
			Constellation:    v.Constellation,  // 字符串字段保持不变
			StarPosition:     v.StarPosition,   // 字符串字段保持不变
			Auspicious:       v.Auspicious,     // 布尔字段保持不变
			AuspiciousInfo:   v.AuspiciousInfo, // 字符串数组保持不变
			DayScore:         config.ApplyConstraint(v.DayScore),
			ConstellationIdx: int(config.ApplyConstraint(float64(v.ConstellationIdx))),
			AuspiciousLevel:  config.ApplyConstraint(v.AuspiciousLevel),
			JulianDay:        config.ApplyConstraint(v.JulianDay),
			TimeCoordinate:   config.ApplyConstraint(v.TimeCoordinate),
		}
	case *SymbolicResult:
		// 专门处理符号计算结果
		return &SymbolicResult{
			OriginalExpression:   v.OriginalExpression, // 字符串保持不变
			ResultExpression:     v.ResultExpression,   // 字符串保持不变
			ParsedTree:           v.ParsedTree,         // 保持不变
			Derivative:           v.Derivative,         // 字符串保持不变
			Simplified:           v.Simplified,         // 字符串保持不变
			Variables:            v.Variables,          // map保持不变
			OperationType:        v.OperationType,      // 字符串保持不变
			NumericValue:         config.ApplyConstraint(v.NumericValue),
			ExpressionComplexity: config.ApplyConstraint(v.ExpressionComplexity),
			VariableCount:        int(config.ApplyConstraint(float64(v.VariableCount))),
			TermCount:            config.ApplyConstraint(v.TermCount),
			TreeDepth:            config.ApplyConstraint(v.TreeDepth),
			EvaluationScore:      config.ApplyConstraint(v.EvaluationScore),
		}
	case *MoonPhaseResult:
		// 专门处理月相计算结果
		return &MoonPhaseResult{
			Date:          v.Date,          // 字符串保持不变
			MoonPhase:     v.MoonPhase,     // 字符串保持不变
			NextPhase:     v.NextPhase,     // 字符串保持不变
			NextPhaseTime: v.NextPhaseTime, // 字符串保持不变
			PhaseAngle:    config.ApplyConstraint(v.PhaseAngle),
			Illumination:  config.ApplyConstraint(v.Illumination),
			Age:           config.ApplyConstraint(v.Age),
			Distance:      config.ApplyConstraint(v.Distance),
			Longitude:     config.ApplyConstraint(v.Longitude),
			Latitude:      config.ApplyConstraint(v.Latitude),
		}
	case *SunriseSunsetResult:
		// 专门处理日出日落计算结果
		return &SunriseSunsetResult{
			Date:      v.Date,      // 字符串保持不变
			Sunrise:   v.Sunrise,   // 字符串保持不变
			Sunset:    v.Sunset,    // 字符串保持不变
			SolarNoon: v.SolarNoon, // 字符串保持不变
			DayLength: config.ApplyConstraint(v.DayLength),
			CivilTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.CivilTwilight),
			NauticalTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.NauticalTwilight),
			AstronomicalTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.AstronomicalTwilight),
		}
	case *ODEResult:
		// 专门处理ODE求解器结果
		return &ODEResult{
			Solution:       config.ApplyConstraint(v.Solution),
			TimePoints:     m.addConstraintViolationToResultWithConfig(v.TimePoints, config).([]float64),
			SolutionPath:   m.addConstraintViolationToResultWithConfig(v.SolutionPath, config).([]float64),
			DerivativePath: m.addConstraintViolationToResultWithConfig(v.DerivativePath, config).([]float64),
			MethodUsed:     v.MethodUsed, // 字符串保持不变
			Stability:      v.Stability,  // 字符串保持不变
			ErrorEstimate:  config.ApplyConstraint(v.ErrorEstimate),
		}
	default:
		if m, ok := result.(map[string]float64); ok {
			newMap := make(map[string]float64)
			for key, val := range m {
				newMap[key] = config.ApplyConstraint(val)
			}
			return newMap
		}
		return result
	}
}

// addPrecisionLossToResultWithConfig 使用动态配置添加精度损失Bug
func (m *CalculatorManager) addPrecisionLossToResultWithConfig(result interface{}, config *bugs.BugDynamicConfig) interface{} {
	switch v := result.(type) {
	case float64:
		return config.ApplyPrecision(v)
	case []float64:
		newSlice := make([]float64, len(v))
		for i, val := range v {
			newSlice[i] = config.ApplyPrecision(val)
		}
		return newSlice
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, val := range v {
			newMap[key] = m.addPrecisionLossToResultWithConfig(val, config)
		}
		return newMap
	case *EquationResult:
		// 专门处理方程求解器的结果 - Solution可能是float64或[]float64
		var newSolution interface{}
		if sol, ok := v.Solution.(float64); ok {
			newSolution = config.ApplyPrecision(sol)
		} else if sol, ok := v.Solution.([]float64); ok {
			newSolution = m.addPrecisionLossToResultWithConfig(sol, config).([]float64)
		} else {
			newSolution = v.Solution // 保持不变
		}
		return &EquationResult{
			Solution:      newSolution,
			Iterations:    v.Iterations,
			Converged:     v.Converged,
			Error:         config.ApplyPrecision(v.Error),
			FunctionValue: config.ApplyPrecision(v.FunctionValue),
			Jacobian:      v.Jacobian,
			TimePoints:    m.addPrecisionLossToResultWithConfig(v.TimePoints, config).([]float64),
			SolutionPath:  m.addPrecisionLossToResultWithConfig(v.SolutionPath, config).([]float64),
		}
	case *PlanetPosition:
		// 专门处理行星位置计算结果
		return &PlanetPosition{
			RightAscension: config.ApplyPrecision(v.RightAscension),
			Declination:    config.ApplyPrecision(v.Declination),
			Distance:       config.ApplyPrecision(v.Distance),
			Magnitude:      config.ApplyPrecision(v.Magnitude),
			Phase:          config.ApplyPrecision(v.Phase),
			Elongation:     config.ApplyPrecision(v.Elongation),
		}
	case *StarResult:
		// 专门处理星曜计算结果
		return &StarResult{
			LunarDate:        v.LunarDate,      // 字符串字段保持不变
			DayGanZhi:        v.DayGanZhi,      // 字符串字段保持不变
			Constellation:    v.Constellation,  // 字符串字段保持不变
			StarPosition:     v.StarPosition,   // 字符串字段保持不变
			Auspicious:       v.Auspicious,     // 布尔字段保持不变
			AuspiciousInfo:   v.AuspiciousInfo, // 字符串数组保持不变
			DayScore:         config.ApplyPrecision(v.DayScore),
			ConstellationIdx: int(config.ApplyPrecision(float64(v.ConstellationIdx))),
			AuspiciousLevel:  config.ApplyPrecision(v.AuspiciousLevel),
			JulianDay:        config.ApplyPrecision(v.JulianDay),
			TimeCoordinate:   config.ApplyPrecision(v.TimeCoordinate),
		}
	case *SymbolicResult:
		// 专门处理符号计算结果
		return &SymbolicResult{
			OriginalExpression:   v.OriginalExpression, // 字符串保持不变
			ResultExpression:     v.ResultExpression,   // 字符串保持不变
			ParsedTree:           v.ParsedTree,         // 保持不变
			Derivative:           v.Derivative,         // 字符串保持不变
			Simplified:           v.Simplified,         // 字符串保持不变
			Variables:            v.Variables,          // map保持不变
			OperationType:        v.OperationType,      // 字符串保持不变
			NumericValue:         config.ApplyPrecision(v.NumericValue),
			ExpressionComplexity: config.ApplyPrecision(v.ExpressionComplexity),
			VariableCount:        int(config.ApplyPrecision(float64(v.VariableCount))),
			TermCount:            config.ApplyPrecision(v.TermCount),
			TreeDepth:            config.ApplyPrecision(v.TreeDepth),
			EvaluationScore:      config.ApplyPrecision(v.EvaluationScore),
		}
	case *MoonPhaseResult:
		// 专门处理月相计算结果
		return &MoonPhaseResult{
			Date:          v.Date,          // 字符串保持不变
			MoonPhase:     v.MoonPhase,     // 字符串保持不变
			NextPhase:     v.NextPhase,     // 字符串保持不变
			NextPhaseTime: v.NextPhaseTime, // 字符串保持不变
			PhaseAngle:    config.ApplyPrecision(v.PhaseAngle),
			Illumination:  config.ApplyPrecision(v.Illumination),
			Age:           config.ApplyPrecision(v.Age),
			Distance:      config.ApplyPrecision(v.Distance),
			Longitude:     config.ApplyPrecision(v.Longitude),
			Latitude:      config.ApplyPrecision(v.Latitude),
		}
	case *SunriseSunsetResult:
		// 专门处理日出日落计算结果
		return &SunriseSunsetResult{
			Date:      v.Date,      // 字符串保持不变
			Sunrise:   v.Sunrise,   // 字符串保持不变
			Sunset:    v.Sunset,    // 字符串保持不变
			SolarNoon: v.SolarNoon, // 字符串保持不变
			DayLength: config.ApplyPrecision(v.DayLength),
			CivilTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.CivilTwilight),
			NauticalTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.NauticalTwilight),
			AstronomicalTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.AstronomicalTwilight),
		}
	case *ODEResult:
		// 专门处理ODE求解器结果
		return &ODEResult{
			Solution:       config.ApplyPrecision(v.Solution),
			TimePoints:     m.addPrecisionLossToResultWithConfig(v.TimePoints, config).([]float64),
			SolutionPath:   m.addPrecisionLossToResultWithConfig(v.SolutionPath, config).([]float64),
			DerivativePath: m.addPrecisionLossToResultWithConfig(v.DerivativePath, config).([]float64),
			MethodUsed:     v.MethodUsed, // 字符串保持不变
			Stability:      v.Stability,  // 字符串保持不变
			ErrorEstimate:  config.ApplyPrecision(v.ErrorEstimate),
		}
	default:
		if m, ok := result.(map[string]float64); ok {
			newMap := make(map[string]float64)
			for key, val := range m {
				newMap[key] = config.ApplyPrecision(val)
			}
			return newMap
		}
		return result
	}
}

// addInstabilityToResultWithConfig 使用动态配置添加不稳定性Bug
func (m *CalculatorManager) addInstabilityToResultWithConfig(result interface{}, config *bugs.BugDynamicConfig) interface{} {
	switch v := result.(type) {
	case float64:
		return config.ApplyInstability(v)
	case int:
		// 对整数添加±10%的波动
		change := 1.0 + (rand.Float64()-0.5)*2*config.InstabilityVariation
		return int(float64(v) * change)
	case []float64:
		newSlice := make([]float64, len(v))
		for i, val := range v {
			newSlice[i] = config.ApplyInstability(val)
		}
		return newSlice
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, val := range v {
			newMap[key] = m.addInstabilityToResultWithConfig(val, config)
		}
		return newMap
	case *EquationResult:
		// 专门处理方程求解器的结果 - Solution可能是float64或[]float64
		var newSolution interface{}
		if sol, ok := v.Solution.(float64); ok {
			newSolution = config.ApplyInstability(sol)
		} else if sol, ok := v.Solution.([]float64); ok {
			newSolution = m.addInstabilityToResultWithConfig(sol, config).([]float64)
		} else {
			newSolution = v.Solution // 保持不变
		}
		// 使用确定性随机源来保证会话一致性
		localRand := config.GetSeededRand(float64(v.Iterations))
		return &EquationResult{
			Solution:      newSolution,
			Iterations:    v.Iterations + int((localRand.Float64()-0.5)*2*config.InstabilityVariation*10),
			Converged:     v.Converged,
			Error:         config.ApplyInstability(v.Error),
			FunctionValue: config.ApplyInstability(v.FunctionValue),
			Jacobian:      v.Jacobian,
			TimePoints:    m.addInstabilityToResultWithConfig(v.TimePoints, config).([]float64),
			SolutionPath:  m.addInstabilityToResultWithConfig(v.SolutionPath, config).([]float64),
		}
	case *PlanetPosition:
		// 专门处理行星位置计算结果
		return &PlanetPosition{
			RightAscension: config.ApplyInstability(v.RightAscension),
			Declination:    config.ApplyInstability(v.Declination),
			Distance:       config.ApplyInstability(v.Distance),
			Magnitude:      config.ApplyInstability(v.Magnitude),
			Phase:          config.ApplyInstability(v.Phase),
			Elongation:     config.ApplyInstability(v.Elongation),
		}
	case *StarResult:
		// 专门处理星曜计算结果
		return &StarResult{
			LunarDate:        v.LunarDate,      // 字符串字段保持不变
			DayGanZhi:        v.DayGanZhi,      // 字符串字段保持不变
			Constellation:    v.Constellation,  // 字符串字段保持不变
			StarPosition:     v.StarPosition,   // 字符串字段保持不变
			Auspicious:       v.Auspicious,     // 布尔字段保持不变
			AuspiciousInfo:   v.AuspiciousInfo, // 字符串数组保持不变
			DayScore:         config.ApplyInstability(v.DayScore),
			ConstellationIdx: int(config.ApplyInstability(float64(v.ConstellationIdx))),
			AuspiciousLevel:  config.ApplyInstability(v.AuspiciousLevel),
			JulianDay:        config.ApplyInstability(v.JulianDay),
			TimeCoordinate:   config.ApplyInstability(v.TimeCoordinate),
		}
	case *SymbolicResult:
		// 专门处理符号计算结果
		return &SymbolicResult{
			OriginalExpression:   v.OriginalExpression, // 字符串保持不变
			ResultExpression:     v.ResultExpression,   // 字符串保持不变
			ParsedTree:           v.ParsedTree,         // 保持不变
			Derivative:           v.Derivative,         // 字符串保持不变
			Simplified:           v.Simplified,         // 字符串保持不变
			Variables:            v.Variables,          // map保持不变
			OperationType:        v.OperationType,      // 字符串保持不变
			NumericValue:         config.ApplyInstability(v.NumericValue),
			ExpressionComplexity: config.ApplyInstability(v.ExpressionComplexity),
			VariableCount:        int(config.ApplyInstability(float64(v.VariableCount))),
			TermCount:            config.ApplyInstability(v.TermCount),
			TreeDepth:            config.ApplyInstability(v.TreeDepth),
			EvaluationScore:      config.ApplyInstability(v.EvaluationScore),
		}
	case *MoonPhaseResult:
		// 专门处理月相计算结果
		return &MoonPhaseResult{
			Date:          v.Date,          // 字符串保持不变
			MoonPhase:     v.MoonPhase,     // 字符串保持不变
			NextPhase:     v.NextPhase,     // 字符串保持不变
			NextPhaseTime: v.NextPhaseTime, // 字符串保持不变
			PhaseAngle:    config.ApplyInstability(v.PhaseAngle),
			Illumination:  config.ApplyInstability(v.Illumination),
			Age:           config.ApplyInstability(v.Age),
			Distance:      config.ApplyInstability(v.Distance),
			Longitude:     config.ApplyInstability(v.Longitude),
			Latitude:      config.ApplyInstability(v.Latitude),
		}
	case *SunriseSunsetResult:
		// 专门处理日出日落计算结果
		return &SunriseSunsetResult{
			Date:      v.Date,      // 字符串保持不变
			Sunrise:   v.Sunrise,   // 字符串保持不变
			Sunset:    v.Sunset,    // 字符串保持不变
			SolarNoon: v.SolarNoon, // 字符串保持不变
			DayLength: config.ApplyInstability(v.DayLength),
			CivilTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.CivilTwilight),
			NauticalTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.NauticalTwilight),
			AstronomicalTwilight: struct {
				Morning string `json:"morning"`
				Evening string `json:"evening"`
			}(v.AstronomicalTwilight),
		}
	case *ODEResult:
		// 专门处理ODE求解器结果
		return &ODEResult{
			Solution:       config.ApplyInstability(v.Solution),
			TimePoints:     m.addInstabilityToResultWithConfig(v.TimePoints, config).([]float64),
			SolutionPath:   m.addInstabilityToResultWithConfig(v.SolutionPath, config).([]float64),
			DerivativePath: m.addInstabilityToResultWithConfig(v.DerivativePath, config).([]float64),
			MethodUsed:     v.MethodUsed, // 字符串保持不变
			Stability:      v.Stability,  // 字符串保持不变
			ErrorEstimate:  config.ApplyInstability(v.ErrorEstimate),
		}
	default:
		if m, ok := result.(map[string]float64); ok {
			newMap := make(map[string]float64)
			for key, val := range m {
				newMap[key] = config.ApplyInstability(val)
			}
			return newMap
		}
		return result
	}
}

// GetCalculator 获取指定类型的计算器
func (m *CalculatorManager) GetCalculator(calcType CalculationType) (Calculator, bool) {
	calculator, exists := m.calculators[calcType]
	return calculator, exists
}

// GetAllCalculators 获取所有注册的计算器
func (m *CalculatorManager) GetAllCalculators() map[CalculationType]Calculator {
	return m.calculators
}

// GetCalculatorInfo 获取计算器信息
func (m *CalculatorManager) GetCalculatorInfo(calcType CalculationType) (map[string]string, error) {
	calculator, exists := m.calculators[calcType]
	if !exists {
		return nil, fmt.Errorf("计算器未注册: %s", calcType.String())
	}

	info := map[string]string{
		"name":        calcType.String(),
		"description": calculator.Description(),
	}

	// 添加支持的Bug类型
	supportedBugs := calculator.GetSupportedBugTypes()
	bugTypes := ""
	for i, bugType := range supportedBugs {
		if i > 0 {
			bugTypes += ", "
		}
		bugTypes += bugType.String()
	}
	info["supported_bugs"] = bugTypes

	return info, nil
}

// GetSupportedCalculationTypes 获取支持的计算类型
func (m *CalculatorManager) GetSupportedCalculationTypes() []CalculationType {
	var supportedTypes []CalculationType
	for calcType := range m.calculators {
		supportedTypes = append(supportedTypes, calcType)
	}
	return supportedTypes
}

// IsCalculationTypeSupported 检查计算类型是否支持
func (m *CalculatorManager) IsCalculationTypeSupported(calcType CalculationType) bool {
	_, exists := m.calculators[calcType]
	return exists
}

// BaseCalculator 基础计算器实现（可被具体计算器嵌入）
type BaseCalculator struct {
	name        string
	description string
}

// NewBaseCalculator 创建基础计算器
func NewBaseCalculator(name, description string) *BaseCalculator {
	return &BaseCalculator{
		name:        name,
		description: description,
	}
}

// Description 返回计算器描述
func (b *BaseCalculator) Description() string {
	return b.description
}

// GetSupportedBugTypes 返回支持的Bug类型（默认支持所有）
func (b *BaseCalculator) GetSupportedBugTypes() []bugs.BugType {
	return []bugs.BugType{
		bugs.BugTypeInstability,
		bugs.BugTypeConstraint,
		bugs.BugTypePrecision,
	}
}

// ValidateParams 通用参数验证函数
func ValidateParams(params interface{}, requiredFields []string) error {
	if params == nil {
		return fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("参数必须是map类型")
	}

	for _, field := range requiredFields {
		if _, exists := paramsMap[field]; !exists {
			return fmt.Errorf("缺少必需参数: %s", field)
		}
	}

	return nil
}

// CalculationResult 通用计算结果结构
type CalculationResult struct {
	Success   bool        `json:"success"`
	Result    interface{} `json:"result"`
	Warnings  []string    `json:"warnings"`
	Timestamp string      `json:"timestamp"`
}

// NewCalculationResult 创建新的计算结果
func NewCalculationResult(success bool, result interface{}, warnings []string) *CalculationResult {
	return &CalculationResult{
		Success:   success,
		Result:    result,
		Warnings:  warnings,
		Timestamp: "", // 将在外部设置
	}
}
