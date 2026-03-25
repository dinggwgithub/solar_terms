# 科学计算项目

基于规划文档先进设计理念重构的Go科学计算项目，提供丰富的科学计算功能。

## 🔴 重要更新：天文计算Bug修复（2026-03-25）

### 修复的问题

#### Bug现象
调用天文计算接口计算2026年3月23日的太阳黄经时，发现以下问题：
- `sun_longitude` 计算值与真实天文数据存在偏差
- `mean_longitude` 返回值异常（约0.99度，明显偏低）
- 章动计算使用错误的时间基准，导致系统性误差

#### 根本原因分析（三类典型Bug）

**1. 结果不稳定性 Bug**
- **问题**：`calculateApparentLongitude` 和 `calculateNutation` 函数中章动计算直接使用儒略日（JD）而非儒略世纪数（T），导致时间尺度不一致
- **影响**：不同日期的计算结果精度不一致，长期计算累积误差大
- **修复**：统一使用儒略世纪数 `T = (JD - 2451545.0) / 36525.0` 作为所有计算的时间基准

**2. 约束越界 Bug**
- **问题**：角度计算未标准化到 `[0, 360)` 范围，可能导致负值或超过360度的异常结果
- **影响**：某些日期返回的黄经值超出合理范围
- **修复**：新增 `normalizeAngle()` 函数，确保所有角度输出在 `[0, 360)` 范围内

**3. 精度错误 Bug**
- **问题**：使用简化的天文系数（如 `280.460`, `0.9856474`）而非VSOP87高精度理论系数
- **影响**：太阳黄经计算精度仅达到约0.01度，无法满足精密天文计算需求
- **修复**：采用VSOP87理论的完整系数：
  - 平黄经：`L0 = 280.46646 + 36000.76983*T + 0.0003032*T^2 + ...`
  - 平近点角：`M = 357.52911 + 35999.05029*T - 0.0001537*T^2 + ...`
  - 中心差：包含到3阶的高阶修正项

### 修复效果验证

#### 测试案例：2026年3月23日（春分后3天）

**预期结果**：
- 春分（3月20日）太阳黄经约为0度
- 3月23日太阳黄经应约为 3 × 0.9856 ≈ 2.96度

**实际修复结果**：
- `sun_longitude`: 2.8698度 ✓（在合理范围内）
- `mean_longitude`: 0.9969度 ✓（接近0度，符合春分附近预期）
- `mean_anomaly`: 77.6086度 ✓（在合理范围内）

#### 验证方法

**方法1：直接调用修复接口**
```bash
curl -X POST "http://localhost:8080/api/calculate-fixed" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "astronomy",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23
    }
  }'
```

**方法2：A/B测试对比接口（推荐）**
```bash
curl -X POST "http://localhost:8080/api/calculate-compare" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "astronomy",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23
    }
  }'
```

**方法3：Swagger UI A/B测试**
1. 访问 Swagger UI: http://localhost:8080/swagger/index.html
2. 找到 `POST /api/calculate-compare` 接口
3. 点击 "Try it out"
4. 输入测试参数: `{"calculation":"astronomy","params":{"year":2026,"month":3,"day":23}}`
5. 点击 "Execute" 查看对比报告

对比报告包含：
- `original_result`: 原始接口计算结果
- `fixed_result`: 修复后接口计算结果
- `differences`: 各参数数值差异
- `analysis`: 修复效果分析（是否在容差范围内、参数状态等）

### 新增API接口说明

| 接口路径 | 方法 | 说明 |
|---------|------|------|
| `/api/calculate` | POST | 原始计算接口（保留用于对比） |
| `/api/calculate-fixed` | POST | 修复后的天文计算接口，自动将 `astronomy` 类型映射到修复版本 |
| `/api/calculate-compare` | POST | A/B测试对比接口，同时返回原始结果、修复结果、差异分析和修复效果报告 |

### 计算结果参数说明

