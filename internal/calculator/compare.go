package calculator

import (
	"fmt"
)

type CompareCalculator struct {
	*BaseCalculator
	oldCalc *PlanetCalculator
	newCalc *PlanetCalculatorFixed
}

func NewCompareCalculator() *CompareCalculator {
	return &CompareCalculator{
		BaseCalculator: NewBaseCalculator(
			"compare",
			"新旧计算结果对比",
		),
		oldCalc: NewPlanetCalculator(),
		newCalc: NewPlanetCalculatorFixed(),
	}
}

type CompareResult struct {
	OldResult *PlanetPosition      `json:"old_result"`
	NewResult *PlanetPositionFixed `json:"new_result"`
	Diff      *PositionDiff        `json:"diff"`
}

type PositionDiff struct {
	RightAscension float64 `json:"right_ascension"`
	Declination    float64 `json:"declination"`
	Distance       float64 `json:"distance"`
	Magnitude      float64 `json:"magnitude"`
	Phase          float64 `json:"phase"`
	Elongation     float64 `json:"elongation"`
}

func (c *CompareCalculator) Calculate(params interface{}) (interface{}, error) {
	oldResult, err := c.oldCalc.Calculate(params)
	if err != nil {
		return nil, fmt.Errorf("旧计算器执行失败: %v", err)
	}

	newResult, err := c.newCalc.Calculate(params)
	if err != nil {
		return nil, fmt.Errorf("新计算器执行失败: %v", err)
	}

	oldPos, ok := oldResult.(*PlanetPosition)
	if !ok {
		return nil, fmt.Errorf("旧计算器返回类型错误")
	}

	newPos, ok := newResult.(*PlanetPositionFixed)
	if !ok {
		return nil, fmt.Errorf("新计算器返回类型错误")
	}

	diff := &PositionDiff{
		RightAscension: newPos.RightAscension - oldPos.RightAscension,
		Declination:    newPos.Declination - oldPos.Declination,
		Distance:       newPos.Distance - oldPos.Distance,
		Magnitude:      newPos.Magnitude - oldPos.Magnitude,
		Phase:          newPos.Phase - oldPos.Phase,
		Elongation:     newPos.Elongation - oldPos.Elongation,
	}

	return &CompareResult{
		OldResult: oldPos,
		NewResult: newPos,
		Diff:      diff,
	}, nil
}

func (c *CompareCalculator) Validate(params interface{}) error {
	return c.oldCalc.Validate(params)
}

func (c *CompareCalculator) Description() string {
	return c.description
}
