package calculator

import (
	"fmt"
	"math/rand"
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
	CalculationTypeSymbolicCalcFixed
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
	case CalculationTypeSymbolicCalcFixed:
		return "symbolic_calc_fixed"
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
	case "symbolic_calc_fixed":
		return CalculationTypeSymbolicCalcFixed, nil
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
}

// CalculatorManager 计算器管理器
type CalculatorManager struct {
	calculators map[CalculationType]Calculator
}

// NewCalculatorManager 创建新的计算器管理器
func NewCalculatorManager() *CalculatorManager {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	return &CalculatorManager{
		calculators: make(map[CalculationType]Calculator),
	}
}

// RegisterCalculator 注册计算器
func (m *CalculatorManager) RegisterCalculator(calcType CalculationType, calculator Calculator) {
	m.calculators[calcType] = calculator
}

// Calculate 执行计算
func (m *CalculatorManager) Calculate(calcType CalculationType, params interface{}) (interface{}, []string, error) {
	calculator, exists := m.calculators[calcType]
	if !exists {
		return nil, nil, fmt.Errorf("计算器未注册: %s", calcType.String())
	}

	// 验证输入参数
	if err := calculator.Validate(params); err != nil {
		return nil, nil, fmt.Errorf("参数验证失败: %v", err)
	}

	// 执行计算
	result, err := calculator.Calculate(params)
	return result, nil, err
}

// CalculateWithSession 执行计算（支持会话级参数一致性）
func (m *CalculatorManager) CalculateWithSession(calcType CalculationType, params interface{}, sessionID string) (interface{}, []string, error) {
	calculator, exists := m.calculators[calcType]
	if !exists {
		return nil, nil, fmt.Errorf("计算器未注册: %s", calcType.String())
	}

	// 验证输入参数
	if err := calculator.Validate(params); err != nil {
		return nil, nil, fmt.Errorf("参数验证失败: %v", err)
	}

	// 执行计算
	result, err := calculator.Calculate(params)
	return result, nil, err
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
			return fmt.Errorf("缺少必填字段: %s", field)
		}
	}

	return nil
}

// GetFloatParam 从参数中获取float64值
func GetFloatParam(params interface{}, key string) (float64, error) {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("参数必须是map类型")
	}

	value, exists := paramsMap[key]
	if !exists {
		return 0, fmt.Errorf("缺少字段: %s", key)
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("字段 %s 必须是数值类型", key)
	}
}

// GetIntParam 从参数中获取int值
func GetIntParam(params interface{}, key string) (int, error) {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("参数必须是map类型")
	}

	value, exists := paramsMap[key]
	if !exists {
		return 0, fmt.Errorf("缺少字段: %s", key)
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("字段 %s 必须是整数类型", key)
	}
}

// GetStringParam 从参数中获取string值
func GetStringParam(params interface{}, key string) (string, error) {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("参数必须是map类型")
	}

	value, exists := paramsMap[key]
	if !exists {
		return "", fmt.Errorf("缺少字段: %s", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("字段 %s 必须是字符串类型", key)
	}

	return str, nil
}
