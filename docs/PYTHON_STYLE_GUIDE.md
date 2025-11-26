# Python Style Guide for pr-faq-validator

This document defines the Python coding standards for tools development in this repository.
Standards are derived from PEP 8, PEP 257, and the Google Python Style Guide, with automated
enforcement via Black, isort, pylint, and mypy.

## References

- [PEP 8 – Style Guide for Python Code](https://peps.python.org/pep-0008/) - Primary baseline
- [PEP 257 – Docstring Conventions](https://peps.python.org/pep-0257/) - Documentation standards
- [Google Python Style Guide](https://google.github.io/styleguide/pyguide.html) - Stricter conventions
- [Black Code Formatter](https://black.readthedocs.io/) - Opinionated formatting
- [Ruff](https://docs.astral.sh/ruff/) - Fast Python linter (optional, replacing flake8)

## Automated Tooling

All code must pass these tools before commit:

```bash
# Formatting (automatic)
black --line-length=120 scripts/
isort --profile black --line-length 120 scripts/

# Linting (must pass with score >= 9.5)
pylint scripts/ --max-line-length=120 --min-similarity-lines=10 --fail-under=9.5

# Type checking (must pass)
mypy scripts/ --ignore-missing-imports

# Tests (must pass with >= 50% coverage)
pytest tests/ --cov=scripts --cov-fail-under=50
```

## Project Structure

```text
scripts/
  prompt_tuning/
    __init__.py          # Package exports
    llm_client.py        # LLM client abstraction
    optimizer.py         # Prompt optimization logic
    evaluator.py         # Evaluation logic
  auto_respond_llm.py    # Auto-responder for testing
  run_*.py               # CLI entry points
tests/
  python/
    test_*.py            # Test modules
    conftest.py          # Shared fixtures
```

## Naming Conventions

### Modules and Packages
- **snake_case**: `llm_client.py`, `prompt_optimizer.py`
- Avoid: `llmClient.py`, `LLMClient.py`

### Classes
- **PascalCase**: `LLMClient`, `PromptOptimizer`, `EvaluationResult`
- Suffix base classes with `Base`: `ClientBase`
- Suffix abstract classes with `ABC`: `OptimizerABC`

### Functions and Methods
- **snake_case**: `optimize_prompt()`, `get_response()`, `_private_helper()`
- Use verbs: `calculate_`, `generate_`, `validate_`, `parse_`
- Boolean methods: `is_valid()`, `has_errors()`, `can_process()`

### Variables
- **snake_case**: `file_path`, `error_count`, `is_valid`
- **UPPER_SNAKE_CASE** for constants: `MAX_ITERATIONS`, `DEFAULT_TIMEOUT`
- Prefix private with underscore: `_internal_state`

## Function and Method Guidelines

### Length
- Target: **≤50 lines** per function
- Maximum: **100 lines** (refactor if approaching)
- Single responsibility principle

### Parameters
- **≤5 parameters** - use dataclass/dict if more needed
- Use keyword arguments for optional parameters
- Type all parameters and return values

```python
# Good
def optimize_prompt(
    prompt: str,
    *,
    max_iterations: int = 10,
    temperature: float = 0.7,
) -> OptimizationResult:
    """Optimize a prompt using evolutionary algorithms."""
    ...

# Bad - too many positional parameters
def optimize_prompt(prompt, iterations, temp, model, timeout, retries):
    ...
```

## Documentation (PEP 257)

### Module Docstrings
```python
"""Prompt optimization using evolutionary algorithms.

This module provides functionality for iterative prompt improvement including:
- Mutation strategies
- Fitness evaluation
- Population management
"""
```

### Function Docstrings
```python
def optimize_prompt(prompt: str, *, max_iterations: int = 10) -> OptimizationResult:
    """Optimize a prompt using evolutionary algorithms.

    Args:
        prompt: The initial prompt to optimize.
        max_iterations: Maximum optimization iterations. Defaults to 10.

    Returns:
        OptimizationResult containing the best prompt and score.

    Raises:
        ValueError: If prompt is empty.
        TimeoutError: If optimization exceeds time limit.
    """
```

## Type Annotations

### Always Type
- All function parameters
- All return values (including `-> None`)
- Class attributes
- Module-level variables

```python
# Good
def calculate_score(content: str, weight: float = 1.0) -> float:
    """Calculate weighted score from content."""
    return len(content) * weight

# Bad - no type annotations
def calculate_score(content, weight=1.0):
    return len(content) * weight
```

## Error Handling

### Exception Hierarchy
```python
class PromptTuningError(Exception):
    """Base exception for prompt tuning."""

class ValidationError(PromptTuningError):
    """Raised when validation fails."""

class LLMError(PromptTuningError):
    """Raised when LLM call fails."""
```

### Error Messages
```python
# Good - descriptive with context
raise ValueError(f"Invalid prompt: {prompt!r} (must not be empty)")

# Bad - vague
raise ValueError("invalid prompt")
```

## Testing

### Test Structure
```python
import pytest
from scripts.prompt_tuning.optimizer import PromptOptimizer

class TestPromptOptimizer:
    """Tests for PromptOptimizer class."""

    def test_optimize_empty_prompt(self) -> None:
        """Empty prompts should raise ValueError."""
        optimizer = PromptOptimizer()
        with pytest.raises(ValueError):
            optimizer.optimize("")
```

### Coverage
- Target: **≥50%** overall (enforced in CI)
- Critical modules: **≥80%** coverage
- Use `# pragma: no cover` sparingly and with justification

## Linting Rules Summary

### pylint
- Required score: **≥9.5/10**
- Key rules: `missing-docstring`, `too-many-arguments`, `unused-import`

### mypy
- Minimum: `--ignore-missing-imports`
- Strict mode recommended for new code

### Black
- Line length: **120** (configured in pre-commit)
- No configuration needed - use defaults

### isort
- Profile: **black** (compatible formatting)
- Line length: **120**

