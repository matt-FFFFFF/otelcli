# Open Telemetry CLI

## Creation of kusto table and ingestion mapping

````kql
// Create table command
////////////////////////////////////////////////////////////
.create table ['eventHubData']  (['records']:dynamic)
// Create mapping command
////////////////////////////////////////////////////////////
.create table ['eventHubData'] ingestion json mapping 'eventHubData_mapping' '[{"column":"records", "Properties":{"Path":"$[\'records\']"}}]'
````

## curl

```bash
curl -X POST 'https://location-0.in.applicationinsights.azure.com/v2.1/track' -H "Content-Type: application/x-json-stream" -d @file.json
````
