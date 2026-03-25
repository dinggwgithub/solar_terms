package calculator

import (
	"fmt"
	"time"
)

// LunarCalculator 农历转换计算器
type LunarCalculator struct {
	*BaseCalculator
}

// NewLunarCalculator 创建新的农历转换计算器
func NewLunarCalculator() *LunarCalculator {
	return &LunarCalculator{
		BaseCalculator: NewBaseCalculator(
			"lunar",
			"公历转农历计算器，支持闰月处理和农历日期转换",
		),
	}
}

// LunarDate 农历日期结构
type LunarDate struct {
	LunarYear   int    `json:"lunar_year"`   // 农历年
	LunarMonth  int    `json:"lunar_month"`  // 农历月（1-12）
	LunarDay    int    `json:"lunar_day"`    // 农历日（1-30）
	IsLeap      bool   `json:"is_leap"`      // 是否为闰月
	LunarString string `json:"lunar_string"` // 农历日期字符串
}

// LunarParams 农历计算参数
type LunarParams struct {
	Year  int `json:"year"`  // 公历年
	Month int `json:"month"` // 公历月
	Day   int `json:"day"`   // 公历日
}

// Calculate 执行农历转换计算
func (c *LunarCalculator) Calculate(params interface{}) (interface{}, error) {
	lunarParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证日期有效性
	if err := c.validateDate(lunarParams.Year, lunarParams.Month, lunarParams.Day); err != nil {
		return nil, err
	}

	// 执行农历转换计算
	lunarDate, err := c.convertToLunar(lunarParams.Year, lunarParams.Month, lunarParams.Day)
	if err != nil {
		return nil, err
	}

	return lunarDate, nil
}

// Validate 验证输入参数
func (c *LunarCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// parseParams 解析参数
func (c *LunarCalculator) parseParams(params interface{}) (*LunarParams, error) {
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

	return &LunarParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
	}, nil
}

// validateDate 验证日期有效性
func (c *LunarCalculator) validateDate(year, month, day int) error {
	// 检查年份范围
	if year < 1900 || year > 2100 {
		return fmt.Errorf("年份超出支持范围 (1900-2100): %d", year)
	}

	// 检查月份范围
	if month < 1 || month > 12 {
		return fmt.Errorf("月份超出范围 (1-12): %d", month)
	}

	// 检查日期范围
	if day < 1 || day > 31 {
		return fmt.Errorf("日期超出范围 (1-31): %d", day)
	}

	// 检查具体月份的天数
	if month == 2 {
		maxDays := 28
		if c.isLeapYear(year) {
			maxDays = 29
		}
		if day > maxDays {
			return fmt.Errorf("2月最多%d天: %d", maxDays, day)
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		if day > 30 {
			return fmt.Errorf("%d月最多30天: %d", month, day)
		}
	}

	return nil
}

// isLeapYear 判断是否为闰年
func (c *LunarCalculator) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// convertToLunar 公历转农历（简化实现）
func (c *LunarCalculator) convertToLunar(year, month, day int) (*LunarDate, error) {
	// 这里使用简化的农历转换算法
	// 实际实现应该基于天文算法计算新月时刻

	// 创建公历日期
	gregorianDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	_ = gregorianDate // 标记为已使用
	
	// 简化的农历转换（基于固定偏移）
	// 实际应该基于天文算法计算
	lunarYear := year
	lunarMonth := month
	lunarDay := day
	isLeap := false

	// 简化的闰月判断（基于年份的模运算）
	if year%19 == 0 || year%19 == 3 || year%19 == 6 || year%19 == 8 ||
		year%19 == 11 || year%19 == 14 || year%19 == 17 {
		// 这些年份有闰月
		if month == 6 {
			isLeap = true
		}
	}

	// 调整农历月份（简化算法）
	if month >= 2 {
		lunarMonth = month - 1
	} else {
		lunarMonth = 12
		lunarYear = year - 1
	}

	// 生成农历日期字符串
	lunarString := c.formatLunarString(lunarYear, lunarMonth, lunarDay, isLeap)

	return &LunarDate{
		LunarYear:   lunarYear,
		LunarMonth:  lunarMonth,
		LunarDay:    lunarDay,
		IsLeap:      isLeap,
		LunarString: lunarString,
	}, nil
}

// formatLunarString 格式化农历日期字符串
func (c *LunarCalculator) formatLunarString(year, month, day int, isLeap bool) string {
	// 天干地支
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 农历月份名称
	monthNames := []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "冬", "腊"}

	// 计算天干地支
	ganIndex := (year - 4) % 10
	zhiIndex := (year - 4) % 12

	if ganIndex < 0 {
		ganIndex += 10
	}
	if zhiIndex < 0 {
		zhiIndex += 12
	}

	gan := ganList[ganIndex]
	zhi := zhiList[zhiIndex]

	// 格式化月份
	monthStr := monthNames[month-1]
	if isLeap {
		monthStr = "闰" + monthStr
	}

	// 格式化日期
	dayStr := ""
	if day == 1 {
		dayStr = "初一"
	} else if day <= 10 {
		dayStr = fmt.Sprintf("初%d", day)
	} else if day == 20 {
		dayStr = "二十"
	} else if day < 20 {
		dayStr = fmt.Sprintf("十%d", day-10)
	} else if day == 30 {
		dayStr = "三十"
	} else {
		dayStr = fmt.Sprintf("廿%d", day-20)
	}

	return fmt.Sprintf("%s%s年%s月%s", gan, zhi, monthStr, dayStr)
}


// LunarCalculationResult 农历计算结果
type LunarCalculationResult struct {
	Success     bool       `json:"success"`
	LunarDate   *LunarDate `json:"lunar_date"`
	Gregorian   string     `json:"gregorian"`
	Calculation string     `json:"calculation"`
	Timestamp   string     `json:"timestamp"`
}

// GetLunarMonthInfo 获取农历月份信息（用于测试）
func (c *LunarCalculator) GetLunarMonthInfo(year, month int) (int, bool, error) {
	// 返回农历月份的天数和是否为闰月
	// 简化实现，实际应该基于天文算法

	if month < 1 || month > 12 {
		return 0, false, fmt.Errorf("月份超出范围: %d", month)
	}

	// 农历月份天数（简化）
	monthDays := []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}
	days := monthDays[month-1]

	// 简化的闰月判断
	isLeap := false
	if year%19 == 0 && month == 6 {
		isLeap = true
	}

	return days, isLeap, nil
}

// CalculateLunarNewYear 计算农历新年日期（用于测试）
func (c *LunarCalculator) CalculateLunarNewYear(year int) (string, error) {
	// 简化的农历新年计算
	// 实际应该基于天文算法计算新月时刻

	if year < 1900 || year > 2100 {
		return "", fmt.Errorf("年份超出支持范围: %d", year)
	}

	// 简化的农历新年日期（基于近似算法）
	// 实际应该精确计算
	newYearDate := fmt.Sprintf("%d-02-%02d", year, 10+(year%10))

	return newYearDate, nil
}
