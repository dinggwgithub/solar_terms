package bugs

// BugType Bug类型枚举
type BugType int

const (
	BugTypeNone BugType = iota
	BugTypeInstability
	BugTypeConstraint
	BugTypePrecision
)

// String 返回Bug类型的字符串表示
func (bt BugType) String() string {
	switch bt {
	case BugTypeNone:
		return "none"
	case BugTypeInstability:
		return "instability"
	case BugTypeConstraint:
		return "constraint"
	case BugTypePrecision:
		return "precision"
	default:
		return "unknown"
	}
}

// Bug Bug接口定义
type Bug interface {
	// Name 返回Bug名称
	Name() string
	// Description 返回Bug描述
	Description() string
	// Apply 应用Bug到计算结果
	Apply(calculationType string, params interface{}) (interface{}, []string)
}

// ParseBugType 从字符串解析Bug类型
func ParseBugType(bugTypeStr string) BugType {
	switch bugTypeStr {
	case "none":
		return BugTypeNone
	case "instability":
		return BugTypeInstability
	case "constraint":
		return BugTypeConstraint
	case "precision":
		return BugTypePrecision
	default:
		return BugTypeNone
	}
}