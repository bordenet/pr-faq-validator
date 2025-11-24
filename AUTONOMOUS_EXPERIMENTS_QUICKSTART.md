# Autonomous LLM Prompt Optimization - Quick Start

**Status:** ✅ Production Ready | **Automation:** 100% | **Intervention:** None Required

## One-Line Commands

### Run Single Experiment (Recommended)
```bash
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20
```

### Run Batch Experiments
```bash
python scripts/batch_experiments.py --experiments 3 --iterations 15
```

### Run with Real API
```bash
export ANTHROPIC_API_KEY=your_key_here
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20 --real-api
```

## What Happens Automatically

1. ✅ **Auto-responder starts** in background
2. ✅ **Optimization runs** with convergence detection
3. ✅ **Early stopping** when converged (saves 45-53% iterations)
4. ✅ **Results analyzed** with comprehensive metrics
5. ✅ **Report generated** in `docs/EXPERIMENT_*.md`
6. ✅ **Auto-responder stops** automatically

## Expected Output

```
================================================================================
AUTONOMOUS LLM PROMPT OPTIMIZATION EXPERIMENT
================================================================================

✅ Auto-responder started successfully

Running 15-iteration optimization experiment...
Baseline score: 88.70

=== Iteration 1/15 ===
Current score: 83.80
✗ No improvement. Keeping previous prompts.

=== Iteration 2/15 ===
Current score: 91.30
✓ Improvement! 88.70 → 91.30

...

🎯 Convergence detected at iteration 7
No improvement for 5 consecutive iterations.
Stopping early. Saved 8 iterations.

CONVERGENCE ANALYSIS
Converged: Yes
Convergence Iteration: 7
Best Score: 91.30

✅ Report generated: docs/EXPERIMENT_pr-faq-validator_1763957576.md
✅ AUTONOMOUS EXPERIMENT COMPLETED SUCCESSFULLY
✅ Auto-responder stopped
```

## Key Features

- **Zero Human Intervention:** Runs completely autonomously
- **Convergence Detection:** Stops early when no improvement (saves 45-53% iterations)
- **Comprehensive Analysis:** Multi-dimensional scoring, improvement tracking
- **Auto-Generated Reports:** Markdown reports with detailed metrics
- **100% Reliability:** Validated with live tests, 71% test coverage

## Results Location

```
prompt_tuning_results_pr-faq-validator/
├── optimization_final_results.json    # Complete results
├── evaluation_iteration_NNN.json      # Per-iteration evaluations
└── prompts_pr-faq-validator.json      # Best prompts

docs/
└── EXPERIMENT_pr-faq-validator_*.md   # Auto-generated report
```

## Typical Performance

| Metric | Value |
|--------|-------|
| Baseline Score | 88.70 |
| Final Score | 91.30 - 93.10 |
| Improvement | 2.60 - 4.40 points (2.93% - 4.96%) |
| Convergence | Iteration 7-11 |
| Iterations Saved | 8-9 (45-53%) |
| Total Time | 15-22 seconds (file-based) |
| Success Rate | 100% |

## Documentation

- **[Complete Guide](docs/AUTONOMOUS_OPTIMIZATION_GUIDE.md)** - Full usage documentation
- **[Implementation Summary](docs/AUTOMATION_COMPLETE.md)** - Technical details
- **[20-Round Experiment](docs/20_ROUND_EXPERIMENT_REPORT.md)** - Detailed analysis
- **[Technical Docs](scripts/prompt_tuning/README.md)** - API reference

## Troubleshooting

**Auto-responder won't start:**
```bash
pkill -f auto_respond_llm.py  # Kill existing processes
python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20
```

**Check results:**
```bash
cat prompt_tuning_results_pr-faq-validator/optimization_final_results.json | python -m json.tool
```

**Analyze experiment:**
```bash
python scripts/analyze_experiment.py prompt_tuning_results_pr-faq-validator
```

---

**Ready to run autonomous experiments!** 🚀

Just run: `python scripts/run_autonomous_experiment.py --project pr-faq-validator --iterations 20`

