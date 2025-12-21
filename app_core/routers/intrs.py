"""
Instagram data routes - handles HTTP requests for Instagram data processing.
"""
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel, Field
from typing import List, Dict, Any, Optional
import pandas as pd

from app_core.services.data_stuff import get_random_posts, process_post_data
from app_core.services.categ import categorize_posts, save_categorization_results


router = APIRouter(prefix="/posts", tags=["posts"])


# Request/Response Models
class SamplePostsRequest(BaseModel):
    mode: str = Field(..., description="Either 'saved' or 'liked'")
    sample_num: int = Field(10, description="Number of posts to sample", ge=1, le=100)
    seed: int = Field(42, description="Random seed for reproducibility")
    caption_map: Optional[Dict[int, str]] = Field(None, description="Optional mapping of post IDs to captions")


class PostData(BaseModel):
    title: str
    link: str
    timestamp: int
    id: int
    type: str
    captions: Optional[str] = ""


class SamplePostsResponse(BaseModel):
    posts: List[PostData]
    count: int


class CategorizeRequest(BaseModel):
    captions: List[str] = Field(..., description="List of captions to categorize")
    model: str = Field("tngtech/deepseek-r1t-chimera:free", description="LLM model to use")
    save_results: bool = Field(False, description="Whether to save results to files")


class CategoryResult(BaseModel):
    raw_reasoning: str
    summary_reasoning: str
    categories: List[str]


class CategorizeResponse(BaseModel):
    results: List[CategoryResult]
    count: int


@router.post("/sample", response_model=SamplePostsResponse)
async def sample_posts(request: SamplePostsRequest):
    """
    Get random sampled posts from Instagram data.
    
    Args:
        request: Sample request parameters
        
    Returns:
        Sampled posts with metadata
        
    Raises:
        HTTPException: If mode is invalid or data cannot be loaded
    """
    try:
        # Get random posts using service layer
        df = get_random_posts(
            mode=request.mode,
            sample_num=request.sample_num,
            seed=request.seed
        )
        
        # Process post data
        df = process_post_data(df, caption_map=request.caption_map)
        
        # Convert to response format
        posts = df.to_dict('records')
        
        return SamplePostsResponse(
            posts=posts,
            count=len(posts)
        )
    
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=f"Data file not found: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")


@router.post("/categorize", response_model=CategorizeResponse)
async def categorize_post_captions(request: CategorizeRequest):
    """
    Categorize post captions using LLM.
    
    Args:
        request: Categorization request with captions
        
    Returns:
        Categorization results for each caption
        
    Raises:
        HTTPException: If categorization fails
    """
    try:
        # Create DataFrame from captions
        df = pd.DataFrame({"captions": request.captions})
        
        # Categorize using service layer
        results = categorize_posts(df, model=request.model)
        
        # Save results if requested
        if request.save_results:
            save_categorization_results(results, df)
        
        return CategorizeResponse(
            results=results,
            count=len(results)
        )
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Categorization failed: {str(e)}")
