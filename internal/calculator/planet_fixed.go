package calculator

import (
	"fmt"
	"math"
)

type PlanetCalculatorFixed struct {
	*BaseCalculator
}

func NewPlanetCalculatorFixed() *PlanetCalculatorFixed {
	return &PlanetCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"planet_fixed",
			"修复版行星位置计算器，修正赤经归一化和轨道参数计算",
		),
	}
}

type PlanetPositionFixed struct {
	RightAscension float64 `json:"right_ascension"`
	Declination    float64 `json:"declination"`
	Distance       float64 `json:"distance"`
	Magnitude      float64 `json:"magnitude"`
	Phase          float64 `json:"phase"`
	Elongation     float64 `json:"elongation"`
}

func (c *PlanetCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	planetParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateParams(planetParams); err != nil {
		return nil, err
	}

	position, err := c.calculatePlanetPosition(planetParams)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func (c *PlanetCalculatorFixed) Validate(params interface{}) error {
	planetParams, err := c.parseParams(params)
	if err != nil {
		return err
	}
	return c.validateParams(planetParams)
}

func (c *PlanetCalculatorFixed) parseParams(params interface{}) (*PlanetParams, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
	}

	planet := ""
	if p, ok := paramsMap["planet"].(string); ok {
		planet = p
	} else if p, ok := paramsMap["planet_name"].(string); ok {
		planet = p
	} else {
		return nil, fmt.Errorf("planet或planet_name参数必须为字符串")
	}

	var jd float64
	if jdValue, ok := paramsMap["jd"].(float64); ok {
		jd = jdValue
	} else {
		year, yOk := paramsMap["year"].(float64)
		month, mOk := paramsMap["month"].(float64)
		day, dOk := paramsMap["day"].(float64)

		if yOk && mOk && dOk {
			jd = c.dateToJulianDay(int(year), int(month), int(day))
		} else {
			return nil, fmt.Errorf("需要提供jd参数或year/month/day参数")
		}
	}

	longitude := 0.0
	if lon, exists := paramsMap["longitude"]; exists {
		if lonFloat, ok := lon.(float64); ok {
			longitude = lonFloat
		}
	}

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

func (c *PlanetCalculatorFixed) dateToJulianDay(year, month, day int) float64 {
	if month <= 2 {
		year--
		month += 12
	}

	a := float64(year / 100)
	b := 2.0 - a + float64(int(a/4))

	jd := math.Floor(365.25*float64(year+4716)) + math.Floor(30.6001*float64(month+1)) + float64(day) + b - 1524.5

	return jd
}

func (c *PlanetCalculatorFixed) validateParams(params *PlanetParams) error {
	_, err := c.parsePlanetType(params.Planet)
	if err != nil {
		return err
	}

	if params.JD < 0 {
		return fmt.Errorf("儒略日不能为负数: %f", params.JD)
	}

	if params.Longitude < -180 || params.Longitude > 180 {
		return fmt.Errorf("经度超出范围 (-180到180): %f", params.Longitude)
	}

	if params.Latitude < -90 || params.Latitude > 90 {
		return fmt.Errorf("纬度超出范围 (-90到90): %f", params.Latitude)
	}

	return nil
}

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

func (c *PlanetCalculatorFixed) calculatePlanetPosition(params *PlanetParams) (*PlanetPositionFixed, error) {
	planetType, _ := c.parsePlanetType(params.Planet)

	t := (params.JD - 2451545.0) / 36525.0

	elements := c.calculateOrbitalElements(planetType, t)

	sunLon := c.calculateSunLongitude(params.JD)

	ra, dec := c.calculateEquatorialCoordinates(elements, sunLon)

	distance := c.calculateDistance(planetType, elements)

	elongation := c.calculateElongation(elements, sunLon)

	phase := c.calculatePhase(elongation)

	magnitude := c.calculateMagnitude(planetType, distance, phase)

	return &PlanetPositionFixed{
		RightAscension: ra,
		Declination:    dec,
		Distance:       distance,
		Magnitude:      magnitude,
		Phase:          phase,
		Elongation:     elongation,
	}, nil
}

type OrbitalElementsFixed struct {
	MeanLongitude float64
	SemimajorAxis float64
	Eccentricity  float64
	Inclination   float64
	LongitudePeri float64
	LongitudeNode float64
	MeanAnomaly   float64
	TrueLongitude float64
	RadiusVector  float64
	EclipticLon   float64
	EclipticLat   float64
}

func (c *PlanetCalculatorFixed) calculateOrbitalElements(planetType PlanetType, t float64) *OrbitalElementsFixed {
	elements := &OrbitalElementsFixed{}

	switch planetType {
	case PlanetTypeMars:
		elements.MeanLongitude = 355.433000 + 19140.2993034*t + 0.0000*t*t
		elements.SemimajorAxis = 1.52366231
		elements.Eccentricity = 0.09341233 + 0.000092064*t
		elements.Inclination = 1.85061 - 0.000675*t
		elements.LongitudePeri = 336.04084 + 0.443901*t
		elements.LongitudeNode = 49.58854 + 0.0000*t
	case PlanetTypeMercury:
		elements.MeanLongitude = 252.250906 + 149472.6746358*t
		elements.SemimajorAxis = 0.38709893
		elements.Eccentricity = 0.20563069 + 0.000020406*t
		elements.Inclination = 7.00487 - 0.0001781*t
		elements.LongitudePeri = 77.45645 + 0.160476*t
		elements.LongitudeNode = 48.33089 - 0.0001261*t
	case PlanetTypeVenus:
		elements.MeanLongitude = 181.979801 + 58517.8156760*t
		elements.SemimajorAxis = 0.72333199
		elements.Eccentricity = 0.00677323 - 0.000049214*t
		elements.Inclination = 3.39471 - 0.0000015*t
		elements.LongitudePeri = 131.56370 + 0.002683*t
		elements.LongitudeNode = 76.67994 - 0.0002778*t
	case PlanetTypeJupiter:
		elements.MeanLongitude = 34.351484 + 3034.9056746*t
		elements.SemimajorAxis = 5.202603191
		elements.Eccentricity = 0.04849485 + 0.000163244*t
		elements.Inclination = 1.30530 - 0.0000051*t
		elements.LongitudePeri = 14.33187 + 0.215552*t
		elements.LongitudeNode = 100.46435 + 0.0000*t
	case PlanetTypeSaturn:
		elements.MeanLongitude = 50.077471 + 1222.1137943*t
		elements.SemimajorAxis = 9.554909595
		elements.Eccentricity = 0.05550862 - 0.000346818*t
		elements.Inclination = 2.48446 - 0.0000019*t
		elements.LongitudePeri = 93.05724 + 0.566541*t
		elements.LongitudeNode = 113.66552 + 0.0000*t
	case PlanetTypeUranus:
		elements.MeanLongitude = 314.055005 + 428.4669983*t
		elements.SemimajorAxis = 19.218446062
		elements.Eccentricity = 0.04629590 - 0.000027337*t
		elements.Inclination = 0.76986 - 0.0000002*t
		elements.LongitudePeri = 173.00529 + 0.089321*t
		elements.LongitudeNode = 74.00600 + 0.0000*t
	case PlanetTypeNeptune:
		elements.MeanLongitude = 304.348665 + 218.4862002*t
		elements.SemimajorAxis = 30.110386869
		elements.Eccentricity = 0.00898809 + 0.000006408*t
		elements.Inclination = 1.76917 - 0.0000002*t
		elements.LongitudePeri = 48.12369 + 0.029158*t
		elements.LongitudeNode = 131.78406 + 0.0000*t
	}

	elements.MeanLongitude = c.normalizeAngle(elements.MeanLongitude)
	elements.LongitudePeri = c.normalizeAngle(elements.LongitudePeri)
	elements.LongitudeNode = c.normalizeAngle(elements.LongitudeNode)

	elements.MeanAnomaly = c.normalizeAngle(elements.MeanLongitude - elements.LongitudePeri)

	E := c.solveKeplerEquation(elements.MeanAnomaly, elements.Eccentricity)

	trueAnomaly := 2 * math.Atan(math.Sqrt((1+elements.Eccentricity)/(1-elements.Eccentricity))*math.Tan(E*math.Pi/180))
	trueAnomaly = c.normalizeAngle(trueAnomaly * 180 / math.Pi)

	elements.TrueLongitude = c.normalizeAngle(trueAnomaly + elements.LongitudePeri)

	elements.RadiusVector = elements.SemimajorAxis * (1 - elements.Eccentricity*math.Cos(E*math.Pi/180))

	lambda := elements.TrueLongitude
	beta := math.Asin(math.Sin((lambda-elements.LongitudeNode)*math.Pi/180) * math.Sin(elements.Inclination*math.Pi/180))

	elements.EclipticLon = lambda
	elements.EclipticLat = beta * 180 / math.Pi

	return elements
}

func (c *PlanetCalculatorFixed) solveKeplerEquation(M, e float64) float64 {
	E := M
	for i := 0; i < 10; i++ {
		dE := (M - E + e*180/math.Pi*math.Sin(E*math.Pi/180)) / (1 - e*math.Cos(E*math.Pi/180))
		E += dE
		if math.Abs(dE) < 1e-8 {
			break
		}
	}
	return E
}

func (c *PlanetCalculatorFixed) calculateSunLongitude(jd float64) float64 {
	t := (jd - 2451545.0) / 36525.0

	L := 280.46646 + 36000.76983*t + 0.0003032*t*t
	M := 357.52911 + 35999.05029*t - 0.0001537*t*t

	L = c.normalizeAngle(L)
	M = c.normalizeAngle(M)

	C := (1.914602-0.004817*t-0.000014*t*t)*math.Sin(M*math.Pi/180) +
		(0.019993-0.000101*t)*math.Sin(2*M*math.Pi/180) +
		0.000289*math.Sin(3*M*math.Pi/180)

	sunLon := L + C
	return c.normalizeAngle(sunLon)
}

func (c *PlanetCalculatorFixed) calculateEquatorialCoordinates(elements *OrbitalElementsFixed, sunLon float64) (float64, float64) {
	lambda := elements.EclipticLon
	beta := elements.EclipticLat

	epsilon := 23.4392911

	lambdaRad := lambda * math.Pi / 180
	betaRad := beta * math.Pi / 180
	epsilonRad := epsilon * math.Pi / 180

	alpha := math.Atan2(
		math.Sin(lambdaRad)*math.Cos(epsilonRad)-math.Tan(betaRad)*math.Sin(epsilonRad),
		math.Cos(lambdaRad),
	)

	delta := math.Asin(
		math.Sin(betaRad)*math.Cos(epsilonRad) +
			math.Cos(betaRad)*math.Sin(epsilonRad)*math.Sin(lambdaRad),
	)

	raHours := alpha * 180 / math.Pi / 15

	if raHours < 0 {
		raHours += 24
	}
	if raHours >= 24 {
		raHours -= 24
	}

	decDeg := delta * 180 / math.Pi

	return raHours, decDeg
}

func (c *PlanetCalculatorFixed) calculateDistance(planetType PlanetType, elements *OrbitalElementsFixed) float64 {
	return elements.RadiusVector
}

func (c *PlanetCalculatorFixed) calculateElongation(elements *OrbitalElementsFixed, sunLon float64) float64 {
	planetLon := elements.EclipticLon

	elongation := math.Abs(planetLon - sunLon)
	if elongation > 180 {
		elongation = 360 - elongation
	}

	return elongation
}

func (c *PlanetCalculatorFixed) calculatePhase(elongation float64) float64 {
	phase := (1 + math.Cos(elongation*math.Pi/180)) / 2
	return phase
}

func (c *PlanetCalculatorFixed) calculateMagnitude(planetType PlanetType, distance, phase float64) float64 {
	baseMagnitude := 0.0
	phaseCoeff := 0.0

	switch planetType {
	case PlanetTypeMercury:
		baseMagnitude = -0.42
		phaseCoeff = 0.038
	case PlanetTypeVenus:
		baseMagnitude = -4.40
		phaseCoeff = 0.0009
	case PlanetTypeMars:
		baseMagnitude = -1.52
		phaseCoeff = 0.016
	case PlanetTypeJupiter:
		baseMagnitude = -9.40
		phaseCoeff = 0.005
	case PlanetTypeSaturn:
		baseMagnitude = -8.88
		phaseCoeff = 0.044
	case PlanetTypeUranus:
		baseMagnitude = -7.19
		phaseCoeff = 0.0028
	case PlanetTypeNeptune:
		baseMagnitude = -6.87
		phaseCoeff = 0.0
	}

	phaseAngle := math.Acos(2*phase-1) * 180 / math.Pi

	distanceAU := distance
	if distanceAU <= 0 {
		distanceAU = 1.0
	}

	magnitude := baseMagnitude + 5*math.Log10(distanceAU) + phaseCoeff*phaseAngle

	return magnitude
}

func (c *PlanetCalculatorFixed) normalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

func (c *PlanetCalculatorFixed) Description() string {
	return c.description
}
