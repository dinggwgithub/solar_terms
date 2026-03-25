package calculator

import (
	"fmt"
	"math"
)

type SunriseSunsetCalculatorFixed struct {
	*BaseCalculator
}

func NewSunriseSunsetCalculatorFixed() *SunriseSunsetCalculatorFixed {
	return &SunriseSunsetCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"sunrise_sunset_fixed",
			"日出日落时间计算器（修复版），基于地理位置和日期计算日出日落时刻",
		),
	}
}

type SunriseSunsetParamsFixed struct {
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	Day       int     `json:"day"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
	Timezone  float64 `json:"timezone"`
}

type SunriseSunsetResultFixed struct {
	Date          string  `json:"date"`
	Sunrise       string  `json:"sunrise"`
	Sunset        string  `json:"sunset"`
	SolarNoon     string  `json:"solar_noon"`
	DayLength     float64 `json:"day_length"`
	CivilTwilight struct {
		Morning string `json:"morning"`
		Evening string `json:"evening"`
	} `json:"civil_twilight"`
	NauticalTwilight struct {
		Morning string `json:"morning"`
		Evening string `json:"evening"`
	} `json:"nautical_twilight"`
	AstronomicalTwilight struct {
		Morning string `json:"morning"`
		Evening string `json:"evening"`
	} `json:"astronomical_twilight"`
}

func (c *SunriseSunsetCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(sunriseParams); err != nil {
		return nil, err
	}

	result, err := c.calculateSunriseSunset(sunriseParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *SunriseSunsetCalculatorFixed) Validate(params interface{}) error {
	sunriseParams, err := c.parseParams(params)
	if err != nil {
		return err
	}
	return c.validateParams(sunriseParams)
}

func (c *SunriseSunsetCalculatorFixed) parseParams(params interface{}) (*SunriseSunsetParamsFixed, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

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

	altitude := 0.0
	if alt, exists := paramsMap["altitude"]; exists {
		if altFloat, ok := alt.(float64); ok {
			altitude = altFloat
		}
	}

	timezone := 8.0
	if tz, exists := paramsMap["timezone"]; exists {
		if tzFloat, ok := tz.(float64); ok {
			timezone = tzFloat
		}
	}

	return &SunriseSunsetParamsFixed{
		Year:      int(year),
		Month:     int(month),
		Day:       int(day),
		Longitude: longitude,
		Latitude:  latitude,
		Altitude:  altitude,
		Timezone:  timezone,
	}, nil
}

func (c *SunriseSunsetCalculatorFixed) validateParams(params *SunriseSunsetParamsFixed) error {
	if err := c.validateDate(params.Year, params.Month, params.Day); err != nil {
		return err
	}

	if params.Longitude < -180 || params.Longitude > 180 {
		return fmt.Errorf("经度超出范围 (-180到180): %f", params.Longitude)
	}

	if params.Latitude < -90 || params.Latitude > 90 {
		return fmt.Errorf("纬度超出范围 (-90到90): %f", params.Latitude)
	}

	if params.Altitude < -1000 || params.Altitude > 10000 {
		return fmt.Errorf("海拔超出合理范围 (-1000到10000米): %f", params.Altitude)
	}

	if params.Timezone < -12 || params.Timezone > 14 {
		return fmt.Errorf("时区超出范围 (-12到14): %f", params.Timezone)
	}

	return nil
}

func (c *SunriseSunsetCalculatorFixed) validateDate(year, month, day int) error {
	if year < 1900 || year > 2100 {
		return fmt.Errorf("年份超出支持范围 (1900-2100): %d", year)
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("月份超出范围 (1-12): %d", month)
	}

	if day < 1 || day > 31 {
		return fmt.Errorf("日期超出范围 (1-31): %d", day)
	}

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

func (c *SunriseSunsetCalculatorFixed) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func (c *SunriseSunsetCalculatorFixed) calculateSunriseSunset(params *SunriseSunsetParamsFixed) (*SunriseSunsetResultFixed, error) {
	jd := c.calculateJulianDay(params.Year, params.Month, params.Day)

	solarPosition := c.calculateSolarPosition(jd)

	sunrise, sunset, err := c.calculateSunTimes(jd, params.Longitude, params.Latitude, params.Altitude)
	if err != nil {
		return nil, err
	}

	solarNoon := c.calculateSolarNoon(jd, params.Longitude)

	dayLength := c.calculateDayLength(sunrise, sunset)

	twilightTimes := c.calculateTwilightTimes(jd, params.Longitude, params.Latitude)

	dateStr := fmt.Sprintf("%d-%02d-%02d", params.Year, params.Month, params.Day)
	sunriseStr := c.formatTime(sunrise, params.Timezone)
	sunsetStr := c.formatTime(sunset, params.Timezone)
	solarNoonStr := c.formatTime(solarNoon, params.Timezone)

	result := &SunriseSunsetResultFixed{
		Date:      dateStr,
		Sunrise:   sunriseStr,
		Sunset:    sunsetStr,
		SolarNoon: solarNoonStr,
		DayLength: dayLength,
	}

	result.CivilTwilight.Morning = c.formatTime(twilightTimes.civilMorning, params.Timezone)
	result.CivilTwilight.Evening = c.formatTime(twilightTimes.civilEvening, params.Timezone)
	result.NauticalTwilight.Morning = c.formatTime(twilightTimes.nauticalMorning, params.Timezone)
	result.NauticalTwilight.Evening = c.formatTime(twilightTimes.nauticalEvening, params.Timezone)
	result.AstronomicalTwilight.Morning = c.formatTime(twilightTimes.astronomicalMorning, params.Timezone)
	result.AstronomicalTwilight.Evening = c.formatTime(twilightTimes.astronomicalEvening, params.Timezone)

	_ = solarPosition
	return result, nil
}

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

type SolarPositionFixed struct {
	RightAscension float64
	Declination    float64
	Distance       float64
}

func (c *SunriseSunsetCalculatorFixed) calculateSolarPosition(jd float64) *SolarPositionFixed {
	t := (jd - 2451545.0) / 36525.0

	L := 280.46646 + 36000.76983*t + 0.0003032*t*t
	L = c.normalizeAngle(L)

	M := 357.52911 + 35999.05029*t - 0.0001537*t*t
	M = c.normalizeAngle(M)

	C := (1.914602-0.004817*t-0.000014*t*t)*math.Sin(M*math.Pi/180) +
		(0.019993-0.000101*t)*math.Sin(2*M*math.Pi/180) +
		0.000289*math.Sin(3*M*math.Pi/180)

	lambda := L + C

	epsilon := 23.4392911 - 0.0130042*t - 0.00000016*t*t + 0.000000504*t*t*t

	alpha := math.Atan2(
		math.Sin(lambda*math.Pi/180)*math.Cos(epsilon*math.Pi/180),
		math.Cos(lambda*math.Pi/180),
	) * 180 / math.Pi

	delta := math.Asin(
		math.Sin(lambda*math.Pi/180)*math.Sin(epsilon*math.Pi/180),
	) * 180 / math.Pi

	return &SolarPositionFixed{
		RightAscension: alpha,
		Declination:    delta,
		Distance:       1.0,
	}
}

func (c *SunriseSunsetCalculatorFixed) calculateSunTimes(jd, longitude, latitude, altitude float64) (float64, float64, error) {
	h0 := -0.8333

	if altitude > 0 {
		h0 -= 0.0347 * math.Sqrt(altitude)
	}

	solarPos := c.calculateSolarPosition(jd)
	declination := solarPos.Declination

	cosH := (math.Sin(h0*math.Pi/180) -
		math.Sin(latitude*math.Pi/180)*math.Sin(declination*math.Pi/180)) /
		(math.Cos(latitude*math.Pi/180) * math.Cos(declination*math.Pi/180))

	if cosH < -1 || cosH > 1 {
		return 0, 0, fmt.Errorf("在该位置该日期没有日出或日落")
	}

	H := math.Acos(cosH) * 180 / math.Pi

	sunriseJD := jd - 0.5 + (720-4*(longitude+H))/1440.0
	sunsetJD := jd - 0.5 + (720-4*(longitude-H))/1440.0

	return sunriseJD, sunsetJD, nil
}

func (c *SunriseSunsetCalculatorFixed) calculateSolarNoon(jd, longitude float64) float64 {
	return jd - longitude/360.0
}

func (c *SunriseSunsetCalculatorFixed) calculateDayLength(sunriseJD, sunsetJD float64) float64 {
	return (sunsetJD - sunriseJD) * 24
}

type TwilightTimesFixed struct {
	civilMorning        float64
	civilEvening        float64
	nauticalMorning     float64
	nauticalEvening     float64
	astronomicalMorning float64
	astronomicalEvening float64
}

func (c *SunriseSunsetCalculatorFixed) calculateTwilightTimes(jd, longitude, latitude float64) *TwilightTimesFixed {
	times := &TwilightTimesFixed{}

	times.civilMorning, times.civilEvening = c.calculateTwilight(jd, longitude, latitude, -6)
	times.nauticalMorning, times.nauticalEvening = c.calculateTwilight(jd, longitude, latitude, -12)
	times.astronomicalMorning, times.astronomicalEvening = c.calculateTwilight(jd, longitude, latitude, -18)

	return times
}

func (c *SunriseSunsetCalculatorFixed) calculateTwilight(jd, longitude, latitude, angle float64) (float64, float64) {
	solarPos := c.calculateSolarPosition(jd)
	declination := solarPos.Declination

	cosH := (math.Sin(angle*math.Pi/180) -
		math.Sin(latitude*math.Pi/180)*math.Sin(declination*math.Pi/180)) /
		(math.Cos(latitude*math.Pi/180) * math.Cos(declination*math.Pi/180))

	if cosH < -1 || cosH > 1 {
		return 0, 0
	}

	H := math.Acos(cosH) * 180 / math.Pi

	morningJD := jd - 0.5 + (720-4*(longitude+H))/1440.0
	eveningJD := jd - 0.5 + (720-4*(longitude-H))/1440.0

	return morningJD, eveningJD
}

func (c *SunriseSunsetCalculatorFixed) formatTime(jd, timezone float64) string {
	fraction := jd - math.Floor(jd)
	hoursUTC := fraction * 24

	hours := hoursUTC + timezone

	for hours < 0 {
		hours += 24
	}
	for hours >= 24 {
		hours -= 24
	}

	hour := int(math.Floor(hours))
	minute := int(math.Floor((hours - float64(hour)) * 60))

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func (c *SunriseSunsetCalculatorFixed) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}
