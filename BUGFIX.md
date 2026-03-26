# 【GLM5.0】模型的北斗七星位置推算 Bug修复报告

## 1. 问题概述

原始API响应（请求日期：2026年3月23日，star_name=big_dipper）存在多处科学错误：

```json
{
  "lunar_date": "丙午年3月23日",
  "day_ganzhi": "丁卯",
  "constellation": "翼",
  "star_position": "朱雀在西方",
  "auspicious": true,
  "day_score": 3,
  "constellation_index": 26,
  "auspicious_level": 9.5,
  "julian_day": 2461122.5
}
```

## 2. 错误分析

### 2.1 北斗七星处理错误

| 错误项 | 原始值 | 正确处理 |
|--------|--------|----------|
| star_name参数 | 未处理 | 应识别并返回北斗七星专属信息 |
| 返回结果 | 返回"翼"宿 | 北斗七星不是二十八宿，应返回北斗七星天文数据 |

**说明**：北斗七星（Big Dipper）是大熊座中的星群，不属于二十八宿体系。当请求`star_name=big_dipper`时，应返回北斗七星的专属天文信息。

### 2.2 星宿方位错误

| 错误项 | 原始值 | 正确值 |
|--------|--------|--------|
| 翼宿方位 | "朱雀在西方" | "朱雀在南方" |

**说明**：二十八宿按方位分为四组：
- 东方青龙七宿：角、亢、氐、房、心、尾、箕
- 北方玄武七宿：斗、牛、女、虚、危、室、壁
- 西方白虎七宿：奎、娄、胃、昴、毕、觜、参
- 南方朱雀七宿：井、鬼、柳、星、张、翼、轸

翼宿属于南方朱雀七宿，方位应为南方，而非西方。

### 2.3 农历日期错误

| 错误项 | 原始值 | 正确值 |
|--------|--------|--------|
| 农历日期 | "丙午年3月23日" | 需精确计算 |

**说明**：原实现直接将公历月日作为农历月日，这是错误的。公历2026年3月23日对应的农历需要通过天文算法精确计算。

### 2.4 日干支计算错误

| 错误项 | 原始值 | 正确值 |
|--------|--------|--------|
| 日干支 | "丁卯" | 需精确计算 |

**说明**：原实现使用错误的基准日计算日干支。

### 2.5 评分体系不自洽

| 错误项 | 原始值 | 问题说明 |
|--------|--------|----------|
| day_score | 3 | 基于日期哈希，无实际意义 |
| auspicious_level | 9.5 | 与day_score=3矛盾，auspicious=true但评分逻辑混乱 |

## 3. 修复方案

### 3.1 儒略日计算公式

```
儒略日(JD) = INT(365.25 × (Y + 4716)) + INT(30.6001 × (M + 1)) + D + B - 1524.5

其中:
- Y = 年份（若月份≤2，则Y = 年份 - 1）
- M = 月份（若月份≤2，则M = 月份 + 12）
- D = 日期（含小数部分）
- A = INT(Y / 100)
- B = 2 - A + INT(A / 4)
```

**验证**：2026年3月23日
- Y = 2026, M = 3, D = 23
- A = 20, B = 2 - 20 + 5 = -13
- JD ≈ 2461122.5

### 3.2 日干支推算规则

```
基准日：1900年1月31日为甲子日（干支周期第1天）

日干支计算：
1. 计算目标日期与基准日的天数差
2. 天干索引 = (天数差 + 4) % 10
3. 地支索引 = (天数差 + 4) % 12

天干：甲(0)、乙(1)、丙(2)、丁(3)、戊(4)、己(5)、庚(6)、辛(7)、壬(8)、癸(9)
地支：子(0)、丑(1)、寅(2)、卯(3)、辰(4)、巳(5)、午(6)、未(7)、申(8)、酉(9)、戌(10)、亥(11)
```

### 3.3 二十八宿值日算法

```
二十八宿值日索引 = (年日序 + 偏移量) % 28

其中：
- 年日序：该日在年中的序号（1月1日=1）
- 偏移量：根据基准年调整，使算法与实际值日对应
```

### 3.4 二十八宿方位映射

```
二十八宿方位映射表：

东方青龙七宿（木）：
  角、亢、氐、房、心、尾、箕

北方玄武七宿（水）：
  斗、牛、女、虚、危、室、壁

西方白虎七宿（金）：
  奎、娄、胃、昴、毕、觜、参

南方朱雀七宿（火）：
  井、鬼、柳、星、张、翼、轸
```

### 3.5 北斗七星天文数据

北斗七星各星的天文坐标（J2000.0历元）：

| 星名 | 传统名称 | 赤经(度) | 赤纬(度) | 星等 | 所属星座 |
|------|----------|----------|----------|------|----------|
| 天枢 | Dubhe | 165.93 | 61.75 | 1.79 | 大熊座 |
| 天璇 | Merak | 165.46 | 56.38 | 2.37 | 大熊座 |
| 天玑 | Phecda | 178.46 | 57.03 | 2.44 | 大熊座 |
| 天权 | Megrez | 183.86 | 57.04 | 3.31 | 大熊座 |
| 玉衡 | Alioth | 193.51 | 55.96 | 1.77 | 大熊座 |
| 开阳 | Mizar | 200.98 | 54.93 | 2.27 | 大熊座 |
| 摇光 | Alkaid | 206.89 | 49.31 | 1.86 | 大熊座 |

### 3.6 农历转换算法

农历转换需要基于天文算法计算新月时刻。简化实现使用查表法：

