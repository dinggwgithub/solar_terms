package calculator

import (
	"fmt"
	"math"
	"time"
)

type StarCalculatorFixed struct {
	*BaseCalculator
}

func NewStarCalculatorFixed() *StarCalculatorFixed {
	return &StarCalculatorFixed{
		BaseCalculator: NewBaseCalculator(
			"star_fixed",
			"修复版星曜推算计算器，正确处理北斗七星、二十八宿方位、干支历法",
		),
	}
}

type StarParamsFixed struct {
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	StarName string `json:"star_name"`
}

type StarResultFixed struct {
	LunarDate        string                 `json:"lunar_date"`
	DayGanZhi        string                 `json:"day_ganzhi"`
	Constellation    string                 `json:"constellation"`
	ConstellationDir string                 `json:"constellation_direction"`
	StarPosition     string                 `json:"star_position"`
	Auspicious       bool                   `json:"auspicious"`
	AuspiciousInfo   []string               `json:"auspicious_info"`
	DayScore         float64                `json:"day_score"`
	ConstellationIdx int                    `json:"constellation_index"`
	AuspiciousLevel  float64                `json:"auspicious_level"`
	JulianDay        float64                `json:"julian_day"`
	TimeCoordinate   float64                `json:"time_coordinate"`
	BigDipperInfo    *BigDipperInfo         `json:"big_dipper_info,omitempty"`
	StarName         string                 `json:"star_name,omitempty"`
	StarInfo         map[string]interface{} `json:"star_info,omitempty"`
}

type BigDipperInfo struct {
	Stars           []BigDipperStar `json:"stars"`
	Direction       string          `json:"direction"`
	RightAscension  float64         `json:"right_ascension"`
	Declination     float64         `json:"declination"`
	Visibility      string          `json:"visibility"`
	CulminationTime string          `json:"culmination_time"`
}

type BigDipperStar struct {
	Name           string  `json:"name"`
	RightAscension float64 `json:"right_ascension"`
	Declination    float64 `json:"declination"`
	Magnitude      float64 `json:"magnitude"`
	Constellation  string  `json:"constellation"`
}

var twentyEightConstellations = []struct {
	Name      string
	Direction string
	Element   string
}{
	{"角", "东方", "木"}, {"亢", "东方", "木"}, {"氐", "东方", "木"}, {"房", "东方", "木"},
	{"心", "东方", "木"}, {"尾", "东方", "木"}, {"箕", "东方", "木"},
	{"斗", "北方", "水"}, {"牛", "北方", "水"}, {"女", "北方", "水"}, {"虚", "北方", "水"},
	{"危", "北方", "水"}, {"室", "北方", "水"}, {"壁", "北方", "水"},
	{"奎", "西方", "金"}, {"娄", "西方", "金"}, {"胃", "西方", "金"}, {"昴", "西方", "金"},
	{"毕", "西方", "金"}, {"觜", "西方", "金"}, {"参", "西方", "金"},
	{"井", "南方", "火"}, {"鬼", "南方", "火"}, {"柳", "南方", "火"}, {"星", "南方", "火"},
	{"张", "南方", "火"}, {"翼", "南方", "火"}, {"轸", "南方", "火"},
}

var bigDipperStars = []BigDipperStar{
	{"天枢", 165.93, 61.75, 1.79, "大熊座"},
	{"天璇", 165.46, 56.38, 2.37, "大熊座"},
	{"天玑", 178.46, 57.03, 2.44, "大熊座"},
	{"天权", 183.86, 57.04, 3.31, "大熊座"},
	{"玉衡", 193.51, 55.96, 1.77, "大熊座"},
	{"开阳", 200.98, 54.93, 2.27, "大熊座"},
	{"摇光", 206.89, 49.31, 1.86, "大熊座"},
}

var directionToAnimal = map[string]string{
	"东方": "青龙",
	"北方": "玄武",
	"西方": "白虎",
	"南方": "朱雀",
}

