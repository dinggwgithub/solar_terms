package models

// CalculationRequest 计算请求结构体
// 支持两种参数格式：
// 格式1: {"calculation": "solar_term", "params": {"year": 2024, "term_index": 2}}
// 格式2: {"calculation": "solar_term", "year": 2024, "term_index": 2}
type CalculationRequest struct {
	Calculation string      `json:"calculation" binding:"required" example:"equation_solver"` // 计算类型
	Params      interface{} `json:"params"`                                                   // 计算参数（格式1）

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
	StartingAge     string                   `json:"starting_age"`     // 起运年齢
	Gender          string                   `json:"gender"`           // 性别
	BirthBazi       string                   `json:"birth_bazi"`       // 生辰八字
	MajorCycles     []map[string]interface{} `json:"major_cycles"`     // 大运周期
	CalculationDate string                   `json:"calculation_date"` // 计算时间
}

// EquationSolverParams 方程求解参数
type EquationSolverParams struct {
	EquationType  string    `json:"equation_type" example:"ode"`  // 方程类型：nonlinear, linear, ode
	Equation      string    `json:"equation" example:"dy/dt = -y"` // 方程表达式
	InitialValue  float64   `json:"initial_value" example:"1.0"`   // 初始值
	InitialGuess  float64   `json:"initial_guess"`                  // 初始猜测值（非线性方程）
	Tolerance     float64   `json:"tolerance" example:"1e-6"`      // 容差
	MaxIterations int       `json:"max_iterations" example:"100"`  // 最大迭代次数
	Coefficients  []float64 `json:"coefficients"`                   // 系数（线性方程组）
	TimeStep      float64   `json:"time_step" example:"0.1"`       // 时间步长（微分方程）
	TimeRange     float64   `json:"time_range" example:"1.0"`      // 时间范围（微分方程）
	Method        string    `json:"method" example:"euler"`        // 求解方法：euler, rk4, rk45
}

// EquationSolverResult 方程求解结果
type EquationSolverResult struct {
	Solution       interface{} `json:"solution"`                  // 解
	Iterations     int         `json:"iterations"`                // 迭代次数
	Converged      bool        `json:"converged"`                 // 是否收敛
	Error          float64     `json:"error"`                     // 误差
	ErrorEstimate  float64     `json:"error_estimate"`            // 误差估计
	FunctionValue  float64     `json:"function_value"`            // 函数值
	Jacobian       [][]float64 `json:"jacobian,omitempty"`        // 雅可比矩阵
	TimePoints     []float64   `json:"time_points,omitempty"`     // 时间点
	SolutionPath   []float64   `json:"solution_path,omitempty"`   // 解路径
	MethodUsed     string      `json:"method_used,omitempty"`     // 使用的方法
	Stability      string      `json:"stability,omitempty"`       // 稳定性
	GlobalError    float64     `json:"global_error,omitempty"`    // 全局误差
	LocalError     float64     `json:"local_error,omitempty"`     // 局部误差
	Analytical     float64     `json:"analytical,omitempty"`      // 解析解
	AbsoluteError  float64     `json:"absolute_error,omitempty"`  // 绝对误差
	RelativeError  float64     `json:"relative_error,omitempty"`  // 相对误差
}
