package bugs

import (
	"math/rand"
	"sync"
	"time"
)

// BugDynamicConfig Bug动态化配置
type BugDynamicConfig struct {
	// Constraint Bug配置
	ConstraintMultiplierMin float64 `json:"constraint_multiplier_min"` // 最小乘数
	ConstraintMultiplierMax float64 `json:"constraint_multiplier_max"` // 最大乘数
	ConstraintPositiveMult  float64 `json:"constraint_positive_mult"`  // 正数实际乘数
	ConstraintNegativeDiv   float64 `json:"constraint_negative_div"`   // 负数实际除数

	// Precision Bug配置
	PrecisionMinDecimals int    `json:"precision_min_decimals"` // 最小小数位数
	PrecisionMaxDecimals int    `json:"precision_max_decimals"` // 最大小数位数
	PrecisionDecimals    int    `json:"precision_decimals"`     // 实际小数位数
	PrecisionRoundMethod string `json:"precision_round_method"` // 舍入方法：round, floor, ceil

	// Instability Bug配置
	InstabilityMinVariation float64 `json:"instability_min_variation"` // 最小波动范围
	InstabilityMaxVariation float64 `json:"instability_max_variation"` // 最大波动范围
	InstabilityVariation    float64 `json:"instability_variation"`     // 实际波动范围(±%)
	InstabilityDistribution string  `json:"instability_distribution"`  // 分布类型：uniform, gaussian, skewed

	// 混合Bug配置
	EnableMixedMode bool      `json:"enable_mixed_mode"` // 是否启用混合模式
	MixedBugTypes   []BugType `json:"mixed_bug_types"`   // 混合的Bug类型
	MixedApplyOrder []BugType `json:"mixed_apply_order"` // 应用顺序

	// 不完美Bug特征
	EnableImperfection bool    `json:"enable_imperfection"` // 是否启用不完美特征
	ImperfectionLevel  float64 `json:"imperfection_level"`  // 不完美程度(±%)

	// 会话级随机种子（用于保证同一会话返回一致的不完美噪声一致）
	SessionID string `json:"session_id"`
	RandSeed  int64  `json:"rand_seed"`
	counter   uint64 `json:"-"` // 用于生成序列随机数序列
}

// NewDefaultBugConfig 创建默认配置
func NewDefaultBugConfig(sessionID string) *BugDynamicConfig {
	config := &BugDynamicConfig{
		// Constraint 默认范围: 100~10000倍
		ConstraintMultiplierMin: 100.0,
		ConstraintMultiplierMax: 10000.0,

		// Precision 默认: 1~4位小数
		PrecisionMinDecimals: 1,
		PrecisionMaxDecimals: 4,

		// Instability 默认范围: ±1%~±10%
		InstabilityMinVariation: 0.01,
		InstabilityMaxVariation: 0.10,

		// 混合模式默认关闭
		EnableMixedMode: false,
		MixedBugTypes:   []BugType{},
		MixedApplyOrder: []BugType{},

		// 默认启用不完美特征
		EnableImperfection: true,
		ImperfectionLevel:  0.02, // ±2%的不完美

		SessionID: sessionID,
	}

	// 随机生成具体参数
	config.Randomize()
	return config
}

