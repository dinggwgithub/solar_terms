package calculator

import (
	"fmt"
	"math"
	"scientific_calc_bugs/internal/bugs"
)

// SunriseSunsetCalculator 日出日落时间计算器
type SunriseSunsetCalculator struct {
	*BaseCalculator
}

// NewSunriseSunsetCalculator 创建新的日出日落时间计算器
func NewSunriseSunsetCalculator() *SunriseSunsetCalculator {
	return &SunriseSunsetCalculator{
		BaseCalculator: NewBaseCalculator(
			"sunrise_sunset",
			"日出日落时间计算器，基于地理位置和日期计算日出日落时刻",
		),
	}
}

// SunriseSunsetParams 日出日落计算参数
type SunriseSunsetParams struct {
	Year      int     `json:"year"`      // 年
	Month     int     `json:"month"`     // 月
	Day       int     `json:"day"`       // 日
	Longitude float64 `json:"longitude"` // 经度（度）
	Latitude  float64 `json:"latitude"`  // 纬度（度）
	Altitude  float64 `json:"altitude"`  // 海拔（米，可选）
	Timezone  float64 `json:"timezone"`  // 时区（小时，可选）
}

// SunriseSunsetResult 日出日落计算结果
type SunriseSunsetResult struct {
	Date          string  `json:"date"`       // 日期
	Sunrise       string  `json:"sunrise"`    // 日出时间
	Sunset        string  `json:"sunset"`     // 日落时间
	SolarNoon     string  `json:"solar_noon"` // 太阳正午
	DayLength     float64 `json:"day_length"` // 白昼长度（小时）
	CivilTwilight struct {
		Morning string `json:"morning"` // 民用晨光开始
		Evening string `json:"evening"` // 民用暮光结束
	} `json:"civil_twilight"`
	NauticalTwilight struct {
		Morning string `json:"morning"` // 航海晨光开始
		Evening string `json:"evening"` // 航海暮光结束
	} `json:"nautical_twilight"`
	AstronomicalTwilight struct {
		Morning string `json:"morning"` // 天文晨光开始
		Evening string `json:"evening"` // 天文暮光结束
	} `json:"astronomical_twilight"`
}

// Calculate 执行日出日落计算
func (c *SunriseSunsetCalculator) Calculate(params interface{}) (interface{}, error) {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(sunriseParams); err != nil {
		return nil, err
	}

	// 执行日出日落计算
	result, err := c.calculateSunriseSunset(sunriseParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *SunriseSunsetCalculator) Validate(params interface{}) error {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return err
	}

	return c.validateParams(sunriseParams)
}

// parseParams 解析参数
func (c *SunriseSunsetCalculator) parseParams(params interface{}) (*SunriseSunsetParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 提取必需参数
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

	longitude, ok := paramsMap["longitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("longitude参数必须为数字")
	}

	latitude, ok := paramsMap["latitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("latitude参数必须为数字")
	}

	// 提取可选参数
	altitude := 0.0
	if alt, exists := paramsMap["altitude"]; exists {
		if altFloat, ok := alt.(float64); ok {
			altitude = altFloat
		}
	}

	timezone := 8.0 // 默认北京时间
	if tz, exists := paramsMap["timezone"]; exists {
		if tzFloat, ok := tz.(float64); ok {
			timezone = tzFloat
		}
	}

	return &SunriseSunsetParams{
		Year:      int(year),
		Month:     int(month),
		Day:       int(day),
		Longitude: longitude,
		Latitude:  latitude,
		Altitude:  altitude,
		Timezone:  timezone,
	}, nil
}

// validateParams 验证参数
func (c *SunriseSunsetCalculator) validateParams(params *SunriseSunsetParams) error {
	// 验证日期
	if err := c.validateDate(params.Year, params.Month, params.Day); err != nil {
		return err
	}

	// 验证经度范围
	if params.Longitude < -180 || params.Longitude > 180 {
		return fmt.Errorf("经度超出范围 (-180到180): %f", params.Longitude)
	}

	// 验证纬度范围
	if params.Latitude < -90 || params.Latitude > 90 {
		return fmt.Errorf("纬度超出范围 (-90到90): %f", params.Latitude)
	}

	// 验证海拔（可选）
	if params.Altitude < -1000 || params.Altitude > 10000 {
		return fmt.Errorf("海拔超出合理范围 (-1000到10000米): %f", params.Altitude)
	}

	// 验证时区（可选）
	if params.Timezone < -12 || params.Timezone > 14 {
		return fmt.Errorf("时区超出范围 (-12到14): %f", params.Timezone)
	}

	return nil
}

