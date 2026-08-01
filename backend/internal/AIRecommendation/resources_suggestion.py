"""
AI Recommendation Module - Real-time resource suggestions using Tavily API.
Fetches flashcard titles per user → batch queries Tavily → returns resources.
No database caching - computed on-demand per request.
"""

import os
import logging
import sys
import json
from typing import List, Dict
import requests

# Configure logging to only show warnings and errors
logging.basicConfig(level=logging.WARNING, format='%(levelname)s: %(message)s', stream=sys.stderr)
logger = logging.getLogger(__name__)

def get_tavily_api_key() -> str:
    """Get Tavily API key from environment."""
    api_key = os.getenv("TAVILY_API_KEY")
    if not api_key:
        raise ValueError("TAVILY_API_KEY environment variable is not set")
    return api_key


def search_resources(title: str, max_results: int = 2) -> List[Dict]:
    """
    Search for resources related to a flashcard title using Tavily API.
    
    Args:
        title: Flashcard title to search for
        max_results: Number of results to return
    
    Returns:
        List of {"title": str, "link": str} resources
    """
    api_key = get_tavily_api_key()
    
    try:
        response = requests.post(
            "https://api.tavily.com/search",
            json={
                "api_key": api_key,
                "query": title,
                "max_results": max_results,
                "include_answer": False,
            },
            timeout=10,
        )
        
        if response.status_code != 200:
            logger.warning(f"Tavily API error for '{title}': {response.status_code}")
            return []
        
        data = response.json()
        results = data.get("results", [])
        
        # Extract title and link from results
        resources = []
        for result in results:
            url = result.get("url", "")
            if url:
                resources.append({
                    "title": result.get("title", ""),
                    "link": url,
                })
        
        return resources
    
    except Exception as e:
        logger.error(f"Error searching Tavily for '{title}': {e}")
        return []


def get_recommendations_for_titles(flashcard_titles: List[str], max_results_per_title: int = 3) -> Dict[str, List[Dict]]:
    """
    Get related resources for multiple flashcard titles.
    Queries Tavily for each title and returns results.
    
    Args:
        flashcard_titles: List of flashcard titles
        max_results_per_title: Number of results per title
    
    Returns:
        Dictionary mapping flashcard_title -> list of {"title", "link"} resources
    """
    if not flashcard_titles:
        return {}
    
    recommendations = {}
    for title in flashcard_titles:
        resources = search_resources(title, max_results_per_title)
        recommendations[title] = resources
    
    return recommendations


# Main entry point for command-line execution
if __name__ == "__main__":
    try:
        if len(sys.argv) < 3:
            logger.error("Usage: python resources_suggestion.py <titles_json> <max_results>")
            print(json.dumps({}))
            sys.exit(0)
        
        titles = json.loads(sys.argv[1])
        max_results = int(sys.argv[2])
        
        result = get_recommendations_for_titles(titles, max_results)
        print(json.dumps(result))
    except Exception as e:
        logger.error(f"Fatal error: {e}", exc_info=True)
        print(json.dumps({}))
