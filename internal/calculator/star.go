package calculator

import (
	"fmt"
	"math"
	"scientific_calc_bugs/internal/bugs"
	"time"
)

// StarCalculator 星曜推算计算器
type StarCalculator struct {
	*BaseCalculator
}

// NewStarCalculator 创建新的星曜推算计算器
func NewStarCalculator() *StarCalculator {
	return &StarCalculator{
		BaseCalculator: NewBaseCalculator(
			"star",
			"星曜推算计算器，计算二十八宿值日和星曜位置",
		),
	}
}

// StarParams 星曜计算参数
type StarParams struct {
	Year  int `json:"year"`  // 年
	Month int `json:"month"` // 月
	Day   int `json:"day"`   // 日
}

// StarResult 星曜计算结果
type StarResult struct {
	LunarDate        string   `json:"lunar_date"`                    // 农历日期
	DayGanZhi        string   `json:"day_ganzhi"`                    // 日干支
	Constellation    string   `json:"constellation"`                 // 二十八宿
	StarPosition     string   `json:"star_position"`                 // 星曜位置
	Auspicious       bool     `json:"auspicious"`                    // 是否吉日
	AuspiciousInfo   []string `json:"auspicious_info"`               // 吉凶信息
	DayScore         float64  `json:"day_score,omitempty"`           // 日分值（用于Bug演示）
	ConstellationIdx int      `json:"constellation_index,omitempty"` // 二十八宿索引（用于Bug演示）
	AuspiciousLevel  float64  `json:"auspicious_level,omitempty"`    // 吉凶程度量化值（用于Bug演示）
	JulianDay        float64  `json:"julian_day,omitempty"`          // 儒略日（用于Bug演示）
	TimeCoordinate   float64  `json:"time_coordinate,omitempty"`     // 时间坐标值（用于Bug演示）
}

// Calculate 执行星曜推算计算
func (c *StarCalculator) Calculate(params interface{}) (interface{}, error) {
	starParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证日期有效性
	if err := c.validateDate(starParams.Year, starParams.Month, starParams.Day); err != nil {
		return nil, err
	}

	// 执行星曜推算
	result, err := c.calculateStarInfo(starParams.Year, starParams.Month, starParams.Day)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *StarCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// parseParams 解析参数
func (c *StarCalculator) parseParams(params interface{}) (*StarParams, error) {
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

	return &StarParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
	}, nil
}

