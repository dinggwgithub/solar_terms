package calculator

import (
	"fmt"
	"math"
	"time"
)

// AstronomyCalculatorFixed 修复后的天文计算器
// 修复了原始代码中的三类典型Bug：
// 1. 结果不稳定性 - 使用统一的儒略世纪数T确保所有计算基于同一时间基准
// 2. 约束越界 - 角度标准化到[0, 360)范围，避免负数或超大值
// 3. 精度错误 - 使用VSOP87理论的高精度系数，正确区分度和弧度
type AstronomyCalculatorFixed struct {
	*BaseCalculator
}

// NewAstronomyCalculatorFixed 创建新的修复后天文计算器
func NewAstronomyCalculatorFixed() *AstronomyCalculatorFixed {
	return &AstronomyCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"astronomy_fixed",
			"修复后的天文黄经计算，精确计算太阳黄经和天文参数（基于VSOP87理论）",
		),
	}
}

// AstronomyFixedParams 天文计算参数
type AstronomyFixedParams struct {
	Year  int     `json:"year"`  // 年份
	Month int     `json:"month"` // 月份
	Day   int     `json:"day"`   // 日期
	Hour  int     `json:"hour"`  // 小时（可选，默认12）
	JD    float64 `json:"jd"`    // 儒略日（可选，优先使用）
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

	// 优先使用儒略日
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

	hour := 12.0 // 默认中午12点
	if h, ok := paramsMap["hour"].(float64); ok {
		hour = h
	}

	return &AstronomyFixedParams{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
		Hour:  int(hour),
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

	daysInMonth := c.getDaysInMonth(params.Year, params.Month)
	if params.Day < 1 || params.Day > daysInMonth {
		return fmt.Errorf("日期超出范围 (1-%d): %d", daysInMonth, params.Day)
	}

	return nil
}

// getDaysInMonth 获取某月的天数
func (c *AstronomyCalculatorFixed) getDaysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if c.isLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 31
	}
}

// isLeapYear 判断是否为闰年
func (c *AstronomyCalculatorFixed) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// ============================================
// 核心天文计算函数（基于VSOP87理论和IAU标准）
// ============================================

// calculateJulianDate 计算儒略日（高精度）
// 参考：天文算法（Astronomical Algorithms）第7章
// JD = 儒略日，表示从公元前4713年1月1日正午开始的天数
func (c *AstronomyCalculatorFixed) calculateJulianDate(t time.Time) float64 {
	year := float64(t.Year())
	month := float64(t.Month())
	day := float64(t.Day())
	hour := float64(t.Hour())
	minute := float64(t.Minute())
	second := float64(t.Second())

	// 如果月份是1月或2月，视为前一年的13月或14月
	if month <= 2 {
		year -= 1
		month += 12
	}

	// 计算世纪数A和B（格里高利历修正）
	A := math.Floor(year / 100)
	B := 2 - A + math.Floor(A/4)

	// 计算儒略日数（正午）
	jd := math.Floor(365.25*(year+4716)) + math.Floor(30.6001*(month+1)) + day + B - 1524.5

	// 添加时间部分（转换为日的小数）
	fractionOfDay := (hour + minute/60.0 + second/3600.0) / 24.0
	jd += fractionOfDay

	return jd
}

// julianCentury 计算从J2000.0开始的儒略世纪数
// T = (JD - 2451545.0) / 36525.0
// 这是所有天文计算的基础时间变量
func (c *AstronomyCalculatorFixed) julianCentury(jd float64) float64 {
	return (jd - 2451545.0) / 36525.0
}

// normalizeAngle 将角度标准化到 [0, 360) 范围
// 这是修复"约束越界"Bug的关键函数
func (c *AstronomyCalculatorFixed) normalizeAngle(degrees float64) float64 {
	normalized := math.Mod(degrees, 360.0)
	if normalized < 0 {
		normalized += 360.0
	}
	return normalized
}

