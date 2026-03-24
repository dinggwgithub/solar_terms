package bugs

import (
	"fmt"
	"math"
	"math/rand"
)

// PrecisionBug 精度错误Bug实现
type PrecisionBug struct{}

// NewPrecisionBug 创建新的精度错误Bug实例
func NewPrecisionBug() *PrecisionBug {
	return &PrecisionBug{}
}

// Name 返回Bug名称
func (b *PrecisionBug) Name() string {
	return "精度错误Bug"
}

// Description 返回Bug描述
func (b *PrecisionBug) Description() string {
	return "计算精度不足，包含浮点数截断、迭代提前终止、算法简化等问题"
}

// Apply 应用精度错误Bug
func (b *PrecisionBug) Apply(calculationType string, params interface{}) (interface{}, []string) {
	warnings := []string{
		"启用精度错误Bug模式",
		"精度损失：计算过程中精度被截断",
	}

	// 根据计算类型应用不同的精度错误策略
	switch calculationType {
	case "solar_term":
		return b.applyToSolarTerm(params), warnings
	case "ganzhi":
		return b.applyToGanZhi(params), warnings
	case "astronomy":
		return b.applyToAstronomy(params), warnings
	case "starting_age":
		return b.applyToStartingAge(params), warnings
	default:
		return nil, []string{"不支持的计算类型: " + calculationType}
	}
}

// applyToSolarTerm 应用到节气计算
func (b *PrecisionBug) applyToSolarTerm(params interface{}) interface{} {
	// 故意降低时间精度，只保留到分钟
	return "2024-02-04 12:00:00"
}

// applyToGanZhi 应用到干支计算
func (b *PrecisionBug) applyToGanZhi(params interface{}) interface{} {
	// 干支计算本身是离散的，精度错误主要体现在边界处理
	// 这里返回正常结果，因为干支计算本身不涉及浮点精度
	return map[string]string{
		"gan_year":  "甲",
		"zhi_year":  "辰",
		"gan_month": "丙",
		"zhi_month": "寅",
		"gan_day":   "甲",
		"zhi_day":   "子",
		"gan_time":  "甲",
		"zhi_time":  "子",
	}
}

// applyToAstronomy 应用到天文计算
func (b *PrecisionBug) applyToAstronomy(params interface{}) interface{} {
	// 故意降低天文计算的精度
	// 使用float32而不是float64，并截断小数位
	
	// 基础值
	baseLongitude := 280.123456789
	
	// 应用精度损失
	// 1. 转换为float32再转回float64
	precisionLost := float64(float32(baseLongitude))
	
	// 2. 进一步截断小数位
	truncated := math.Floor(precisionLost*100) / 100 // 保留2位小数
	
	// 3. 添加随机的小误差
	rand.Seed(42) // 固定种子确保可复现
	randomError := rand.Float64() * 0.01 // 0-0.01度的随机误差
	
	finalLongitude := truncated + randomError
	
	return map[string]float64{
		"sun_longitude":      finalLongitude,
		"julian_date":        2459580.5,
		"apparent_longitude": finalLongitude + 0.002, // 故意的小误差
		"true_longitude":     finalLongitude + 0.001, // 故意的小误差
		"mean_longitude":     280.12, // 截断精度
		"mean_anomaly":       357.53, // 截断精度
		"equation_of_center": 1.91,   // 截断精度
		"nutation":           0.004,  // 保持相对精度
	}
}

// applyToStartingAge 应用到起运岁数计算
func (b *PrecisionBug) applyToStartingAge(params interface{}) interface{} {
	// 起运岁数计算涉及小数处理，故意降低精度
	return "5岁"
}

// applyToLunar 应用到农历计算（新增功能）
func (b *PrecisionBug) applyToLunar(params interface{}) interface{} {
	// 农历计算涉及新月时刻的精确计算
	// 故意降低精度
	return map[string]interface{}{
		"lunar_year":   2024,
		"lunar_month":   1,
		"lunar_day":     1,
		"is_leap":       false,
		"lunar_string": "甲辰年正月初一",
	}
}

// applyToPlanet 应用到行星位置计算（新增功能）
func (b *PrecisionBug) applyToPlanet(params interface{}) interface{} {
	// 行星位置计算需要高精度
	// 故意降低精度
	return map[string]float64{
		"right_ascension": 12.34,  // 降低精度
		"declination":     23.45,  // 降低精度
		"distance":        5.20,   // 降低精度
		"magnitude":       -2.0,   // 降低精度
	}
}

