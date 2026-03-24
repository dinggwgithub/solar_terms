用户想评估三个AI模型在科学计算领域的质量对比，需要从`solar_terms.py`中找出适合的任务。我需要先仔细分析这个文件的内容，看看有哪些科学计算相关的部分可以用于模型对比。

这个文件看起来是一个八字排盘系统，里面包含了很多科学计算元素：节气计算、天文计算（太阳黄经、儒略日转换）、干支计算、数值算法（二分法求解方程）、数学函数等。这些都是典型的科学计算任务，适合用来测试AI模型的数值计算精度、算法实现能力和复杂数学问题的处理能力。

我需要先全面了解文件内容，找出所有潜在的科学计算任务。让我仔细阅读文件，特别是其中的计算模块。文件包含几个主要模块：加密验证、常量管理、干支计算、节气计算、大运计算、纳音计算、流年计算、宫位计算等。其中节气计算模块涉及天文计算，干支计算涉及模运算和日期处理，这些都是科学计算任务。

让我先快速浏览整个文件，然后提取出具体的计算任务。我需要找到那些有明确输入输出、可以独立测试的计算函数。例如：
1. `SolarTermCalculator.calculate_solar_term` - 计算节气精确时间
2. `SolarTermCalculator._sun_longitude` - 计算太阳视黄经
3. `SolarTermCalculator._to_julian_date` - 转换为儒略日
4. `GanZhiCalculator.get_gan_day` - 计算日天干
5. `GanZhiCalculator.get_zhi_day` - 计算日地支
6. `SolarTermCalculator.s_term` - 计算节气日期
7. `DaYunCalculator.compute_starting_age` - 计算起运岁数
8. `LiuNianCalculator.compute_liu_nian` - 计算流年信息
9. `SolarTermCalculator._find_exact_time` - 二分法查找精确时间

这些函数都有数学计算、天文算法或数值方法，适合用于AI模型对比。我需要为每个任务列出详细信息：任务描述、输入输出、计算复杂度、评估指标等。

让我先读取文件，然后提取函数定义。我需要并行读取多个相关文件吗？目前只需要分析`solar_terms.py`，但可能还需要查看其他文件来理解上下文。不过用户明确要求分析这个文件，所以先集中在这里。

我需要创建一个任务清单，列出每个潜在的科学计算任务，包括：
- 任务名称
- 函数/方法
- 输入参数
- 输出结果
- 计算类型（数值计算、天文算法、日期计算等）
- 评估指标（精度、速度、正确性等）
- 测试用例建议

让我开始分析文件内容。首先，我需要仔细阅读文件，提取所有函数定义和计算逻辑。我将使用搜索工具来快速定位函数定义。同时，我需要理解每个函数的数学原理，以便设计合适的测试。

我将使用`search_content`工具来查找函数定义，并读取关键函数的具体实现。现在开始并行搜索。
搜索结果太多，我需要更精确地分析`solar_terms.py`文件。让我直接读取整个文件，然后提取出其中的函数定义。我已经有了文件内容，但需要更系统地分析。让我仔细查看文件内容，重点关注计算相关的函数。

从文件内容来看，有以下关键计算模块：

1. **节气计算模块** (`SolarTermCalculator`):
   - `calculate_solar_term(year, term_index)` - 计算节气精确时间
   - `_sun_longitude(date)` - 计算太阳视黄经
   - `_to_julian_date(date)` - 转换为儒略日
   - `_find_exact_time(start_date, target_longitude)` - 二分法查找精确时间
   - `s_term(year, n)` - 计算节气日期

2. **干支计算模块** (`GanZhiCalculator`):
   - `get_gan_year(year, solar_date)` - 计算年天干
   - `get_zhi_year(year, solar_date)` - 计算年地支
   - `get_gan_month(year, month)` - 计算月天干
   - `get_zhi_month(year, month, solar_date)` - 计算月地支
   - `get_gan_day(year, month, day)` - 计算日天干
   - `get_zhi_day(year, month, day)` - 计算日地支
   - `get_gan_time(hour)` - 计算时天干
   - `get_zhi_time(hour)` - 计算时地支

