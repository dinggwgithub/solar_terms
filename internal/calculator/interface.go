package calculator

import (
	"fmt"
	"scientific_calc_bugs/internal/bugs"
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
	return result, warnings, nil
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
