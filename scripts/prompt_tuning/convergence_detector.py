"""
Convergence detection for evolutionary optimization.

Detects when optimization has converged and should stop early.
"""

from typing import List, Dict, Any, Optional
from dataclasses import dataclass


@dataclass
class ConvergenceConfig:
    """Configuration for convergence detection."""
    
    # Stop after N consecutive iterations without improvement
    no_improvement_threshold: int = 5
    
    # Minimum improvement to count as significant (percentage)
    min_improvement_percent: float = 0.1
    
    # Enable early stopping
    enable_early_stop: bool = True
    
    # Track plateau variance
    track_plateau_variance: bool = True


class ConvergenceDetector:
    """Detects convergence in optimization experiments."""
    
    def __init__(self, config: Optional[ConvergenceConfig] = None):
        self.config = config or ConvergenceConfig()
        self.iteration_history: List[Dict[str, Any]] = []
        self.best_score: float = 0.0
        self.no_improvement_count: int = 0
        self.converged: bool = False
        self.convergence_iteration: Optional[int] = None
        
    def update(self, iteration: int, score: float, improved: bool) -> Dict[str, Any]:
        """
        Update convergence state with new iteration results.
        
        Args:
            iteration: Current iteration number
            score: Current iteration score
            improved: Whether this iteration improved the best score
            
        Returns:
            Dict with convergence status and recommendations
        """
        # Track iteration
        self.iteration_history.append({
            "iteration": iteration,
            "score": score,
            "improved": improved,
            "best_score": max(self.best_score, score)
        })
        
        # Update best score
        if score > self.best_score:
            improvement_percent = ((score - self.best_score) / self.best_score * 100) if self.best_score > 0 else 100
            
            # Check if improvement is significant
            if improvement_percent >= self.config.min_improvement_percent:
                self.best_score = score
                self.no_improvement_count = 0
            else:
                self.no_improvement_count += 1
        else:
            self.no_improvement_count += 1
        
        # Check for convergence
        if self.no_improvement_count >= self.config.no_improvement_threshold:
            if not self.converged:
                self.converged = True
                self.convergence_iteration = iteration
        
        return self.get_status()
    
    def should_stop(self) -> bool:
        """Check if optimization should stop early."""
        if not self.config.enable_early_stop:
            return False
        
        return self.converged
    
    def get_status(self) -> Dict[str, Any]:
        """Get current convergence status."""
        plateau_stats = self._calculate_plateau_stats()
        
        return {
            "converged": self.converged,
            "convergence_iteration": self.convergence_iteration,
            "no_improvement_count": self.no_improvement_count,
            "best_score": self.best_score,
            "should_stop": self.should_stop(),
            "plateau_variance": plateau_stats["variance"],
            "plateau_avg": plateau_stats["avg"],
            "total_iterations": len(self.iteration_history),
            "wasted_iterations": self._calculate_wasted_iterations()
        }
    
    def _calculate_plateau_stats(self) -> Dict[str, float]:
        """Calculate statistics for plateau period."""
        if not self.converged or not self.convergence_iteration:
            return {"variance": 0.0, "avg": 0.0, "min": 0.0, "max": 0.0}
        
        # Get scores after convergence
        plateau_scores = [
            it["score"] for it in self.iteration_history
            if it["iteration"] > self.convergence_iteration
        ]
        
        if not plateau_scores:
            return {"variance": 0.0, "avg": 0.0, "min": 0.0, "max": 0.0}
        
        return {
            "variance": max(plateau_scores) - min(plateau_scores),
            "avg": sum(plateau_scores) / len(plateau_scores),
            "min": min(plateau_scores),
            "max": max(plateau_scores)
        }
    
    def _calculate_wasted_iterations(self) -> int:
        """Calculate number of wasted iterations after convergence."""
        if not self.converged or not self.convergence_iteration:
            return 0
        
        total = len(self.iteration_history)
        return total - self.convergence_iteration
    
    def get_recommendations(self) -> List[str]:
        """Get recommendations based on convergence analysis."""
        recommendations = []
        
        if self.converged:
            wasted = self._calculate_wasted_iterations()
            if wasted > 0:
                recommendations.append(
                    f"Convergence detected at iteration {self.convergence_iteration}. "
                    f"Could have stopped {wasted} iterations earlier."
                )
        
        plateau_stats = self._calculate_plateau_stats()
        if plateau_stats["variance"] > 5.0:
            recommendations.append(
                f"High plateau variance (±{plateau_stats['variance']:.2f}) suggests "
                "mutations are too aggressive. Consider reducing mutation strength."
            )
        
        if len(self.iteration_history) >= 3:
            # Check for diminishing returns
            improvements = [
                it for it in self.iteration_history
                if it["improved"] and it["score"] > 0
            ]
            
            if len(improvements) >= 2:
                first_delta = improvements[0]["score"] - (improvements[0]["best_score"] - improvements[0]["score"])
                last_delta = improvements[-1]["score"] - improvements[-2]["best_score"]
                
                if last_delta < first_delta / 2:
                    recommendations.append(
                        "Diminishing returns detected. Consider more aggressive mutations "
                        "or different optimization strategy."
                    )
        
        return recommendations
    
    def print_summary(self):
        """Print convergence summary."""
        status = self.get_status()
        recommendations = self.get_recommendations()
        
        print("\n" + "="*80)
        print("CONVERGENCE ANALYSIS")
        print("="*80)
        
        print(f"\nConverged: {'Yes' if status['converged'] else 'No'}")
        if status['converged']:
            print(f"Convergence Iteration: {status['convergence_iteration']}")
            print(f"Wasted Iterations: {status['wasted_iterations']}")
            print(f"Plateau Variance: ±{status['plateau_variance']:.2f}")
            print(f"Plateau Average: {status['plateau_avg']:.2f}")
        
        print(f"\nBest Score: {status['best_score']:.2f}")
        print(f"Total Iterations: {status['total_iterations']}")
        print(f"No Improvement Count: {status['no_improvement_count']}")
        
        if recommendations:
            print("\nRecommendations:")
            for i, rec in enumerate(recommendations, 1):
                print(f"{i}. {rec}")
        
        print("="*80 + "\n")

