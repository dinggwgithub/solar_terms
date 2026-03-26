package calculator

import (
	"fmt"
	"math"
)

// planetHelper 内部辅助结构，共享方法
type planetHelper struct{}

// PlanetCalculator 行星位置计算器（原始版本）
type PlanetCalculator struct {
	*BaseCalculator
	helper planetHelper
}

// NewPlanetCalculator 创建新的行星位置计算器
func NewPlanetCalculator() *PlanetCalculator {
	return &PlanetCalculator{
		BaseCalculator: NewBaseCalculator(
			"planet",
			"行星位置计算器，基于VSOP87理论计算行星的赤经赤纬",
		),
		helper: planetHelper{},
	}
}

// PlanetCalculatorFixed 修复版行星位置计算器
type PlanetCalculatorFixed struct {
	*BaseCalculator
	helper planetHelper
}

// NewPlanetCalculatorFixed 创建修复版行星位置计算器
func NewPlanetCalculatorFixed() *PlanetCalculatorFixed {
	return &PlanetCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"planet_fixed",
			"修复版行星位置计算器，基于开普勒轨道要素计算行星的赤经赤纬",
		),
		helper: planetHelper{},
	}
}

// Calculate 执行行星位置计算（修复版）
func (c *PlanetCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	planetParams, err := c.helper.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.helper.validateParams(planetParams); err != nil {
		return nil, err
	}

	// 执行行星位置计算
	position, err := c.helper.calculatePlanetPosition(planetParams)
	if err != nil {
		return nil, err
	}

	return position, nil
}

// Validate 验证输入参数（修复版）
func (c *PlanetCalculatorFixed) Validate(params interface{}) error {
	planetParams, err := c.helper.parseParams(params)
	if err != nil {
		return err
	}

	return c.helper.validateParams(planetParams)
}

// PlanetType 行星类型枚举
type PlanetType int

const (
	PlanetTypeMercury PlanetType = iota // 水星
	PlanetTypeVenus                     // 金星
	PlanetTypeMars                      // 火星
	PlanetTypeJupiter                   // 木星
	PlanetTypeSaturn                    // 土星
	PlanetTypeUranus                    // 天王星
	PlanetTypeNeptune                   // 海王星
)

// String 返回行星类型的字符串表示
func (pt PlanetType) String() string {
	switch pt {
	case PlanetTypeMercury:
		return "mercury"
	case PlanetTypeVenus:
		return "venus"
	case PlanetTypeMars:
		return "mars"
	case PlanetTypeJupiter:
		return "jupiter"
	case PlanetTypeSaturn:
		return "saturn"
	case PlanetTypeUranus:
		return "uranus"
	case PlanetTypeNeptune:
		return "neptune"
	default:
		return "unknown"
	}
}

// PlanetParams 行星计算参数
type PlanetParams struct {
	Planet    string  `json:"planet"`    // 行星名称
	JD        float64 `json:"jd"`        // 儒略日
	Longitude float64 `json:"longitude"` // 观测者经度（度）
	Latitude  float64 `json:"latitude"`  // 观测者纬度（度）
}

// PlanetPosition 行星位置结果
type PlanetPosition struct {
	RightAscension float64 `json:"right_ascension"` // 赤经（小时）
	Declination    float64 `json:"declination"`     // 赤纬（度）
	Distance       float64 `json:"distance"`        // 距离（天文单位）
	Magnitude      float64 `json:"magnitude"`       // 星等
	Phase          float64 `json:"phase"`           // 相位（0-1）
	Elongation     float64 `json:"elongation"`      // 距角（度）
}

// Calculate 执行行星位置计算（原始版本 - 保留旧逻辑用于对比）
func (c *PlanetCalculator) Calculate(params interface{}) (interface{}, error) {
	planetParams, err := c.helper.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.helper.validateParams(planetParams); err != nil {
		return nil, err
	}

	// 执行原始的行星位置计算（旧逻辑，用于对比）
	position, err := c.helper.calculatePlanetPositionOriginal(planetParams)
	if err != nil {
		return nil, err
	}

	return position, nil
}

// Validate 验证输入参数
func (c *PlanetCalculator) Validate(params interface{}) error {
	planetParams, err := c.helper.parseParams(params)
	if err != nil {
		return err
	}

	return c.helper.validateParams(planetParams)
}

