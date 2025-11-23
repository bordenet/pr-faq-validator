"""Evolutionary prompt tuner for PR-FAQ Validator."""

import json
import asyncio
from pathlib import Path
from typing import Dict, List, Any, Optional
from datetime import datetime

from prompt_tuning_config import PromptTuningConfig, load_prompts, save_prompts
from llm_client import create_llm_client
from prompt_simulator import PromptSimulator
from quality_evaluator import QualityEvaluator


class EvolutionaryTuner:
    """Evolutionary optimization for PR-FAQ prompts."""
    
    def __init__(self, config: PromptTuningConfig):
        self.config = config
        self.mutation_client = create_llm_client(
            provider=config.llm_provider,
            model=config.llm_model,
            mock=config.mock_mode
        )
        self.simulator = PromptSimulator(config)
        self.evaluator = QualityEvaluator(config)
        
        self.best_score = 0.0
        self.best_prompts = {}
        self.iteration_history = []
    
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
        baseline_results = await self._evaluate_prompts(current_prompts, iteration=0)
        self.best_score = baseline_results["aggregate_scores"]["overall"]
        self.best_prompts = current_prompts.copy()
        
        print(f"Baseline score: {self.best_score:.2f}")
        
        # Evolutionary loop
        for iteration in range(1, max_iterations + 1):
            print(f"\n=== Iteration {iteration}/{max_iterations} ===")
            
            # Mutate prompts
            mutated_prompts = await self._mutate_prompts(current_prompts, iteration)
            
            # Evaluate mutated prompts
            results = await self._evaluate_prompts(mutated_prompts, iteration)
            current_score = results["aggregate_scores"]["overall"]
            
            print(f"Current score: {current_score:.2f}")
            
            # Keep or discard mutation
            if current_score > self.best_score:
                print(f"✓ Improvement! {self.best_score:.2f} → {current_score:.2f}")
                self.best_score = current_score
                self.best_prompts = mutated_prompts.copy()
                current_prompts = mutated_prompts
                
                # Save improved prompts
                save_prompts(self.config, self.best_prompts)
            else:
                print(f"✗ No improvement. Keeping previous prompts.")
            
            # Record iteration
            self.iteration_history.append({
                "iteration": iteration,
                "score": current_score,
                "best_score": self.best_score,
                "improved": current_score > self.best_score
            })
        
        # Save final results
        final_results = {
            "project": self.config.project_name,
            "timestamp": datetime.now().isoformat(),
            "max_iterations": max_iterations,
            "baseline_score": baseline_results["aggregate_scores"]["overall"],
            "final_score": self.best_score,
            "improvement": self.best_score - baseline_results["aggregate_scores"]["overall"],
            "iteration_history": self.iteration_history,
            "best_prompts": self.best_prompts
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
            # Restore original prompts if evaluation failed
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
1. Generate higher quality press releases following Amazon's PR-FAQ format
2. Create more comprehensive and useful FAQ sections
3. Better address stakeholder questions and concerns
4. Maintain clarity and structure

Provide ONLY the improved prompt text, without any explanation or meta-commentary.
"""
            
            mutated_content = await self.mutation_client.generate(
                mutation_prompt,
                temperature=0.8
            )
            
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
        
        with open(output_file, 'w') as f:
            json.dump(results, f, indent=2)
        
        return output_file

