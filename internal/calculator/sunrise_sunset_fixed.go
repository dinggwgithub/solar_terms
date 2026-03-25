package calculator

import (
	"fmt"
	"math"
)

// SunriseSunsetCalculatorFixed 修复后的日出日落时间计算器
// 修复了原版本中的以下Bug：
// 1. 使用固定的太阳赤纬0度，而不是根据日期计算实际赤纬
// 2. 时区转换逻辑错误
// 3. 日出日落公式中的符号处理问题
type SunriseSunsetCalculatorFixed struct {
	*BaseCalculator
}

// NewSunriseSunsetCalculatorFixed 创建修复后的日出日落时间计算器
func NewSunriseSunsetCalculatorFixed() *SunriseSunsetCalculatorFixed {
	return &SunriseSunsetCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"sunrise_sunset_fixed",
			"日出日落时间计算器（修复版），基于专业天文算法",
		),
	}
}

// SunriseSunsetResultFixed 日出日落计算结果（与原版兼容）
type SunriseSunsetResultFixed struct {
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

// Calculate 执行日出日落计算（修复版）
func (c *SunriseSunsetCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(sunriseParams); err != nil {
		return nil, err
	}

	// 执行修复后的日出日落计算
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

// ============================================
// 修复后的核心天文计算算法
// ============================================

// calculateSunriseSunsetFixed 修复后的日出日落时间计算
func (c *SunriseSunsetCalculatorFixed) calculateSunriseSunsetFixed(params *SunriseSunsetParams) (*SunriseSunsetResultFixed, error) {
	// 使用改进的天文算法计算日出日落
	// 基于NOAA的日出日落计算公式

	year := params.Year
	month := params.Month
	day := params.Day
	lat := params.Latitude
	lng := params.Longitude
	timezone := params.Timezone

	// 计算儒略日（正午UTC）
	jd := c.calculateJulianDayNoon(year, month, day)

	// 计算太阳位置参数
	solarPos := c.calculateSolarPositionFixed(jd)

	// 计算日出日落时间（UTC，以小时为单位，从午夜开始）
	sunriseUTC, sunsetUTC, solarNoonUTC, err := c.calculateSunTimesFixed(
		jd, lng, lat, solarPos.Declination, params.Altitude,
	)
	if err != nil {
		return nil, err
	}

	// 计算白昼长度
	dayLength := c.calculateDayLengthFixed(sunriseUTC, sunsetUTC)

	// 计算晨昏蒙影时间
	twilightTimes := c.calculateTwilightTimesFixed(jd, lng, lat, solarPos.Declination)

	// 将UTC时间转换为本地时间
	dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	result := &SunriseSunsetResultFixed{
		Date:      dateStr,
		Sunrise:   c.hoursToTimeString(sunriseUTC, timezone),
		Sunset:    c.hoursToTimeString(sunsetUTC, timezone),
		SolarNoon: c.hoursToTimeString(solarNoonUTC, timezone),
		DayLength: dayLength,
	}

	// 设置晨昏蒙影时间
	result.CivilTwilight.Morning = c.hoursToTimeString(twilightTimes.civilMorning, timezone)
	result.CivilTwilight.Evening = c.hoursToTimeString(twilightTimes.civilEvening, timezone)
	result.NauticalTwilight.Morning = c.hoursToTimeString(twilightTimes.nauticalMorning, timezone)
	result.NauticalTwilight.Evening = c.hoursToTimeString(twilightTimes.nauticalEvening, timezone)
	result.AstronomicalTwilight.Morning = c.hoursToTimeString(twilightTimes.astronomicalMorning, timezone)
	result.AstronomicalTwilight.Evening = c.hoursToTimeString(twilightTimes.astronomicalEvening, timezone)

	return result, nil
}

// SolarPositionFixed 太阳位置（修复版）
type SolarPositionFixed struct {
	RightAscension float64 // 赤经（度）
	Declination    float64 // 赤纬（度）- 这是关键参数！
	Distance       float64 // 距离（天文单位）
}

// calculateJulianDayNoon 计算儒略日（当天正午UTC）
func (c *SunriseSunsetCalculatorFixed) calculateJulianDayNoon(year, month, day int) float64 {
	// 如果月份是1月或2月，转换为上一年的13月或14月
	if month <= 2 {
		year -= 1
		month += 12
	}

	A := year / 100
	B := 2 - A + A/4

	// 计算当天正午（12:00 UTC）的儒略日
	jd := math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		float64(day) + float64(B) - 1524.5

	return jd
}

// calculateSolarPositionFixed 计算太阳位置（修复版）
// 这是关键修复：正确计算太阳赤纬，而不是使用固定的0度
func (c *SunriseSunsetCalculatorFixed) calculateSolarPositionFixed(jd float64) *SolarPositionFixed {
	// 基于NOAA的太阳能计算公式
	// https://gml.noaa.gov/grad/solcalc/calcdetails.html

	// 儒略世纪数（从J2000.0起算）
	julianCentury := (jd - 2451545.0) / 36525.0

	// 几何平黄经（度）
	geomMeanLongSun := 280.46646 + julianCentury*(36000.76983+0.0003032*julianCentury)
	geomMeanLongSun = c.normalizeAngleFixed(geomMeanLongSun)

	// 几何平近点角（度）
	geomMeanAnomalySun := 357.52911 + julianCentury*(35999.05029-0.0001537*julianCentury)
	geomMeanAnomalySun = c.normalizeAngleFixed(geomMeanAnomalySun)

	// 地球轨道离心率
	eccentEarthOrbit := 0.016708634 - julianCentury*(0.000042037+0.0000001267*julianCentury)

	// 太阳中心差（度）
	sunEqOfCtr := math.Sin(geomMeanAnomalySun*math.Pi/180.0)*
		(1.914602-julianCentury*(0.004817+0.000014*julianCentury)) +
		math.Sin(2*geomMeanAnomalySun*math.Pi/180.0)*(0.019993-0.000101*julianCentury) +
		math.Sin(3*geomMeanAnomalySun*math.Pi/180.0)*0.000289

	// 太阳真黄经（度）
	sunTrueLong := geomMeanLongSun + sunEqOfCtr

	// 太阳真近点角（度）
	sunTrueAnomaly := geomMeanAnomalySun + sunEqOfCtr

	// 太阳-地球距离（天文单位）
	sunRadVector := (1.000001018 * (1 - eccentEarthOrbit*eccentEarthOrbit)) /
		(1 + eccentEarthOrbit*math.Cos(sunTrueAnomaly*math.Pi/180.0))

	// 太阳视黄经（度）
	sunApparentLong := sunTrueLong - 0.00569 - 0.00478*
		math.Sin((125.04-1934.136*julianCentury)*math.Pi/180.0)

	// 平均黄赤交角（度）
	meanObliqEcliptic := 23.0 + (26.0+
		((21.448-julianCentury*(46.815+julianCentury*(0.00059-julianCentury*0.001813))))/60.0)/60.0

	// 修正黄赤交角（度）
	obliquityCorrection := meanObliqEcliptic +
		0.00256*math.Cos((125.04-1934.136*julianCentury)*math.Pi/180.0)

	// 太阳赤纬（度）- 这是关键参数！
	sunDeclination := math.Asin(math.Sin(obliquityCorrection*math.Pi/180.0)*
		math.Sin(sunApparentLong*math.Pi/180.0)) * 180.0 / math.Pi

	// 太阳赤经（度）
	sunRightAscension := math.Atan2(
		math.Cos(obliquityCorrection*math.Pi/180.0)*math.Sin(sunApparentLong*math.Pi/180.0),
		math.Cos(sunApparentLong*math.Pi/180.0),
	) * 180.0 / math.Pi
	sunRightAscension = c.normalizeAngleFixed(sunRightAscension)

	return &SolarPositionFixed{
		RightAscension: sunRightAscension,
		Declination:    sunDeclination,
		Distance:       sunRadVector,
	}
}

// calculateSunTimesFixed 修复后的日出日落时间计算
func (c *SunriseSunsetCalculatorFixed) calculateSunTimesFixed(
	jd, longitude, latitude, declination, altitude float64,
) (sunrise, sunset, solarNoon float64, err error) {
	// 日出日落时太阳在地平线下的角度（度）
	// 标准大气折射下的日出日落：-0.8333度
	sunriseSunsetAngle := -0.8333

	// 海拔修正
	if altitude > 0 {
		// 海拔越高，日出越早，日落越晚
		// 修正角度（度）
		altitudeCorrection := -0.0347 * math.Sqrt(altitude)
		sunriseSunsetAngle += altitudeCorrection
	}

	// 计算时角（度）
	// cos(H) = (cos(zenith) - sin(lat)*sin(dec)) / (cos(lat)*cos(dec))
	latRad := latitude * math.Pi / 180.0
	decRad := declination * math.Pi / 180.0
	zenithRad := (90.0 + sunriseSunsetAngle) * math.Pi / 180.0

	cosH := (math.Cos(zenithRad) - math.Sin(latRad)*math.Sin(decRad)) /
		(math.Cos(latRad) * math.Cos(decRad))

	if cosH > 1.0 {
		// 极夜：该日没有日出
		return 0, 0, 0, fmt.Errorf("该日期在该位置为极夜，没有日出")
	}
	if cosH < -1.0 {
		// 极昼：该日没有日落
		return 0, 0, 0, fmt.Errorf("该日期在该位置为极昼，没有日落")
	}

	// 时角（度）
	H := math.Acos(cosH) * 180.0 / math.Pi

	// 计算太阳正午的本地视太阳时（LST）
	// 首先计算方程时差（equation of time）
	equationOfTime := c.calculateEquationOfTime(jd)

	// 太阳正午的本地视太阳时（小时）= 12:00
	// 转换为平太阳时需要减去方程时差
	solarNoonLST := 12.0 - equationOfTime/60.0

	// 将本地视太阳时转换为UTC时间
	// UTC = LST - 经度/15（因为地球每小时转15度）
	solarNoonUTC := solarNoonLST - longitude/15.0

	// 日出和日落的本地视太阳时（小时）
	sunriseLST := solarNoonLST - H/15.0
	sunsetLST := solarNoonLST + H/15.0

	// 转换为UTC时间
	sunriseUTC := sunriseLST - longitude/15.0
	sunsetUTC := sunsetLST - longitude/15.0

	// 确保时间在0-24范围内
	sunriseUTC = c.normalizeHour(sunriseUTC)
	sunsetUTC = c.normalizeHour(sunsetUTC)
	solarNoonUTC = c.normalizeHour(solarNoonUTC)

	return sunriseUTC, sunsetUTC, solarNoonUTC, nil
}

// calculateEquationOfTime 计算方程时差（分钟）
// 方程时差 = 视太阳时 - 平太阳时
func (c *SunriseSunsetCalculatorFixed) calculateEquationOfTime(jd float64) float64 {
	julianCentury := (jd - 2451545.0) / 36525.0

	// 几何平黄经
	geomMeanLongSun := 280.46646 + julianCentury*(36000.76983+0.0003032*julianCentury)
	geomMeanLongSun = c.normalizeAngleFixed(geomMeanLongSun)

	// 几何平近点角
	geomMeanAnomalySun := 357.52911 + julianCentury*(35999.05029-0.0001537*julianCentury)
	geomMeanAnomalySun = c.normalizeAngleFixed(geomMeanAnomalySun)

	// 地球轨道离心率
	eccentEarthOrbit := 0.016708634 - julianCentury*(0.000042037+0.0000001267*julianCentury)

	// 黄赤交角修正
	meanObliqEcliptic := 23.0 + (26.0+
		((21.448-julianCentury*(46.815+julianCentury*(0.00059-julianCentury*0.001813))))/60.0)/60.0
	obliquityCorrection := meanObliqEcliptic +
		0.00256*math.Cos((125.04-1934.136*julianCentury)*math.Pi/180.0)

	// 计算方程时差（分钟）
	var y = math.Tan(obliquityCorrection*math.Pi/180.0/2.0) *
		math.Tan(obliquityCorrection*math.Pi/180.0/2.0)

	sin2l0 := math.Sin(2.0 * geomMeanLongSun * math.Pi / 180.0)
	sinm := math.Sin(geomMeanAnomalySun * math.Pi / 180.0)
	cos2l0 := math.Cos(2.0 * geomMeanLongSun * math.Pi / 180.0)
	sin4l0 := math.Sin(4.0 * geomMeanLongSun * math.Pi / 180.0)
	sin2m := math.Sin(2.0 * geomMeanAnomalySun * math.Pi / 180.0)

	equationOfTime := y*sin2l0 - 2.0*eccentEarthOrbit*sinm +
		4.0*eccentEarthOrbit*y*sinm*cos2l0 -
		0.5*y*y*sin4l0 - 1.25*eccentEarthOrbit*eccentEarthOrbit*sin2m

	return equationOfTime * 180.0 / math.Pi * 4.0 // 转换为分钟
}

// calculateDayLengthFixed 计算白昼长度（小时）
func (c *SunriseSunsetCalculatorFixed) calculateDayLengthFixed(sunrise, sunset float64) float64 {
	if sunset < sunrise {
		sunset += 24.0
	}
	return sunset - sunrise
}

// TwilightTimesFixed 晨昏蒙影时间
type TwilightTimesFixed struct {
	civilMorning        float64
	civilEvening        float64
	nauticalMorning     float64
	nauticalEvening     float64
	astronomicalMorning float64
	astronomicalEvening float64
}

// calculateTwilightTimesFixed 计算晨昏蒙影时间（修复版）
func (c *SunriseSunsetCalculatorFixed) calculateTwilightTimesFixed(
	jd, longitude, latitude, declination float64,
) *TwilightTimesFixed {
	times := &TwilightTimesFixed{}

	// 民用晨昏蒙影（太阳在地平线下6度）
	times.civilMorning, times.civilEvening = c.calculateTwilightFixed(
		jd, longitude, latitude, declination, -6.0,
	)

	// 航海晨昏蒙影（太阳在地平线下12度）
	times.nauticalMorning, times.nauticalEvening = c.calculateTwilightFixed(
		jd, longitude, latitude, declination, -12.0,
	)

	// 天文晨昏蒙影（太阳在地平线下18度）
	times.astronomicalMorning, times.astronomicalEvening = c.calculateTwilightFixed(
		jd, longitude, latitude, declination, -18.0,
	)

	return times
}

// calculateTwilightFixed 计算特定角度的晨昏蒙影（修复版）
func (c *SunriseSunsetCalculatorFixed) calculateTwilightFixed(
	jd, longitude, latitude, declination, angle float64,
) (morning, evening float64) {
	// 计算时角
	latRad := latitude * math.Pi / 180.0
	decRad := declination * math.Pi / 180.0
	zenithRad := (90.0 + angle) * math.Pi / 180.0

	cosH := (math.Cos(zenithRad) - math.Sin(latRad)*math.Sin(decRad)) /
		(math.Cos(latRad) * math.Cos(decRad))

	if cosH < -1.0 || cosH > 1.0 {
		return 0, 0 // 该纬度该日期没有晨昏蒙影
	}

	H := math.Acos(cosH) * 180.0 / math.Pi

	// 计算方程时差
	equationOfTime := c.calculateEquationOfTime(jd)
	solarNoonLST := 12.0 - equationOfTime/60.0

	// 计算晨昏蒙影的本地视太阳时
	morningLST := solarNoonLST - H/15.0
	eveningLST := solarNoonLST + H/15.0

	// 转换为UTC
	morningUTC := morningLST - longitude/15.0
	eveningUTC := eveningLST - longitude/15.0

	return c.normalizeHour(morningUTC), c.normalizeHour(eveningUTC)
}

// hoursToTimeString 将小时数转换为时间字符串（支持时区转换）
func (c *SunriseSunsetCalculatorFixed) hoursToTimeString(hoursUTC, timezone float64) string {
	// 将UTC小时转换为本地时间
	localHours := hoursUTC + timezone
	localHours = c.normalizeHour(localHours)

	hour := int(math.Floor(localHours))
	minute := int(math.Floor((localHours - float64(hour)) * 60.0))

	// 处理边界情况
	if minute >= 60 {
		minute -= 60
		hour += 1
	}
	if hour >= 24 {
		hour -= 24
	}

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// normalizeHour 将小时数归一化到0-24范围
func (c *SunriseSunsetCalculatorFixed) normalizeHour(hour float64) float64 {
	for hour < 0 {
		hour += 24.0
	}
	for hour >= 24.0 {
		hour -= 24.0
	}
	return hour
}

// normalizeAngleFixed 归一化角度到0-360度
func (c *SunriseSunsetCalculatorFixed) normalizeAngleFixed(angle float64) float64 {
	for angle < 0 {
		angle += 360.0
	}
	for angle >= 360.0 {
		angle -= 360.0
	}
	return angle
}
