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

        print("="*80, flush=True)
        print("AUTO-RESPONDER VERSION 4.0 - Prompt-Aware Content Generation", flush=True)
        print("="*80, flush=True)

    def detect_request_type(self, request_data: Dict[str, Any]) -> str:
        """Detect the type of request based on prompt content."""
        print(f"    [DEBUG-V4] detect_request_type called - VERSION 4.0", flush=True)
        prompt = request_data.get("prompt", "")
        system_prompt = (request_data.get("system_prompt") or "")
        request_id = request_data.get('request_id', 'unknown')

        print(f"    [DEBUG-V4] Prompt first 150 chars: {repr(prompt[:150])}", flush=True)
        print(f"    [DEBUG-V4] Has 'Current prompt:': {'Current prompt:' in prompt}", flush=True)
        print(f"    [DEBUG-V4] Has 'current prompt:': {'current prompt:' in prompt}", flush=True)

        # For mutation detection, only check the INSTRUCTION part (before "Current prompt:" or "current prompt:")
        # to avoid false matches in the embedded prompt content
        if "Current prompt:" in prompt:
            instruction_part = prompt.split("Current prompt:")[0]
            print(f"    [DEBUG-V4] Split on 'Current prompt:' - instruction length: {len(instruction_part)}", flush=True)
        elif "current prompt:" in prompt:
            instruction_part = prompt.split("current prompt:")[0]
            print(f"    [DEBUG-V4] Split on 'current prompt:' - instruction length: {len(instruction_part)}", flush=True)
        else:
            instruction_part = prompt
            print(f"    [DEBUG-V4] No split - using full prompt", flush=True)

        instruction_lower = instruction_part.lower()
        prompt_lower = prompt.lower()
        system_lower = system_prompt.lower()

        print(f"    [DEBUG] Request {request_id}: Prompt length = {len(prompt)}, instruction length = {len(instruction_part)}", flush=True)

        # Mutation/optimization requests (check FIRST before evaluation)
        # These ask for improved prompts, not evaluations
        mutation_keywords = [
            "suggest an improved version",
            "optimizing llm prompts",
            "improved prompt text",
            "provide only the improved prompt",
            "expert at optimizing"
        ]
        found_mutation = [k for k in mutation_keywords if k in instruction_lower]
        print(f"    [DEBUG] Request {request_id}: Checking mutation keywords in instruction... found: {found_mutation}", flush=True)
        if found_mutation:
            print(f"    [DEBUG] Request {request_id}: RETURNING mutation", flush=True)
            return "mutation"

        # Evaluation requests contain specific keywords
        eval_keywords = [
            "evaluate", "score", "assessment",
            "press_release_quality", "faq_completeness", "clarity_score"
        ]
        found_eval = [k for k in eval_keywords if k in instruction_lower]
        if found_eval:
            print(f"    [DEBUG] Request {request_id}: Found evaluation keywords in instruction: {found_eval[:2]}", flush=True)
            return "evaluation"

        # Press release generation (check instruction part only)
        if "press release" in instruction_lower or "press release" in system_lower:
            return "press_release"

        # FAQ generation (check instruction part only)
        if "faq" in instruction_lower or "frequently asked" in instruction_lower:
            return "faq"

        return "generic"

    def generate_press_release(self, request_data: Dict[str, Any]) -> str:
        """Generate varied press release based on prompt requirements."""
        request_id = request_data.get("request_id", 0)
        prompt = request_data.get("prompt", "")
        prompt_lower = prompt.lower()
        seed = hashlib.md5(str(request_id).encode()).hexdigest()[:8]

        # Use prompt hash to create deterministic but varied outputs for different prompts
        # This ensures that improved prompts generate different (hopefully better) content
        prompt_hash = int(hashlib.md5(prompt.encode()).hexdigest()[:16], 16)
        quality_score = (prompt_hash % 100) / 100.0  # 0.00 to 0.99

        # Extract key requirements from prompt to vary the output
        has_specificity = "specific" in prompt_lower or "quantitative" in prompt_lower or "measurable" in prompt_lower
        has_evidence = "evidence" in prompt_lower or "credibility" in prompt_lower or "validation" in prompt_lower
        has_customer_focus = "customer" in prompt_lower or "stakeholder" in prompt_lower
        has_metrics = "metric" in prompt_lower or "number" in prompt_lower or "percentage" in prompt_lower

        # Boost quality score based on prompt requirements
        if has_evidence:
            quality_score += 0.15
        if has_specificity:
            quality_score += 0.10
        if has_metrics:
            quality_score += 0.10
        if has_customer_focus:
            quality_score += 0.05

        # Cap at 1.0
        quality_score = min(quality_score, 1.0)

        print(f"    [DEBUG-V4-PR] Request {request_id}: quality_score={quality_score:.2f}, has_evidence={has_evidence}, has_specificity={has_specificity}, has_metrics={has_metrics}")

        # Generate content based on quality score
        variation = int(quality_score * 10) % 3  # 0, 1, or 2

        # Generate content with quality proportional to quality_score
        if quality_score >= 0.8:
            # Highest quality PR with strong evidence and specific metrics
            if variation == 0:
                metrics = f"{int(45 + quality_score * 10)}% improvement"
                customer_quote = '"This solution reduced our processing time from 8 hours to 3.8 hours per day, saving our team $147,000 annually," said Sarah Chen, Operations Director at TechCorp. "The ROI was evident within the first quarter."'
                specifics = "The platform processes 18,000 transactions per hour with 99.98% accuracy, handling peak loads of 65,000 concurrent users across 12 geographic regions."
            elif variation == 1:
                metrics = f"{int(43 + quality_score * 10)}% improvement"
                customer_quote = '"We achieved a 4.2-hour reduction in daily processing time, translating to $127,000 in annual savings," said Michael Rodriguez, VP of Operations.'
                specifics = "The platform processes 15,000 transactions per hour with 99.97% accuracy, handling peak loads of 50,000 concurrent users."
            else:
                metrics = f"{int(42 + quality_score * 10)}% improvement"
                customer_quote = '"Processing time dropped from 8 to 4.5 hours daily, saving $115,000 per year," said Jennifer Park, Director of Engineering.'
                specifics = "The system handles 14,000 transactions per hour with 99.95% accuracy and supports 45,000 concurrent users."
        elif quality_score >= 0.6:
            # Good PR with metrics but less evidence
            metrics = f"{int(30 + quality_score * 15)}% improvement"
            customer_quote = '"We\'ve seen measurable productivity gains in our workflows," said Product Leader. "The metrics speak for themselves."'
            specifics = "The solution integrates with existing tools and provides real-time analytics with sub-second response times."
        elif quality_score >= 0.4:
            # Moderate PR with some customer focus
            metrics = f"{int(25 + quality_score * 15)}% improvement"
            customer_quote = '"This addresses what our customers have been asking for," said Customer Success Lead.'
            specifics = "Built based on feedback from customer interviews and beta testing."
        else:
            # Basic PR with minimal specifics
            metrics = f"{int(20 + quality_score * 20)}% improvement"
            customer_quote = '"This represents progress in our capabilities," said Team Lead.'
            specifics = "The solution provides automation features."

        return f"""# Revolutionary Product Launch

**FOR IMMEDIATE RELEASE**

## New Solution Transforms Industry Landscape

SEATTLE, WA - Today marks a significant milestone with the launch of an innovative solution that addresses critical customer challenges. This groundbreaking approach delivers measurable value through automation and intelligent design.

{customer_quote}

The solution provides three core benefits:
- **{metrics}** in workflow efficiency through intelligent automation
- **Enhanced user experience** with intuitive, data-driven design
- **25% cost reduction** through optimized resource allocation

{specifics}

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
        """Generate varied FAQ based on prompt requirements."""
        request_id = request_data.get("request_id", 0)
        prompt = request_data.get("prompt", "").lower()
        seed = hashlib.md5(str(request_id).encode()).hexdigest()[:8]

        # Extract key requirements from prompt to vary the output
        has_comprehensive = "comprehensive" in prompt or "detailed" in prompt
        has_stakeholder = "stakeholder" in prompt or "concern" in prompt
        has_technical = "technical" in prompt or "implementation" in prompt

        # Build FAQ sections based on prompt emphasis
        technical_section = ""
        if has_technical:
            technical_section = """

