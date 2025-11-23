"""Prompt simulator for PR-FAQ Validator."""

import json
import asyncio
from pathlib import Path
from typing import Dict, List, Any
from datetime import datetime

from prompt_tuning_config import PromptTuningConfig, load_test_cases, load_prompts
from llm_client import create_llm_client


class PromptSimulator:
    """Simulates PR-FAQ generation using current prompts."""
    
    def __init__(self, config: PromptTuningConfig):
        self.config = config
        self.llm_client = create_llm_client(
            provider=config.llm_provider,
            model=config.llm_model,
            mock=config.mock_mode
        )
    
    async def run_simulation(self, iteration: int = 0) -> Dict[str, Any]:
        """Run simulation for all test cases."""
        test_cases_data = load_test_cases(self.config)
        test_cases = test_cases_data.get("test_cases", [])
        prompts = load_prompts(self.config)
        
        results = {
            "iteration": iteration,
            "timestamp": datetime.now().isoformat(),
            "project": self.config.project_name,
            "test_case_results": []
        }
        
        for test_case in test_cases:
            result = await self._run_test_case(test_case, prompts, iteration)
            results["test_case_results"].append(result)
        
        return results
    
    async def _run_test_case(self, test_case: Dict[str, Any], prompts: Dict[str, str], iteration: int) -> Dict[str, Any]:
        """Run a single test case."""
        test_case_id = test_case.get("id", "unknown")
        inputs = test_case.get("inputs", {})
        
        # Generate PR-FAQ using prompts
        pr_faq_content = await self._generate_pr_faq(inputs, prompts)
        
        result = {
            "test_case_id": test_case_id,
            "test_case_name": test_case.get("name", ""),
            "iteration": iteration,
            "timestamp": datetime.now().isoformat(),
            "inputs": inputs,
            "generated_content": pr_faq_content,
            "metadata": {
                "industry": test_case.get("industry", ""),
                "project_type": test_case.get("project_type", ""),
                "scope": test_case.get("scope", "")
            }
        }
        
        return result
    
    async def _generate_pr_faq(self, inputs: Dict[str, Any], prompts: Dict[str, str]) -> Dict[str, str]:
        """Generate PR-FAQ content using LLM."""
        # Get the main PR-FAQ generation prompt
        pr_faq_prompt = prompts.get("pr_faq_generation", self._get_default_prompt())
        
        # Format prompt with inputs
        formatted_prompt = self._format_prompt(pr_faq_prompt, inputs)
        
        # Generate content
        content = await self.llm_client.generate(formatted_prompt, temperature=self.config.temperature)
        
        # Parse into sections
        return {
            "full_content": content,
            "press_release": self._extract_section(content, "press release"),
            "faq": self._extract_section(content, "faq")
        }
    
    def _format_prompt(self, prompt: str, inputs: Dict[str, Any]) -> str:
        """Format prompt template with inputs."""
        formatted = prompt
        for key, value in inputs.items():
            placeholder = f"{{{key}}}"
            formatted = formatted.replace(placeholder, str(value))
        return formatted
    
    def _extract_section(self, content: str, section_name: str) -> str:
        """Extract a specific section from generated content."""
        # Simple extraction - look for section headers
        lines = content.split('\n')
        section_lines = []
        in_section = False
        
        for line in lines:
            if section_name.lower() in line.lower() and (line.startswith('#') or line.startswith('##')):
                in_section = True
                continue
            elif in_section and line.startswith('#'):
                break
            elif in_section:
                section_lines.append(line)
        
        return '\n'.join(section_lines).strip()
    
    def _get_default_prompt(self) -> str:
        """Get default PR-FAQ generation prompt."""
        return """Generate a comprehensive PR-FAQ document for the following project:

Project Name: {projectName}
Problem Description: {problemDescription}
Business Context: {businessContext}

Please create:
1. A compelling press release that follows Amazon's PR-FAQ format
2. A comprehensive FAQ section addressing key stakeholder questions

The press release should:
- Start with a clear headline
- Include a compelling opening paragraph
- Explain the customer problem and solution
- Include a quote from leadership
- Describe key benefits and features

The FAQ should address:
- What problem does this solve?
- Who is this for?
- How does it work?
- What are the key benefits?
- When will this be available?
- How much does it cost?
- What support is available?
"""
    
    def save_results(self, results: Dict[str, Any], iteration: int) -> Path:
        """Save simulation results to file."""
        self.config.results_dir.mkdir(parents=True, exist_ok=True)
        
        output_file = self.config.results_dir / f"simulation_iteration_{iteration:03d}.json"
        
        with open(output_file, 'w') as f:
            json.dump(results, f, indent=2)
        
        return output_file

