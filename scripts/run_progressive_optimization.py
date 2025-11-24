#!/usr/bin/env python3
"""
Run successive progressive optimization rounds.

Each round starts from the best prompt of the previous round,
allowing for progressive improvement over multiple experiments.
"""

import json
import subprocess
import sys
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Any


class ProgressiveOptimizer:
    """Runs multiple successive optimization rounds."""

    def __init__(self, project: str, rounds: int, iterations_per_round: int = 20):
        self.project = project
        self.rounds = rounds
        self.iterations_per_round = iterations_per_round
        self.results_dir = Path.cwd() / f"prompt_tuning_results_{project}"
        self.prompts_dir = Path.cwd() / "prompts"
        self.round_history = []

    def run_single_round(self, round_num: int) -> Dict[str, Any]:
        """Run a single optimization round."""
        print(f"\n{'='*80}")
        print(f"ROUND {round_num}/{self.rounds}")
        print(f"{'='*80}\n")

        # Run optimization
        result = subprocess.run(
            [
                sys.executable,
                "scripts/run_autonomous_experiment.py",
                "--project", self.project,
                "--iterations", str(self.iterations_per_round)
            ],
            capture_output=True,
            text=True,
            timeout=300
        )

        if result.returncode != 0:
            print(f"❌ Round {round_num} failed:")
            print(result.stderr)
            return None

        # Load results
        results_file = self.results_dir / "optimization_final_results.json"
        if not results_file.exists():
            print(f"❌ Results file not found: {results_file}")
            return None

        with open(results_file, 'r') as f:
            round_results = json.load(f)

        round_data = {
            "round": round_num,
            "baseline_score": round_results.get("baseline_score"),
            "final_score": round_results.get("final_score"),
            "improvement": round_results.get("improvement"),
            "iterations": len(round_results.get("iteration_history", [])),
            "converged": round_results.get("convergence", {}).get("converged", False),
            "convergence_iteration": round_results.get("convergence", {}).get("convergence_iteration"),
        }

        self.round_history.append(round_data)
        return round_data

    def run_progressive_optimization(self) -> Dict[str, Any]:
        """Run all rounds progressively."""
        print(f"\n{'='*80}")
        print(f"PROGRESSIVE OPTIMIZATION: {self.rounds} ROUNDS")
        print(f"Project: {self.project}")
        print(f"Iterations per round: {self.iterations_per_round}")
        print(f"{'='*80}\n")

        initial_score = None
        final_score = None

        for round_num in range(1, self.rounds + 1):
            round_data = self.run_single_round(round_num)

            if round_data is None:
                print(f"❌ Stopping at round {round_num} due to error")
                break

            if round_num == 1:
                initial_score = round_data["baseline_score"]

            final_score = round_data["final_score"]

            # Print round summary
            print(f"\n📊 Round {round_num} Summary:")
            print(f"   Baseline: {round_data['baseline_score']:.2f}")
            print(f"   Final:    {round_data['final_score']:.2f}")
            print(f"   Change:   {round_data['improvement']:+.2f}")
            print(f"   Converged: {'Yes' if round_data['converged'] else 'No'} "
                  f"(iteration {round_data['convergence_iteration']})")

        # Generate final report
        return self.generate_final_report(initial_score, final_score)

    def generate_final_report(self, initial_score: float, final_score: float) -> Dict[str, Any]:
        """Generate comprehensive report of all rounds."""
        print(f"\n{'='*80}")
        print(f"PROGRESSIVE OPTIMIZATION COMPLETE")
        print(f"{'='*80}\n")

        total_improvement = final_score - initial_score if initial_score and final_score else 0
        improvement_pct = (total_improvement / initial_score * 100) if initial_score else 0

        print(f"📈 OVERALL RESULTS:")
        print(f"   Initial Score (Round 1 baseline): {initial_score:.2f}")
        print(f"   Final Score (Round {self.rounds}):   {final_score:.2f}")
        print(f"   Total Improvement:                {total_improvement:+.2f} ({improvement_pct:+.2f}%)")
        print(f"   Rounds Completed:                 {len(self.round_history)}/{self.rounds}")

        print(f"\n📊 ROUND-BY-ROUND BREAKDOWN:")
        print(f"   {'Round':<8} {'Baseline':<10} {'Final':<10} {'Change':<10} {'Converged':<12}")
        print(f"   {'-'*60}")

        for round_data in self.round_history:
            converged_str = f"Yes (iter {round_data['convergence_iteration']})" if round_data['converged'] else "No"
            print(f"   {round_data['round']:<8} "
                  f"{round_data['baseline_score']:<10.2f} "
                  f"{round_data['final_score']:<10.2f} "
                  f"{round_data['improvement']:<+10.2f} "
                  f"{converged_str:<12}")

        # Save report
        report = {
            "project": self.project,
            "timestamp": datetime.now().isoformat(),
            "rounds": self.rounds,
            "iterations_per_round": self.iterations_per_round,
            "initial_score": initial_score,
            "final_score": final_score,
            "total_improvement": total_improvement,
            "improvement_percentage": improvement_pct,
            "round_history": self.round_history
        }

        report_file = Path.cwd() / f"progressive_optimization_report_{int(datetime.now().timestamp())}.json"
        with open(report_file, 'w') as f:
            json.dump(report, f, indent=2)

        print(f"\n✅ Report saved: {report_file}")

        return report


def main():
    """Main entry point."""
    import argparse

    parser = argparse.ArgumentParser(description="Run progressive optimization rounds")
    parser.add_argument("--project", default="pr-faq-validator", help="Project name")
    parser.add_argument("--rounds", type=int, default=20, help="Number of rounds")
    parser.add_argument("--iterations", type=int, default=20, help="Iterations per round")

    args = parser.parse_args()

    optimizer = ProgressiveOptimizer(
        project=args.project,
        rounds=args.rounds,
        iterations_per_round=args.iterations
    )

    optimizer.run_progressive_optimization()


if __name__ == "__main__":
    main()

