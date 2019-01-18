# DellStorageCenter Input Plugin

The dellstoragecenter input plugin gathers volume performance data

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
    - readIops (integer)
    - writeIops (integer)
    - totalIops (integer)
    - ioPending (integer)
    - readKbPerSecond (integer, kilobytes)
    - writeKbPerSecond (integer, kilobytes)
    - totalKbPerSecond (integer, kilobytes)
    - averageKbPerIo (integer, kilobytes)
    - readLatency (integer, milliseconds)
    - writeLatency (integer, milliseconds)
    - xferLatency (integer, milliseconds)
    - activeSpace (integer, bytes)
    - activeSpaceOnDisk (integer, bytes)
    - actualSpace (integer, bytes)
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
    - name (string)
    - objectType (string)
    - raidOverhead (integer, bytes)
    - replaySpace (integer, bytes)
    - savingsVsRaidTen (integer, bytes)
    - scName (string)
    - scSerialNumber (integer)
    - sharedSpace (integer, bytes)
    - snapshotOverheadOnDisk (integer, bytes)
    - time (timestamp)
    - totalDiskSpace (integer, bytes)

#### `readIops`, `writeIops` & `totalIops`:

* Read IO/Second for the Volume
* Write IO/Second for the Volume
* Total IO/Second for the Volume	

#### `ioPending`

* IO that is pending for the Volume

#### `readKbPerSecond`, `writeKbPerSecond` & `totalKbPerSecond`

* Read KB/Second for the Volume	
* Write KB/Second for the Volume	
* Total KB/Second for the Volume	

#### `averageKbPerIo`

* The average IO size in KB for the Volume

#### `readLatency`, `writeLatency` & `xferLatency`

* Read Latency (in microseconds) for the Volume
* Write Latency (in microseconds) for the Volume
* Transfer Latency for the Volume

#### `activeSpace`, ``, ``, ``, ``,  




### Sample Queries:
#### Calculate percent IO utilization per disk and host:
* **todo**
```
SELECT non_negative_derivative(last("io_time"),1ms) FROM "diskio" WHERE time > now() - 30m GROUP BY "host","name",time(60s)
```

#### Calculate average queue depth:
* **todo**
`iops_in_progress` will give you an instantaneous value. This will give you the average between polling intervals.
```
SELECT non_negative_derivative(last("weighted_io_time",1ms)) from "diskio" WHERE time > now() - 30m GROUP BY "host","name",time(60s)
```

### Example Output:

* **todo**
```
diskio,name=sda weighted_io_time=8411917i,read_time=7446444i,write_time=971489i,io_time=866197i,write_bytes=5397686272i,iops_in_progress=0i,reads=2970519i,writes=361139i,read_bytes=119528903168i 1502467254359000000
diskio,name=sda1 reads=2149i,read_bytes=10753536i,write_bytes=20697088i,write_time=346i,weighted_io_time=505i,writes=2110i,read_time=161i,io_time=208i,iops_in_progress=0i 1502467254359000000
diskio,name=sda2 reads=2968279i,writes=359029i,write_bytes=5376989184i,iops_in_progress=0i,weighted_io_time=8411250i,read_bytes=119517334528i,read_time=7446249i,write_time=971143i,io_time=866010i 1502467254359000000
diskio,name=sdb writes=99391856i,write_time=466700894i,io_time=630259874i,weighted_io_time=4245949844i,reads=2750773828i,read_bytes=80667939499008i,write_bytes=6329347096576i,read_time=3783042534i,iops_in_progress=2i 1502467254359000000
diskio,name=centos/root read_time=7472461i,write_time=950014i,iops_in_progress=0i,weighted_io_time=8424447i,writes=298543i,read_bytes=119510105088i,io_time=837421i,reads=2971769i,write_bytes=5192795648i 1502467254359000000
diskio,name=centos/var_log reads=1065i,writes=69711i,read_time=1083i,write_time=35376i,read_bytes=6828032i,write_bytes=184193536i,io_time=29699i,iops_in_progress=0i,weighted_io_time=36460i 1502467254359000000
diskio,name=postgresql/pgsql write_time=478267417i,io_time=631098730i,iops_in_progress=2i,weighted_io_time=4263637564i,reads=2750777151i,writes=110044361i,read_bytes=80667939288064i,write_bytes=6329347096576i,read_time=3784499336i 1502467254359000000
```
