import requests

# 你的 GCP API Key
api_key = ""

def test_google_maps(api_key):
    url = "https://maps.googleapis.com/maps/api/geocode/json"
    params = {
        "address": "1600 Amphitheatre Parkway, Mountain View, CA",
        "key": api_key
    }
    response = requests.get(url, params=params)
    return response.status_code, response.json()

def test_youtube_data(api_key):
    url = "https://www.googleapis.com/youtube/v3/search"
    params = {
        "part": "snippet",
        "q": "Google",
        "key": api_key
    }
    response = requests.get(url, params=params)
    return response.status_code, response.json()

def test_custom_search(api_key):
    url = "https://www.googleapis.com/customsearch/v1"
    params = {
        "q": "Google",
        "cx": "YOUR_SEARCH_ENGINE_ID",  # 替换为你的自定义搜索引擎ID
        "key": api_key
    }
    response = requests.get(url, params=params)
    return response.status_code, response.json()

def test_cloud_sql(api_key):
    url = "https://sqladmin.googleapis.com/sql/v1beta4/projects/YOUR_PROJECT_ID/instances"
    params = {
        "key": api_key
    }
    response = requests.get(url, params=params)
    return response.status_code, response.json()

def test_compute_engine(api_key):
    url = "https://compute.googleapis.com/compute/v1/projects/YOUR_PROJECT_ID/zones/YOUR_ZONE/instances"
    params = {
        "key": api_key
    }
    response = requests.get(url, params=params)
    return response.status_code, response.json()

def test_cloud_storage(api_key):
    url = "https://www.googleapis.com/storage/v1/b"
    params = {
        "project": "YOUR_PROJECT_ID",
        "key": api_key
    }
    response = requests.get(url, params=params)
    return response.status_code, response.json()

def main(api_key):
    tests = {
        "Google Maps API": test_google_maps,
        "YouTube Data API": test_youtube_data,
        "Custom Search API": test_custom_search,
        "Cloud SQL Admin API": test_cloud_sql,
        "Compute Engine API": test_compute_engine,
        "Cloud Storage API": test_cloud_storage,
    }

    for api_name, test_func in tests.items():
        try:
            status_code, result = test_func(api_key)
            if status_code == 200:
                print(f"{api_name} 可以通过此 API Key 访问。")
            else:
                print(f"{api_name} 无法通过此 API Key 访问。状态码: {status_code}")
        except Exception as e:
            print(f"测试 {api_name} 时发生错误: {e}")

if __name__ == "__main__":
    main(api_key)
