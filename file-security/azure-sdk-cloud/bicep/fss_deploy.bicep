param storageAccountNames array
param storageAccountName string
param functionCodeBlobUrl string
param location string = 'eastus2'
param storageAccountLocations array
param amaasRegion string
@secure()
param amaasApiKey string
param appServiceSku string = 'Y1'
param appServiceTier string = 'Dynamic'
@secure()
param storageAccountKey string

var functionAppName = 'fss-scanner-${uniqueString(resourceGroup().id)}'

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' existing = {
  name: storageAccountName
}

resource storageQueueService 'Microsoft.Storage/storageAccounts/queueServices@2021-09-01' = {
  name: 'default'
  parent: storageAccount
}

resource storageQueue 'Microsoft.Storage/storageAccounts/queueServices/queues@2021-09-01' = {
  name: 'scan-queue'
  parent: storageQueueService
}

resource eventGridTopics 'Microsoft.EventGrid/systemTopics@2021-12-01' = [for (name, i) in storageAccountNames: {
  name: '${name}-eventgrid-topic'
  location: storageAccountLocations[i]
  properties: {
    source: resourceId('Microsoft.Storage/storageAccounts', name)
    topicType: 'Microsoft.Storage.StorageAccounts'
  }
}]

resource eventGridSubscriptions 'Microsoft.EventGrid/systemTopics/eventSubscriptions@2021-12-01' = [for (name, i) in storageAccountNames: {
  name: '${name}-event-subscription'
  parent: eventGridTopics[i]
  properties: {
    destination: {
      endpointType: 'StorageQueue'
      properties: {
        resourceId: resourceId('Microsoft.Storage/storageAccounts', storageAccountName)
        queueName: 'scan-queue'
        queueMessageTimeToLiveInSeconds: 3600
      }
    }
    filter: {
      includedEventTypes: [
        'Microsoft.Storage.BlobCreated'
      ]
    }
  }
}]

// Updated the plan name to include a unique string to avoid conflicts with an existing plan.
resource appServicePlan 'Microsoft.Web/serverfarms@2021-02-01' = {
  name: 'myFunctionPlan-${uniqueString(resourceGroup().id)}'
  properties: {
    reserved: true
  }
  location: location
  sku: {
    name: appServiceSku
    tier: appServiceTier
  }
  kind: 'linux'
}

resource appInsights 'Microsoft.Insights/components@2020-02-02' = {
  name: 'fss-scanner-insights-${uniqueString(resourceGroup().id)}'
  location: location
  kind: 'web'
  properties: {
    Application_Type: 'web'
  }
}

resource functionApp 'Microsoft.Web/sites@2021-02-01' = {
  name: functionAppName
  location: location
  kind: 'functionapp'
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    serverFarmId: appServicePlan.id
    siteConfig: {
      pythonVersion: '3.10'
      linuxFxVersion: 'PYTHON|3.10'
      appSettings: [
        {
          name: 'AzureWebJobsStorage'
          value: 'DefaultEndpoints=protocol=https;AccountName=${storageAccount.name};EndpointSuffix=${environment().suffixes.storage};AccountKey=${storageAccount.listKeys().keys[0].value}'
        }
        {
          name: 'FUNCTIONS_WORKER_RUNTIME'
          value: 'python'
        }
        {
          name: 'WEBSITE_RUN_FROM_PACKAGE'
          value: functionCodeBlobUrl
        }
        {
          name: 'AMAAS_REGION'
          value: amaasRegion
        }
        {
          name: 'AMAAS_API_KEY'
          value: amaasApiKey
        }
        {
          name: 'APPLICATIONINSIGHTS_CONNECTION_STRING'
          value: appInsights.properties.ConnectionString
        }
        {
          name: 'STORAGE_ACCOUNT_KEY'
          value: storageAccountKey
        }
        {
          name: 'FUNCTIONS_EXTENSION_VERSION'
          value: '~4'
        }
        {
          name: 'AzureWebJobsSecretStorageType'
          value: 'Files'
        }
      ]
    }
  }
}
