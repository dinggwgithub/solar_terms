# 【Kimi Km2.5】模型的符号微分计算错误修复文档

## 1. 问题概述

### 1.1 异常现象

针对请求：
```json
{
  "calculation": "symbolic_calc",
  "params": {
    "operation": "differentiate",
    "expression": "x^3 + 2*x^2 + x",
    "variable": "x"
  }
}
```

原版 API 返回结果：
```json
{
  "success": true,
  "result": {
    "original_expression": "x^3 + 2*x^2 + x",
    "result_expression": "d/dx(x^3 + 2*x^2 + x)",
    "parsed_tree": null,
    "derivative": "d/dx(x^3 + 2*x^2 + x)",
    "simplified": "",
    "numeric_value": 0,
    "variables": null,
    "operation_type": "differentiate"
  }
}
```

### 1.2 异常字段分析

| 字段 | 当前值 | 预期值 | 问题描述 |
|------|--------|--------|----------|
| `result_expression` | `"d/dx(x^3 + 2*x^2 + x)"` | `"3*x^2 + 4*x + 1"` | 仅返回占位符字符串，未实际计算 |
| `derivative` | `"d/dx(x^3 + 2*x^2 + x)"` | `"3*x^2 + 4*x + 1"` | 导数未实际计算 |
| `simplified` | `""` | `"3*x^2 + 4*x + 1"` | 化简结果为空 |
| `numeric_value` | `0` | 应代入变量值计算 | 数值求值恒为0 |
| `parsed_tree` | `null` | 应包含表达式解析树 | 解析树为空 |

### 1.3 预期正确结果

根据微积分基本规则：
- `d/dx(x^3) = 3*x^2`
- `d/dx(2*x^2) = 4*x`
- `d/dx(x) = 1`

所以 `d/dx(x^3 + 2*x^2 + x) = 3*x^2 + 4*x + 1`

---

## 2. 根本原因分析

### 2.1 代码问题定位

问题出在 `internal/calculator/symbolic_calc.go` 的以下函数：

#### 2.1.1 `calculateDerivative` 函数（第334-358行）

```go
func (c *SymbolicCalcCalculator) calculateDerivative(expr, variable string) string {
    // 常数导数
    if matched, _ := regexp.MatchString(`^\d+$`, expr); matched {
        return "0"
    }
    // 变量导数
    if expr == variable {
        return "1"
    }
    // 幂函数导数: d/dx(x^n) = n*x^(n-1)
    if match := regexp.MustCompile(`^` + variable + `\^(\d+)$`).FindStringSubmatch(expr); match != nil {
        // ... 仅处理简单幂函数
    }
    // 默认返回原始表达式（表示无法求导）
    return "d/d" + variable + "(" + expr + ")"  // <-- 问题所在
}
```

**问题**：对于复杂的多项式表达式（如 `x^3 + 2*x^2 + x`），函数直接返回占位符字符串，没有实现递归求导。

#### 2.1.2 `simplify` 函数（第361-385行）

```go
func (c *SymbolicCalcCalculator) simplify(expr string) string {
    // 去除多余空格
    expr = strings.ReplaceAll(expr, " ", "")
    // 简单的常数运算化简
    if strings.Contains(expr, "0*") || strings.Contains(expr, "*0") {
        return "0"
    }
    // ... 仅实现简单替换
    return expr
}
```

**问题**：没有实现多项式合并同类项、系数运算等复杂化简功能。

#### 2.1.3 `evaluate` 函数（第388-420行）

```go
func (c *SymbolicCalcCalculator) evaluate(expr string, variables map[string]float64) float64 {
    // 替换变量
    for varName, value := range variables {
        expr = strings.ReplaceAll(expr, varName, fmt.Sprintf("%.6f", value))
    }
    // 简单的数值计算
    if strings.Contains(expr, "+") {
        parts := strings.Split(expr, "+")
        // ... 简单分割处理
    }
    // 默认返回0
    return 0.0  // <-- 问题所在
}
```

**问题**：仅支持简单的加减乘幂运算，对于复杂表达式返回0。

#### 2.1.4 `buildParseTree` 函数（第279-295行）

```go
func (c *SymbolicCalcCalculator) buildParseTree(expr string) map[string]interface{} {
    tree := map[string]interface{}{
        "type":     "expression",
        "value":    expr,
        "tokens":   c.tokenize(expr),
        "operator": c.findMainOperator(expr),
    }
    return tree
}
```

**问题**：只是简单分词，没有构建真正的抽象语法树(AST)。