3. **大运计算模块** (`DaYunCalculator`):
   - `compute_starting_age(days_from_term)` - 计算起运岁数
   - `get_wang_chai(tian_gan, di_zhi)` - 获取旺衰状态

4. **其他计算模块**:
   - `NaYinCalculator.get_na_yin(gan_zhi)` - 计算纳音
   - `LiuNianCalculator.compute_liu_nian(start_year, gan_zhi)` - 计算流年
   - `GongCalculator` - 宫位计算

我需要为每个任务提供详细描述。让我创建一个任务清单，包括以下信息：
- 任务名称
- 函数/方法
- 输入参数
- 输出结果
- 计算类型
- 评估指标
- 测试用例建议

首先，我需要提取文件中的所有函数定义。让我搜索`solar_terms.py`文件中的函数定义。用更精确的搜索。
搜索结果仍然太多。我需要专注于`solar_terms.py`文件。让我直接分析我已经有的文件内容。我可以从文件内容中手动提取函数定义。

从文件内容看，有以下函数（按类分组）：

**TokenValidator类**：
- `aes_decrypt(cipher_text)` - AES解密
- `validate_api_token(encrypted_token)` - API令牌验证

**GanZhiCalculator类**：
- `get_gan_year(year, solar_date)`
- `get_zhi_year(year, solar_date)`
- `get_gan_month(year, month)`
- `get_zhi_month(year, month, solar_date)`
- `get_gan_day(year, month, day)`
- `get_zhi_day(year, month, day)`
- `get_gan_time(hour)`
- `get_zhi_time(hour)`
- `calculator_get_zhi(year, month, day, hour, solar_date)`
- `get_ganzhi_str()`
- `get_ganzhi()`

**SolarTermCalculator类**：
- `get_term_info(year, month, day, solar_date)`
- `set_lunar_term(year, month, day, sex)`
- `calculate_solar_term(year, term_index)`
- `is_forward(elements1, sex)`
- `_estimate_start_date(year, term_index)`
- `_find_exact_time(start_date, target_longitude)`
- `_sun_longitude(date)`
- `_to_julian_date(date)`
- `s_term(year, n)`
- `special_year_offset(year, n)`
- `testAllTermTime(args, year)`

**DaYunCalculator类**：
- `compute_starting_age(days_from_term)`
- `find_index(gan_zhi)`
- `get_next_nine_gan_zhi(gan_zhi, direction)`
- `get_da_yun(gan_zhi, direction)`
- `get_age(gan_zhi, direction, starting_age, solar_year)`
- `get_wang_chai(tian_gan, di_zhi)`
- `compute_da_yun_wang_chai(da_yun)`
- `compute_rigan_ws()`
- `compute_sizhu_ws()`

**NaYinCalculator类**：
- `get_na_yin(gan_zhi)`
- `get_nayin()`

**LiuNianCalculator类**：
- `compute_liu_nian(start_year, gan_zhi)`

**GongCalculator类**：
- `calculate_tai_yuan(month_gan, month_zhi)`
- `calculate_tai_xi(ri_gan, ri_zhi)`

**CalendarHandler类**：
- `_get_ten_god(ri_zhu_gan, hidden_gan)`
- `_format_bazi_pan(ganzhi, gender)`
- `_format_ages(ages)`
- `_format_da_yun_years(da_yun_str, ages, year, da_yun_wang_chai)`
- `_count_wuxing(gan_zhi)`
- `tai_yuan_calculator(calendar_info)`
- `bazi_ming_pan(args)`

现在我需要从中筛选出适合AI模型对比的科学计算任务。科学计算任务应该涉及数值计算、算法实现、数学建模等。以下是一些候选：

1. **天文计算**：节气计算、太阳黄经计算、儒略日转换
2. **数值算法**：二分法求解方程
3. **日期计算**：干支计算、模运算
4. **统计计算**：五行数量统计

让我为每个任务提供详细描述。我需要创建一个清单，列出8-10个任务。每个任务包括：
- 任务名称
- 函数/方法
- 输入参数
- 输出结果
- 计算类型
- 评估指标
- 测试用例建议

现在开始创建清单。我将基于文件内容提取信息。让我先仔细查看几个关键函数的实现细节。

