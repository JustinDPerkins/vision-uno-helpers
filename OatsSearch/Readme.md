# Trend Micro XDR Detections Retrieval Script

## Overview
This Python script retrieves Observed Attack Techniques (OAT) events from the Trend Micro XDR API, with flexible filtering options.

## Prerequisites
- Python 3.8+
- `requests` library
- Trend Micro XDR API Token

## Installation
1. Clone the repository
2. Install required dependencies:
```bash
pip install requests
```

## Usage

### Python Script
Run the script with your API token:

```bash
# Long-form argument
python txone.py --APIToken YOUR_API_TOKEN_HERE

# Short-form argument
python txone.py -t YOUR_API_TOKEN_HERE
```

### Curl Command
For direct API testing, use this curl command:

```bash
curl -X GET \
  "https://api.xdr.trendmicro.com/v3.0/oat/detections?detectedStartDateTime=$(date -u -v-30d +%Y-%m-%dT%H:%M:%SZ)&detectedEndDateTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)&ingestedStartDateTime=$(date -u -v-30d +%Y-%m-%dT%H:%M:%SZ)&ingestedEndDateTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)&top=100" \
  -H "Authorization: Bearer YOUR_API_TOKEN_HERE" \
  -H "TMV1-Filter: (riskLevel eq 'critical' or riskLevel eq 'high' or riskLevel eq 'medium' or riskLevel eq 'low' or riskLevel eq 'info') and (productCode eq 'ptn' or productCode eq 'pts')"
```

## Features
- Retrieves detections from the past 30 days
- Filters by multiple risk levels
- Supports multiple product codes
- Command-line API token input

## Filter Details
The script uses the following filter:
- Risk Levels: critical, high, medium, low, info
- Product Codes: ptn, pts

## Notes
- Ensure you have a valid Trend Micro API Token
- The script uses UTC time for date ranges
- Maximum data retrieval is limited to 365 days

## Troubleshooting
- Verify your API token is correct
- Check network connectivity
- Ensure you have the necessary API permissions

