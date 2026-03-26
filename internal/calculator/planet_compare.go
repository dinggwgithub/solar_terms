package calculator

import (
	"fmt"
)

// PlanetComparisonResult 行星计算对比结果
type PlanetComparisonResult struct {
	Original *PlanetPosition       `json:"original"` // 原始计算结果
	Fixed    *PlanetPositionFixed  `json:"fixed"`    // 修复后结果
	Diff     *PlanetPositionDiff   `json:"diff"`     // 差异
	Analysis string                `json:"analysis"` // 分析说明
}

// PlanetPositionDiff 位置差异
type PlanetPositionDiff struct {
	RightAscensionDiff float64 `json:"right_ascension_diff"` // 赤经差异（小时）
	DeclinationDiff    float64 `json:"declination_diff"`     // 赤纬差异（度）
	DistanceDiff       float64 `json:"distance_diff"`        // 距离差异（AU）
	MagnitudeDiff      float64 `json:"magnitude_diff"`       // 星等差异
	PhaseDiff          float64 `json:"phase_diff"`           // 相位差异
	ElongationDiff     float64 `json:"elongation_diff"`      // 距角差异
}

// PlanetCalculatorCompare 行星计算对比计算器
type PlanetCalculatorCompare struct {
	*BaseCalculator
	originalCalc *PlanetCalculator
	fixedCalc    *PlanetCalculatorFixed
}

// NewPlanetCalculatorCompare 创建新的对比计算器
func NewPlanetCalculatorCompare() *PlanetCalculatorCompare {
	return &PlanetCalculatorCompare{
		BaseCalculator: NewBaseCalculator(
			"planet_compare",
			"行星位置计算对比工具，对比原始版和修复版的结果差异",
		),
		originalCalc: NewPlanetCalculator(),
		fixedCalc:    NewPlanetCalculatorFixed(),
	}
}

// Calculate 执行对比计算
func (c *PlanetCalculatorCompare) Calculate(params interface{}) (interface{}, error) {
	// 解析参数
	planetParams, err := c.parseParams(params)
	if err != nil {
		return nil, err
	}

	// 验证参数
	if err := c.validateParams(planetParams); err != nil {
		return nil, err
	}

	// 执行原始计算
	originalResult, err := c.originalCalc.Calculate(params)
	if err != nil {
		return nil, fmt.Errorf("原始计算失败: %v", err)
	}

	// 执行修复后计算
	fixedResult, err := c.fixedCalc.Calculate(params)
	if err != nil {
		return nil, fmt.Errorf("修复后计算失败: %v", err)
	}

	// 类型转换
	originalPos, ok := originalResult.(*PlanetPosition)
	if !ok {
		return nil, fmt.Errorf("原始结果类型错误")
	}

	fixedPos, ok := fixedResult.(*PlanetPositionFixed)
	if !ok {
		return nil, fmt.Errorf("修复后结果类型错误")
	}

	// 计算差异
	diff := c.calculateDiff(originalPos, fixedPos)

	// 生成分析
	analysis := c.generateAnalysis(originalPos, fixedPos, diff, planetParams.Planet)

	return &PlanetComparisonResult{
		Original: originalPos,
		Fixed:    fixedPos,
		Diff:     diff,
		Analysis: analysis,
	}, nil
}

// Validate 验证输入参数
func (c *PlanetCalculatorCompare) Validate(params interface{}) error {
	return c.originalCalc.Validate(params)
}

// parseParams 解析参数
func (c *PlanetCalculatorCompare) parseParams(params interface{}) (*PlanetParams, error) {
	return c.fixedCalc.parseParams(params)
}

// validateParams 验证参数
func (c *PlanetCalculatorCompare) validateParams(params *PlanetParams) error {
	return c.fixedCalc.validateParams(params)
}

// calculateDiff 计算差异
func (c *PlanetCalculatorCompare) calculateDiff(original *PlanetPosition, fixed *PlanetPositionFixed) *PlanetPositionDiff {
	return &PlanetPositionDiff{
		RightAscensionDiff: fixed.RightAscension - original.RightAscension,
		DeclinationDiff:    fixed.Declination - original.Declination,
		DistanceDiff:       fixed.Distance - original.Distance,
		MagnitudeDiff:      fixed.Magnitude - original.Magnitude,
		PhaseDiff:          fixed.Phase - original.Phase,
		ElongationDiff:     fixed.Elongation - original.Elongation,
	}
}

