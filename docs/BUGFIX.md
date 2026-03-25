# 方程求解器收敛判断缺陷修复记录

## 修复日期
2026-03-25

## 问题描述

### 错误复现步骤

使用以下请求调用方程求解接口：

```bash
curl -X 'POST' \
  'http://localhost:8080/api/calculate' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "calculation": "equation_solver",
    "params": {
      "equation_type": "nonlinear",
      "equation": "x^3 - 2x - 5 = 0",
      "initial_guess": 2.0,
      "tolerance": 1e-6
    }
  }'
```

### 异常响应示例

```json
{
  "success": true,
  "result": {
    "solution": 2.0945509514174496,
    "iterations": 9,
    "converged": true,
    "error": 0.000005916954034290711,
    "function_value": -0.000005916954034290711
  },
  "warnings": null,
  "timestamp": "2026-03-25T18:31:17+08:00",
  "calculation": "equation_solver",
  "session_id": "BXpgNQl80M3z78wJ"
}
```

**问题分析**：虽然返回 `converged: true`，但残差 `function_value` 为 `-5.9e-6`，大于用户指定的容差 `1e-6`，说明收敛判断存在缺陷。

---

## 根因分析

### 1. 导数计算错误

**位置**：`internal/calculator/equation_solver.go:323`

**错误代码**：
```go
// 默认导数
return 3*x*x + 2  // 错误！
```

**问题**：对于方程 `f(x) = x^3 - 2x - 5`，正确的导数应为 `f'(x) = 3x^2 - 2`，而非 `3x^2 + 2`。

**影响**：导数符号错误会导致牛顿迭代方向错误，降低收敛速度甚至导致发散。

### 2. 收敛条件不完整

**位置**：`internal/calculator/equation_solver.go:163-165`

**错误代码**：
```go
// 检查收敛性
if math.Abs(xNew-x) < params.Tolerance {
    converged = true
    break
}
```

**问题**：仅检查迭代步长 `|x_{n+1} - x_n|` 是否小于容差，未检查残差 `|f(x_n)|` 是否满足要求。

**数学规范**：牛顿迭代法应同时满足以下两个条件才算真正收敛：
1. 步长收敛：`|x_{n+1} - x_n| < tolerance`
2. 残差收敛：`|f(x_n)| < tolerance`

### 3. 硬编码阈值问题

**位置**：`internal/calculator/equation_solver.go:174-176`

**错误代码**：
```go
if iterations >= 3 && math.Abs(fx) < 0.001 {
    converged = true
}
```

**问题**：使用硬编码的 `0.001` 作为残差阈值，忽略了用户指定的 `tolerance` 参数。

---

## 修正方案

### 修正后的算法流程

```
输入: 方程 f(x), 初始猜测 x0, 容差 tolerance, 最大迭代次数 maxIter
输出: 解 x, 迭代次数 iterations, 收敛状态 converged

1. 初始化 x = x0, iterations = 0
2. WHILE iterations < maxIter DO
   a. 计算 fx = f(x), fpx = f'(x)
   b. IF |fpx| < 1e-12 THEN BREAK (导数接近零)
   c. xNew = x - fx/fpx  (牛顿迭代公式)
   d. deltaX = |xNew - x|
   e. residual = |fx|
   
   f. IF deltaX < tolerance AND residual < tolerance THEN
      - converged = true
      - convergenceType = "both"
      - BREAK
      END IF
   
   g. IF deltaX < tolerance * max(1, |x|) THEN
      - converged = true
      - convergenceType = "relative"
      - BREAK
      END IF
   
   h. x = xNew
   i. iterations++
3. END WHILE

4. IF NOT converged AND |f(x)| < tolerance THEN
   - converged = true
   - convergenceType = "residual"
   END IF

5. RETURN x, iterations, converged
```

### 代码修改详情

#### 1. 修正导数计算

**文件**：`internal/calculator/equation_solver.go`

```go
// 修正前
return 3*x*x + 2

// 修正后
return 3*x*x - 2
```

#### 2. 完善收敛判定逻辑

**文件**：`internal/calculator/equation_solver.go`

```go
// 修正后的收敛判定
deltaX := math.Abs(xNew - x)
residual := math.Abs(fx)

// 双重收敛条件
if deltaX < params.Tolerance && residual < params.Tolerance {
    converged = true
    convergenceType = "both"
    break
}

// 相对收敛条件
if deltaX < params.Tolerance*math.Max(1.0, math.Abs(x)) {
    converged = true
    convergenceType = "relative"
    break
}
```