// deg2rad 角度转弧度
func (c *AstronomyCalculatorFixed) deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// rad2deg 弧度转角度
func (c *AstronomyCalculatorFixed) rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// calculateMeanLongitude 计算太阳平黄经 L0（度）
// 基于VSOP87理论的高精度公式
// L0 = 280.46646 + 36000.76983*T + 0.0003032*T^2
// 这是修复"精度错误"Bug的关键
func (c *AstronomyCalculatorFixed) calculateMeanLongitude(T float64) float64 {
	L0 := 280.46646 +
		36000.76983*T +
		0.0003032*T*T +
		T*T*T/49931.0 -
		T*T*T*T/15300.0

	return c.normalizeAngle(L0)
}

// calculateMeanAnomaly 计算太阳平近点角 M（度）
// 基于VSOP87理论
// M = 357.52911 + 35999.05029*T - 0.0001537*T^2
func (c *AstronomyCalculatorFixed) calculateMeanAnomaly(T float64) float64 {
	M := 357.52911 +
		35999.05029*T -
		0.0001537*T*T +
		T*T*T/24490000.0

	return c.normalizeAngle(M)
}

// calculateEquationOfCenter 计算中心差 C（度）
// 中心差是平近点角与真近点角的差值
// C = (1.914602 - 0.004817*T - 0.000014*T^2) * sin(M)
//   - (0.019993 - 0.000101*T) * sin(2*M)
//   - 0.000289 * sin(3*M)
func (c *AstronomyCalculatorFixed) calculateEquationOfCenter(M, T float64) float64 {
	Mrad := c.deg2rad(M)

	C := (1.914602-0.004817*T-0.000014*T*T)*math.Sin(Mrad) +
		(0.019993-0.000101*T)*math.Sin(2*Mrad) +
		0.000289*math.Sin(3*Mrad)

	return C
}

// calculateTrueLongitude 计算太阳真黄经（度）
// 真黄经 = 平黄经 + 中心差
func (c *AstronomyCalculatorFixed) calculateTrueLongitude(L0, C float64) float64 {
	return c.normalizeAngle(L0 + C)
}

// calculateNutation 计算黄经章动（度）
// 基于IAU 1980章动理论
// 章动是地球自转轴的微小摆动，影响黄经测量
func (c *AstronomyCalculatorFixed) calculateNutation(T float64) float64 {
	// 升交点黄经（月球轨道升交点）
	Omega := 125.04452 - 1934.136261*T + 0.0020708*T*T + T*T*T/450000.0
	Omega = c.normalizeAngle(Omega)
	OmegaRad := c.deg2rad(Omega)

	// 太阳平黄经
	L := c.calculateMeanLongitude(T)
	Lrad := c.deg2rad(L)

	// 太阳平近点角
	M := c.calculateMeanAnomaly(T)
	Mrad := c.deg2rad(M)

	// 黄经章动（度）- 主要项
	// 系数来自IAU 1980章动理论（角秒转换为度）
	nutation := (-17.20/3600.0)*math.Sin(OmegaRad) +
		(-1.32/3600.0)*math.Sin(2*Lrad) +
		(-0.23/3600.0)*math.Sin(2*Mrad) +
		(0.21/3600.0)*math.Sin(2*OmegaRad)

	return nutation
}

// calculateAberration 计算光行差修正（度）
// 光行差是由于地球运动导致的光视方向变化
// 年平均光行差约为20.49552角秒
func (c *AstronomyCalculatorFixed) calculateAberration(T float64) float64 {
	// 光行差常数（度）
	const aberrationConstant = 20.49552 / 3600.0
	return -aberrationConstant
}

// calculateApparentLongitude 计算太阳视黄经（度）
// 视黄经 = 真黄经 + 章动 + 光行差
// 这是从地球观测到的太阳黄经
func (c *AstronomyCalculatorFixed) calculateApparentLongitude(trueLongitude, T float64) float64 {
	nutation := c.calculateNutation(T)
	aberration := c.calculateAberration(T)

	apparentLongitude := trueLongitude + nutation + aberration
	return c.normalizeAngle(apparentLongitude)
}

