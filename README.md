
# pr-faq-validator

[![CI](https://github.com/bordenet/pr-faq-validator/actions/workflows/ci.yml/badge.svg)](https://github.com/bordenet/pr-faq-validator/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/bordenet/pr-faq-validator/branch/main/graph/badge.svg)](https://codecov.io/gh/bordenet/pr-faq-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/bordenet/pr-faq-validator)](https://goreportcard.com/report/github.com/bordenet/pr-faq-validator)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bordenet/pr-faq-validator)](https://github.com/bordenet/pr-faq-validator/blob/main/go.mod)
[![License](https://img.shields.io/github/license/bordenet/pr-faq-validator)](https://github.com/bordenet/pr-faq-validator/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/bordenet/pr-faq-validator)](https://github.com/bordenet/pr-faq-validator/releases)
[![Dependabot](https://img.shields.io/badge/Dependabot-enabled-brightgreen?logo=dependabot)](https://github.com/bordenet/pr-faq-validator/security/dependabot)

A CLI tool to analyze and score [PR-FAQ documents](https://github.com/bordenet/Engineering_Culture/blob/main/SDLC/The_PR-FAQ.md) using deterministic quality metrics.

## Overview

pr-faq-validator analyzes PR-FAQ (Press Release - Frequently Asked Questions) documents and provides quality scoring with detailed feedback. The tool detects press release and FAQ sections and evaluates content against journalistic standards.

## Features

- 100-point scoring system across 4 dimensions
- Automatic section detection for flexible document structures
- Press release evaluation against journalistic standards
- Quote metric analysis - identifies quantitative data in testimonials
- Marketing language detection - flags hyperbolic claims
- 5 Ws validation - ensures WHO, WHAT, WHEN, WHERE, WHY coverage
- Interactive terminal UI with detailed breakdowns
- Optional AI feedback via OpenAI API

## Installation

```bash
git clone https://github.com/bordenet/pr-faq-validator.git
cd pr-faq-validator
go mod tidy
go build
```

**Requirements:** Go 1.21+, OpenAI API key (optional, for AI feedback)

## Usage

```bash
# Set API key for AI feedback (optional)
export OPENAI_API_KEY=your_openai_api_key_here

# Analyze a document
./pr-faq-validator -file path/to/your/prfaq.md
```

### Examples

Analyze any of the included sample documents:

```bash
./pr-faq-validator -file testdata/example_prfaq_1.md
./pr-faq-validator -file testdata/example_prfaq_2.txt  
./pr-faq-validator -file testdata/example_prfaq_3.md
./pr-faq-validator -file testdata/example_prfaq_4.md
```

## Input Format

Works with any document structure. Recommended format:

```markdown
# Your PR-FAQ Title

## Press Release
Your press release content...

## FAQ
Q: Question?
A: Answer...
```

Automatically detects sections regardless of headers ("Press Release", "Announcement", "Q&A", etc.).

## Output

Provides interactive terminal UI with:

- Score breakdown across 4 categories (Structure, Content, Professional, Evidence)
- Strengths and improvements with specific recommendations
- Quote analysis with individual scoring and metric detection
- AI feedback for detailed insights (requires OpenAI API key)

## Scoring Methodology

**Deterministic Scoring (100% of numerical score):** Rule-based algorithms analyze text patterns for consistent results. AI does not influence scores.

**AI Feedback (qualitative only):** Optional GPT-4 insights for improvement suggestions.

**Scoring Breakdown:**

- **Structure & Hook (30 pts):** Headline quality, newsworthy hook, release date
- **Content Quality (35 pts):** 5 Ws coverage, credibility, structure
- **Professional Quality (20 pts):** Tone, readability, marketing language detection
- **Customer Evidence (15 pts):** Quote quality with quantitative metrics

**Note:** Scoring is strict - high scores require well-crafted documents. Focus on actionable feedback to improve quality.

---

## Code Coverage

pr-faq-validator maintains **76.3% test coverage** with comprehensive testing of core functionality. The coverage visualization below shows detailed coverage by module:

[![Coverage Grid](https://codecov.io/gh/bordenet/pr-faq-validator/graphs/tree.svg)](https://codecov.io/gh/bordenet/pr-faq-validator)

**What this means:**

- **Green**: Well-tested code (>80% coverage)
- **Yellow**: Moderate coverage (60-80%)
- **Red**: Needs more tests (<60%)
- **Size**: Larger boxes = more lines of code

Click the image to explore detailed coverage reports on Codecov, including line-by-line coverage, branch coverage, and historical trends.

---

## License

MIT License - see LICENSE file for details.
