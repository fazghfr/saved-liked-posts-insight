"""
Data persistence repository - handles file I/O operations for Instagram data.
"""
import json
import pandas as pd


def load_instagram_data(mode: str, base_path_data: str = 'playground/data/your_instagram_activity/') -> dict:
    """
    Load Instagram data from JSON files.
    
    Args:
        mode: Either 'saved' or 'liked'
        base_path_data: Base path to Instagram data directory
        
    Returns:
        Dictionary containing the loaded JSON data
        
    Raises:
        ValueError: If mode is not 'saved' or 'liked'
        FileNotFoundError: If the data file doesn't exist
    """
    if mode == 'saved':
        path = base_path_data + 'saved/saved_posts.json'
    elif mode == 'liked':
        path = base_path_data + 'likes/liked_posts.json'
    else:
        raise ValueError("Unexpected Input: bad request")
    
    with open(path, 'r') as f:
        data_json = json.load(f)
    
    return data_json


def save_dataframe_to_json(df: pd.DataFrame, filepath: str) -> None:
    """
    Save a DataFrame to a JSON file.
    
    Args:
        df: DataFrame to save
        filepath: Path where the JSON file will be saved
    """
    df.to_json(filepath)
