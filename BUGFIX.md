# 【Kimi Km2.5】模型的火星位置计算 API 修复文档

## 问题概述

在 2026 年 3 月 23 日的火星位置计算中，API 返回了不合理的数据，主要问题包括：
- 赤经 (right_ascension) 出现负值 (-1.60 小时)
- 距离、相位、距角等参数使用固定值而非实际计算

## 异常数据分析

### API 返回数据（2026-03-23）
```json
{
  "right_ascension": -1.6004243650179677,  // ❌ 异常：负值
  "declination": -8.016832970529133,       // ⚠️ 需验证
  "distance": 1.2,                          // ⚠️ 固定值
  "magnitude": -1.119093769761876,         // ⚠️ 基于固定值计算
  "phase": 0.95,                            // ⚠️ 固定值
  "elongation": 120                         // ⚠️ 固定值
}
```

### 预期正确范围（2026年3月火星位置）
| 参数 | 预期范围 | 说明 |
|------|----------|------|
| 赤经 | 22h - 23h | 火星位于室女座/天秤座附近 |
| 赤纬 | -10° ~ -15° | 南天球 |
| 距离 | 1.0 - 1.5 AU | 地火距离变化范围 |
| 星等 | +1.0 ~ +1.5 | 火星较暗时的星等 |
| 相位 | ~0.98 | 外行星接近满相 |
| 距角 | 0° - 180° | 与太阳的角距离 |

## 错误原因分析

### 1. 赤经负值问题

**根本原因**：`math.Atan2(y, x)` 函数返回值范围是 (-π, π)，即 (-180°, 180°)，转换为小时后范围是 (-12h, 12h)。当黄经位于第四象限时，赤经计算结果为负值。

**错误代码**（planet.go:340-345）：
```go
alpha := math.Atan2(
    math.Sin(lambda*math.Pi/180)*math.Cos(epsilon*math.Pi/180)-
        math.Tan(beta*math.Pi/180)*math.Sin(epsilon*math.Pi/180),
    math.Cos(lambda*math.Pi/180),
)
ra := alpha * 180 / math.Pi / 15 // 赤经（小时）
// ❌ 缺少归一化处理
```

**修复方案**：添加赤经归一化函数
```go
func (c *PlanetCalculatorFixed) normalizeHour(hour float64) float64 {
    for hour < 0 {
        hour += 24
    }
    for hour >= 24 {
        hour -= 24
    }
    return hour
}
```

### 2. 距离计算问题

**根本原因**：`calculateDistance` 函数直接返回固定值，没有基于轨道要素计算实际距离。

**错误代码**（planet.go:360-377）：
```go
func (c *PlanetCalculator) calculateDistance(planetType PlanetType, elements *OrbitalElements) float64 {
    switch planetType {
    case PlanetTypeMars:
        return 1.2  // ❌ 固定值
    // ...
    }
}
```

**修复方案**：使用开普勒方程计算实际轨道位置，然后计算地心距离
```go
// 计算日心坐标
planetHelio := c.calculateHeliocentricCoordinates(planetType, jd)
earthHelio := c.calculateHeliocentricCoordinates(PlanetTypeEarth, jd)

// 计算地心坐标
distance := math.Sqrt(
    math.Pow(planetHelio.X - earthHelio.X, 2) +
    math.Pow(planetHelio.Y - earthHelio.Y, 2) +
    math.Pow(planetHelio.Z - earthHelio.Z, 2),
)
```

### 3. 相位和距角计算问题

**根本原因**：使用固定值而非基于几何关系计算。

**修复方案**：基于太阳-地球-行星的几何关系计算
```go
// 计算相位角（太阳-行星-地球的夹角）
func calculatePhaseFixed(planet, earth *HeliocentricCoordinates) float64 {
    // 向量计算...
    phaseAngle := math.Acos(dot/(magD*magS)) * 180 / math.Pi
    phase := (1 + math.Cos(phaseAngle*math.Pi/180)) / 2
    return phase
}
```

## 修复内容总结

### 文件变更

| 文件 | 操作 | 说明 |
|------|------|------|
| `internal/calculator/planet_fixed.go` | 新增 | 修复版行星计算器 |
| `internal/calculator/planet_compare.go` | 新增 | 对比计算器 |
| `internal/calculator/interface.go` | 修改 | 添加新计算类型枚举 |
| `internal/api/handler.go` | 修改 | 添加新接口端点 |
| `cmd/server/main.go` | 修改 | 注册新计算器 |