// generateAnalysis 生成分析说明
func (c *PlanetCalculatorCompare) generateAnalysis(original *PlanetPosition, fixed *PlanetPositionFixed, diff *PlanetPositionDiff, planetName string) string {
	analysis := fmt.Sprintf("=== %s 行星位置计算对比分析 ===\n\n", planetName)

	// 赤经分析
	analysis += "【赤经 (Right Ascension)】\n"
	analysis += fmt.Sprintf("  原始值: %.4f 小时 (%.2f°)\n", original.RightAscension, original.RightAscension*15)
	analysis += fmt.Sprintf("  修复值: %.4f 小时 (%.2f°)\n", fixed.RightAscension, fixed.RightAscension*15)
	analysis += fmt.Sprintf("  差异: %+.4f 小时\n", diff.RightAscensionDiff)
	if original.RightAscension < 0 {
		analysis += "  ⚠️ 原始值出现负值，这是错误的！赤经应该在 0-24h 范围内。\n"
		analysis += "  ✅ 修复后已归一化到正确范围。\n"
	}
	if fixed.RightAscension >= 0 && fixed.RightAscension < 24 {
		analysis += "  ✅ 修复后赤经在有效范围内。\n"
	}
	analysis += "\n"

	// 赤纬分析
	analysis += "【赤纬 (Declination)】\n"
	analysis += fmt.Sprintf("  原始值: %.4f°\n", original.Declination)
	analysis += fmt.Sprintf("  修复值: %.4f°\n", fixed.Declination)
	analysis += fmt.Sprintf("  差异: %+.4f°\n", diff.DeclinationDiff)
	if fixed.Declination >= -90 && fixed.Declination <= 90 {
		analysis += "  ✅ 修复后赤纬在有效范围内 (-90° 到 +90°)。\n"
	}
	analysis += "\n"

	// 距离分析
	analysis += "【距离 (Distance)】\n"
	analysis += fmt.Sprintf("  原始值: %.4f AU\n", original.Distance)
	analysis += fmt.Sprintf("  修复值: %.4f AU\n", fixed.Distance)
	analysis += fmt.Sprintf("  差异: %+.4f AU\n", diff.DistanceDiff)
	if original.Distance == 1.2 {
		analysis += "  ⚠️ 原始值是固定值 1.2，没有实际计算！\n"
		analysis += "  ✅ 修复后基于轨道要素计算实际距离。\n"
	}
	analysis += "\n"

	// 星等分析
	analysis += "【星等 (Magnitude)】\n"
	analysis += fmt.Sprintf("  原始值: %.4f\n", original.Magnitude)
	analysis += fmt.Sprintf("  修复值: %.4f\n", fixed.Magnitude)
	analysis += fmt.Sprintf("  差异: %+.4f\n", diff.MagnitudeDiff)
	analysis += "\n"

	// 相位分析
	analysis += "【相位 (Phase)】\n"
	analysis += fmt.Sprintf("  原始值: %.4f\n", original.Phase)
	analysis += fmt.Sprintf("  修复值: %.4f\n", fixed.Phase)
	analysis += fmt.Sprintf("  差异: %+.4f\n", diff.PhaseDiff)
	if original.Phase == 0.95 {
		analysis += "  ⚠️ 原始值是固定值 0.95，没有实际计算！\n"
		analysis += "  ✅ 修复后基于几何关系计算实际相位。\n"
	}
	analysis += "\n"

	// 距角分析
	analysis += "【距角 (Elongation)】\n"
	analysis += fmt.Sprintf("  原始值: %.4f°\n", original.Elongation)
	analysis += fmt.Sprintf("  修复值: %.4f°\n", fixed.Elongation)
	analysis += fmt.Sprintf("  差异: %+.4f°\n", diff.ElongationDiff)
	if original.Elongation == 120 {
		analysis += "  ⚠️ 原始值是固定值 120°，没有实际计算！\n"
		analysis += "  ✅ 修复后基于几何关系计算实际距角。\n"
	}
	analysis += "\n"

	// 总结
	analysis += "【修复总结】\n"
	analysis += "1. 赤经归一化：修复了 math.Atan2 返回负值的问题，确保赤经在 0-24h 范围内\n"
	analysis += "2. 轨道计算：使用开普勒方程计算真实轨道位置，而非简化模型\n"
	analysis += "3. 坐标转换：正确实现日心坐标到地心坐标的转换\n"
	analysis += "4. 距离计算：基于轨道要素计算实际距离，而非固定值\n"
	analysis += "5. 相位计算：基于太阳-地球-行星几何关系计算\n"
	analysis += "6. 距角计算：基于实际位置计算距角，而非固定值\n"

	return analysis
}