// parseParams 解析参数
func (h planetHelper) parseParams(params interface{}) (*PlanetParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	// 辅助函数：转换为字符串
	toString := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}

	// 提取行星名称（支持 planet 和 planet_name 两种字段名）
	planet := ""
	if p, ok := paramsMap["planet"]; ok {
		planet = toString(p)
	} else if p, ok := paramsMap["planet_name"]; ok {
		planet = toString(p)
	} else {
		return nil, fmt.Errorf("planet或planet_name参数必须为字符串")
	}

	// 辅助函数：转换为float64
	toFloat := func(v interface{}) (float64, bool) {
		switch val := v.(type) {
		case float64:
			return val, true
		case int:
			return float64(val), true
		case int64:
			return float64(val), true
		case int32:
			return float64(val), true
		case float32:
			return float64(val), true
		default:
			return 0, false
		}
	}

	// 提取儒略日（支持直接传入jd或通过year/month/day计算）
	var jd float64
	if jdValue, ok := paramsMap["jd"]; ok {
		if jdFloat, ok := toFloat(jdValue); ok {
			jd = jdFloat
		}
	} else {
		// 尝试从year/month/day计算儒略日
		var year, month, day float64
		var yOk, mOk, dOk bool

		if y, ok := paramsMap["year"]; ok {
			year, yOk = toFloat(y)
		}
		if m, ok := paramsMap["month"]; ok {
			month, mOk = toFloat(m)
		}
		if d, ok := paramsMap["day"]; ok {
			day, dOk = toFloat(d)
		}

		if yOk && mOk && dOk {
			jd = h.dateToJulianDay(int(year), int(month), int(day))
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
func (h planetHelper) dateToJulianDay(year, month, day int) float64 {
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
func (h planetHelper) validateParams(params *PlanetParams) error {
	// 验证行星名称
	planetType, err := h.parsePlanetType(params.Planet)
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
func (h planetHelper) parsePlanetType(planet string) (PlanetType, error) {
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

// calculatePlanetPositionOriginal 原始的行星位置计算（旧逻辑，用于对比测试）
func (h planetHelper) calculatePlanetPositionOriginal(params *PlanetParams) (*PlanetPosition, error) {
	planetType, _ := h.parsePlanetType(params.Planet)

	// 原始简化计算逻辑
	orbitalElements := h.calculateOrbitalElementsOriginal(planetType, params.JD)

	// 计算赤经赤纬（旧逻辑）
	ra, dec := h.calculateEquatorialCoordinatesOriginal(orbitalElements)

	// 计算距离和星等
	distance := h.calculateDistanceOriginal(planetType)
	magnitude := h.calculateMagnitudeOriginal(planetType, distance, orbitalElements.Phase)

	return &PlanetPosition{
		RightAscension: ra,
		Declination:    dec,
		Distance:       distance,
		Magnitude:      magnitude,
		Phase:          orbitalElements.Phase,
		Elongation:     orbitalElements.Elongation,
	}, nil
}

// calculateOrbitalElementsOriginal 原始轨道要素计算（旧逻辑）
func (h planetHelper) calculateOrbitalElementsOriginal(planetType PlanetType, jd float64) *OrbitalElements {
	t := (jd - 2451545.0) / 36525.0 // 儒略世纪数

	elements := &OrbitalElements{}

	switch planetType {
	case PlanetTypeMercury:
		elements.MeanLongitude = 252.250906 + 149472.6746358*t
		elements.SemimajorAxis = 0.38709893
		elements.Eccentricity = 0.20563069 + 0.000020406*t
		elements.Inclination = 7.00487
		elements.LongitudePeri = 77.45645 + 0.160476*t

	case PlanetTypeVenus:
		elements.MeanLongitude = 181.979801 + 58517.8156760*t
		elements.SemimajorAxis = 0.72333199
		elements.Eccentricity = 0.00677323 - 0.000049214*t
		elements.Inclination = 3.39471
		elements.LongitudePeri = 131.56370 + 0.002683*t

	case PlanetTypeMars:
		elements.MeanLongitude = 355.433000 + 19140.2993034*t
		elements.SemimajorAxis = 1.52366231
		elements.Eccentricity = 0.09341233 + 0.000092064*t
		elements.Inclination = 1.85061
		elements.LongitudePeri = 336.04084 + 0.443901*t

	case PlanetTypeJupiter:
		elements.MeanLongitude = 34.351484 + 3034.9056746*t
		elements.SemimajorAxis = 5.202603191
		elements.Eccentricity = 0.04849485 + 0.000163244*t
		elements.Inclination = 1.30530
		elements.LongitudePeri = 14.33187 + 0.215552*t

	case PlanetTypeSaturn:
		elements.MeanLongitude = 50.077471 + 1222.1137943*t
		elements.SemimajorAxis = 9.554909595
		elements.Eccentricity = 0.05550862 - 0.000346818*t
		elements.Inclination = 2.48446
		elements.LongitudePeri = 93.05724 + 0.566541*t

	case PlanetTypeUranus:
		elements.MeanLongitude = 314.055005 + 428.4669983*t
		elements.SemimajorAxis = 19.218446062
		elements.Eccentricity = 0.04629590 - 0.000027337*t
		elements.Inclination = 0.76986
		elements.LongitudePeri = 173.00529 + 0.089321*t

	case PlanetTypeNeptune:
		elements.MeanLongitude = 304.348665 + 218.4862002*t
		elements.SemimajorAxis = 30.110386869
		elements.Eccentricity = 0.00898809 + 0.000006408*t
		elements.Inclination = 1.76917
		elements.LongitudePeri = 48.12369 + 0.029158*t
	}

	// 归一化角度
	elements.MeanLongitude = h.normalizeAngle(elements.MeanLongitude)
	elements.LongitudePeri = h.normalizeAngle(elements.LongitudePeri)

	// 原始简化的相位和距角
	elements.Phase = h.calculatePhaseOriginal(planetType)
	elements.Elongation = h.calculateElongationOriginal(planetType)

	return elements
}

// calculateEquatorialCoordinatesOriginal 原始赤道坐标计算（旧逻辑）
func (h planetHelper) calculateEquatorialCoordinatesOriginal(elements *OrbitalElements) (float64, float64) {
	lambda := elements.MeanLongitude // 简化：使用平黄经
	beta := elements.Inclination     // 简化：使用轨道倾角

	epsilon := 23.4392911 // 黄赤交角

	alpha := math.Atan2(
		math.Sin(lambda*math.Pi/180)*math.Cos(epsilon*math.Pi/180)-
			math.Tan(beta*math.Pi/180)*math.Sin(epsilon*math.Pi/180),
		math.Cos(lambda*math.Pi/180),
	)
	delta := math.Asin(
		math.Sin(beta*math.Pi/180)*math.Cos(epsilon*math.Pi/180) +
			math.Cos(beta*math.Pi/180)*math.Sin(epsilon*math.Pi/180)*math.Sin(lambda*math.Pi/180),
	)

	// 转换为度和小时（注意：旧逻辑未进行归一化，导致负值）
	ra := alpha * 180 / math.Pi / 15 // 赤经（小时，可能为负）
	dec := delta * 180 / math.Pi     // 赤纬（度）

	return ra, dec
}

// calculateDistanceOriginal 原始距离计算（旧逻辑，硬编码）
func (h planetHelper) calculateDistanceOriginal(planetType PlanetType) float64 {
	switch planetType {
	case PlanetTypeMercury:
		return 0.5
	case PlanetTypeVenus:
		return 0.7
	case PlanetTypeMars:
		return 1.2
	case PlanetTypeJupiter:
		return 4.5
	case PlanetTypeSaturn:
		return 9.0
	case PlanetTypeUranus:
		return 18.0
	case PlanetTypeNeptune:
		return 30.0
	default:
		return 1.0
	}
}

// calculateMagnitudeOriginal 原始星等计算（旧逻辑）
func (h planetHelper) calculateMagnitudeOriginal(planetType PlanetType, distance, phase float64) float64 {
	baseMagnitude := 0.0
	switch planetType {
	case PlanetTypeMercury:
		baseMagnitude = -0.42
	case PlanetTypeVenus:
		baseMagnitude = -4.40
	case PlanetTypeMars:
		baseMagnitude = -1.52
	case PlanetTypeJupiter:
		baseMagnitude = -2.94
	case PlanetTypeSaturn:
		baseMagnitude = 0.43
	case PlanetTypeUranus:
		baseMagnitude = 5.32
	case PlanetTypeNeptune:
		baseMagnitude = 7.78
	}

	magnitude := baseMagnitude + 5*math.Log10(distance) + 0.1*(1-phase)

	return magnitude
}

// calculatePhaseOriginal 原始相位计算（旧逻辑，硬编码）
func (h planetHelper) calculatePhaseOriginal(planetType PlanetType) float64 {
	if planetType == PlanetTypeMercury || planetType == PlanetTypeVenus {
		return 0.5
	} else {
		return 0.95
	}
}

// calculateElongationOriginal 原始距角计算（旧逻辑，硬编码）
func (h planetHelper) calculateElongationOriginal(planetType PlanetType) float64 {
	if planetType == PlanetTypeMercury || planetType == PlanetTypeVenus {
		return 30.0
	} else {
		return 120.0
	}
}

// OrbitalElements 轨道要素
type OrbitalElements struct {
	MeanLongitude float64 // 平黄经（度）
	TrueLongitude float64 // 真黄经（度）
	SemimajorAxis float64 // 半长轴（天文单位）
	Eccentricity  float64 // 偏心率
	Inclination   float64 // 轨道倾角（度）
	LongitudePeri float64 // 近日点经度（度）
	TrueAnomaly   float64 // 真近点角（度）
	RadiusVector  float64 // 向径（天文单位）
	Phase         float64 // 相位（0-1）
	Elongation    float64 // 距角（度）
}

// calculatePlanetPosition 计算行星位置（基于开普勒轨道要素）
func (h planetHelper) calculatePlanetPosition(params *PlanetParams) (*PlanetPosition, error) {
	planetType, _ := h.parsePlanetType(params.Planet)

	// 计算轨道参数
	orbitalElements := h.calculateOrbitalElements(planetType, params.JD)

	// 计算赤经赤纬
	ra, dec := h.calculateEquatorialCoordinates(orbitalElements)

	// 获取地球轨道要素
	earthElements := h.calculateEarthOrbitalElements(params.JD)

	// 计算距离（从地球到行星）
	distance := h.calculateDistance(orbitalElements, earthElements)

	// 计算星等
	magnitude := h.calculateMagnitude(planetType, orbitalElements, earthElements, distance)

	return &PlanetPosition{
		RightAscension: ra,
		Declination:    dec,
		Distance:       distance,
		Magnitude:      magnitude,
		Phase:          orbitalElements.Phase,
		Elongation:     orbitalElements.Elongation,
	}, nil
}

// calculateOrbitalElements 计算轨道要素
func (h planetHelper) calculateOrbitalElements(planetType PlanetType, jd float64) *OrbitalElements {
	t := (jd - 2451545.0) / 36525.0 // 儒略世纪数（J2000历元）

	elements := &OrbitalElements{}

	// 行星轨道要素（J2000）
	switch planetType {
	case PlanetTypeMercury:
		elements.MeanLongitude = h.normalizeAngle(252.250906 + 149472.6746358*t)
		elements.SemimajorAxis = 0.38709893
		elements.Eccentricity = 0.20563069 + 0.000020406*t
		elements.Inclination = 7.00487
		elements.LongitudePeri = h.normalizeAngle(77.45645 + 0.160476*t)

	case PlanetTypeVenus:
		elements.MeanLongitude = h.normalizeAngle(181.979801 + 58517.8156760*t)
		elements.SemimajorAxis = 0.72333199
		elements.Eccentricity = 0.00677323 - 0.000049214*t
		elements.Inclination = 3.39471
		elements.LongitudePeri = h.normalizeAngle(131.56370 + 0.002683*t)

	case PlanetTypeMars:
		elements.MeanLongitude = h.normalizeAngle(355.433000 + 19140.2993034*t)
		elements.SemimajorAxis = 1.52366231
		elements.Eccentricity = 0.09341233 + 0.000092064*t
		elements.Inclination = 1.85061
		elements.LongitudePeri = h.normalizeAngle(336.04084 + 0.443901*t)

	case PlanetTypeJupiter:
		elements.MeanLongitude = h.normalizeAngle(34.351484 + 3034.9056746*t)
		elements.SemimajorAxis = 5.202603191
		elements.Eccentricity = 0.04849485 + 0.000163244*t
		elements.Inclination = 1.30530
		elements.LongitudePeri = h.normalizeAngle(14.33187 + 0.215552*t)

	case PlanetTypeSaturn:
		elements.MeanLongitude = h.normalizeAngle(50.077471 + 1222.1137943*t)
		elements.SemimajorAxis = 9.554909595
		elements.Eccentricity = 0.05550862 - 0.000346818*t
		elements.Inclination = 2.48446
		elements.LongitudePeri = h.normalizeAngle(93.05724 + 0.566541*t)

	case PlanetTypeUranus:
		elements.MeanLongitude = h.normalizeAngle(314.055005 + 428.4669983*t)
		elements.SemimajorAxis = 19.218446062
		elements.Eccentricity = 0.04629590 - 0.000027337*t
		elements.Inclination = 0.76986
		elements.LongitudePeri = h.normalizeAngle(173.00529 + 0.089321*t)

	case PlanetTypeNeptune:
		elements.MeanLongitude = h.normalizeAngle(304.348665 + 218.4862002*t)
		elements.SemimajorAxis = 30.110386869
		elements.Eccentricity = 0.00898809 + 0.000006408*t
		elements.Inclination = 1.76917
		elements.LongitudePeri = h.normalizeAngle(48.12369 + 0.029158*t)
	}

	// 计算平近点角 M = L - ω
	M := elements.MeanLongitude - elements.LongitudePeri
	M = h.normalizeAngle(M)

	// 求解开普勒方程 E - e*sin(E) = M
	MRad := M * math.Pi / 180
	ERad := h.solveKepler(MRad, elements.Eccentricity)

	// 计算真近点角 ν（使用 tan(ν/2) = sqrt((1+e)/(1-e)) * tan(E/2)）
	nuRad := 2 * math.Atan2(
		math.Sqrt((1+elements.Eccentricity)/(1-elements.Eccentricity))*math.Sin(ERad/2),
		math.Cos(ERad/2),
	)
	elements.TrueAnomaly = h.normalizeAngle(nuRad * 180 / math.Pi)

	// 计算向径 r = a*(1-e²)/(1+e*cos(ν))
	elements.RadiusVector = elements.SemimajorAxis * (1 - elements.Eccentricity*elements.Eccentricity) /
		(1 + elements.Eccentricity*math.Cos(nuRad))

	// 计算真黄经 θ = ω + ν
	elements.TrueLongitude = h.normalizeAngle(elements.LongitudePeri + elements.TrueAnomaly)

	// 计算相位和距角（需要地球位置）
	earthElements := h.calculateEarthOrbitalElements(jd)
	elements.Phase = h.calculatePhase(planetType, elements, earthElements)
	elements.Elongation = h.calculateElongation(planetType, elements, earthElements)

	return elements
}

// calculateEarthOrbitalElements 计算地球轨道要素（用于相位和距角计算）
func (h planetHelper) calculateEarthOrbitalElements(jd float64) *OrbitalElements {
	t := (jd - 2451545.0) / 36525.0 // 儒略世纪数

	elements := &OrbitalElements{}

	// 地球轨道要素（J2000）
	elements.MeanLongitude = h.normalizeAngle(100.46435 + 129597740.63*t/3600)
	elements.SemimajorAxis = 1.00000011
	elements.Eccentricity = 0.01671022 - 0.00003804*t
	elements.Inclination = 0.00005
	elements.LongitudePeri = h.normalizeAngle(102.94719 + 0.3225654*t)

	// 计算平近点角
	M := elements.MeanLongitude - elements.LongitudePeri
	M = h.normalizeAngle(M)

	// 求解开普勒方程
	MRad := M * math.Pi / 180
	ERad := h.solveKepler(MRad, elements.Eccentricity)

	// 计算真近点角
	nuRad := 2 * math.Atan2(
		math.Sqrt((1+elements.Eccentricity)/(1-elements.Eccentricity))*math.Sin(ERad/2),
		math.Cos(ERad/2),
	)
	elements.TrueAnomaly = h.normalizeAngle(nuRad * 180 / math.Pi)

	// 计算向径
	elements.RadiusVector = elements.SemimajorAxis * (1 - elements.Eccentricity*elements.Eccentricity) /
		(1 + elements.Eccentricity*math.Cos(nuRad))

	// 计算真黄经
	elements.TrueLongitude = h.normalizeAngle(elements.LongitudePeri + elements.TrueAnomaly)

	return elements
}

// calculateEquatorialCoordinates 计算赤道坐标
func (h planetHelper) calculateEquatorialCoordinates(elements *OrbitalElements) (float64, float64) {
	// 黄道坐标转赤道坐标
	lambda := elements.TrueLongitude // 真黄经
	beta := elements.Inclination     // 黄纬

	epsilon := 23.4392911 // 黄赤交角（J2000）

	// 转换为弧度
	lambdaRad := lambda * math.Pi / 180
	betaRad := beta * math.Pi / 180
	epsilonRad := epsilon * math.Pi / 180

	// 计算赤道坐标（赤经α，赤纬δ）
	alpha := math.Atan2(
		math.Sin(lambdaRad)*math.Cos(epsilonRad)-
			math.Tan(betaRad)*math.Sin(epsilonRad),
		math.Cos(lambdaRad),
	)
	delta := math.Asin(
		math.Sin(betaRad)*math.Cos(epsilonRad) +
			math.Cos(betaRad)*math.Sin(epsilonRad)*math.Sin(lambdaRad),
	)

	// 转换为小时（赤经）和度（赤纬），并归一化范围
	ra := h.normalizeHours(alpha * 180 / math.Pi / 15) // 转换为小时并归一化到0-24
	dec := delta * 180 / math.Pi                       // 赤纬范围-90到+90，无需归一化

	return ra, dec
}

// calculateDistance 计算从地球到行星的距离
func (h planetHelper) calculateDistance(planetElements, earthElements *OrbitalElements) float64 {
	// 行星的日心黄道坐标
	r1 := planetElements.RadiusVector
	theta1 := planetElements.TrueLongitude * math.Pi / 180
	phi1 := planetElements.Inclination * math.Pi / 180

	// 地球的日心黄道坐标
	r2 := earthElements.RadiusVector
	theta2 := earthElements.TrueLongitude * math.Pi / 180
	phi2 := earthElements.Inclination * math.Pi / 180 // 接近0

	// 转换为笛卡尔坐标（黄道坐标系）
	x1 := r1 * math.Cos(phi1) * math.Cos(theta1)
	y1 := r1 * math.Cos(phi1) * math.Sin(theta1)
	z1 := r1 * math.Sin(phi1)

	x2 := r2 * math.Cos(phi2) * math.Cos(theta2)
	y2 := r2 * math.Cos(phi2) * math.Sin(theta2)
	z2 := r2 * math.Sin(phi2)

	// 计算距离
	dx := x1 - x2
	dy := y1 - y2
	dz := z1 - z2

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// calculateMagnitude 计算星等
func (h planetHelper) calculateMagnitude(planetType PlanetType, planetElements, earthElements *OrbitalElements, distance float64) float64 {
	// 计算相位角
	phaseAngle := h.calculatePhaseAngle(planetElements, earthElements)

	// 基V星等（冲日时视星等）
	V0 := 0.0
	// 相位系数（用于相位修正）
	phaseCoeff := 0.0

	switch planetType {
	case PlanetTypeMercury:
		V0 = -0.42
		phaseCoeff = 0.035 // 简化的相位系数
	case PlanetTypeVenus:
		V0 = -4.40
		phaseCoeff = 0.013
	case PlanetTypeMars:
		V0 = -1.52
		phaseCoeff = 0.016
	case PlanetTypeJupiter:
		V0 = -2.94
		phaseCoeff = 0.014
	case PlanetTypeSaturn:
		V0 = 0.43
		phaseCoeff = 0.012
		// 土星需要考虑环的影响，这里简化处理
	case PlanetTypeUranus:
		V0 = 5.32
		phaseCoeff = 0.005
	case PlanetTypeNeptune:
		V0 = 7.78
		phaseCoeff = 0.003
	}

	// 星等计算公式：V = V0 + 5*log10(r*R/au²) + f(相位角)
	// r: 行星到太阳的距离（au）
	// R: 地球到行星的距离（au）
	r := planetElements.RadiusVector

	// 距离修正项
	distanceCorrection := 5.0 * math.Log10(r*distance)

	// 相位修正项（简化模型）
	phaseCorrection := phaseCoeff * phaseAngle

	magnitude := V0 + distanceCorrection + phaseCorrection

	return magnitude
}

// calculatePhaseAngle 计算相位角（太阳-行星-地球的夹角，度）
func (h planetHelper) calculatePhaseAngle(planetElements, earthElements *OrbitalElements) float64 {
	r := planetElements.RadiusVector                        // 行星到太阳距离
	R := earthElements.RadiusVector                         // 地球到太阳距离
	d := h.calculateDistance(planetElements, earthElements) // 地球到行星距离

	// 使用余弦定理计算相位角
	// cos(phaseAngle) = (r² + d² - R²) / (2rd)
	cosPhase := (r*r + d*d - R*R) / (2 * r * d)

	// 确保在[-1, 1]范围内
	cosPhase = math.Max(-1.0, math.Min(1.0, cosPhase))

	phaseAngle := math.Acos(cosPhase) * 180 / math.Pi

	return phaseAngle
}

// calculatePhase 计算相位（被照亮部分的比例，0-1）
func (h planetHelper) calculatePhase(planetType PlanetType, planetElements, earthElements *OrbitalElements) float64 {
	phaseAngle := h.calculatePhaseAngle(planetElements, earthElements)

	// 相位 = (1 + cos(相位角)) / 2
	phase := (1 + math.Cos(phaseAngle*math.Pi/180)) / 2

	return phase
}

// calculateElongation 计算距角（太阳-地球-行星的夹角，度）
func (h planetHelper) calculateElongation(planetType PlanetType, planetElements, earthElements *OrbitalElements) float64 {
	r := planetElements.RadiusVector                        // 行星到太阳距离
	R := earthElements.RadiusVector                         // 地球到太阳距离
	d := h.calculateDistance(planetElements, earthElements) // 地球到行星距离

	// 使用余弦定理计算距角
	// cos(elongation) = (R² + d² - r²) / (2Rd)
	cosElong := (R*R + d*d - r*r) / (2 * R * d)

	// 确保在[-1, 1]范围内
	cosElong = math.Max(-1.0, math.Min(1.0, cosElong))

	elongation := math.Acos(cosElong) * 180 / math.Pi

	// 计算黄经差以确定行星在太阳的东侧还是西侧
	deltaLong := planetElements.TrueLongitude - earthElements.TrueLongitude
	deltaLong = h.normalizeAngle(deltaLong)

	// 如果黄经差大于180度，距角应该是360 - 计算值（表示西侧）
	if deltaLong > 180 {
		elongation = 360 - elongation
	}

	return math.Abs(elongation)
}

// normalizeAngle 归一化角度到0-360度
func (h planetHelper) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// normalizeHours 归一化小时数到0-24范围
func (h planetHelper) normalizeHours(hours float64) float64 {
	for hours < 0 {
		hours += 24
	}
	for hours >= 24 {
		hours -= 24
	}
	return hours
}

// solveKepler 求解开普勒方程（迭代法）
func (h planetHelper) solveKepler(M, e float64) float64 {
	// M: 平近点角（弧度）
	// e: 偏心率
	// 返回: 偏近点角E（弧度）

	E := M          // 初始值
	epsilon := 1e-8 // 收敛精度
	maxIter := 100

	for i := 0; i < maxIter; i++ {
		dE := (E - e*math.Sin(E) - M) / (1 - e*math.Cos(E))
		E -= dE
		if math.Abs(dE) < epsilon {
			break
		}
	}

	return E
}

// GetPlanetInfo 获取行星基本信息（用于测试）
func (h planetHelper) GetPlanetInfo(planetType PlanetType) map[string]interface{} {
	info := map[string]interface{}{
		"name": planetType.String(),
		"type": "planet",
	}

	switch planetType {
	case PlanetTypeMercury:
		info["radius"] = 2439.7        // 公里
		info["mass"] = 3.301e23        // 千克
		info["orbital_period"] = 87.97 // 天

	case PlanetTypeVenus:
		info["radius"] = 6051.8
		info["mass"] = 4.867e24
		info["orbital_period"] = 224.70

	case PlanetTypeMars:
		info["radius"] = 3389.5
		info["mass"] = 6.417e23
		info["orbital_period"] = 686.98

	case PlanetTypeJupiter:
		info["radius"] = 69911
		info["mass"] = 1.898e27
		info["orbital_period"] = 4332.59

	case PlanetTypeSaturn:
		info["radius"] = 58232
		info["mass"] = 5.683e26
		info["orbital_period"] = 10759.22

	case PlanetTypeUranus:
		info["radius"] = 25362
		info["mass"] = 8.681e25
		info["orbital_period"] = 30688.5

	case PlanetTypeNeptune:
		info["radius"] = 24622
		info["mass"] = 1.024e26
		info["orbital_period"] = 60182
	}

	return info
}