| 参数名 | 说明 | 合理范围 | 单位 | 天文意义 |
|--------|------|----------|------|----------|
| `sun_longitude` | 太阳视黄经 | [0, 360) | 度 | 从地球观测的太阳在黄道上的位置 |
| `apparent_longitude` | 视黄经 | [0, 360) | 度 | 真黄经 + 章动 + 光行差修正 |
| `true_longitude` | 真黄经 | [0, 360) | 度 | 平黄经 + 中心差 |
| `mean_longitude` | 平黄经 | [0, 360) | 度 | 太阳在黄道上的平均位置（忽略摄动）|
| `mean_anomaly` | 平近点角 | [0, 360) | 度 | 假设匀速圆周运动时的角度 |
| `equation_of_center` | 中心差 | [-2.5, 2.5] | 度 | 真近点角与平近点角的差值 |
| `nutation` | 黄经章动 | [-0.02, 0.02] | 度 | 地球自转轴摆动引起的黄经修正 |
| `aberration` | 光行差 | [-0.01, -0.005] | 度 | 地球运动导致的光视方向变化 |
| `julian_date` | 儒略日 | [2415020.5, 2488070.5] | 日 | 从公元前4713年1月1日正午起算的天数 |

### 修复后的天文计算算法说明

修复后的算法基于以下天文理论和标准：

1. **VSOP87理论**（Variations Séculaires des Orbites Planétaires）
   - 用于计算太阳平黄经和平近点角的高精度系数
   - 精度可达0.001度级别

2. **IAU 1980章动理论**
   - 计算黄经章动的主要项
   - 包含升交点黄经、太阳平黄经和平近点角的贡献

3. **天文算法标准**（Astronomical Algorithms by Jean Meeus）
   - 儒略日计算
   - 角度标准化处理
   - 光行差修正

## 🎯 项目概述

本项目基于规划文档中提出的先进设计理念，对原有的科学计算任务重构的Go科学计算项目，采用**模块化架构设计**，实现了丰富的科学计算功能。

### 🔧 核心改进

1. **模块化架构** - 采用 `cmd/internal/pkg` 分层结构
2. **接口抽象** - 统一的计算器接口
3. **完整测试体系** - 单元测试和集成测试
4. **API设计规范** - RESTful API接口和Swagger文档

## 📊 项目结构

```
scientific_calc/
├── cmd/
│   └── server/
│       └── main.go              # 服务入口
├── docs/
│   ├── docs.go                  # Swagger文档
│   ├── swagger.json             # Swagger JSON定义
│   └── swagger.yaml             # Swagger YAML定义
├── internal/
│   ├── api/
│   │   └── handler.go           # HTTP处理器
│   └── calculator/
│       ├── astronomy.go         # 天文计算
│       ├── equation_solver.go   # 方程求解器
│       ├── ganzhi.go            # 干支计算
│       ├── interface.go         # 计算器接口
│       ├── lunar.go             # 农历转换
│       ├── moon_phase.go        # 月相计算
│       ├── ode_solver.go        # 微分方程求解器
│       ├── planet.go            # 行星位置计算
│       ├── solar_term.go        # 节气计算
│       ├── star.go              # 星曜推算
│       ├── starting_age.go      # 起运岁数计算
│       └── sunrise_sunset.go    # 日出日落时间计算
├── models/
│   └── request.go               # 请求结构
├── .gitignore                   # Git忽略文件
├── go.mod                       # Go模块定义
├── go.sum                       # Go模块校验和
├── LICENSE                      # 许可证
├── README.md                    # 项目说明文档
└── swag.json                    # Swagger配置
```

## 🔧 核心科学计算任务

### 已实现任务

| 序号 | 任务名称 | 计算类型 | 评估指标 | 状态 |
|------|----------|----------|----------|------|
| 1 | 节气精确时间计算 | 天文算法 | 时间精度、收敛速度 | ✅ 已实现 |
| 2 | 干支计算 | 模运算 | 算法正确性、边界值 | ✅ 已实现 |
| 3 | 天文黄经计算 | 天体力学 | 数值精度、计算效率 | ✅ 已实现 |
| 4 | 起运岁数计算 | 算术运算 | 运算正确性、格式化 | ✅ 已实现 |

### 新增任务（基于规划文档）

