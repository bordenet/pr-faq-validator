"""End-to-end integration tests."""

import pytest
import json
import tempfile
from pathlib import Path
from prompt_tuning_config import load_project_config, save_test_cases
from evolutionary_tuner import EvolutionaryTuner


@pytest.mark.integration
class TestEndToEndWorkflow:
    """End-to-end workflow tests."""

    async def test_complete_optimization_cycle(self):
        """Test a complete optimization cycle from start to finish."""
        with tempfile.TemporaryDirectory() as tmpdir:
            # Setup
            config = load_project_config("e2e-test", base_dir=Path(tmpdir))
            config.mock_mode = True
            config.max_iterations = 3
            
            # Create test cases
            test_data = {
                "test_cases": [
                    {
                        "id": "test1",
                        "name": "E2E Test Case",
                        "industry": "Technology",
                        "project_type": "Product Launch",
                        "scope": "Enterprise",
                        "inputs": {
                            "projectName": "Revolutionary Platform",
                            "problemDescription": "Customers struggle with fragmented workflows",
                            "businessContext": "Enterprise SaaS market opportunity"
                        }
                    },
                    {
                        "id": "test2",
                        "name": "Second Test Case",
                        "industry": "Healthcare",
                        "project_type": "Service Launch",
                        "scope": "SMB",
                        "inputs": {
                            "projectName": "Healthcare Solution",
                            "problemDescription": "Patient data management challenges",
                            "businessContext": "Growing healthcare IT market"
                        }
                    }
                ]
            }
            save_test_cases(config, test_data)
            
            # Run evolution
            tuner = EvolutionaryTuner(config)
            results = await tuner.run_evolution()
            
            # Verify results structure
            assert "project" in results
            assert results["project"] == "e2e-test"
            assert "timestamp" in results
            assert "max_iterations" in results
            assert results["max_iterations"] == 3
            assert "baseline_score" in results
            assert "final_score" in results
            assert "improvement" in results
            assert "iteration_history" in results
            assert "best_prompts" in results
            
            # Verify iteration history
            assert len(results["iteration_history"]) == 3
            for iteration in results["iteration_history"]:
                assert "iteration" in iteration
                assert "score" in iteration
                assert "best_score" in iteration
                assert "improved" in iteration
            
            # Verify scores are reasonable
            assert 0 <= results["baseline_score"] <= 100
            assert 0 <= results["final_score"] <= 100
            
            # Verify prompts exist
            assert len(results["best_prompts"]) > 0
            
            # Verify files were created
            assert config.results_dir.exists()
            
            # Check for simulation files
            simulation_files = list(config.results_dir.glob("simulation_*.json"))
            assert len(simulation_files) > 0
            
            # Check for evaluation files
            evaluation_files = list(config.results_dir.glob("evaluation_*.json"))
            assert len(evaluation_files) > 0
            
            # Check for final results file
            final_results_file = config.results_dir / "optimization_final_results.json"
            assert final_results_file.exists()
            
            # Verify final results file content
            with open(final_results_file) as f:
                saved_results = json.load(f)
            
            assert saved_results["project"] == "e2e-test"
            assert saved_results["final_score"] == results["final_score"]

    async def test_multiple_iterations_improve_or_maintain(self):
        """Test that multiple iterations either improve or maintain score."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("multi-iter-test", base_dir=Path(tmpdir))
            config.mock_mode = True
            config.max_iterations = 5
            
            test_data = {
                "test_cases": [
                    {
                        "id": "test1",
                        "name": "Test",
                        "inputs": {
                            "projectName": "Test Project",
                            "problemDescription": "Test problem",
                            "businessContext": "Test context"
                        }
                    }
                ]
            }
            save_test_cases(config, test_data)
            
            tuner = EvolutionaryTuner(config)
            results = await tuner.run_evolution()
            
            # Final score should be >= baseline (evolution should never regress)
            assert results["final_score"] >= results["baseline_score"]
            
            # Best score should be monotonically non-decreasing
            best_scores = [iter["best_score"] for iter in results["iteration_history"]]
            for i in range(1, len(best_scores)):
                assert best_scores[i] >= best_scores[i-1]

    async def test_prompts_are_saved_and_loadable(self):
        """Test that prompts are saved and can be loaded."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("prompt-save-test", base_dir=Path(tmpdir))
            config.mock_mode = True
            config.max_iterations = 2
            
            test_data = {
                "test_cases": [
                    {
                        "id": "test1",
                        "name": "Test",
                        "inputs": {
                            "projectName": "Test",
                            "problemDescription": "Problem",
                            "businessContext": "Context"
                        }
                    }
                ]
            }
            save_test_cases(config, test_data)
            
            tuner = EvolutionaryTuner(config)
            results = await tuner.run_evolution()
            
            # Verify prompts directory exists
            assert config.prompts_dir.exists()
            
            # Verify prompt files exist
            prompt_files = list(config.prompts_dir.glob("*.txt"))
            assert len(prompt_files) > 0
            
            # Verify we can load the prompts
            from prompt_tuning_config import load_prompts
            loaded_prompts = load_prompts(config)
            
            assert len(loaded_prompts) > 0
            assert loaded_prompts == results["best_prompts"]


if __name__ == "__main__":
    pytest.main([__file__, "-v"])

