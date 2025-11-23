# PR-FAQ Validator Prompts

This directory contains all LLM prompts used by the pr-faq-validator tool, organized in YAML format for easy maintenance and optimization.

## Structure

```
prompts/
├── README.md                    # This file
├── analysis/
│   └── section_review.yaml      # Prompt for analyzing PR-FAQ sections
└── generation/
    └── pr_faq_generation.yaml   # Prompt for generating PR-FAQs (used by prompt tuning)
```

## Prompt Format

Each YAML file follows this structure:

```yaml
name: "prompt-name"
version: "1.0.0"
description: "Brief description of what this prompt does"

context: |
  Detailed context about when and how this prompt is used.

system_prompt: |
  System-level instructions that set the LLM's role and constraints.
  Can include {{variable}} placeholders for dynamic content.

user_prompt_template: |
  The actual request with {{variable}} substitution.
  Uses Jinja2 template syntax for flexibility.

parameters:
  temperature: 0.7
  max_tokens: 2000

quality_criteria:
  - "Criterion 1"
  - "Criterion 2"
```

## Variables

Prompts use Jinja2 template syntax for variable substitution:
- `{{variable}}` - Simple variable substitution
- `{% if condition %}...{% endif %}` - Conditional blocks
- `{% for item in list %}...{% endfor %}` - Loops

## Versioning

Each prompt has a version number following semantic versioning:
- **Major**: Breaking changes to prompt structure or expected output
- **Minor**: Improvements or additions that don't break existing usage
- **Patch**: Bug fixes or minor wording improvements

## Optimization

Prompts in this directory can be optimized using the prompt tuning tool:

```bash
python prompt_tuning_tool.py evolve pr-faq-validator
```

The evolutionary optimization process will:
1. Load current prompts from this directory
2. Generate test cases and evaluate quality
3. Iteratively improve prompts using LLM-based mutations
4. Save improved prompts back to this directory

## Best Practices

1. **Keep prompts focused**: Each prompt should have a single, clear purpose
2. **Document variables**: List all expected variables in the context section
3. **Include examples**: Add example inputs/outputs in comments
4. **Version carefully**: Update version numbers when making changes
5. **Test thoroughly**: Run prompt tuning tests after modifications

