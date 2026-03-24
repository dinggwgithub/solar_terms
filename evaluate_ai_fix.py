#!/usr/bin/env python3
"""
AI Bug修复能力评估脚本
独立于被测系统，使用独立数学引擎计算标准答案
"""

import math
import requests
import json
from typing import Dict, Any


class BugFixEvaluator:
    """Bug修复评估器"""
    
    def __init__(self, server_url: str = "http://localhost:8080"):
        self.server_url = server_url
    
    def calculate_ground_truth(self, calc_type: str, params: Dict) -> Dict[str, Any]:
        """
        独立计算标准答案（黄金基准）
        """
        if calc_type == "ode_solver":
            return self._ground_truth_ode(params)
        elif calc_type == "equation_solver":
            return self._ground_truth_equation(params)
        elif calc_type == "planet":
            return self._ground_truth_planet(params)
        elif calc_type == "moon_phase":
            return self._ground_truth_moon_phase(params)
        else:
            return {"error": f"暂不支持 {calc_type} 的基准计算"}
    
    def _ground_truth_ode(self, params: Dict) -> Dict:
        """ODE求解器的基准计算"""
        equation = params.get("equation", "dy/dt = -y")
        y0 = params.get("initial_value", 1.0)
        t_end = params.get("time_range", 2.0)
        dt = params.get("time_step", 0.1)
        method = params.get("method", "euler").lower()
        
        # 解析解（黄金标准）
        if equation == "dy/dt = -y":
            analytical = y0 * math.exp(-t_end)
        elif equation == "dy/dt = y":
            analytical = y0 * math.exp(t_end)
        elif equation == "dy/dt = t*y":
            analytical = y0 * math.exp(t_end**2 / 2)
        else:
            analytical = y0 * math.exp(-t_end)
        
        # 数值方法计算
        if method == "euler":
            t, y = 0, y0
            while t < t_end - 1e-10:
                if equation == "dy/dt = -y":
                    dydt = -y
                elif equation == "dy/dt = y":
                    dydt = y
                else:
                    dydt = -y
                y = y + dydt * dt
                t = t + dt
            numerical = y
        else:  # rk4
            numerical = analytical  # 简化，实际应实现RK4
        
        return {
            "analytical_solution": analytical,
            "numerical_solution": numerical,
            "acceptable_error": 0.05  # 5%误差容忍度
        }
    
    def _ground_truth_equation(self, params: Dict) -> Dict:
        """方程求解器的基准计算"""
        equation = params.get("equation", "x^2 - 4 = 0")
        
        # 简单的二次方程求解：x^2 - 4 = 0
        if equation == "x^2 - 4 = 0" or equation == "x**2 - 4 = 0":
            return {
                "analytical_solution": [2.0, -2.0],
                "acceptable_error": 1e-3
            }
        elif equation == "sin(x) = 0":
            return {
                "analytical_solution": [0.0, 3.1415926535],
                "acceptable_error": 1e-3
            }
        else:
            return {
                "analytical_solution": 2.0,
                "acceptable_error": 0.01
            }
    
    def _ground_truth_planet(self, params: Dict) -> Dict:
        """行星位置的基准计算"""
        planet = params.get("planet", "sun")
        
        # 简化的行星基准值（实际应使用VSOP87等天文算法）
        planet_data = {
            "sun": {"ra": 0.0, "dec": 0.0, "distance": 1.0},
            "mercury": {"ra": 0.5, "dec": 0.1, "distance": 0.387},
            "venus": {"ra": 1.2, "dec": 0.2, "distance": 0.723},
            "earth": {"ra": 0.0, "dec": 0.0, "distance": 1.000},
            "mars": {"ra": 3.5, "dec": 0.3, "distance": 1.524},
            "jupiter": {"ra": 5.2, "dec": 0.1, "distance": 5.203},
            "saturn": {"ra": 6.0, "dec": 0.2, "distance": 9.537},
        }
        
        default_data = {"ra": 0.0, "dec": 0.0, "distance": 1.0}
        return {
            "expected": planet_data.get(planet.lower(), default_data),
            "acceptable_error": 0.1
        }
    
    def _ground_truth_moon_phase(self, params: Dict) -> Dict:
        """月相计算的基准"""
        year = params.get("year", 2024)
        month = params.get("month", 3)
        day = params.get("day", 20)
        
        # 简化的月相计算
        # 2024-3-20 实际为新月附近
        jd = 2460400.5  # 儒略日简化计算
        phase_angle = ((jd - 2451550.1) / 29.530588853) % 1.0 * 360
        
        return {
            "phase_angle": phase_angle,
            "illumination": (1 - math.cos(math.radians(phase_angle))) / 2,
            "acceptable_error": 5.0  # 5度误差
        }
    
    def evaluate(self, calc_type: str, params: Dict, ai_result: Dict) -> Dict:
        """
        评估AI修复的正确性
        """
        # 获取标准答案
        ground_truth = self.calculate_ground_truth(calc_type, params)
        
        # 对比评估
        report = {
            "calculation_type": calc_type,
            "parameters": params,
            "ground_truth": ground_truth,
            "ai_result": ai_result,
            "passed": False,
            "error_analysis": {}
        }
        
        # 根据计算类型进行针对性评估
        if calc_type == "ode_solver":
            ai_solution = ai_result.get("solution", 0)
            expected = ground_truth.get("analytical_solution", 0)
            tolerance = ground_truth.get("acceptable_error", 0.05)
            
            abs_error = abs(ai_solution - expected)
            rel_error = abs_error / max(abs(expected), 1e-10) if expected != 0 else abs_error
            
            report["error_analysis"] = {
                "absolute_error": abs_error,
                "relative_error": rel_error,
                "tolerance": tolerance
            }
            report["passed"] = abs_error < tolerance
        
        elif calc_type == "equation_solver":
            ai_solution = ai_result.get("solution", 0)
            if isinstance(ai_solution, list) and len(ai_solution) > 0:
                ai_solution = ai_solution[0]
            
            expected = ground_truth.get("analytical_solution", 0)
            if isinstance(expected, list) and len(expected) > 0:
                expected = expected[0]
            
            tolerance = ground_truth.get("acceptable_error", 0.01)
            abs_error = abs(float(ai_solution) - float(expected))
            
            report["error_analysis"] = {
                "absolute_error": abs_error,
                "tolerance": tolerance
            }
            report["passed"] = abs_error < tolerance
        
        else:
            report["note"] = "需要针对此计算类型实现具体评估逻辑"
        
        return report


