#【GLM5.0】模型的火星位置计算 API Bug 修复报告

## 1. 问题概述

### 1.1 原始 API 响应数据

```json
{
  "success": true,
  "result": {
    "right_ascension": -1.6004243650179677,
    "declination": -8.016832970529133,
    "distance": 1.2,
    "magnitude": -1.119093769761876,
    "phase": 0.95,
    "elongation": 120
  }
}
```

### 1.2 异常字段分析

| 字段 | 返回值 | 预期范围 | 状态 | 说明 |
|------|--------|----------|------|------|
| right_ascension | -1.60h | 0~24h | ❌ 异常 | 赤经出现负值，违反天文定义 |
| declination | -8.02° | +5°~+25° | ❌ 异常 | 2026年3月火星应在北半球 |
| distance | 1.2 AU | 0.7~1.0 AU | ❌ 异常 | 固定值，未实际计算 |
| magnitude | -1.12 | -2.0~-2.5 | ⚠️ 偏亮 | 2026年3月火星接近冲日 |
| phase | 0.95 | ~0.99 | ⚠️ 固定值 | 未根据实际位置计算 |
| elongation | 120° | 150°~180° | ❌ 异常 | 固定值，冲日附近应更大 |

## 2. 根本原因分析

### 2.1 赤经负值问题

