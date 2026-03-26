package api

import (
	"fmt"
	"net/http"
	"scientific_calc/internal/calculator"
	"time"

	"github.com/gin-gonic/gin"
)

// CompareHandler 对比处理器
type CompareHandler struct {
	calculatorManager *calculator.CalculatorManager
}

// NewCompareHandler 创建新的对比处理器
func NewCompareHandler(calculatorManager *calculator.CalculatorManager) *CompareHandler {
	return &CompareHandler{
		calculatorManager: calculatorManager,
	}
}

// CompareRequest 对比请求
type CompareRequest struct {
	Calculation string      `json:"calculation" binding:"required"` // 计算类型
	Params      interface{} `json:"params"`                         // 计算参数
}

// CompareResponse 对比响应
type CompareResponse struct {
	Success     bool             `json:"success"`
	Original    interface{}      `json:"original"`    // 原始响应
	Fixed       interface{}      `json:"fixed"`       // 修复后响应
	Differences []DifferenceInfo `json:"differences"` // 差异列表
	Timestamp   string           `json:"timestamp"`
	SessionID   string           `json:"session_id,omitempty"`
}

// DifferenceInfo 差异信息
type DifferenceInfo struct {
	Field       string      `json:"field"`       // 差异字段
	Original    interface{} `json:"original"`    // 原始值
	Fixed       interface{} `json:"fixed"`       // 修复值
	Description string      `json:"description"` // 差异说明
	Severity    string      `json:"severity"`    // 严重程度: critical, major, minor
}

// CompareCalculation 对比计算接口
// @Summary 对比原始与修复后的计算结果
// @Description 接收相同的原始请求参数，返回原始缺陷响应与修复后响应的结构化对比结果
// @Tags 科学计算
// @Accept json
// @Produce json
// @Param request body CompareRequest true "对比请求参数"
// @Success 200 {object} CompareResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/solver/compare [post]
func (h *CompareHandler) CompareCalculation(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 只支持star类型的对比
	if req.Calculation != "star" {
		h.sendError(c, http.StatusBadRequest, "当前仅支持star类型的对比")
		return
	}

	sessionID := GenerateSessionID()

	// 检查计算器是否存在
	_, exists := h.calculatorManager.GetCalculator(calculator.CalculationTypeStar)
	if !exists {
		h.sendError(c, http.StatusInternalServerError, "原始计算器未找到")
		return
	}

	_, exists = h.calculatorManager.GetCalculator(calculator.CalculationTypeStarFixed)
	if !exists {
		h.sendError(c, http.StatusInternalServerError, "修复版计算器未找到")
		return
	}

	// 执行原始计算
	originalResult, _, err := h.calculatorManager.CalculateWithSession(calculator.CalculationTypeStar, req.Params, sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "原始计算失败: "+err.Error())
		return
	}

	// 执行修复版计算
	fixedResult, _, err := h.calculatorManager.CalculateWithSession(calculator.CalculationTypeStarFixed, req.Params, sessionID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "修复版计算失败: "+err.Error())
		return
	}

	// 计算差异
	differences := h.calculateDifferences(originalResult, fixedResult)

	response := CompareResponse{
		Success:     true,
		Original:    originalResult,
		Fixed:       fixedResult,
		Differences: differences,
		Timestamp:   time.Now().Format(time.RFC3339),
		SessionID:   sessionID,
	}

	c.JSON(http.StatusOK, response)
}

