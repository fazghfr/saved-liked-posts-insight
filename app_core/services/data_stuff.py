"""
Data processing service - contains business logic for Instagram data processing.
"""
import pandas as pd
from typing import Dict, Any, Optional
from app_core.repositories.data_prs import load_instagram_data


def get_random_posts(mode: str, sample_num: int = 10, seed: int = 42) -> pd.DataFrame:
    """
    Get random sampled posts from Instagram data.
    
    Args:
        mode: Either 'saved' or 'liked'
        sample_num: Number of posts to sample
        seed: Random seed for reproducibility
        
    Returns:
        DataFrame with sampled posts containing title, link, and timestamp
        
    Raises:
        ValueError: If mode is not 'saved' or 'liked'
    """
    if mode == 'saved':
        attr = 'saved_saved_media'
        getter = lambda item: (
            item["title"],
            item["string_map_data"]["Saved on"]["href"],
            item["string_map_data"]["Saved on"]["timestamp"],
        )
    elif mode == 'liked':
        attr = 'likes_media_likes'
        getter = lambda item: (
            item["title"],
            item["string_list_data"][0]["href"],
            item["string_list_data"][0]["timestamp"],
        )
    else:
        raise ValueError("Unexpected Input: bad request")
    
    # Load data using repository
    data_json = load_instagram_data(mode)
    
    # Transform into DataFrame
    df = pd.DataFrame([
        {
            "title": title,
            "link": link,
            "timestamp": timestamp
        }
        for title, link, timestamp in (getter(item) for item in data_json[attr])
    ])
    
    return df.sample(sample_num, random_state=seed)


def process_post_data(df: pd.DataFrame, caption_map: Optional[Dict[int, str]] = None) -> pd.DataFrame:
    """
    Process post data by adding id, type, and captions.
    
    Args:
        df: DataFrame with post data
        caption_map: Optional mapping of post IDs to captions
        
    Returns:
        Processed DataFrame with id, type, and captions columns
    """
    # Add id column from index
    df['id'] = df.index
    
    # Add type column based on link pattern
    df['type'] = df['link'].apply(
        lambda x: 'post' if '/p/' in x else 'reel'
    )
    
    # Add captions if caption_map is provided
    if caption_map is not None:
        df['captions'] = df['id'].map(caption_map).fillna("")
    
    return df


def parse_llm_output(text: str) -> Dict[str, Any]:
    """
    Parse LLM output into structured format.
    
    Args:
        text: Raw LLM response text
        
    Returns:
        Dictionary with raw_reasoning, summary_reasoning, and categories
    """
    import re
    import ast
    
    # 1. Split by the separator pattern (=-=-=-=) - handle both escaped and literal newlines
    separator_pattern = r"=-=-=-=\s*(?:\\n|\n)"
    parts = re.split(separator_pattern, text, maxsplit=1)
    
    if len(parts) > 1:
        # Separator found - first part is raw reasoning, second is structured
        raw_reasoning = parts[0].strip()
        structured = parts[1].strip()
    else:
        # No separator - try to split by REASONING keyword
        reasoning_split = re.split(r"\bREASONING\s*:", text, maxsplit=1, flags=re.IGNORECASE)
        raw_reasoning = reasoning_split[0].strip()
        structured = reasoning_split[1] if len(reasoning_split) > 1 else ""
    
    # 2. Extract summary reasoning (after REASONING: and before CATEGORIES:)
    # Handle both escaped \n and literal newlines
    summary_match = re.search(
        r"REASONING\s*:\s*(.*?)\s*(?:\\n|\n)\s*CATEGORIES\s*:",
        structured,
        re.DOTALL | re.IGNORECASE
    )
    
    summary_reasoning = summary_match.group(1).strip() if summary_match else ""
    
    # 3. Extract categories list
    categories_match = re.search(
        r"CATEGORIES\s*:\s*(\[[^\]]*\])",
        structured,
        re.IGNORECASE
    )
    
    categories = []
    if categories_match:
        try:
            categories = ast.literal_eval(categories_match.group(1))
        except:
            categories = []
    
    return {
        "raw_reasoning": raw_reasoning,
        "summary_reasoning": summary_reasoning,
        "categories": categories
    }
