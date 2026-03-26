#【GLM5.0】模型的符号微分计算错误修复报告

## 问题描述

### 原始请求
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

### 错误响应
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

### 预期正确结果
```json
{
  "success": true,
  "result": {
    "original_expression": "x^3 + 2*x^2 + x",
    "result_expression": "3*x^2 + 4*x + 1",
    "parsed_tree": { ... },
    "derivative": "3*x^2 + 4*x + 1",
    "simplified": "3*x^2 + 4*x + 1",
    "numeric_value": 0,
    "variables": {"x": 0},
    "operation_type": "differentiate"
  }
}
```

## 异常字段分析

| 字段 | 当前值 | 预期值 | 问题描述 |
|------|--------|--------|----------|
| `derivative` | `"d/dx(x^3 + 2*x^2 + x)"` | `"3*x^2 + 4*x + 1"` | 未实际计算，仅返回占位符字符串 |
| `simplified` | `""` | `"3*x^2 + 4*x + 1"` | 为空，化简功能未实现 |
| `numeric_value` | `0` | 需代入变量值计算 | 恒为0，数值求值未实现 |
| `parsed_tree` | `null` | 表达式树结构 | 为null，解析失败 |
| `variables` | `null` | `{"x": 0}` | 为null，变量提取失败 |

## 根本原因分析

### 1. 求导算法未递归实现

**问题代码位置**: `symbolic_calc.go` 第 331-357 行

```go
func (c *SymbolicCalcCalculator) calculateDerivative(expr, variable string) string {
    // 只能处理单个项，无法处理多项式
    if match := regexp.MustCompile(`^` + variable + `\^(\d+)$`).FindStringSubmatch(expr); match != nil {
        n, _ := strconv.Atoi(match[1])
        if n == 2 {
            return "2*" + variable
        }
        return fmt.Sprintf("%d*%s^%d", n, variable, n-1)
    }
    // 默认返回原始表达式（表示无法求导）
    return "d/d" + variable + "(" + expr + ")"
}
```

**问题**:
- 缺少加法法则 `(f+g)' = f' + g'`
- 缺少乘法法则 `(f*g)' = f'*g + f*g'`
- 无法递归处理复合表达式

### 2. 缺少表达式树解析

**问题**: 没有将表达式解析为AST（抽象语法树），无法递归应用求导规则

**原版 buildParseTree 函数**:
```go
func (c *SymbolicCalcCalculator) buildParseTree(expr string) map[string]interface{} {
    tree := map[string]interface{}{
        "type":     "expression",
        "value":    expr,  // 仅存储原始字符串
        "tokens":   c.tokenize(expr),
        "operator": c.findMainOperator(expr),
    }
    return tree
}
```

### 3. 化简功能过于简化

**问题**: 只能做简单字符串替换，无法进行代数化简

### 4. 数值求值功能不完整

**问题**: 只能处理简单二元运算，无法处理复杂表达式

### 5. 为何返回 `success: true`？

代码没有检测到错误，只是返回了占位符字符串作为"结果"，没有验证计算结果的有效性。

## 修复方案

### 符号微分算法要点

#### 1. 表达式解析

将数学表达式解析为抽象语法树（AST），支持以下节点类型：

```go
type ExprNode interface {
    String() string
    Equals(other ExprNode) bool
}

type NumberNode struct { Value float64 }      // 数字常量
type VariableNode struct { Name string }       // 变量
type BinaryOpNode struct {                     // 二元运算
    Op    string
    Left  ExprNode
    Right ExprNode
}
type UnaryOpNode struct {                      // 一元运算
    Op      string
    Operand ExprNode
}
type FunctionNode struct {                     // 函数调用
    Name string
    Arg  ExprNode
}
```

**解析器实现**:
- 词法分析：将表达式字符串分解为token序列
- 语法分析：根据运算符优先级构建AST
  - 优先级：`+ -` < `* /` < `^` < 一元运算 < 函数调用 < 括号

#### 2. 求导规则递归应用

**基本规则**:
```
常数规则: d/dx(c) = 0
变量规则: d/dx(x) = 1, d/dx(y) = 0 (y ≠ x)
```

**二元运算规则**:
```
加法法则: d/dx(f + g) = f' + g'
减法法则: d/dx(f - g) = f' - g'
乘法法则: d/dx(f * g) = f' * g + f * g'
除法法则: d/dx(f / g) = (f' * g - f * g') / g²
幂法则: d/dx(f^n) = n * f^(n-1) * f'  (n为常数)
```

**函数求导规则**:
```
sin(x)' = cos(x)
cos(x)' = -sin(x)
tan(x)' = 1/cos²(x)
exp(x)' = exp(x)
ln(x)' = 1/x
sqrt(x)' = 1/(2*sqrt(x))
```

**实现示例**:
```go
func Differentiate(node ExprNode, variable string) ExprNode {
    switch n := node.(type) {
    case *NumberNode:
        return &NumberNode{Value: 0}
    case *VariableNode:
        if n.Name == variable {
            return &NumberNode{Value: 1}
        }
        return &NumberNode{Value: 0}
    case *BinaryOpNode:
        switch n.Op {
        case "+":
            return &BinaryOpNode{
                Op:    "+",
                Left:  Differentiate(n.Left, variable),
                Right: Differentiate(n.Right, variable),
            }
        case "*":
            // 乘法法则: (f*g)' = f'*g + f*g'
            return &BinaryOpNode{
                Op: "+",
                Left: &BinaryOpNode{
                    Op:    "*",
                    Left:  Differentiate(n.Left, variable),
                    Right: n.Right,
                },
                Right: &BinaryOpNode{
                    Op:    "*",
                    Left:  n.Left,
                    Right: Differentiate(n.Right, variable),
                },
            }
        // ... 其他运算符
        }
    // ... 其他节点类型
    }
}
```

