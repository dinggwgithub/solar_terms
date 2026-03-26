package calculator

import (
	"fmt"
	"math"
)

// PlanetCalculatorFixed 修复版行星位置计算器
// 修复了原计算器中的以下问题：
// 1. 赤经输出范围归一化到 0-24h
// 2. 使用更精确的开普勒轨道计算
// 3. 正确计算日心坐标到地心坐标的转换
// 4. 正确计算距离、相位和距角
type PlanetCalculatorFixed struct {
	*BaseCalculator
}

// NewPlanetCalculatorFixed 创建新的修复版行星位置计算器
func NewPlanetCalculatorFixed() *PlanetCalculatorFixed {
	return &PlanetCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"planet_fixed",
			"行星位置计算器（修复版），修正了赤经归一化和轨道计算",
		),
	}
}

// PlanetPositionFixed 行星位置结果（与原接口保持一致）
type PlanetPositionFixed struct {
	RightAscension float64 `json:"right_ascension"` // 赤经（小时，0-24）
	Declination    float64 `json:"declination"`     // 赤纬（度）
	Distance       float64 `json:"distance"`        // 距离（天文单位）
	Magnitude      float64 `json:"magnitude"`       // 星等
	Phase          float64 `json:"phase"`           // 相位（0-1）
	Elongation     float64 `json:"elongation"`      // 距角（度）
}

// Calculate 执行行星位置计算
func (c *PlanetCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	planetParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(planetParams); err != nil {
		return nil, err
	}

	// 执行行星位置计算
	position, err := c.calculatePlanetPosition(planetParams)
	if err != nil {
		return nil, err
	}

	return position, nil
}

// Validate 验证输入参数
func (c *PlanetCalculatorFixed) Validate(params interface{}) error {
	planetParams, err := c.parseParams(params)
	if err != nil {
		return err
	}

	return c.validateParams(planetParams)
}

// parseParams 解析参数（与原版兼容）
func (c *PlanetCalculatorFixed) parseParams(params interface{}) (*PlanetParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 提取行星名称（支持 planet 和 planet_name 两种字段名）
	planet := ""
	if p, ok := paramsMap["planet"].(string); ok {
		planet = p
	} else if p, ok := paramsMap["planet_name"].(string); ok {
		planet = p
	} else {
		return nil, fmt.Errorf("planet或planet_name参数必须为字符串")
	}

	// 提取儒略日（支持直接传入jd或通过year/month/day计算）
	var jd float64
	if jdValue, ok := paramsMap["jd"].(float64); ok {
		jd = jdValue
	} else {
		// 尝试从year/month/day计算儒略日
		year, yOk := paramsMap["year"].(float64)
		month, mOk := paramsMap["month"].(float64)
		day, dOk := paramsMap["day"].(float64)

		if yOk && mOk && dOk {
			jd = c.dateToJulianDay(int(year), int(month), int(day))
		} else {
			return nil, fmt.Errorf("需要提供jd参数或year/month/day参数")
		}
	}

	// 提取经度（可选）
	longitude := 0.0
	if lon, exists := paramsMap["longitude"]; exists {
		if lonFloat, ok := lon.(float64); ok {
			longitude = lonFloat
		}
	}

	// 提取纬度（可选）
	latitude := 0.0
	if lat, exists := paramsMap["latitude"]; exists {
		if latFloat, ok := lat.(float64); ok {
			latitude = latFloat
		}
	}

	return &PlanetParams{
		Planet:    planet,
		JD:        jd,
		Longitude: longitude,
		Latitude:  latitude,
	}, nil
}

// dateToJulianDay 将年月日转换为儒略日（简化实现）
func (c *PlanetCalculatorFixed) dateToJulianDay(year, month, day int) float64 {
	// 简化的儒略日计算（基于公历1900-2100年范围）
	if month <= 2 {
		year--
		month += 12
	}

	// 计算儒略日
	jd := 365.25*float64(year) + 30.6001*float64(month+1) + float64(day) + 1720994.5

	// 修正格里高利历（1582年10月15日之后）
	if year > 1582 || (year == 1582 && month > 10) || (year == 1582 && month == 10 && day >= 15) {
		a := float64(year / 100)
		b := 2.0 - a + float64(int(a/4))
		jd += b
	}

	return jd
}

