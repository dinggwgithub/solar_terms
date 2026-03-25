package calculator

import (
	"fmt"
	"time"
)

// StartingAgeCalculator 起运岁数计算器
type StartingAgeCalculator struct {
	*BaseCalculator
}

// NewStartingAgeCalculator 创建新的起运岁数计算器
func NewStartingAgeCalculator() *StartingAgeCalculator {
	return &StartingAgeCalculator{
		BaseCalculator: NewBaseCalculator(
			"starting_age",
			"起运岁数计算，根据生辰八字计算起运年龄",
		),
	}
}

// StartingAgeParams 起运岁数计算参数
type StartingAgeParams struct {
	Year  int `json:"year"`  // 出生年份
	Month int `json:"month"` // 出生月份
	Day   int `json:"day"`   // 出生日期
	Hour  int `json:"hour"`  // 出生小时
}

// Calculate 执行起运岁数计算
func (c *StartingAgeCalculator) Calculate(params interface{}) (interface{}, error) {
	ageParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数有效性
	if err := c.validateParams(ageParams); err != nil {
		return nil, err
	}

	// 执行起运岁数计算
	result, err := c.calculateStartingAge(ageParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *StartingAgeCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// parseParams 解析参数
func (c *StartingAgeCalculator) parseParams(params interface{}) (*StartingAgeParams, error) {
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

	return &StartingAgeParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
		Hour:  int(hour),
	}, nil
}

// validateParams 验证参数有效性
func (c *StartingAgeCalculator) validateParams(params *StartingAgeParams) error {
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

// calculateStartingAge 计算起运岁数
func (c *StartingAgeCalculator) calculateStartingAge(params *StartingAgeParams) (map[string]interface{}, error) {
	// 计算出生日期
	birthDate := time.Date(params.Year, time.Month(params.Month), params.Day, params.Hour, 0, 0, 0, time.UTC)

	// 计算年干支
	ganYear, zhiYear := c.calculateYearGanZhi(params.Year)

	// 计算月干支
	ganMonth, zhiMonth := c.calculateMonthGanZhi(params.Year, params.Month)

	// 计算日干支
	ganDay, zhiDay := c.calculateDayGanZhi(params.Year, params.Month, params.Day)

	// 计算时干支
	ganTime, zhiTime := c.calculateTimeGanZhi(params.Hour, ganDay)

	// 计算性别（简化：根据年份尾数判断）
	gender := "男"
	if params.Year%2 == 0 {
		gender = "女"
	}

	// 计算起运年龄
	startingAge := c.calculateAge(birthDate, gender, ganYear)

	// 计算大运
	majorCycles := c.calculateMajorCycles(startingAge, gender, ganYear)

	return map[string]interface{}{
		"starting_age":     startingAge,
		"gender":           gender,
		"birth_bazi":       fmt.Sprintf("%s%s %s%s %s%s %s%s", ganYear, zhiYear, ganMonth, zhiMonth, ganDay, zhiDay, ganTime, zhiTime),
		"major_cycles":     majorCycles,
		"calculation_date": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// calculateYearGanZhi 计算年干支
func (c *StartingAgeCalculator) calculateYearGanZhi(year int) (string, string) {
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	ganIndex := (year - 4) % 10
	if ganIndex < 0 {
		ganIndex += 10
	}

	zhiIndex := (year - 4) % 12
	if zhiIndex < 0 {
		zhiIndex += 12
	}

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateMonthGanZhi 计算月干支
func (c *StartingAgeCalculator) calculateMonthGanZhi(year, month int) (string, string) {
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	ganIndex := (month - 1) % 10
	zhiIndex := (month - 1) % 12

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateDayGanZhi 计算日干支
func (c *StartingAgeCalculator) calculateDayGanZhi(year, month, day int) (string, string) {
	// 简化的日干支计算
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	baseDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	days := int(t.Sub(baseDate).Hours() / 24)

	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	ganIndex := days % 10
	zhiIndex := days % 12

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateTimeGanZhi 计算时干支
func (c *StartingAgeCalculator) calculateTimeGanZhi(hour int, dayGan string) (string, string) {
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	timeZhiMap := map[int]int{
		23: 0, 0: 0, 1: 1, 2: 1, 3: 2, 4: 2, 5: 3, 6: 3,
		7: 4, 8: 4, 9: 5, 10: 5, 11: 6, 12: 6, 13: 7, 14: 7,
		15: 8, 16: 8, 17: 9, 18: 9, 19: 10, 20: 10, 21: 11, 22: 11,
	}

	zhiIndex := timeZhiMap[hour]

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

	timeGanMap := map[int]int{0: 0, 1: 2, 2: 4, 3: 6, 4: 8}
	ganIndex := (timeGanMap[dayGanIndex%5] + zhiIndex) % 10

	return ganList[ganIndex], zhiList[zhiIndex]
}

// calculateAge 计算起运年龄
func (c *StartingAgeCalculator) calculateAge(birthDate time.Time, gender, yearGan string) int {
	// 简化的起运年龄计算
	// 实际应该基于精确的节气计算

	// 计算从出生到下一个节气的时间
	currentYear := birthDate.Year()
	nextTermDate := time.Date(currentYear, time.March, 20, 0, 0, 0, 0, time.UTC)
	if birthDate.After(nextTermDate) {
		nextTermDate = time.Date(currentYear+1, time.March, 20, 0, 0, 0, 0, time.UTC)
	}

	daysToNextTerm := nextTermDate.Sub(birthDate).Hours() / 24

	// 3天为1岁
	age := int(daysToNextTerm / 3)

	// 根据性别和年干调整
	if gender == "女" {
		// 女性逆排
		age = 10 - age
	}

	// 确保年龄在合理范围内
	if age < 1 {
		age = 1
	}
	if age > 10 {
		age = 10
	}

	return age
}

// calculateMajorCycles 计算大运
func (c *StartingAgeCalculator) calculateMajorCycles(startingAge int, gender, yearGan string) []map[string]interface{} {
	var cycles []map[string]interface{}

	// 天干列表
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

	// 地支列表
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 找到年干的位置
	yearGanIndex := -1
	for i, gan := range ganList {
		if gan == yearGan {
			yearGanIndex = i
			break
		}
	}

	if yearGanIndex == -1 {
		yearGanIndex = 0
	}

	// 计算大运
	for i := 0; i < 8; i++ {
		// 计算天干地支
		ganIndex := (yearGanIndex + i) % 10
		zhiIndex := (yearGanIndex + i) % 12

		if gender == "女" {
			// 女性逆排
			ganIndex = (yearGanIndex - i + 10) % 10
			zhiIndex = (yearGanIndex - i + 12) % 12
		}

		cycle := map[string]interface{}{
			"cycle":       i + 1,
			"age_start":   startingAge + i*10,
			"age_end":     startingAge + (i+1)*10 - 1,
			"ganzhi":      fmt.Sprintf("%s%s", ganList[ganIndex], zhiList[zhiIndex]),
			"description": c.getCycleDescription(ganList[ganIndex], zhiList[zhiIndex]),
		}

		cycles = append(cycles, cycle)
	}

	return cycles
}

// getCycleDescription 获取大运描述
func (c *StartingAgeCalculator) getCycleDescription(gan, zhi string) string {
	descriptions := map[string]string{
		"甲子": "事业起步，贵人相助",
		"乙丑": "稳步发展，积累经验",
		"丙寅": "创新突破，机遇增多",
		"丁卯": "人际关系，合作共赢",
		"戊辰": "稳定发展，财运提升",
		"己巳": "自我提升，学习成长",
		"庚午": "挑战增多，需要坚持",
		"辛未": "收获成果，享受生活",
		"壬申": "变动较大，适应变化",
		"癸酉": "反思总结，准备新阶段",
	}

	if desc, exists := descriptions[gan+zhi]; exists {
		return desc
	}

	return "平稳发展，顺其自然"
}

