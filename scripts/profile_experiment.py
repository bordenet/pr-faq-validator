#!/usr/bin/env python3
"""
Profile autonomous experiment performance.

Measures timing, memory usage, and system metrics.
"""

import argparse
import json
import subprocess
import sys
import time
from pathlib import Path
from typing import Any, Dict


class ExperimentProfiler:
    """Profile experiment performance."""

    def __init__(self):
        self.metrics: Dict[str, Any] = {
            "start_time": None,
            "end_time": None,
            "duration_seconds": 0,
            "requests_processed": 0,
            "iterations_completed": 0,
            "convergence_iteration": None,
            "early_stop_savings": 0,
            "baseline_score": 0,
            "final_score": 0,
            "improvement": 0,
        }

    def profile_experiment(self, project: str, iterations: int) -> Dict[str, Any]:
        """Profile a single experiment."""
        print("=" * 80)
        print("EXPERIMENT PROFILER")
        print("=" * 80)
        print(f"Project: {project}")
        print(f"Iterations: {iterations}")
        print()

        # Record start metrics
        self.metrics["start_time"] = time.time()

        # Run experiment
        print("Starting experiment...")
        start = time.time()

        result = subprocess.run(
            [
                sys.executable,
                "scripts/run_autonomous_experiment.py",
                "--project",
                project,
                "--iterations",
                str(iterations),
            ],
            capture_output=True,
            text=True,
            check=False,
        )

        elapsed = time.time() - start

        # Record end metrics
        self.metrics["end_time"] = time.time()
        self.metrics["duration_seconds"] = elapsed

        # Parse results
        results_file = Path(f"prompt_tuning_results_{project}/optimization_final_results.json")
        if results_file.exists():
            with open(results_file, "r", encoding="utf-8") as f:
                exp_results = json.load(f)

            self.metrics["iterations_completed"] = len(exp_results.get("iteration_history", []))
            self.metrics["baseline_score"] = exp_results.get("baseline_score", 0)
            self.metrics["final_score"] = exp_results.get("final_score", 0)
            self.metrics["improvement"] = exp_results.get("improvement", 0)

            convergence = exp_results.get("convergence", {})
            self.metrics["convergence_iteration"] = convergence.get("convergence_iteration")
            self.metrics["early_stop_savings"] = convergence.get("wasted_iterations", 0)

            # Count requests (2 per iteration: simulation + evaluation)
            self.metrics["requests_processed"] = self.metrics["iterations_completed"] * 2 + 2  # +2 for baseline

        # Count LLM requests from files
        llm_response_dir = Path(".pr-faq-validator/llm_responses")
        if llm_response_dir.exists():
            self.metrics["requests_processed"] = len(list(llm_response_dir.glob("*.json")))

        self.metrics["success"] = result.returncode == 0
        self.metrics["stdout"] = result.stdout
        self.metrics["stderr"] = result.stderr

        return self.metrics

    def print_report(self):
        """Print profiling report."""
        print("\n" + "=" * 80)
        print("PROFILING REPORT")
        print("=" * 80)

        print("\n⏱️  TIMING METRICS")
        print(f"  Total Duration:        {self.metrics['duration_seconds']:.2f} seconds")
        print(f"  Iterations Completed:  {self.metrics['iterations_completed']}")
        if self.metrics["iterations_completed"] > 0:
            secs_per_iter = self.metrics["duration_seconds"] / self.metrics["iterations_completed"]
            print(f"  Seconds per Iteration: {secs_per_iter:.2f}s")

        print("\n📈 QUALITY METRICS")
        print(f"  Baseline Score:        {self.metrics['baseline_score']:.2f}")
        print(f"  Final Score:           {self.metrics['final_score']:.2f}")
        print(f"  Improvement:           +{self.metrics['improvement']:.2f} points")
        if self.metrics["baseline_score"] > 0:
            improvement_pct = (self.metrics["improvement"] / self.metrics["baseline_score"]) * 100
            print(f"  Improvement %:         {improvement_pct:.2f}%")

        if self.metrics["convergence_iteration"]:
            print("\n🎯 CONVERGENCE METRICS")
            print(f"  Convergence Iteration: {self.metrics['convergence_iteration']}")
            print(f"  Early Stop Savings:    {self.metrics['early_stop_savings']} iterations")
            if self.metrics["iterations_completed"] > 0:
                savings_percent = (self.metrics["early_stop_savings"] / self.metrics["iterations_completed"]) * 100
                print(f"  Efficiency Gain:       {savings_percent:.1f}%")

        print("\n📊 REQUEST METRICS")
        print(f"  Total Requests:        {self.metrics['requests_processed']}")
        if self.metrics["duration_seconds"] > 0:
            print(
                f"  Requests per Second:   {self.metrics['requests_processed'] / self.metrics['duration_seconds']:.2f}"
            )

        print("\n✅ SUCCESS")
        print(f"  Experiment Succeeded:  {self.metrics['success']}")

        print("\n" + "=" * 80)

    def save_report(self, output_file: str):
        """Save profiling report to JSON."""
        with open(output_file, "w", encoding="utf-8") as f:
            json.dump(self.metrics, f, indent=2)
        print(f"\n✅ Profiling report saved to: {output_file}")


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Profile autonomous experiment")
    parser.add_argument("--project", type=str, default="pr-faq-validator", help="Project name")
    parser.add_argument("--iterations", type=int, default=10, help="Number of iterations")
    parser.add_argument("--output", type=str, default="profiling_report.json", help="Output file")

    args = parser.parse_args()

    profiler = ExperimentProfiler()
    profiler.profile_experiment(args.project, args.iterations)
    profiler.print_report()
    profiler.save_report(args.output)


if __name__ == "__main__":
    main()