// validateParams 验证参数
func (c *PlanetCalculatorFixed) validateParams(params *PlanetParams) error {
	// 验证行星名称
	planetType, err := c.parsePlanetType(params.Planet)
	if err != nil {
		return err
	}

	// 验证儒略日范围
	if params.JD < 0 {
		return fmt.Errorf("儒略日不能为负数: %f", params.JD)
	}

	// 验证经度范围
	if params.Longitude < -180 || params.Longitude > 180 {
		return fmt.Errorf("经度超出范围 (-180到180): %f", params.Longitude)
	}

	// 验证纬度范围
	if params.Latitude < -90 || params.Latitude > 90 {
		return fmt.Errorf("纬度超出范围 (-90到90): %f", params.Latitude)
	}

	_ = planetType // 使用变量避免警告
	return nil
}

// parsePlanetType 解析行星类型
func (c *PlanetCalculatorFixed) parsePlanetType(planet string) (PlanetType, error) {
	switch planet {
	case "mercury":
		return PlanetTypeMercury, nil
	case "venus":
		return PlanetTypeVenus, nil
	case "mars":
		return PlanetTypeMars, nil
	case "jupiter":
		return PlanetTypeJupiter, nil
	case "saturn":
		return PlanetTypeSaturn, nil
	case "uranus":
		return PlanetTypeUranus, nil
	case "neptune":
		return PlanetTypeNeptune, nil
	default:
		return PlanetTypeMercury, fmt.Errorf("不支持的行星: %s", planet)
	}
}

// calculatePlanetPosition 计算行星位置（修复版）
func (c *PlanetCalculatorFixed) calculatePlanetPosition(params *PlanetParams) (*PlanetPositionFixed, error) {
	planetType, _ := c.parsePlanetType(params.Planet)

	// 计算行星和地球的日心坐标
	planetHelio := c.calculateHeliocentricCoordinates(planetType, params.JD)
	earthHelio := c.calculateHeliocentricCoordinates(PlanetTypeEarth, params.JD)

	// 计算地心坐标（行星相对于地球的位置）
	geoX := planetHelio.X - earthHelio.X
	geoY := planetHelio.Y - earthHelio.Y
	geoZ := planetHelio.Z - earthHelio.Z

	// 计算距离
	distance := math.Sqrt(geoX*geoX + geoY*geoY + geoZ*geoZ)

	// 转换为黄道坐标
	lambda := math.Atan2(geoY, geoX) * 180 / math.Pi                         // 黄经（度）
	beta := math.Atan2(geoZ, math.Sqrt(geoX*geoX+geoY*geoY)) * 180 / math.Pi // 黄纬（度）

	// 归一化黄经到 0-360
	lambda = c.normalizeAngle(lambda)

	// 黄道坐标转赤道坐标
	ra, dec := c.eclipticToEquatorial(lambda, beta, params.JD)

	// 计算相位和距角
	phase := c.calculatePhaseFixed(planetHelio, earthHelio)
	elongation := c.calculateElongationFixed(planetHelio, earthHelio)

	// 计算星等
	magnitude := c.calculateMagnitudeFixed(planetType, distance, phase)

	return &PlanetPositionFixed{
		RightAscension: ra,
		Declination:    dec,
		Distance:       distance,
		Magnitude:      magnitude,
		Phase:          phase,
		Elongation:     elongation,
	}, nil
}

// HeliocentricCoordinates 日心坐标
type HeliocentricCoordinates struct {
	X float64 // 日心黄道坐标 X (AU)
	Y float64 // 日心黄道坐标 Y (AU)
	Z float64 // 日心黄道坐标 Z (AU)
}

// PlanetTypeEarth 地球（用于计算）
const PlanetTypeEarth PlanetType = 7