```
1. 存储各年农历新年对应的公历日期偏移
2. 计算目标日期与农历新年的天数差
3. 根据农历月份天数表计算农历月日
4. 处理闰月情况
```

### 3.7 评分体系修正

**日分值计算（0-100分）：**
```
基础分 = 50
若日干支吉利 → +25分
若星宿吉利 → +25分
最终分数 = min(100, max(0, 计算分数))
```

**吉凶程度计算（0-10分）：**
```
基础分 = 5.0
若auspicious=true → +2.0
若"日干支吉利" → +1.0
若"星宿吉利" → +1.0
若"日干支平常" → -0.5
若"星宿平常" → -0.5
最终分数 = min(10, max(0, 计算分数))
```

## 4. 修复后API响应示例

### 请求
```bash
POST /api/calculate-fixed
Content-Type: application/json

{
  "calculation": "star",
  "params": {
    "year": 2026,
    "month": 3,
    "day": 23,
    "star_name": "big_dipper"
  }
}
```

### 响应
```json
{
  "success": true,
  "result": {
    "lunar_date": "丙午年二月初五",
    "day_ganzhi": "庚辰",
    "constellation": "奎",
    "constellation_direction": "白虎西方",
    "star_position": "北斗七星在北方天空",
    "auspicious": true,
    "auspicious_info": ["日干支吉利", "星宿吉利"],
    "day_score": 100,
    "constellation_index": 14,
    "auspicious_level": 9.0,
    "julian_day": 2461122.5,
    "time_coordinate": 337.25,
    "star_name": "big_dipper",
    "big_dipper_info": {
      "stars": [
        {"name": "天枢", "right_ascension": 165.93, "declination": 61.75, "magnitude": 1.79, "constellation": "大熊座"},
        {"name": "天璇", "right_ascension": 165.46, "declination": 56.38, "magnitude": 2.37, "constellation": "大熊座"},
        {"name": "天玑", "right_ascension": 178.46, "declination": 57.03, "magnitude": 2.44, "constellation": "大熊座"},
        {"name": "天权", "right_ascension": 183.86, "declination": 57.04, "magnitude": 3.31, "constellation": "大熊座"},
        {"name": "玉衡", "right_ascension": 193.51, "declination": 55.96, "magnitude": 1.77, "constellation": "大熊座"},
        {"name": "开阳", "right_ascension": 200.98, "declination": 54.93, "magnitude": 2.27, "constellation": "大熊座"},
        {"name": "摇光", "right_ascension": 206.89, "declination": 49.31, "magnitude": 1.86, "constellation": "大熊座"}
      ],
      "direction": "北方天空",
      "right_ascension": 184.0,
      "declination": 55.0,
      "visibility": "前半夜可见",
      "culmination_time": "约22时中天"
    },
    "star_info": {
      "type": "big_dipper",
      "name": "北斗七星",
      "direction": "北方天空",
      "visibility": "前半夜可见",
      "description": "北斗七星属大熊座，是北半球最著名的星群之一"
    }
  },
  "warnings": null,
  "timestamp": "2026-03-26T17:30:00+08:00",
  "calculation": "star_fixed",
  "session_id": "xxxxxxxxxxxxxxxx"
}
```

## 5. 对比接口使用

### 请求
```bash
POST /api/solver/compare
Content-Type: application/json

{
  "calculation": "star",
  "params": {
    "year": 2026,
    "month": 3,
    "day": 23,
    "star_name": "big_dipper"
  }
}
```

### 响应
```json
{
  "success": true,
  "original_result": { ... 原始结果 ... },
  "fixed_result": { ... 修复后结果 ... },
  "differences": [
    {
      "field": "lunar_date",
      "original_value": "丙午年3月23日",
      "fixed_value": "丙午年二月初五",
      "description": "农历日期修正：原实现直接使用公历月日作为农历，修复版实现正确的公历转农历算法"
    },
    {
      "field": "day_ganzhi",
      "original_value": "丁卯",
      "fixed_value": "庚辰",
      "description": "日干支修正：原实现使用错误的基准日，修复版使用正确的1900年1月31日（甲子日）作为基准"
    },
    {
      "field": "star_position",
      "original_value": "朱雀在西方",
      "fixed_value": "北斗七星在北方天空",
      "description": "星曜位置修正：原实现将南方朱雀标注为西方，修复版正确标注方位"
    },
    {
      "field": "big_dipper_info",
      "original_value": null,
      "fixed_value": { ... },
      "description": "新增北斗七星专属信息：原实现未处理star_name=big_dipper参数，修复版返回北斗七星详细天文数据"
    }
  ],
  "summary": "共发现 4 处差异",
  "timestamp": "2026-03-26T17:30:00+08:00"
}
```

## 6. 文件变更清单

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/calculator/star_fixed.go` | 新增 | 修复版星曜计算器 |
| `internal/api/handler.go` | 修改 | 添加修复版接口和对比接口 |
| `BUGFIX.md` | 新增 | 本修复报告文档 |

## 7. 测试验证

### 7.1 启动服务
```bash
go run cmd/server/main.go
```

### 7.2 测试修复版接口
```bash
curl -X POST 'http://localhost:8080/api/calculate-fixed' \
  -H 'Content-Type: application/json' \
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

### 7.3 测试对比接口
```bash
curl -X POST 'http://localhost:8080/api/solver/compare' \
  -H 'Content-Type: application/json' \
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

## 8. 参考资料

1. 《天文算法》(Astronomical Algorithms) - Jean Meeus
2. 《中国天文年历》- 中国科学院紫金山天文台
3. 二十八宿值日推算方法 - 中国传统历法
4. 北斗七星天文数据 - SIMBAD Astronomical Database
