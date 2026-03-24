我们理解您的需求：需要规划一个基于 Go 的科学计算项目，该项目预先植入了三类典型 Bug（结果不稳定性、约束越界、精度错误），用于后续三个 AI 模型的修复能力对比实验。规划应涵盖项目结构、核心功能、Bug 植入方式、API 设计、实验流程等，但不涉及具体代码实现。

---

## 1. 项目概述

### 1.1 项目名称
`scientific-calc-bugs`（暂定）

### 1.2 项目目标
提供一套包含 8 个核心科学计算任务（源自 `solar_terms.py`）的 Go 服务，每个任务均植入三类可开关的 Bug，用于评估 AI 模型在 Bug 识别、修复、代码质量等方面的能力。

### 1.3 技术栈
- **语言**：Go 1.20+  
- **Web 框架**：Gin / Echo（轻量级，便于 API 开发）  
- **API 文档**：Swagger（通过 swaggo 自动生成）  
- **日志**：logrus 或标准库 slog  
- **测试**：Go 标准 testing 包 + testify 断言库  
- **依赖管理**：Go Modules  

---

## 2. 核心科学计算任务（8 个）

参考原始 `solar_terms.py`，规划以下 8 个科学计算任务（可根据实际天文算法调整）：

| 序号 | 任务名称 | 英文标识 | 说明 |
|------|----------|----------|------|
| 1 | 节气计算 | `solar_term` | 给定年份和节气索引，返回精确到秒的节气时间 |
| 2 | 干支计算 | `ganzhi` | 给定日期，返回该日的干支（天干地支） |
| 3 | 天文黄经计算 | `astronomy` | 给定儒略日，返回太阳黄经（精度要求 0.000001°） |
| 4 | 农历日期转换 | `lunar` | 公历转农历，含闰月处理 |
| 5 | 星曜推算 | `star` | 计算某日的二十八宿值日 |
| 6 | 日出日落时间 | `sunrise_sunset` | 给定经纬度、日期，返回日出/日落时刻 |
| 7 | 月相计算 | `moon_phase` | 给定日期，返回月相（新月、上弦等）及准确时间 |
| 8 | 行星位置计算 | `planet` | 计算某行星（如木星）的赤经/赤纬 |

每个任务均提供 **有 Bug 版本** 和 **无 Bug（正确）版本**，实验中通过参数选择是否启用 Bug。

---

## 3. 项目结构规划

```
scientific-calc-bugs/
├── cmd/
│   └── server/
│       └── main.go               # 服务入口，路由注册
├── internal/
│   ├── api/
│   │   ├── handler.go            # HTTP 处理器
│   │   └── dto.go                # 请求/响应数据结构
│   ├── calculator/
│   │   ├── solar_term.go         # 节气计算（含正确版 + Bug 版）
│   │   ├── ganzhi.go             # 干支计算
│   │   ├── astronomy.go          # 黄经计算
│   │   ├── lunar.go              # 农历转换
│   │   ├── star.go               # 星曜计算
│   │   ├── sunrise_sunset.go     # 日出日落
│   │   ├── moon_phase.go         # 月相计算
│   │   └── planet.go             # 行星位置
│   ├── bugs/
│   │   ├── instability.go        # 结果不稳定性 Bug 的实现（注入随机种子/迭代不足）
│   │   ├── constraint.go         # 约束越界 Bug 的实现（边界未检查/返回非法值）
│   │   └── precision.go          # 精度错误 Bug 的实现（浮点截断/过早终止迭代）
│   └── utils/
│       ├── math.go               # 公用数学函数（如角度归一化）
│       └── time.go               # 时间处理辅助函数
├── pkg/
│   └── models/
│       └── request.go            # 统一的请求结构（含 bug_type 字段）
├── tests/
│   ├── integration/              # 集成测试（针对 API）
│   ├── unit/                     # 单元测试（针对每个计算函数）
│   └── testdata/                 # 测试数据（如已知节气时间）
├── docs/                         # Swagger 生成的文档
├── scripts/                      # 实验辅助脚本（如批量调用 curl）
├── go.mod
├── go.sum
├── README.md                     # 项目说明、快速启动
├── EXPERIMENT.md                 # 实验指南（基于提供的实验指南.md）
└── .gitignore
```

---

## 4. Bug 植入设计

### 4.1 统一的 Bug 开关机制

- 每个计算函数接收一个可选的 `BugType` 参数（枚举：`none`, `instability`, `constraint`, `precision`）。
- 在函数内部，根据 BugType 决定调用“正确实现”还是“Bug 实现”。
- Bug 实现通过组合或替换部分逻辑完成，尽量保持与正确实现相似的代码结构，便于 AI 识别。

### 4.2 三类 Bug 具体植入方式

| Bug 类型 | 英文标识 | 植入手法示例 |
|----------|----------|--------------|
| **结果不稳定性** | `instability` | - 使用未初始化的随机数种子（每次调用结果不同）<br>- 迭代算法中故意减少迭代次数（如二分法只迭代 5 次）<br>- 依赖全局状态且被并发修改 |
| **约束越界** | `constraint` | - 不检查输入边界（如年份为负数时仍计算）<br>- 计算结果直接返回超出合理范围的值（如干支返回“无”、月相返回空字符串）<br>- 对无效输入不报错，返回虚假结果 |
| **精度错误** | `precision` | - 浮点数使用 float32 而非 float64<br>- 截断小数位（如只保留 2 位小数）<br>- 迭代算法未达到收敛精度就提前返回（如角度误差 >0.1°） |

