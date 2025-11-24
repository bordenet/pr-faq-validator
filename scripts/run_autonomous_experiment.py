#!/usr/bin/env python3
"""
Autonomous LLM Prompt Optimization Experiment Runner

This script runs a complete optimization experiment autonomously:
1. Starts auto-responder in background
2. Runs optimization experiment
3. Analyzes results
4. Generates comprehensive report
5. Cleans up processes

Usage:
    python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20
    python scripts/run_autonomous_experiment.py --project my-project --iterations 10 --real-api
"""

import argparse
import asyncio
import json
import os
import signal
import subprocess
import sys
import time
from pathlib import Path
from typing import Optional, Dict, Any


class AutonomousExperiment:
    """Manages autonomous optimization experiments."""
    
    def __init__(self, project: str, max_iterations: int, use_real_api: bool = False):
        self.project = project
        self.max_iterations = max_iterations
        self.use_real_api = use_real_api
        self.auto_responder_process: Optional[subprocess.Popen] = None
        self.results_dir = Path(f"prompt_tuning_results_{project}")
        
    def start_auto_responder(self) -> bool:
        """Start auto-responder in background."""
        if self.use_real_api:
            print("Using real API - auto-responder not needed")
            return True

        print("Starting auto-responder in background...")
        try:
            # Ensure log directory exists
            Path(".pr-faq-validator").mkdir(parents=True, exist_ok=True)

            # Write logs to files for debugging
            stdout_log = open(".pr-faq-validator/auto_responder_stdout.log", "w")
            stderr_log = open(".pr-faq-validator/auto_responder_stderr.log", "w")

            self.auto_responder_process = subprocess.Popen(
                [sys.executable, "-u", "scripts/auto_respond_llm.py", "--continuous", "--interval", "0.2"],
                stdout=stdout_log,
                stderr=stderr_log,
                text=True,
                bufsize=1  # Line buffered
            )
            time.sleep(2)  # Give it time to start

            if self.auto_responder_process.poll() is None:
                print("✅ Auto-responder started successfully")
                print("   Logs: .pr-faq-validator/auto_responder_stdout.log")
                return True
            else:
                print("❌ Auto-responder failed to start")
                # Read error logs
                stderr_log.flush()
                with open(".pr-faq-validator/auto_responder_stderr.log", "r") as f:
                    error_output = f.read()
                if error_output:
                    print(f"   Error output: {error_output}")
                return False
        except Exception as e:
            print(f"❌ Error starting auto-responder: {e}")
            return False
    
    def stop_auto_responder(self):
        """Stop auto-responder process."""
        if self.auto_responder_process:
            print("Stopping auto-responder...")
            self.auto_responder_process.terminate()
            try:
                self.auto_responder_process.wait(timeout=5)
                print("✅ Auto-responder stopped")
            except subprocess.TimeoutExpired:
                self.auto_responder_process.kill()
                print("⚠️  Auto-responder killed (timeout)")
    
    def initialize_project(self) -> bool:
        """Initialize project if needed."""
        # Check if project is already initialized
        if self.results_dir.exists():
            test_cases_file = self.results_dir / f"test_cases_{self.project}.json"
            if test_cases_file.exists():
                print(f"✅ Project already initialized: {self.project}")
                return True

        print(f"Initializing project: {self.project}...")
        try:
            result = subprocess.run(
                [sys.executable, "prompt_tuning_tool.py", "init", self.project],
                capture_output=True,
                text=True,
                timeout=30
            )

            if result.returncode == 0:
                print(f"✅ Project initialized successfully")
                return True
            else:
                print(f"❌ Project initialization failed")
                print(result.stderr)
                return False
        except Exception as e:
            print(f"❌ Error initializing project: {e}")
            return False

    def run_optimization(self) -> bool:
        """Run optimization experiment."""
        print(f"\nRunning {self.max_iterations}-iteration optimization experiment...")
        print(f"Project: {self.project}")
        print(f"Mode: {'Real API' if self.use_real_api else 'File-based LLM'}")
        print()

        env = os.environ.copy()
        if not self.use_real_api:
            env["LLM_PROVIDER"] = "assistant"

        cmd = [
            sys.executable,
            "prompt_tuning_tool.py",
            "evolve",
            self.project,
            "--max-iterations",
            str(self.max_iterations)
        ]
        
        try:
            result = subprocess.run(
                cmd,
                env=env,
                capture_output=True,
                text=True,
                timeout=600  # 10 minute timeout
            )
            
            print(result.stdout)
            if result.stderr:
                print("STDERR:", result.stderr, file=sys.stderr)
            
            if result.returncode == 0:
                print("✅ Optimization completed successfully")
                return True
            else:
                print(f"❌ Optimization failed with code {result.returncode}")
                return False
        except subprocess.TimeoutExpired:
            print("❌ Optimization timed out after 10 minutes")
            return False
        except Exception as e:
            print(f"❌ Error running optimization: {e}")
            return False
    
    def analyze_results(self) -> Optional[Dict[str, Any]]:
        """Analyze experiment results."""
        print("\nAnalyzing results...")
        
        try:
            result = subprocess.run(
                [sys.executable, "scripts/analyze_experiment.py", str(self.results_dir)],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            print(result.stdout)
            if result.stderr:
                print("STDERR:", result.stderr, file=sys.stderr)
            
            # Load results for return
            results_file = self.results_dir / "optimization_final_results.json"
            if results_file.exists():
                with open(results_file, 'r') as f:
                    return json.load(f)
            else:
                print("⚠️  Results file not found")
                return None
        except Exception as e:
            print(f"❌ Error analyzing results: {e}")
            return None
    
    def generate_report(self, results: Dict[str, Any]) -> bool:
        """Generate comprehensive experiment report."""
        print("\nGenerating experiment report...")
        
        try:
            report_path = Path(f"docs/EXPERIMENT_{self.project}_{int(time.time())}.md")
            
            # Extract key metrics
            baseline = results.get('baseline_score', 0)
            final = results.get('final_score', 0)
            improvement = results.get('improvement', 0)
            iterations = results.get('max_iterations', 0)
            
            report_content = f"""# Autonomous Optimization Experiment Report

**Project:** {self.project}  
**Date:** {time.strftime('%Y-%m-%d %H:%M:%S')}  
**Mode:** {'Real API' if self.use_real_api else 'File-based LLM'}

## Results Summary

- **Baseline Score:** {baseline:.2f}
- **Final Score:** {final:.2f}
- **Improvement:** +{improvement:.2f} points ({improvement/baseline*100:.2f}%)
- **Iterations:** {iterations}

## Iteration History

"""
            
            # Add iteration details
            for it in results.get('iteration_history', []):
                status = "✓" if it['score'] == it['best_score'] else "✗"
                report_content += f"- Iteration {it['iteration']:2d}: {it['score']:.2f} {status}\n"
            
            report_content += f"""

## Files Generated

- Results directory: `{self.results_dir}/`
- Total files: {len(list(self.results_dir.glob('*.json')))} JSON files

## Next Steps

1. Review detailed evaluation files in `{self.results_dir}/`
2. Examine best prompts in `prompts_{self.project}.json`
3. Consider running with more iterations if not converged
4. Implement recommended improvements from evolutionary strategy

---

*Generated automatically by autonomous experiment runner*
"""
            
            with open(report_path, 'w') as f:
                f.write(report_content)
            
            print(f"✅ Report generated: {report_path}")
            return True
        except Exception as e:
            print(f"❌ Error generating report: {e}")
            return False
    
    def cleanup(self):
        """Clean up resources."""
        self.stop_auto_responder()
    
    def run(self) -> bool:
        """Run complete autonomous experiment."""
        print("="*80)
        print("AUTONOMOUS LLM PROMPT OPTIMIZATION EXPERIMENT")
        print("="*80)
        print()

        try:
            # Step 1: Initialize project
            if not self.initialize_project():
                return False

            # Step 2: Start auto-responder
            if not self.start_auto_responder():
                return False

            # Step 3: Run optimization
            if not self.run_optimization():
                return False
            
            # Step 4: Analyze results
            results = self.analyze_results()
            if not results:
                return False

            # Step 5: Generate report
            if not self.generate_report(results):
                return False
            
            print()
            print("="*80)
            print("✅ AUTONOMOUS EXPERIMENT COMPLETED SUCCESSFULLY")
            print("="*80)
            return True
            
        except KeyboardInterrupt:
            print("\n⚠️  Experiment interrupted by user")
            return False
        except Exception as e:
            print(f"\n❌ Unexpected error: {e}")
            return False
        finally:
            self.cleanup()


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Run autonomous LLM prompt optimization experiment"
    )
    parser.add_argument(
        "--project",
        type=str,
        default="pr-faq-validator",
        help="Project name (default: pr-faq-validator)"
    )
    parser.add_argument(
        "--iterations",
        type=int,
        default=20,
        help="Maximum iterations (default: 20)"
    )
    parser.add_argument(
        "--real-api",
        action="store_true",
        help="Use real Anthropic API instead of file-based LLM"
    )
    
    args = parser.parse_args()
    
    # Validate API key if using real API
    if args.real_api and not os.getenv("ANTHROPIC_API_KEY"):
        print("❌ Error: ANTHROPIC_API_KEY not set")
        print("Set it with: export ANTHROPIC_API_KEY=your_key_here")
        sys.exit(1)
    
    # Run experiment
    experiment = AutonomousExperiment(
        project=args.project,
        max_iterations=args.iterations,
        use_real_api=args.real_api
    )
    
    success = experiment.run()
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()