// Randomize 随机生成所有动态参数
func (c *BugDynamicConfig) Randomize() {
	// 基于会话ID和当前时间生成随机种子，保证同一会话配置相同
	seedHash := hashString(c.SessionID)
	c.RandSeed = time.Now().UnixNano() + seedHash
	c.counter = 0

	// 创建确定性随机源，确保同一会话生成相同配置
	localRand := rand.New(rand.NewSource(c.RandSeed))

	// Constraint: 随机选择乘数
	c.ConstraintPositiveMult = c.ConstraintMultiplierMin +
		localRand.Float64()*(c.ConstraintMultiplierMax-c.ConstraintMultiplierMin)
	c.ConstraintNegativeDiv = c.ConstraintMultiplierMin +
		localRand.Float64()*(c.ConstraintMultiplierMax-c.ConstraintMultiplierMin)

	// Precision: 随机选择小数位数和舍入方法
	c.PrecisionDecimals = localRand.Intn(c.PrecisionMaxDecimals-c.PrecisionMinDecimals+1) +
		c.PrecisionMinDecimals
	roundMethods := []string{"round", "floor", "ceil"}
	c.PrecisionRoundMethod = roundMethods[localRand.Intn(len(roundMethods))]

	// Instability: 随机选择波动范围和分布
	c.InstabilityVariation = c.InstabilityMinVariation +
		localRand.Float64()*(c.InstabilityMaxVariation-c.InstabilityMinVariation)
	distributions := []string{"uniform", "gaussian", "skewed"}
	c.InstabilityDistribution = distributions[localRand.Intn(len(distributions))]

	// 如果启用混合模式，随机选择Bug类型组合
	if c.EnableMixedMode {
		allBugTypes := []BugType{BugTypeConstraint, BugTypePrecision, BugTypeInstability}
		numTypes := localRand.Intn(2) + 2 // 选择2-3种Bug类型
		c.MixedBugTypes = make([]BugType, numTypes)
		c.MixedApplyOrder = make([]BugType, numTypes)

		// 随机选择Bug类型
		perm := localRand.Perm(len(allBugTypes))
		for i := 0; i < numTypes; i++ {
			c.MixedBugTypes[i] = allBugTypes[perm[i]]
			c.MixedApplyOrder[i] = allBugTypes[perm[i]]
		}

		// 随机打乱应用顺序
		localRand.Shuffle(len(c.MixedApplyOrder), func(i, j int) {
			c.MixedApplyOrder[i], c.MixedApplyOrder[j] = c.MixedApplyOrder[j], c.MixedApplyOrder[i]
		})
	}
}

// Clone 克隆配置
func (c *BugDynamicConfig) Clone() *BugDynamicConfig {
	clone := *c
	clone.MixedBugTypes = make([]BugType, len(c.MixedBugTypes))
	clone.MixedApplyOrder = make([]BugType, len(c.MixedApplyOrder))
	copy(clone.MixedBugTypes, c.MixedBugTypes)
	copy(clone.MixedApplyOrder, c.MixedApplyOrder)
	return &clone
}

// ConfigManager 配置管理器（会话级）
type ConfigManager struct {
	sessions map[string]*BugDynamicConfig
	mu       sync.RWMutex
}

var (
	globalConfigManager *ConfigManager
	configManagerOnce   sync.Once
)

// GetGlobalConfigManager 获取全局配置管理器单例
func GetGlobalConfigManager() *ConfigManager {
	configManagerOnce.Do(func() {
		globalConfigManager = &ConfigManager{
			sessions: make(map[string]*BugDynamicConfig),
		}
	})
	return globalConfigManager
}

// GetOrCreateSessionConfig 获取或创建会话配置
func (m *ConfigManager) GetOrCreateSessionConfig(sessionID string) *BugDynamicConfig {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config, exists := m.sessions[sessionID]; exists {
		return config
	}

	config := NewDefaultBugConfig(sessionID)
	m.sessions[sessionID] = config
	return config
}

// GetSessionConfig 获取会话配置（如果不存在则返回nil）
func (m *ConfigManager) GetSessionConfig(sessionID string) *BugDynamicConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID]
}

// SetSessionConfig 设置会话配置
func (m *ConfigManager) SetSessionConfig(sessionID string, config *BugDynamicConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[sessionID] = config
}

// ClearSessionConfig 清除会话配置
func (m *ConfigManager) ClearSessionConfig(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, sessionID)
}

// ClearAllSessions 清除所有会话配置
func (m *ConfigManager) ClearAllSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions = make(map[string]*BugDynamicConfig)
}

// GenerateSessionID 生成会话ID
func GenerateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// hashString 简单的字符串哈希函数
func hashString(s string) int64 {
	var h int64 = 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int64(s[i])
	}
	return h
}

// GetSeededRand 获取基于当前值的确定性随机源（暴露给外部使用）
func (c *BugDynamicConfig) GetSeededRand(value float64) *rand.Rand {
	// 使用配置种子 + 值的哈希来生成确定性随机数
	valueHash := int64(value * 1000000) // 将浮点数转换为整数哈希
	combinedSeed := c.RandSeed + valueHash
	return rand.New(rand.NewSource(combinedSeed))
}

// getSeededRand 获取基于当前值的确定性随机源（内部使用）
func (c *BugDynamicConfig) getSeededRand(value float64) *rand.Rand {
	return c.GetSeededRand(value)
}

