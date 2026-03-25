package calculator

import (
	"fmt"
	"math"
)

// MoonPhaseCalculator 月相计算器
type MoonPhaseCalculator struct {
	*BaseCalculator
}

// NewMoonPhaseCalculator 创建新的月相计算器
func NewMoonPhaseCalculator() *MoonPhaseCalculator {
	return &MoonPhaseCalculator{
		BaseCalculator: NewBaseCalculator(
			"moon_phase",
			"月相计算器，计算月相类型和精确时间",
		),
	}
}

// MoonPhaseType 月相类型枚举
type MoonPhaseType int

const (
	MoonPhaseNew          MoonPhaseType = iota // 新月
	MoonPhaseFirstQuarter                      // 上弦月
	MoonPhaseFull                              // 满月
	MoonPhaseLastQuarter                       // 下弦月
)

// String 返回月相类型的字符串表示
func (mpt MoonPhaseType) String() string {
	switch mpt {
	case MoonPhaseNew:
		return "new"
	case MoonPhaseFirstQuarter:
		return "first_quarter"
	case MoonPhaseFull:
		return "full"
	case MoonPhaseLastQuarter:
		return "last_quarter"
	default:
		return "unknown"
	}
}

// MoonPhaseParams 月相计算参数
type MoonPhaseParams struct {
	Year  int `json:"year"`  // 年
	Month int `json:"month"` // 月
	Day   int `json:"day"`   // 日（可选，用于指定日期）
}

// MoonPhaseResult 月相计算结果
type MoonPhaseResult struct {
	Date          string  `json:"date"`            // 日期
	MoonPhase     string  `json:"moon_phase"`      // 月相类型
	PhaseAngle    float64 `json:"phase_angle"`     // 相位角（度）
	Illumination  float64 `json:"illumination"`    // 照明比例（0-1）
	Age           float64 `json:"age"`             // 月龄（天）
	NextPhase     string  `json:"next_phase"`      // 下一个月相
	NextPhaseTime string  `json:"next_phase_time"` // 下一个月相时间
	Distance      float64 `json:"distance"`        // 地月距离（公里）
	Longitude     float64 `json:"longitude"`       // 月球黄经（度）
	Latitude      float64 `json:"latitude"`        // 月球黄纬（度）
}

// Calculate 执行月相计算
func (c *MoonPhaseCalculator) Calculate(params interface{}) (interface{}, error) {
	moonParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证日期有效性
	if err := c.validateDate(moonParams.Year, moonParams.Month, moonParams.Day); err != nil {
		return nil, err
	}

	// 执行月相计算
	result, err := c.calculateMoonPhase(moonParams.Year, moonParams.Month, moonParams.Day)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *MoonPhaseCalculator) Validate(params interface{}) error {
	moonParams, err := c.parseParams(params)
	if err != nil {
		return err
	}

	return c.validateDate(moonParams.Year, moonParams.Month, moonParams.Day)
}

// parseParams 解析参数
func (c *MoonPhaseCalculator) parseParams(params interface{}) (*MoonPhaseParams, error) {
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

	// 日参数可选，默认为1
	day := 1
	if d, exists := paramsMap["day"]; exists {
		if dFloat, ok := d.(float64); ok {
			day = int(dFloat)
		}
	}

	return &MoonPhaseParams{
		Year:  int(year),
		Month: int(month),
		Day:   day,
	}, nil
}