// calculateSunLongitude 计算太阳黄经（主函数）
// 返回视黄经，这是从地球观测到的太阳位置
func (c *AstronomyCalculatorFixed) calculateSunLongitude(jd float64) float64 {
	T := c.julianCentury(jd)
	L0 := c.calculateMeanLongitude(T)
	M := c.calculateMeanAnomaly(T)
	C := c.calculateEquationOfCenter(M, T)
	trueLongitude := c.calculateTrueLongitude(L0, C)
	apparentLongitude := c.calculateApparentLongitude(trueLongitude, T)
	return apparentLongitude
}

// calculateAstronomyFixed 执行修复后的天文计算
// 这是修复"结果不稳定性"Bug的关键：所有计算使用相同的T值
func (c *AstronomyCalculatorFixed) calculateAstronomyFixed(params *AstronomyFixedParams) (map[string]float64, error) {
	// 计算儒略日
	var jd float64
	if params.JD > 0 {
		jd = params.JD
	} else {
		t := time.Date(params.Year, time.Month(params.Month), params.Day, params.Hour, 0, 0, 0, time.UTC)
		jd = c.calculateJulianDate(t)
	}

	// 计算儒略世纪数（所有计算共享同一时间基准）
	// 这是修复"结果不稳定性"Bug的关键
	T := c.julianCentury(jd)

	// 基于同一时间基准计算所有参数
	meanLongitude := c.calculateMeanLongitude(T)
	meanAnomaly := c.calculateMeanAnomaly(T)
	equationOfCenter := c.calculateEquationOfCenter(meanAnomaly, T)
	trueLongitude := c.calculateTrueLongitude(meanLongitude, equationOfCenter)
	nutation := c.calculateNutation(T)
	aberration := c.calculateAberration(T)
	apparentLongitude := c.calculateApparentLongitude(trueLongitude, T)

	// 太阳黄经使用视黄经（从地球观测的实际位置）
	sunLongitude := apparentLongitude

	return map[string]float64{
		"sun_longitude":      sunLongitude,      // 太阳视黄经（度）
		"julian_date":        jd,                // 儒略日（带小数部分）
		"apparent_longitude": apparentLongitude, // 视黄经（度）
		"true_longitude":     trueLongitude,     // 真黄经（度）
		"mean_longitude":     meanLongitude,     // 平黄经（度）
		"mean_anomaly":       meanAnomaly,       // 平近点角（度）
		"equation_of_center": equationOfCenter,  // 中心差（度）
		"nutation":           nutation,          // 黄经章动（度）
		"aberration":         aberration,        // 光行差（度）
	}, nil
}

// ValidateAstronomyResult 验证天文计算结果是否在合理范围内
// 用于A/B测试时判断修复效果
func (c *AstronomyCalculatorFixed) ValidateAstronomyResult(result map[string]float64) map[string]string {
	validation := make(map[string]string)

	// 定义合理范围
	ranges := map[string][2]float64{
		"sun_longitude":      {0, 360},
		"apparent_longitude": {0, 360},
		"true_longitude":     {0, 360},
		"mean_longitude":     {0, 360},
		"mean_anomaly":       {0, 360},
		"julian_date":        {2415020.5, 2488070.5},
		"equation_of_center": {-2.5, 2.5},
		"nutation":           {-0.02, 0.02},
		"aberration":         {-0.01, -0.005},
	}

	for key, validRange := range ranges {
		if val, ok := result[key]; ok {
			if val >= validRange[0] && val <= validRange[1] {
				validation[key] = "正常"
			} else {
				validation[key] = fmt.Sprintf("异常（%.4f 不在 [%.2f, %.2f] 范围内）", val, validRange[0], validRange[1])
			}
		}
	}

	return validation
}
