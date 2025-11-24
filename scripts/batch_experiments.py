#!/usr/bin/env python3
"""
Batch Experiment Runner

Run multiple optimization experiments with different configurations.
Useful for comparing convergence detection, mutation strategies, etc.

Usage:
    python scripts/batch_experiments.py --experiments 3 --iterations 20
    python scripts/batch_experiments.py --config experiments_config.json
"""

import argparse
import json
import subprocess
import sys
import time
from pathlib import Path
from typing import List, Dict, Any


class BatchExperimentRunner:
    """Run multiple experiments in batch."""
    
    def __init__(self, output_dir: str = "batch_results"):
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(exist_ok=True)
        self.results: List[Dict[str, Any]] = []
    
    def run_experiment(self, config: Dict[str, Any]) -> Dict[str, Any]:
        """Run a single experiment with given configuration."""
        project = config.get("project", "pr-faq-validator")
        iterations = config.get("iterations", 20)
        use_real_api = config.get("use_real_api", False)
        enable_convergence = config.get("enable_convergence", True)
        
        print("\n" + "="*80)
        print(f"Running Experiment: {config.get('name', 'Unnamed')}")
        print("="*80)
        print(f"Project: {project}")
        print(f"Iterations: {iterations}")
        print(f"Real API: {use_real_api}")
        print(f"Convergence Detection: {enable_convergence}")
        print()
        
        start_time = time.time()
        
        # Build command
        cmd = [
            sys.executable,
            "scripts/run_autonomous_experiment.py",
            "--project", project,
            "--iterations", str(iterations)
        ]
        
        if use_real_api:
            cmd.append("--real-api")
        
        # Run experiment
        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=1200  # 20 minute timeout
            )
            
            elapsed = time.time() - start_time
            
            # Parse results
            results_file = Path(f"prompt_tuning_results_{project}/optimization_final_results.json")
            if results_file.exists():
                with open(results_file, 'r') as f:
                    experiment_results = json.load(f)
                
                return {
                    "name": config.get("name", "Unnamed"),
                    "config": config,
                    "success": result.returncode == 0,
                    "elapsed_seconds": elapsed,
                    "results": experiment_results,
                    "stdout": result.stdout,
                    "stderr": result.stderr
                }
            else:
                return {
                    "name": config.get("name", "Unnamed"),
                    "config": config,
                    "success": False,
                    "elapsed_seconds": elapsed,
                    "error": "Results file not found",
                    "stdout": result.stdout,
                    "stderr": result.stderr
                }
        except subprocess.TimeoutExpired:
            return {
                "name": config.get("name", "Unnamed"),
                "config": config,
                "success": False,
                "error": "Timeout after 20 minutes"
            }
        except Exception as e:
            return {
                "name": config.get("name", "Unnamed"),
                "config": config,
                "success": False,
                "error": str(e)
            }
    
    def run_batch(self, experiments: List[Dict[str, Any]]):
        """Run batch of experiments."""
        print("\n" + "="*80)
        print(f"BATCH EXPERIMENT RUNNER - {len(experiments)} experiments")
        print("="*80)
        
        for i, config in enumerate(experiments, 1):
            print(f"\n[{i}/{len(experiments)}] Starting experiment...")
            result = self.run_experiment(config)
            self.results.append(result)
            
            if result["success"]:
                print(f"✅ Experiment completed in {result['elapsed_seconds']:.1f}s")
            else:
                print(f"❌ Experiment failed: {result.get('error', 'Unknown error')}")
        
        # Save batch results
        self.save_batch_results()
        self.print_summary()
    
    def save_batch_results(self):
        """Save batch results to file."""
        timestamp = int(time.time())
        results_file = self.output_dir / f"batch_results_{timestamp}.json"
        
        with open(results_file, 'w') as f:
            json.dump({
                "timestamp": time.strftime('%Y-%m-%d %H:%M:%S'),
                "total_experiments": len(self.results),
                "successful": sum(1 for r in self.results if r["success"]),
                "failed": sum(1 for r in self.results if not r["success"]),
                "experiments": self.results
            }, f, indent=2)
        
        print(f"\n✅ Batch results saved to: {results_file}")
    
    def print_summary(self):
        """Print batch summary."""
        print("\n" + "="*80)
        print("BATCH SUMMARY")
        print("="*80)
        
        successful = [r for r in self.results if r["success"]]
        failed = [r for r in self.results if not r["success"]]
        
        print(f"\nTotal Experiments: {len(self.results)}")
        print(f"Successful: {len(successful)}")
        print(f"Failed: {len(failed)}")
        
        if successful:
            print("\n📊 Successful Experiments:")
            for r in successful:
                exp_results = r.get("results", {})
                baseline = exp_results.get("baseline_score", 0)
                final = exp_results.get("final_score", 0)
                improvement = exp_results.get("improvement", 0)
                convergence = exp_results.get("convergence", {})
                
                print(f"\n  {r['name']}:")
                print(f"    Baseline: {baseline:.2f}")
                print(f"    Final: {final:.2f}")
                print(f"    Improvement: +{improvement:.2f} ({improvement/baseline*100:.2f}%)")
                print(f"    Time: {r['elapsed_seconds']:.1f}s")
                
                if convergence.get("converged"):
                    print(f"    Converged: Iteration {convergence.get('convergence_iteration')}")
                    print(f"    Wasted: {convergence.get('wasted_iterations')} iterations")
        
        if failed:
            print("\n❌ Failed Experiments:")
            for r in failed:
                print(f"  {r['name']}: {r.get('error', 'Unknown error')}")
        
        print("\n" + "="*80)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Run batch optimization experiments")
    parser.add_argument("--experiments", type=int, default=1, help="Number of experiments to run")
    parser.add_argument("--iterations", type=int, default=20, help="Iterations per experiment")
    parser.add_argument("--config", type=str, help="JSON config file with experiment definitions")
    parser.add_argument("--output-dir", type=str, default="batch_results", help="Output directory")
    
    args = parser.parse_args()
    
    runner = BatchExperimentRunner(output_dir=args.output_dir)
    
    if args.config:
        # Load experiments from config file
        with open(args.config, 'r') as f:
            config_data = json.load(f)
        experiments = config_data.get("experiments", [])
    else:
        # Generate default experiments
        experiments = [
            {
                "name": f"Experiment {i+1}",
                "project": "pr-faq-validator",
                "iterations": args.iterations,
                "use_real_api": False,
                "enable_convergence": True
            }
            for i in range(args.experiments)
        ]
    
    runner.run_batch(experiments)


if __name__ == "__main__":
    main()

