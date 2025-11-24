# LLM Prompt Tuning Infrastructure

**Status:** ✅ Production Ready  
**Test Coverage:** 71% (31 passing tests)  
**Architecture:** Bloginator-inspired file-based LLM communication

## Overview

A comprehensive prompt optimization system using evolutionary algorithms and LLM-as-judge evaluation. Supports both real API calls (Anthropic Claude) and file-based communication for cost-free development and testing.

## Features

### Core Capabilities
- ✅ **Evolutionary Optimization** - Iterative prompt improvement using mutations and selection
- ✅ **Multi-Dimensional Scoring** - 8-dimensional quality evaluation (clarity, depth, nuance, specificity, etc.)
- ✅ **File-Based LLM Communication** - Zero-cost development using AI assistant as LLM
- ✅ **Auto-Responder System** - Autonomous request processing for experiments
- ✅ **Comprehensive Evaluation** - LLM-as-judge with detailed feedback
- ✅ **Evolutionary Strategy** - AI-driven recommendations for prompt improvements

### Validated Performance
- ✅ **4.96% quality improvement** in 20-round experiment
- ✅ **100% system reliability** (42/42 requests processed)
- ✅ **Clear convergence behavior** (plateau detection)
- ✅ **22-second execution** for 20 iterations (file-based mode)

## Quick Start

### 1. Initialize Project

```bash
python prompt_tuning_tool.py init my-project
```

This creates:
- `prompt_tuning_results_my-project/` directory
- `test_cases_my-project.json` with sample test cases
- `prompts_my-project.json` with default prompts

### 2. Customize Test Cases

Edit `prompt_tuning_results_my-project/test_cases_my-project.json`:

```json
{
  "test_cases": [
    {
      "id": "tc001",
      "name": "Your Test Case",
      "inputs": {
        "projectName": "Your Project",
        "problemDescription": "Problem you're solving",
        "businessContext": "Business goals and metrics"
      },
      "metadata": {
        "industry": "Your Industry",
        "project_type": "Feature/Product/Tool",
        "scope": "Small/Medium/Large"
      }
    }
  ]
}
```

### 3. Run Optimization

**With Real API (Anthropic Claude):**
```bash
export ANTHROPIC_API_KEY=your_key_here
python prompt_tuning_tool.py evolve my-project --max-iterations 20
```

**With File-Based LLM (Free):**
```bash
# Terminal 1: Start auto-responder
python scripts/auto_respond_llm.py --continuous --interval 0.2

# Terminal 2: Run optimization
LLM_PROVIDER=assistant python prompt_tuning_tool.py evolve my-project --max-iterations 20
```

**With Mock Mode (Testing):**
```bash
python prompt_tuning_tool.py evolve my-project --mock --max-iterations 5
```

## Architecture

### Components

```
scripts/prompt_tuning/
├── llm_client.py              # LLM client abstraction (Anthropic, Assistant, Mock)
├── prompt_simulator.py        # PR-FAQ content generation
├── quality_evaluator.py       # Multi-dimensional evaluation
├── evolutionary_tuner.py      # Evolutionary optimization algorithm
├── prompt_tuning_config.py    # Configuration management
└── prompt_tuning_cli.py       # CLI interface

scripts/
└── auto_respond_llm.py        # Auto-responder for file-based LLM

prompt_tuning_tool.py          # Main entry point
```

### Data Flow

```
1. Load test cases and prompts
2. For each iteration:
   a. Simulate: Generate PR-FAQ content using current prompts
   b. Evaluate: Score content using LLM-as-judge
   c. Mutate: Create prompt variations
   d. Select: Keep best-performing prompts
3. Save results and best prompts
```

## File-Based LLM Communication

### How It Works

1. **Request Phase:**
   - Optimization writes JSON request to `.pr-faq-validator/llm_requests/request_NNNN.json`
   - Request includes prompt, temperature, max_tokens, etc.