// validateDate 验证日期有效性
func (c *MoonPhaseCalculator) validateDate(year, month, day int) error {
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
func (c *MoonPhaseCalculator) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// calculateMoonPhase 计算月相信息
func (c *MoonPhaseCalculator) calculateMoonPhase(year, month, day int) (*MoonPhaseResult, error) {
	// 计算儒略日
	jd := c.calculateJulianDay(year, month, day)

	// 计算月球位置
	moonPosition := c.calculateMoonPosition(jd)

	// 计算月相
	phaseAngle, illumination := c.calculatePhase(jd)

	// 计算月龄
	age := c.calculateMoonAge(jd)

	// 判断月相类型
	moonPhase := c.determineMoonPhase(phaseAngle)

	// 计算下一个主要月相
	nextPhase, nextPhaseTime := c.calculateNextMajorPhase(jd)

	// 格式化日期
	dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	return &MoonPhaseResult{
		Date:          dateStr,
		MoonPhase:     moonPhase,
		PhaseAngle:    phaseAngle,
		Illumination:  illumination,
		Age:           age,
		NextPhase:     nextPhase,
		NextPhaseTime: nextPhaseTime,
		Distance:      moonPosition.Distance,
		Longitude:     moonPosition.Longitude,
		Latitude:      moonPosition.Latitude,
	}, nil
}

// calculateJulianDay 计算儒略日
func (c *MoonPhaseCalculator) calculateJulianDay(year, month, day int) float64 {
	if month <= 2 {
		year -= 1
		month += 12
	}

	a := year / 100
	b := 2 - a + a/4

	jd := math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		float64(day) + float64(b) - 1524.5

	return jd
}

// MoonPosition 月球位置
type MoonPosition struct {
	Longitude float64 // 黄经（度）
	Latitude  float64 // 黄纬（度）
	Distance  float64 // 距离（公里）
	RA        float64 // 赤经（度）
	Dec       float64 // 赤纬（度）
}

// calculateMoonPosition 计算月球位置（简化实现）
func (c *MoonPhaseCalculator) calculateMoonPosition(jd float64) *MoonPosition {
	// 简化的月球位置计算
	// 实际应该基于完整的月球理论

	t := (jd - 2451545.0) / 36525.0 // 儒略世纪数

	// 平黄经
	L := 218.3164591 + 481267.88134236*t - 0.0013268*t*t + t*t*t/538841 - t*t*t*t/65194000
	L = c.normalizeAngle(L)

	// 平近点角
	M := 134.9634114 + 477198.8676313*t + 0.008997*t*t + t*t*t/69699 - t*t*t*t/14712000
	M = c.normalizeAngle(M)

	// 平升交角距
	F := 93.2720993 + 483202.0175273*t - 0.0034029*t*t - t*t*t/3526000 + t*t*t*t/863310000
	F = c.normalizeAngle(F)

	// 简化位置计算
	longitude := L + 6.289*math.Sin(M*math.Pi/180)
	latitude := 5.128 * math.Sin(F*math.Pi/180)
	distance := 385000.56 + 20905.355*math.Cos(M*math.Pi/180) // 平均距离+修正

	return &MoonPosition{
		Longitude: c.normalizeAngle(longitude),
		Latitude:  latitude,
		Distance:  distance,
		RA:        c.normalizeAngle(longitude), // 简化
		Dec:       latitude,                    // 简化
	}
}

// calculatePhase 计算月相
func (c *MoonPhaseCalculator) calculatePhase(jd float64) (float64, float64) {
	// 计算月相角
	t := (jd - 2451545.0) / 36525.0

	// 太阳平黄经
	sunL := 280.46645 + 36000.76983*t + 0.0003032*t*t
	sunL = c.normalizeAngle(sunL)

	// 月球平黄经
	moonL := 218.3164591 + 481267.88134236*t
	moonL = c.normalizeAngle(moonL)

	// 相位角 = 月球黄经 - 太阳黄经
	phaseAngle := moonL - sunL
	if phaseAngle < 0 {
		phaseAngle += 360
	}

	// 照明比例 = (1 + cos(相位角)) / 2
	illumination := (1 + math.Cos(phaseAngle*math.Pi/180)) / 2

	return phaseAngle, illumination
}

// calculateMoonAge 计算月龄
func (c *MoonPhaseCalculator) calculateMoonAge(jd float64) float64 {
	// 计算上次新月的时间
	lastNewMoonJD := c.calculateLastNewMoon(jd)

	// 月龄 = 当前时间 - 上次新月时间
	age := jd - lastNewMoonJD

	return age
}

// calculateLastNewMoon 计算上次新月时间
func (c *MoonPhaseCalculator) calculateLastNewMoon(jd float64) float64 {
	// 简化的新月计算
	// 实际应该基于精确的新月算法

	synodicMonth := 29.530588853 // 朔望月长度（天）

	// 近似计算上次新月
	cycles := math.Floor((jd - 2451550.1) / synodicMonth)
	lastNewMoon := 2451550.1 + cycles*synodicMonth

	return lastNewMoon
}

// determineMoonPhase 判断月相类型
func (c *MoonPhaseCalculator) determineMoonPhase(phaseAngle float64) string {
	// 根据相位角判断月相类型

	if phaseAngle < 45 || phaseAngle >= 315 {
		return "new" // 新月
	} else if phaseAngle >= 45 && phaseAngle < 135 {
		return "first_quarter" // 上弦月
	} else if phaseAngle >= 135 && phaseAngle < 225 {
		return "full" // 满月
	} else {
		return "last_quarter" // 下弦月
	}
}

// calculateNextMajorPhase 计算下一个主要月相
func (c *MoonPhaseCalculator) calculateNextMajorPhase(jd float64) (string, string) {
	// 简化的下一个主要月相计算

	synodicMonth := 29.530588853 // 朔望月长度（天）

	// 计算当前月相
	phaseAngle, _ := c.calculatePhase(jd)
	currentPhase := c.determineMoonPhase(phaseAngle)

	// 计算下一个主要月相
	var nextPhase string
	var daysToNext float64

	switch currentPhase {
	case "new":
		nextPhase = "first_quarter"
		daysToNext = synodicMonth / 4
	case "first_quarter":
		nextPhase = "full"
		daysToNext = synodicMonth / 4
	case "full":
		nextPhase = "last_quarter"
		daysToNext = synodicMonth / 4
	case "last_quarter":
		nextPhase = "new"
		daysToNext = synodicMonth / 4
	default:
		nextPhase = "new"
		daysToNext = 7.0
	}

	// 计算下一个月相时间
	nextPhaseJD := jd + daysToNext
	nextPhaseTime := c.formatJulianDay(nextPhaseJD)

	return nextPhase, nextPhaseTime
}

// formatJulianDay 格式化儒略日
func (c *MoonPhaseCalculator) formatJulianDay(jd float64) string {
	// 儒略日转公历日期（简化）

	jd += 0.5
	Z := math.Floor(jd)
	F := jd - Z

	var A float64
	if Z < 2299161 {
		A = Z
	} else {
		alpha := math.Floor((Z - 1867216.25) / 36524.25)
		A = Z + 1 + alpha - math.Floor(alpha/4)
	}

	B := A + 1524
	C := math.Floor((B - 122.1) / 365.25)
	D := math.Floor(365.25 * C)
	E := math.Floor((B - D) / 30.6001)

	day := B - D - math.Floor(30.6001*E) + F
	month := 0
	if E < 14 {
		month = int(E - 1)
	} else {
		month = int(E - 13)
	}

	year := 0
	if month > 2 {
		year = int(C - 4716)
	} else {
		year = int(C - 4715)
	}

	// 提取时间部分
	hours := (day - math.Floor(day)) * 24
	hour := int(math.Floor(hours))
	minutes := (hours - float64(hour)) * 60
	minute := int(math.Floor(minutes))

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d", year, month, int(day), hour, minute)
}

// normalizeAngle 归一化角度到0-360度
func (c *MoonPhaseCalculator) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}


