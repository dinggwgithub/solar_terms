package models

// CalculationRequest 璁＄畻璇锋眰缁撴瀯浣?
// 鏀寔涓ょ鍙傛暟鏍煎紡锛?
// 鏍煎紡1: {"calculation": "solar_term", "params": {"year": 2024, "term_index": 2}}
// 鏍煎紡2: {"calculation": "solar_term", "year": 2024, "term_index": 2}
type CalculationRequest struct {
	Calculation string      `json:"calculation" binding:"required"` // 璁＄畻绫诲瀷
	Params      interface{} `json:"params"`                         // 璁＄畻鍙傛暟锛堟牸寮?锛?

	// 鐩存帴鍙傛暟锛堟牸寮?锛?
	Year            int     `json:"year"`             // 骞翠唤
	Month           int     `json:"month"`            // 鏈堜唤
	Day             int     `json:"day"`              // 鏃ユ湡
	Hour            int     `json:"hour"`             // 灏忔椂
	Minute          int     `json:"minute"`           // 鍒嗛挓
	Second          int     `json:"second"`           // 绉?
	TermIndex       int     `json:"term_index"`       // 鑺傛皵绱㈠紩 (0-23)
	SolarDate       string  `json:"solar_date"`       // 闃冲巻鏃ユ湡
	DaysFromTerm    int     `json:"days_from_term"`   // 璺濊妭姘斿ぉ鏁?
	TargetLongitude float64 `json:"target_longitude"` // 鐩爣榛勭粡
}