我需要读取`solar_terms.py`中的特定函数。让我搜索`calculate_solar_term`函数的具体实现。
现在搜索`_sun_longitude`函数：
搜索`_to_julian_date`：
搜索`_find_exact_time`：
搜索`get_gan_day`：
搜索`s_term`：
搜索`compute_starting_age`：
搜索`compute_liu_nian`：
现在我有足够的信息来创建任务清单。我需要列出8-10个任务，涵盖不同类型的科学计算。让我组织一下：

1. **节气精确时间计算** - 天文计算，涉及太阳黄经和二分法
2. **太阳视黄经计算** - 天文计算，涉及复杂的数学公式
3. **儒略日转换** - 天文历法计算
4. **日天干计算** - 模运算和日期算法
5. **节气日期计算** - 数值计算和查表法
6. **起运岁数计算** - 简单的算术运算
7. **流年信息计算** - 循环和条件逻辑
8. **五行数量统计** - 字符串处理和统计

让我为每个任务创建详细的条目。现在开始撰写回答。我将提供一个表格形式的清单。
基于对 `solar_terms.py` 文件的深入分析，该文件包含丰富的科学计算任务，非常适合用于AI模型在**数值精度、算法实现、天文计算和复杂数学建模**方面的质量对比。以下是精心筛选的8个科学计算任务清单：

| 序号 | 任务名称 | 输入参数 | 输出结果 | 计算类型 | 评估指标 | 函数/方法 | 测试用例建议 |
|------|----------|----------|----------|----------|----------|-----------|--------------|
| 1 | **节气精确时间计算** | `year`: int (年份)<br>`term_index`: int (0-23，节气索引) | 字符串，格式: `"YYYY-MM-DD HH:MM:SS"` (北京时间) | **天文算法**<br>1. 太阳黄经目标值计算<br>2. 二分法迭代求解<br>3. 时区转换 | 1. **时间精度** (秒级)<br>2. **收敛速度** (迭代次数)<br>3. **边界处理** (闰年、特殊年份) | `SolarTermCalculator.calculate_solar_term(year, term_index)` | 1. 2024年立春 (term_index=2)<br>2. 2000年冬至 (term_index=23)<br>3. 2026年大寒 (term_index=1) |
| 2 | **太阳视黄经计算** | `date`: datetime (UTC时间) | float (度，范围 [0, 360)) | **天体力学计算**<br>1. 儒略日转换<br>2. 平黄经、平近点角计算<br>3. 中心差、章动修正 | 1. **数值精度** (与NASA JPL DE星历对比)<br>2. **计算效率** (浮点运算量)<br>3. **周期性处理** (模360度) | `SolarTermCalculator._sun_longitude(date)` | 1. 2025-01-01 00:00:00 UTC<br>2. 2000-01-01 12:00:00 UTC<br>3. 1990-06-21 06:00:00 UTC |
| 3 | **儒略日转换** | `date`: datetime (UTC时间) | float (儒略日) | **历法转换算法**<br>1. 格里高利历转儒略日<br>2. 年月日时分秒分解<br>3. 浮点累加计算 | 1. **转换精度** (与标准儒略日表对比)<br>2. **时区处理** (UTC确保)<br>3. **整数溢出风险** | `SolarTermCalculator._to_julian_date(date)` | 1. 2000-01-01 12:00:00 UTC<br>2. 1900-03-01 00:00:00 UTC<br>3. 2100-12-31 23:59:59 UTC |
| 4 | **日天干计算** | `year`: int, `month`: int, `day`: int | 字符串 (十天干之一) | **模运算与日期算法**<br>1. 年份范围判断<br>2. 闰年修正<br>3. 复杂公式: `((year%100)*5 + temp + base + day + month_double + cor) % 60 % 10` | 1. **算法正确性** (与万年历对比)<br>2. **边界值测试** (世纪交界)<br>3. **闰年特殊处理** | `GanZhiCalculator.get_gan_day(year, month, day)` | 1. 2024-02-29 (闰年)<br>2. 2000-01-01 (世纪交界)<br>3. 1997-10-22 (用户案例) |
| 5 | **节气日期计算** | `year`: int, `n`: int (节气索引) | int (当月日期) | **数值计算 + 查表法**<br>1. 世纪数组查表<br>2. 线性公式计算<br>3. 特殊年份偏移调整 | 1. **查表准确性**<br>2. **特殊年份处理** (如1982年立春偏移)<br>3. **年份范围限制** (1900-2100) | `SolarTermCalculator.s_term(year, n)` | 1. 2025年惊蛰 (n=3)<br>2. 1982年立春 (n=2，有偏移)<br>3. 2089年寒露 (n=19，有偏移) |
| 6 | **起运岁数计算** | `days_from_term`: int (距节气天数) | 字符串 (如 "5岁8个月") | **简单算术运算**<br>1. 整数除法与取余<br>2. 单位转换 (余数→月) | 1. **运算正确性**<br>2. **格式化输出**<br>3. **边界值** (整除情况) | `DaYunCalculator.compute_starting_age(days_from_term)` | 1. 17天 → "5岁8个月"<br>2. 15天 → "5岁"<br>3. 0天 → "0岁" |
| 7 | **流年信息批量计算** | `start_year`: int, `gan_zhi`: str (年柱) | List[Dict] (每年干支、纳音、旺衰、十神) | **循环与条件逻辑**<br>1. 干支序列循环<br>2. 纳音查表<br>3. 旺衰状态判断<br>4. 十神关系映射 | 1. **数据完整性** (30年范围)<br>2. **关系映射正确性**<br>3. **计算效率** (批量处理) | `LiuNianCalculator.compute_liu_nian(start_year, gan_zhi)` | 1. 1997年甲子年开始<br>2. 2024年甲辰年开始<br>3. 随机年份随机年柱 |
| 8 | **五行数量统计** | `gan_zhi`: str (四柱字符串，如 "甲子 丙寅 戊辰 庚午") | 字符串 (五行统计，如 "木(2) 火(1) 土(3) 金(1) 水(0)") | **字符串处理与统计**<br>1. 干支解析<br>2. 天干五行映射<br>3. 地支藏干展开<br>4. 频次统计 | 1. **解析准确性** (含藏干)<br>2. **统计完整性**<br>3. **复杂案例处理** (重复五行) | `CalendarHandler._count_wuxing(gan_zhi)` | 1. 标准八字案例<br>2. 全阳干案例<br>3. 地支藏干复杂案例 |

