"""LLM client with mock support for PR-FAQ Validator prompt tuning."""

import os
import json
import hashlib
import time
import asyncio
from pathlib import Path
from typing import Optional, Dict, Any
from abc import ABC, abstractmethod


class LLMClient(ABC):
    """Abstract base class for LLM clients."""
    
    @abstractmethod
    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate text from prompt."""
        pass


class MockLLMClient(LLMClient):
    """Mock LLM client for testing without API keys."""
    
    def __init__(self, model: str = "mock"):
        self.model = model
        self.call_count = 0
    
    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate deterministic mock response based on prompt hash."""
        self.call_count += 1
        
        # Create deterministic response based on prompt content
        prompt_hash = hashlib.md5(prompt.encode()).hexdigest()[:8]
        
        # Detect what kind of response is needed based on prompt keywords
        if "press release" in prompt.lower():
            return self._generate_mock_press_release(prompt_hash)
        elif "faq" in prompt.lower() or "frequently asked" in prompt.lower():
            return self._generate_mock_faq(prompt_hash)
        elif "evaluate" in prompt.lower() or "score" in prompt.lower():
            return self._generate_mock_evaluation(prompt_hash)
        else:
            return self._generate_mock_generic(prompt_hash)
    
    def _generate_mock_press_release(self, seed: str) -> str:
        """Generate mock press release."""
        return f"""# Mock Press Release ({seed})

**FOR IMMEDIATE RELEASE**

## Revolutionary Product Transforms Industry

SEATTLE, WA - Today marks a significant milestone in product development with the launch of our innovative solution. This groundbreaking approach addresses key customer pain points while delivering exceptional value.

"This represents a fundamental shift in how we approach the problem," said Mock Executive. "Our customers have been asking for this capability, and we're excited to deliver it."

The solution provides three core benefits:
- Improved efficiency through automation
- Enhanced user experience with intuitive design
- Reduced costs through optimized processes

Early customer feedback has been overwhelmingly positive, with beta testers reporting significant improvements in their workflows.

For more information, visit our website or contact our press team.

### About the Company
We are committed to innovation and customer success.
"""
    
    def _generate_mock_faq(self, seed: str) -> str:
        """Generate mock FAQ."""
        return f"""# Frequently Asked Questions ({seed})

## Q: What problem does this solve?
A: This solution addresses the core challenge of inefficiency in current workflows by providing automated, intelligent assistance.

## Q: Who is this for?
A: This is designed for teams and individuals who need to streamline their processes and improve productivity.

## Q: How does it work?
A: The system uses advanced algorithms to analyze inputs and generate optimized outputs based on best practices.

## Q: What are the key benefits?
A: Users can expect faster turnaround times, higher quality results, and reduced manual effort.

## Q: When will this be available?
A: The solution is currently in beta and will be generally available in Q2.

## Q: How much does it cost?
A: Pricing will be announced closer to general availability, with options for different team sizes.

## Q: What support is available?
A: Comprehensive documentation, tutorials, and dedicated support channels will be provided.
"""
    
    def _generate_mock_evaluation(self, seed: str) -> str:
        """Generate mock evaluation score."""
        # Use seed to generate consistent but varied scores
        score_base = int(seed[:2], 16) % 30 + 70  # Score between 70-100
        
        return json.dumps({
            "overall_score": score_base,
            "press_release_quality": score_base + 5,
            "faq_completeness": score_base - 3,
            "clarity_score": score_base + 2,
            "structure_adherence": score_base,
            "feedback": f"Mock evaluation with seed {seed}. Good structure and clarity.",
            "strengths": [
                "Clear problem statement",
                "Well-structured content",
                "Comprehensive FAQ coverage"
            ],
            "improvements": [
                "Add more specific metrics",
                "Include customer quotes",
                "Expand on technical details"
            ]
        }, indent=2)
    
    def _generate_mock_generic(self, seed: str) -> str:
        """Generate generic mock response."""
        return f"Mock LLM response (seed: {seed}). This is a simulated output for testing purposes."


# Global request counter shared across all AssistantLLMClient instances
# to prevent request ID collisions
_global_request_counter = 0

class AssistantLLMClient(LLMClient):
    """
    File-based LLM client for AI assistant communication.

    Enables AI assistant (Claude) to act as the LLM without API keys.
    Writes requests to .pr-faq-validator/llm_requests/ and polls for
    responses in .pr-faq-validator/llm_responses/.

    Based on bloginator's AssistantLLMClient implementation.
    """

    def __init__(self, model: str = "assistant-llm", timeout: int = 300):
        self.model = model
        self.timeout = timeout

        # Create request/response directories
        self.base_dir = Path.cwd() / ".pr-faq-validator"
        self.request_dir = self.base_dir / "llm_requests"
        self.response_dir = self.base_dir / "llm_responses"

        self.request_dir.mkdir(parents=True, exist_ok=True)
        self.response_dir.mkdir(parents=True, exist_ok=True)

    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate text via file-based communication with AI assistant."""
        global _global_request_counter
        _global_request_counter += 1
        request_id = _global_request_counter

        # Prepare request
        request_data = {
            "request_id": request_id,
            "model": self.model,
            "temperature": kwargs.get("temperature", 0.3),
            "max_tokens": kwargs.get("max_tokens", 3000),
            "system_prompt": kwargs.get("system_prompt"),
            "prompt": prompt,
            "timestamp": time.time()
        }

        # Write request file
        request_file = self.request_dir / f"request_{request_id:04d}.json"
        with open(request_file, 'w') as f:
            json.dump(request_data, f, indent=2)

        # Poll for response
        response_file = self.response_dir / f"response_{request_id:04d}.json"
        start_time = time.time()

        while time.time() - start_time < self.timeout:
            if response_file.exists():
                with open(response_file, 'r') as f:
                    response_data = json.load(f)
                return response_data["content"]

            await asyncio.sleep(0.5)  # Poll every 500ms

        raise TimeoutError(
            f"No response received for request {request_id} after {self.timeout}s. "
            f"Expected response file: {response_file}"
        )


class AnthropicClient(LLMClient):
    """Anthropic Claude client."""

    def __init__(self, model: str = "claude-3-5-sonnet-20241022", api_key: Optional[str] = None):
        self.model = model
        self.api_key = api_key or os.getenv("ANTHROPIC_API_KEY")

        if not self.api_key:
            raise ValueError("ANTHROPIC_API_KEY not found in environment")

    async def generate(self, prompt: str, **kwargs) -> str:
        """Generate text using Anthropic API."""
        try:
            import anthropic
        except ImportError:
            raise ImportError("anthropic package not installed. Run: pip install anthropic")

        client = anthropic.AsyncAnthropic(api_key=self.api_key)

        temperature = kwargs.get("temperature", 1.0)
        max_tokens = kwargs.get("max_tokens", 4096)

        message = await client.messages.create(
            model=self.model,
            max_tokens=max_tokens,
            temperature=temperature,
            messages=[{"role": "user", "content": prompt}]
        )

        return message.content[0].text


def create_llm_client(provider: str = "anthropic", model: str = "claude-3-5-sonnet-20241022", mock: bool = False) -> LLMClient:
    """Factory function to create LLM client."""
    if mock or os.getenv("AI_AGENT_MOCK_MODE", "false").lower() == "true":
        return MockLLMClient(model=model)

    # Check for assistant mode (file-based communication)
    if provider == "assistant" or os.getenv("LLM_PROVIDER") == "assistant":
        return AssistantLLMClient(model=model)

    if provider == "anthropic":
        return AnthropicClient(model=model)
    else:
        raise ValueError(f"Unsupported LLM provider: {provider}")