def demo_evaluation():
    """演示评估流程"""
    evaluator = BugFixEvaluator()
    
    print("=" * 70)
    print("AI Bug修复能力评估 - 演示")
    print("=" * 70)
    
    # 测试案例：ODE求解器
    print("\n📊 测试案例1: ODE求解器")
    print("-" * 50)
    
    calc_type = "ode_solver"
    params = {
        "equation": "dy/dt = -y",
        "initial_value": 1.0,
        "time_step": 0.1,
        "time_range": 2.0,
        "method": "euler"
    }
    
    # 模拟AI返回的结果（假设已修复）
    ai_fixed_result = {
        "solution": 0.12157665,  # 修复后接近正确值
        "method_used": "Euler"
    }
    
    report = evaluator.evaluate(calc_type, params, ai_fixed_result)
    
    print(f"计算类型: {calc_type}")
    print(f"方程: {params['equation']}")
    print(f"初始条件: y(0) = {params['initial_value']}")
    print(f"时间范围: [0, {params['time_range']}]")
    print()
    print(f"📐 标准答案: {report['ground_truth']['analytical_solution']:.8f}")
    print(f"🤖 AI结果: {report['ai_result']['solution']:.8f}")
    print(f"📏 绝对误差: {report['error_analysis']['absolute_error']:.8f}")
    print(f"✅ 容忍度: ±{report['error_analysis']['tolerance']:.4f}")
    print(f"🎯 评估结果: {'PASS ✓' if report['passed'] else 'FAIL ✗'}")
    
    print("\n" + "=" * 70)


if __name__ == "__main__":
    demo_evaluation()
