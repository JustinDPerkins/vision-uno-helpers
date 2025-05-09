import requests
import json
import argparse
from datetime import datetime, timedelta, timezone

# Set up argument parsing
parser = argparse.ArgumentParser(description='Trend Micro XDR Detections Retrieval')
parser.add_argument('-t', '--APIToken', type=str, required=True, help='Trend Micro XDR API Token')
args = parser.parse_args()

# Use timezone-aware datetime to address deprecation warning
end_date = datetime.now(timezone.utc)
start_date = end_date - timedelta(days=30)

url_base = 'https://api.xdr.trendmicro.com'
url_path = '/v3.0/oat/detections'
token = args.APIToken

query_params = {
    'detectedStartDateTime': start_date.strftime('%Y-%m-%dT%H:%M:%SZ'),
    'detectedEndDateTime': end_date.strftime('%Y-%m-%dT%H:%M:%SZ'),
    'ingestedStartDateTime': start_date.strftime('%Y-%m-%dT%H:%M:%SZ'),
    'ingestedEndDateTime': end_date.strftime('%Y-%m-%dT%H:%M:%SZ'),
    'top': 100
}
headers = {
    'Authorization': 'Bearer ' + token,
    'TMV1-Filter': "(riskLevel eq 'critical' or riskLevel eq 'high' or riskLevel eq 'medium' or riskLevel eq 'low' or riskLevel eq 'info') and (productCode eq 'ptn' or productCode eq 'pts')"
}

r = requests.get(url_base + url_path, params=query_params, headers=headers)

# Print headers for debugging
print(f"Status Code: {r.status_code}")
for k, v in r.headers.items():
    print(f'{k}: {v}')
print('')

# Enhanced error handling
if r.status_code != 200:
    print("Error occurred:")
    try:
        error_details = r.json()
        print(json.dumps(error_details, indent=4))
    except ValueError:
        print(r.text)
else:
    # Process successful response
    if 'application/json' in r.headers.get('Content-Type', '') and len(r.content):
        print(json.dumps(r.json(), indent=4))
    else:
        print(r.text)
