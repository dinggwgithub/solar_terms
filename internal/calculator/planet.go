package calculator

import (
	"fmt"
	"math"
)

// PlanetCalculator 行星位置计算器
type PlanetCalculator struct {
	*BaseCalculator
}

// NewPlanetCalculator 创建新的行星位置计算器
func NewPlanetCalculator() *PlanetCalculator {
	return &PlanetCalculator{
		BaseCalculator: NewBaseCalculator(
			"planet",
			"行星位置计算器，基于VSOP87理论计算行星的赤经赤纬",
		),
	}
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

// Calculate 执行行星位置计算
func (c *PlanetCalculator) Calculate(params interface{}) (interface{}, error) {
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
func (c *PlanetCalculator) Validate(params interface{}) error {
	planetParams, err := c.parseParams(params)
	if err != nil {
		return err
	}

	return c.validateParams(planetParams)
}

// parseParams 解析参数
func (c *PlanetCalculator) parseParams(params interface{}) (*PlanetParams, error) {
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
func (c *PlanetCalculator) dateToJulianDay(year, month, day int) float64 {
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
func (c *PlanetCalculator) validateParams(params *PlanetParams) error {
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
func (c *PlanetCalculator) parsePlanetType(planet string) (PlanetType, error) {
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

// calculatePlanetPosition 计算行星位置（简化VSOP87实现）
func (c *PlanetCalculator) calculatePlanetPosition(params *PlanetParams) (*PlanetPosition, error) {
	planetType, _ := c.parsePlanetType(params.Planet)

	// 简化的行星位置计算
	// 实际应该基于完整的VSOP87理论

	// 计算轨道参数
	orbitalElements := c.calculateOrbitalElements(planetType, params.JD)

	// 计算赤经赤纬
	ra, dec := c.calculateEquatorialCoordinates(orbitalElements)

	// 计算距离和星等
	distance := c.calculateDistance(planetType, orbitalElements)
	magnitude := c.calculateMagnitude(planetType, distance, orbitalElements.Phase)

	return &PlanetPosition{
		RightAscension: ra,
		Declination:    dec,
		Distance:       distance,
		Magnitude:      magnitude,
		Phase:          orbitalElements.Phase,
		Elongation:     orbitalElements.Elongation,
	}, nil
}

// OrbitalElements 轨道要素
type OrbitalElements struct {
	MeanLongitude float64 // 平黄经（度）
	SemimajorAxis float64 // 半长轴（天文单位）
	Eccentricity  float64 // 偏心率
	Inclination   float64 // 轨道倾角（度）
	LongitudePeri float64 // 近日点经度（度）
	Phase         float64 // 相位（0-1）
	Elongation    float64 // 距角（度）
}

// calculateOrbitalElements 计算轨道要素
func (c *PlanetCalculator) calculateOrbitalElements(planetType PlanetType, jd float64) *OrbitalElements {
	// 简化的轨道要素计算
	// 实际应该基于VSOP87理论

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
	elements.MeanLongitude = c.normalizeAngle(elements.MeanLongitude)
	elements.LongitudePeri = c.normalizeAngle(elements.LongitudePeri)

	// 计算相位和距角
	elements.Phase = c.calculatePhase(planetType, elements)
	elements.Elongation = c.calculateElongation(planetType, elements)

	return elements
}

// calculateEquatorialCoordinates 计算赤道坐标
func (c *PlanetCalculator) calculateEquatorialCoordinates(elements *OrbitalElements) (float64, float64) {
	// 简化的赤道坐标计算
	// 实际应该基于完整的球面天文公式

	// 计算黄道坐标
	lambda := elements.MeanLongitude // 简化：使用平黄经
	beta := elements.Inclination     // 简化：使用轨道倾角

	// 黄道坐标转赤道坐标（简化）
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

	// 转换为度和小时
	ra := alpha * 180 / math.Pi / 15 // 赤经（小时）
	dec := delta * 180 / math.Pi     // 赤纬（度）

	return ra, dec
}

// calculateDistance 计算距离
func (c *PlanetCalculator) calculateDistance(planetType PlanetType, elements *OrbitalElements) float64 {
	// 简化的距离计算
	// 实际应该基于开普勒方程

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

// calculateMagnitude 计算星等
func (c *PlanetCalculator) calculateMagnitude(planetType PlanetType, distance, phase float64) float64 {
	// 简化的星等计算
	// 实际应该基于光度公式

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

	// 距离和相位修正（简化）
	magnitude := baseMagnitude + 5*math.Log10(distance) + 0.1*(1-phase)

	return magnitude
}

// calculatePhase 计算相位
func (c *PlanetCalculator) calculatePhase(planetType PlanetType, elements *OrbitalElements) float64 {
	// 简化的相位计算
	// 实际应该基于行星-地球-太阳的几何关系

	// 内行星和外行星的相位计算不同
	if planetType == PlanetTypeMercury || planetType == PlanetTypeVenus {
		// 内行星相位
		return 0.5 + 0.5*math.Cos(elements.Elongation*math.Pi/180)
	} else {
		// 外行星相位（接近满相）
		return 0.95
	}
}

// calculateElongation 计算距角
func (c *PlanetCalculator) calculateElongation(planetType PlanetType, elements *OrbitalElements) float64 {
	// 简化的距角计算
	// 实际应该基于行星和太阳的黄经差

	if planetType == PlanetTypeMercury || planetType == PlanetTypeVenus {
		// 内行星最大距角约28度（水星）或47度（金星）
		return 30.0
	} else {
		// 外行星距角可以接近180度
		return 120.0
	}
}

// normalizeAngle 归一化角度到0-360度
func (c *PlanetCalculator) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}


// GetPlanetInfo 获取行星基本信息（用于测试）
func (c *PlanetCalculator) GetPlanetInfo(planetType PlanetType) map[string]interface{} {
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