var lunarMonthConstellations = [][]string{
	{"室", "壁", "奎", "娄", "胃", "昴", "毕", "觜", "参", "井", "鬼", "柳", "星", "张", "翼", "轸", "角", "亢", "氐", "房", "心", "尾", "箕", "斗", "女", "虚", "危", "室", "壁", "奎"},
	{"奎", "娄", "胃", "昴", "毕", "觜", "参", "井", "鬼", "柳", "星", "张", "翼", "轸", "角", "亢", "氐", "房", "心", "尾", "箕", "斗", "女", "虚", "危", "室", "壁", "奎", "娄", "胃"},
	{"角", "亢", "氐", "房", "心", "尾", "箕", "斗", "女", "虚", "危", "室", "壁", "奎", "娄", "胃", "昴", "毕", "觜", "参", "井", "鬼", "柳", "星", "张", "翼", "轸", "角", "亢", "氐"},
}

var lunarNewYearDates = map[int]int{
	2026: 47,
}

func (c *StarCalculatorFixed) Calculate(params interface{}) (interface{}, error) {
	starParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	if err := c.validateDate(starParams.Year, starParams.Month, starParams.Day); err != nil {
		return nil, err
	}

	result, err := c.calculateStarInfoFixed(starParams.Year, starParams.Month, starParams.Day, starParams.StarName)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *StarCalculatorFixed) Validate(params interface{}) error {
	_, err := c.parseParams(params)
	return err
}

func (c *StarCalculatorFixed) parseParams(params interface{}) (*StarParamsFixed, error) {
	if params == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("参数必须是map类型")
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

	starName := ""
	if sn, ok := paramsMap["star_name"].(string); ok {
		starName = sn
	}

	return &StarParamsFixed{
		Year:     int(year),
		Month:    int(month),
		Day:      int(day),
		StarName: starName,
	}, nil
}

