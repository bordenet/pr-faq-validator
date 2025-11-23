#!/usr/bin/env python3
"""
PR-FAQ Validator Prompt Tuning Tool - Main Entry Point

This tool implements the evolutionary prompt tuning methodology for the
pr-faq-validator project. It provides automated LLM prompt optimization
with mutation-based keep/discard logic.

Usage:
    python prompt_tuning_tool.py init pr-faq-validator
    python prompt_tuning_tool.py simulate pr-faq-validator [--mock]
    python prompt_tuning_tool.py evolve pr-faq-validator [--mock]
    python prompt_tuning_tool.py ai-agent-optimize pr-faq-validator
    python prompt_tuning_tool.py status pr-faq-validator

Example:
    # Initialize new project
    python prompt_tuning_tool.py init pr-faq-validator

    # Run evolutionary tuning (with API keys)
    python prompt_tuning_tool.py evolve pr-faq-validator

    # AI Agent autonomous optimization (no API keys required)
    python prompt_tuning_tool.py ai-agent-optimize pr-faq-validator
"""

import sys
import os
from pathlib import Path

# Add scripts directory to Python path
scripts_dir = Path(__file__).parent / "scripts" / "prompt_tuning"
sys.path.insert(0, str(scripts_dir))

# Import and run CLI
from prompt_tuning_cli import cli

if __name__ == '__main__':
    # Load environment variables from .env file if it exists
    env_file = Path('.env')
    if env_file.exists():
        try:
            from dotenv import load_dotenv
            load_dotenv()
        except ImportError:
            # dotenv is optional - continue without it
            pass

    cli()