### 关键修复点

1. **赤经归一化**：添加 `normalizeHour` 函数，确保赤经输出在 0-24h 范围内
2. **开普勒轨道计算**：实现开普勒方程迭代求解，计算真实近点角
3. **坐标转换**：正确实现日心坐标到地心坐标的转换
4. **距离计算**：基于轨道要素计算实际距离
5. **相位计算**：基于太阳-地球-行星几何关系计算
6. **距角计算**：基于实际位置计算距角

## 天文学公式参考

### 1. 开普勒方程
$$M = E - e \cdot \sin(E)$$

其中：
- $M$ = 平近点角
- $E$ = 偏近点角
- $e$ = 轨道偏心率

使用牛顿迭代法求解：
$$E_{n+1} = E_n + \frac{M + e \cdot \sin(E_n) - E_n}{1 - e \cdot \cos(E_n)}$$

### 2. 真近点角
$$\nu = 2 \cdot \arctan\left(\sqrt{\frac{1+e}{1-e}} \cdot \tan\frac{E}{2}\right)$$

### 3. 日心坐标
$$\begin{aligned}
x &= r \cdot (\cos\Omega \cdot \cos u - \sin\Omega \cdot \sin u \cdot \cos i) \\
y &= r \cdot (\sin\Omega \cdot \cos u + \cos\Omega \cdot \sin u \cdot \cos i) \\
z &= r \cdot \sin u \cdot \sin i
\end{aligned}$$

其中：
- $r = a \cdot (1 - e \cdot \cos E)$ = 日心距离
- $u = \omega + \nu$ = 升交点角距 + 真近点角
- $\Omega$ = 升交点经度
- $i$ = 轨道倾角

### 4. 黄道坐标转赤道坐标
$$\begin{aligned}
\tan(\alpha) &= \frac{\sin\lambda \cdot \cos\varepsilon - \tan\beta \cdot \sin\varepsilon}{\cos\lambda} \\
\sin(\delta) &= \sin\beta \cdot \cos\varepsilon + \cos\beta \cdot \sin\varepsilon \cdot \sin\lambda
\end{aligned}$$

其中：
- $\lambda$ = 黄经
- $\beta$ = 黄纬
- $\varepsilon$ = 黄赤交角
- $\alpha$ = 赤经
- $\delta$ = 赤纬

### 5. 黄赤交角
$$\varepsilon = 23.4392911° - 0.0130042° \cdot T - 1.64 \times 10^{-7} \cdot T^2 + 5.04 \times 10^{-7} \cdot T^3$$

其中 $T$ = 儒略世纪数

### 6. 相位计算
$$\text{相位} = \frac{1 + \cos(\theta)}{2}$$

其中 $\theta$ = 相位角（太阳-行星-地球的夹角）

### 7. 距角计算
$$\cos(E) = \frac{\vec{r}_{地日} \cdot \vec{r}_{地行}}{|\vec{r}_{地日}| \cdot |\vec{r}_{地行}|}$$

其中：
- $E$ = 距角
- $\vec{r}_{地日}$ = 地球到太阳的向量
- $\vec{r}_{地行}$ = 地球到行星的向量

## 接口说明

### 修复版接口

**端点**：`POST /api/calculate-fixed`

**请求格式**：与原 `/api/calculate` 完全一致

```json
{
  "calculation": "planet",
  "params": {
    "year": 2026,
    "month": 3,
    "day": 23,
    "planet_name": "mars"
  }
}
```

**响应格式**：与原接口一致，但数值已修复

### 对比接口

**端点**：`POST /api/solver/compare`

**功能**：返回原始计算和修复后计算的对比结果

**响应示例**：
```json
{
  "success": true,
  "result": {
    "original": { /* 原始计算结果 */ },
    "fixed": { /* 修复后结果 */ },
    "diff": { /* 差异值 */ },
    "analysis": "分析说明文本"
  }
}
```

## 测试建议

1. 使用对比接口验证修复效果
2. 测试不同日期的火星位置计算
3. 验证赤经始终在 0-24h 范围内
4. 验证其他行星（金星、木星等）的计算结果

## 参考资源

- VSOP87 行星理论
- 《天文算法》(Astronomical Algorithms) - Jean Meeus
- NASA JPL Horizons 系统
