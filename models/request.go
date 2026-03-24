package models

// CalculationRequest 计算请求结构体
// 支持两种参数格式：
// 格式1: {"calculation": "solar_term", "params": {"year": 2024, "term_index": 2}}
// 格式2: {"calculation": "solar_term", "year": 2024, "term_index": 2}
type CalculationRequest struct {
	Calculation string      `json:"calculation" binding:"required"` // 计算类型
	Params      interface{} `json:"params"`                         // 计算参数（格式1）

	// 直接参数（格式2）
	Year            int     `json:"year"`             // 年份
	Month           int     `json:"month"`            // 月份
	Day             int     `json:"day"`              // 日期
	Hour            int     `json:"hour"`             // 小时
	Minute          int     `json:"minute"`           // 分钟
	Second          int     `json:"second"`           // 秒
	TermIndex       int     `json:"term_index"`       // 节气索引 (0-23)
	SolarDate       string  `json:"solar_date"`       // 阳历日期
	DaysFromTerm    int     `json:"days_from_term"`   // 距节气天数
	TargetLongitude float64 `json:"target_longitude"` // 目标黄经
}

// GetParams 获取计算参数（支持两种格式）
func (r *CalculationRequest) GetParams() interface{} {
	// 如果提供了Params字段，优先使用
	if r.Params != nil {
		return r.Params
	}

	// 否则使用直接参数构建参数map
	params := make(map[string]interface{})

	if r.Year != 0 {
		params["year"] = r.Year
	}
	if r.Month != 0 {
		params["month"] = r.Month
	}
	if r.Day != 0 {
		params["day"] = r.Day
	}
	if r.Hour != 0 {
		params["hour"] = r.Hour
	}
	if r.Minute != 0 {
		params["minute"] = r.Minute
	}
	if r.Second != 0 {
		params["second"] = r.Second
	}
	if r.TermIndex != 0 {
		params["term_index"] = r.TermIndex
	}
	if r.SolarDate != "" {
		params["solar_date"] = r.SolarDate
	}
	if r.DaysFromTerm != 0 {
		params["days_from_term"] = r.DaysFromTerm
	}
	if r.TargetLongitude != 0 {
		params["target_longitude"] = r.TargetLongitude
	}

	return params
}

// CalculationResult 计算结果结构体
type CalculationResult struct {
	SolarTermTime  string       `json:"solar_term_time"` // 节气精确时间
	SunLongitude   float64      `json:"sun_longitude"`   // 太阳黄经
	JulianDate     float64      `json:"julian_date"`     // 儒略日
	GanZhi         GanZhiResult `json:"gan_zhi"`         // 干支结果
	StartingAge    string       `json:"starting_age"`    // 起运岁数
	TermDate       int          `json:"term_date"`       // 节气日期
	Iterations     int          `json:"iterations"`      // 迭代次数
	Converged      bool         `json:"converged"`       // 是否收敛
	PrecisionError float64      `json:"precision_error"` // 精度误差
}

// GanZhiResult 干支计算结果
type GanZhiResult struct {
	GanYear  string `json:"gan_year"`  // 年天干
	ZhiYear  string `json:"zhi_year"`  // 年地支
	GanMonth string `json:"gan_month"` // 月天干
	ZhiMonth string `json:"zhi_month"` // 月地支
	GanDay   string `json:"gan_day"`   // 日天干
	ZhiDay   string `json:"zhi_day"`   // 日地支
	GanTime  string `json:"gan_time"`  // 时天干
	ZhiTime  string `json:"zhi_time"`  // 时地支
}

// AstronomyResult 天文计算结果
type AstronomyResult struct {
	SunLongitude      float64 `json:"sun_longitude"`      // 太阳黄经
	JulianDate        float64 `json:"julian_date"`        // 儒略日
	ApparentLongitude float64 `json:"apparent_longitude"` // 视黄经
	TrueLongitude     float64 `json:"true_longitude"`     // 真黄经
	MeanLongitude     float64 `json:"mean_longitude"`     // 平黄经
	MeanAnomaly       float64 `json:"mean_anomaly"`       // 平近点角
	EquationOfCenter  float64 `json:"equation_of_center"` // 中心差
	Nutation          float64 `json:"nutation"`           // 章动
}

// LunarDate 农历日期结构
type LunarDate struct {
	LunarYear   int    `json:"lunar_year"`   // 农历年
	LunarMonth  int    `json:"lunar_month"`  // 农历月（1-12）
	LunarDay    int    `json:"lunar_day"`    // 农历日（1-30）
	IsLeap      bool   `json:"is_leap"`      // 是否为闰月
	LunarString string `json:"lunar_string"` // 农历日期字符串
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

// StarConstellation 星曜推算结果
type StarConstellation struct {
	Constellation  string  `json:"constellation"`   // 星座
	RightAscension float64 `json:"right_ascension"` // 赤经
	Declination    float64 `json:"declination"`     // 赤纬
	Magnitude      float64 `json:"magnitude"`       // 星等
	Visibility     string  `json:"visibility"`      // 可见性
}

// SunriseSunsetResult 日出日落时间结果
type SunriseSunsetResult struct {
	Sunrise            string `json:"sunrise"`              // 日出时间
	Sunset             string `json:"sunset"`               // 日落时间
	DayLength          string `json:"day_length"`           // 白昼时长
	SolarNoon          string `json:"solar_noon"`           // 正午时间
	CivilTwilightBegin string `json:"civil_twilight_begin"` // 民用晨光开始
	CivilTwilightEnd   string `json:"civil_twilight_end"`   // 民用晨光结束
}

// MoonPhaseResult 月相计算结果
type MoonPhaseResult struct {
	Phase        string  `json:"phase"`          // 月相类型
	Illumination float64 `json:"illumination"`   // 光照比例
	Age          float64 `json:"age"`            // 月龄（天）
	NextNewMoon  string  `json:"next_new_moon"`  // 下次新月时间
	NextFullMoon string  `json:"next_full_moon"` // 下次满月时间
}

// StartingAgeResult 起运岁数结果
type StartingAgeResult struct {
	StartingAge     string                   `json:"starting_age"`     // 起运年龄
	Gender          string                   `json:"gender"`           // 性别
	BirthBazi       string                   `json:"birth_bazi"`       // 生辰八字
	MajorCycles     []map[string]interface{} `json:"major_cycles"`     // 大运周期
	CalculationDate string                   `json:"calculation_date"` // 计算时间
}
