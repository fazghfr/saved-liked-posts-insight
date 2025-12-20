"""
Categorization service - contains business logic for LLM-based categorization.
"""
import os
import json
import pandas as pd
from typing import List, Dict, Any
from app.repositories.llm_obj import call_llm
from app.services.data_stuff import parse_llm_output


def categorize_posts(df: pd.DataFrame, model: str = "tngtech/deepseek-r1t-chimera:free") -> List[Dict[str, Any]]:
    """
    Categorize posts using LLM.
    
    Args:
        df: DataFrame with posts containing 'captions' column
        model: LLM model to use for categorization
        
    Returns:
        List of parsed categorization results
    """
    results = []
    
    for caption in df["captions"]:
        prompt = "This is the caption that you need to categorize\\n" + caption
        raw, content = call_llm(prompt, model)
        
        # Parse the LLM output
        parsed = parse_llm_output(content)
        results.append(parsed)
    
    return results


def save_categorization_results(results: List[Dict[str, Any]], df: pd.DataFrame, 
                                output_dir: str = "outputs") -> None:
    """
    Save categorization results to files.
    
    Args:
        results: List of parsed categorization results
        df: Original DataFrame with captions
        output_dir: Directory to save results
    """
    os.makedirs(output_dir, exist_ok=True)
    
    # Save individual results
    for i, (result, caption) in enumerate(zip(results, df["captions"])):
        # Save parsed result as JSON
        with open(f"{output_dir}/content_{i}.txt", "w", encoding="utf-8") as f:
            f.write(str(result))
    
    # Save all results as JSON
    with open(f"{output_dir}/parsed_results.json", "w", encoding="utf-8") as f:
        json.dump(results, f, ensure_ascii=False, indent=2)
