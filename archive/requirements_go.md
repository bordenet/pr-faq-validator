# Purpose of file
Yes-- this file has been included on purpose to illustrate vibe coding from a guy who doesn't have a lot of free time to code from scratch. Here, I'm using ChatGPT to bootstrap project creation.

# Body
Note, we will implement as much of this project in Go as makes sense. Python for those bits and pieces which require it-- e.g. MCP integration and/or LangGraph integration.
⸻

🔧 Updated Tech Stack Plan

Component	Language	Notes
CLI & file parsing	Go	Fast, clean, and easily deployable
Markdown section splitter	Go	Regex or parser library
LLM interaction (OpenAI/Groq)	Go (via HTTP) or Python (LangGraph module)	Go if keeping simple; Python if using LangGraph
LangGraph workflow (optional)	Python	Only if you want multi-step node orchestration
Report rendering	Go	Colorized terminal + optional JSON
Unit tests	Go	Use go test with fixtures
LLM prompt templates	Text or Go constants	Reuse across components


⸻

🧱 Project Scaffold

prfaq-reviewer/
├── cmd/
│   └── prfaq-reviewer/          # CLI entrypoint (Go)
├── internal/
│   ├── parser/                  # Markdown parsing & section tagging (Go)
│   ├── analyzer/                # Flags/suggestions logic (Go)
│   └── llm/
│       ├── client.go            # Calls out to LLM service (Go)
│       └── langgraph_bridge.py # (Optional) Python fallback for LangGraph
├── prompts/
│   └── review_prompt.txt       # System/user prompt templates
├── testdata/
│   └── example_prfaq.md
├── go.mod / go.sum
├── README.md
└── Makefile


⸻

🪜 Development Phases (Go-First)

⸻

✅ Phase 1: All-in-Go MVP
	•	CLI: prfaq-reviewer input.md
	•	Parse Markdown into:
	•	Press Release
	•	FAQ
	•	Success Metrics
	•	Send entire content as one prompt to OpenAI API (via Go HTTP client)
	•	Output:
	•	Summary
	•	Flags
	•	Suggestions
	•	Readiness Score

⸻

🧪 Phase 2: Modular LLM Node Workflow (Optional via LangGraph)

If you find that a single-shot prompt becomes messy, move this logic into Python/LangGraph:
	•	Go sends markdown chunks to Python as JSON
	•	Python LangGraph routes each section to different LLM nodes
	•	Returns structured output to Go

Simple interface contract:

POST /analyze
{
  "press_release": "...",
  "faq": "...",
  "metrics": "..."
}
→
{
  "summary": "...",
  "flags": [...],
  "suggestions": [...],
  "score": 7.5
}

You can call this as a subprocess from Go (or run a local HTTP server in Python).