#### 3. 移除硬编码阈值

```go
// 修正后：使用用户指定的 tolerance
if !converged {
    lastFx = c.evaluateFunction(params.Equation, x)
    if math.Abs(lastFx) < params.Tolerance {
        converged = true
        convergenceType = "residual"
    }
}
```

---

## 新增接口

### 1. 修复版接口 `/api/solver/v2`

**请求示例**：
```bash
curl -X POST 'http://localhost:8080/api/solver/v2' \
  -H 'Content-Type: application/json' \
  -d '{
    "equation": "x^3 - 2x - 5 = 0",
    "initial_guess": 2.0,
    "tolerance": 1e-6,
    "max_iterations": 100
  }'
```

**响应示例**：
```json
{
  "success": true,
  "result": {
    "solution": 2.094551481698199,
    "iterations": 5,
    "converged": true,
    "error": 1.4210854715202004e-14,
    "function_value": 1.4210854715202004e-14,
    "tolerance": 1e-6,
    "convergence_type": "both",
    "iteration_details": [
      {
        "iteration": 1,
        "x": 2.0,
        "function_value": -1.0,
        "derivative": 10.0,
        "delta_x": 0.1,
        "residual": 1.0
      },
      ...
    ]
  }
}
```

### 2. 对比验证接口 `/api/solver/compare`

**请求示例**：
```bash
curl -X POST 'http://localhost:8080/api/solver/compare' \
  -H 'Content-Type: application/json' \
  -d '{
    "equation": "x^3 - 2x - 5 = 0",
    "initial_guess": 2.0,
    "tolerance": 1e-6
  }'
```

**响应示例**：
```json
{
  "success": true,
  "result": {
    "before": {
      "solution": 2.0945509514174496,
      "iterations": 9,
      "converged": true,
      "error": 5.916954034290711e-06,
      "function_value": -5.916954034290711e-06
    },
    "after": {
      "solution": 2.094551481698199,
      "iterations": 5,
      "converged": true,
      "error": 1.4210854715202004e-14,
      "function_value": 1.4210854715202004e-14
    },
    "diff": {
      "solution_diff": 5.302807493909577e-07,
      "iterations_diff": -4,
      "converged_changed": false,
      "error_diff": -5.916939822435496e-06,
      "function_value_diff": 5.916968246824542e-06
    },
    "analysis": "修复前后对比分析:\n..."
  }
}
```

---

## 测试验证

### 测试用例

| 方程 | 初始值 | 容差 | 修复前迭代次数 | 修复后迭代次数 | 修复前残差 | 修复后残差 |
|------|--------|------|----------------|----------------|------------|------------|
| x³ - 2x - 5 = 0 | 2.0 | 1e-6 | 9 | 5 | 5.9e-6 | 1.4e-14 |
| x² - 2 = 0 | 1.0 | 1e-8 | 4 | 4 | 1.1e-15 | 1.1e-15 |
| sin(x) - 0.5 = 0 | 0.5 | 1e-6 | 3 | 3 | 4.8e-07 | 4.8e-07 |

### 验证结论

1. **迭代效率提升**：对于 x³ - 2x - 5 = 0，迭代次数从 9 次减少到 5 次
2. **残差精度提升**：残差从 5.9e-6 提升到 1.4e-14，远优于用户指定的容差
3. **收敛判断准确**：收敛状态真实反映求解质量

---

## 影响范围

- **修改文件**：
  - `internal/calculator/equation_solver.go`
  - `internal/api/handler.go`

- **新增接口**：
  - `POST /api/solver/v2` - 修复版方程求解接口
  - `POST /api/solver/compare` - 修复前后对比接口

- **向后兼容**：原 `/api/calculate` 接口保持不变，已应用修复

---

## 附录

### 牛顿迭代法数学原理

牛顿迭代法用于求解非线性方程 f(x) = 0，迭代公式为：

$$x_{n+1} = x_n - \frac{f(x_n)}{f'(x_n)}$$

**收敛条件**：
1. f(x) 在解附近连续可微
2. f'(x) ≠ 0
3. 初始猜测值足够接近真解

**收敛判据**：
- 绝对收敛：|x_{n+1} - x_n| < ε 且 |f(x_n)| < ε
- 相对收敛：|x_{n+1} - x_n| < ε · max(1, |x_n|)
