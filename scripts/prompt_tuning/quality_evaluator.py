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
        """Evaluate a single test case result using comprehensive LLM evaluation."""
        generated_content = test_result.get("generated_content", {})
        press_release = generated_content.get("press_release", "")
        faq = generated_content.get("faq", "")

        if not press_release and not faq:
            return {
                "test_case_id": test_result.get("test_case_id"),
                "test_case_name": test_result.get("test_case_name"),
                "overall_score": 0.0,
                "press_release_score": {"score": 0, "clarity": 0, "structure": 0, "feedback": "No content generated"},
                "faq_score": {"score": 0, "completeness": 0, "feedback": "No content generated"},
                "timestamp": datetime.now().isoformat()
            }

        # Comprehensive evaluation prompt (bloginator-style)
        evaluation_prompt = f"""Evaluate the following PR-FAQ content comprehensively.

PRESS RELEASE:
{press_release}

FAQ:
{faq}

Provide a comprehensive evaluation with:
1. overall_score (0-5 scale, floating point)
2. press_release_quality (0-5 scale)
3. faq_completeness (0-5 scale)
4. clarity_score (0-5 scale)
5. structure_adherence (0-5 scale)
6. content_quality breakdown (clarity, depth, nuance, specificity - each 0-5)
7. slop_violations (critical, high, medium, low - arrays)
8. voice_analysis (authenticity_score, strengths, concerns)
9. feedback (detailed text)
10. strengths (array of strings)
11. improvements (array of strings)
12. evolutionary_strategy (prompt_to_modify, specific_changes, priority, expected_impact)

Return ONLY valid JSON matching this structure.
"""

        response = await self.evaluator_client.generate(evaluation_prompt, temperature=0.3)

        try:
            # Parse JSON response (may be wrapped in content field for file-based LLM)
            evaluation = json.loads(response)
        except json.JSONDecodeError:
            # Fallback scoring
            evaluation = {
                "overall_score": 3.5,
                "press_release_quality": 3.5,
                "faq_completeness": 3.5,
                "clarity_score": 3.5,
                "structure_adherence": 3.5,
                "feedback": "Evaluation parsing failed",
                "strengths": [],
                "improvements": []
            }

        # Convert 0-5 scores to 0-100 for compatibility
        overall_score = evaluation.get("overall_score", 0) * 20  # 5.0 -> 100

        return {
            "test_case_id": test_result.get("test_case_id"),
            "test_case_name": test_result.get("test_case_name"),
            "overall_score": overall_score,
            "press_release_score": {
                "score": evaluation.get("press_release_quality", 0) * 20,
                "clarity": evaluation.get("clarity_score", 0) * 20,
                "structure": evaluation.get("structure_adherence", 0) * 20,
                "feedback": evaluation.get("feedback", "")
            },
            "faq_score": {
                "score": evaluation.get("faq_completeness", 0) * 20,
                "completeness": evaluation.get("faq_completeness", 0) * 20,
                "feedback": evaluation.get("feedback", "")
            },
            "content_quality": evaluation.get("content_quality", {}),
            "slop_violations": evaluation.get("slop_violations", {}),
            "voice_analysis": evaluation.get("voice_analysis", {}),
            "strengths": evaluation.get("strengths", []),
            "improvements": evaluation.get("improvements", []),
            "evolutionary_strategy": evaluation.get("evolutionary_strategy", {}),
            "timestamp": datetime.now().isoformat()
        }
    

    
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