// validateDate 验证日期有效性
func (c *StarCalculator) validateDate(year, month, day int) error {
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
func (c *StarCalculator) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// calculateStarInfo 计算星曜信息
func (c *StarCalculator) calculateStarInfo(year, month, day int) (*StarResult, error) {
	// 计算农历日期（简化）
	lunarDate := c.calculateLunarDate(year, month, day)

	// 计算日干支
	dayGanZhi := c.calculateDayGanZhi(year, month, day)

	// 计算二十八宿
	constellation := c.calculateConstellation(year, month, day)

	// 计算星曜位置
	starPosition := c.calculateStarPosition(year, month, day)

	// 判断吉凶
	auspicious, auspiciousInfo := c.judgeAuspicious(dayGanZhi, constellation)

	// ========== 新增：计算数值字段（用于Bug演示） ==========

	// 计算儒略日
	jd := c.dateToJulianDay(year, month, day)

	// 日分值：基于年月日的哈希值（0-100分）
	dayScore := float64((year*10000 + month*100 + day) % 10000 / 100.0)

	// 二十八宿索引值（0-27）
	constellationIdx := c.calculateConstellationIndex(year, month, day)

	// 吉凶程度量化值（0-10分）
	auspiciousLevel := 5.0 // 基础分5分
	if auspicious {
		auspiciousLevel += 3.0
	}
	// 根据星宿吉凶调整
	if len(auspiciousInfo) > 0 && auspiciousInfo[len(auspiciousInfo)-1] == "星宿吉利" {
		auspiciousLevel += 1.5
	}
	// 根据干支吉凶调整
	if len(auspiciousInfo) > 0 && auspiciousInfo[0] == "日干支不吉" {
		auspiciousLevel -= 2.0
	}
	auspiciousLevel = math.Max(0, math.Min(10, auspiciousLevel)) // 限制在0-10之间

	// 时间坐标值（基于儒略日的归一化值）
	timeCoordinate := math.Mod(jd, 365.25)

	return &StarResult{
		LunarDate:        lunarDate,
		DayGanZhi:        dayGanZhi,
		Constellation:    constellation,
		StarPosition:     starPosition,
		Auspicious:       auspicious,
		AuspiciousInfo:   auspiciousInfo,
		DayScore:         dayScore,
		ConstellationIdx: constellationIdx,
		AuspiciousLevel:  auspiciousLevel,
		JulianDay:        jd,
		TimeCoordinate:   timeCoordinate,
	}, nil
}

// calculateLunarDate 计算农历日期（简化实现）
func (c *StarCalculator) calculateLunarDate(year, month, day int) string {
	// 简化的农历日期计算
	// 实际应该基于天文算法

	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

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

	return fmt.Sprintf("%s%s年%d月%d日", gan, zhi, month, day)
}

// calculateDayGanZhi 计算日干支
func (c *StarCalculator) calculateDayGanZhi(year, month, day int) string {
	// 简化的日干支计算
	// 实际应该基于精确的干支周期

	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 基于日期的简单哈希
	dayIndex := (year*10000 + month*100 + day) % 60

	ganIndex := dayIndex % 10
	zhiIndex := dayIndex % 12

	return fmt.Sprintf("%s%s", ganList[ganIndex], zhiList[zhiIndex])
}

// calculateConstellation 计算二十八宿
func (c *StarCalculator) calculateConstellation(year, month, day int) string {
	// 二十八宿列表
	constellations := []string{
		"角", "亢", "氐", "房", "心", "尾", "箕", // 东方青龙七宿
		"斗", "牛", "女", "虚", "危", "室", "壁", // 北方玄武七宿
		"奎", "娄", "胃", "昴", "毕", "觜", "参", // 西方白虎七宿
		"井", "鬼", "柳", "星", "张", "翼", "轸", // 南方朱雀七宿
	}

	// 基于日期的简单映射
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	yearDay := date.YearDay()

	constellationIndex := yearDay % len(constellations)

	return constellations[constellationIndex]
}

// calculateStarPosition 计算星曜位置
func (c *StarCalculator) calculateStarPosition(year, month, day int) string {
	// 主要星曜列表
	stars := []string{
		"太阳", "太阴", "岁星", "荧惑", "镇星", "太白", "辰星", // 七曜
		"紫气", "月孛", "罗睺", "计都", // 四余
		"天乙", "太乙", "青龙", "朱雀", "白虎", "玄武", // 六神
	}

	// 基于日期的简单映射
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	yearDay := date.YearDay()

	starIndex := yearDay % len(stars)

	// 简化的位置描述
	positions := []string{"在东方", "在南方", "在西方", "在北方", "在中天", "在地平"}
	positionIndex := (yearDay / 10) % len(positions)

	return fmt.Sprintf("%s%s", stars[starIndex], positions[positionIndex])
}

// judgeAuspicious 判断吉凶
func (c *StarCalculator) judgeAuspicious(dayGanZhi, constellation string) (bool, []string) {
	var auspiciousInfo []string

	// 吉日判断规则（简化）
	auspicious := true

	// 根据日干支判断
	if c.isAuspiciousGanZhi(dayGanZhi) {
		auspiciousInfo = append(auspiciousInfo, "日干支吉利")
	} else {
		auspicious = false
		auspiciousInfo = append(auspiciousInfo, "日干支不吉")
	}

	// 根据二十八宿判断
	if c.isAuspiciousConstellation(constellation) {
		auspiciousInfo = append(auspiciousInfo, "星宿吉利")
	} else {
		auspicious = false
		auspiciousInfo = append(auspiciousInfo, "星宿不吉")
	}

	// 特殊吉日判断
	if c.isSpecialAuspiciousDay(dayGanZhi, constellation) {
		auspicious = true
		auspiciousInfo = append(auspiciousInfo, "特殊吉日")
	}

	return auspicious, auspiciousInfo
}

// isAuspiciousGanZhi 判断日干支是否吉利
func (c *StarCalculator) isAuspiciousGanZhi(ganZhi string) bool {
	// 吉利的日干支（简化）
	auspiciousGanZhi := []string{
		"甲子", "乙丑", "丙寅", "丁卯", "戊辰", "己巳", "庚午", "辛未", "壬申", "癸酉",
		"甲戌", "乙亥", "丙子", "丁丑", "戊寅", "己卯", "庚辰", "辛巳", "壬午", "癸未",
	}

	for _, gz := range auspiciousGanZhi {
		if gz == ganZhi {
			return true
		}
	}

	return false
}

// isAuspiciousConstellation 判断二十八宿是否吉利
func (c *StarCalculator) isAuspiciousConstellation(constellation string) bool {
	// 吉利的二十八宿（简化）
	auspiciousConstellations := []string{
		"角", "房", "尾", "斗", "牛", "虚", "危", "室", "壁",
		"奎", "胃", "昴", "毕", "参", "井", "鬼", "柳", "星", "张", "翼", "轸",
	}

	for _, cons := range auspiciousConstellations {
		if cons == constellation {
			return true
		}
	}

	return false
}

// ========== 新增：Bug演示用的辅助函数 ==========

// dateToJulianDay 将日期转换为儒略日（天文计算用）
func (c *StarCalculator) dateToJulianDay(year, month, day int) float64 {
	Y := year
	M := month
	D := float64(day)

	if M <= 2 {
		Y--
		M += 12
	}

	A := float64(Y / 100)
	B := 2 - A + float64(int(A)/4)

	jd := float64(int(365.25*float64(Y))) + float64(int(30.6001*float64(M+1))) + D + B + 1720994.5

	return jd
}

// calculateConstellationIndex 计算二十八宿的索引值（0-27）
func (c *StarCalculator) calculateConstellationIndex(year, month, day int) int {
	// 二十八宿列表
	constellations := []string{
		"角", "亢", "氐", "房", "心", "尾", "箕", // 东方青龙七宿
		"斗", "牛", "女", "虚", "危", "室", "壁", // 北方玄武七宿
		"奎", "娄", "胃", "昴", "毕", "觜", "参", // 西方白虎七宿
		"井", "鬼", "柳", "星", "张", "翼", "轸", // 南方朱雀七宿
	}

	// 基于日期的简单映射
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	yearDay := date.YearDay()

	constellationIndex := yearDay % len(constellations)

	return constellationIndex
}

// isSpecialAuspiciousDay 判断特殊吉日
func (c *StarCalculator) isSpecialAuspiciousDay(dayGanZhi, constellation string) bool {
	// 特殊吉日组合（简化）
	specialDays := []string{
		"甲子角", "乙丑房", "丙寅尾", "丁卯斗", "戊辰牛",
		"己巳虚", "庚午危", "辛未室", "壬申壁", "癸酉奎",
	}

	combination := dayGanZhi + constellation
	for _, day := range specialDays {
		if day == combination {
			return true
		}
	}

	return false
}

// GetSupportedBugTypes 返回支持的Bug类型
func (c *StarCalculator) GetSupportedBugTypes() []bugs.BugType {
	return []bugs.BugType{
		bugs.BugTypeInstability,
		bugs.BugTypeConstraint,
		bugs.BugTypePrecision,
	}
}

// GetConstellationInfo 获取二十八宿信息（用于测试）
func (c *StarCalculator) GetConstellationInfo(constellation string) map[string]interface{} {
	info := map[string]interface{}{
		"name":       constellation,
		"group":      "",
		"element":    "",
		"auspicious": false,
	}

	// 分组信息
	switch constellation {
	case "角", "亢", "氐", "房", "心", "尾", "箕":
		info["group"] = "东方青龙"
		info["element"] = "木"
	case "斗", "牛", "女", "虚", "危", "室", "壁":
		info["group"] = "北方玄武"
		info["element"] = "水"
	case "奎", "娄", "胃", "昴", "毕", "觜", "参":
		info["group"] = "西方白虎"
		info["element"] = "金"
	case "井", "鬼", "柳", "星", "张", "翼", "轸":
		info["group"] = "南方朱雀"
		info["element"] = "火"
	}

	// 吉凶信息
	info["auspicious"] = c.isAuspiciousConstellation(constellation)

	return info
}

// GetStarCalendar 获取星历信息（用于测试）
func (c *StarCalculator) GetStarCalendar(year, month int) ([]map[string]interface{}, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("月份超出范围: %d", month)
	}

	// 计算该月的天数
	daysInMonth := 31
	if month == 2 {
		if c.isLeapYear(year) {
			daysInMonth = 29
		} else {
			daysInMonth = 28
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		daysInMonth = 30
	}

	var calendar []map[string]interface{}

	for day := 1; day <= daysInMonth; day++ {
		dayInfo := map[string]interface{}{
			"date":          fmt.Sprintf("%d-%02d-%02d", year, month, day),
			"ganzhi":        c.calculateDayGanZhi(year, month, day),
			"constellation": c.calculateConstellation(year, month, day),
			"star_position": c.calculateStarPosition(year, month, day),
		}

		auspicious, auspiciousInfo := c.judgeAuspicious(
			dayInfo["ganzhi"].(string),
			dayInfo["constellation"].(string),
		)

		dayInfo["auspicious"] = auspicious
		dayInfo["auspicious_info"] = auspiciousInfo

		calendar = append(calendar, dayInfo)
	}

	return calendar, nil
}
