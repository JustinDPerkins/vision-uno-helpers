# File Security Scanner Azure Function

This repository contains an Azure Function deployment that sets up a file scanning service using Azure Functions, Event Grid, and Storage Queues.

## Prerequisites

- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- [Azure Bicep CLI](https://docs.microsoft.com/en-us/azure/azure-resource-manager/bicep/install)
- An Azure subscription
- A storage account containing the function code package (zip)
- Python 3.11 or later
- Vision One File Security SDK access (AMAAS API key)

## Architecture

The deployment creates the following resources:
- Azure Function App (Python)
- Application Insights
- Storage Queue
- Event Grid System Topics and Subscriptions
- App Service Plan

### Function Components

The solution consists of two Azure Functions:

1. **BlobCreatedTrigger**
   - Triggered by Event Grid blob creation events
   - Generates SAS URLs for secure blob access
   - Queues scan requests to the storage queue

2. **MalwareScanner**
   - Triggered by messages in the storage queue
   - Downloads files using SAS URLs
   - Scans files using Vision One File Security SDK
   - Logs scan results

### Code Structure

```
functions/
├── BlobCreatedTrigger/
│   ├── __init__.py          # Event Grid trigger handler
│   └── function.json        # Function binding configuration
├── MalwareScanner/
│   ├── __init__.py          # Queue trigger handler
│   ├── function.json        # Function binding configuration
│   └── requirements.txt     # Python dependencies
├── Shared/
│   ├── __init__.py
│   └── Scanner.py          # Scanner class implementation
├── host.json               # Function app configuration
└── requirements.txt        # Project-level dependencies
```

## Dependencies

```txt
azure-functions
azure-storage-blob
azure-storage-queue
visionone-filesecurity
requests
```

## Deployment Steps

1. **Login to Azure**
   ```bash
   az login
   ```

2. **Set your subscription**
   ```bash
   az account set --subscription <your-subscription-id>
   ```

3. **Create a Resource Group**
   ```bash
   az group create --name <resource-group-name> --location <location>
   ```

4. **Zip the function code**

    ```bash
    cd fss-sdk-azure-stack
    zip -r fss_scanner.zip functions/*
    ```

5. **Upload the function code to the storage account**

    ```bash
    az storage blob upload --account-name <storage-account> \
                         --container-name <container> \
                         --name fss_scanner.zip \
                         --file fss_scanner.zip \
                         --auth-mode key 
    ``` 

6. **List Blob URL**

    ```bash
    az storage blob url --account-name <storage-account> \
                         --container-name <container> \
                         --name fss_scanner.zip 
    ``` 

7. **List storage account keys of code blob account**

    ```bash
    az storage account keys list --resource-group <resource-group-name> --account-name <storage-account>
    ```

7. SAS token Generation for file retrieval

fss-sdk-azure-stack % az storage blob generate-sas \  --account-name <place holder> \
  --container-name <place holder> \
  --name fss_scanner.zip \
  --permissions r \
  --expiry 2025-<month>-<date>T00:00:00Z \
  --output tsv

8. **Update Parameters File**
   
   Edit `fss_deploy.parameters.json` with your values:
   ```json
   {
       "parameters": {
           "storageAccountNames": {
               "value": ["monitoredstorage1", "monitoredstorage2"]  // Storage accounts to monitor for files being upload too.
           },
           "storageAccountName": {
               "value": "mystorageaccountcodeholder"  // Storage account where the function code zip file is stored.
           },
           "functionCodeBlobUrl": {
               "value": "https://<yourstorage>.blob.core.windows.net/<container>/fss_scanner.zip"
           },
           "location": {
               "value": "<eastus2>"  // Location of the storage accounts and function app
           },
           "storageAccountLocations": {
               "value": "<eastus2>"  // Locations matching storageAccountNames
           },
           "amaasRegion": {
               "value": "<your-amaas-region>" // This is the region for the AMAAS API
           },
           "amaasApiKey": {
               "value": "<your-amaas-api-key>" // This is the key for the AMAAS API
           },
           "storageAccountKey": {
               "value": "your-storage-account-key" // This is the key for the storage account that will be used to store the function code 
           }
       }
   }
   ```

7. **Deploy using Bicep**
   ```bash
   cd bicep
   az deployment group create \
     --resource-group <resource-group-name> \
     --template-file fss_deploy.bicep \
     --parameters fss_deploy.parameters.json
   ```

## Configuration Details

### Storage Accounts
- The deployment will create Event Grid subscriptions for each storage account specified in `storageAccountNames`
- Each subscription will monitor for blob creation events
- Events will be routed to a storage queue named `scan-queue`

### Function App
- Python runtime
- System-assigned managed identity
- Connected to Application Insights
- Configured with necessary environment variables:
  - AMAAS_REGION
  - AMAAS_API_KEY
  - STORAGE_ACCOUNT_KEY
  - Application Insights connection

### Security
- Sensitive parameters (amaasApiKey, storageAccountKey) are marked as secure
- Function app uses managed identity to access the storage queue
- Storage Queue Data Reader role assigned to the function app
- SAS URLs are generated with minimum required permissions and 1-hour expiry

## Function Implementation Details

### BlobCreatedTrigger
This function is triggered when new blobs are created in the monitored storage accounts:
```python
def main(event: func.EventGridEvent, outputQueueItem: func.Out[str]) -> None:
    # Extracts blob information from the event
    # Generates a SAS URL for secure access
    # Queues a scan request
```

### MalwareScanner
This function processes queued scan requests:
```python
def main(msg: str) -> None:
    # Processes queue messages
    # Downloads files using SAS URLs
    # Scans files using Vision One File Security SDK
    # Logs results
```

## Post-Deployment

After deployment:
1. Verify the function app is running
2. Test by uploading a file to one of the monitored storage accounts
3. Check the storage queue for the scan message
4. Monitor function execution in Application Insights

## Troubleshooting

Common issues:
- **Function not triggering**: Check Event Grid subscription and queue permissions
- **Function failing**: Check Application Insights logs
- **Authentication errors**: Verify the AMAAS API key and storage account key
- **Scanner errors**: Check AMAAS_REGION and AMAAS_API_KEY configuration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

[Add your license information here]