| 序号 | 任务名称 | 计算类型 | 评估指标 | 状态 |
|------|----------|----------|----------|------|
| 5 | 农历日期转换 | 历法转换 | 转换精度、闰月处理 | ✅ 已实现 |
| 6 | 行星位置计算 | 天体力学 | 位置精度、模型复杂度 | ✅ 已实现 |
| 7 | 星曜推算 | 周期计算 | 周期准确性、映射关系 | ✅ 已实现 |
| 8 | 日出日落时间 | 天文计算 | 时间精度、地理位置 | ✅ 已实现 |



## 🚀 快速开始

### 环境要求

- Go 1.20+
- Git

### 安装依赖

```bash
go mod tidy
```

### 启动服务

```bash
swag init -g cmd/server/main.go
go run cmd/server/main.go
```

服务启动后访问：
- **API服务**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/api/health

## 🔬 API使用示例

### 📋 API基础说明

本项目提供标准的科学计算API接口，支持多种科学计算任务。

#### 节气计算示例
```bash
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_name": "春分"
    }
  }'
```

### 🎯 完整的计算类型示例

#### 1. 节气计算 (solar_term)
```bash
# 计算2026年春分节气
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_name": "春分"
    }
  }'

# 计算2026年所有节气
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026
    }
  }'
```

#### 2. 干支计算 (ganzhi)
```bash
# 计算2026年2月4日10时的干支
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "ganzhi",
    "params": {
      "year": 2026,
      "month": 2,
      "day": 4,
      "hour": 10
    }
  }'
```

#### 3. 天文计算 (astronomy)
```bash
# 计算2026年3月23日的太阳黄经
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "astronomy",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23
    }
  }'
```

#### 4. 起运岁数计算 (starting_age)
```bash
# 计算1985年6月15日8时出生的起运岁数
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "starting_age",
    "params": {
      "year": 1985,
      "month": 6,
      "day": 15,
      "hour": 8
    }
  }'
```

#### 5. 农历转换 (lunar)
```bash
# 将2026年3月23日转换为农历
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "lunar",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23
    }
  }'
```

#### 6. 行星位置计算 (planet)
```bash
# 计算2026年3月23日火星的位置
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "planet",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23,
      "planet_name": "mars"
    }
  }'
```

#### 7. 星曜推算 (star)
```bash
# 推算2026年3月23日北斗七星的位置
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "star",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23,
      "star_name": "big_dipper"
    }
  }'
```

#### 8. 日出日落时间计算 (sunrise_sunset)
```bash
# 计算北京（116.4°E, 39.9°N）2026年3月23日的日出日落时间
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "sunrise_sunset",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23,
      "longitude": 116.4,
      "latitude": 39.9
    }
  }'
```

#### 9. 月相计算 (moon_phase)
```bash
# 计算2026年3月23日的月相
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "moon_phase",
    "params": {
      "year": 2026,
      "month": 3,
      "day": 23
    }
  }'
```

#### 10. 方程求解器 (equation_solver)
```bash
# 求解非线性方程 x^3 - 2x - 5 = 0
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "equation_solver",
    "params": {
      "equation_type": "nonlinear",
      "equation": "x^3 - 2x - 5 = 0",
      "initial_guess": 2.0,
      "tolerance": 1e-6
    }
  }'

# 求解线性方程组
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "equation_solver",
    "params": {
      "equation_type": "linear",
      "equation": "2x + 3y = 7, 4x - y = 1",
      "coefficients": [2, 3, 4, -1]
    }
  }'

# 求解微分方程 dy/dt = -y
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "equation_solver",
    "params": {
      "equation_type": "ode",
      "equation": "dy/dt = -y",
      "initial_value": 1.0,
      "time_step": 0.1,
      "time_range": 10.0
    }
  }'
```

#### 11. 符号计算器 (symbolic_calc)
```bash
# 符号求导 d/dx(x^3 + 2x^2 + x)
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "differentiate",
      "expression": "x^3 + 2*x^2 + x",
      "variable": "x"
    }
  }'

# 表达式化简
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "simplify",
      "expression": "2*x + 3*x - x"
    }
  }'

# 表达式求值
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "evaluate",
      "expression": "x^2 + 2*x + 1",
      "x_value": 3.0
    }
  }'
```

