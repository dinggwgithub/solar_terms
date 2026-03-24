package calculator

import (
	"fmt"
	"math"
	"scientific_calc_bugs/internal/bugs"
	"time"
)

// AstronomyCalculator 天文计算器
type AstronomyCalculator struct {
	*BaseCalculator
}

// NewAstronomyCalculator 创建新的天文计算器
func NewAstronomyCalculator() *AstronomyCalculator {
	return &AstronomyCalculator{
		BaseCalculator: NewBaseCalculator(
			"astronomy",
			"天文黄经计算，计算太阳黄经和天文参数",
		),
	}
}

// AstronomyParams 天文计算参数
type AstronomyParams struct {
	Year  int     `json:"year"`  // 年份
	Month int     `json:"month"` // 月份
	Day   int     `json:"day"`   // 日期
	JD    float64 `json:"jd"`    // 儒略日
}

// Calculate 执行天文计算
func (c *AstronomyCalculator) Calculate(params interface{}) (interface{}, error) {
	astroParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数有效性
	if err := c.validateParams(astroParams); err != nil {
		return nil, err
	}

	// 执行天文计算
	result, err := c.calculateAstronomy(astroParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *AstronomyCalculator) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// parseParams 解析参数
func (c *AstronomyCalculator) parseParams(params interface{}) (*AstronomyParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 优先使用儒略日
	if jd, ok := paramsMap["jd"].(float64); ok && jd > 0 {
		return &AstronomyParams{JD: jd}, nil
	}

	// 使用年月日计算儒略日
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

	return &AstronomyParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
	}, nil
}

// validateParams 验证参数有效性
func (c *AstronomyCalculator) validateParams(params *AstronomyParams) error {
	// 如果提供了儒略日，验证其范围
	if params.JD > 0 {
		if params.JD < 2415020.5 || params.JD > 2488070.5 {
			return fmt.Errorf("儒略日超出支持范围 (2415020.5-2488070.5): %f", params.JD)
		}
		return nil
	}

	// 验证年月日范围
	if params.Year < 1900 || params.Year > 2100 {
		return fmt.Errorf("年份超出支持范围 (1900-2100): %d", params.Year)
	}

	if params.Month < 1 || params.Month > 12 {
		return fmt.Errorf("月份超出范围 (1-12): %d", params.Month)
	}

	if params.Day < 1 || params.Day > 31 {
		return fmt.Errorf("日期超出范围 (1-31): %d", params.Day)
	}

	return nil
}

// calculateAstronomy 执行天文计算
func (c *AstronomyCalculator) calculateAstronomy(params *AstronomyParams) (map[string]float64, error) {
	// 计算儒略日
	jd := params.JD
	if jd == 0 {
		t := time.Date(params.Year, time.Month(params.Month), params.Day, 12, 0, 0, 0, time.UTC)
		jd = c.calculateJulianDate(t)
	}

	// 计算太阳黄经
	sunLongitude := c.calculateSunLongitude(jd)

	// 计算视黄经
	apparentLongitude := c.calculateApparentLongitude(sunLongitude, jd)

	// 计算真黄经
	trueLongitude := c.calculateTrueLongitude(apparentLongitude)

	// 计算平黄经
	meanLongitude := c.calculateMeanLongitude(jd)

	// 计算平近点角
	meanAnomaly := c.calculateMeanAnomaly(jd)

	// 计算中心差
	equationOfCenter := c.calculateEquationOfCenter(meanAnomaly)

	// 计算章动
	nutation := c.calculateNutation(jd)

	return map[string]float64{
		"sun_longitude":      sunLongitude,
		"julian_date":        jd,
		"apparent_longitude": apparentLongitude,
		"true_longitude":     trueLongitude,
		"mean_longitude":     meanLongitude,
		"mean_anomaly":       meanAnomaly,
		"equation_of_center": equationOfCenter,
		"nutation":           nutation,
	}, nil
}

// calculateJulianDate 计算儒略日
func (c *AstronomyCalculator) calculateJulianDate(t time.Time) float64 {
	year := float64(t.Year())
	month := float64(t.Month())
	day := float64(t.Day())
	hour := float64(t.Hour())
	minute := float64(t.Minute())
	second := float64(t.Second())

	if month <= 2 {
		year -= 1
		month += 12
	}

	A := math.Floor(year / 100)
	B := 2 - A + math.Floor(A/4)

	jd := math.Floor(365.25*(year+4716)) + math.Floor(30.6001*(month+1)) + day + B - 1524.5
	jd += (hour + minute/60 + second/3600) / 24

	return jd
}

