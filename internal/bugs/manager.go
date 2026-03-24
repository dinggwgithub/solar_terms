package bugs

import (
	"fmt"
)

// BugManager Bug管理器
type BugManager struct {
	bugs map[BugType]Bug
}

// NewBugManager 创建新的Bug管理器
func NewBugManager() *BugManager {
	manager := &BugManager{
		bugs: make(map[BugType]Bug),
	}
	
	// 注册所有Bug实现
	manager.RegisterBug(BugTypeInstability, NewInstabilityBug())
	manager.RegisterBug(BugTypeConstraint, NewConstraintBug())
	manager.RegisterBug(BugTypePrecision, NewPrecisionBug())
	
	return manager
}

// RegisterBug 注册Bug实现
func (m *BugManager) RegisterBug(bugType BugType, bug Bug) {
	m.bugs[bugType] = bug
}

// ApplyBug 应用指定类型的Bug
func (m *BugManager) ApplyBug(bugType BugType, calculationType string, params interface{}) (interface{}, []string) {
	if bug, exists := m.bugs[bugType]; exists {
		return bug.Apply(calculationType, params)
	}
	return nil, []string{"未知的Bug类型: " + bugType.String()}
}

// GetBug 获取指定类型的Bug
func (m *BugManager) GetBug(bugType BugType) (Bug, bool) {
	bug, exists := m.bugs[bugType]
	return bug, exists
}

// GetAllBugs 获取所有注册的Bug
func (m *BugManager) GetAllBugs() map[BugType]Bug {
	return m.bugs
}

// GetBugInfo 获取Bug的详细信息
func (m *BugManager) GetBugInfo(bugType BugType) (map[string]string, error) {
	bug, exists := m.bugs[bugType]
	if !exists {
		return nil, fmt.Errorf("Bug类型不存在: %s", bugType.String())
	}
	
	info := map[string]string{
		"name":        bug.Name(),
		"description": bug.Description(),
		"type":        bugType.String(),
	}
	
	// 添加Bug特征信息
	if instabilityBug, ok := bug.(*InstabilityBug); ok {
		characteristics := instabilityBug.GetBugCharacteristics()
		for k, v := range characteristics {
			info["characteristic_"+k] = v
		}
	}
	
	if constraintBug, ok := bug.(*ConstraintBug); ok {
		characteristics := constraintBug.GetBugCharacteristics()
		for k, v := range characteristics {
			info["characteristic_"+k] = v
		}
	}
	
	if precisionBug, ok := bug.(*PrecisionBug); ok {
		characteristics := precisionBug.GetBugCharacteristics()
		for k, v := range characteristics {
			info["characteristic_"+k] = v
		}
	}
	
	return info, nil
}

// GetBugFixSuggestions 获取Bug修复建议
func (m *BugManager) GetBugFixSuggestions(bugType BugType) ([]string, error) {
	bug, exists := m.bugs[bugType]
	if !exists {
		return nil, fmt.Errorf("Bug类型不存在: %s", bugType.String())
	}
	
	var suggestions []string
	
	// 获取特定Bug的修复建议
	if instabilityBug, ok := bug.(*InstabilityBug); ok {
		suggestions = instabilityBug.GetFixSuggestions()
	}
	
	if constraintBug, ok := bug.(*ConstraintBug); ok {
		suggestions = constraintBug.GetFixSuggestions()
	}
	
	if precisionBug, ok := bug.(*PrecisionBug); ok {
		suggestions = precisionBug.GetFixSuggestions()
	}
	
	// 添加通用修复建议
	generalSuggestions := []string{
		"添加单元测试覆盖边界情况",
		"实现输入参数验证",
		"添加错误处理和日志记录",
		"进行代码审查和重构",
		"参考权威算法实现",
	}
	
	suggestions = append(suggestions, generalSuggestions...)
	
	return suggestions, nil
}

