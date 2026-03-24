package bugs

import (
	"math/rand"
	"time"
)

// InstabilityBug 结果不稳定性Bug实现
type InstabilityBug struct{}

// NewInstabilityBug 创建新的结果不稳定性Bug实例
func NewInstabilityBug() *InstabilityBug {
	return &InstabilityBug{}
}

// Name 返回Bug名称
func (b *InstabilityBug) Name() string {
	return "结果不稳定性Bug"
}

// Description 返回Bug描述
func (b *InstabilityBug) Description() string {
	return "同一参数多次调用返回不同结果，包含随机种子未固定、迭代次数不足等问题"
}

// Apply 应用结果不稳定性Bug
func (b *InstabilityBug) Apply(calculationType string, params interface{}) (interface{}, []string) {
	warnings := []string{
		"启用结果不稳定Bug模式",
		"随机性干扰：优化过程受随机因素影响",
	}

	// 根据计算类型应用不同的不稳定性策略
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
func (b *InstabilityBug) applyToSolarTerm(params interface{}) interface{} {
	// 使用当前时间作为随机种子，确保每次调用不同
	rand.Seed(time.Now().UnixNano())

	// 解析参数获取年份
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return "参数格式错误"
	}

	year, ok := paramsMap["year"].(float64)
	if !ok {
		return "年份参数错误"
	}

	// 支持两种参数格式：数字索引或中文名称
	var termIndex int
	if termIndexFloat, ok := paramsMap["term_index"].(float64); ok {
		// 数字索引格式
		termIndex = int(termIndexFloat)
	} else if termName, ok := paramsMap["term_name"].(string); ok {
		// 中文名称格式
		termIndex = b.getTermIndexByName(termName)
		if termIndex == -1 {
			return "不支持的节气名称: " + termName
		}
	} else {
		return "必须提供term_index(数字索引)或term_name(中文名称)参数"
	}

	// 基于输入年份和节气索引计算基础时间
	// 立春节气（索引0）通常在2月4日左右，作为基准
	baseTime := time.Date(int(year), time.February, 4, 0, 0, 0, 0, time.UTC)

	// 每个节气间隔约15.2天
	termOffset := float64(int(termIndex)) * 15.2
	termTime := baseTime.Add(time.Duration(termOffset*24) * time.Hour)

	// 添加随机偏移：±1天内的随机时间
	randomDays := rand.Intn(2) - 1 // -1, 0, 1
	randomHours := rand.Intn(24)
	randomMinutes := rand.Intn(60)

	resultTime := termTime.Add(time.Duration(randomDays)*24*time.Hour +
		time.Duration(randomHours)*time.Hour +
		time.Duration(randomMinutes)*time.Minute)

	return resultTime.Format("2006-01-02 15:04:05")
}

// applyToGanZhi 应用到干支计算
func (b *InstabilityBug) applyToGanZhi(params interface{}) interface{} {
	// 使用当前时间作为随机种子
	rand.Seed(time.Now().UnixNano())

	// 随机修改天干地支
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 30%的概率随机修改结果
	if rand.Float64() < 0.3 {
		return map[string]string{
			"gan_year":  ganList[rand.Intn(len(ganList))],
			"zhi_year":  zhiList[rand.Intn(len(zhiList))],
			"gan_month": ganList[rand.Intn(len(ganList))],
			"zhi_month": zhiList[rand.Intn(len(zhiList))],
			"gan_day":   ganList[rand.Intn(len(ganList))],
			"zhi_day":   zhiList[rand.Intn(len(zhiList))],
			"gan_time":  ganList[rand.Intn(len(ganList))],
			"zhi_time":  zhiList[rand.Intn(len(zhiList))],
		}
	}

	// 返回正常的计算结果（这里应该调用正确的计算函数）
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
func (b *InstabilityBug) applyToAstronomy(params interface{}) interface{} {
	// 使用当前时间作为随机种子
	rand.Seed(time.Now().UnixNano())

	// 随机偏移黄经 (0-0.1度)
	randomOffset := rand.Float64() * 0.1

	// 返回带随机偏移的天文计算结果
	return map[string]float64{
		"sun_longitude":      280.123 + randomOffset,
		"julian_date":        2459580.5,
		"apparent_longitude": 280.123 + randomOffset,
		"true_longitude":     280.123 + randomOffset,
		"mean_longitude":     280.123,
		"mean_anomaly":       357.529,
		"equation_of_center": 1.914,
		"nutation":           0.004,
	}
}

// applyToStartingAge 应用到起运岁数计算
func (b *InstabilityBug) applyToStartingAge(params interface{}) interface{} {
	// 使用当前时间作为随机种子
	rand.Seed(time.Now().UnixNano())

	// 随机修改起运岁数 (±1岁)
	baseAge := 5
	randomOffset := rand.Intn(3) - 1 // -1, 0, 1

	if randomOffset == 0 {
		return "5岁"
	}

	_ = baseAge // 标记为已使用
	return "5岁"
}

// getTermIndexByName 根据节气名称获取索引
func (b *InstabilityBug) getTermIndexByName(name string) int {
	termNames := []string{
		"立春", "雨水", "惊蛰", "春分", "清明", "谷雨",
		"立夏", "小满", "芒种", "夏至", "小暑", "大暑",
		"立秋", "处暑", "白露", "秋分", "寒露", "霜降",
		"立冬", "小雪", "大雪", "冬至", "小寒", "大寒",
	}

	for i, termName := range termNames {
		if termName == name {
			return i
		}
	}
	return -1
}

// GetBugCharacteristics 返回Bug的特征描述
func (b *InstabilityBug) GetBugCharacteristics() map[string]string {
	return map[string]string{
		"类型":   "结果不稳定性",
		"表现形式": "同一参数多次调用返回不同结果",
		"根本原因": "随机种子未固定、迭代次数不足、全局状态依赖",
		"修复难度": "中等",
		"影响范围": "所有涉及随机性或迭代的计算",
		"测试方法": "多次调用相同参数，检查结果一致性",
	}
}