// calculateSunLongitude 计算太阳黄经
func (c *AstronomyCalculator) calculateSunLongitude(jd float64) float64 {
	// 简化的太阳黄经计算
	// 实际应该基于精确的天文算法

	// 计算从J2000.0开始的天数
	daysSinceJ2000 := jd - 2451545.0

	// 平均黄经（度/天）
	meanLongitude := 280.460 + 0.9856474*daysSinceJ2000
	meanLongitude = math.Mod(meanLongitude, 360)
	if meanLongitude < 0 {
		meanLongitude += 360
	}

	// 平近点角
	meanAnomaly := 357.528 + 0.9856003*daysSinceJ2000
	meanAnomaly = math.Mod(meanAnomaly, 360)
	if meanAnomaly < 0 {
		meanAnomaly += 360
	}

	// 中心差（度）
	equationOfCenter := 1.915*math.Sin(meanAnomaly*math.Pi/180) + 0.020*math.Sin(2*meanAnomaly*math.Pi/180)

	// 真黄经
	trueLongitude := meanLongitude + equationOfCenter
	trueLongitude = math.Mod(trueLongitude, 360)
	if trueLongitude < 0 {
		trueLongitude += 360
	}

	return trueLongitude
}

// calculateApparentLongitude 计算视黄经
func (c *AstronomyCalculator) calculateApparentLongitude(sunLongitude float64, jd float64) float64 {
	// 简化的视黄经计算
	// 实际应该考虑章动和光行差

	// 章动修正（简化）
	nutation := 0.004 * math.Sin((125.04-0.052954*jd)*math.Pi/180)

	// 光行差修正（简化）
	aberration := 0.0057 * math.Sin(sunLongitude*math.Pi/180)

	apparentLongitude := sunLongitude + nutation + aberration
	apparentLongitude = math.Mod(apparentLongitude, 360)
	if apparentLongitude < 0 {
		apparentLongitude += 360
	}

	return apparentLongitude
}

// calculateTrueLongitude 计算真黄经
func (c *AstronomyCalculator) calculateTrueLongitude(apparentLongitude float64) float64 {
	// 真黄经与视黄经基本相同（简化）
	return apparentLongitude
}

// calculateMeanLongitude 计算平黄经
func (c *AstronomyCalculator) calculateMeanLongitude(jd float64) float64 {
	// 简化的平黄经计算
	daysSinceJ2000 := jd - 2451545.0
	meanLongitude := 280.460 + 0.9856474*daysSinceJ2000
	meanLongitude = math.Mod(meanLongitude, 360)
	if meanLongitude < 0 {
		meanLongitude += 360
	}
	return meanLongitude
}

// calculateMeanAnomaly 计算平近点角
func (c *AstronomyCalculator) calculateMeanAnomaly(jd float64) float64 {
	// 简化的平近点角计算
	daysSinceJ2000 := jd - 2451545.0
	meanAnomaly := 357.528 + 0.9856003*daysSinceJ2000
	meanAnomaly = math.Mod(meanAnomaly, 360)
	if meanAnomaly < 0 {
		meanAnomaly += 360
	}
	return meanAnomaly
}

// calculateEquationOfCenter 计算中心差
func (c *AstronomyCalculator) calculateEquationOfCenter(meanAnomaly float64) float64 {
	// 简化的中心差计算
	meanAnomalyRad := meanAnomaly * math.Pi / 180
	equationOfCenter := 1.915*math.Sin(meanAnomalyRad) + 0.020*math.Sin(2*meanAnomalyRad)
	return equationOfCenter
}

// calculateNutation 计算章动
func (c *AstronomyCalculator) calculateNutation(jd float64) float64 {
	// 简化的章动计算
	nutation := 0.004 * math.Sin((125.04-0.052954*jd)*math.Pi/180)
	return nutation
}

// GetSupportedBugTypes 返回支持的Bug类型
func (c *AstronomyCalculator) GetSupportedBugTypes() []bugs.BugType {
	return []bugs.BugType{
		bugs.BugTypeInstability,
		bugs.BugTypeConstraint,
		bugs.BugTypePrecision,
	}
}
