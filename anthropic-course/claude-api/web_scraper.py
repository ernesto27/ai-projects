import requests
from requests.exceptions import RequestException

def scrape_page(url):
    """
    Scrape a web page and print its HTML content.
    
    Args:
        url (str): The URL of the web page to scrape
        
    Returns:
        str: The HTML content of the page if successful, None otherwise
    """
    try:
        # Send a GET request to the URL
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(url, headers=headers)
        
        # Check if the request was successful
        response.raise_for_status()
        
        # Get the HTML content
        html_content = response.text
        
        # Print the HTML content
        print(f"Response status code: {response.status_code}")
        print("HTML Content:")
        print(html_content)
        
        return html_content
    
    except RequestException as e:
        print(f"Error during request: {e}")
        return None

# Example usage
if __name__ == "__main__":
    example_url = "https://example.com"
    scrape_page(example_url)