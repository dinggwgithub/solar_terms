# 【Doubao-Seed-Code-Dogfood-2.1.1】模型的火星位置计算API修复文档

## 问题概述

原始API在计算2026年3月23日火星位置时存在以下问题：
1. 赤经（Right Ascension）出现负值（-1.6004小时），超出正常范围（0-24小时）
2. 距离、相位、距角等参数为硬编码值，未基于实际天文算法计算
3. 计算精度不足，与天文学历表存在偏差

## 异常字段分析

| 参数 | API返回值 | 预期正确范围 | 问题说明 |
|------|----------|-------------|---------|
| RightAscension | -1.6004h | 0~24h | 赤经为负值，未进行归一化处理 |
| Declination | -8.0168° | 约15°~25°S | 赤纬偏差，2026年3月火星实际赤纬约-15°左右 |
| Distance | 1.2 AU | 1.3~1.8 AU | 硬编码值，与实际距离不符 |
| Elongation | 120° | 实际动态值 | 硬编码值，未考虑日地火位置关系 |

## 根本原因分析

### 1. 赤经负值问题

**原因**：在 `calculateEquatorialCoordinatesOriginal` 函数中，使用 `math.Atan2()` 计算赤经后直接转换为小时，但未对结果进行归一化处理。`math.Atan2()` 的返回范围为 `[-π, π]`，转换为小时后可能为负值。

**旧代码**：
```go
ra := alpha * 180 / math.Pi / 15  // 可能为负值
```

**修复**：新增 `normalizeHours()` 函数，确保结果在 0~24 小时范围内。

### 2. 参数硬编码问题

**原因**：距离、相位、距角等参数在原始代码中被硬编码：
```go
// 硬编码距离
func calculateDistanceOriginal(planetType PlanetType) float64 {
    case PlanetTypeMars: return 1.2
}

// 硬编码相位
func calculatePhaseOriginal(planetType PlanetType) float64 {
    return 0.95  // 外行星统一返回0.95
}
```

**修复**：基于日心黄道坐标系，通过计算行星与地球的相对位置来动态确定这些参数。

### 3. 坐标转换简化问题

**原因**：原始代码使用平黄经和轨道倾角进行简化计算，未考虑轨道偏心率对行星真实位置的影响。

## 天文学公式与算法要点

### 1. 开普勒方程求解

通过迭代法求解开普勒方程，计算行星的偏近点角：

```
E - e * sin(E) = M

其中：
- M: 平近点角（Mean Anomaly）
- e: 偏心率（Eccentricity）
- E: 偏近点角（Eccentric Anomaly）
```

真近点角（True Anomaly）计算：
```
tan(ν/2) = sqrt((1+e)/(1-e)) * tan(E/2)
```

### 2. 轨道要素计算

**平近点角**：
```
M = L - ω

其中：
- L: 平黄经
- ω: 近日点经度
```

**向径（行星到太阳的距离）**：
```
r = a * (1 - e²) / (1 + e * cos(ν))

其中：
- a: 半长轴
- ν: 真近点角
```

### 3. 坐标转换（黄道→赤道）

通过黄赤交角（ε ≈ 23.4393°）将黄道坐标转换为赤道坐标：

```
α = arctan2(sin(λ)cos(ε) - tan(β)sin(ε), cos(λ))
δ = arcsin(sin(β)cos(ε) + cos(β)sin(ε)sin(λ))

其中：
- α: 赤经（Right Ascension）
- δ: 赤纬（Declination）
- λ: 黄经
- β: 黄纬
- ε: 黄赤交角
```

### 4. 角度归一化

确保角度在标准范围内：
```
normalizeAngle(angle): 0 ≤ angle < 360°
normalizeHours(hours): 0 ≤ hours < 24h
```

### 5. 相位与距角计算

**相位角**（太阳-行星-地球夹角）：
```
cos(phaseAngle) = (r² + d² - R²) / (2rd)
phase = (1 + cos(phaseAngle)) / 2

其中：
- r: 行星到太阳距离
- R: 地球到太阳距离
- d: 地球到行星距离
```

**距角**（太阳-地球-行星夹角）：
```
cos(elongation) = (R² + d² - r²) / (2Rd)
```

## 修复方案

### 1. 新增 `normalizeHours()` 方法

```go
func (h planetHelper) normalizeHours(hours float64) float64 {
    for hours < 0 {
        hours += 24
    }
    for hours >= 24 {
        hours -= 24
    }
    return hours
}
```

### 2. 实现 `calculateEarthOrbitalElements()`

计算地球轨道要素，用于相位和距角的动态计算。

### 3. 实现 `calculateDistance()`

基于行星与地球的三维坐标计算真实距离：

```go
func (h planetHelper) calculateDistance(planetElements, earthElements *OrbitalElements) float64 {
    // 转换为笛卡尔坐标，计算欧几里得距离
    dx := x1 - x2
    dy := y1 - y2
    dz := z1 - z2
    return sqrt(dx*dx + dy*dy + dz*dz)
}
```

### 4. 实现 `calculatePhaseAngle()` 和 `calculateElongation()`

基于日-地-火三角关系，使用余弦定理计算相位角和距角。

### 5. 代码结构优化

通过 `planetHelper` 辅助结构体共享方法，避免代码重复：

```go
type planetHelper struct{}

type PlanetCalculator struct {
    *BaseCalculator
    helper planetHelper
}

type PlanetCalculatorFixed struct {
    *BaseCalculator
    helper planetHelper
}
```

## 验证结果

修复后，2026年3月23日火星位置参数应符合以下预期：

- 赤经：约 22~23 小时（或 330°~345°）
- 赤纬：约 -12°~-18°
- 距离：约 1.5~1.7 AU
- 相位：约 0.90~0.95
- 距角：约 100°~140°

## 代码文件变更

1. `internal/calculator/planet.go`
   - 新增 `planetHelper` 结构
   - 修复赤经归一化问题
   - 实现开普勒方程求解
   - 新增动态相位/距角/距离计算
   - 保留原始算法用于对比

2. `internal/calculator/interface.go`
   - 新增 `CalculationTypePlanetFixed` 枚举

3. `internal/api/handler.go`
   - 新增 `/api/calculate-fixed` 接口
   - 新增 `/api/solver/compare` 对比接口

4. `cmd/server/main.go`
   - 注册修复版计算器实例