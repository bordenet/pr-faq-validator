#!/usr/bin/env python3
"""
Analyze 20-round optimization experiment results.
Generates detailed statistics, convergence analysis, and recommendations.
"""

import json
import sys
from pathlib import Path
from typing import Dict, List, Any


def load_results(results_dir: Path) -> Dict[str, Any]:
    """Load optimization results."""
    results_file = results_dir / "optimization_final_results.json"
    with open(results_file, 'r') as f:
        return json.load(f)


def analyze_convergence(iteration_history: List[Dict[str, Any]]) -> Dict[str, Any]:
    """Analyze convergence behavior."""
    improvements = [it for it in iteration_history if it.get('improved', False)]
    
    # Find last improvement
    last_improvement_idx = -1
    for i, it in enumerate(iteration_history):
        if it['score'] == it['best_score']:
            last_improvement_idx = i
    
    # Calculate post-convergence iterations
    total_iterations = len(iteration_history)
    post_convergence = total_iterations - last_improvement_idx - 1 if last_improvement_idx >= 0 else 0
    
    # Calculate score variance in plateau
    if post_convergence > 0:
        plateau_scores = [it['score'] for it in iteration_history[last_improvement_idx + 1:]]
        plateau_min = min(plateau_scores) if plateau_scores else 0
        plateau_max = max(plateau_scores) if plateau_scores else 0
        plateau_avg = sum(plateau_scores) / len(plateau_scores) if plateau_scores else 0
        plateau_variance = plateau_max - plateau_min
    else:
        plateau_min = plateau_max = plateau_avg = plateau_variance = 0
    
    return {
        "converged": post_convergence >= 5,
        "convergence_iteration": last_improvement_idx + 1 if last_improvement_idx >= 0 else None,
        "post_convergence_iterations": post_convergence,
        "wasted_iterations": post_convergence if post_convergence >= 5 else 0,
        "plateau_variance": plateau_variance,
        "plateau_avg": plateau_avg,
        "plateau_min": plateau_min,
        "plateau_max": plateau_max,
        "total_improvements": len([it for it in iteration_history if it['score'] == it['best_score']])
    }


def calculate_improvement_deltas(iteration_history: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """Calculate improvement deltas for each successful iteration."""
    deltas = []
    prev_best = iteration_history[0]['best_score'] if iteration_history else 0
    
    for it in iteration_history:
        if it['score'] == it['best_score'] and it['score'] > prev_best:
            delta = it['score'] - prev_best
            deltas.append({
                "iteration": it['iteration'],
                "score": it['score'],
                "delta": delta,
                "percent": (delta / prev_best * 100) if prev_best > 0 else 0
            })
            prev_best = it['score']
    
    return deltas


def print_analysis(results: Dict[str, Any]):
    """Print comprehensive analysis."""
    print("=" * 80)
    print("20-ROUND OPTIMIZATION EXPERIMENT ANALYSIS")
    print("=" * 80)
    print()
    
    # Basic metrics
    print("📊 BASIC METRICS")
    print(f"  Project: {results['project']}")
    print(f"  Max Iterations: {results['max_iterations']}")
    print(f"  Baseline Score: {results['baseline_score']:.2f}")
    print(f"  Final Score: {results['final_score']:.2f}")
    print(f"  Improvement: +{results['improvement']:.2f} points ({results['improvement']/results['baseline_score']*100:.2f}%)")
    print()
    
    # Convergence analysis
    convergence = analyze_convergence(results['iteration_history'])
    print("🎯 CONVERGENCE ANALYSIS")
    print(f"  Converged: {'Yes' if convergence['converged'] else 'No'}")
    if convergence['convergence_iteration']:
        print(f"  Convergence Iteration: {convergence['convergence_iteration']}")
        print(f"  Post-Convergence Iterations: {convergence['post_convergence_iterations']}")
        print(f"  Wasted Iterations: {convergence['wasted_iterations']} ({convergence['wasted_iterations']/results['max_iterations']*100:.1f}%)")
    print(f"  Total Improvements: {convergence['total_improvements']}")
    print(f"  Plateau Variance: ±{convergence['plateau_variance']:.2f} points")
    print(f"  Plateau Average: {convergence['plateau_avg']:.2f}")
    print()
    
    # Improvement deltas
    deltas = calculate_improvement_deltas(results['iteration_history'])
    print("📈 IMPROVEMENT TIMELINE")
    for delta in deltas:
        print(f"  Iteration {delta['iteration']:2d}: {delta['score']:.2f} (+{delta['delta']:.2f}, +{delta['percent']:.2f}%)")
    print()
    
    # Recommendations
    print("💡 RECOMMENDATIONS")
    if convergence['wasted_iterations'] > 0:
        print(f"  ⚠️  CRITICAL: Implement auto-stop after 5 consecutive non-improvements")
        print(f"     Potential savings: {convergence['wasted_iterations']} iterations ({convergence['wasted_iterations']/results['max_iterations']*100:.1f}%)")
    
    if len(deltas) > 1:
        # Check for diminishing returns
        if deltas[-1]['delta'] < deltas[0]['delta'] / 2:
            print(f"  ⚠️  Diminishing returns detected:")
            print(f"     First improvement: +{deltas[0]['delta']:.2f}")
            print(f"     Last improvement: +{deltas[-1]['delta']:.2f}")
            print(f"     Suggests approaching local optimum")
    
    if convergence['plateau_variance'] > 5:
        print(f"  ℹ️  High plateau variance (±{convergence['plateau_variance']:.2f}) suggests:")
        print(f"     - Evolutionary mutations exploring diverse prompt space")
        print(f"     - Consider reducing mutation strength near convergence")
    
    print()
    
    # Score distribution
    print("📊 SCORE DISTRIBUTION")
    scores = [it['score'] for it in results['iteration_history']]
    print(f"  Min: {min(scores):.2f}")
    print(f"  Max: {max(scores):.2f}")
    print(f"  Avg: {sum(scores)/len(scores):.2f}")
    print(f"  Range: {max(scores) - min(scores):.2f}")
    print()
    
    print("=" * 80)


def main():
    """Main entry point."""
    if len(sys.argv) > 1:
        results_dir = Path(sys.argv[1])
    else:
        results_dir = Path("prompt_tuning_results_pr-faq-validator")
    
    if not results_dir.exists():
        print(f"Error: Results directory not found: {results_dir}")
        sys.exit(1)
    
    results = load_results(results_dir)
    print_analysis(results)


if __name__ == "__main__":
    main()