// calculateHeliocentricCoordinates 计算行星的日心坐标
func (c *PlanetCalculatorFixed) calculateHeliocentricCoordinates(planetType PlanetType, jd float64) *HeliocentricCoordinates {
	t := (jd - 2451545.0) / 36525.0 // 儒略世纪数

	// 轨道要素（基于VSOP87简化理论）
	var L, a, e, i, omega, varpi float64

	switch planetType {
	case PlanetTypeMercury:
		L = 252.250906 + 149472.6746358*t
		a = 0.38709893
		e = 0.20563069 + 0.000020406*t
		i = 7.00487
		omega = 48.33167 + 0.003984*t // 升交点经度
		varpi = 77.45645 + 0.160476*t // 近日点经度
	case PlanetTypeVenus:
		L = 181.979801 + 58517.8156760*t
		a = 0.72333199
		e = 0.00677323 - 0.000049214*t
		i = 3.39471
		omega = 76.68069 - 0.001302*t
		varpi = 131.56370 + 0.002683*t
	case PlanetTypeMars:
		L = 355.433000 + 19140.2993034*t
		a = 1.52366231
		e = 0.09341233 + 0.000092064*t
		i = 1.85061
		omega = 49.57854 - 0.002860*t
		varpi = 336.04084 + 0.443901*t
	case PlanetTypeJupiter:
		L = 34.351484 + 3034.9056746*t
		a = 5.202603191
		e = 0.04849485 + 0.000163244*t
		i = 1.30530
		omega = 100.55615 - 0.005691*t
		varpi = 14.33187 + 0.215552*t
	case PlanetTypeSaturn:
		L = 50.077471 + 1222.1137943*t
		a = 9.554909595
		e = 0.05550862 - 0.000346818*t
		i = 2.48446
		omega = 113.71504 - 0.012658*t
		varpi = 93.05724 + 0.566541*t
	case PlanetTypeUranus:
		L = 314.055005 + 428.4669983*t
		a = 19.218446062
		e = 0.04629590 - 0.000027337*t
		i = 0.76986
		omega = 74.22988 - 0.002844*t
		varpi = 173.00529 + 0.089321*t
	case PlanetTypeNeptune:
		L = 304.348665 + 218.4862002*t
		a = 30.110386869
		e = 0.00898809 + 0.000006408*t
		i = 1.76917
		omega = 131.72169 - 0.002653*t
		varpi = 48.12369 + 0.029158*t
	case PlanetTypeEarth:
		L = 100.466457 + 36000.7698278*t
		a = 1.000001018
		e = 0.01670863 - 0.000042037*t
		i = 0.0
		omega = 0.0
		varpi = 102.937348 + 0.322565*t
	}

	// 归一化角度
	L = c.normalizeAngle(L)
	omega = c.normalizeAngle(omega)
	varpi = c.normalizeAngle(varpi)

	// 近心点角距
	omegaBar := varpi - omega

	// 计算偏近点角（开普勒方程迭代求解）
	M := c.normalizeAngle(L-varpi) * math.Pi / 180 // 平近点角（弧度）
	E := M                                         // 初始值
	for iter := 0; iter < 10; iter++ {
		E = E + (M+e*math.Sin(E)-E)/(1-e*math.Cos(E))
	}

	// 计算真近点角
	trueAnomaly := 2 * math.Atan2(math.Sqrt(1+e)*math.Sin(E/2), math.Sqrt(1-e)*math.Cos(E/2)) * 180 / math.Pi

	// 计算日心距离
	r := a * (1 - e*math.Cos(E))

	// 计算日心黄道坐标
	cosO := math.Cos(omega * math.Pi / 180)
	sinO := math.Sin(omega * math.Pi / 180)
	cosI := math.Cos(i * math.Pi / 180)
	sinI := math.Sin(i * math.Pi / 180)

	u := omegaBar + trueAnomaly // 升交点角距 + 真近点角
	cosU := math.Cos(u * math.Pi / 180)
	sinU := math.Sin(u * math.Pi / 180)

	x := r * (cosO*cosU - sinO*sinU*cosI)
	y := r * (sinO*cosU + cosO*sinU*cosI)
	z := r * sinU * sinI

	return &HeliocentricCoordinates{X: x, Y: y, Z: z}
}