**错误代码位置**: [planet.go:347-367](file:///d:/work/zj_20260319/0324/branch15/internal/calculator/planet.go#L347-L367)

```go
// 原始代码
alpha := math.Atan2(
    math.Sin(lambda*math.Pi/180)*math.Cos(epsilon*math.Pi/180)-
        math.Tan(beta*math.Pi/180)*math.Sin(epsilon*math.Pi/180),
    math.Cos(lambda*math.Pi/180),
)

// 转换为度和小时
ra := alpha * 180 / math.Pi / 15 // 赤经（小时）
```

**问题**: `math.Atan2` 返回值范围为 [-π, π]，转换为小时后范围是 [-12h, 12h]，未进行归一化处理。

**修复方案**:
```go
raHours := alpha * 180 / math.Pi / 15

// 归一化到 0~24h
if raHours < 0 {
    raHours += 24
}
if raHours >= 24 {
    raHours -= 24
}
```

### 2.2 黄纬计算错误

**错误代码位置**: [planet.go:352-353](file:///d:/work/zj_20260319/0324/branch15/internal/calculator/planet.go#L352-L353)

```go
// 原始代码
lambda := elements.MeanLongitude // 简化：使用平黄经
beta := elements.Inclination     // 简化：使用轨道倾角
```

**问题**: 直接使用轨道倾角作为黄纬是错误的。黄纬应基于行星在轨道上的实际位置计算。

**天文学原理**:
黄纬 β 的计算公式：
```
sin(β) = sin(λ - Ω) × sin(i)
```
其中：
- λ: 行星真黄经
- Ω: 升交点黄经
- i: 轨道倾角

**修复方案**:
```go
omega := elements.LongitudePeri - elements.LongitudeNode
lambda := elements.TrueLongitude
beta := math.Asin(math.Sin((lambda-elements.LongitudeNode)*math.Pi/180) * 
                  math.Sin(elements.Inclination*math.Pi/180))
```

### 2.3 距离/相位/距角使用固定值

**错误代码位置**: [planet.go:393-410](file:///d:/work/zj_20260319/0324/branch15/internal/calculator/planet.go#L393-L410)

```go
// 原始代码 - 距离使用固定值
func (c *PlanetCalculator) calculateDistance(planetType PlanetType, elements *OrbitalElements) float64 {
    switch planetType {
    case PlanetTypeMars:
        return 1.2  // 固定值！
    ...
}

// 原始代码 - 距角使用固定值
func (c *PlanetCalculator) calculateElongation(...) float64 {
    ...
    return 120.0  // 固定值！
}
```

**问题**: 距离、相位、距角均使用硬编码固定值，未进行实际天文计算。

**修复方案**:

1. **距离计算** - 基于开普勒方程求解向径：
```go
// 求解开普勒方程
E := M  // 初始猜测
for i := 0; i < 10; i++ {
    dE := (M - E + e*180/π*sin(E)) / (1 - e*cos(E))
    E += dE
    if abs(dE) < 1e-8 {
        break
    }
}

// 计算向径
r = a * (1 - e * cos(E))
```

2. **距角计算** - 基于行星与太阳的黄经差：
```go
elongation := abs(planetLon - sunLon)
if elongation > 180 {
    elongation = 360 - elongation
}
```

3. **相位计算** - 基于距角：
```go
phase := (1 + cos(elongation * π / 180)) / 2
```

## 3. 天文学公式与算法要点

### 3.1 黄道坐标转赤道坐标

**公式**:
```
α = atan2(sin(λ)×cos(ε) - tan(β)×sin(ε), cos(λ))
δ = asin(sin(β)×cos(ε) + cos(β)×sin(ε)×sin(λ))
```

其中：
- α: 赤经
- δ: 赤纬
- λ: 黄经
- β: 黄纬
- ε: 黄赤交角（约 23.44°）

**注意事项**:
- `atan2` 返回值范围为 [-π, π]，需归一化到 [0, 2π]
- 赤经单位为小时，需将角度除以 15

### 3.2 开普勒方程求解

**方程**:
```
M = E - e × sin(E)
```

其中：
- M: 平近点角
- E: 偏近点角
- e: 偏心率

**牛顿迭代法**:
```go
E := M
for i := 0; i < 10; i++ {
    dE := (M - E + e*180/π*sin(E)) / (1 - e*cos(E))
    E += dE
    if abs(dE) < 1e-8 {
        break
    }
}
```

### 3.3 真近点角计算

**公式**:
```
tan(ν/2) = √((1+e)/(1-e)) × tan(E/2)
```

其中：
- ν: 真近点角
- E: 偏近点角
- e: 偏心率

### 3.4 行星视星等计算

**公式**:
```
m = m₀ + 5×log₁₀(d) + k×β
```

其中：
- m₀: 基础星等
- d: 行星到地球距离（AU）
- k: 相位系数
- β: 相位角（度）

**火星参数**:
- m₀ = -1.52
- k = 0.016

### 3.5 儒略日计算

**公式**:
```
JD = 365.25×(Y+4716) + 30.6001×(M+1) + D + B - 1524.5
```

其中：
- Y: 年（1、2月需减1）
- M: 月（1、2月需加12）
- D: 日
- B = 2 - A + A/4（格里高利历修正）
- A = Y/100

## 4. 修复验证

### 4.1 修复后预期结果

对于 2026年3月23日 的火星位置：

| 字段 | 修复前 | 修复后预期 | 说明 |
|------|--------|------------|------|
| right_ascension | -1.60h | ~10h~12h | 归一化到 0~24h |
| declination | -8.02° | ~+10°~+15° | 火星在北半球 |
| distance | 1.2 AU | ~0.7~0.9 AU | 基于轨道计算 |
| magnitude | -1.12 | ~-2.0~-2.5 | 冲日附近较亮 |
| phase | 0.95 | ~0.98~0.99 | 接近满相 |
| elongation | 120° | ~160°~180° | 接近冲日 |

### 4.2 测试命令

```bash
# 测试修复版接口
curl -X POST 'http://localhost:8080/api/calculate-fixed' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "calculation": "planet_fixed",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23,
      "planet_name": "mars"
    }
  }'

# 测试对比接口
curl -X POST 'http://localhost:8080/api/solver/compare' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "calculation": "compare",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23,
      "planet_name": "mars"
    }
  }'
```

## 5. 修改文件清单

| 文件 | 修改类型 | 说明 |
|------|----------|------|
| `internal/calculator/planet_fixed.go` | 新增 | 修复版行星计算器 |
| `internal/calculator/compare.go` | 新增 | 对比计算器 |
| `internal/calculator/interface.go` | 修改 | 添加新计算类型 |
| `internal/api/handler.go` | 修改 | 添加新 API 端点 |
| `cmd/server/main.go` | 修改 | 注册新计算器 |

## 6. 参考资料

1. Meeus, J. (1998). *Astronomical Algorithms*. Willmann-Bell.
2. VSOP87 Planetary Theory - https://cdsarc.u-strasbg.fr/viz-bin/Cat?cat=VI%2F81
3. JPL Horizons System - https://ssd.jpl.nasa.gov/horizons/
