# 科学计算项目

基于规划文档先进设计理念重构的Go科学计算项目，提供丰富的科学计算功能。

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

## 🐛 Bug修复说明

### 日出日落时间计算Bug修复

#### Bug现象描述

在调用 `http://localhost:8080/api/calculate` 接口计算北京（116.4°E, 39.9°N）2026年3月23日的日出日落时间时，返回结果存在严重的时间逻辑错误：

**错误返回示例：**
```json
{
  "result": {
    "sunrise": "18:10",    // 错误：日出时间显示为傍晚
    "sunset": "06:18",     // 错误：日落时间显示为清晨
    "solar_noon": "00:14"  // 错误：正午时间显示为午夜
  }
}
```

**问题特征：**
- 日出日落时间完全颠倒
- 时间范围完全不符合正常逻辑
- 晨昏蒙影时间也存在同样的颠倒问题

---

#### 根本原因分析

通过分析 `internal/calculator/sunrise_sunset.go` 代码，发现以下核心问题：

**1. 太阳赤纬值硬编码错误**
```go
// 错误代码（第349-351行）：赤纬硬编码为0
cosH := (math.Sin(h0*math.Pi/180) -
    math.Sin(latitude*math.Pi/180)*math.Sin(0*math.Pi/180)) /
    (math.Cos(latitude*math.Pi/180) * math.Cos(0*math.Pi/180))
```
- 赤纬（Declination）应该根据日期动态计算
- 硬编码为0导致时角计算完全错误

**2. 缺少真太阳时修正（时差方程）**
```go
// 错误代码（第360-361行）：缺少太阳赤经修正
sunriseJD := jd + (720-4*(longitude+H)-0)/1440.0
sunsetJD := jd + (720-4*(longitude-H)-0)/1440.0
```
- 公式中的最后一个 `0` 应为太阳赤经修正值
- 缺少这一修正导致时间计算偏差12小时以上

**3. 太阳正午计算公式错误**
```go
// 错误公式
noonUTC := jd + (longitude*4 - E) / 1440.0

// 正确公式（NOAA标准算法）
noonMinutesUTC := 720 - 4*longitude - E  // 720分钟=正午12点
noonUTC := jd + noonMinutesUTC / 1440.0
```
- 公式符号错误导致时间偏差4-8小时
- 是导致日出日落时间整体偏移的根本原因

**4. 儒略日基准转换错误**
- NOAA算法要求使用午夜基准的儒略日（xxxx.0表示午夜）
- 计算未正确区分中午基准与午夜基准
- 时间格式化时儒略日小数部分处理逻辑错误

---

#### 验证方法及效果

**新增接口：**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/calculate-fixed` | POST | 修复后的日出日落计算接口 |
| `/api/calculate-compare` | POST | 原接口与修复后接口的结果对比 |

**使用修复接口验证：**
```bash
curl -X 'POST' \
  'http://localhost:8080/api/calculate-fixed' \
  -H 'Content-Type: application/json' \
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

**修复后返回结果（示例，2026年3月25日验证）：**
```json
{
  "success": true,
  "result": {
    "date": "2026-03-23",
    "sunrise": "06:14",    // 正确：清晨日出
    "sunset": "18:28",     // 正确：傍晚日落
    "solar_noon": "12:21", // 正确：中午正午
    "day_length": 12.23
  }
}
```

**使用对比接口验证：**
```bash
curl -X 'POST' \
  'http://localhost:8080/api/calculate-compare' \
  -H 'Content-Type: application/json' \
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

**对比结果示例：**
```json
{
  "success": true,
  "differences": {
    "sunrise": {
      "original": "18:10",
      "fixed": "06:14"
    },
    "sunset": {
      "original": "06:18",
      "fixed": "18:28"
    },
    "solar_noon": {
      "original": "00:14",
      "fixed": "12:21"
    }
  },
  "summary": "修复了以下问题：[日出时间 日落时间 正午时间]"
}
```

---

#### 修复后的算法说明

**修复后的计算器文件：** `internal/calculator/sunrise_sunset_fixed.go`

**采用的算法：** **NOAA日出日落计算方法**（美国国家海洋和大气管理局标准算法）

**算法步骤：**

1. **计算儒略日（Julian Day）**
   - 正确转换公历日期到儒略日
   - 基准点为UTC中午12点

2. **计算儒略世纪（Julian Century）**
   - 以J2000.0历元为基准
   - 用于后续的太阳位置计算

3. **太阳位置计算**
   - 太阳几何平黄经（GMLS）
   - 太阳平近点角（MNAS）
   - 太阳中心差（Equation of Center）
   - 太阳真黄经和赤纬

4. **时差方程（Equation of Time）**
   - 计算真太阳时与平太阳时的差值
   - 修正地球轨道椭圆和黄赤交角的影响

5. **时角计算（Hour Angle）**
   - 基于观测者纬度和太阳赤纬
   - 考虑大气折射修正值（0.8333度）
   - 支持海拔高度修正

6. **日出日落UTC时间计算**
   - 太阳正午时间精确计算
   - 日出 = 正午 - 半昼弧
   - 日落 = 正午 + 半昼弧

7. **时区转换**
   - UTC时间转当地时区
   - 正确处理24小时制的取模运算

**主要修复点（最终版本）：**
```go
// 1. 儒略日基准转换：从中午基准转为午夜基准（NOAA算法要求）
jdMidnight := jd - 0.5  // xxxx.5中午 -> xxxx.0午夜

// 2. 使用动态计算的太阳赤纬
delta := math.Asin(math.Sin(epsilon*math.Pi/180)*math.Sin(longitudeSun*math.Pi/180)) * 180 / math.Pi

// 3. 实现时差方程修正
E := (L - alpha) * 4  // 时差（分钟）

// 4. 正确的太阳正午计算（NOAA标准公式）
noonMinutesUTC := 720 - 4*longitude - E  // 720分钟=正午12点
noonUTC := jdMidnight + noonMinutesUTC / 1440.0

// 5. 日出日落时间计算
sunriseUTC := noonUTC - H*4/1440.0  // H为时角（度）
sunsetUTC := noonUTC + H*4/1440.0

// 6. 正确的时区转换（午夜基准）
dayFraction := jd - math.Floor(jd)  // .0=午夜, .5=中午
hoursUTC := dayFraction * 24        // UTC小时
hours := hoursUTC + timezone        // 当地时间
```

**参考文献：**
- NOAA Solar Calculator: https://www.esrl.noaa.gov/gmd/grad/solcalc/
- Astronomical Algorithms by Jean Meeus (第2版)

---

## 📋 项目状态

### 已完成功能

- ✅ 基础项目结构重构
- ✅ 核心科学计算任务（12个）
- ✅ API接口和Swagger文档
- ✅ 基础测试框架
- ✅ 日出日落时间Bug修复（2026-03-25）

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