### 4.3 代码示例结构（伪代码）

```go
// calculator/solar_term.go
func CalculateSolarTerm(year int, termIndex int, bugType string) (string, error) {
    switch bugType {
    case "instability":
        return instabilityBugSolarTerm(year, termIndex)
    case "constraint":
        return constraintBugSolarTerm(year, termIndex)
    case "precision":
        return precisionBugSolarTerm(year, termIndex)
    default:
        return correctSolarTerm(year, termIndex)
    }
}
```

各 Bug 函数可置于 `bugs/` 包中，但为了方便 AI 分析，建议将 Bug 实现与正确实现放在同一个文件内，通过条件分支区分，这样 AI 可以看到完整上下文。

---

## 5. API 设计

### 5.1 通用端点

```
POST /api/calculate-with-bugs
```

### 5.2 请求体

```json
{
  "calculation": "solar_term",    // 对应8个任务之一的标识
  "bug_type": "instability",      // none / instability / constraint / precision
  "params": {                     // 各任务所需的参数
    "year": 2024,
    "term_index": 2,
    ...
  }
}
```

### 5.3 响应体

```json
{
  "success": true,
  "result": "2024-02-04 12:03:15",
  "warnings": ["启用结果不稳定Bug模式"]
}
```

若计算失败（如参数非法），返回 4xx 并包含错误信息。

### 5.4 辅助端点

- `GET /api/health`：健康检查
- `GET /api/swagger/*`：Swagger UI
- `GET /api/tasks`：返回所有支持的 calculation 类型及其所需参数说明

### 5.5 参数映射表

为每个任务设计独立的结构体，统一通过 `params` 字段传递，由后端根据 `calculation` 解析。

---

## 6. 实验流程支持

### 6.1 启动服务

```bash
go run cmd/server/main.go
```

默认监听 `8080` 端口。

### 6.2 验证环境

提供 `scripts/verify.sh` 脚本，自动调用 `/api/health` 和 `/api/tasks` 验证服务正常。

### 6.3 批量测试脚本

提供 `scripts/run_experiment.sh` 或 Python 脚本，循环调用 API 并记录结果，方便对比。

### 6.4 测试数据

在 `tests/testdata/` 下提供每个任务的标准答案（基于正确版本计算或权威数据），用于验证修复后的结果。

---

## 7. 核心计算任务实现要点

### 7.1 节气计算
- 基于天文算法（如 Jean Meeus 的《天文算法》）计算太阳黄经达到特定值的时间。
- 正确版本：使用高精度迭代（如牛顿迭代），误差小于 1 秒。
- Bug 版本：使用二分法但迭代次数不足，导致结果不稳定或精度差。

### 7.2 干支计算
- 公历日期转天干地支，需考虑年柱以立春为界。
- 正确版本：严格按天文年历规则。
- Bug 版本：不处理立春边界，或对负数年份直接返回“无”。

### 7.3 天文黄经计算
- 基于 VSOP87 理论计算太阳黄经。
- 正确版本：使用高精度多项式，返回 float64 精度。
- Bug 版本：使用 float32，或截断高次项。

### 7.4 农历转换
- 采用天文算法计算新月时刻，确定农历日期。
- 正确版本：处理闰月、大小月。
- Bug 版本：忽略闰月或错误计算新月。

### 7.5 星曜推算
- 根据日干支确定二十八宿值日。
- 正确版本：按历书规则。
- Bug 版本：随机返回星宿。

### 7.6 日出日落时间
- 根据经纬度、儒略日计算。
- 正确版本：考虑大气折射，精度至分钟。
- Bug 版本：忽略折射，或使用错误公式。

### 7.7 月相计算
- 计算给定日期的月相及精确时间。
- 正确版本：基于高精度新月算法。
- Bug 版本：使用近似周期，误差可达一天。

### 7.8 行星位置计算
- 采用 VSOP87 或简化模型计算赤经赤纬。
- 正确版本：高精度。
- Bug 版本：忽略岁差或使用低精度模型。

---

## 8. 实验数据与评估

- **基准数据**：使用正确版本的计算结果作为基准，或者引用外部权威数据（如紫金山天文台发布的节气时刻）。
- **评估指标**：参考《AI模型对比实验指南.md》中的量化评估体系，在项目内提供评分脚本（可选），辅助人工评估。

---

## 9. 扩展性与注意事项

- **可配置性**：Bug 类型可通过环境变量或配置文件全局开启，方便批量测试。
- **日志**：记录每次请求使用的 Bug 类型和计算结果，便于复现。
- **并发安全**：正确版本应避免使用全局可变状态；Bug 版本可有意识引入竞态条件以测试 AI 修复能力。
- **文档**：每个计算函数需有详细注释，说明算法来源和预期行为，帮助 AI 理解。

---

## 10. 项目交付清单

1. **源代码**：完整项目代码（符合上述结构）。
2. **API 文档**：Swagger UI 可访问。
3. **实验指南**：更新 `EXPERIMENT.md`，融入《AI模型对比实验指南.md》的内容，并说明如何使用本项目的 API 进行实验。
4. **测试脚本**：提供快速验证和批量调用的脚本。
5. **示例请求/响应**：在 README 中给出各任务的 curl 示例。
6. **基准数据**：提供 `tests/testdata/` 中的正确答案，用于验证修复。

---

通过以上规划，我们构建了一个结构清晰、Bug 植入可控、便于 AI 模型进行对比实验的 Go 项目基础。后续开发可依据此规划逐步实现，确保实验的顺利进行。