// GetBugCharacteristics 返回Bug的特征描述
func (b *PrecisionBug) GetBugCharacteristics() map[string]string {
	return map[string]string{
		"类型":         "精度错误",
		"表现形式":     "计算结果精度不足，与权威数据存在较大误差",
		"根本原因":     "浮点数精度损失、迭代次数不足、算法简化、数值截断",
		"修复难度":     "中等",
		"影响范围":     "所有涉及浮点计算和数值近似的算法",
		"测试方法":     "与权威数据对比，检查误差是否在可接受范围内",
	}
}

// GetFixSuggestions 返回修复建议
func (b *PrecisionBug) GetFixSuggestions() []string {
	return []string{
		"使用更高精度的浮点类型（float64）",
		"增加迭代次数确保收敛精度",
		"实现高精度数学库",
		"使用有理数或定点数表示",
		"优化数值算法减少累积误差",
	}
}

// CalculatePrecisionError 计算精度误差（用于评估）
func (b *PrecisionBug) CalculatePrecisionError(actual, expected interface{}, calculationType string) (float64, error) {
	switch calculationType {
	case "astronomy":
		return b.calculateAstronomyPrecisionError(actual, expected)
	case "solar_term":
		return b.calculateSolarTermPrecisionError(actual, expected)
	default:
		return 0, fmt.Errorf("不支持的计算类型: %s", calculationType)
	}
}

// calculateAstronomyPrecisionError 计算天文计算的精度误差
func (b *PrecisionBug) calculateAstronomyPrecisionError(actual, expected interface{}) (float64, error) {
	actualMap, ok1 := actual.(map[string]float64)
	expectedMap, ok2 := expected.(map[string]float64)
	
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("参数类型错误")
	}
	
	// 计算太阳黄经的误差
	actualLongitude, ok1 := actualMap["sun_longitude"]
	expectedLongitude, ok2 := expectedMap["sun_longitude"]
	
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("缺少太阳黄经数据")
	}
	
	return math.Abs(actualLongitude - expectedLongitude), nil
}

// calculateSolarTermPrecisionError 计算节气计算的精度误差
func (b *PrecisionBug) calculateSolarTermPrecisionError(actual, expected interface{}) (float64, error) {
	actualStr, ok1 := actual.(string)
	expectedStr, ok2 := expected.(string)
	
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("参数类型错误")
	}
	
	// 解析时间字符串，计算时间差（秒）
	// 这里简化实现，实际应该解析时间字符串
	// 假设格式为 "YYYY-MM-DD HH:MM:SS"
	
	if len(actualStr) != 19 || len(expectedStr) != 19 {
		return 0, fmt.Errorf("时间格式错误")
	}
	
	// 简单比较字符串差异
	// 实际实现应该解析时间并计算差值
	diff := 0
	for i := 0; i < 19; i++ {
		if actualStr[i] != expectedStr[i] {
			diff++
		}
	}
	
	return float64(diff), nil
}

// GetPrecisionRequirements 获取精度要求
func (b *PrecisionBug) GetPrecisionRequirements(calculationType string) map[string]float64 {
	switch calculationType {
	case "astronomy":
		return map[string]float64{
			"sun_longitude":      0.000001, // 0.000001度
			"julian_date":        0.000001, // 0.000001天
			"apparent_longitude": 0.000001, // 0.000001度
		}
	case "solar_term":
		return map[string]float64{
			"time_precision": 1.0, // 1秒
		}
	case "planet":
		return map[string]float64{
			"right_ascension": 0.0001, // 0.0001小时
			"declination":     0.0001, // 0.0001度
		}
	default:
		return map[string]float64{
			"general_precision": 0.001, // 默认精度要求
		}
	}
}

// IsPrecisionAcceptable 检查精度是否可接受
func (b *PrecisionBug) IsPrecisionAcceptable(actual, expected interface{}, calculationType string) (bool, float64, error) {
	error, err := b.CalculatePrecisionError(actual, expected, calculationType)
	if err != nil {
		return false, 0, err
	}
	
	req := b.GetPrecisionRequirements(calculationType)
	
	// 获取相关精度要求
	var maxError float64
	switch calculationType {
	case "astronomy":
		maxError = req["sun_longitude"]
	case "solar_term":
		maxError = req["time_precision"]
	default:
		maxError = req["general_precision"]
	}
	
	return error <= maxError, error, nil
}