// ApplyConstraint 应用Constraint Bug（使用动态配置）
func (c *BugDynamicConfig) ApplyConstraint(value float64) float64 {
	if value == 0 {
		return 0
	}

	// 基础乘数/除数
	var result float64
	if value > 0 {
		result = value * c.ConstraintPositiveMult
	} else {
		result = value / c.ConstraintNegativeDiv
	}

	// 添加不完美噪声（使用确定性随机源）
	if c.EnableImperfection {
		localRand := c.getSeededRand(value)
		noise := 1.0 + (localRand.Float64()-0.5)*2*c.ImperfectionLevel
		result *= noise
	}

	return result
}

// ApplyPrecision 应用Precision Bug（使用动态配置）
func (c *BugDynamicConfig) ApplyPrecision(value float64) float64 {
	factor := 1.0
	for i := 0; i < c.PrecisionDecimals; i++ {
		factor *= 10
	}

	var result float64
	switch c.PrecisionRoundMethod {
	case "round":
		result = round(value*factor) / factor
	case "floor":
		result = floor(value*factor) / factor
	case "ceil":
		result = ceil(value*factor) / factor
	default:
		result = round(value*factor) / factor
	}

	// 添加不完美噪声（使用确定性随机源）
	if c.EnableImperfection {
		localRand := c.getSeededRand(value)
		noise := 1.0 + (localRand.Float64()-0.5)*2*c.ImperfectionLevel
		epsilon := (localRand.Float64() - 0.5) * 0.1 / factor
		result += epsilon * noise
	}

	return result
}

// ApplyInstability 应用Instability Bug（使用动态配置）
func (c *BugDynamicConfig) ApplyInstability(value float64) float64 {
	if value == 0 {
		return 0
	}

	// 使用确定性随机源
	localRand := c.getSeededRand(value)

	var change float64
	switch c.InstabilityDistribution {
	case "uniform":
		change = 1.0 + (localRand.Float64()-0.5)*2*c.InstabilityVariation
	case "gaussian":
		// 简单的高斯近似
		u1 := localRand.Float64()
		u2 := localRand.Float64()
		z := sqrt(-2*ln(u1)) * cos(2*pi*u2)
		change = 1.0 + z*c.InstabilityVariation/3 // 除以3让大部分落在±3σ内
	case "skewed":
		// 偏态分布
		u := localRand.Float64()
		if u < 0.7 { // 70%正向偏移
			change = 1.0 + localRand.Float64()*c.InstabilityVariation
		} else { // 30%负向偏移
			change = 1.0 - localRand.Float64()*c.InstabilityVariation*0.5
		}
	default:
		change = 1.0 + (localRand.Float64()-0.5)*2*c.InstabilityVariation
	}

	return value * change
}

// 辅助数学函数
func round(x float64) float64 {
	if x < 0 {
		return -floor(-x + 0.5)
	}
	return floor(x + 0.5)
}

func floor(x float64) float64 {
	if x >= 0 {
		return float64(int64(x))
	}
	return float64(int64(x) - 1)
}

func ceil(x float64) float64 {
	if x >= 0 {
		return float64(int64(x) + 1)
	}
	return -floor(-x)
}

func sqrt(x float64) float64 {
	// 简单的平方根近似
	if x <= 0 {
		return 0
	}
	z := x / 2
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

func ln(x float64) float64 {
	// 简单的自然对数近似
	if x <= 0 {
		return -1e10
	}
	// 使用泰勒级数近似（简化版）
	result := 0.0
	for x > 2 {
		result += 0.69314718056
		x /= 2
	}
	x--
	z := x
	term := x
	n := 1
	for i := 0; i < 20; i++ {
		term *= -x
		n++
		z += term / float64(n)
	}
	return result + z
}

func cos(x float64) float64 {
	// 简单的余弦近似（使用泰勒级数）
	x = mod2pi(x)
	result := 1.0
	term := 1.0
	for n := 1; n < 10; n++ {
		term *= -x * x / float64((2*n-1)*(2*n))
		result += term
	}
	return result
}

func mod2pi(x float64) float64 {
	pi := 3.141592653589793
	twoPi := 2 * pi
	for x >= twoPi {
		x -= twoPi
	}
	for x < 0 {
		x += twoPi
	}
	return x
}

var pi = 3.141592653589793
