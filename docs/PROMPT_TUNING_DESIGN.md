# PR-FAQ Validator Prompt Tuning Infrastructure - Design Document

## Overview

This document describes the LLM prompt tuning infrastructure for pr-faq-validator, implementing the evolutionary optimization methodology used in bloginator, one-pager, and product-requirements-assistant.

## Architecture

### Components

1. **Prompt Tuning Tool** (`prompt_tuning_tool.py`)
   - Main entry point for all prompt tuning operations
   - CLI interface for initialization, simulation, evaluation, and evolution

2. **Configuration Management** (`scripts/prompt_tuning/prompt_tuning_config.py`)
   - Project configuration loading and validation
   - API key management
   - Environment variable handling
   - Mock mode support

3. **LLM Client** (`scripts/prompt_tuning/llm_client.py`)
   - Abstract LLM client interface
   - Mock LLM client for testing without API keys
   - Anthropic Claude client implementation
   - Deterministic mock responses for reproducible testing

4. **Prompt Simulator** (`scripts/prompt_tuning/prompt_simulator.py`)
   - Executes test cases using current prompts
   - Generates PR-FAQ content via LLM
   - Saves simulation results for evaluation

5. **Quality Evaluator** (`scripts/prompt_tuning/quality_evaluator.py`)
   - Evaluates generated PR-FAQ content
   - Scores press release quality, FAQ completeness, clarity, and structure
   - Uses LLM-as-judge pattern for evaluation
   - Calculates aggregate scores across test cases

6. **Evolutionary Tuner** (`scripts/prompt_tuning/evolutionary_tuner.py`)
   - Implements evolutionary optimization loop
   - Mutates prompts using LLM
   - Keep/discard logic based on quality scores
   - Tracks iteration history and best results

### Workflow

```
1. Initialize Project
   └─> Create directory structure, test cases template, .env file

2. Run Baseline
   └─> Simulate with current prompts
   └─> Evaluate quality
   └─> Record baseline score

3. Evolutionary Loop (N iterations)
   ├─> Mutate prompts using LLM
   ├─> Simulate with mutated prompts
   ├─> Evaluate quality
   └─> Keep if improved, discard otherwise

4. Save Results
   └─> Best prompts, iteration history, final scores
```

## Quality Metrics

### Evaluation Criteria

- **Press Release Quality** (30%): Adherence to Amazon PR-FAQ format, clarity, compelling narrative
- **FAQ Completeness** (25%): Coverage of key stakeholder questions
- **Clarity Score** (25%): Readability and comprehension
- **Structure Adherence** (20%): Proper formatting and organization

### Scoring

- All scores on 0-100 scale
- Aggregate score weighted by criteria percentages
- LLM-as-judge evaluation with structured JSON output

## Mock Mode

### Purpose
Enable prompt tuning without API keys for:
- Development and testing
- CI/CD integration
- Demonstration and training

### Implementation
- Deterministic responses based on prompt hash
- Simulates press release, FAQ, and evaluation outputs
- Consistent scoring for reproducibility

## Usage

### Initialize Project
```bash
python prompt_tuning_tool.py init pr-faq-validator
```

### Run Simulation (Mock Mode)
```bash
python prompt_tuning_tool.py simulate pr-faq-validator --mock
```

### Run Evolutionary Optimization (Mock Mode)
```bash
python prompt_tuning_tool.py evolve pr-faq-validator --mock
```

### Run with Real API
```bash
# Set ANTHROPIC_API_KEY in .env
python prompt_tuning_tool.py evolve pr-faq-validator
```

## File Structure

```
pr-faq-validator/
├── prompt_tuning_tool.py              # Main entry point
├── docs/
│   ├── PROMPT_TUNING_DESIGN.md        # This file
│   └── PROMPT_TUNING_PRD.md           # Product requirements
├── scripts/
│   └── prompt_tuning/
│       ├── __init__.py
│       ├── prompt_tuning_cli.py       # CLI implementation
│       ├── prompt_tuning_config.py    # Configuration
│       ├── llm_client.py              # LLM clients
│       ├── prompt_simulator.py        # Simulation engine
│       ├── quality_evaluator.py       # Evaluation engine
│       └── evolutionary_tuner.py      # Optimization engine
├── prompts/
│   └── pr_faq_generation.txt          # Current prompts
└── prompt_tuning_results_pr-faq-validator/
    ├── test_cases_pr-faq-validator.json
    ├── simulation_iteration_*.json
    ├── evaluation_iteration_*.json
    └── optimization_final_results.json
```

## Testing Strategy

### Unit Tests
- Mock LLM client deterministic behavior
- Configuration loading and validation
- Prompt formatting and section extraction
- Score calculation and aggregation

### Integration Tests
- End-to-end simulation with mock LLM
- Evolutionary optimization convergence
- File I/O and persistence

### Coverage Target
- Minimum 85% code coverage
- All critical paths tested
- Edge cases and error handling

## Dependencies

### Required
- Python 3.11+
- click (CLI framework)
- rich (terminal UI)

### Optional
- anthropic (for real API calls)
- python-dotenv (for .env file support)

## Future Enhancements

1. **Multi-model Support**: OpenAI, local models via Ollama
2. **Advanced Mutations**: Crossover, ensemble strategies
3. **Automated Test Case Generation**: From existing PR-FAQs
4. **Web UI**: Interactive prompt tuning dashboard
5. **A/B Testing**: Compare multiple prompt variants