// calculateDifferences 计算差异
func (h *CompareHandler) calculateDifferences(original, fixed interface{}) []DifferenceInfo {
	var differences []DifferenceInfo

	// 类型断言
	origMap, ok1 := original.(*calculator.StarResult)
	fixedMap, ok2 := fixed.(*calculator.StarFixedResult)

	if !ok1 || !ok2 {
		// 尝试从map中获取
		return h.calculateDifferencesFromMap(original, fixed)
	}

	// 对比农历日期
	if origMap.LunarDate != fixedMap.LunarDate {
		differences = append(differences, DifferenceInfo{
			Field:       "lunar_date",
			Original:    origMap.LunarDate,
			Fixed:       fixedMap.LunarDate,
			Description: "修正农历日期格式：原格式错误使用公历月日，修正后使用正确的农历月日",
			Severity:    "critical",
		})
	}

	// 对比日干支
	if origMap.DayGanZhi != fixedMap.DayGanZhi {
		differences = append(differences, DifferenceInfo{
			Field:       "day_ganzhi",
			Original:    origMap.DayGanZhi,
			Fixed:       fixedMap.DayGanZhi,
			Description: "修正日干支计算：原算法基于简单哈希，修正后基于儒略日精确计算",
			Severity:    "critical",
		})
	}

	// 对比星宿
	if origMap.Constellation != fixedMap.Constellation {
		differences = append(differences, DifferenceInfo{
			Field:       "constellation",
			Original:    origMap.Constellation,
			Fixed:       fixedMap.Constellation,
			Description: "修正二十八宿计算：原算法基于年内天数取模，修正后基于节气周期计算",
			Severity:    "critical",
		})
	}

	// 对比星宿方位描述
	if origMap.StarPosition != fixedMap.StarPosition {
		differences = append(differences, DifferenceInfo{
			Field:       "star_position",
			Original:    origMap.StarPosition,
			Fixed:       fixedMap.StarPosition,
			Description: "修正星曜位置描述：原描述存在方位矛盾（如朱雀在西方），修正后确保四象与方位对应正确",
			Severity:    "major",
		})
	}

	// 对比儒略日
	if origMap.JulianDay != fixedMap.JulianDay {
		differences = append(differences, DifferenceInfo{
			Field:       "julian_day",
			Original:    origMap.JulianDay,
			Fixed:       fixedMap.JulianDay,
			Description: "修正儒略日计算：原计算公式有误，修正后使用标准儒略日公式",
			Severity:    "critical",
		})
	}

	// 对比日分值
	if origMap.DayScore != fixedMap.DayScore {
		differences = append(differences, DifferenceInfo{
			Field:       "day_score",
			Original:    origMap.DayScore,
			Fixed:       fixedMap.DayScore,
			Description: "修正日分值计算：原算法范围不合理，修正后确保在0-100范围内",
			Severity:    "minor",
		})
	}

	// 对比吉凶程度
	if origMap.AuspiciousLevel != fixedMap.AuspiciousLevel {
		differences = append(differences, DifferenceInfo{
			Field:       "auspicious_level",
			Original:    origMap.AuspiciousLevel,
			Fixed:       fixedMap.AuspiciousLevel,
			Description: "修正吉凶程度计算：原评分与吉凶状态不一致，修正后确保逻辑自洽",
			Severity:    "major",
		})
	}

	// 对比二十八宿索引
	if origMap.ConstellationIdx != fixedMap.ConstellationIdx {
		differences = append(differences, DifferenceInfo{
			Field:       "constellation_index",
			Original:    origMap.ConstellationIdx,
			Fixed:       fixedMap.ConstellationIdx,
			Description: "修正二十八宿索引计算：原算法使用年内天数取模，修正后基于节气周期计算",
			Severity:    "minor",
		})
	}

	// 新增：北斗七星支持检测
	if fixedMap.IsBigDipper {
		differences = append(differences, DifferenceInfo{
			Field:       "star_name",
			Original:    "无北斗七星信息",
			Fixed:       fixedMap.BigDipperInfo,
			Description: "新增北斗七星专属信息：原版不支持big_dipper参数，修复版提供完整的北斗七星天文数据",
			Severity:    "critical",
		})
	}

	// 新增：星宿分组和方位信息
	if fixedMap.ConstellationGroup != "" {
		differences = append(differences, DifferenceInfo{
			Field:       "constellation_group",
			Original:    "无分组信息",
			Fixed:       fixedMap.ConstellationGroup,
			Description: "新增星宿所属四象信息：提供东方青龙、北方玄武、西方白虎、南方朱雀的准确归属",
			Severity:    "major",
		})

		differences = append(differences, DifferenceInfo{
			Field:       "constellation_direction",
			Original:    "无方位信息",
			Fixed:       fixedMap.ConstellationDirection,
			Description: "新增星宿方位信息：修正后方位与四象正确对应（青龙-东，朱雀-南，白虎-西，玄武-北）",
			Severity:    "major",
		})
	}

	// 新增：年干支信息
	if fixedMap.YearGanZhi != "" {
		differences = append(differences, DifferenceInfo{
			Field:       "year_ganzhi",
			Original:    "无年干支",
			Fixed:       fixedMap.YearGanZhi,
			Description: "新增年干支信息：原版缺少年干支，修复版提供完整的年干支数据",
			Severity:    "minor",
		})
	}

	// 新增：修正说明
	if len(fixedMap.FixesApplied) > 0 {
		differences = append(differences, DifferenceInfo{
			Field:       "fixes_applied",
			Original:    "无修正记录",
			Fixed:       fixedMap.FixesApplied,
			Description: "新增修正说明：列出所有应用的修正项，便于追踪问题修复过程",
			Severity:    "minor",
		})
	}

	return differences
}

// calculateDifferencesFromMap 从map类型计算差异（备选方案）
func (h *CompareHandler) calculateDifferencesFromMap(original, fixed interface{}) []DifferenceInfo {
	var differences []DifferenceInfo

	origMap, ok1 := original.(map[string]interface{})
	fixedMap, ok2 := fixed.(map[string]interface{})

	if !ok1 || !ok2 {
		differences = append(differences, DifferenceInfo{
			Field:       "general",
			Original:    "无法解析原始响应",
			Fixed:       "无法解析修复后响应",
			Description: "无法进行对比分析",
			Severity:    "critical",
		})
		return differences
	}

	// 对比关键字段
	fields := []string{"lunar_date", "day_ganzhi", "constellation", "star_position", "julian_day", "day_score", "auspicious_level"}

	for _, field := range fields {
		origVal, origExists := origMap[field]
		fixedVal, fixedExists := fixedMap[field]

		if origExists && fixedExists && fmt.Sprintf("%v", origVal) != fmt.Sprintf("%v", fixedVal) {
			differences = append(differences, DifferenceInfo{
				Field:       field,
				Original:    origVal,
				Fixed:       fixedVal,
				Description: fmt.Sprintf("字段 %s 存在差异", field),
				Severity:    "major",
			})
		}
	}

	return differences
}

// sendError 发送错误响应
func (h *CompareHandler) sendError(c *gin.Context, code int, message string) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}

	c.JSON(code, response)
}
