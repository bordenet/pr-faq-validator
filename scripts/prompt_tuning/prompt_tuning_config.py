"""Configuration management for PR-FAQ Validator prompt tuning."""

import os
import json
from pathlib import Path
from typing import Dict, List, Any, Optional
from dataclasses import dataclass, field


@dataclass
class PromptTuningConfig:
    """Configuration for prompt tuning."""
    
    project_name: str
    base_dir: Path
    results_dir: Path
    prompts_dir: Path
    test_cases_file: Path
    llm_provider: str = "anthropic"
    llm_model: str = "claude-3-5-sonnet-20241022"
    evaluator_model: str = "claude-3-5-sonnet-20241022"
    max_iterations: int = 20
    temperature: float = 1.0
    mock_mode: bool = False
    
    # PR-FAQ specific settings
    validation_criteria: Dict[str, Any] = field(default_factory=lambda: {
        "press_release_quality": 0.3,
        "faq_completeness": 0.25,
        "clarity_score": 0.25,
        "structure_adherence": 0.2
    })


def load_project_config(project_name: str, base_dir: Optional[Path] = None) -> PromptTuningConfig:
    """Load configuration for a prompt tuning project."""
    if base_dir is None:
        base_dir = Path.cwd()
    
    results_dir = base_dir / f"prompt_tuning_results_{project_name}"
    prompts_dir = base_dir / "prompts"
    test_cases_file = results_dir / f"test_cases_{project_name}.json"
    
    # Check if running in mock mode
    mock_mode = os.getenv("AI_AGENT_MOCK_MODE", "false").lower() == "true"
    
    config = PromptTuningConfig(
        project_name=project_name,
        base_dir=base_dir,
        results_dir=results_dir,
        prompts_dir=prompts_dir,
        test_cases_file=test_cases_file,
        llm_provider=os.getenv("LLM_PROVIDER", "anthropic"),
        llm_model=os.getenv("LLM_MODEL", "claude-3-5-sonnet-20241022"),
        evaluator_model=os.getenv("EVALUATOR_MODEL", "claude-3-5-sonnet-20241022"),
        mock_mode=mock_mode
    )
    
    return config


def validate_api_keys(config: PromptTuningConfig) -> List[str]:
    """Validate that required API keys are present."""
    if config.mock_mode:
        return []
    
    missing_keys = []
    
    if config.llm_provider == "anthropic":
        if not os.getenv("ANTHROPIC_API_KEY"):
            missing_keys.append("ANTHROPIC_API_KEY")
    elif config.llm_provider == "openai":
        if not os.getenv("OPENAI_API_KEY"):
            missing_keys.append("OPENAI_API_KEY")
    
    return missing_keys


def get_env_file_template() -> str:
    """Get template for .env file."""
    return """# PR-FAQ Validator Prompt Tuning Configuration

# LLM Provider (anthropic or openai)
LLM_PROVIDER=anthropic

# API Keys (uncomment and fill in the one you're using)
# ANTHROPIC_API_KEY=your_key_here
# OPENAI_API_KEY=your_key_here

# Model Configuration
LLM_MODEL=claude-3-5-sonnet-20241022
EVALUATOR_MODEL=claude-3-5-sonnet-20241022

# AI Agent Mock Mode (set to true to run without API keys)
AI_AGENT_MOCK_MODE=false
"""


def load_test_cases(config: PromptTuningConfig) -> Dict[str, Any]:
    """Load test cases from JSON file."""
    if not config.test_cases_file.exists():
        raise FileNotFoundError(
            f"Test cases file not found: {config.test_cases_file}\n"
            f"Run 'python prompt_tuning_tool.py init {config.project_name}' first"
        )
    
    with open(config.test_cases_file, 'r') as f:
        return json.load(f)


def save_test_cases(config: PromptTuningConfig, test_cases: Dict[str, Any]) -> None:
    """Save test cases to JSON file."""
    config.results_dir.mkdir(parents=True, exist_ok=True)
    
    with open(config.test_cases_file, 'w') as f:
        json.dump(test_cases, f, indent=2)


def load_prompts(config: PromptTuningConfig) -> Dict[str, str]:
    """Load prompts from prompts directory."""
    prompts = {}
    
    if not config.prompts_dir.exists():
        return prompts
    
    for prompt_file in config.prompts_dir.glob("*.txt"):
        prompt_name = prompt_file.stem
        with open(prompt_file, 'r') as f:
            prompts[prompt_name] = f.read()
    
    return prompts


def save_prompts(config: PromptTuningConfig, prompts: Dict[str, str]) -> None:
    """Save prompts to prompts directory."""
    config.prompts_dir.mkdir(parents=True, exist_ok=True)
    
    for prompt_name, prompt_content in prompts.items():
        prompt_file = config.prompts_dir / f"{prompt_name}.txt"
        with open(prompt_file, 'w') as f:
            f.write(prompt_content)

