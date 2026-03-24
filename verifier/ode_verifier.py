#!/usr/bin/env python3
"""
独立验证脚本：用于评估AI修复的ODE求解器结果正确性
完全独立于被测系统，使用自己的数学实现来计算标准答案
"""

import math
import numpy as np
from scipy.integrate import solve_ivp
from typing import Dict, List, Union


class ODEVerifier:
    """ODE求解器独立验证器"""
    
    def solve_euler(self, equation: str, y0: float, t0: float, t_end: float, dt: float) -> Dict:
        """
        使用欧拉方法求解ODE（独立实现）
        """
        t_points = []
        y_points = []
        t = t0
        y = y0
        
        t_points.append(t)
        y_points.append(y)
        
        # 解析简单的 dy/dt = -y 形式
        # 这里简化处理，实际应用需要表达式解析器
        while t < t_end:
            if equation == "dy/dt = -y":
                dydt = -y
            elif equation == "dy/dt = y":
                dydt = y
            elif equation == "dy/dt = t*y":
                dydt = t * y
            else:
                # 默认用解析解计算导数
                dydt = self.analytical_derivative(equation, t, y)
            
            y = y + dydt * dt
            t = t + dt
            
            t_points.append(float(t))
            y_points.append(float(y))
        
        return {
            "solution": float(y_points[-1]),
            "time_points": t_points,
            "solution_path": y_points,
            "method_used": "Euler (Verifier)"
        }
    
    def solve_rk4(self, equation: str, y0: float, t0: float, t_end: float, dt: float) -> Dict:
        """
        使用RK4方法求解ODE（独立实现）
        """
        def f(t, y):
            if equation == "dy/dt = -y":
                return -y
            elif equation == "dy/dt = y":
                return y
            elif equation == "dy/dt = t*y":
                return t * y
            return -y  # 默认
        
        t_points = []
        y_points = []
        t = t0
        y = y0
        
        t_points.append(t)
        y_points.append(y)
        
        while t < t_end:
            k1 = dt * f(t, y)
            k2 = dt * f(t + dt/2, y + k1/2)
            k3 = dt * f(t + dt/2, y + k2/2)
            k4 = dt * f(t + dt, y + k3)
            
            y = y + (k1 + 2*k2 + 2*k3 + k4) / 6
            t = t + dt
            
            t_points.append(float(t))
            y_points.append(float(y))
        
        return {
            "solution": float(y_points[-1]),
            "time_points": t_points,
            "solution_path": y_points,
            "method_used": "RK4 (Verifier)"
        }
    
    def analytical_solution(self, equation: str, y0: float, t: float) -> float:
        """
        计算解析解（黄金标准）
        """
        if equation == "dy/dt = -y":
            return y0 * math.exp(-t)
        elif equation == "dy/dt = y":
            return y0 * math.exp(t)
        elif equation == "dy/dt = t*y":
            return y0 * math.exp(t**2 / 2)
        return y0 * math.exp(-t)  # 默认
    
    def analytical_derivative(self, equation: str, t: float, y: float) -> float:
        """
        计算导数的解析值
        """
        if equation == "dy/dt = -y":
            return -y
        elif equation == "dy/dt = y":
            return y
        elif equation == "dy/dt = t*y":
            return t * y
        return -y  # 默认
    
    def verify_result(self, 
                     ai_result: Dict, 
                     equation: str, 
                     y0: float, 
                     t_end: float,
                     tolerance: float = 1e-3) -> Dict:
        """
        验证AI修复的结果正确性
        
        Args:
            ai_result: AI修复后返回的结果字典
            equation: ODE方程
            y0: 初始值
            t_end: 结束时间
            tolerance: 误差容忍度
            
        Returns:
            验证报告
        """
        # 计算解析解（黄金标准）
        analytical_value = self.analytical_solution(equation, y0, t_end)
        
        # AI返回值
        ai_value = ai_result.get("solution", 0)
        
        # 计算绝对误差和相对误差
        abs_error = abs(ai_value - analytical_value)
        rel_error = abs_error / max(abs(analytical_value), 1e-10)
        
        # 判定是否修复成功
        is_fixed = abs_error < tolerance
        
        return {
            "analytical_solution": analytical_value,
            "ai_returned_value": ai_value,
            "absolute_error": abs_error,
            "relative_error": rel_error,
            "tolerance": tolerance,
            "is_fixed": is_fixed,
            "verdict": "PASS" if is_fixed else f"FAIL (误差超出容忍度: {abs_error:.6f})"
        }


def evaluate_ode_fix(candidate_result: Dict) -> Dict:
    """
    评估ODE求解器修复结果的入口函数
    
    Args:
        candidate_result: 候选系统返回的结果
        
    Returns:
        评估报告
    """
    verifier = ODEVerifier()
    
    # 测试用例参数
    test_cases = [
        {
            "equation": "dy/dt = -y",
            "initial_value": 1.0,
            "time_end": 2.0,
            "description": "指数衰减方程"
        },
        {
            "equation": "dy/dt = y",
            "initial_value": 1.0,
            "time_end": 1.0,
            "description": "指数增长方程"
        }
    ]
    
    results = []
    all_passed = True
    
    for test_case in test_cases:
        report = verifier.verify_result(
            candidate_result,
            test_case["equation"],
            test_case["initial_value"],
            test_case["time_end"]
        )
        report["test_case"] = test_case
        results.append(report)
        
        if not report["is_fixed"]:
            all_passed = False
    
    return {
        "overall_verdict": "PASS" if all_passed else "FAIL",
        "passed_cases": sum(1 for r in results if r["is_fixed"]),
        "total_cases": len(results),
        "details": results
    }


if __name__ == "__main__":
    # 示例用法
    verifier = ODEVerifier()
    
    # 测试独立求解
    print("=" * 60)
    print("独立验证脚本示例")
    print("=" * 60)
    
    equation = "dy/dt = -y"
    y0 = 1.0
    t_end = 2.0
    
    analytical = verifier.analytical_solution(equation, y0, t_end)
    euler_result = verifier.solve_euler(equation, y0, 0, t_end, 0.1)
    rk4_result = verifier.solve_rk4(equation, y0, 0, t_end, 0.1)
    
    print(f"\n测试方程: {equation}")
    print(f"初始条件: y(0) = {y0}, t ∈ [0, {t_end}]")
    print(f"\n解析解（黄金标准）: {analytical:.10f}")
    print(f"欧拉方法（独立实现）: {euler_result['solution']:.10f}, 误差: {abs(euler_result['solution'] - analytical):.10f}")
    print(f"RK4方法（独立实现）: {rk4_result['solution']:.10f}, 误差: {abs(rk4_result['solution'] - analytical):.10f}")
    
    print("\n" + "=" * 60)
