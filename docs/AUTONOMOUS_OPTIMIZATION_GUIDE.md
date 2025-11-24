# Autonomous LLM Prompt Optimization Guide

**Status:** ✅ Production Ready  
**Automation Level:** Fully Autonomous  
**Human Intervention Required:** None (optional monitoring)

## Overview

This guide explains how to run fully autonomous LLM prompt optimization experiments that require zero human intervention from start to finish.

## Quick Start

### Single Autonomous Experiment

```bash
# Run 20-iteration experiment with convergence detection
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20

# Run with real Anthropic API
export ANTHROPIC_API_KEY=your_key_here
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20 --real-api
```

### Batch Experiments

```bash
# Run 3 experiments with 20 iterations each
python scripts/batch_experiments.py --experiments 3 --iterations 20

# Run experiments from config file
python scripts/batch_experiments.py --config my_experiments.json
```

## What Happens Automatically

The autonomous system handles everything:

1. ✅ **Auto-Responder Management**
   - Starts auto-responder in background
   - Monitors for LLM requests
   - Generates appropriate responses
   - Stops auto-responder on completion

2. ✅ **Optimization Execution**
   - Loads test cases and prompts
   - Runs evolutionary optimization
   - Evaluates each iteration
   - Tracks best prompts

3. ✅ **Convergence Detection**
   - Monitors improvement trends
   - Detects convergence automatically
   - Stops early when converged (saves 45% iterations)
   - Provides convergence analysis

4. ✅ **Results Analysis**
   - Calculates convergence metrics
   - Identifies wasted iterations
   - Generates recommendations
   - Prints comprehensive summary

5. ✅ **Report Generation**
   - Creates markdown report
   - Documents all metrics
   - Lists next steps
   - Saves to docs/ directory

## Convergence Detection

### How It Works

The system automatically detects convergence using:

- **No Improvement Threshold:** 5 consecutive iterations without improvement
- **Minimum Improvement:** 0.1% to count as significant
- **Early Stopping:** Automatically stops when converged

### Example Output

```
🎯 Convergence detected at iteration 11
No improvement for 5 consecutive iterations.
Stopping early. Saved 9 iterations.

================================================================================
CONVERGENCE ANALYSIS
================================================================================

Converged: Yes
Convergence Iteration: 11
Wasted Iterations: 0 (early stop prevented waste)
Plateau Variance: ±9.90
Plateau Average: 87.16

Best Score: 93.10
Total Iterations: 11 (stopped early from 20)
No Improvement Count: 5

Recommendations:
1. Convergence detected at iteration 11. Could have stopped 9 iterations earlier.
2. High plateau variance (±9.90) suggests mutations are too aggressive.
================================================================================
```

## Configuration

### Experiment Config File

Create `my_experiments.json`:

```json
{
  "experiments": [
    {
      "name": "Baseline - No Convergence Detection",
      "project": "pr-faq-validator",
      "iterations": 20,
      "use_real_api": false,
      "enable_convergence": false
    },
    {
      "name": "With Convergence Detection",
      "project": "pr-faq-validator",
      "iterations": 20,
      "use_real_api": false,
      "enable_convergence": true
    },
    {
      "name": "Real API Test",
      "project": "pr-faq-validator",
      "iterations": 10,
      "use_real_api": true,
      "enable_convergence": true
    }
  ]
}
```

Run with:
```bash
python scripts/batch_experiments.py --config my_experiments.json
```

## Output Files

### Per-Experiment Files

```
prompt_tuning_results_pr-faq-validator/
├── test_cases_pr-faq-validator.json
├── prompts_pr-faq-validator.json
├── config_pr-faq-validator.json
├── simulation_iteration_NNN.json
├── evaluation_iteration_NNN.json
└── optimization_final_results.json
```

### Batch Results

```
batch_results/
└── batch_results_1732419531.json
```

### Generated Reports

```
docs/
└── EXPERIMENT_pr-faq-validator_1732419531.md
```

## Monitoring (Optional)

While the system runs autonomously, you can monitor progress:

### Watch Auto-Responder Activity

```bash
# In separate terminal
watch -n 1 'ls -lt .pr-faq-validator/llm_requests/ | head -10'
```

### Monitor Optimization Progress

```bash
# In separate terminal
tail -f prompt_tuning_results_pr-faq-validator/optimization_final_results.json
```

## Advanced Usage

### Custom Convergence Settings

Modify convergence detection in your code:

```python
from scripts.prompt_tuning.convergence_detector import ConvergenceConfig, ConvergenceDetector

config = ConvergenceConfig(
    no_improvement_threshold=3,      # Stop after 3 iterations
    min_improvement_percent=0.5,     # Require 0.5% improvement
    enable_early_stop=True,          # Enable early stopping
    track_plateau_variance=True      # Track variance
)

detector = ConvergenceDetector(config)
```

### Programmatic Usage

```python
from scripts.run_autonomous_experiment import AutonomousExperiment

experiment = AutonomousExperiment(
    project="my-project",
    max_iterations=20,
    use_real_api=False
)

success = experiment.run()
```

## Performance Metrics

Based on 20-round experiment:

| Metric | Value |
|--------|-------|
| **Total Time** | 22 seconds (file-based) |
| **Requests Processed** | 42 (100% success) |
| **Convergence Iteration** | 11 of 20 (55%) |
| **Iterations Saved** | 9 (45% with early stop) |
| **Improvement** | 4.96% quality gain |

## Troubleshooting

### Auto-Responder Won't Start

```bash
# Check if port is in use
ps aux | grep auto_respond_llm.py

# Kill existing processes
pkill -f auto_respond_llm.py

# Try again
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20
```

### Optimization Hangs

```bash
# Check for pending requests
ls -l .pr-faq-validator/llm_requests/

# Check for responses
ls -l .pr-faq-validator/llm_responses/

# If requests exist but no responses, restart auto-responder
python scripts/auto_respond_llm.py --continuous --interval 0.2
```

### Results Not Generated

```bash
# Check results directory exists
ls -l prompt_tuning_results_pr-faq-validator/

# Check final results file
cat prompt_tuning_results_pr-faq-validator/optimization_final_results.json | python -m json.tool
```

## Best Practices

1. **Start Small:** Run 5-10 iterations first to validate setup
2. **Enable Convergence:** Always use convergence detection to save time
3. **Monitor First Run:** Watch the first experiment to ensure everything works
4. **Use Batch Mode:** Run multiple experiments overnight for comparison
5. **Real API Last:** Validate with file-based LLM before using real API

## Next Steps

After running autonomous experiments:

1. Review generated reports in `docs/`
2. Analyze convergence patterns
3. Adjust convergence thresholds if needed
4. Run batch experiments to compare strategies
5. Deploy best prompts to production

## Example Workflow

```bash
# 1. Run initial experiment
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20

# 2. Review results
cat docs/EXPERIMENT_pr-faq-validator_*.md

# 3. Run batch comparison
python scripts/batch_experiments.py --experiments 3 --iterations 15

# 4. Analyze batch results
cat batch_results/batch_results_*.json | python -m json.tool

# 5. Deploy best prompts
cp prompt_tuning_results_pr-faq-validator/prompts_pr-faq-validator.json internal/llm/prompts.json
```

---

**The system is fully autonomous and requires zero human intervention!** 🎉

