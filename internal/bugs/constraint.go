package bugs

import (
	"fmt"
	"math/rand"
)

// ConstraintBug 约束越界Bug实现
type ConstraintBug struct{}

// NewConstraintBug 创建新的约束越界Bug实例
func NewConstraintBug() *ConstraintBug {
	return &ConstraintBug{}
}

// Name 返回Bug名称
func (b *ConstraintBug) Name() string {
	return "约束越界Bug"
}

// Description 返回Bug描述
func (b *ConstraintBug) Description() string {
	return "计算结果超出合理范围，包含边界检查缺失、返回无效值等问题"
}

// Apply 应用约束越界Bug
func (b *ConstraintBug) Apply(calculationType string, params interface{}) (interface{}, []string) {
	warnings := []string{
		"启用约束越界Bug模式",
		"约束越界：计算结果超出合理范围",
	}

	// 根据计算类型应用不同的约束越界策略
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
func (b *ConstraintBug) applyToSolarTerm(params interface{}) interface{} {
	// 故意返回超出范围的日期
	// 20%的概率返回过去或未来的极端日期
	if rand.Float64() < 0.2 {
		if rand.Float64() < 0.5 {
			return "1900-01-01 00:00:00" // 过去极端日期
		} else {
			return "2100-12-31 23:59:59" // 未来极端日期
		}
	}
	
	// 返回正常范围内的日期
	return "2024-02-04 12:00:00"
}

// applyToGanZhi 应用到干支计算
func (b *ConstraintBug) applyToGanZhi(params interface{}) interface{} {
	// 30%的概率返回无效的干支
	if rand.Float64() < 0.3 {
		return map[string]string{
			"gan_year":  "无",
			"zhi_year":  "无",
			"gan_month": "无",
			"zhi_month": "无",
			"gan_day":   "无",
			"zhi_day":   "无",
			"gan_time":  "无",
			"zhi_time":  "无",
		}
	}
	
	// 返回正常的计算结果
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
func (b *ConstraintBug) applyToAstronomy(params interface{}) interface{} {
	// 返回超出正常范围的天文数值
	return map[string]float64{
		"sun_longitude":      720.0, // 超出360度范围
		"julian_date":        -1000.0, // 负的儒略日
		"apparent_longitude": 500.0, // 超出正常范围
		"true_longitude":     400.0, // 超出正常范围
		"mean_longitude":     600.0, // 超出正常范围
		"mean_anomaly":       800.0, // 超出正常范围
		"equation_of_center": 50.0, // 超出正常范围
		"nutation":           10.0, // 超出正常范围
	}
}

// applyToStartingAge 应用到起运岁数计算
func (b *ConstraintBug) applyToStartingAge(params interface{}) interface{} {
	// 故意返回负数的年龄
	// 检查输入参数，如果是小数值则返回负数
	if paramsMap, ok := params.(map[string]interface{}); ok {
		if daysFromTerm, ok := paramsMap["days_from_term"].(float64); ok {
			if daysFromTerm < 10 {
				return "-1岁" // 返回负数年龄
			}
		}
	}
	
	// 返回正常的计算结果
	return "5岁"
}

// applyToLunar 应用到农历计算（新增功能）
func (b *ConstraintBug) applyToLunar(params interface{}) interface{} {
	// 返回无效的农历日期
	return map[string]interface{}{
		"lunar_year":   -1, // 负的农历年份
		"lunar_month":  13, // 超出1-12月范围
		"lunar_day":    32, // 超出1-30日范围
		"is_leap":      true, // 错误的闰月标记
		"lunar_string": "无效农历日期",
	}
}

// applyToPlanet 应用到行星位置计算（新增功能）
func (b *ConstraintBug) applyToPlanet(params interface{}) interface{} {
	// 返回超出正常范围的行星位置
	return map[string]float64{
		"right_ascension":  25.0, // 超出0-24小时范围
		"declination":      100.0, // 超出-90到90度范围
		"distance":         -1.0, // 负的距离
		"magnitude":        20.0, // 异常的星等值
	}
}

// GetBugCharacteristics 返回Bug的特征描述
func (b *ConstraintBug) GetBugCharacteristics() map[string]string {
	return map[string]string{
		"类型":         "约束越界",
		"表现形式":     "计算结果超出合理范围，返回无效值",
		"根本原因":     "边界条件检查缺失、输入验证不足、算法实现错误",
		"修复难度":     "简单到中等",
		"影响范围":     "所有涉及数值范围和边界条件的计算",
		"测试方法":     "输入边界值和异常值，检查结果是否在合理范围内",
	}
}

// GetFixSuggestions 返回修复建议
func (b *ConstraintBug) GetFixSuggestions() []string {
	return []string{
		"添加输入参数验证",
		"实现边界条件检查",
		"使用断言确保计算结果有效性",
		"添加默认值处理",
		"实现错误处理机制",
	}
}

// ValidateConstraints 验证约束条件（修复后的示例）
func (b *ConstraintBug) ValidateConstraints(result interface{}, calculationType string) error {
	switch calculationType {
	case "solar_term":
		return b.validateSolarTermConstraints(result)
	case "ganzhi":
		return b.validateGanZhiConstraints(result)
	case "astronomy":
		return b.validateAstronomyConstraints(result)
	case "starting_age":
		return b.validateStartingAgeConstraints(result)
	default:
		return fmt.Errorf("不支持的计算类型: %s", calculationType)
	}
}

// validateSolarTermConstraints 验证节气计算约束
func (b *ConstraintBug) validateSolarTermConstraints(result interface{}) error {
	if resultStr, ok := result.(string); ok {
		// 检查日期格式和范围
		// 这里应该实现具体的验证逻辑
		if len(resultStr) != 19 { // "YYYY-MM-DD HH:MM:SS"
			return fmt.Errorf("日期格式错误: %s", resultStr)
		}
		return nil
	}
	return fmt.Errorf("结果类型错误: %T", result)
}

// validateGanZhiConstraints 验证干支计算约束
func (b *ConstraintBug) validateGanZhiConstraints(result interface{}) error {
	if resultMap, ok := result.(map[string]string); ok {
		validGan := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
		validZhi := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
		
		for key, value := range resultMap {
			if key == "gan_year" || key == "gan_month" || key == "gan_day" || key == "gan_time" {
				if !contains(validGan, value) {
					return fmt.Errorf("无效的天干: %s", value)
				}
			} else if key == "zhi_year" || key == "zhi_month" || key == "zhi_day" || key == "zhi_time" {
				if !contains(validZhi, value) {
					return fmt.Errorf("无效的地支: %s", value)
				}
			}
		}
		return nil
	}
	return fmt.Errorf("结果类型错误: %T", result)
}

// validateAstronomyConstraints 验证天文计算约束
func (b *ConstraintBug) validateAstronomyConstraints(result interface{}) error {
	if resultMap, ok := result.(map[string]float64); ok {
		// 检查黄经范围 (0-360度)
		if longitude, exists := resultMap["sun_longitude"]; exists {
			if longitude < 0 || longitude >= 360 {
				return fmt.Errorf("黄经超出范围: %f", longitude)
			}
		}
		
		// 检查儒略日范围（合理的日期范围）
		if jd, exists := resultMap["julian_date"]; exists {
			if jd < 0 {
				return fmt.Errorf("儒略日不能为负: %f", jd)
			}
		}
		
		return nil
	}
	return fmt.Errorf("结果类型错误: %T", result)
}

// validateStartingAgeConstraints 验证起运岁数约束
func (b *ConstraintBug) validateStartingAgeConstraints(result interface{}) error {
	if resultStr, ok := result.(string); ok {
		// 检查年龄不能为负数
		if len(resultStr) > 0 && resultStr[0] == '-' {
			return fmt.Errorf("起运岁数不能为负: %s", resultStr)
		}
		return nil
	}
	return fmt.Errorf("结果类型错误: %T", result)
}

// contains 检查切片是否包含元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}