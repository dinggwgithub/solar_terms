# 科学计算Bug测试项目

基于规划文档先进设计理念重构的Go科学计算项目，用于AI模型在科学计算领域的质量对比实验。

## 🎯 项目概述

本项目基于 `solar_terms.py` 科学计算任务，采用**模块化架构设计**，实现了完整的Bug管理系统和科学计算框架。

### 🔧 核心改进

1. **模块化架构** - 采用 `cmd/internal/pkg` 分层结构
2. **Bug系统分离** - 独立的Bug模块，支持扩展
3. **接口抽象** - 统一的计算器和Bug接口
4. **完整测试体系** - 单元测试和集成测试
5. **实验脚本支持** - 自动化实验流程

## 📊 项目结构（改进后）

```
scientific_calc_bugs/
├── cmd/
│   └── server/
│       └── main.go              # 服务入口（重构）
├── internal/
│   ├── api/
│   │   ├── handler.go           # HTTP处理器（重构）
│   │   └── dto.go               # 请求响应结构
│   ├── calculator/
│   │   ├── solar_term.go        # 节气计算
│   │   ├── ganzhi.go            # 干支计算
│   │   ├── astronomy.go         # 天文计算
│   │   ├── lunar.go             # 农历转换（新增）
│   │   ├── planet.go            # 行星位置（新增）
│   │   └── interface.go         # 计算器接口（新增）
│   ├── bugs/
│   │   ├── interface.go         # Bug接口（新增）
│   │   ├── manager.go           # Bug管理器（新增）
│   │   ├── instability.go       # 结果不稳定性Bug（重构）
│   │   ├── constraint.go        # 约束越界Bug（重构）
│   │   └── precision.go         # 精度错误Bug（重构）
│   └── utils/
│       ├── math.go              # 数学工具
│       └── time.go              # 时间工具
├── pkg/
│   └── models/
│       ├── request.go           # 请求结构
│       └── enums.go             # 枚举定义（新增）
├── tests/
│   ├── integration/             # 集成测试（新增）
│   ├── unit/                    # 单元测试（新增）
│   └── testdata/                # 测试数据（新增）
├── scripts/                     # 实验脚本（新增）
└── docs/                        # 文档
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

## 🐛 Bug管理系统

### Bug类型枚举

```go
type BugType int

const (
    BugTypeNone BugType = iota
    BugTypeInstability    // 结果不稳定性
    BugTypeConstraint     // 约束越界
    BugTypePrecision      // 精度错误
)
```

### Bug接口设计

```go
type Bug interface {
    Name() string                    // Bug名称
    Description() string             // Bug描述
    Apply(calculationType string, params interface{}) (interface{}, []string)
}
```

### Bug管理器功能

- **Bug注册管理** - 支持动态注册新的Bug类型
- **Bug应用控制** - 统一的Bug应用接口
- **Bug信息查询** - 获取Bug详细信息和修复建议
- **Bug验证** - 验证Bug应用结果的合理性

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

本项目提供两种API调用方式：

#### 方式1：通过查询参数指定Bug类型（推荐）
```bash
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_index": 2
    }
  }'
```

#### 方式2：通过请求体指定Bug类型
```bash
curl -X POST "http://localhost:8080/api/calculate-with-bugs" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "bug_type": "instability",
    "params": {
      "year": 2026,
      "term_index": 2
    }
  }'
```

**注意**：两种方式等效，推荐使用方式1（查询参数），因为更符合RESTful API设计规范。

**重要说明**：为了确保AI模型修复Bug时不会"作弊"，本项目**只提供带Bug的计算接口**。AI模型必须通过分析Bug代码来修复问题，而不是直接获取标准答案。

### 🎯 完整的计算类型示例

#### 1. 节气计算 (solar_term)
```bash
# 方式1：使用数字索引（原有方式）
# 计算2024年春分节气（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_index": 3
    }
  }'

# 方式2：使用中文名称（新增方式，推荐）
# 计算2024年春分节气（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_name": "春分"
    }
  }'

# 计算2024年所有节气（Bug类型：约束越界）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=constraint" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2024
    }
  }'
```

#### 2. 干支计算 (ganzhi)
```bash
# 计算2024年2月4日10时的干支（Bug类型：精度错误）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=precision" \
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
# 计算2024年3月20日的太阳黄经（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
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
# 计算1985年6月15日8时出生的起运岁数（Bug类型：约束越界）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=constraint" \
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
# 将2024年3月20日转换为农历（Bug类型：精度错误）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=precision" \
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
# 计算2024年3月20日火星的位置（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
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
# 推算2024年3月20日北斗七星的位置（Bug类型：约束越界）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=constraint" \
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
# 计算北京（116.4°E, 39.9°N）2024年3月20日的日出日落时间（Bug类型：精度错误）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=precision" \
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
# 计算2024年3月20日的月相（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
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
# 求解非线性方程 x^3 - 2x - 5 = 0（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
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

# 求解线性方程组（Bug类型：约束越界）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=constraint" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "equation_solver",
    "params": {
      "equation_type": "linear",
      "equation": "2x + 3y = 7, 4x - y = 1",
      "coefficients": [2, 3, 4, -1]
    }
  }'

# 求解微分方程 dy/dt = -y（Bug类型：精度错误）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=precision" \
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
# 符号求导 d/dx(x^3 + 2x^2 + x)（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "differentiate",
      "expression": "x^3 + 2*x^2 + x",
      "variable": "x"
    }
  }'

