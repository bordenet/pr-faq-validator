# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based CLI tool that analyzes PR-FAQ (Press Release - Frequently Asked Questions) documents using deterministic rule-based scoring (100% of numerical score) plus optional OpenAI GPT-4 qualitative feedback. The tool provides quality scoring, interactive terminal UI, and improvement recommendations.

## Development Commands

### Build and Run
```bash
# Install dependencies
go mod tidy

# Build the binary
go build

# Run the tool
./pr-faq-validator -file path/to/your/prfaq.md

# Test with example files
./pr-faq-validator -file testdata/example_prfaq_1.md
./pr-faq-validator -file testdata/example_prfaq_4.md  # High-quality example (77/100)

# Run without TUI (legacy output)
./pr-faq-validator -file testdata/example_prfaq_1.md -no-tui

# Generate markdown report
./pr-faq-validator -file testdata/example_prfaq_1.md -report
```

### Environment Setup
```bash
# Optional: Set OpenAI API key for AI feedback
export OPENAI_API_KEY=your_openai_api_key_here

# Tool works without API key - scoring is 100% deterministic
# AI feedback provides supplementary qualitative insights only
```

## Architecture

The codebase follows a clean modular structure:

- **main.go**: CLI entry point with interactive TUI (Bubble Tea) and legacy stdout modes
- **internal/parser/**: Document parsing, quality analysis, and deterministic scoring algorithms
- **internal/llm/**: OpenAI API integration for supplementary qualitative feedback
- **internal/ui/**: Interactive terminal UI components using Lip Gloss and Bubble Tea

### Key Components

**Parser (internal/parser/parser.go)**:
- `PRQualityBreakdown` struct with comprehensive scoring across 4 dimensions
- `ParsePRFAQ()` extracts sections with flexible header detection (H1/H2, various naming)
- `comprehensivePRAnalysis()` performs 100-point deterministic scoring:
  - Structure & Hook (30 pts): Headlines, newsworthy hooks, release dates
  - Content Quality (35 pts): 5 Ws coverage, credibility, structure
  - Professional Quality (20 pts): Tone, readability, fluff detection
  - Customer Evidence (15 pts): Quote quality with quantitative metrics
- Advanced features: quote metric extraction, marketing fluff detection, 5Ws validation, quote count optimization

**LLM Integration (internal/llm/llm.go)**:
- `AnalyzeSection()` provides qualitative feedback only (no score influence)
- Exponential backoff with jitter for API resilience
- Gracefully handles missing API keys

**UI Components (internal/ui/)**:
- `components.go`: Renders score breakdowns, strengths, improvements, quote analysis
- `model.go`: Bubble Tea model with tabbed interface (Overview, Breakdown, Quotes, AI Feedback)
- `styles.go`: Lip Gloss styling with color-coded scoring and progress bars

### Data Flow
1. CLI parses input file path and flags (`-no-tui`, `-report`)
2. Parser extracts structured sections and performs deterministic scoring
3. Interactive TUI displays results with tabbed navigation
4. Optional: LLM provides qualitative feedback (requires API key)
5. Results available in interactive UI, legacy stdout, or markdown report formats

## Dependencies

- **github.com/charmbracelet/bubbletea**: Terminal UI framework
- **github.com/charmbracelet/lipgloss**: Terminal styling
- **github.com/sashabaranov/go-openai**: OpenAI API client (optional)
- Standard library packages for file I/O, regex, and text analysis

## Recent Updates

### Recent Updates (November 2025)
- Fixed TUI alignment issues with emoji rendering
- Enhanced README with clearer structure
- Created high-quality example document (`example_prfaq_1.md`) scoring 77/100
- Added quote count feedback for documents with excessive quotes
- Implemented comprehensive Go linting infrastructure with golangci-lint
- Added pre-commit hooks for code quality enforcement
- Achieved 76%+ test coverage with comprehensive unit tests

### Architectural Decisions
- Deterministic scoring (100%) for consistency, optional AI for qualitative insights
- Interactive TUI as primary interface, with legacy stdout and markdown report modes
- Flexible document parsing supporting various header styles
- Modular structure with clear separation of concerns

## Testing

Sample documents demonstrate scoring range:
- `example_prfaq_4.md`: ~51/100 (typical first draft)
- `example_prfaq_3.md`: ~38/100 (needs significant improvement)
- `example_prfaq_1.md`: ~77/100 (high-quality example with metrics-rich quotes)

The validator is intentionally demanding - scores above 80 are rare and indicate publication-ready quality.

## Coding Conventions (MANDATORY)

**All code in this repository MUST follow the style guides without exception.**

### Go Code

See [docs/GO_STYLE_GUIDE.md](docs/GO_STYLE_GUIDE.md) for complete requirements.

**Key Requirements:**
- All code must pass `golangci-lint run` with zero errors
- Test coverage must be **≥80%** for core logic
- Functions must be **≤50 lines** (max 100)
- Files must be **≤400 lines**
- All exported functions must have doc comments
- Error wrapping with `fmt.Errorf("context: %w", err)`
- Table-driven tests for all test cases

**Pre-commit checks:**
```bash
golangci-lint run
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

### Python Code

See [docs/PYTHON_STYLE_GUIDE.md](docs/PYTHON_STYLE_GUIDE.md) for complete requirements.

**Key Requirements:**
- All code must pass `pylint` with score **≥9.5/10**
- All code must pass `mypy --ignore-missing-imports`
- Test coverage must be **≥50%** overall
- Functions must be **≤50 lines** (max 100)
- Files must be **≤400 lines**
- All functions must have type annotations
- All public functions must have docstrings (Google style)
- Line length: **120 characters**

**Pre-commit checks:**
```bash
black --line-length=120 --check scripts/
isort --profile black --line-length 120 --check-only scripts/
pylint scripts/ --fail-under=9.5
mypy scripts/ --ignore-missing-imports
pytest tests/python/ --cov=scripts --cov-fail-under=50
```

### CI Quality Gates

All PRs must pass:
1. **Go linting** - golangci-lint with zero errors
2. **Go tests** - All tests pass with race detection
3. **Go coverage** - ≥80% coverage on core packages
4. **Python linting** - pylint ≥9.5, mypy passes
5. **Python tests** - ≥50% coverage
6. **Security scan** - gosec with no high-severity issues

### Enforcement

- Pre-commit hooks enforce formatting and linting
- CI blocks merges that fail quality gates
- Coverage reports uploaded to Codecov
- No exceptions without explicit approval and documented justification