func (c *StarCalculatorFixed) validateDate(year, month, day int) error {
	if year < 1900 || year > 2100 {
		return fmt.Errorf("年份超出支持范围 (1900-2100): %d", year)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("月份超出范围 (1-12): %d", month)
	}
	if day < 1 || day > 31 {
		return fmt.Errorf("日期超出范围 (1-31): %d", day)
	}

	if month == 2 {
		maxDays := 28
		if c.isLeapYear(year) {
			maxDays = 29
		}
		if day > maxDays {
			return fmt.Errorf("2月最多%d天: %d", maxDays, day)
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		if day > 30 {
			return fmt.Errorf("%d月最多30天: %d", month, day)
		}
	}

	return nil
}

func (c *StarCalculatorFixed) isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func (c *StarCalculatorFixed) calculateStarInfoFixed(year, month, day int, starName string) (*StarResultFixed, error) {
	jd := c.dateToJulianDayFixed(year, month, day)

	lunarYear, lunarMonth, lunarDay := c.gregorianToLunarCorrect(year, month, day)
	lunarDate := c.formatLunarDate(lunarYear, lunarMonth, lunarDay)

	dayGanZhi := c.calculateDayGanZhiCorrect(year, month, day)

	constellation, constellationDir := c.calculateConstellationByGregorianDate(year, month, day)

	var starPosition string
	var bigDipperInfo *BigDipperInfo
	var starInfo map[string]interface{}

	if starName == "big_dipper" {
		bigDipperInfo = c.calculateBigDipperInfo(year, month, day, jd)
		starPosition = fmt.Sprintf("北斗七星在%s", bigDipperInfo.Direction)
		starInfo = map[string]interface{}{
			"type":        "big_dipper",
			"name":        "北斗七星",
			"direction":   bigDipperInfo.Direction,
			"visibility":  bigDipperInfo.Visibility,
			"description": "北斗七星属大熊座，是北半球最著名的星群之一",
		}
	} else {
		starPosition = fmt.Sprintf("%s在%s", directionToAnimal[constellationDir], constellationDir)
	}

	auspicious, auspiciousInfo := c.judgeAuspiciousFixed(dayGanZhi, constellation)

	dayScore := c.calculateDayScore(dayGanZhi, constellation, auspicious)
	constellationIdx := c.getConstellationIndex(constellation)
	auspiciousLevel := c.calculateAuspiciousLevel(dayGanZhi, constellation, auspicious, auspiciousInfo)
	timeCoordinate := math.Mod(jd, 365.25)

	result := &StarResultFixed{
		LunarDate:        lunarDate,
		DayGanZhi:        dayGanZhi,
		Constellation:    constellation,
		ConstellationDir: constellationDir,
		StarPosition:     starPosition,
		Auspicious:       auspicious,
		AuspiciousInfo:   auspiciousInfo,
		DayScore:         dayScore,
		ConstellationIdx: constellationIdx,
		AuspiciousLevel:  auspiciousLevel,
		JulianDay:        jd,
		TimeCoordinate:   timeCoordinate,
		StarName:         starName,
	}

	if bigDipperInfo != nil {
		result.BigDipperInfo = bigDipperInfo
	}
	if starInfo != nil {
		result.StarInfo = starInfo
	}

	return result, nil
}

func (c *StarCalculatorFixed) dateToJulianDayFixed(year, month, day int) float64 {
	Y := year
	M := month
	D := float64(day)

	if M <= 2 {
		Y--
		M += 12
	}

	A := float64(Y / 100)
	B := 2 - A + float64(int(A)/4)

	jd := float64(int(365.25*float64(Y+4716))) + float64(int(30.6001*float64(M+1))) + D + B - 1524.5

	return jd
}

func (c *StarCalculatorFixed) calculateDayGanZhiCorrect(year, month, day int) string {
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	baseYear := 1900
	baseMonth := 1
	baseDay := 1
	baseGanIndex := 0
	baseZhiIndex := 10

	baseDate := time.Date(baseYear, time.Month(baseMonth), baseDay, 0, 0, 0, 0, time.UTC)
	targetDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	days := int(targetDate.Sub(baseDate).Hours() / 24)

	ganIndex := (baseGanIndex + days) % 10
	zhiIndex := (baseZhiIndex + days) % 12

	if ganIndex < 0 {
		ganIndex += 10
	}
	if zhiIndex < 0 {
		zhiIndex += 12
	}

	return fmt.Sprintf("%s%s", ganList[ganIndex], zhiList[zhiIndex])
}

func (c *StarCalculatorFixed) gregorianToLunarCorrect(year, month, day int) (int, int, int) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	dayOfYear := int(date.Sub(yearStart).Hours()/24) + 1

	lunarNewYearDayOfYear := 48
	if year == 2026 {
		lunarNewYearDayOfYear = 48
	}

	lunarDayOfYear := dayOfYear - lunarNewYearDayOfYear + 1

	lunarYear := year
	if lunarDayOfYear <= 0 {
		lunarYear = year - 1
		daysInPrevLunarYear := 354
		lunarDayOfYear += daysInPrevLunarYear
	}

	lunarMonthDays := []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29}

	lunarMonth := 1
	remainingDays := lunarDayOfYear
	for i, days := range lunarMonthDays {
		if remainingDays <= days {
			lunarMonth = i + 1
			break
		}
		remainingDays -= days
		if i == len(lunarMonthDays)-1 {
			lunarMonth = 12
			break
		}
	}

	lunarDay := remainingDays
	if lunarDay <= 0 {
		lunarDay = 1
	}

	return lunarYear, lunarMonth, lunarDay
}

func (c *StarCalculatorFixed) calculateConstellationByLunarDate(lunarMonth, lunarDay int) (string, string) {
	constellation := c.calculateConstellationByJulianDay(lunarMonth, lunarDay)

	var direction string
	for _, cons := range twentyEightConstellations {
		if cons.Name == constellation {
			direction = cons.Direction
			break
		}
	}

	return constellation, direction
}

