#!/usr/bin/env python3
"""PR-FAQ Validator Prompt Tuning Tool - Command Line Interface."""

import asyncio
import json
import os
import sys
from pathlib import Path

import click
from rich.console import Console
from rich.progress import Progress, SpinnerColumn, TextColumn
from rich.table import Table

# Add scripts directory to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from evolutionary_tuner import EvolutionaryTuner  # pylint: disable=wrong-import-position
from prompt_simulator import PromptSimulator  # pylint: disable=wrong-import-position
from prompt_tuning_config import (  # pylint: disable=wrong-import-position
    get_env_file_template,
    load_project_config,
    validate_api_keys,
)
from quality_evaluator import QualityEvaluator  # pylint: disable=wrong-import-position

console = Console()


@click.group()
@click.option("--verbose", "-v", is_flag=True, help="Enable verbose output")
@click.pass_context
def cli(ctx, verbose):
    """PR-FAQ Validator Prompt Tuning Tool"""
    ctx.ensure_object(dict)
    ctx.obj["verbose"] = verbose


@cli.command()
@click.argument("project_name")
@click.option("--base-dir", type=click.Path(exists=True), help="Base directory (default: current)")
def init(project_name, base_dir):
    """Initialize a new prompt tuning project"""
    base_path = Path(base_dir) if base_dir else Path.cwd()

    console.print(f"[bold blue]Initializing prompt tuning project: {project_name}[/bold blue]")

    # Create directory structure
    results_dir = base_path / f"prompt_tuning_results_{project_name}"
    results_dir.mkdir(exist_ok=True)

    prompts_dir = base_path / "prompts"
    if not prompts_dir.exists():
        console.print(f"[yellow]Warning: prompts directory not found at {prompts_dir}[/yellow]")
        console.print("Please ensure your prompts are in the 'prompts/' directory")

    # Create .env template if it doesn't exist
    env_file = base_path / ".env"
    if not env_file.exists():
        with open(env_file, "w", encoding="utf-8") as f:
            f.write(get_env_file_template())
        console.print(f"[green]Created .env template at {env_file}[/green]")
        console.print("[yellow]Please fill in your API keys in the .env file[/yellow]")

    # Create sample test cases file
    test_cases_file = results_dir / f"test_cases_{project_name}.json"
    if not test_cases_file.exists():
        sample_test_cases = {
            "project": project_name,
            "description": "Test cases for PR-FAQ prompt tuning",
            "test_cases": [
                {
                    "id": "tc001",
                    "name": "E-commerce Checkout Optimization",
                    "description": "PR-FAQ for improving checkout conversion rates",
                    "industry": "E-commerce",
                    "project_type": "Product Feature",
                    "scope": "Medium",
                    "stakeholder_complexity": "High",
                    "inputs": {
                        "projectName": "One-Click Checkout",
                        "problemDescription": (
                            "Customers abandon carts due to lengthy checkout process with 15+ form fields"
                        ),
                        "businessContext": (
                            "Reduce cart abandonment from 68% to under 40%, increase conversion by 25%"
                        ),
                    },
                },
                {
                    "id": "tc002",
                    "name": "Developer Productivity Tool",
                    "description": "PR-FAQ for automated code review assistant",
                    "industry": "Technology",
                    "project_type": "Internal Tool",
                    "scope": "Large",
                    "stakeholder_complexity": "Medium",
                    "inputs": {
                        "projectName": "AI Code Review Assistant",
                        "problemDescription": (
                            "Engineers spend 8-10 hours per week on manual code reviews, delaying releases"
                        ),
                        "businessContext": (
                            "Reduce review time by 60%, improve code quality, accelerate release cycles"
                        ),
                    },
                },
            ],
        }

        with open(test_cases_file, "w", encoding="utf-8") as f:
            json.dump(sample_test_cases, f, indent=2)

        console.print(f"[green]Created sample test cases at {test_cases_file}[/green]")
        console.print("[yellow]Please customize the test cases for your project[/yellow]")

    console.print(f"[bold green]Project {project_name} initialized successfully![/bold green]")
    console.print(f"Results directory: {results_dir}")


@cli.command()
@click.argument("project_name")
@click.option("--iteration", "-i", default=0, help="Iteration number (0 for baseline)")
@click.option("--mock", is_flag=True, help="Run in AI agent mock mode (no API keys required)")
def simulate(project_name, iteration, mock):
    """Run prompt simulation for test cases"""
    config = load_project_config(project_name)

    # Set mock mode if requested
    if mock:
        os.environ["AI_AGENT_MOCK_MODE"] = "true"
        console.print("[yellow]Running in AI Agent Mock Mode - no API keys required[/yellow]")

    # Validate API keys (skip if in mock mode)
    if not mock:
        missing_keys = validate_api_keys(config)
        if missing_keys:
            console.print("[bold red]Missing API keys:[/bold red]")
            for key in missing_keys:
                console.print(f"  - {key}")
            console.print("[yellow]Tip: Use --mock flag to run without API keys[/yellow]")
            console.print("Please set these environment variables or update your .env file")
            sys.exit(1)

    console.print(f"[bold blue]Running simulation for {project_name} (iteration {iteration})[/bold blue]")

    async def run_simulation():
        simulator = PromptSimulator(config)

        with Progress(
            SpinnerColumn(), TextColumn("[progress.description]{task.description}"), console=console
        ) as progress:
            task = progress.add_task("Running simulation...", total=None)
            results = await simulator.run_simulation(iteration)
            progress.update(task, description="Saving results...")
            output_file = simulator.save_results(results, iteration)

        console.print("[bold green]Simulation complete![/bold green]")
        console.print(f"Results saved to: {output_file}")

        return results

    return asyncio.run(run_simulation())


