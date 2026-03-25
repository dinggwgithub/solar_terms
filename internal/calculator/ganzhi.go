package calculator

import (
	"fmt"
	"time"
)

// GanZhiCalculator 干支计算器
type GanZhiCalculator struct {
	*BaseCalculator
}

// NewGanZhiCalculator 创建新的干支计算器
func NewGanZhiCalculator() *GanZhiCalculator {
	return &GanZhiCalculator{
		BaseCalculator: NewBaseCalculator(
			"ganzhi",
			"干支计算器，计算年月日时的天干地支",
		),
	}
}

// GanZhiParams 干支计算参数
type GanZhiParams struct {
	Year  int `json:"year"`  // 年份
	Month int `json:"month"` // 月份
	Day   int `json:"day"`   // 日期
	Hour  int `json:"hour"`  // 小时
}

// Calculate 执行干支计算
func (c *GanZhiCalculator) Calculate(params interface{}) (interface{}, error) {
	ganzhiParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数有效性
	if err := c.validateParams(ganzhiParams); err != nil {
		return nil, err
	}

	// 执行干支计算
	result, err := c.calculateGanZhi(ganzhiParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *GanZhiCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// parseParams 解析参数
func (c *GanZhiCalculator) parseParams(params interface{}) (*GanZhiParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 提取参数
	year, ok := paramsMap["year"].(float64)
	if !ok {
		return nil, fmt.Errorf("year参数必须为数字")
	}

	month, ok := paramsMap["month"].(float64)
	if !ok {
		return nil, fmt.Errorf("month参数必须为数字")
	}

	day, ok := paramsMap["day"].(float64)
	if !ok {
		return nil, fmt.Errorf("day参数必须为数字")
	}

	hour, ok := paramsMap["hour"].(float64)
	if !ok {
		hour = 12 // 默认中午12点
	}

	return &GanZhiParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
		Hour:  int(hour),
	}, nil
}

// validateParams 验证参数有效性
func (c *GanZhiCalculator) validateParams(params *GanZhiParams) error {
	// 检查年份范围
	if params.Year < 1900 || params.Year > 2100 {
		return fmt.Errorf("年份超出支持范围 (1900-2100): %d", params.Year)
	}

	// 检查月份范围
	if params.Month < 1 || params.Month > 12 {
		return fmt.Errorf("月份超出范围 (1-12): %d", params.Month)
	}

	// 检查日期范围
	if params.Day < 1 || params.Day > 31 {
		return fmt.Errorf("日期超出范围 (1-31): %d", params.Day)
	}

	// 检查小时范围
	if params.Hour < 0 || params.Hour > 23 {
		return fmt.Errorf("小时超出范围 (0-23): %d", params.Hour)
	}

	return nil
}

// calculateGanZhi 计算干支
func (c *GanZhiCalculator) calculateGanZhi(params *GanZhiParams) (map[string]string, error) {
	// 计算年干支
	ganYear, zhiYear := c.calculateYearGanZhi(params.Year)

	// 计算月干支
	ganMonth, zhiMonth := c.calculateMonthGanZhi(params.Year, params.Month)

	// 计算日干支
	ganDay, zhiDay := c.calculateDayGanZhi(params.Year, params.Month, params.Day)

	// 计算时干支
	ganTime, zhiTime := c.calculateTimeGanZhi(params.Hour, ganDay)

	return map[string]string{
		"gan_year":  ganYear,
		"zhi_year":  zhiYear,
		"gan_month": ganMonth,
		"zhi_month": zhiMonth,
		"gan_day":   ganDay,
		"zhi_day":   zhiDay,
		"gan_time":  ganTime,
		"zhi_time":  zhiTime,
	}, nil
}

// calculateYearGanZhi 计算年干支
func (c *GanZhiCalculator) calculateYearGanZhi(year int) (string, string) {
	// 天干：年份尾数对应天干
	ganIndex := (year - 4) % 10
	if ganIndex < 0 {
		ganIndex += 10
	}

	// 地支：年份对应地支
	zhiIndex := (year - 4) % 12
	if zhiIndex < 0 {
		zhiIndex += 12
	}

	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateMonthGanZhi 计算月干支
func (c *GanZhiCalculator) calculateMonthGanZhi(year, month int) (string, string) {
	// 根据年干支计算月干支
	_, zhiYear := c.calculateYearGanZhi(year)
	_ = zhiYear // 标记为已使用

	// 地支对应月份
	zhiMonthMap := map[string]int{
		"寅": 1, "卯": 2, "辰": 3, "巳": 4, "午": 5, "未": 6,
		"申": 7, "酉": 8, "戌": 9, "亥": 10, "子": 11, "丑": 12,
	}
	_ = zhiMonthMap // 标记为已使用

	// 计算月干支
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 简化的月干支计算
	ganIndex := (month - 1) % 10
	zhiIndex := (month - 1) % 12

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateDayGanZhi 计算日干支
func (c *GanZhiCalculator) calculateDayGanZhi(year, month, day int) (string, string) {
	// 简化的日干支计算
	// 实际应该基于精确的儒略日计算

	// 创建日期
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// 计算从某个基准日开始的日数
	baseDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	days := int(t.Sub(baseDate).Hours() / 24)

	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	ganIndex := days % 10
	zhiIndex := days % 12

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateTimeGanZhi 计算时干支
func (c *GanZhiCalculator) calculateTimeGanZhi(hour int, dayGan string) (string, string) {
	// 根据日干计算时干支
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 时地支：小时对应地支
	timeZhiMap := map[int]int{
		23: 0, 0: 0, // 子时
		1: 1, 2: 1, // 丑时
		3: 2, 4: 2, // 寅时
		5: 3, 6: 3, // 卯时
		7: 4, 8: 4, // 辰时
		9: 5, 10: 5, // 巳时
		11: 6, 12: 6, // 午时
		13: 7, 14: 7, // 未时
		15: 8, 16: 8, // 申时
		17: 9, 18: 9, // 酉时
		19: 10, 20: 10, // 戌时
		21: 11, 22: 11, // 亥时
	}

	zhiIndex := timeZhiMap[hour]

	// 时天干：根据日干计算
	dayGanIndex := -1
	for i, gan := range ganList {
		if gan == dayGan {
			dayGanIndex = i
			break
		}
	}

	if dayGanIndex == -1 {
		return "甲", zhiList[zhiIndex]
	}

	// 五鼠遁口诀：甲己还加甲，乙庚丙作初，丙辛从戊起，丁壬庚子居，戊癸何方发，壬子是真途
	timeGanMap := map[int]int{
		0: 0, // 甲己日：甲子时
		1: 2, // 乙庚日：丙子时
		2: 4, // 丙辛日：戊子时
		3: 6, // 丁壬日：庚子时
		4: 8, // 戊癸日：壬子时
	}

	ganIndex := (timeGanMap[dayGanIndex%5] + zhiIndex) % 10

	return ganList[ganIndex], zhiList[zhiIndex]
}

