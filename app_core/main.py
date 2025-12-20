"""
FastAPI application entry point.
"""
from dotenv import load_dotenv
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.routers import intrs

# Load environment variables from .env file
load_dotenv()


# Create FastAPI application
app = FastAPI(
    title="Instagram Data Processor API",
    description="API for processing and categorizing Instagram posts using LLM",
    version="1.0.0"
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(intrs.router)


@app.get("/")
async def root():
    """Root endpoint with API information."""
    return {
        "message": "Instagram Data Processor API",
        "docs": "/docs",
        "redoc": "/redoc"
    }


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy"}