### 2.2 为何返回 `success: true`

代码中没有对计算失败情况进行检测和返回错误：

```go
func (c *SymbolicCalcCalculator) differentiateExpression(params *SymbolicParams) (*SymbolicResult, error) {
    derivative := c.calculateDerivative(params.Expression, params.Variable)
    
    result := &SymbolicResult{
        OriginalExpression: params.Expression,
        ResultExpression:   derivative,  // 即使是占位符也直接返回
        Derivative:         derivative,
        OperationType:      "differentiate",
    }
    // 所有路径都返回成功，没有错误检查
    return result, nil
}
```

---

## 3. 修复方案

### 3.1 修复策略

1. **实现完整的表达式解析器**：将表达式字符串解析为抽象语法树(AST)
2. **实现递归求导算法**：基于AST应用求导规则
3. **实现表达式化简**：合并同类项、常数运算等
4. **实现正确的数值求值**：基于AST进行变量替换和计算

### 3.2 关键代码实现

#### 3.2.1 AST节点定义

```go
// ASTNode 抽象语法树节点
type ASTNode struct {
    Type     string      // 节点类型: "number", "variable", "add", "sub", "mul", "div", "pow", "neg"
    Value    float64     // 数值（仅用于number类型）
    Name     string      // 变量名（仅用于variable类型）
    Children []*ASTNode  // 子节点
}
```

#### 3.2.2 表达式解析

```go
// parseToAST 将表达式字符串解析为AST
func (c *SymbolicCalcCalculatorFixed) parseToAST(expr string) *ASTNode {
    expr = strings.ReplaceAll(expr, " ", "")
    tokens := c.tokenize(expr)
    return c.parseExpressionTokens(tokens, 0, len(tokens)-1)
}

// 使用递归下降解析，处理运算符优先级
// 优先级: + - (最低) > * / > ^ (最高)
```

#### 3.2.3 符号求导（递归实现）

```go
// diffAST 对AST进行符号求导
func (c *SymbolicCalcCalculatorFixed) diffAST(node *ASTNode, variable string) *ASTNode {
    switch node.Type {
    case "number":
        return NewNumberNode(0)  // 常数导数为0
        
    case "variable":
        if node.Name == variable {
            return NewNumberNode(1)  // 变量导数为1
        }
        return NewNumberNode(0)
        
    case "add":
        // (u + v)' = u' + v'
        left := c.diffAST(node.Children[0], variable)
        right := c.diffAST(node.Children[1], variable)
        return NewBinaryNode("add", left, right)
        
    case "mul":
        // (u * v)' = u' * v + u * v'
        u := node.Children[0]
        v := node.Children[1]
        uPrime := c.diffAST(u, variable)
        vPrime := c.diffAST(v, variable)
        term1 := NewBinaryNode("mul", uPrime, v)
        term2 := NewBinaryNode("mul", u, vPrime)
        return NewBinaryNode("add", term1, term2)
        
    case "pow":
        // d/dx(x^n) = n * x^(n-1)
        base := node.Children[0]
        exponent := node.Children[1]
        if base.Type == "variable" && exponent.Type == "number" {
            n := exponent.Value
            newExp := NewNumberNode(n - 1)
            newPow := NewBinaryNode("pow", base, newExp)
            return NewBinaryNode("mul", NewNumberNode(n), newPow)
        }
        // 处理复合函数...
    }
    // ...
}
```

#### 3.2.4 多项式化简

```go
// simplifyPolynomial 化简多项式表达式
func (c *SymbolicCalcCalculatorFixed) simplifyPolynomial(expr string) string {
    // 1. 解析为AST
    ast := c.parseToAST(expr)
    
    // 2. 收集同类项
    terms := c.collectTerms(ast)
    
    // 3. 合并同类项
    merged := c.mergeTerms(terms)
    
    // 4. 构建化简后的表达式
    return c.buildPolynomialString(merged)
}

// Term 表示多项式中的一项
type Term struct {
    Coefficient float64 // 系数
    Power       float64 // 变量的幂次
    Variable    string  // 变量名
}
```

#### 3.2.5 AST求值

```go
// evaluateAST 对AST进行数值求值
func (c *SymbolicCalcCalculatorFixed) evaluateAST(node *ASTNode, variables map[string]float64) float64 {
    switch node.Type {
    case "number":
        return node.Value
    case "variable":
        if val, ok := variables[node.Name]; ok {
            return val
        }
        return 0
    case "add":
        return c.evaluateAST(node.Children[0], variables) + 
               c.evaluateAST(node.Children[1], variables)
    case "mul":
        return c.evaluateAST(node.Children[0], variables) * 
               c.evaluateAST(node.Children[1], variables)
    case "pow":
        base := c.evaluateAST(node.Children[0], variables)
        exp := c.evaluateAST(node.Children[1], variables)
        return math.Pow(base, exp)
    // ...
    }
}
```