// GetMoonPhaseCalendar 获取月相日历（用于测试）
func (c *MoonPhaseCalculator) GetMoonPhaseCalendar(year, month int) ([]map[string]interface{}, error) {
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
		jd := c.calculateJulianDay(year, month, day)

		phaseAngle, illumination := c.calculatePhase(jd)
		moonPhase := c.determineMoonPhase(phaseAngle)
		age := c.calculateMoonAge(jd)

		dayInfo := map[string]interface{}{
			"date":         fmt.Sprintf("%d-%02d-%02d", year, month, day),
			"moon_phase":   moonPhase,
			"phase_angle":  phaseAngle,
			"illumination": illumination,
			"age":          age,
		}

		calendar = append(calendar, dayInfo)
	}

	return calendar, nil
}

// CalculateEclipse 计算月食信息（用于测试）
func (c *MoonPhaseCalculator) CalculateEclipse(year, month int) (map[string]interface{}, error) {
	// 简化的月食计算
	// 实际应该基于精确的月食算法

	if month < 1 || month > 12 {
		return nil, fmt.Errorf("月份超出范围: %d", month)
	}

	// 计算该月中间的儒略日
	midJD := c.calculateJulianDay(year, month, 15)
	_ = midJD // 标记为已使用
	
	// 简化的月食判断
	hasEclipse := false
	eclipseType := "none"

	// 基于月份和年份的简单模式
	if (year%4 == 0 && month == 6) || (year%7 == 0 && month == 12) {
		hasEclipse = true
		eclipseType = "partial"
	}

	return map[string]interface{}{
		"year":         year,
		"month":        month,
		"has_eclipse":  hasEclipse,
		"eclipse_type": eclipseType,
		"eclipse_date": fmt.Sprintf("%d-%02d-15", year, month),
	}, nil
}