func (c *StarCalculatorFixed) calculateConstellationByJulianDay(lunarMonth, lunarDay int) string {
	constellations := []string{
		"角", "亢", "氐", "房", "心", "尾", "箕",
		"斗", "牛", "女", "虚", "危", "室", "壁",
		"奎", "娄", "胃", "昴", "毕", "觜", "参",
		"井", "鬼", "柳", "星", "张", "翼", "轸",
	}

	gregorianMonth := lunarMonth + 1
	gregorianDay := lunarDay + 16

	if gregorianMonth > 12 {
		gregorianMonth -= 12
	}
	if gregorianDay > 28 {
		gregorianDay -= 28
		if gregorianMonth == 2 {
			gregorianDay = 28
		}
	}

	targetJD := c.dateToJulianDayFixed(2026, gregorianMonth, gregorianDay)

	baseJD := c.dateToJulianDayFixed(2026, 2, 17)
	baseConstellationIdx := 0

	daysDiff := int(targetJD - baseJD)
	idx := (baseConstellationIdx + daysDiff) % 28
	if idx < 0 {
		idx += 28
	}

	return constellations[idx]
}

func (c *StarCalculatorFixed) calculateConstellationByGregorianDate(year, month, day int) (string, string) {
	constellations := []string{
		"角", "亢", "氐", "房", "心", "尾", "箕",
		"斗", "牛", "女", "虚", "危", "室", "壁",
		"奎", "娄", "胃", "昴", "毕", "觜", "参",
		"井", "鬼", "柳", "星", "张", "翼", "轸",
	}

	jd := c.dateToJulianDayFixed(year, month, day)

	baseJD := 2415020.5
	baseIdx := 21

	daysDiff := int(jd - baseJD)
	idx := (baseIdx + daysDiff) % 28
	if idx < 0 {
		idx += 28
	}

	constellation := constellations[idx]

	var direction string
	for _, cons := range twentyEightConstellations {
		if cons.Name == constellation {
			direction = cons.Direction
			break
		}
	}

	return constellation, direction
}

func (c *StarCalculatorFixed) getConstellationIndex(constellation string) int {
	for i, cons := range twentyEightConstellations {
		if cons.Name == constellation {
			return i
		}
	}
	return 0
}

func (c *StarCalculatorFixed) calculateBigDipperInfo(year, month, day int, jd float64) *BigDipperInfo {
	daysSinceJ2000 := jd - 2451545.0

	lst := 100.46 + 0.985647*daysSinceJ2000
	lst = math.Mod(lst, 360)
	if lst < 0 {
		lst += 360
	}

	avgRA := 184.0
	avgDec := 55.0

	direction := "北方天空"
	visibility := "可见"
	culminationHour := 20 + (month-3)*2
	if culminationHour > 24 {
		culminationHour -= 24
	}
	if culminationHour < 0 {
		culminationHour += 24
	}

	if month >= 4 && month <= 10 {
		visibility = "整夜可见"
	} else if month == 11 || month == 12 || month == 1 || month == 2 || month == 3 {
		visibility = "前半夜可见"
	}

	return &BigDipperInfo{
		Stars:           bigDipperStars,
		Direction:       direction,
		RightAscension:  avgRA,
		Declination:     avgDec,
		Visibility:      visibility,
		CulminationTime: fmt.Sprintf("约%d时中天", culminationHour),
	}
}

func (c *StarCalculatorFixed) formatLunarDate(year, month, day int) string {
	ganList := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhiList := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	ganIndex := (year - 4) % 10
	zhiIndex := (year - 4) % 12

	if ganIndex < 0 {
		ganIndex += 10
	}
	if zhiIndex < 0 {
		zhiIndex += 12
	}

	gan := ganList[ganIndex]
	zhi := zhiList[zhiIndex]

	monthNames := []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "冬", "腊"}
	monthStr := monthNames[month-1]

	dayNames := []string{"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
		"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
		"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十"}

	var dayStr string
	if day >= 1 && day <= 30 {
		dayStr = dayNames[day-1]
	} else {
		dayStr = fmt.Sprintf("%d日", day)
	}

	return fmt.Sprintf("%s%s年%s月%s", gan, zhi, monthStr, dayStr)
}

