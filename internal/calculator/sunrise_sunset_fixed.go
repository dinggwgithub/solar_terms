package calculator

import (
	"fmt"
	"math"
)

// SunriseSunsetCalculatorFixed 修复后的日出日落时间计算器
type SunriseSunsetCalculatorFixed struct {
	*BaseCalculator
}

// NewSunriseSunsetCalculatorFixed 创建新的修复后的日出日落时间计算器
func NewSunriseSunsetCalculatorFixed() *SunriseSunsetCalculatorFixed {
	return &SunriseSunsetCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"sunrise_sunset_fixed",
			"修复后的日出日落时间计算器，基于正确的天文算法",
		),
	}
}

// Calculate 执行日出日落计算
func (c *SunriseSunsetCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(sunriseParams); err != nil {
		return nil, err
	}

	// 执行日出日落计算
	result, err := c.calculateSunriseSunsetFixed(sunriseParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *SunriseSunsetCalculatorFixed) Validate(params interface{}) error {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return err
	}

	return c.validateParams(sunriseParams)
}

// parseParams 解析参数
func (c *SunriseSunsetCalculatorFixed) parseParams(params interface{}) (*SunriseSunsetParams, error) {
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
func (c *SunriseSunsetCalculatorFixed) validateParams(params *SunriseSunsetParams) error {
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
func (c *SunriseSunsetCalculatorFixed) validateDate(year, month, day int) error {
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
func (c *SunriseSunsetCalculatorFixed) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// calculateSunriseSunsetFixed 修复后的日出日落计算
func (c *SunriseSunsetCalculatorFixed) calculateSunriseSunsetFixed(params *SunriseSunsetParams) (*SunriseSunsetResult, error) {
	// 计算儒略日（中午12点）
	jd := c.calculateJulianDay(params.Year, params.Month, params.Day)

	// 计算日出日落时应该使用午夜作为基准（NOAA算法要求）
	// 标准儒略日: xxxx.5 表示中午，xxxx.0 表示午夜
	// NOAA 算法需要的是从午夜开始的儒略日（减去0.5天）
	jdMidnight := jd - 0.5 // 转换为午夜基准（xxxx.0）

	// 计算日出日落时间（使用午夜基准）
	sunrise, sunset, err := c.calculateSunTimesFixed(jdMidnight, params.Longitude, params.Latitude, params.Altitude)
	if err != nil {
		return nil, err
	}

	// 计算太阳正午（使用午夜基准）
	solarNoon := c.calculateSolarNoonFixed(jdMidnight, params.Longitude)

	// 计算白昼长度
	dayLength := c.calculateDayLength(sunrise, sunset)

	// 计算晨昏蒙影（使用午夜基准）
	twilightTimes := c.calculateTwilightTimesFixed(jdMidnight, params.Longitude, params.Latitude)

	// 格式化时间
	dateStr := fmt.Sprintf("%d-%02d-%02d", params.Year, params.Month, params.Day)
	sunriseStr := c.formatTimeFixed(sunrise, params.Timezone)
	sunsetStr := c.formatTimeFixed(sunset, params.Timezone)
	solarNoonStr := c.formatTimeFixed(solarNoon, params.Timezone)

	result := &SunriseSunsetResult{
		Date:      dateStr,
		Sunrise:   sunriseStr,
		Sunset:    sunsetStr,
		SolarNoon: solarNoonStr,
		DayLength: dayLength,
	}

	// 设置晨昏蒙影时间
	result.CivilTwilight.Morning = c.formatTimeFixed(twilightTimes.civilMorning, params.Timezone)
	result.CivilTwilight.Evening = c.formatTimeFixed(twilightTimes.civilEvening, params.Timezone)
	result.NauticalTwilight.Morning = c.formatTimeFixed(twilightTimes.nauticalMorning, params.Timezone)
	result.NauticalTwilight.Evening = c.formatTimeFixed(twilightTimes.nauticalEvening, params.Timezone)
	result.AstronomicalTwilight.Morning = c.formatTimeFixed(twilightTimes.astronomicalMorning, params.Timezone)
	result.AstronomicalTwilight.Evening = c.formatTimeFixed(twilightTimes.astronomicalEvening, params.Timezone)

	return result, nil
}

// calculateJulianDay 计算儒略日（中午12点）
func (c *SunriseSunsetCalculatorFixed) calculateJulianDay(year, month, day int) float64 {
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

// calculateSunTimesFixed 修复后的日出日落时间计算
func (c *SunriseSunsetCalculatorFixed) calculateSunTimesFixed(jd, longitude, latitude, altitude float64) (float64, float64, error) {
	// 使用更精确的算法：NOAA日出日落计算方法
	// 参考：https://www.esrl.noaa.gov/gmd/grad/solcalc/calcdetails.html

	// 日出日落时太阳中心在地平线下0.833度（考虑大气折射和太阳视直径）
	zenith := 90.8333

	// 修正海拔影响
	if altitude > 0 {
		zenith += 0.0347 * math.Sqrt(altitude) / 60.0 // 转换为度
	}

	// 计算儒略世纪
	T := (jd - 2451545.0) / 36525.0

	// 计算太阳几何平黄经
	L := c.normalizeAngle(280.46646 + 36000.76983*T + 0.0003032*T*T)

	// 计算太阳平近点角
	M := c.normalizeAngle(357.52911 + 35999.05029*T - 0.0001537*T*T)

	// 计算太阳中心差
	C := (1.914602-0.004817*T-0.000014*T*T)*math.Sin(M*math.Pi/180) +
		(0.019993-0.000101*T)*math.Sin(2*M*math.Pi/180) +
		0.000289*math.Sin(3*M*math.Pi/180)

	// 太阳真黄经
	longitudeSun := L + C

	// 黄赤交角
	epsilon := 23 + (26+(21.448-46.8150*T-0.00059*T*T+0.001813*T*T*T)/60)/60

	// 太阳赤纬
	delta := math.Asin(math.Sin(epsilon*math.Pi/180)*math.Sin(longitudeSun*math.Pi/180)) * 180 / math.Pi

	// 计算时角
	cosH := (math.Cos(zenith*math.Pi/180) -
		math.Sin(latitude*math.Pi/180)*math.Sin(delta*math.Pi/180)) /
		(math.Cos(latitude*math.Pi/180) * math.Cos(delta*math.Pi/180))

	if cosH < -1 || cosH > 1 {
		return 0, 0, fmt.Errorf("在该位置该日期没有日出或日落（极昼或极夜）")
	}

	H := math.Acos(cosH) * 180 / math.Pi

	// 太阳正午（相对于本初子午线）
	noonUTC := c.calculateSolarNoonUTC(jd, longitude)

	// 日出日落的UTC时间
	sunriseUTC := noonUTC - H*4/1440.0
	sunsetUTC := noonUTC + H*4/1440.0

	return sunriseUTC, sunsetUTC, nil
}

// calculateSolarNoonFixed 修复后的太阳正午计算
func (c *SunriseSunsetCalculatorFixed) calculateSolarNoonFixed(jd, longitude float64) float64 {
	return c.calculateSolarNoonUTC(jd, longitude)
}

// calculateSolarNoonUTC 计算太阳正午（UTC时间，儒略日）
func (c *SunriseSunsetCalculatorFixed) calculateSolarNoonUTC(jd, longitude float64) float64 {
	// 计算儒略世纪
	T := (jd - 2451545.0) / 36525.0

	// 计算太阳几何平黄经
	L := c.normalizeAngle(280.46646 + 36000.76983*T + 0.0003032*T*T)

	// 计算太阳平近点角
	M := c.normalizeAngle(357.52911 + 35999.05029*T - 0.0001537*T*T)

	// 计算太阳中心差
	C := (1.914602-0.004817*T-0.000014*T*T)*math.Sin(M*math.Pi/180) +
		(0.019993-0.000101*T)*math.Sin(2*M*math.Pi/180) +
		0.000289*math.Sin(3*M*math.Pi/180)

	// 太阳真黄经
	longitudeSun := L + C

	// 黄赤交角
	epsilon := 23 + (26+(21.448-46.8150*T-0.00059*T*T+0.001813*T*T*T)/60)/60

	// 太阳赤经
	alpha := math.Atan2(math.Cos(epsilon*math.Pi/180)*math.Sin(longitudeSun*math.Pi/180), math.Cos(longitudeSun*math.Pi/180)) * 180 / math.Pi
	alpha = c.normalizeAngle(alpha)

	// 赤经和赤纬的四象限修正
	if math.Floor(L/90) != math.Floor(alpha/90) {
		if alpha < L {
			alpha += 360
		} else {
			alpha -= 360
		}
	}

	// 时差（分钟）
	E := (L - alpha) * 4

	// 太阳正午（UTC时间，儒略日）
	// 正确公式: 分钟数 = 720 - 4*longitude - E （UTC分钟）
	// 转换为儒略日: jd + 分钟数 / 1440.0
	noonMinutesUTC := 720 - 4*longitude - E
	noonUTC := jd + noonMinutesUTC/1440.0

	return noonUTC
}

// calculateDayLength 计算白昼长度
func (c *SunriseSunsetCalculatorFixed) calculateDayLength(sunriseJD, sunsetJD float64) float64 {
	return (sunsetJD - sunriseJD) * 24
}

// TwilightTimesFixed 晨昏蒙影时间
type TwilightTimesFixed struct {
	civilMorning        float64 // 民用晨光开始
	civilEvening        float64 // 民用暮光结束
	nauticalMorning     float64 // 航海晨光开始
	nauticalEvening     float64 // 航海暮光结束
	astronomicalMorning float64 // 天文晨光开始
	astronomicalEvening float64 // 天文暮光结束
}

// calculateTwilightTimesFixed 修复后的晨昏蒙影计算
func (c *SunriseSunsetCalculatorFixed) calculateTwilightTimesFixed(jd, longitude, latitude float64) *TwilightTimesFixed {
	times := &TwilightTimesFixed{}

	// 民用晨昏蒙影（太阳在地平线下6度）
	times.civilMorning, times.civilEvening = c.calculateTwilight(jd, longitude, latitude, 96.0)

	// 航海晨昏蒙影（太阳在地平线下12度）
	times.nauticalMorning, times.nauticalEvening = c.calculateTwilight(jd, longitude, latitude, 102.0)

	// 天文晨昏蒙影（太阳在地平线下18度）
	times.astronomicalMorning, times.astronomicalEvening = c.calculateTwilight(jd, longitude, latitude, 108.0)

	return times
}

// calculateTwilight 修复后的晨昏蒙影计算
func (c *SunriseSunsetCalculatorFixed) calculateTwilight(jd, longitude, latitude, zenith float64) (float64, float64) {
	// 计算儒略世纪
	T := (jd - 2451545.0) / 36525.0

	// 计算太阳几何平黄经
	L := c.normalizeAngle(280.46646 + 36000.76983*T + 0.0003032*T*T)

	// 计算太阳平近点角
	M := c.normalizeAngle(357.52911 + 35999.05029*T - 0.0001537*T*T)

	// 计算太阳中心差
	C := (1.914602-0.004817*T-0.000014*T*T)*math.Sin(M*math.Pi/180) +
		(0.019993-0.000101*T)*math.Sin(2*M*math.Pi/180) +
		0.000289*math.Sin(3*M*math.Pi/180)

	// 太阳真黄经
	longitudeSun := L + C

	// 黄赤交角
	epsilon := 23 + (26+(21.448-46.8150*T-0.00059*T*T+0.001813*T*T*T)/60)/60

	// 太阳赤纬
	delta := math.Asin(math.Sin(epsilon*math.Pi/180)*math.Sin(longitudeSun*math.Pi/180)) * 180 / math.Pi

	// 计算时角
	cosH := (math.Cos(zenith*math.Pi/180) -
		math.Sin(latitude*math.Pi/180)*math.Sin(delta*math.Pi/180)) /
		(math.Cos(latitude*math.Pi/180) * math.Cos(delta*math.Pi/180))

	if cosH < -1 || cosH > 1 {
		return 0, 0
	}

	H := math.Acos(cosH) * 180 / math.Pi

	// 太阳正午
	noonUTC := c.calculateSolarNoonUTC(jd, longitude)

	// 晨昏蒙影时间（UTC）
	morningUTC := noonUTC - H*4/1440.0
	eveningUTC := noonUTC + H*4/1440.0

	return morningUTC, eveningUTC
}

// formatTimeFixed 修复后的时间格式化（输入为UTC儒略日，午夜基准）
func (c *SunriseSunsetCalculatorFixed) formatTimeFixed(jd, timezone float64) string {
	// 儒略日转UTC时间
	// 输入jd是午夜基准: xxxx.0 表示午夜0点，xxxx.5 表示中午12点
	dayFraction := jd - math.Floor(jd)
	hoursUTC := dayFraction * 24

	// 时区修正
	hours := hoursUTC + timezone

	// 确保在0-24小时范围内
	if hours < 0 {
		hours += 24
	} else if hours >= 24 {
		hours -= 24
	}

	hour := int(math.Floor(hours))
	minute := int(math.Round((hours - float64(hour)) * 60))

	// 处理分钟进位
	if minute >= 60 {
		minute -= 60
		hour++
		if hour >= 24 {
			hour -= 24
		}
	}

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// normalizeAngle 归一化角度到0-360度
func (c *SunriseSunsetCalculatorFixed) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}