@cli.command()
@click.argument("project_name")
@click.option("--iteration", "-i", default=0, help="Iteration number (0 for baseline)")
@click.option("--mock", is_flag=True, help="Run in AI agent mock mode (no API keys required)")
def evaluate(project_name, iteration, mock):
    """Evaluate simulation results"""
    config = load_project_config(project_name)

    # Set mock mode if requested
    if mock:
        os.environ["AI_AGENT_MOCK_MODE"] = "true"

    console.print(f"[bold blue]Evaluating results for {project_name} (iteration {iteration})[/bold blue]")

    async def run_evaluation():
        # Load simulation results
        sim_file = config.results_dir / f"simulation_iteration_{iteration:03d}.json"
        if not sim_file.exists():
            console.print(f"[bold red]Simulation results not found: {sim_file}[/bold red]")
            console.print("Run 'simulate' command first")
            sys.exit(1)

        with open(sim_file, "r", encoding="utf-8") as f:
            simulation_results = json.load(f)

        evaluator = QualityEvaluator(config)

        with Progress(
            SpinnerColumn(), TextColumn("[progress.description]{task.description}"), console=console
        ) as progress:
            task = progress.add_task("Evaluating quality...", total=None)
            evaluation_results = await evaluator.evaluate_results(simulation_results)
            progress.update(task, description="Saving evaluation...")
            output_file = evaluator.save_evaluation(evaluation_results, iteration)

        # Display results
        console.print("\n[bold green]Evaluation complete![/bold green]")
        console.print(f"Results saved to: {output_file}")

        agg_scores = evaluation_results.get("aggregate_scores", {})
        console.print("\n[bold]Aggregate Scores:[/bold]")
        console.print(f"  Overall: {agg_scores.get('overall', 0):.2f}")
        console.print(f"  Press Release: {agg_scores.get('press_release', 0):.2f}")
        console.print(f"  FAQ: {agg_scores.get('faq', 0):.2f}")

        return evaluation_results

    return asyncio.run(run_evaluation())


@cli.command()
@click.argument("project_name")
@click.option("--max-iterations", "-n", default=20, help="Maximum number of iterations")
@click.option("--mock", is_flag=True, help="Run in AI agent mock mode (no API keys required)")
def evolve(project_name, max_iterations, mock):
    """Run evolutionary prompt optimization"""
    config = load_project_config(project_name)
    config.max_iterations = max_iterations

    # Set mock mode if requested
    if mock:
        os.environ["AI_AGENT_MOCK_MODE"] = "true"
        console.print("[yellow]Running in AI Agent Mock Mode - no API keys required[/yellow]")

    # Validate API keys (skip if in mock mode)
    if not mock:
        missing_keys = validate_api_keys(config)
        if missing_keys:
            console.print("[bold red]Missing API keys:[/bold red]")
            for key in missing_keys:
                console.print(f"  - {key}")
            console.print("[yellow]Tip: Use --mock flag to run without API keys[/yellow]")
            sys.exit(1)

    console.print(f"[bold blue]Starting evolutionary optimization for {project_name}[/bold blue]")
    console.print(f"Max iterations: {max_iterations}")

    async def run_evolution():
        tuner = EvolutionaryTuner(config)
        results = await tuner.run_evolution(max_iterations)

        console.print("\n[bold green]Optimization complete![/bold green]")
        console.print(f"Baseline score: {results['baseline_score']:.2f}")
        console.print(f"Final score: {results['final_score']:.2f}")
        console.print(f"Improvement: {results['improvement']:.2f}")

        return results

    return asyncio.run(run_evolution())


@cli.command()
@click.argument("project_name")
def status(project_name):
    """Show optimization status and results"""
    config = load_project_config(project_name)

    console.print(f"[bold blue]Status for {project_name}[/bold blue]\n")

    # Check for final results
    final_results_file = config.results_dir / "optimization_final_results.json"
    if final_results_file.exists():
        with open(final_results_file, "r", encoding="utf-8") as f:
            results = json.load(f)

        console.print("[bold green]Optimization completed[/bold green]")
        console.print(f"Timestamp: {results.get('timestamp', 'N/A')}")
        console.print(f"Iterations: {results.get('max_iterations', 'N/A')}")
        console.print("\n[bold]Scores:[/bold]")
        console.print(f"  Baseline: {results.get('baseline_score', 0):.2f}")
        console.print(f"  Final: {results.get('final_score', 0):.2f}")
        console.print(f"  Improvement: {results.get('improvement', 0):.2f}")

        # Show iteration history
        history = results.get("iteration_history", [])
        if history:
            console.print("\n[bold]Recent Iterations:[/bold]")
            table = Table()
            table.add_column("Iteration", style="cyan")
            table.add_column("Score", style="magenta")
            table.add_column("Best", style="green")
            table.add_column("Improved", style="yellow")

            for h in history[-10:]:  # Last 10 iterations
                improved = "✓" if h.get("improved", False) else "✗"
                table.add_row(
                    str(h.get("iteration", 0)), f"{h.get('score', 0):.2f}", f"{h.get('best_score', 0):.2f}", improved
                )

            console.print(table)
    else:
        console.print("[yellow]No optimization results found[/yellow]")
        console.print("Run 'evolve' command to start optimization")


if __name__ == "__main__":
    cli()