// ValidateBugApplication 验证Bug应用结果
func (m *BugManager) ValidateBugApplication(result interface{}, bugType BugType, calculationType string) (bool, []string) {
	var warnings []string
	
	bug, exists := m.bugs[bugType]
	if !exists {
		return false, []string{"Bug类型不存在: " + bugType.String()}
	}
	
	// 根据Bug类型进行验证
	switch bugType {
	case BugTypeConstraint:
		if constraintBug, ok := bug.(*ConstraintBug); ok {
			err := constraintBug.ValidateConstraints(result, calculationType)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("约束验证失败: %v", err))
				return false, warnings
			}
		}
	
	case BugTypePrecision:
		// 精度验证需要基准数据对比
		warnings = append(warnings, "精度验证需要基准数据对比")
	
	case BugTypeInstability:
		// 不稳定性Bug的验证需要多次调用
		warnings = append(warnings, "不稳定性验证需要多次调用测试")
	}
	
	return true, warnings
}

// TestCase 测试用例结构
type TestCase struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Input       map[string]interface{} `json:"input"`
	Expected    string                 `json:"expected"`
}

// GetBugTestCases 获取Bug测试用例
func (m *BugManager) GetBugTestCases(bugType BugType, calculationType string) ([]TestCase, error) {
	var testCases []TestCase
	
	switch bugType {
	case BugTypeInstability:
		testCases = append(testCases, TestCase{
			Name:        "多次调用测试",
			Description: "同一参数连续调用10次，检查结果一致性",
			Input:       map[string]interface{}{"year": 2024, "term_index": 2},
			Expected:    "结果应该完全一致",
		})
	
	case BugTypeConstraint:
		testCases = append(testCases, TestCase{
			Name:        "边界值测试",
			Description: "输入边界值，检查结果是否在合理范围内",
			Input:       map[string]interface{}{"year": -1000, "term_index": 0},
			Expected:    "应该返回错误或默认值",
		})
	
	case BugTypePrecision:
		testCases = append(testCases, TestCase{
			Name:        "精度对比测试",
			Description: "与权威数据对比，检查误差是否在可接受范围内",
			Input:       map[string]interface{}{"julian_date": 2459580.5},
			Expected:    "误差小于0.000001度",
		})
	}
	
	// 添加通用测试用例
	generalTestCases := []TestCase{
		{
			Name:        "正常输入测试",
			Description: "输入正常参数，检查基本功能",
			Input:       map[string]interface{}{"year": 2024, "term_index": 2},
			Expected:    "返回有效结果",
		},
		{
			Name:        "异常输入测试",
			Description: "输入异常参数，检查错误处理",
			Input:       map[string]interface{}{"year": "invalid", "term_index": -1},
			Expected:    "返回错误信息",
		},
	}
	
	testCases = append(testCases, generalTestCases...)
	
	return testCases, nil
}

// GetBugStatistics 获取Bug统计信息
func (m *BugManager) GetBugStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"total_bugs": len(m.bugs),
		"bug_types":  []string{},
		"details":    map[string]interface{}{},
	}
	
	for bugType, bug := range m.bugs {
		bugTypes := stats["bug_types"].([]string)
		bugTypes = append(bugTypes, bugType.String())
		stats["bug_types"] = bugTypes
		
		details := stats["details"].(map[string]interface{})
		details[bugType.String()] = map[string]string{
			"name":        bug.Name(),
			"description": bug.Description(),
		}
	}
	
	return stats
}

// IsBugTypeValid 检查Bug类型是否有效
func (m *BugManager) IsBugTypeValid(bugType BugType) bool {
	_, exists := m.bugs[bugType]
	return exists
}

// GetSupportedBugTypes 获取支持的Bug类型
func (m *BugManager) GetSupportedBugTypes() []BugType {
	var supportedTypes []BugType
	for bugType := range m.bugs {
		supportedTypes = append(supportedTypes, bugType)
	}
	return supportedTypes
}

// GetBugTypeFromString 从字符串获取Bug类型
func (m *BugManager) GetBugTypeFromString(bugTypeStr string) (BugType, error) {
	bugType := ParseBugType(bugTypeStr)
	if bugType == BugTypeNone && bugTypeStr != "none" {
		return BugTypeNone, fmt.Errorf("无效的Bug类型: %s", bugTypeStr)
	}
	
	if !m.IsBugTypeValid(bugType) && bugTypeStr != "none" {
		return BugTypeNone, fmt.Errorf("Bug类型未注册: %s", bugTypeStr)
	}
	
	return bugType, nil
}