## Q: What are the technical requirements?

A: The solution requires minimal infrastructure changes. It integrates via REST API with existing systems, supports OAuth 2.0 authentication, and runs on standard cloud platforms (AWS, Azure, GCP). Average implementation time is 2-3 weeks for most organizations.
"""

        stakeholder_section = ""
        if has_stakeholder:
            stakeholder_section = """

## Q: How does this address stakeholder concerns about cost?

A: The solution pays for itself within 6 months through efficiency gains. Organizations save an average of $150,000 annually through reduced manual processing time and error correction costs. We offer flexible pricing tiers to match different organization sizes.
"""

        comprehensive_detail = "15 hours per week" if has_comprehensive else "significant time"
        reduction_detail = "60%" if has_comprehensive else "substantially"

        return f"""# Frequently Asked Questions

## Q: What specific problem does this solve?

A: This solution addresses the core challenge of inefficiency in current workflows by providing automated, intelligent assistance. Teams spend an average of {comprehensive_detail} on manual tasks that can be automated, and this solution reduces that by {reduction_detail}.

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
{technical_section}{stakeholder_section}

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
        """Generate evaluation based on actual content quality."""
        request_id = request_data.get("request_id", 0)
        prompt = request_data.get("prompt", "")

        # Extract content quality indicators from the PR-FAQ content
        import re

        # Extract improvement percentage
        percentage_match = re.search(r'(\d+)% improvement', prompt)
        improvement_pct = int(percentage_match.group(1)) if percentage_match else 30

        # Check for quantitative metrics in customer quotes
        has_dollar_amount = bool(re.search(r'\$[\d,]+', prompt))
        has_time_savings = bool(re.search(r'\d+\.?\d* hours?', prompt))
        has_specific_metrics = bool(re.search(r'\d+,\d+ (transactions|users|concurrent)', prompt))
        has_accuracy = bool(re.search(r'\d+\.\d+% accuracy', prompt))

        # Calculate base score from content quality (3.5-4.8 range)
        # Higher improvement percentage = higher score
        base_score = 3.5 + (improvement_pct / 100.0) * 1.3  # 30% → 3.89, 50% → 4.15

        # Boost for quantitative evidence
        if has_dollar_amount:
            base_score += 0.15
        if has_time_savings:
            base_score += 0.10
        if has_specific_metrics:
            base_score += 0.15
        if has_accuracy:
            base_score += 0.10

        # Cap at 4.8
        base_score = min(base_score, 4.8)

        # Use request_id for minor variation to avoid identical scores
        random.seed(request_id)
        variation = random.uniform(-0.05, 0.05)
        base_score += variation

        print(f"    [DEBUG-V4-EVAL] Request {request_id}: improvement={improvement_pct}%, $={has_dollar_amount}, time={has_time_savings}, metrics={has_specific_metrics}, accuracy={has_accuracy}, base_score={base_score:.2f}")

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

    def generate_mutation(self, request_data: Dict[str, Any]) -> str:
        """Generate an improved/mutated version of a prompt."""
        request_id = request_data.get("request_id", 0)
        prompt = request_data.get("prompt", "")

        # Extract the current prompt from the mutation request
        # It's typically after "Current prompt:" in the request
        current_prompt_marker = "Current prompt:"
        if current_prompt_marker in prompt:
            parts = prompt.split(current_prompt_marker, 1)
            if len(parts) > 1:
                # Extract everything between "Current prompt:" and "Based on iteration"
                current_prompt_section = parts[1].split("Based on iteration")[0].strip()
            else:
                current_prompt_section = ""
        else:
            current_prompt_section = ""

        # Generate variations based on request_id for diversity
        variations = [
            # Variation 1: Add more specificity requirements
            """Generate a comprehensive PR-FAQ document following Amazon's format with strong emphasis on quantitative specificity and measurable outcomes.

Project: {projectName}
Problem: {problemDescription}
Context: {businessContext}

CRITICAL REQUIREMENTS:

1. QUANTITATIVE SPECIFICITY (MANDATORY)
   - Every claim must include concrete numbers, percentages, or metrics
   - Provide measurable outcomes (e.g., "reduces time from X to Y by Z%")
   - Reference specific data points (e.g., "68% cart abandonment rate")
   - Use data-driven examples with real numbers throughout
   - Avoid vague terms like "significant", "substantial", "many"

2. AUTHENTIC CUSTOMER EVIDENCE (REQUIRED)
   - Include 2-3 customer quotes with quantitative, measurable outcomes
   - Format: "Specific metric-driven quote" - Full Name, Title, Company
   - Example: "We reduced deployment time from 4 hours to 15 minutes, saving 20 hours per week" - Jane Smith, CTO, TechCorp
   - Each quote must contain at least one specific number or percentage

3. TECHNICAL DEPTH AND MECHANISM
   - Explain HOW the solution works (technical mechanism)
   - Include implementation details and architecture insights
   - Address technical concerns comprehensively in FAQ
   - Provide workflow diagrams or process descriptions
   - Specify technologies, integrations, or methodologies used

4. TIMELINE CLARITY AND AVAILABILITY
   - Specify exact availability dates or clear phases (e.g., "Q2 2025", "March 15, 2025")
   - Outline implementation timeline with milestones
   - Set clear expectations for rollout and adoption
   - Include beta/pilot program details if applicable

PRESS RELEASE STRUCTURE:
- Headline: Clear, compelling, captures essence (8-12 words)
- Opening: Who, what, when, where, why with SPECIFIC METRICS
- Problem: Customer pain with QUANTIFIED impact (numbers required)
- Solution: How it works with MEASURABLE benefits (percentages required)
- Leadership Quote: Authentic with strategic context and vision
- Key Benefits: 3-5 benefits, each with SPECIFIC NUMBERS
- Customer Quote: Real outcome with METRICS (required)
- Availability: Clear timeline with specific dates

FAQ REQUIREMENTS:
- Address in priority order: problem, audience, mechanism, benefits, availability, cost, support
- Provide SPECIFIC answers with numbers wherever applicable
- Include detailed technical implementation information
- Anticipate objections with data-driven, evidence-based responses
- Minimum 7 questions, maximum 15 questions
- Each answer should be 2-4 sentences with concrete details

QUALITY STANDARDS:
- Zero generic statements - be SPECIFIC in every sentence
- Every claim must have supporting data or metrics
- Use clear, precise, professional language
- No marketing jargon, hype, or superlatives
- Maintain authentic, credible tone throughout
- Ensure perfect consistency between PR and FAQ
- Include real-world examples and use cases

Create a press release and FAQ that clearly articulates the customer problem, solution, and benefits with maximum specificity, measurable outcomes, and quantitative evidence.""",

            # Variation 2: Emphasize customer-centricity
            """Generate a customer-focused PR-FAQ document following Amazon's working backwards methodology with emphasis on measurable customer outcomes.

Project: {projectName}
Problem: {problemDescription}
Context: {businessContext}

CORE PRINCIPLES:

1. START WITH THE CUSTOMER
   - Lead with customer pain points, not product features
   - Quantify customer problems with specific metrics
   - Show measurable customer outcomes (before/after comparisons)
   - Include real customer voices with quantitative results

2. MEASURABLE OUTCOMES (NON-NEGOTIABLE)
   - Every benefit must have a number: percentage, time saved, cost reduced
   - Use before/after comparisons: "from X to Y"
   - Include customer success metrics with specific data
   - Provide ROI or value metrics where applicable

3. AUTHENTIC EVIDENCE
   - 2-3 customer quotes with real names, titles, companies
   - Each quote must include specific, measurable outcomes
   - Format: "We achieved [specific metric]" - Name, Title, Company
   - Example: "Cut onboarding time from 2 weeks to 3 days" - Sarah Johnson, VP Operations

4. TECHNICAL CREDIBILITY
   - Explain the mechanism: HOW does it work?
   - Include architecture, workflow, or process details
   - Address technical implementation in FAQ
   - Specify integrations, technologies, standards

PRESS RELEASE FORMAT:
- Headline: Customer benefit-focused (not feature-focused)
- Opening: Customer problem with quantified impact
- Problem Statement: Specific pain points with metrics
- Solution: How it solves the problem (mechanism + benefits)
- Leadership Quote: Vision and customer commitment
- Key Benefits: 3-5 benefits, all with specific numbers
- Customer Quote: Real success story with metrics
- Availability: Specific dates and access information

FAQ STRUCTURE:
- Q1: What customer problem does this solve? (with metrics)
- Q2: Who is this for? (specific personas/segments)
- Q3: How does it work? (technical mechanism)
- Q4: What are the benefits? (quantified outcomes)
- Q5: When is it available? (specific timeline)
- Q6: How much does it cost? (pricing/value)
- Q7: What support is provided? (implementation/training)
- Additional questions as needed (max 15 total)

QUALITY REQUIREMENTS:
- Customer-first language throughout
- Specific metrics in every major claim
- No marketing fluff or hype
- Professional, authentic tone
- Technical depth where appropriate
- Perfect PR-FAQ consistency

Generate a compelling PR-FAQ that demonstrates clear customer value with quantitative evidence and measurable outcomes.""",

            # Variation 3: Focus on credibility and evidence
            """Generate an evidence-based PR-FAQ document following Amazon's format with emphasis on credibility, specificity, and quantitative validation.

Project: {projectName}
Problem: {problemDescription}
Context: {businessContext}

EVIDENCE-BASED REQUIREMENTS:

1. QUANTITATIVE VALIDATION
   - Support every claim with specific numbers or data
   - Include metrics: percentages, time savings, cost reductions, efficiency gains
   - Use comparative data: before/after, baseline/improved
   - Provide statistical evidence where possible
   - Example: "Reduces processing time by 73% (from 45 minutes to 12 minutes)"

2. CREDIBLE CUSTOMER EVIDENCE
   - Include 2-3 detailed customer quotes with full attribution
   - Each quote must contain specific, verifiable metrics
   - Format: "Detailed outcome with numbers" - Full Name, Title, Company Name
   - Example: "Our team deployed 40% faster, cutting release cycles from 10 days to 6 days" - Michael Chen, Director of Engineering, DataFlow Inc.

3. TECHNICAL SUBSTANCE
   - Explain the underlying mechanism and how it works
   - Include technical architecture or workflow details
   - Address implementation approach and methodology
   - Specify technologies, standards, or frameworks used
   - Provide integration and compatibility information

4. CLEAR TIMELINE AND AVAILABILITY
   - Specify exact dates: "Available March 15, 2025" or "Q2 2025"
   - Outline rollout phases with specific milestones
   - Include beta/pilot program details with dates
   - Set clear expectations for general availability

PRESS RELEASE COMPONENTS:
- Headline: Specific, benefit-focused, newsworthy (8-12 words)
- Opening Paragraph: Who, what, when, where, why - all with specific details
- Problem Statement: Customer pain quantified with real metrics
- Solution Description: How it works + measurable benefits
- Leadership Quote: Strategic vision with commitment to customers
- Key Benefits: 3-5 benefits, each with specific quantitative outcomes
- Customer Success Quote: Real results with specific metrics
- Availability Information: Exact dates and access details

FAQ REQUIREMENTS:
- Minimum 7, maximum 15 questions
- Priority order: problem (quantified), audience (specific), mechanism (detailed), benefits (measured), availability (dated), cost (transparent), support (comprehensive)
- Every answer must be specific and detailed (2-4 sentences)
- Include numbers and metrics wherever applicable
- Address potential objections with evidence
- Provide technical depth for implementation questions

CREDIBILITY STANDARDS:
- No unsubstantiated claims - everything must be specific
- No marketing superlatives or hype words
- Professional, authentic, trustworthy tone
- Consistent messaging between PR and FAQ
- Real-world examples and use cases
- Verifiable metrics and outcomes

Create a credible, evidence-based PR-FAQ that builds trust through specificity, quantitative validation, and measurable customer outcomes."""
        ]

        # Select variation based on request_id
        variation_index = request_id % len(variations)
        return variations[variation_index]

    def generate_response(self, request_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate appropriate response based on request type."""
        request_id = request_data.get('request_id', 'unknown')
        request_type = self.detect_request_type(request_data)

        print(f"    [DEBUG] Request {request_id}: Detected type = {request_type}", flush=True)

        if request_type == "mutation":
            content = self.generate_mutation(request_data)
            print(f"    [DEBUG] Request {request_id}: Generated mutation ({len(content)} chars)", flush=True)
            print(f"    [DEBUG] Request {request_id}: Content starts with: {content[:80]}", flush=True)
        elif request_type == "evaluation":
            content = self.generate_evaluation(request_data)
            print(f"    [DEBUG] Request {request_id}: Generated evaluation ({len(content)} chars)", flush=True)
            print(f"    [DEBUG] Request {request_id}: Content starts with: {content[:80]}", flush=True)
        elif request_type == "press_release":
            content = self.generate_press_release(request_data)
            print(f"    [DEBUG] Request {request_id}: Generated press_release ({len(content)} chars)", flush=True)
        elif request_type == "faq":
            content = self.generate_faq(request_data)
            print(f"    [DEBUG] Request {request_id}: Generated faq ({len(content)} chars)", flush=True)
        else:
            content = f"Generic response for request {request_data.get('request_id')}"
            print(f"    [DEBUG] Request {request_id}: Generated generic ({len(content)} chars)", flush=True)

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
        print(f"🔧 Auto-responder VERSION 3.0 - Instruction-only keyword detection", flush=True)
        print(f"Auto-responder monitoring: {self.request_dir}", flush=True)
        print(f"Responses will be written to: {self.response_dir}", flush=True)

        iteration = 0
        while True:
            iteration += 1

            # Find unprocessed requests
            request_files = sorted(self.request_dir.glob("request_*.json"))
            new_requests = [f for f in request_files if f.name not in self.processed_requests]

            if new_requests:
                print(f"\n[Iteration {iteration}] Found {len(new_requests)} new request(s)", flush=True)

            for request_file in new_requests:
                try:
                    # Read request
                    with open(request_file, 'r') as f:
                        request_data = json.load(f)

                    request_id = request_data.get("request_id")
                    request_type = self.detect_request_type(request_data)

                    print(f"  Processing request {request_id} (type: {request_type})", flush=True)

                    # Generate response
                    response_data = self.generate_response(request_data)

                    # Write response
                    response_file = self.response_dir / f"response_{request_id:04d}.json"
                    with open(response_file, 'w') as f:
                        json.dump(response_data, f, indent=2)

                    print(f"  ✓ Response written: {response_file.name}", flush=True)

                    # Mark as processed
                    self.processed_requests.add(request_file.name)

                except Exception as e:
                    print(f"  ✗ Error processing {request_file.name}: {e}", flush=True)

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