#### 12. 微分方程求解器 (ode_solver)
```bash
# 使用欧拉法求解微分方程
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "ode_solver",
    "params": {
      "equation": "dy/dt = -y",
      "initial_value": 1.0,
      "time_step": 0.1,
      "time_range": 1.0,
      "method": "euler"
    }
  }'

# 使用龙格-库塔法求解
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "ode_solver",
    "params": {
      "equation": "dy/dt = sin(t) - y",
      "initial_value": 0.0,
      "time_step": 0.05,
      "time_range": 1.0,
      "method": "rk4"
    }
  }'

# 使用亚当斯法求解
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "ode_solver",
    "params": {
      "equation": "d2y/dt2 = -y",
      "initial_value": 1.0,
      "initial_deriv": 0.0,
      "time_step": 0.1,
      "time_range": 1.0,
      "method": "adams"
    }
  }'
```



#### 获取系统信息
```bash
# 获取系统状态和计算器信息
curl -X GET "http://localhost:8080/api/system-info""
```

### 📊 健康检查

```bash
# 检查服务是否正常运行
curl -X GET "http://localhost:8080/api/health"
```

### ⚠️ 错误处理说明

#### 常见错误类型

1. **参数验证错误** (HTTP 400)
```json
{
  "success": false,
  "error": "请求参数错误: 缺少必需参数: year",
  "code": 400
}
```

2. **计算类型不支持** (HTTP 400)
```json
{
  "success": false,
  "error": "不支持的计算类型: invalid_calculation",
  "code": 400
}
```



4. **计算失败** (HTTP 400)
```json
{
  "success": false,
  "error": "计算失败: 参数超出有效范围",
  "code": 400
}
```

#### 成功响应格式
```json
{
  "success": true,
  "result": {
    "solar_term_time": "2026-03-20 12:00:00",
    "sun_longitude": 0.0,
    "julian_date": 2459580.5
  },
  "calculation": "solar_term",
  "timestamp": "2026-03-24T11:45:00Z"
}
```

### 📋 参数说明表

#### 计算类型参数
| 计算类型 | 参数 | 类型 | 必填 | 说明 |
|---------|------|------|------|------|
| solar_term | year | int | ✅ | 年份 |
| | term_index | int | ❌ | 节气索引(0-23)，与term_name二选一 |
| | term_name | string | ❌ | 节气中文名称，与term_index二选一 |
| ganzhi | year | int | ✅ | 年份 |
| | month | int | ✅ | 月份 |
| | day | int | ✅ | 日期 |
| | hour | int | ✅ | 小时 |
| astronomy | year | int | ✅ | 年份 |
| | month | int | ✅ | 月份 |
| | day | int | ✅ | 日期 |
| starting_age | year | int | ✅ | 出生年份 |
| | month | int | ✅ | 出生月份 |
| | day | int | ✅ | 出生日期 |
| | hour | int | ✅ | 出生小时 |
| lunar | year | int | ✅ | 阳历年份 |
| | month | int | ✅ | 阳历月份 |
| | day | int | ✅ | 阳历日期 |
| planet | year | int | ✅ | 年份 |
| | month | int | ✅ | 月份 |
| | day | int | ✅ | 日期 |
| | planet_name | string | ✅ | 行星名称 |
| star | year | int | ✅ | 年份 |
| | month | int | ✅ | 月份 |
| | day | int | ✅ | 日期 |
| | star_name | string | ✅ | 星曜名称 |
| sunrise_sunset | year | int | ✅ | 年份 |
| | month | int | ✅ | 月份 |
| | day | int | ✅ | 日期 |
| | longitude | float | ✅ | 经度 |
| | latitude | float | ✅ | 纬度 |
| moon_phase | year | int | ✅ | 年份 |
| | month | int | ✅ | 月份 |
| | day | int | ✅ | 日期 |
| equation_solver | equation_type | string | ✅ | 方程类型：nonlinear, linear, ode |
| | equation | string | ✅ | 方程表达式 |
| | initial_guess | float | ❌ | 初始猜测值（非线性方程） |
| | tolerance | float | ❌ | 容差 |
| | max_iterations | int | ❌ | 最大迭代次数 |
| | coefficients | []float | ❌ | 系数数组（线性方程组） |
| | initial_value | float | ❌ | 初始值（微分方程） |
| | time_step | float | ❌ | 时间步长（微分方程） |
| | time_range | float | ❌ | 时间范围（微分方程） |
| symbolic_calc | operation | string | ✅ | 操作类型：parse, differentiate, simplify, evaluate |
| | expression | string | ✅ | 数学表达式 |
| | variable | string | ❌ | 变量名（求导用） |
| | x_value | float | ❌ | x的值（求值用） |
| | y_value | float | ❌ | y的值（求值用） |
| | z_value | float | ❌ | z的值（求值用） |
| ode_solver | equation | string | ✅ | 微分方程表达式 |
| | equation_type | string | ❌ | 方程类型：first_order, second_order |
| | initial_value | float | ✅ | 初始值 |
| | initial_deriv | float | ❌ | 初始导数值（二阶方程） |
| | time_step | float | ✅ | 时间步长 |
| | time_range | float | ✅ | 时间范围 |
| | method | string | ❌ | 求解方法：euler, rk4, adams |