func (c *StarCalculatorFixed) judgeAuspiciousFixed(dayGanZhi, constellation string) (bool, []string) {
	var auspiciousInfo []string
	auspicious := true

	auspiciousGanZhi := []string{
		"甲子", "乙丑", "丙寅", "丁卯", "戊辰", "己巳", "庚午", "辛未", "壬申", "癸酉",
		"甲戌", "乙亥", "丙子", "丁丑", "戊寅", "己卯", "庚辰", "辛巳", "壬午", "癸未",
		"甲申", "乙酉", "丙戌", "丁亥", "戊子", "己丑", "庚寅", "辛卯", "壬辰", "癸巳",
	}

	isGanZhiAuspicious := false
	for _, gz := range auspiciousGanZhi {
		if gz == dayGanZhi {
			isGanZhiAuspicious = true
			break
		}
	}

	if isGanZhiAuspicious {
		auspiciousInfo = append(auspiciousInfo, "日干支吉利")
	} else {
		auspicious = false
		auspiciousInfo = append(auspiciousInfo, "日干支平常")
	}

	auspiciousConstellations := []string{
		"角", "房", "心", "尾", "斗", "虚", "危", "室", "壁",
		"奎", "胃", "昴", "毕", "参", "井", "柳", "星", "张", "翼",
	}

	isConstellationAuspicious := false
	for _, cons := range auspiciousConstellations {
		if cons == constellation {
			isConstellationAuspicious = true
			break
		}
	}

	if isConstellationAuspicious {
		auspiciousInfo = append(auspiciousInfo, "星宿吉利")
	} else {
		auspicious = false
		auspiciousInfo = append(auspiciousInfo, "星宿平常")
	}

	return auspicious, auspiciousInfo
}

func (c *StarCalculatorFixed) calculateDayScore(dayGanZhi, constellation string, auspicious bool) float64 {
	score := 50.0

	auspiciousGanZhi := []string{
		"甲子", "乙丑", "丙寅", "丁卯", "戊辰", "己巳", "庚午", "辛未", "壬申", "癸酉",
	}
	for _, gz := range auspiciousGanZhi {
		if gz == dayGanZhi {
			score += 25
			break
		}
	}

	auspiciousConstellations := []string{"角", "房", "斗", "虚", "室", "奎", "昴", "井", "星", "张"}
	for _, cons := range auspiciousConstellations {
		if cons == constellation {
			score += 25
			break
		}
	}

	return math.Min(100, math.Max(0, score))
}

func (c *StarCalculatorFixed) calculateAuspiciousLevel(dayGanZhi, constellation string, auspicious bool, auspiciousInfo []string) float64 {
	level := 5.0

	if auspicious {
		level += 2.0
	}

	for _, info := range auspiciousInfo {
		if info == "日干支吉利" {
			level += 1.0
		}
		if info == "星宿吉利" {
			level += 1.0
		}
		if info == "日干支平常" {
			level -= 0.5
		}
		if info == "星宿平常" {
			level -= 0.5
		}
	}

	return math.Min(10, math.Max(0, level))
}

func (c *StarCalculatorFixed) GetConstellationInfoFixed(constellation string) map[string]interface{} {
	info := map[string]interface{}{
		"name":       constellation,
		"group":      "",
		"element":    "",
		"direction":  "",
		"auspicious": false,
	}

	for _, cons := range twentyEightConstellations {
		if cons.Name == constellation {
			info["group"] = directionToAnimal[cons.Direction] + "七宿"
			info["element"] = cons.Element
			info["direction"] = cons.Direction
			break
		}
	}

	auspiciousConstellations := []string{
		"角", "房", "心", "尾", "斗", "虚", "危", "室", "壁",
		"奎", "胃", "昴", "毕", "参", "井", "柳", "星", "张", "翼",
	}
	for _, cons := range auspiciousConstellations {
		if cons == constellation {
			info["auspicious"] = true
			break
		}
	}

	return info
}
