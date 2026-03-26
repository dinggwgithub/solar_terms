package main

import (
	"fmt"
	"log"
	"scientific_calc/internal/api"
	"scientific_calc/internal/calculator"

	_ "scientific_calc/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 科学计算工具集 API
// @version 1.0
// @description 功能丰富的Go语言科学计算库，提供天文计算、历法转换、方程求解等多种科学计算功能。
// @host localhost:8080
// @BasePath /
func main() {
	// 初始化计算器管理器
	calculatorManager := calculator.NewCalculatorManager()

	// 注册所有科学计算器
	registerCalculators(calculatorManager)

	// 初始化API处理器
	apiHandler := api.NewAPIHandler(calculatorManager)

	// 初始化Gin路由
	r := gin.Default()

	// 设置Swagger UI路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册API路由
	apiHandler.RegisterRoutes(r)

	// 启动服务器
	fmt.Println("========================================")
	fmt.Println("科学计算工具集 ScientificCalc v1.0")
	fmt.Println("========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	fmt.Println("健康检查: http://localhost:8080/api/health")
	fmt.Println("系统信息: http://localhost:8080/api/system-info")
	fmt.Println("========================================")
	fmt.Println("支持的计算任务:")
	printSupportedCalculations(calculatorManager)
	fmt.Println("========================================")

	log.Fatal(r.Run(":8080"))
}

// registerCalculators 注册所有科学计算器
func registerCalculators(manager *calculator.CalculatorManager) {
	// 注册节气计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeSolarTerm,
		calculator.NewSolarTermCalculator(),
	)

	// 注册干支计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeGanZhi,
		calculator.NewGanZhiCalculator(),
	)

	// 注册天文计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeAstronomy,
		calculator.NewAstronomyCalculator(),
	)

	// 注册起运岁数计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeStartingAge,
		calculator.NewStartingAgeCalculator(),
	)

	// 注册农历转换计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeLunar,
		calculator.NewLunarCalculator(),
	)

	// 注册行星位置计算器
	manager.RegisterCalculator(
		calculator.CalculationTypePlanet,
		calculator.NewPlanetCalculator(),
	)

	// 注册星曜推算计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeStar,
		calculator.NewStarCalculator(),
	)

	// 注册日出日落时间计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeSunriseSunset,
		calculator.NewSunriseSunsetCalculator(),
	)

	// 注册月相计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeMoonPhase,
		calculator.NewMoonPhaseCalculator(),
	)

	// 注册方程求解器
	manager.RegisterCalculator(
		calculator.CalculationTypeEquationSolver,
		calculator.NewEquationSolverCalculator(),
	)

	// 注册符号计算器
	manager.RegisterCalculator(
		calculator.CalculationTypeSymbolicCalc,
		calculator.NewSymbolicCalcCalculator(),
	)

	// 注册微分方程求解器
	manager.RegisterCalculator(
		calculator.CalculationTypeODESolver,
		calculator.NewODESolverCalculator(),
	)

	// 注册修复版方程求解器
	manager.RegisterCalculator(
		calculator.CalculationTypeEquationSolverFixed,
		calculator.NewFixedEquationSolverCalculator(),
	)
}

// printSupportedCalculations 打印支持的计算任务
func printSupportedCalculations(manager *calculator.CalculatorManager) {
	supportedTypes := manager.GetSupportedCalculationTypes()

	for _, calcType := range supportedTypes {
		calculator, exists := manager.GetCalculator(calcType)
		if exists {
			fmt.Printf("  • %s: %s\n", calcType.String(), calculator.Description())
		}
	}
}

// 需要创建的计算器实现（占位符）

// SolarTermCalculator 节气计算器
type SolarTermCalculator struct {
	*calculator.BaseCalculator
}

func NewSolarTermCalculator() *SolarTermCalculator {
	return &SolarTermCalculator{
		BaseCalculator: calculator.NewBaseCalculator(
			"solar_term",
			"节气精确时间计算，基于天文算法",
		),
	}
}

func (c *SolarTermCalculator) Calculate(params interface{}) (interface{}, error) {
	// 实现节气计算
	return "2024-02-04 12:00:00", nil
}

func (c *SolarTermCalculator) Validate(params interface{}) error {
	return nil
}

// GanZhiCalculator 干支计算器
type GanZhiCalculator struct {
	*calculator.BaseCalculator
}

func NewGanZhiCalculator() *GanZhiCalculator {
	return &GanZhiCalculator{
		BaseCalculator: calculator.NewBaseCalculator(
			"ganzhi",
			"干支计算，支持年月日时干支",
		),
	}
}

func (c *GanZhiCalculator) Calculate(params interface{}) (interface{}, error) {
	// 实现干支计算
	return map[string]string{
		"gan_year":  "甲",
		"zhi_year":  "辰",
		"gan_month": "丙",
		"zhi_month": "寅",
		"gan_day":   "甲",
		"zhi_day":   "子",
		"gan_time":  "甲",
		"zhi_time":  "子",
	}, nil
}

func (c *GanZhiCalculator) Validate(params interface{}) error {
	return nil
}

// AstronomyCalculator 天文计算器
type AstronomyCalculator struct {
	*calculator.BaseCalculator
}

func NewAstronomyCalculator() *AstronomyCalculator {
	return &AstronomyCalculator{
		BaseCalculator: calculator.NewBaseCalculator(
			"astronomy",
			"天文黄经计算，基于VSOP87理论",
		),
	}
}

func (c *AstronomyCalculator) Calculate(params interface{}) (interface{}, error) {
	// 实现天文计算
	return map[string]float64{
		"sun_longitude":      280.123,
		"julian_date":        2459580.5,
		"apparent_longitude": 280.123,
		"true_longitude":     280.123,
		"mean_longitude":     280.123,
		"mean_anomaly":       357.529,
		"equation_of_center": 1.914,
		"nutation":           0.004,
	}, nil
}

func (c *AstronomyCalculator) Validate(params interface{}) error {
	return nil
}

// StartingAgeCalculator 起运岁数计算器
type StartingAgeCalculator struct {
	*calculator.BaseCalculator
}

func NewStartingAgeCalculator() *StartingAgeCalculator {
	return &StartingAgeCalculator{
		BaseCalculator: calculator.NewBaseCalculator(
			"starting_age",
			"起运岁数计算，基于生辰八字",
		),
	}
}

func (c *StartingAgeCalculator) Calculate(params interface{}) (interface{}, error) {
	// 实现起运岁数计算
	return "5岁", nil
}

func (c *StartingAgeCalculator) Validate(params interface{}) error {
	return nil
}