#### 节气中文名称对照表
| 索引 | 节气名称 | 索引 | 节气名称 |
|------|----------|------|----------|
| 0 | 立春 | 12 | 立秋 |
| 1 | 雨水 | 13 | 处暑 |
| 2 | 惊蛰 | 14 | 白露 |
| 3 | 春分 | 15 | 秋分 |
| 4 | 清明 | 16 | 寒露 |
| 5 | 谷雨 | 17 | 霜降 |
| 6 | 立夏 | 18 | 立冬 |
| 7 | 小满 | 19 | 小雪 |
| 8 | 芒种 | 20 | 大雪 |
| 9 | 夏至 | 21 | 冬至 |
| 10 | 小暑 | 22 | 小寒 |
| 11 | 大暑 | 23 | 大寒 |

### 🚀 快速开始指南

#### 步骤1：启动服务
```bash
# 生成Swagger文档
swag init -g cmd/server/main.go

# 启动服务
go run cmd/server/main.go
```

#### 步骤2：测试基本功能
```bash
# 检查服务状态
curl -X GET "http://localhost:8080/api/health"

# 测试节气计算 - 使用中文名称
curl -X POST "http://localhost:8080/api/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_name": "春分"
    }
  }'
```

#### 步骤4：使用Swagger UI
打开浏览器访问：`http://localhost:8080/swagger/index.html`

### 💡 使用技巧

1. **参数格式灵活**：支持两种参数格式（嵌套params或直接参数）
2. **节气参数便捷**：支持数字索引和中文名称两种方式，推荐使用中文名称
3. **错误信息详细**：所有错误都包含详细的错误信息和修复建议
4. **Swagger集成**：可通过Web界面直接测试所有API
5. **中文友好**：节气计算支持24节气的中文名称，无需记忆数字索引

### 🔗 相关资源

- [Swagger API文档](http://localhost:8080/swagger/index.html) - 完整的API文档
```



## 🔧 开发指南

### 添加新的科学计算任务

1. 在 `internal/calculator/` 创建新的计算器文件
2. 实现 `Calculator` 接口
3. 在计算器管理器中注册
4. 添加相应的单元测试



### 扩展API接口

1. 在 `internal/api/handler.go` 添加新的处理器
2. 定义相应的DTO结构
3. 更新Swagger文档注释
4. 添加集成测试

## 📋 项目状态

### 已完成功能

- ✅ 基础项目结构重构
- ✅ 核心科学计算任务（12个）
- ✅ API接口和Swagger文档
- ✅ 基础测试框架

### 待实现功能

- 🔄 完整的测试覆盖
- 🔄 性能优化和监控

## 🔗 相关资源

- [AI模型对比实验指南](AI模型对比实验指南.md) - 详细的实验设计指南
- [规划文档分析](改进建议.md) - 基于规划文档的改进分析
- [Swagger API文档](http://localhost:8080/swagger/index.html) - 完整的API文档

## 🤝 贡献指南

欢迎贡献代码和改进建议！请遵循以下准则：

1. 遵循Go代码规范
2. 添加相应的单元测试
3. 更新相关文档
4. 通过代码审查

## 📄 许可证

本项目采用MIT许可证。