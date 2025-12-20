"""
LLM repository - handles external LLM API calls.
"""
import os
from openai import OpenAI
from typing import Tuple, Any


def get_llm_client(api_key: str = None, base_url: str = "https://openrouter.ai/api/v1") -> OpenAI:
    """
    Create and return an OpenAI client configured for OpenRouter.
    
    Args:
        api_key: API key for authentication (defaults to env variable)
        base_url: Base URL for the API
        
    Returns:
        Configured OpenAI client
    """
    if api_key is None:
        # Try OPENROUTER_API_KEY first, then fall back to OPENAI_API_KEY
        api_key = os.getenv("OPENROUTER_API_KEY") or os.getenv("OPENAI_API_KEY", "from-env")
    
    import httpx
    
    return OpenAI(
        api_key=api_key,
        base_url=base_url,
        http_client=httpx.Client()
    )


def call_llm(prompt: str, model: str = "google/gemini-2.0-flash:free", 
             client: OpenAI = None) -> Tuple[Any, str]:
    """
    Call the LLM API with a categorization prompt.
    
    Args:
        prompt: The prompt to send to the LLM
        model: Model identifier to use
        client: Optional pre-configured OpenAI client
        
    Returns:
        Tuple of (raw_response, content_string)
    """
    if client is None:
        client = get_llm_client()
    
    response = client.chat.completions.create(
        model=model,
        messages=[
            {
                "role": "system", 
                "content": "You are a core AI that will categorize captions into the right categories. You will be given captions of instagram posts, and categorize them accordingly."
            },
            {
                "role": "system", 
                "content": "Please respond with this format. =-=-=-=\\n\\nREASONING : {YOUR REASONING} \\nCATEGORIES: {consisting array of categories, if multiple. if singular, use array with one item}"
            },
            {
                "role": "user", 
                "content": prompt
            }
        ],
        temperature=0.2
    )
    
    return response, response.choices[0].message.content