# 表达式化简（Bug类型：约束越界）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=constraint" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "simplify",
      "expression": "2*x + 3*x - x"
    }
  }'

# 表达式求值（Bug类型：精度错误）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=precision" \
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
# 使用欧拉法求解微分方程（Bug类型：结果不稳定）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
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

# 使用龙格-库塔法求解（Bug类型：约束越界）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=constraint" \
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

# 使用亚当斯法求解（Bug类型：精度错误）
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=precision" \
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

### 🐛 Bug管理API

#### 获取Bug信息
```bash
# 获取所有Bug类型信息
curl -X GET "http://localhost:8080/api/bug-info"

# 获取特定Bug类型信息
curl -X GET "http://localhost:8080/api/bug-info?bug_type=instability"
```

#### 获取系统信息
```bash
# 获取系统状态和计算器信息
curl -X GET "http://localhost:8080/api/system-info"
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

3. **Bug类型不支持** (HTTP 400)
```json
{
  "success": false,
  "error": "不支持的Bug类型: invalid_bug",
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
    "solar_term_time": "2024-03-20 12:00:00",
    "sun_longitude": 0.0,
    "julian_date": 2459580.5
  },
  "warnings": ["启用结果不稳定Bug模式"],
  "bug_type": "instability",
  "calculation": "solar_term",
  "timestamp": "2024-03-24T11:45:00Z"
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

#### Bug类型参数
| Bug类型 | 查询参数 | 说明 |
|---------|----------|------|
| 结果不稳定 | `bug_type=instability` | 计算结果存在随机波动 |
| 约束越界 | `bug_type=constraint` | 参数或结果超出有效范围 |
| 精度错误 | `bug_type=precision` | 计算精度不足或错误 |
| 无Bug | `bug_type=none` 或不指定 | 正确计算结果 |

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

# 测试节气计算（无Bug模式）- 使用中文名称
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=none" \
  -H "Content-Type: application/json" \
  -d '{
    "calculation": "solar_term",
    "params": {
      "year": 2026,
      "term_name": "春分"
    }
  }'
```

#### 步骤3：测试Bug功能
```bash
# 测试结果不稳定Bug - 使用中文名称
curl -X POST "http://localhost:8080/api/calculate-with-bugs?bug_type=instability" \
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
2. **Bug类型灵活**：支持查询参数和请求体两种方式指定
3. **节气参数便捷**：支持数字索引和中文名称两种方式，推荐使用中文名称
4. **错误信息详细**：所有错误都包含详细的错误信息和修复建议
5. **Swagger集成**：可通过Web界面直接测试所有API
6. **中文友好**：节气计算支持24节气的中文名称，无需记忆数字索引
7. **防作弊设计**：只提供带Bug的计算接口，确保AI模型必须通过分析代码修复问题

### 🔗 相关资源

- [AI模型对比实验指南](AI模型对比实验指南.md) - 详细的实验设计指南
- [Swagger API文档](http://localhost:8080/swagger/index.html) - 完整的API文档
```

## 🧪 实验支持

### 自动化测试脚本

```bash
# 运行结果不稳定性测试
./scripts/test_instability.sh

# 运行约束越界测试
./scripts/test_constraint.sh

# 运行精度错误测试
./scripts/test_precision.sh

# 运行完整实验套件
./scripts/run_experiment.sh
```

### 测试用例管理

项目提供完整的测试用例，包括：
- **单元测试** - 针对每个计算函数的测试
- **集成测试** - API接口的端到端测试
- **性能测试** - 计算性能基准测试
- **Bug验证测试** - Bug效果的验证测试

## 📊 评估指标体系

### Bug修复能力评估（60%权重）

| 指标 | 权重 | 评估标准 |
|------|------|----------|
| 结果稳定性 | 30% | 同一输入多次调用结果一致性 |
| 约束合规性 | 30% | 计算结果在合理范围内 |
| 精度准确性 | 40% | 与权威数据的误差程度 |

### 代码质量评估（40%权重）

| 指标 | 权重 | 评估标准 |
|------|------|----------|
| 可读性 | 25% | 代码结构、注释、命名规范 |
| 可维护性 | 25% | 模块化设计、错误处理 |
| 性能优化 | 25% | 计算效率、内存使用 |
| 测试覆盖 | 25% | 单元测试完整性 |

## 🔧 开发指南

### 添加新的科学计算任务

1. 在 `internal/calculator/` 创建新的计算器文件
2. 实现 `Calculator` 接口
3. 在计算器管理器中注册
4. 添加相应的单元测试

### 添加新的Bug类型

1. 在 `internal/bugs/` 创建新的Bug文件
2. 实现 `Bug` 接口
3. 在Bug管理器中注册
4. 添加Bug特征描述和修复建议

### 扩展API接口

1. 在 `internal/api/handler.go` 添加新的处理器
2. 定义相应的DTO结构
3. 更新Swagger文档注释
4. 添加集成测试

## 📋 项目状态

### 已完成功能

- ✅ 基础项目结构重构
- ✅ Bug管理系统实现
- ✅ 核心科学计算任务（4个）
- ✅ API接口和Swagger文档
- ✅ 基础测试框架

### 待实现功能

- 🔄 新增科学计算任务（4个）
- 🔄 完整的测试覆盖
- 🔄 实验脚本和工具
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

---

**注意**：本项目专为AI模型科学计算质量对比实验设计，所有Bug均为故意引入，用于评估AI模型的修复能力。