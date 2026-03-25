package main

import (
	"fmt"
	"scientific_calc/internal/calculator"
)

func main() {
	// 创建修复后的计算器
	calc := calculator.NewSunriseSunsetCalculatorFixed()

	// 测试参数：北京 2026年3月23日
	// 注意：参数需要使用float64类型以匹配parseParams函数的期望
	params := map[string]interface{}{
		"year":      float64(2026),
		"month":     float64(3),
		"day":       float64(23),
		"longitude": 116.4,
		"latitude":  39.9,
		"timezone":  8.0,
	}

	// 执行计算
	result, err := calc.Calculate(params)
	if err != nil {
		fmt.Printf("计算错误: %v\n", err)
		return
	}

	// 输出结果
	sunriseResult, ok := result.(*calculator.SunriseSunsetResult)
	if !ok {
		fmt.Println("结果类型错误")
		return
	}

	fmt.Println("=== 北京 2026年3月23日 日出日落时间 ===")
	fmt.Printf("日期: %s\n", sunriseResult.Date)
	fmt.Printf("日出: %s\n", sunriseResult.Sunrise)
	fmt.Printf("日落: %s\n", sunriseResult.Sunset)
	fmt.Printf("正午: %s\n", sunriseResult.SolarNoon)
	fmt.Printf("昼长: %.2f小时\n", sunriseResult.DayLength)
	fmt.Println()
	fmt.Println("=== 晨昏蒙影 ===")
	fmt.Printf("民用晨光: %s\n", sunriseResult.CivilTwilight.Morning)
	fmt.Printf("民用暮光: %s\n", sunriseResult.CivilTwilight.Evening)
	fmt.Printf("航海晨光: %s\n", sunriseResult.NauticalTwilight.Morning)
	fmt.Printf("航海暮光: %s\n", sunriseResult.NauticalTwilight.Evening)
	fmt.Printf("天文晨光: %s\n", sunriseResult.AstronomicalTwilight.Morning)
	fmt.Printf("天文暮光: %s\n", sunriseResult.AstronomicalTwilight.Evening)
	fmt.Println()
	fmt.Println("预期日出时间应在 06:00-06:30 之间")
	fmt.Println("预期日落时间应在 18:00-18:30 之间")
}