#### 3. 化简策略

**常数折叠**:
```
2 + 3 → 5
2 * 3 → 6
x + 0 → x
x * 1 → x
x * 0 → 0
x^0 → 1
x^1 → x
```

**代数化简**:
```
x + (-y) → x - y
(-x) + y → y - x
```

**递归化简实现**:
```go
func Simplify(node ExprNode) ExprNode {
    switch n := node.(type) {
    case *BinaryOpNode:
        // 先化简子节点
        n.Left = Simplify(n.Left)
        n.Right = Simplify(n.Right)
        
        // 常数折叠
        if leftNum, ok := n.Left.(*NumberNode); ok {
            if rightNum, ok := n.Right.(*NumberNode); ok {
                switch n.Op {
                case "+":
                    return &NumberNode{Value: leftNum.Value + rightNum.Value}
                case "*":
                    return &NumberNode{Value: leftNum.Value * rightNum.Value}
                // ...
                }
            }
        }
        
        // 代数化简
        // x + 0 → x
        // x * 1 → x
        // ...
    }
}
```

#### 4. 数值求值正确性

**变量替换**:
```go
func Evaluate(node ExprNode, variables map[string]float64) (float64, error) {
    switch n := node.(type) {
    case *NumberNode:
        return n.Value, nil
    case *VariableNode:
        if val, ok := variables[n.Name]; ok {
            return val, nil
        }
        return 0, fmt.Errorf("undefined variable: %s", n.Name)
    case *BinaryOpNode:
        left, _ := Evaluate(n.Left, variables)
        right, _ := Evaluate(n.Right, variables)
        switch n.Op {
        case "+": return left + right, nil
        case "*": return left * right, nil
        case "^": return math.Pow(left, right), nil
        // ...
        }
    }
}
```

## 新增接口

### 1. 修复版接口 `/api/calculate-fixed`

**请求格式**与原 `/api/calculate` 完全一致：

```bash
curl -X 'POST' \
  'http://localhost:8080/api/calculate-fixed' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "differentiate",
      "expression": "x^3 + 2*x^2 + x",
      "variable": "x"
    }
  }'
```

**响应示例**:
```json
{
  "success": true,
  "result": {
    "original_expression": "x^3 + 2*x^2 + x",
    "result_expression": "3*x^2+4*x+1",
    "parsed_tree": {
      "type": "binary_op",
      "op": "+",
      "left": { ... },
      "right": { ... }
    },
    "derivative": "3*x^2+4*x+1",
    "simplified": "3*x^2+4*x+1",
    "numeric_value": 0,
    "variables": {"x": 0},
    "operation_type": "differentiate"
  }
}
```

### 2. 对比接口 `/api/solver/compare`

**请求格式**:
```bash
curl -X 'POST' \
  'http://localhost:8080/api/solver/compare' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "calculation": "symbolic_calc",
    "params": {
      "operation": "differentiate",
      "expression": "x^3 + 2*x^2 + x",
      "variable": "x"
    }
  }'
```

**响应示例**:
```json
{
  "original_result": { ... },
  "fixed_result": { ... },
  "differences": [
    {
      "field": "derivative",
      "original": "d/dx(x^3 + 2*x^2 + x)",
      "fixed": "3*x^2+4*x+1",
      "issue": "导数计算结果不同"
    },
    {
      "field": "simplified",
      "original": "",
      "fixed": "3*x^2+4*x+1",
      "issue": "化简结果不同"
    }
  ],
  "is_fixed": true,
  "comparison_time": "2026-03-26T19:00:00+08:00"
}
```

## 测试用例

### 测试1: 多项式求导
```
输入: x^3 + 2*x^2 + x
预期导数: 3*x^2 + 4*x + 1
```

### 测试2: 三角函数求导
```
输入: sin(x) + cos(x)
预期导数: cos(x) - sin(x)
```

### 测试3: 复合函数求导
```
输入: x^2 * sin(x)
预期导数: 2*x*sin(x) + x^2*cos(x)
```

### 测试4: 数值求值
```
输入: x^2 + 2*x + 1, x = 3
预期值: 16
```

## 文件变更

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/calculator/symbolic_engine.go` | 新增 | 符号计算引擎核心实现 |
| `internal/calculator/symbolic_calc_fixed.go` | 新增 | 修复版符号计算器 |
| `internal/api/handler.go` | 修改 | 添加新接口路由 |

## 后续建议

1. **单元测试**: 为符号计算引擎添加完整的单元测试覆盖
2. **性能优化**: 对于复杂表达式，考虑使用缓存优化
3. **错误处理**: 增强错误信息，提供更详细的错误位置
4. **扩展功能**: 支持更多数学函数（如 tan, cot, sec, csc 等）
5. **积分运算**: 基于现有框架扩展积分功能

## 版本信息

- 修复日期: 2026-03-26
- 修复版本: v1.1.0
- 影响范围: symbolic_calc 计算类型