// eclipticToEquatorial 黄道坐标转赤道坐标
func (c *PlanetCalculatorFixed) eclipticToEquatorial(lambda, beta, jd float64) (float64, float64) {
	// 计算黄赤交角
	t := (jd - 2451545.0) / 36525.0
	epsilon := 23.4392911 - 0.0130042*t - 1.64e-7*t*t + 5.04e-7*t*t*t

	// 转换为弧度
	lambdaRad := lambda * math.Pi / 180
	betaRad := beta * math.Pi / 180
	epsilonRad := epsilon * math.Pi / 180

	// 黄道坐标转赤道坐标
	sinLambda := math.Sin(lambdaRad)
	cosLambda := math.Cos(lambdaRad)
	sinBeta := math.Sin(betaRad)
	cosBeta := math.Cos(betaRad)
	sinEps := math.Sin(epsilonRad)
	cosEps := math.Cos(epsilonRad)

	// 赤经
	y := sinLambda*cosEps - math.Tan(betaRad)*sinEps
	x := cosLambda
	alpha := math.Atan2(y, x) * 180 / math.Pi / 15 // 转换为小时

	// 归一化赤经到 0-24 小时
	alpha = c.normalizeHour(alpha)

	// 赤纬
	delta := math.Asin(sinBeta*cosEps+cosBeta*sinEps*sinLambda) * 180 / math.Pi

	return alpha, delta
}

// normalizeHour 归一化小时到 0-24 范围
func (c *PlanetCalculatorFixed) normalizeHour(hour float64) float64 {
	for hour < 0 {
		hour += 24
	}
	for hour >= 24 {
		hour -= 24
	}
	return hour
}

// normalizeAngle 归一化角度到 0-360 度
func (c *PlanetCalculatorFixed) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// calculatePhaseFixed 计算相位（修复版）
func (c *PlanetCalculatorFixed) calculatePhaseFixed(planet, earth *HeliocentricCoordinates) float64 {
	// 计算行星-地球-太阳的夹角（相位角）
	// 向量：地球到行星
	dx := planet.X - earth.X
	dy := planet.Y - earth.Y
	dz := planet.Z - earth.Z

	// 向量：地球到太阳（即负的地球坐标）
	sx := -earth.X
	sy := -earth.Y
	sz := -earth.Z

	// 计算点积
	dot := dx*sx + dy*sy + dz*sz
	magD := math.Sqrt(dx*dx + dy*dy + dz*dz)
	magS := math.Sqrt(sx*sx + sy*sy + sz*sz)

	// 相位角（太阳-地球-行星的夹角）
	phaseAngle := math.Acos(dot/(magD*magS)) * 180 / math.Pi

	// 相位 = (1 + cos(相位角)) / 2
	phase := (1 + math.Cos(phaseAngle*math.Pi/180)) / 2

	return phase
}

// calculateElongationFixed 计算距角（修复版）
func (c *PlanetCalculatorFixed) calculateElongationFixed(planet, earth *HeliocentricCoordinates) float64 {
	// 计算距角（太阳-地球-行星的夹角）
	// 向量：地球到行星
	dx := planet.X - earth.X
	dy := planet.Y - earth.Y
	dz := planet.Z - earth.Z

	// 向量：地球到太阳
	sx := -earth.X
	sy := -earth.Y
	sz := -earth.Z

	// 计算点积
	dot := dx*sx + dy*sy + dz*sz
	magD := math.Sqrt(dx*dx + dy*dy + dz*dz)
	magS := math.Sqrt(sx*sx + sy*sy + sz*sz)

	// 距角
	elongation := math.Acos(dot/(magD*magS)) * 180 / math.Pi

	return elongation
}

// calculateMagnitudeFixed 计算星等（修复版）
func (c *PlanetCalculatorFixed) calculateMagnitudeFixed(planetType PlanetType, distance, phase float64) float64 {
	// 行星绝对星等参数（简化模型）
	var H float64
	switch planetType {
	case PlanetTypeMercury:
		H = -0.42
	case PlanetTypeVenus:
		H = -4.40
	case PlanetTypeMars:
		H = -1.52
	case PlanetTypeJupiter:
		H = -9.40
	case PlanetTypeSaturn:
		H = -8.88
	case PlanetTypeUranus:
		H = -7.19
	case PlanetTypeNeptune:
		H = -6.87
	}

	// 简化星等计算：m = H + 5*log10(distance) + 相位修正
	// 距离包括日地距离和地行距离，这里简化处理
	phaseCorrection := 0.1 * (1 - phase)
	magnitude := H + 5*math.Log10(distance) + phaseCorrection

	return magnitude
}
