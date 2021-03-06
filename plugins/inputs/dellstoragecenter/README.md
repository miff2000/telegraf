# DellStorageCenter Input Plugin

The dellstoragecenter input plugin gathers volume performance data.

## A Few Notable Things...
### Storage Center Timestamps Used  
The timestamp returned by the Storage Center's JSON responses is used as the timestamp for the metric that Telegraf returns. The result of this is, although Telegraf polls the Storage Center API on a 1 minute interval by default, new metrics will only be produced if the timestamp differs. This stops points being added every minute when your metrics data hasn't changed, and makes your graphing apps like Grafana draw more meaningful graphs.

### History Filter For I/O Usage 
The Storage Center uses the concept of a History Filter to select the IO Usage statistics you require. Within that you specify a StartTime and EndTime. The StartTime is presently hard-coded at **`now() - 15 minutes`** as, from what I can tell, the Storage Center produces its IO Usage statistics at least once every 15 minutes. If you know this to be incorrect, please raise an issue for it.

This value is used to limit the response size from the Storage Center, so we can retrieve only the last I/O Usage information.

### History Filter For Storage Usage
Similarly to the History Filter For I/O Usage, the History Filter for Storage Usage is configured with a StartTime of **`now() - 240 minutes`** (4 hours). This is because the Storage Usage metrics are produced every 4 hours. Again, if you know this to be incorrect, please raise an issue for it.  

### Configuration:

```toml
# Return performance data for all volumes in a Dell Storage Center endpoint.
[[inputs.dellstoragecenter]]
  ## IP address the Data Collector is listening on
  # ip_address = "192.168.192.168"

  ## The port number the Data Collector is listening on
  # port = 3033

  ## The username to log into the Data Collector with
  # username = "admin"

  ## The password to log into the Data Collector with
  # password = "admin"

  ## Version of the Dell API to use
  dell-api-version = 4.1

  ## Interval to poll for stats
  interval = "60s"
```

### Measurements & Fields:

- dellstoragecenter
  - tags:
    - scName (Storage Center name)
    - scVolume (Volume name, as shown in the Storage Center)
    - instanceId (Volume unique identifier, in the format 12345.1)
  - fields:
    - activeSpace (integer, bytes)
    - activeSpaceOnDisk (integer, bytes)
    - actualSpace (integer, bytes)
    - averageKbPerIo (integer, kilobytes)
    - configuredSpace (integer, bytes)
    - estimatedDataReductionSpaceSavings (integer, bytes)
    - estimatedDiskSpaceSavedByCompression (integer, bytes)
    - estimatedDiskSpaceSavedByDeduplicated (integer, bytes)
    - estimatedNonDeduplicatedToDuplicatedPageRatio (float64)
    - estimatedPercentCompressed (float64)
    - estimatedPercentDeduplicated (float)
    - estimatedUncompressedToCompressedPageRatio (float64)
    - freeSpace (integer)
    - instanceId (string)
    - instanceName (string)
    - ioPending (integer)
    - name (string)
    - objectType (string)
    - raidOverhead (integer, bytes)
    - readIops (integer)
    - readKbPerSecond (integer, kilobytes)
    - readLatency (integer, milliseconds)
    - replaySpace (integer, bytes)
    - savingsVsRaidTen (integer, bytes)
    - scName (string)
    - scSerialNumber (integer)
    - sharedSpace (integer, bytes)
    - snapshotOverheadOnDisk (integer, bytes)
    - time (timestamp)
    - totalDiskSpace (integer, bytes)
    - totalIops (integer)
    - totalKbPerSecond (integer, kilobytes)
    - writeIops (integer)
    - writeKbPerSecond (integer, kilobytes)
    - writeLatency (integer, milliseconds)
    - xferLatency (integer, milliseconds)


### Example Output:
```
dellstoragecenter,host=telegraf_proxy,instanceId=54321.1,scName=INT-SAN-01,scVolume=VMware_Datastore_01 averageKbPerIo=0i,instanceId="54321.1",ioPending=0i,readIops=0i,readKbPerSecond=0i,readLatency=0i,scName="INT-SAN-01",totalIops=0i,totalKbPerSecond=0i,writeIops=0i,writeKbPerSecond=0i,writeLatency=0i,xferLatency=0i 1547807198000000000
dellstoragecenter,host=telegraf_proxy,instanceId=54321.1,scName=INT-SAN-01,scVolume=VMware_Datastore_01 activeSpace=1607141949440i,activeSpaceOnDisk=0i,actualSpace=1607141949440i,configuredSpace=1610612736000i,estimatedDataReductionSpaceSavings=0i,estimatedDiskSpaceSavedByCompression=0i,estimatedDiskSpaceSavedByDeduplicated=0i,estimatedNonDeduplicatedToDuplicatedPageRatio=0,estimatedPercentCompressed=0,estimatedPercentDeduplicated=0,estimatedUncompressedToCompressedPageRatio=0,freeSpace=3470786560i,instanceId="54321.1",instanceName="VMware_Datastore_01",name="VMware_Datastore_01",objectType="ScVolumeStorageUsage",raidOverhead=803570974720i,replaySpace=0i,savingsVsRaidTen=803570974720i,scName="INT-SAN-01",scSerialNumber=54321i,sharedSpace=0i,snapshotOverheadOnDisk=0i,time="2019-01-18T08:00:01Z",totalDiskSpace=2410712924160i 1547798401000000000
dellstoragecenter,host=telegraf_proxy,instanceId=54321.2,scName=INT-SAN-01,scVolume=SQL_Data_01 averageKbPerIo=0i,instanceId="54321.2",ioPending=0i,readIops=0i,readKbPerSecond=0i,readLatency=0i,scName="INT-SAN-01",totalIops=0i,totalKbPerSecond=0i,writeIops=0i,writeKbPerSecond=0i,writeLatency=0i,xferLatency=0i 1547807198000000000
dellstoragecenter,host=telegraf_proxy,instanceId=54321.2,scName=INT-SAN-01,scVolume=SQL_Data_01 activeSpace=459657969664i,activeSpaceOnDisk=0i,actualSpace=459657969664i,configuredSpace=858993459200i,estimatedDataReductionSpaceSavings=0i,estimatedDiskSpaceSavedByCompression=0i,estimatedDiskSpaceSavedByDeduplicated=0i,estimatedNonDeduplicatedToDuplicatedPageRatio=0,estimatedPercentCompressed=0,estimatedPercentDeduplicated=0,estimatedUncompressedToCompressedPageRatio=0,freeSpace=399335489536i,instanceId="54321.2",instanceName="SQL_Data_01",name="SQL_Data_01",objectType="ScVolumeStorageUsage",raidOverhead=229828984832i,replaySpace=0i,savingsVsRaidTen=229828984832i,scName="INT-SAN-01",scSerialNumber=54321i,sharedSpace=0i,snapshotOverheadOnDisk=0i,time="2019-01-18T08:00:01Z",totalDiskSpace=689486954496i 1547798401000000000```
