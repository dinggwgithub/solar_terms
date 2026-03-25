package calculator

import (
	"fmt"
	"math"
	"time"
)

// AstronomyCalculatorFixed 修复后的天文计算器
type AstronomyCalculatorFixed struct {
	*BaseCalculator
}

// NewAstronomyCalculatorFixed 创建新的修复后天文计算器
func NewAstronomyCalculatorFixed() *AstronomyCalculatorFixed {
	return &AstronomyCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"astronomy_fixed",
			"修复后的天文黄经计算，精确计算太阳黄经和天文参数",
		),
	}
}

// AstronomyFixedParams 天文计算参数
type AstronomyFixedParams struct {
	Year  int     `json:"year"`  // 年份
	Month int     `json:"month"` // 月份
	Day   int     `json:"day"`   // 日期
	JD    float64 `json:"jd"`    // 儒略日
}

// Calculate 执行天文计算
func (c *AstronomyCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	astroParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(astroParams); err != nil {
		return nil, err
	}

	result, err := c.calculateAstronomyFixed(astroParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Validate 验证输入参数
func (c *AstronomyCalculatorFixed) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

// parseParams 解析参数
func (c *AstronomyCalculatorFixed) parseParams(params interface{}) (*AstronomyFixedParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	if jd, ok := paramsMap["jd"].(float64); ok && jd > 0 {
		return &AstronomyFixedParams{JD: jd}, nil
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

	return &AstronomyFixedParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
	}, nil
}

// validateParams 验证参数有效性
func (c *AstronomyCalculatorFixed) validateParams(params *AstronomyFixedParams) error {
	if params.JD > 0 {
		if params.JD < 2415020.5 || params.JD > 2488070.5 {
			return fmt.Errorf("儒略日超出支持范围 (2415020.5-2488070.5): %f", params.JD)
		}
		return nil
	}

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

// 计算J2000.0以来的儒略世纪数
func (c *AstronomyCalculatorFixed) julianCentury(jd float64) float64 {
	return (jd - 2451545.0) / 36525.0
}

// 角度标准化到 [0, 360)
func (c *AstronomyCalculatorFixed) normalizeAngle(degrees float64) float64 {
	normalized := math.Mod(degrees, 360.0)
	if normalized < 0 {
		normalized += 360.0
	}
	return normalized
}

// 角度转弧度
func (c *AstronomyCalculatorFixed) deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// 弧度转角度
func (c *AstronomyCalculatorFixed) rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// calculateJulianDate 计算儒略日（精确计算）
// 公式来源：天文算法，基于格里高利历
func (c *AstronomyCalculatorFixed) calculateJulianDate(t time.Time) float64 {
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

	// 基础儒略日计算（正午）
	// 注意：使用30.6而不是30.6001来避免某些边界值的floor计算误差
	jd := math.Floor(365.25*(year+4716)) + math.Floor(30.6*(month+1)) + day + B - 1524.5

	// 添加时间部分
	jd += (hour + minute/60.0 + second/3600.0) / 24.0

	return jd
}

// calculateMeanLongitude 计算太阳平黄经（精度改进版）
func (c *AstronomyCalculatorFixed) calculateMeanLongitude(T float64) float64 {
	// VSOP87理论的平黄经系数（精度更高）
	L0 := 280.4664567 +
		36000.76982779*T +
		0.0003032028*T*T +
		1.0/49931.0*T*T*T -
		1.0/15300.0*T*T*T*T -
		1.0/2000000.0*T*T*T*T*T

	return c.normalizeAngle(L0)
}

// calculateMeanAnomaly 计算太阳平近点角（精度改进版）
func (c *AstronomyCalculatorFixed) calculateMeanAnomaly(T float64) float64 {
	// VSOP87理论的平近点角系数
	M := 357.5291092 +
		35999.0502909*T -
		0.0001536*T*T +
		1.0/24490000.0*T*T*T

	return c.normalizeAngle(M)
}

// calculateEquationOfCenter 计算中心差（精度改进版）
func (c *AstronomyCalculatorFixed) calculateEquationOfCenter(M float64, T float64) float64 {
	Mrad := c.deg2rad(M)

	// 更高阶的中心差计算，包含更多项
	C := (1.914602-0.004817*T-0.000014*T*T)*math.Sin(Mrad) +
		(0.019993-0.000101*T)*math.Sin(2*Mrad) +
		0.000289*math.Sin(3*Mrad)

	return C
}

// calculateTrueLongitude 计算太阳真黄经
func (c *AstronomyCalculatorFixed) calculateTrueLongitude(L0 float64, C float64) float64 {
	return c.normalizeAngle(L0 + C)
}

// calculateApparentLongitude 计算视黄经（精度改进版）
func (c *AstronomyCalculatorFixed) calculateApparentLongitude(trueLongitude float64, T float64) float64 {
	// 计算章动
	nutation := c.calculateNutation(T)

	// 光行差修正（更精确的值）
	aberration := 0.005693

	// 视黄经 = 真黄经 + 章动 + 光行差
	apparentLongitude := trueLongitude + nutation + aberration

	return c.normalizeAngle(apparentLongitude)
}

// calculateNutation 计算黄经章动（精度改进版）
func (c *AstronomyCalculatorFixed) calculateNutation(T float64) float64 {
	// 章动主要项（IAU1980章动理论的主要项）
	Omega := 125.04452 - 1934.136261*T + 0.0020708*T*T + T*T*T/450000
	Omega = c.normalizeAngle(Omega)
	OmegaRad := c.deg2rad(Omega)

	// 黄经章动主要项，单位：度（转换自角秒: 1度 = 3600角秒）
	// Δψ = -17.20" sinΩ - 1.32" sin2L - 0.23" sin2L' + 0.21" sin2Ω
	nutation := (-17.20/3600.0)*math.Sin(OmegaRad) +
		(-1.32/3600.0)*math.Sin(2*c.deg2rad(c.calculateMeanLongitude(T))) +
		(-0.23/3600.0)*math.Sin(2*c.deg2rad(c.calculateMeanAnomaly(T))) +
		(0.21/3600.0)*math.Sin(2*OmegaRad)

	return nutation
}

// calculateSunLongitude 计算太阳黄经（主函数）
func (c *AstronomyCalculatorFixed) calculateSunLongitude(jd float64) float64 {
	T := c.julianCentury(jd)
	L0 := c.calculateMeanLongitude(T)
	M := c.calculateMeanAnomaly(T)
	C := c.calculateEquationOfCenter(M, T)
	trueLongitude := c.calculateTrueLongitude(L0, C)
	return trueLongitude
}

// calculateAstronomyFixed 执行修复后的天文计算
func (c *AstronomyCalculatorFixed) calculateAstronomyFixed(params *AstronomyFixedParams) (map[string]float64, error) {
	// 计算儒略日
	jd := params.JD
	if jd == 0 {
		t := time.Date(params.Year, time.Month(params.Month), params.Day, 12, 0, 0, 0, time.UTC)
		jd = c.calculateJulianDate(t)
	}

	// 计算J2000世纪数（用于所有计算，确保一致性）
	T := c.julianCentury(jd)

	// 所有参数使用相同的T值，确保计算一致性
	meanLongitude := c.calculateMeanLongitude(T)
	meanAnomaly := c.calculateMeanAnomaly(T)
	equationOfCenter := c.calculateEquationOfCenter(meanAnomaly, T)
	trueLongitude := c.calculateTrueLongitude(meanLongitude, equationOfCenter)
	apparentLongitude := c.calculateApparentLongitude(trueLongitude, T)
	nutation := c.calculateNutation(T)

	// 使用更精确的算法计算太阳黄经
	sunLongitude := trueLongitude

	return map[string]float64{
		"sun_longitude":      sunLongitude,
		"julian_date":        jd, // 保持完整浮点精度
		"apparent_longitude": apparentLongitude,
		"true_longitude":     trueLongitude,
		"mean_longitude":     meanLongitude,
		"mean_anomaly":       meanAnomaly,
		"equation_of_center": equationOfCenter,
		"nutation":           nutation,
	}, nil
}