// GetParams 鑾峰彇璁＄畻鍙傛暟锛堟敮鎸佷袱绉嶆牸寮忥級
func (r *CalculationRequest) GetParams() interface{} {
	// 濡傛灉鎻愪緵浜哖arams瀛楁锛屼紭鍏堜娇鐢?
	if r.Params != nil {
		return r.Params
	}

	// 鍚﹀垯浣跨敤鐩存帴鍙傛暟鏋勫缓鍙傛暟map
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

// CalculationResult 璁＄畻缁撴灉缁撴瀯浣?
type CalculationResult struct {
	SolarTermTime  string       `json:"solar_term_time"` // 鑺傛皵绮剧‘鏃堕棿
	SunLongitude   float64      `json:"sun_longitude"`   // 澶槼榛勭粡
	JulianDate     float64      `json:"julian_date"`     // 鍎掔暐鏃?
	GanZhi         GanZhiResult `json:"gan_zhi"`         // 骞叉敮缁撴灉
	StartingAge    string       `json:"starting_age"`    // 璧疯繍宀佹暟
	TermDate       int          `json:"term_date"`       // 鑺傛皵鏃ユ湡
	Iterations     int          `json:"iterations"`      // 杩唬娆℃暟
	Converged      bool         `json:"converged"`       // 鏄惁鏀舵暃
	PrecisionError float64      `json:"precision_error"` // 绮惧害璇樊
}

// GanZhiResult 骞叉敮璁＄畻缁撴灉
type GanZhiResult struct {
	GanYear  string `json:"gan_year"`  // 骞村ぉ骞?
	ZhiYear  string `json:"zhi_year"`  // 骞村湴鏀?
	GanMonth string `json:"gan_month"` // 鏈堝ぉ骞?
	ZhiMonth string `json:"zhi_month"` // 鏈堝湴鏀?
	GanDay   string `json:"gan_day"`   // 鏃ュぉ骞?
	ZhiDay   string `json:"zhi_day"`   // 鏃ュ湴鏀?
	GanTime  string `json:"gan_time"`  // 鏃跺ぉ骞?
	ZhiTime  string `json:"zhi_time"`  // 鏃跺湴鏀?
}

// AstronomyResult 澶╂枃璁＄畻缁撴灉
type AstronomyResult struct {
	SunLongitude      float64 `json:"sun_longitude"`      // 澶槼榛勭粡
	JulianDate        float64 `json:"julian_date"`        // 鍎掔暐鏃?
	ApparentLongitude float64 `json:"apparent_longitude"` // 瑙嗛粍缁?
	TrueLongitude     float64 `json:"true_longitude"`     // 鐪熼粍缁?
	MeanLongitude     float64 `json:"mean_longitude"`     // 骞抽粍缁?
	MeanAnomaly       float64 `json:"mean_anomaly"`       // 骞宠繎鐐硅
	EquationOfCenter  float64 `json:"equation_of_center"` // 涓績宸?
	Nutation          float64 `json:"nutation"`           // 绔犲姩
}

// LunarDate 鍐滃巻鏃ユ湡缁撴瀯
type LunarDate struct {
	LunarYear   int    `json:"lunar_year"`   // 鍐滃巻骞?
	LunarMonth  int    `json:"lunar_month"`  // 鍐滃巻鏈堬紙1-12锛?
	LunarDay    int    `json:"lunar_day"`    // 鍐滃巻鏃ワ紙1-30锛?
	IsLeap      bool   `json:"is_leap"`      // 鏄惁涓洪棸鏈?
	LunarString string `json:"lunar_string"` // 鍐滃巻鏃ユ湡瀛楃涓?
}

// PlanetPosition 琛屾槦浣嶇疆缁撴灉
type PlanetPosition struct {
	RightAscension float64 `json:"right_ascension"` // 璧ょ粡锛堝皬鏃讹級
	Declination    float64 `json:"declination"`     // 璧ょ含锛堝害锛?
	Distance       float64 `json:"distance"`        // 璺濈锛堝ぉ鏂囧崟浣嶏級
	Magnitude      float64 `json:"magnitude"`       // 鏄熺瓑
	Phase          float64 `json:"phase"`           // 鐩镐綅锛?-1锛?
	Elongation     float64 `json:"elongation"`      // 璺濊锛堝害锛?
}

// StarConstellation 鏄熸洔鎺ㄧ畻缁撴灉
type StarConstellation struct {
	Constellation  string  `json:"constellation"`   // 鏄熷骇
	RightAscension float64 `json:"right_ascension"` // 璧ょ粡
	Declination    float64 `json:"declination"`     // 璧ょ含
	Magnitude      float64 `json:"magnitude"`       // 鏄熺瓑
	Visibility     string  `json:"visibility"`      // 鍙鎬?
}

// SunriseSunsetResult 鏃ュ嚭鏃ヨ惤鏃堕棿缁撴灉
type SunriseSunsetResult struct {
	Sunrise            string `json:"sunrise"`              // 鏃ュ嚭鏃堕棿
	Sunset             string `json:"sunset"`               // 鏃ヨ惤鏃堕棿
	DayLength          string `json:"day_length"`           // 鐧芥樇鏃堕暱
	SolarNoon          string `json:"solar_noon"`           // 姝ｅ崍鏃堕棿
	CivilTwilightBegin string `json:"civil_twilight_begin"` // 姘戠敤鏅ㄥ厜寮€濮?
	CivilTwilightEnd   string `json:"civil_twilight_end"`   // 姘戠敤鏅ㄥ厜缁撴潫
}

// MoonPhaseResult 鏈堢浉璁＄畻缁撴灉
type MoonPhaseResult struct {
	Phase        string  `json:"phase"`          // 鏈堢浉绫诲瀷
	Illumination float64 `json:"illumination"`   // 鍏夌収姣斾緥
	Age          float64 `json:"age"`            // 鏈堥緞锛堝ぉ锛?
	NextNewMoon  string  `json:"next_new_moon"`  // 涓嬫鏂版湀鏃堕棿
	NextFullMoon string  `json:"next_full_moon"` // 涓嬫婊℃湀鏃堕棿
}

// StartingAgeResult 璧疯繍宀佹暟缁撴灉
type StartingAgeResult struct {
	StartingAge     string                   `json:"starting_age"`     // 璧疯繍骞撮緞
	Gender          string                   `json:"gender"`           // 鎬у埆
	BirthBazi       string                   `json:"birth_bazi"`       // 鐢熻景鍏瓧
	MajorCycles     []map[string]interface{} `json:"major_cycles"`     // 澶ц繍鍛ㄦ湡
	CalculationDate string                   `json:"calculation_date"` // 璁＄畻鏃堕棿
}