2. **Response Phase:**
   - Auto-responder monitors request directory
   - Generates appropriate response (evaluation, content, etc.)
   - Writes JSON response to `.pr-faq-validator/llm_responses/response_NNNN.json`

3. **Polling:**
   - Optimization polls for response file (500ms interval, 300s timeout)
   - Reads response and continues processing

### Auto-Responder Usage

```bash
# Continuous mode (runs until stopped)
python scripts/auto_respond_llm.py --continuous --interval 0.2

# One-shot mode (process existing requests and exit)
python scripts/auto_respond_llm.py

# Help
python scripts/auto_respond_llm.py --help
```

## Configuration

### Environment Variables

```bash
# LLM Provider (anthropic, assistant, or mock)
export LLM_PROVIDER=assistant

# Anthropic API Key (for real API calls)
export ANTHROPIC_API_KEY=your_key_here

# Mock mode (for testing)
export AI_AGENT_MOCK_MODE=true
```

### Config File

Edit `prompt_tuning_results_my-project/config_my-project.json`:

```json
{
  "project_name": "my-project",
  "llm_provider": "anthropic",
  "generator_model": "claude-3-5-sonnet-20241022",
  "evaluator_model": "claude-3-5-sonnet-20241022",
  "temperature": 0.7,
  "max_tokens": 3000,
  "validation_criteria": {
    "press_release_quality": 0.3,
    "faq_completeness": 0.25,
    "clarity_score": 0.25,
    "structure_adherence": 0.2
  }
}
```

## Results Analysis

### View Results

```bash
# Analyze experiment results
python scripts/analyze_experiment.py prompt_tuning_results_my-project

# View final results
cat prompt_tuning_results_my-project/optimization_final_results.json | python -m json.tool

# View specific iteration
cat prompt_tuning_results_my-project/evaluation_iteration_011.json | python -m json.tool
```

### Result Files

```
prompt_tuning_results_my-project/
├── test_cases_my-project.json           # Test cases
├── prompts_my-project.json              # Current prompts
├── config_my-project.json               # Configuration
├── simulation_iteration_NNN.json        # Generated content per iteration
├── evaluation_iteration_NNN.json        # Evaluation scores per iteration
└── optimization_final_results.json      # Final summary
```

## Testing

```bash
cd scripts/prompt_tuning

# Run all tests
python -m pytest -v

# Run with coverage
python -m pytest -v --cov=. --cov-report=term-missing

# Run specific test file
python -m pytest test_llm_client.py -v
```

## Documentation

- **[PRD](../../docs/PROMPT_TUNING_PRD.md)** - Product requirements
- **[Design](../../docs/PROMPT_TUNING_DESIGN.md)** - Technical design
- **[Bloginator Comparison](../../docs/BLOGINATOR_COMPARISON.md)** - Feature comparison
- **[20-Round Experiment](../../docs/20_ROUND_EXPERIMENT_REPORT.md)** - Detailed experiment results
- **[Experiment Summary](../../docs/EXPERIMENT_SUMMARY.md)** - Executive summary
- **[Test Summary](TEST_SUMMARY.md)** - Test coverage details

## Known Limitations

1. **No Convergence Detection** - System doesn't auto-stop when converged (45% wasted iterations in 20-round test)
2. **Fixed Mutation Strength** - Doesn't adapt mutation aggressiveness near convergence
3. **Limited Iteration Tracking** - Per-iteration multi-dimensional trends not tracked

## Roadmap

### Priority 1: Convergence Detection
- Auto-stop after 5 consecutive non-improvements
- Expected savings: 45% fewer iterations

### Priority 2: Adaptive Mutation
- Reduce mutation strength near convergence
- Better exploration/exploitation balance

### Priority 3: Enhanced Metrics
- Per-iteration multi-dimensional score tracking
- Trend analysis and visualization

## License

Part of the pr-faq-validator project.

