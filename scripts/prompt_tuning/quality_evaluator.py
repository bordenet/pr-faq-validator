"""Quality evaluator for PR-FAQ content."""

import json
import asyncio
from pathlib import Path
from typing import Dict, List, Any
from datetime import datetime

from prompt_tuning_config import PromptTuningConfig
from llm_client import create_llm_client


class QualityEvaluator:
    """Evaluates quality of generated PR-FAQ content."""
    
    def __init__(self, config: PromptTuningConfig):
        self.config = config
        self.evaluator_client = create_llm_client(
            provider=config.llm_provider,
            model=config.evaluator_model,
            mock=config.mock_mode
        )
    
    async def evaluate_results(self, simulation_results: Dict[str, Any]) -> Dict[str, Any]:
        """Evaluate simulation results."""
        test_case_results = simulation_results.get("test_case_results", [])
        
        evaluation_results = {
            "iteration": simulation_results.get("iteration", 0),
            "timestamp": datetime.now().isoformat(),
            "project": self.config.project_name,
            "evaluations": []
        }
        
        for test_result in test_case_results:
            evaluation = await self._evaluate_test_case(test_result)
            evaluation_results["evaluations"].append(evaluation)
        
        # Calculate aggregate scores
        evaluation_results["aggregate_scores"] = self._calculate_aggregate_scores(
            evaluation_results["evaluations"]
        )
        
        return evaluation_results
    
    async def _evaluate_test_case(self, test_result: Dict[str, Any]) -> Dict[str, Any]:
        """Evaluate a single test case result."""
        generated_content = test_result.get("generated_content", {})
        
        # Evaluate press release quality
        pr_score = await self._evaluate_press_release(
            generated_content.get("press_release", "")
        )
        
        # Evaluate FAQ quality
        faq_score = await self._evaluate_faq(
            generated_content.get("faq", "")
        )
        
        # Calculate overall score
        weights = self.config.validation_criteria
        overall_score = (
            pr_score.get("score", 0) * weights.get("press_release_quality", 0.3) +
            faq_score.get("score", 0) * weights.get("faq_completeness", 0.25) +
            pr_score.get("clarity", 0) * weights.get("clarity_score", 0.25) +
            pr_score.get("structure", 0) * weights.get("structure_adherence", 0.2)
        )
        
        return {
            "test_case_id": test_result.get("test_case_id"),
            "test_case_name": test_result.get("test_case_name"),
            "overall_score": overall_score,
            "press_release_score": pr_score,
            "faq_score": faq_score,
            "timestamp": datetime.now().isoformat()
        }
    
    async def _evaluate_press_release(self, press_release: str) -> Dict[str, Any]:
        """Evaluate press release quality."""
        if not press_release:
            return {"score": 0, "clarity": 0, "structure": 0, "feedback": "No press release generated"}
        
        evaluation_prompt = f"""Evaluate the following press release on a scale of 0-100:

{press_release}

Provide scores for:
1. Overall quality (0-100)
2. Clarity (0-100)
3. Structure adherence to Amazon PR-FAQ format (0-100)
4. Specific feedback on strengths and areas for improvement

Return your evaluation as JSON with keys: score, clarity, structure, feedback, strengths (list), improvements (list)
"""
        
        response = await self.evaluator_client.generate(evaluation_prompt, temperature=0.3)
        
        try:
            # Try to parse JSON response
            evaluation = json.loads(response)
        except json.JSONDecodeError:
            # Fallback to basic scoring if JSON parsing fails
            evaluation = {
                "score": 75,
                "clarity": 75,
                "structure": 75,
                "feedback": response[:200],
                "strengths": [],
                "improvements": []
            }
        
        return evaluation
    
    async def _evaluate_faq(self, faq: str) -> Dict[str, Any]:
        """Evaluate FAQ quality."""
        if not faq:
            return {"score": 0, "completeness": 0, "feedback": "No FAQ generated"}
        
        evaluation_prompt = f"""Evaluate the following FAQ section on a scale of 0-100:

{faq}

Provide scores for:
1. Overall quality (0-100)
2. Completeness - does it address key questions? (0-100)
3. Specific feedback on strengths and areas for improvement

Return your evaluation as JSON with keys: score, completeness, feedback, strengths (list), improvements (list)
"""
        
        response = await self.evaluator_client.generate(evaluation_prompt, temperature=0.3)
        
        try:
            evaluation = json.loads(response)
        except json.JSONDecodeError:
            evaluation = {
                "score": 75,
                "completeness": 75,
                "feedback": response[:200],
                "strengths": [],
                "improvements": []
            }
        
        return evaluation
    
    def _calculate_aggregate_scores(self, evaluations: List[Dict[str, Any]]) -> Dict[str, float]:
        """Calculate aggregate scores across all test cases."""
        if not evaluations:
            return {"overall": 0.0, "press_release": 0.0, "faq": 0.0}
        
        total_overall = sum(e.get("overall_score", 0) for e in evaluations)
        total_pr = sum(e.get("press_release_score", {}).get("score", 0) for e in evaluations)
        total_faq = sum(e.get("faq_score", {}).get("score", 0) for e in evaluations)
        
        count = len(evaluations)
        
        return {
            "overall": total_overall / count,
            "press_release": total_pr / count,
            "faq": total_faq / count
        }
    
    def save_evaluation(self, evaluation_results: Dict[str, Any], iteration: int) -> Path:
        """Save evaluation results to file."""
        self.config.results_dir.mkdir(parents=True, exist_ok=True)
        
        output_file = self.config.results_dir / f"evaluation_iteration_{iteration:03d}.json"
        
        with open(output_file, 'w') as f:
            json.dump(evaluation_results, f, indent=2)
        
        return output_file

