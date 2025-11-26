"""Tests for core components (simulator, evaluator, tuner)."""

import json
import tempfile
from pathlib import Path

import pytest

from scripts.prompt_tuning.evolutionary_tuner import EvolutionaryTuner
from scripts.prompt_tuning.prompt_simulator import PromptSimulator
from scripts.prompt_tuning.prompt_tuning_config import load_project_config, save_prompts, save_test_cases
from scripts.prompt_tuning.quality_evaluator import QualityEvaluator


class TestPromptSimulator:
    """Tests for PromptSimulator class."""

    @pytest.fixture
    def config_with_test_cases(self):
        """Create a configuration with test cases."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("test", base_dir=Path(tmpdir))
            config.mock_mode = True

            # Create test cases
            test_data = {
                "test_cases": [
                    {
                        "id": "test1",
                        "name": "Test Case 1",
                        "inputs": {
                            "projectName": "Test Project",
                            "problemDescription": "Test problem",
                            "businessContext": "Test context",
                        },
                    }
                ]
            }
            save_test_cases(config, test_data)

            # Create prompts
            prompts = {"pr_faq_generation": "Generate PR-FAQ for {projectName}"}
            save_prompts(config, prompts)

            yield config

    async def test_simulator_initialization(self, config_with_test_cases):
        """Test that simulator initializes correctly."""
        simulator = PromptSimulator(config_with_test_cases)

        assert simulator.config == config_with_test_cases
        assert simulator.llm_client is not None

    async def test_run_simulation(self, config_with_test_cases):
        """Test running a simulation."""
        simulator = PromptSimulator(config_with_test_cases)

        results = await simulator.run_simulation(iteration=0)

        assert "iteration" in results
        assert results["iteration"] == 0
        assert "test_case_results" in results
        assert len(results["test_case_results"]) == 1

    async def test_save_results(self, config_with_test_cases):
        """Test saving simulation results."""
        simulator = PromptSimulator(config_with_test_cases)

        results = await simulator.run_simulation(iteration=0)
        output_file = simulator.save_results(results, iteration=0)

        assert output_file.exists()
        assert "simulation_iteration_000.json" in str(output_file)

        # Verify content
        with open(output_file, encoding="utf-8") as f:
            saved_data = json.load(f)

        assert saved_data["iteration"] == 0

    def test_format_prompt(self, config_with_test_cases):
        """Test prompt formatting."""
        simulator = PromptSimulator(config_with_test_cases)

        template = "Project: {projectName}, Problem: {problemDescription}"
        inputs = {"projectName": "MyProject", "problemDescription": "MyProblem"}

        formatted = simulator._format_prompt(template, inputs)

        assert "MyProject" in formatted
        assert "MyProblem" in formatted
        assert "{projectName}" not in formatted

    def test_extract_section(self, config_with_test_cases):
        """Test section extraction."""
        simulator = PromptSimulator(config_with_test_cases)

        content = """# Press Release

This is the press release content.

## FAQ

This is the FAQ content.
"""

        pr_section = simulator._extract_section(content, "press release")
        faq_section = simulator._extract_section(content, "faq")

        assert "press release content" in pr_section
        assert "FAQ content" in faq_section


class TestQualityEvaluator:
    """Tests for QualityEvaluator class."""

    @pytest.fixture
    def config(self):
        """Create a test configuration."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("test", base_dir=Path(tmpdir))
            config.mock_mode = True
            yield config

    async def test_evaluator_initialization(self, config):
        """Test that evaluator initializes correctly."""
        evaluator = QualityEvaluator(config)

        assert evaluator.config == config
        assert evaluator.evaluator_client is not None

    async def test_evaluate_results(self, config):
        """Test evaluating simulation results."""
        evaluator = QualityEvaluator(config)

        simulation_results = {
            "iteration": 0,
            "test_case_results": [
                {
                    "test_case_id": "test1",
                    "generated_content": {
                        "full_content": "# Press Release\n\nTest content\n\n# FAQ\n\nQ: Test?\nA: Answer."
                    },
                }
            ],
        }

        eval_results = await evaluator.evaluate_results(simulation_results)

        assert "iteration" in eval_results
        assert "evaluations" in eval_results
        assert "aggregate_scores" in eval_results


class TestEvolutionaryTuner:
    """Tests for EvolutionaryTuner class."""

    @pytest.fixture
    def config_with_data(self):
        """Create a configuration with test data."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("test", base_dir=Path(tmpdir))
            config.mock_mode = True
            config.max_iterations = 2  # Small number for testing

            # Create test cases
            test_data = {
                "test_cases": [
                    {
                        "id": "test1",
                        "name": "Test Case 1",
                        "inputs": {
                            "projectName": "Test",
                            "problemDescription": "Problem",
                            "businessContext": "Context",
                        },
                    }
                ]
            }
            save_test_cases(config, test_data)

            yield config

    async def test_tuner_initialization(self, config_with_data):
        """Test that tuner initializes correctly."""
        tuner = EvolutionaryTuner(config_with_data)

        assert tuner.config == config_with_data
        assert tuner.mutation_client is not None
        assert tuner.simulator is not None
        assert tuner.evaluator is not None


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
