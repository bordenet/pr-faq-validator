#!/usr/bin/env python3
"""
Auto-responder for file-based LLM communication.

Monitors .pr-faq-validator/llm_requests/ for new requests and generates
appropriate responses in .pr-faq-validator/llm_responses/.

Based on bloginator's auto_respond_llm.py implementation.
"""

import json
import time
import random
import hashlib
from pathlib import Path
from typing import Dict, Any


class AutoResponder:
    """Autonomous LLM response generator for prompt optimization experiments."""

    def __init__(self, base_dir: Path = None):
        self.base_dir = base_dir or Path.cwd() / ".pr-faq-validator"
        self.request_dir = self.base_dir / "llm_requests"
        self.response_dir = self.base_dir / "llm_responses"
        self.processed_requests = set()

        # Ensure directories exist
        self.request_dir.mkdir(parents=True, exist_ok=True)
        self.response_dir.mkdir(parents=True, exist_ok=True)

    def detect_request_type(self, request_data: Dict[str, Any]) -> str:
        """Detect the type of request based on prompt content."""
        prompt = request_data.get("prompt", "").lower()
        system_prompt = (request_data.get("system_prompt") or "").lower()

        # Evaluation requests contain specific keywords
        if any(keyword in prompt for keyword in [
            "evaluate", "score", "quality", "assessment",
            "press_release_quality", "faq_completeness", "clarity_score"
        ]):
            return "evaluation"

        # Press release generation
        if "press release" in prompt or "press release" in system_prompt:
            return "press_release"

        # FAQ generation
        if "faq" in prompt or "frequently asked" in prompt:
            return "faq"

        return "generic"

    def generate_press_release(self, request_data: Dict[str, Any]) -> str:
        """Generate mock press release with variation."""
        request_id = request_data.get("request_id", 0)
        seed = hashlib.md5(str(request_id).encode()).hexdigest()[:8]

        return f"""# Revolutionary Product Launch

**FOR IMMEDIATE RELEASE**

## New Solution Transforms Industry Landscape

SEATTLE, WA - Today marks a significant milestone with the launch of an innovative solution that addresses critical customer challenges. This groundbreaking approach delivers measurable value through automation and intelligent design.

"This represents a fundamental shift in how teams approach complex problems," said Product Leader. "Our customers have been requesting this capability, and we're delivering it with precision and scale."

The solution provides three core benefits:
- **40% improvement** in workflow efficiency through intelligent automation
- **Enhanced user experience** with intuitive, data-driven design
- **25% cost reduction** through optimized resource allocation

Early customer feedback has been overwhelmingly positive. Beta testers report significant improvements in productivity and quality metrics.

### Key Features

- Real-time analytics and insights
- Seamless integration with existing tools
- Enterprise-grade security and compliance
- Scalable architecture for teams of all sizes

### Availability

The solution enters general availability in Q2 2025, with early access programs available now.

For more information, visit our website or contact our team.

### About the Company

We are committed to innovation, customer success, and delivering measurable business value.

**Contact:** press@company.com

(Seed: {seed})
"""

    def generate_faq(self, request_data: Dict[str, Any]) -> str:
        """Generate mock FAQ with variation."""
        request_id = request_data.get("request_id", 0)
        seed = hashlib.md5(str(request_id).encode()).hexdigest()[:8]

        return f"""# Frequently Asked Questions

## Q: What specific problem does this solve?

A: This solution addresses the core challenge of inefficiency in current workflows by providing automated, intelligent assistance. Teams spend an average of 15 hours per week on manual tasks that can be automated, and this solution reduces that by 60%.

## Q: Who is this designed for?

A: This is built for engineering teams, product managers, and technical leaders who need to streamline their processes and improve productivity. It's particularly valuable for teams managing complex projects with multiple stakeholders.

## Q: How does the system work?

A: The platform uses advanced algorithms to analyze inputs, identify patterns, and generate optimized outputs based on industry best practices. It integrates seamlessly with existing tools via APIs and webhooks.

## Q: What are the measurable benefits?

A: Users can expect:
- 40% faster turnaround times on key deliverables
- 30% higher quality scores in peer reviews
- 50% reduction in manual rework
- 25% improvement in team satisfaction scores

## Q: When will this be available?

A: The solution is currently in beta with select customers and will be generally available in Q2 2025. Early access programs are available for qualified teams.

## Q: What is the pricing model?

A: Pricing starts at $99/user/month for teams of 10+, with volume discounts available. Enterprise pricing includes dedicated support and custom integrations.

## Q: What support and training is provided?

A: Comprehensive documentation, video tutorials, live training sessions, and dedicated support channels are included. Enterprise customers receive a dedicated customer success manager.

## Q: How does this integrate with our existing tools?

A: The platform offers native integrations with popular tools (Jira, GitHub, Slack, etc.) and a robust REST API for custom integrations.

(Seed: {seed})
"""

    def generate_evaluation(self, request_data: Dict[str, Any]) -> str:
        """Generate mock evaluation with floating-point scores and variation."""
        request_id = request_data.get("request_id", 0)

        # Use request_id to generate varied but consistent scores
        random.seed(request_id)

        # Generate scores with realistic variation (4.0-4.8 range for good content)
        base_score = 4.0 + random.random() * 0.8

        evaluation = {
            "overall_score": round(base_score, 2),
            "press_release_quality": round(base_score + random.uniform(-0.2, 0.3), 2),
            "faq_completeness": round(base_score + random.uniform(-0.3, 0.2), 2),
            "clarity_score": round(base_score + random.uniform(-0.2, 0.2), 2),
            "structure_adherence": round(base_score + random.uniform(-0.1, 0.2), 2),
            "content_quality": {
                "clarity": round(base_score + random.uniform(-0.2, 0.2), 2),
                "depth": round(base_score + random.uniform(-0.1, 0.3), 2),
                "nuance": round(base_score + random.uniform(-0.3, 0.2), 2),
                "specificity": round(base_score + random.uniform(-0.4, 0.1), 2)
            },
            "slop_violations": {
                "critical": [],
                "high": [],
                "medium": [],
                "low": []
            },
            "voice_analysis": {
                "authenticity_score": round(base_score, 2),
                "strengths": [
                    "Direct, concrete language",
                    "Specific metrics and examples",
                    "Clear problem-solution framing"
                ],
                "concerns": []
            },
            "feedback": f"Score: {base_score:.2f}/5.0. Content demonstrates good clarity and practical focus. "
                       f"Could benefit from more specific examples and quantitative data.",
            "strengths": [
                "Clear problem statement",
                "Well-structured content",
                "Comprehensive FAQ coverage",
                "Specific metrics included"
            ],
            "improvements": [
                "Add more concrete customer quotes",
                "Include specific technical details",
                "Expand on implementation timeline",
                "Provide more quantitative evidence"
            ],
            "evolutionary_strategy": {
                "prompt_to_modify": random.choice(["press_release", "faq", "refinement"]),
                "specific_changes": [
                    {
                        "section": "introduction",
                        "current_issue": "Lacks specific quantitative metrics",
                        "proposed_change": "Add concrete numbers and percentages to claims",
                        "rationale": "Specificity scores are lowest dimension"
                    },
                    {
                        "section": "benefits",
                        "current_issue": "Generic benefit statements",
                        "proposed_change": "Include customer quotes with measurable outcomes",
                        "rationale": "Voice authenticity needs concrete evidence"
                    }
                ],
                "priority": "high",
                "expected_impact": "Improve specificity score by 0.3-0.5 points"
            }
        }

        return json.dumps(evaluation, indent=2)

    def generate_response(self, request_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate appropriate response based on request type."""
        request_type = self.detect_request_type(request_data)

        if request_type == "evaluation":
            content = self.generate_evaluation(request_data)
        elif request_type == "press_release":
            content = self.generate_press_release(request_data)
        elif request_type == "faq":
            content = self.generate_faq(request_data)
        else:
            content = f"Generic response for request {request_data.get('request_id')}"

        # Estimate token counts
        prompt_tokens = len(request_data.get("prompt", "").split())
        completion_tokens = len(content.split())

        return {
            "content": content,
            "prompt_tokens": prompt_tokens,
            "completion_tokens": completion_tokens,
            "finish_reason": "stop"
        }

    def process_requests(self, continuous: bool = False, sleep_interval: float = 1.0):
        """Process pending requests."""
        print(f"Auto-responder monitoring: {self.request_dir}")
        print(f"Responses will be written to: {self.response_dir}")

        iteration = 0
        while True:
            iteration += 1

            # Find unprocessed requests
            request_files = sorted(self.request_dir.glob("request_*.json"))
            new_requests = [f for f in request_files if f.name not in self.processed_requests]

            if new_requests:
                print(f"\n[Iteration {iteration}] Found {len(new_requests)} new request(s)")

            for request_file in new_requests:
                try:
                    # Read request
                    with open(request_file, 'r') as f:
                        request_data = json.load(f)

                    request_id = request_data.get("request_id")
                    request_type = self.detect_request_type(request_data)

                    print(f"  Processing request {request_id} (type: {request_type})")

                    # Generate response
                    response_data = self.generate_response(request_data)

                    # Write response
                    response_file = self.response_dir / f"response_{request_id:04d}.json"
                    with open(response_file, 'w') as f:
                        json.dump(response_data, f, indent=2)

                    print(f"  ✓ Response written: {response_file.name}")

                    # Mark as processed
                    self.processed_requests.add(request_file.name)

                except Exception as e:
                    print(f"  ✗ Error processing {request_file.name}: {e}")

            if not continuous:
                break

            time.sleep(sleep_interval)


def main():
    """Main entry point."""
    import argparse

    parser = argparse.ArgumentParser(description="Auto-responder for LLM requests")
    parser.add_argument("--continuous", action="store_true", help="Run continuously")
    parser.add_argument("--interval", type=float, default=1.0, help="Sleep interval in seconds")
    parser.add_argument("--base-dir", type=Path, default=None, help="Base directory")

    args = parser.parse_args()

    responder = AutoResponder(base_dir=args.base_dir)

    try:
        responder.process_requests(continuous=args.continuous, sleep_interval=args.interval)
    except KeyboardInterrupt:
        print("\n\nAuto-responder stopped by user")


if __name__ == "__main__":
    main()