// validateDate 验证日期有效性
func (c *SunriseSunsetCalculator) validateDate(year, month, day int) error {
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
func (c *SunriseSunsetCalculator) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// calculateSunriseSunset 计算日出日落时间
func (c *SunriseSunsetCalculator) calculateSunriseSunset(params *SunriseSunsetParams) (*SunriseSunsetResult, error) {
	// 计算儒略日
	jd := c.calculateJulianDay(params.Year, params.Month, params.Day)

	// 计算太阳位置
	solarPosition := c.calculateSolarPosition(jd)

	// 计算日出日落时间
	sunrise, sunset, err := c.calculateSunTimes(jd, params.Longitude, params.Latitude, params.Altitude)
	if err != nil {
		return nil, err
	}

	// 计算太阳正午
	solarNoon := c.calculateSolarNoon(jd, params.Longitude)

	// 计算白昼长度
	dayLength := c.calculateDayLength(sunrise, sunset)

	// 计算晨昏蒙影
	twilightTimes := c.calculateTwilightTimes(jd, params.Longitude, params.Latitude)

	// 格式化时间
	dateStr := fmt.Sprintf("%d-%02d-%02d", params.Year, params.Month, params.Day)
	sunriseStr := c.formatTime(sunrise, params.Timezone)
	sunsetStr := c.formatTime(sunset, params.Timezone)
	solarNoonStr := c.formatTime(solarNoon, params.Timezone)

	result := &SunriseSunsetResult{
		Date:      dateStr,
		Sunrise:   sunriseStr,
		Sunset:    sunsetStr,
		SolarNoon: solarNoonStr,
		DayLength: dayLength,
	}

	// 设置晨昏蒙影时间
	result.CivilTwilight.Morning = c.formatTime(twilightTimes.civilMorning, params.Timezone)
	result.CivilTwilight.Evening = c.formatTime(twilightTimes.civilEvening, params.Timezone)
	result.NauticalTwilight.Morning = c.formatTime(twilightTimes.nauticalMorning, params.Timezone)
	result.NauticalTwilight.Evening = c.formatTime(twilightTimes.nauticalEvening, params.Timezone)
	result.AstronomicalTwilight.Morning = c.formatTime(twilightTimes.astronomicalMorning, params.Timezone)
	result.AstronomicalTwilight.Evening = c.formatTime(twilightTimes.astronomicalEvening, params.Timezone)

	_ = solarPosition // 使用变量避免警告
	return result, nil
}

// calculateJulianDay 计算儒略日
func (c *SunriseSunsetCalculator) calculateJulianDay(year, month, day int) float64 {
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

// SolarPosition 太阳位置
type SolarPosition struct {
	RightAscension float64 // 赤经（度）
	Declination    float64 // 赤纬（度）
	Distance       float64 // 距离（天文单位）
}

// calculateSolarPosition 计算太阳位置
func (c *SunriseSunsetCalculator) calculateSolarPosition(jd float64) *SolarPosition {
	// 简化的太阳位置计算
	// 实际应该基于VSOP87理论

	t := (jd - 2451545.0) / 36525.0 // 儒略世纪数

	// 平黄经
	L := 280.46646 + 36000.76983*t + 0.0003032*t*t
	L = c.normalizeAngle(L)

	// 平近点角
	M := 357.52911 + 35999.05029*t - 0.0001537*t*t
	M = c.normalizeAngle(M)

	// 中心差
	C := (1.914602-0.004817*t-0.000014*t*t)*math.Sin(M*math.Pi/180) +
		(0.019993-0.000101*t)*math.Sin(2*M*math.Pi/180) +
		0.000289*math.Sin(3*M*math.Pi/180)

	// 真黄经
	lambda := L + C

	// 黄赤交角
	epsilon := 23.4392911 - 0.0130042*t - 0.00000016*t*t + 0.000000504*t*t*t

	// 赤经赤纬
	alpha := math.Atan2(
		math.Sin(lambda*math.Pi/180)*math.Cos(epsilon*math.Pi/180),
		math.Cos(lambda*math.Pi/180),
	) * 180 / math.Pi

	delta := math.Asin(
		math.Sin(lambda*math.Pi/180)*math.Sin(epsilon*math.Pi/180),
	) * 180 / math.Pi

	return &SolarPosition{
		RightAscension: alpha,
		Declination:    delta,
		Distance:       1.0, // 简化
	}
}

// calculateSunTimes 计算日出日落时间
func (c *SunriseSunsetCalculator) calculateSunTimes(jd, longitude, latitude, altitude float64) (float64, float64, error) {
	// 简化的日出日落计算
	// 实际应该基于精确的天文算法

	// 计算太阳时角
	h0 := -0.8333 // 日出日落时太阳在地平线下的角度（度）

	// 修正海拔影响
	if altitude > 0 {
		h0 -= 0.0347 * math.Sqrt(altitude)
	}

	// 计算时角
	cosH := (math.Sin(h0*math.Pi/180) -
		math.Sin(latitude*math.Pi/180)*math.Sin(0*math.Pi/180)) /
		(math.Cos(latitude*math.Pi/180) * math.Cos(0*math.Pi/180))

	if cosH < -1 || cosH > 1 {
		return 0, 0, fmt.Errorf("在该位置该日期没有日出或日落")
	}

	H := math.Acos(cosH) * 180 / math.Pi

	// 计算日出日落时间（儒略日）
	sunriseJD := jd + (720-4*(longitude+H)-0)/1440.0
	sunsetJD := jd + (720-4*(longitude-H)-0)/1440.0

	return sunriseJD, sunsetJD, nil
}

// calculateSolarNoon 计算太阳正午
func (c *SunriseSunsetCalculator) calculateSolarNoon(jd, longitude float64) float64 {
	// 太阳正午（儒略日）
	return jd + (720-4*longitude-0)/1440.0
}

// calculateDayLength 计算白昼长度
func (c *SunriseSunsetCalculator) calculateDayLength(sunriseJD, sunsetJD float64) float64 {
	return (sunsetJD - sunriseJD) * 24
}

// TwilightTimes 晨昏蒙影时间
type TwilightTimes struct {
	civilMorning        float64 // 民用晨光开始
	civilEvening        float64 // 民用暮光结束
	nauticalMorning     float64 // 航海晨光开始
	nauticalEvening     float64 // 航海暮光结束
	astronomicalMorning float64 // 天文晨光开始
	astronomicalEvening float64 // 天文暮光结束
}

// calculateTwilightTimes 计算晨昏蒙影时间
func (c *SunriseSunsetCalculator) calculateTwilightTimes(jd, longitude, latitude float64) *TwilightTimes {
	// 简化的晨昏蒙影计算

	times := &TwilightTimes{}

	// 民用晨昏蒙影（太阳在地平线下6度）
	times.civilMorning, times.civilEvening = c.calculateTwilight(jd, longitude, latitude, -6)

	// 航海晨昏蒙影（太阳在地平线下12度）
	times.nauticalMorning, times.nauticalEvening = c.calculateTwilight(jd, longitude, latitude, -12)

	// 天文晨昏蒙影（太阳在地平线下18度）
	times.astronomicalMorning, times.astronomicalEvening = c.calculateTwilight(jd, longitude, latitude, -18)

	return times
}

// calculateTwilight 计算特定角度的晨昏蒙影
func (c *SunriseSunsetCalculator) calculateTwilight(jd, longitude, latitude, angle float64) (float64, float64) {
	// 简化的晨昏蒙影计算

	cosH := (math.Sin(angle*math.Pi/180) -
		math.Sin(latitude*math.Pi/180)*math.Sin(0*math.Pi/180)) /
		(math.Cos(latitude*math.Pi/180) * math.Cos(0*math.Pi/180))

	if cosH < -1 || cosH > 1 {
		return 0, 0
	}

	H := math.Acos(cosH) * 180 / math.Pi

	morningJD := jd + (720-4*(longitude+H)-0)/1440.0
	eveningJD := jd + (720-4*(longitude-H)-0)/1440.0

	return morningJD, eveningJD
}

// formatTime 格式化时间
func (c *SunriseSunsetCalculator) formatTime(jd, timezone float64) string {
	// 儒略日转时间
	hours := (jd - math.Floor(jd)) * 24
	hours += timezone // 时区修正

	if hours < 0 {
		hours += 24
	} else if hours >= 24 {
		hours -= 24
	}

	hour := int(math.Floor(hours))
	minute := int(math.Floor((hours - float64(hour)) * 60))

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// normalizeAngle 归一化角度到0-360度
func (c *SunriseSunsetCalculator) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// GetSupportedBugTypes 返回支持的Bug类型
func (c *SunriseSunsetCalculator) GetSupportedBugTypes() []bugs.BugType {
	return []bugs.BugType{
		bugs.BugTypeInstability,
		bugs.BugTypeConstraint,
		bugs.BugTypePrecision,
	}
}

// GetLocationInfo 获取地理位置信息（用于测试）
func (c *SunriseSunsetCalculator) GetLocationInfo(longitude, latitude float64) map[string]interface{} {
	info := map[string]interface{}{
		"longitude":           longitude,
		"latitude":            latitude,
		"hemisphere":          "",
		"timezone_suggestion": 0.0,
	}

	// 判断半球
	if latitude >= 0 {
		info["hemisphere"] = "北半球"
	} else {
		info["hemisphere"] = "南半球"
	}

	// 建议时区（基于经度）
	info["timezone_suggestion"] = math.Round(longitude / 15)

	return info
}
