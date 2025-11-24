"""Tests for configuration and utility functions."""

import pytest
import os
import json
import tempfile
from pathlib import Path
from prompt_tuning_config import (
    PromptTuningConfig,
    load_project_config,
    load_test_cases,
    save_test_cases,
    load_prompts,
    save_prompts,
    validate_api_keys,
    get_env_file_template
)


class TestPromptTuningConfig:
    """Tests for PromptTuningConfig dataclass."""

    def test_config_creation(self):
        """Test creating a configuration object."""
        config = PromptTuningConfig(
            project_name="test",
            base_dir=Path("/tmp"),
            results_dir=Path("/tmp/results"),
            prompts_dir=Path("/tmp/prompts"),
            test_cases_file=Path("/tmp/test_cases.json")
        )
        
        assert config.project_name == "test"
        assert config.llm_provider == "anthropic"
        assert config.max_iterations == 20
        assert config.mock_mode is False

    def test_config_with_custom_values(self):
        """Test configuration with custom values."""
        config = PromptTuningConfig(
            project_name="custom",
            base_dir=Path("/tmp"),
            results_dir=Path("/tmp/results"),
            prompts_dir=Path("/tmp/prompts"),
            test_cases_file=Path("/tmp/test.json"),
            mock_mode=True,
            max_iterations=10
        )
        
        assert config.mock_mode is True
        assert config.max_iterations == 10


class TestLoadProjectConfig:
    """Tests for load_project_config function."""

    def test_load_config_basic(self):
        """Test loading basic configuration."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("test", base_dir=Path(tmpdir))
            
            assert config.project_name == "test"
            assert config.base_dir == Path(tmpdir)
            assert "test" in str(config.results_dir)

    def test_load_config_mock_mode_from_env(self, monkeypatch):
        """Test loading configuration with mock mode from environment."""
        monkeypatch.setenv("AI_AGENT_MOCK_MODE", "true")
        
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("test", base_dir=Path(tmpdir))
            
            assert config.mock_mode is True

    def test_load_config_creates_paths(self):
        """Test that configuration has correct paths."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("myproject", base_dir=Path(tmpdir))
            
            assert "myproject" in str(config.results_dir)
            assert config.prompts_dir.name == "prompts"
            assert "test_cases_myproject.json" in str(config.test_cases_file)


class TestValidateAPIKeys:
    """Tests for validate_api_keys function."""

    def test_validate_mock_mode_no_keys_required(self):
        """Test that mock mode doesn't require API keys."""
        config = PromptTuningConfig(
            project_name="test",
            base_dir=Path("/tmp"),
            results_dir=Path("/tmp/results"),
            prompts_dir=Path("/tmp/prompts"),
            test_cases_file=Path("/tmp/test.json"),
            mock_mode=True
        )
        
        missing = validate_api_keys(config)
        
        assert missing == []

    def test_validate_anthropic_key_missing(self, monkeypatch):
        """Test validation when Anthropic key is missing."""
        monkeypatch.delenv("ANTHROPIC_API_KEY", raising=False)
        
        config = PromptTuningConfig(
            project_name="test",
            base_dir=Path("/tmp"),
            results_dir=Path("/tmp/results"),
            prompts_dir=Path("/tmp/prompts"),
            test_cases_file=Path("/tmp/test.json"),
            llm_provider="anthropic",
            mock_mode=False
        )
        
        missing = validate_api_keys(config)
        
        assert "ANTHROPIC_API_KEY" in missing

    def test_validate_anthropic_key_present(self, monkeypatch):
        """Test validation when Anthropic key is present."""
        monkeypatch.setenv("ANTHROPIC_API_KEY", "test-key")
        
        config = PromptTuningConfig(
            project_name="test",
            base_dir=Path("/tmp"),
            results_dir=Path("/tmp/results"),
            prompts_dir=Path("/tmp/prompts"),
            test_cases_file=Path("/tmp/test.json"),
            llm_provider="anthropic",
            mock_mode=False
        )
        
        missing = validate_api_keys(config)
        
        assert missing == []


class TestTestCasesIO:
    """Tests for test cases loading and saving."""

    def test_save_and_load_test_cases(self):
        """Test saving and loading test cases."""
        with tempfile.TemporaryDirectory() as tmpdir:
            config = load_project_config("test", base_dir=Path(tmpdir))
            
            test_data = {
                "test_cases": [
                    {"id": "test1", "name": "Test 1"},
                    {"id": "test2", "name": "Test 2"}
                ]
            }
            
            save_test_cases(config, test_data)
            loaded = load_test_cases(config)
            
            assert loaded == test_data
            assert len(loaded["test_cases"]) == 2


if __name__ == "__main__":
    pytest.main([__file__, "-v"])

