"""Tests for LLM client with mock support."""

import pytest

from scripts.prompt_tuning.llm_client import MockLLMClient, create_llm_client


class TestMockLLMClient:
    """Tests for MockLLMClient."""

    @pytest.mark.asyncio
    async def test_generate_press_release(self):
        """Test mock press release generation."""
        client = MockLLMClient()
        prompt = "Generate a press release for our new product"

        response = await client.generate(prompt)

        assert response is not None
        assert len(response) > 0
        assert "press release" in response.lower() or "mock" in response.lower()

    @pytest.mark.asyncio
    async def test_generate_faq(self):
        """Test mock FAQ generation."""
        client = MockLLMClient()
        prompt = "Generate FAQ for our product"

        response = await client.generate(prompt)

        assert response is not None
        assert len(response) > 0
        assert "faq" in response.lower() or "question" in response.lower() or "mock" in response.lower()

    @pytest.mark.asyncio
    async def test_generate_evaluation(self):
        """Test mock evaluation generation."""
        client = MockLLMClient()
        prompt = "Evaluate the following content and provide a score"

        response = await client.generate(prompt)

        assert response is not None
        assert len(response) > 0
        # Should return JSON for evaluation
        assert "score" in response.lower() or "mock" in response.lower()

    @pytest.mark.asyncio
    async def test_deterministic_responses(self):
        """Test that mock responses are deterministic."""
        client = MockLLMClient()
        prompt = "Generate a press release"

        response1 = await client.generate(prompt)
        response2 = await client.generate(prompt)

        # Same prompt should give same response
        assert response1 == response2

    @pytest.mark.asyncio
    async def test_different_prompts_different_responses(self):
        """Test that different prompts give different responses."""
        client = MockLLMClient()

        response1 = await client.generate("Generate a press release")
        response2 = await client.generate("Generate a FAQ")

        # Different prompts should give different responses
        assert response1 != response2

    @pytest.mark.asyncio
    async def test_call_count(self):
        """Test that call count is tracked."""
        client = MockLLMClient()

        assert client.call_count == 0

        await client.generate("Test prompt 1")
        assert client.call_count == 1

        await client.generate("Test prompt 2")
        assert client.call_count == 2


class TestCreateLLMClient:
    """Tests for create_llm_client factory function."""

    def test_create_mock_client_explicit(self):
        """Test creating mock client explicitly."""
        client = create_llm_client(mock=True)

        assert isinstance(client, MockLLMClient)

    def test_create_mock_client_from_env(self, monkeypatch):
        """Test creating mock client from environment variable."""
        monkeypatch.setenv("AI_AGENT_MOCK_MODE", "true")

        client = create_llm_client()

        assert isinstance(client, MockLLMClient)

    def test_create_anthropic_client_requires_api_key(self):
        """Test that Anthropic client requires API key."""
        with pytest.raises(ValueError, match="ANTHROPIC_API_KEY"):
            create_llm_client(provider="anthropic", mock=False)

    def test_unsupported_provider(self, monkeypatch):
        """Test that unsupported provider raises error."""
        monkeypatch.setenv("ANTHROPIC_API_KEY", "test-key")
        with pytest.raises(ValueError, match="Unsupported LLM provider"):
            create_llm_client(provider="unsupported", mock=False)


@pytest.mark.asyncio
async def test_mock_client_integration():
    """Integration test for mock client workflow."""
    client = MockLLMClient()

    # Simulate a complete workflow
    pr_prompt = "Generate a press release for Project X that solves problem Y"
    pr_response = await client.generate(pr_prompt)
    assert len(pr_response) > 100  # Should be substantial

    faq_prompt = "Generate FAQ for Project X"
    faq_response = await client.generate(faq_prompt)
    assert len(faq_response) > 100

    eval_prompt = "Evaluate this content and score it"
    eval_response = await client.generate(eval_prompt)
    assert "score" in eval_response.lower() or "mock" in eval_response.lower()

    # All responses should be different
    assert pr_response != faq_response
    assert faq_response != eval_response
    assert pr_response != eval_response


if __name__ == "__main__":
    # Run tests
    pytest.main([__file__, "-v"])