---

## 4. 符号微分算法要点

### 4.1 表达式解析

1. **分词(Tokenization)**：将表达式字符串分割为token序列
2. **语法分析(Parsing)**：根据运算符优先级构建AST
3. **优先级处理**：括号 > 幂运算 > 乘除 > 加减

### 4.2 求导规则递归应用

| 规则 | 数学表达式 | 实现方式 |
|------|-----------|---------|
| 常数规则 | `d/dx(c) = 0` | 返回数值节点0 |
| 变量规则 | `d/dx(x) = 1` | 返回数值节点1 |
| 加法规则 | `d/dx(u+v) = u' + v'` | 递归求导后相加 |
| 乘法规则 | `d/dx(u*v) = u'*v + u*v'` | 乘积法则实现 |
| 幂函数规则 | `d/dx(x^n) = n*x^(n-1)` | 构建新的幂节点 |
| 链式法则 | `d/dx(f(g(x))) = f'(g(x))*g'(x)` | 复合函数处理 |

### 4.3 化简策略

1. **常数折叠**：`2 + 3` → `5`
2. **零元消除**：`x + 0` → `x`, `x * 0` → `0`
3. **单位元消除**：`x * 1` → `x`, `x^1` → `x`
4. **同类项合并**：`3*x^2 + 2*x^2` → `5*x^2`
5. **降幂排序**：按幂次从高到低排列

### 4.4 数值求值正确性保证

1. **变量替换**：在AST层面进行变量值替换，避免字符串替换带来的精度问题
2. **递归求值**：自底向上计算每个节点的值
3. **除零保护**：检查除数是否为零
4. **定义域检查**：对数、开方等运算检查参数有效性

---

## 5. 修复验证

### 5.1 测试用例

```bash
# 测试修复版接口
curl -X 'POST' \
  'http://localhost:8080/api/calculate-fixed' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "operation": "differentiate",
    "expression": "x^3 + 2*x^2 + x",
    "variable": "x"
  }'
```

预期响应：
```json
{
  "success": true,
  "result": {
    "original_expression": "x^3 + 2*x^2 + x",
    "result_expression": "3*x^2 + 4*x + 1",
    "derivative": "3*x^2 + 4*x + 1",
    "simplified": "3*x^2 + 4*x + 1",
    "numeric_value": 0,
    "operation_type": "differentiate"
  },
  "fixed": true
}
```

### 5.2 对比接口测试

```bash
# 对比修复前后的结果
curl -X 'POST' \
  'http://localhost:8080/api/solver/compare' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "operation": "differentiate",
    "expression": "x^3 + 2*x^2 + x",
    "variable": "x"
  }'
```

---

## 6. 接口说明

### 6.1 修复版接口 `/api/calculate-fixed`

- **请求格式**：与原 `/api/calculate` 完全一致
- **响应格式**：与原接口一致，新增 `fixed: true` 标记
- **修复内容**：
  - 正确的符号求导计算
  - 完整的多项式化简
  - 准确的数值求值
  - 有效的解析树生成

### 6.2 对比接口 `/api/solver/compare`

- **用途**：对比修复前后的计算结果
- **响应字段**：
  - `original`: 原版计算结果
  - `fixed`: 修复版计算结果
  - `differences`: 差异详情列表
  - `is_fixed`: 是否已修复
  - `expected`: 预期正确结果

---

## 7. 代码审查清单

- [x] 实现完整的表达式解析器（AST构建）
- [x] 实现递归求导算法
- [x] 实现多项式化简（合并同类项）
- [x] 实现AST数值求值
- [x] 保持与原接口的请求/响应格式兼容
- [x] 添加对比接口用于验证修复效果
- [x] 更新 Swagger 文档

---

## 8. 后续优化建议

1. **支持更多函数**：sin, cos, exp, log 等初等函数
2. **多元函数求导**：支持偏导数计算
3. **高阶导数**：支持二阶及以上导数
4. **积分计算**：实现符号积分功能
5. **方程求解**：结合求导实现牛顿迭代法

---

**修复日期**: 2026-03-26  
**修复版本**: v1.1  
**作者**: ScientificCalc Team
