"""Evolutionary prompt tuner for PR-FAQ Validator."""

import json
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional

try:
    from convergence_detector import ConvergenceConfig, ConvergenceDetector
    from llm_client import create_llm_client
    from prompt_simulator import PromptSimulator
    from prompt_tuning_config import PromptTuningConfig, load_prompts, save_prompts
    from quality_evaluator import QualityEvaluator
except ImportError:
    from scripts.prompt_tuning.convergence_detector import ConvergenceConfig, ConvergenceDetector
    from scripts.prompt_tuning.llm_client import create_llm_client
    from scripts.prompt_tuning.prompt_simulator import PromptSimulator
    from scripts.prompt_tuning.prompt_tuning_config import (
        PromptTuningConfig,
        load_prompts,
        save_prompts,
    )
    from scripts.prompt_tuning.quality_evaluator import QualityEvaluator


class EvolutionaryTuner:
    """Evolutionary optimization for PR-FAQ prompts."""

    def __init__(self, config: PromptTuningConfig, enable_convergence_detection: bool = True):
        self.config = config
        self.mutation_client = create_llm_client(
            provider=config.llm_provider, model=config.llm_model, mock=config.mock_mode
        )
        self.simulator = PromptSimulator(config)
        self.evaluator = QualityEvaluator(config)

        self.best_score = 0.0
        self.best_prompts: Dict[str, str] = {}
        self.iteration_history: List[Dict[str, Any]] = []

        # Convergence detection
        self.convergence_detector = (
            ConvergenceDetector(
                ConvergenceConfig(
                    no_improvement_threshold=5,
                    min_improvement_percent=0.1,
                    enable_early_stop=enable_convergence_detection,
                )
            )
            if enable_convergence_detection
            else None
        )

    async def run_evolution(self, max_iterations: Optional[int] = None) -> Dict[str, Any]:
        """Run evolutionary optimization."""
        if max_iterations is None:
            max_iterations = self.config.max_iterations

        # Load initial prompts
        current_prompts = load_prompts(self.config)
        if not current_prompts:
            current_prompts = {"pr_faq_generation": self._get_default_prompt()}
            save_prompts(self.config, current_prompts)

        # Run baseline evaluation
        print(f"[DEBUG] Before baseline eval - current_prompts keys: {list(current_prompts.keys())}")
        print(f"[DEBUG] Before baseline eval - prompt length: {len(current_prompts.get('pr_faq_generation', ''))}")
        print(
            f"[DEBUG] Before baseline eval - prompt starts with: {current_prompts.get('pr_faq_generation', '')[:100]}"
        )

        baseline_results = await self._evaluate_prompts(current_prompts, iteration=0)
        self.best_score = baseline_results["aggregate_scores"]["overall"]
        self.best_prompts = current_prompts.copy()

        print(
            f"[DEBUG] After baseline eval - current_prompts length: {len(current_prompts.get('pr_faq_generation', ''))}"
        )
        print(f"[DEBUG] After baseline eval - prompt starts with: {current_prompts.get('pr_faq_generation', '')[:100]}")
        print(
            f"[DEBUG] After baseline eval - best_prompts length: {len(self.best_prompts.get('pr_faq_generation', ''))}"
        )
        print(f"Baseline score: {self.best_score:.2f}")

        # Evolutionary loop
        for iteration in range(1, max_iterations + 1):
            print(f"\n=== Iteration {iteration}/{max_iterations} ===")

            # Mutate prompts
            prompt_len = len(current_prompts.get("pr_faq_generation", ""))
            prompt_start = current_prompts.get("pr_faq_generation", "")[:100]
            print(f"[DEBUG] Iteration {iteration} - Before mutation - length: {prompt_len}")
            print(f"[DEBUG] Iteration {iteration} - Before mutation - starts with: {prompt_start}")

            mutated_prompts = await self._mutate_prompts(current_prompts, iteration)

            mutated_len = len(mutated_prompts.get("pr_faq_generation", ""))
            mutated_start = mutated_prompts.get("pr_faq_generation", "")[:100]
            print(f"[DEBUG] Iteration {iteration} - After mutation - length: {mutated_len}")
            print(f"[DEBUG] Iteration {iteration} - After mutation - starts with: {mutated_start}")

            # Evaluate mutated prompts
            results = await self._evaluate_prompts(mutated_prompts, iteration)
            current_score = results["aggregate_scores"]["overall"]

            print(f"Current score: {current_score:.2f}")

            # Keep or discard mutation
            improved = current_score > self.best_score
            if improved:
                print(f"✓ Improvement! {self.best_score:.2f} → {current_score:.2f}")
                self.best_score = current_score
                self.best_prompts = mutated_prompts.copy()
                current_prompts = mutated_prompts

                # Save improved prompts
                save_prompts(self.config, self.best_prompts)
            else:
                print("✗ No improvement. Keeping previous prompts.")

            # Record iteration
            self.iteration_history.append(
                {"iteration": iteration, "score": current_score, "best_score": self.best_score, "improved": improved}
            )

            # Update convergence detector
            if self.convergence_detector:
                convergence_status = self.convergence_detector.update(iteration, current_score, improved)

                # Check for early stopping
                if self.convergence_detector.should_stop():
                    print(f"\n🎯 Convergence detected at iteration {convergence_status['convergence_iteration']}")
                    print(f"No improvement for {convergence_status['no_improvement_count']} consecutive iterations.")
                    print(f"Stopping early. Saved {max_iterations - iteration} iterations.")
                    break

        # Print convergence summary
        if self.convergence_detector:
            self.convergence_detector.print_summary()

        # Save final results
        convergence_status = self.convergence_detector.get_status() if self.convergence_detector else {}

        final_results = {
            "project": self.config.project_name,
            "timestamp": datetime.now().isoformat(),
            "max_iterations": max_iterations,
            "baseline_score": baseline_results["aggregate_scores"]["overall"],
            "final_score": self.best_score,
            "improvement": self.best_score - baseline_results["aggregate_scores"]["overall"],
            "iteration_history": self.iteration_history,
            "best_prompts": self.best_prompts,
            "convergence": convergence_status,
            "recommendations": self.convergence_detector.get_recommendations() if self.convergence_detector else [],
        }

        self._save_final_results(final_results)

        return final_results

    async def _evaluate_prompts(self, prompts: Dict[str, str], iteration: int) -> Dict[str, Any]:
        """Evaluate a set of prompts."""
        # Temporarily save prompts for simulation
        original_prompts = load_prompts(self.config)
        save_prompts(self.config, prompts)

        try:
            # Run simulation
            simulation_results = await self.simulator.run_simulation(iteration)

            # Evaluate results
            evaluation_results = await self.evaluator.evaluate_results(simulation_results)

            # Save results
            self.simulator.save_results(simulation_results, iteration)
            self.evaluator.save_evaluation(evaluation_results, iteration)

            return evaluation_results
        finally:
            # Restore original prompts after evaluation (will be overwritten if this was an improvement)
            if iteration > 0:
                save_prompts(self.config, original_prompts)

    async def _mutate_prompts(self, prompts: Dict[str, str], iteration: int) -> Dict[str, str]:
        """Generate mutated version of prompts."""
        mutated = {}

        for prompt_name, prompt_content in prompts.items():
            mutation_prompt = f"""You are an expert at optimizing LLM prompts for PR-FAQ generation.

Current prompt:
{prompt_content}

Based on iteration {iteration}, suggest an improved version of this prompt that will:
1. Generate higher quality press releases following Amazon's format
2. Create more comprehensive and useful FAQ sections
3. Better address stakeholder questions and concerns
4. Maintain clarity and structure

Provide ONLY the improved prompt text, without any explanation or meta-commentary.
"""

            print(f"[DEBUG-TUNER] Mutation prompt length: {len(mutation_prompt)}")
            print(f"[DEBUG-TUNER] Mutation prompt first 150 chars: {mutation_prompt[:150]}")
            print(f"[DEBUG-TUNER] Has 'Current prompt:': {'Current prompt:' in mutation_prompt}")

            mutated_content = await self.mutation_client.generate(mutation_prompt, temperature=0.8)

            print(f"[DEBUG-TUNER] Mutated content length: {len(mutated_content)}")
            print(f"[DEBUG-TUNER] Mutated content first 150 chars: {mutated_content[:150]}")

            mutated[prompt_name] = mutated_content.strip()

        return mutated

    def _get_default_prompt(self) -> str:
        """Get default PR-FAQ generation prompt."""
        return """Generate a comprehensive PR-FAQ document following Amazon's format.

Project: {projectName}
Problem: {problemDescription}
Context: {businessContext}

Create a press release and FAQ that clearly articulates the customer problem, solution, and benefits."""

    def _save_final_results(self, results: Dict[str, Any]) -> Path:
        """Save final optimization results."""
        output_file = self.config.results_dir / "optimization_final_results.json"

        with open(output_file, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2)

        return output_file