## 🎯 任务设计要点

### **科学计算维度覆盖**
1. **高精度天文计算** (任务1-3)：涉及太阳黄经、儒略日等天体力学计算，考验浮点精度和迭代算法稳定性。
2. **复杂日期算法** (任务4-5)：干支计算包含闰年修正、世纪交界、特殊偏移等多重逻辑，测试条件分支处理能力。
3. **批量数据处理** (任务7)：30年流年信息生成，评估循环效率和内存管理。
4. **映射与统计** (任务8)：五行统计需要多层映射（天干→五行、地支→藏干），测试数据结构处理能力。

### **评估指标建议**
- **正确性**：与权威历法数据（如紫金山天文台节气表）对比。
- **精度**：浮点计算结果与高精度参考值（如NASA星历）的误差。
- **鲁棒性**：对异常输入（如超出年份范围、非法日期）的处理。
- **效率**：执行时间和内存占用（特别对批量计算任务）。

### **实验设计**
1. **输入多样性**：涵盖正常值、边界值、特殊年份（闰年、世纪交界）。
2. **黄金标准**：使用文件中已实现的算法作为基线（若已验证正确），或引用权威第三方数据。
3. **量化对比**：记录每个AI模型输出的数值结果、执行时间、错误率。

## 📊 推荐优先级

1. **必做核心任务**：任务1（节气时间）、任务4（日天干）——最具代表性，涵盖天文和传统历法。
2. **算法深度任务**：任务2（太阳黄经）——数学公式复杂，考验模型对科学公式的理解和实现。
3. **批量处理任务**：任务7（流年计算）——测试模型在循环和条件分支上的综合能力。
4. **快速验证任务**：任务6（起运岁数）——简单但可快速验证基本算术能力。

这些任务从简单算术到复杂天文算法均有覆盖，能够全面评估AI模型在**数值计算精度、算法实现能力、科学公式理解**等方面的表现，非常适合您的对比实验需求。