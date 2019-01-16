package dellstoragecenter

import (
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type (
	// Dellstoragecenter Defines the interface for connecting to a Dell Storage Center
	Dellstoragecenter struct {
		IPAddress      string `toml:"ip_address"`
		Port           int    `toml:"port"`
		Username       string `toml:"username"`
		Password       string `toml:"password"`
		DellAPIVersion string `toml:"dell-api-version"`
		BaseURL        string
		connection     *apiConnection
	}
)

// Description Returns the description of the plugin
func (d *Dellstoragecenter) Description() string {
	return "Return performance data for all volumes in a Dell Storage Center endpoint."
}

// SampleConfig Returns an sample configuration for the Telegraf config file
func (d *Dellstoragecenter) SampleConfig() string {
	return `
  ## IP address the Data Collector is listening on
  # ip_address = "192.168.192.168"

  ## The port number the Data Collector is listening on
  # port = 3033

  ## The username to log into the Data Collector with
  # username = "admin"

  ## The password to log into the Data Collector with
  # password = "admin"

  ## Version of the Dell API to use
  # dell-api-version = 4.1

  ## Interval to poll for stats
  interval = "60s"
	`
}

// Gather Collects Input data and pushes to Telegraf
func (d *Dellstoragecenter) Gather(acc telegraf.Accumulator) error {
	baseURL := "https://" + d.IPAddress + ":" + strconv.Itoa(d.Port)
	apiConn := newAPIConnection(baseURL, d.DellAPIVersion, d.Username, d.Password)
	err := apiConn.Login()
	if err != nil {
		return err
	}

	scVolumeList, err := apiConn.GetVolumeList()
	if err != nil {
		return err
	}
	for _, scVolume := range scVolumeList {

		tags := map[string]string{
			"scName":     scVolume.SCName,
			"scVolume":   scVolume.Name,
			"instanceId": scVolume.InstanceID,
		}

		fields, timestamp, err := d.gatherIoUsageStat(apiConn, scVolume)
		if err != nil {
			return err
		}
		acc.AddFields("dellstoragecenter", fields, tags, timestamp)

		fields, timestamp, err = d.gatherStorageUsageStat(apiConn, scVolume)
		if err != nil {
			return err
		}
		acc.AddFields("dellstoragecenter", fields, tags, timestamp)
	}

	return nil
}

func (d *Dellstoragecenter) gatherIoUsageStat(apiConn *apiConnection, scVolume ScVolume) (map[string]interface{}, time.Time, error) {

	volumeIoUsageStats, err := apiConn.GetVolumeIoUsageStats(scVolume.InstanceID)
	if err != nil {
		return nil, time.Time{}, err
	}

	volStat := volumeIoUsageStats[len(volumeIoUsageStats)-1]

	fields := map[string]interface{}{
		"readIops":         volStat.ReadIOPS,
		"writeIops":        volStat.WriteIOPS,
		"totalIops":        volStat.TotalIOPS,
		"ioPending":        volStat.IOPending,
		"readKbPerSecond":  volStat.ReadKbPerSecond,
		"writeKbPerSecond": volStat.WriteKbPerSecond,
		"totalKbPerSecond": volStat.TotalKbPerSecond,
		"averageKbPerIo":   volStat.AverageKbPerIO,
		"readLatency":      volStat.ReadLatency,
		"writeLatency":     volStat.WriteLatency,
		"xferLatency":      volStat.XferLatency,
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", volStat.Time)
	if err != nil {
		return nil, time.Time{}, err
	}

	return fields, timestamp, nil
}

func (d *Dellstoragecenter) gatherStorageUsageStat(apiConn *apiConnection, scVolume ScVolume) (map[string]interface{}, time.Time, error) {

	volumeStorageUsageStats, err := apiConn.GetVolumeStorageUsageStats(scVolume.InstanceID)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(volumeStorageUsageStats) < 1 {
		return map[string]interface{}{}, time.Time{}, nil
	}

	volStat := volumeStorageUsageStats[len(volumeStorageUsageStats)-1]

	fields := map[string]interface{}{
		"activeSpace":                                   bytesStringToInt(volStat.ActiveSpace),
		"activeSpaceOnDisk":                             bytesStringToInt(volStat.ActiveSpaceOnDisk),
		"actualSpace":                                   bytesStringToInt(volStat.ActualSpace),
		"configuredSpace":                               bytesStringToInt(volStat.ConfiguredSpace),
		"estimatedDataReductionSpaceSavings":            bytesStringToInt(volStat.EstimatedDataReductionSpaceSavings),
		"estimatedDiskSpaceSavedByCompression":          bytesStringToInt(volStat.EstimatedDiskSpaceSavedByCompression),
		"estimatedDiskSpaceSavedByDeduplicated":         bytesStringToInt(volStat.EstimatedDiskSpaceSavedByDeduplication),
		"estimatedNonDeduplicatedToDuplicatedPageRatio": volStat.EstimatedNonDeduplicatedToDuplicatedPageRatio,
		"estimatedPercentCompressed":                    volStat.EstimatedPercentCompressed,
		"estimatedPercentDeduplicated":                  volStat.EstimatedPercentDeduplicated,
		"estimatedUncompressedToCompressedPageRatio":    volStat.EstimatedUncompressedToCompressedPageRatio,
		"freeSpace":                                     bytesStringToInt(volStat.FreeSpace),
		"instanceId":                                    volStat.InstanceID,
		"instanceName":                                  volStat.InstanceName,
		"name":                                          volStat.Name,
		"objectType":                                    volStat.ObjectType,
		"raidOverhead":                                  bytesStringToInt(volStat.RaidOverhead),
		"replaySpace":                                   bytesStringToInt(volStat.ReplaySpace),
		"savingsVsRaidTen":                              bytesStringToInt(volStat.SavingsVsRaidTen),
		"scName":                                        volStat.ScName,
		"scSerialNumber":                                volStat.ScSerialNumber,
		"sharedSpace":                                   bytesStringToInt(volStat.SharedSpace),
		"snapshotOverheadOnDisk":                        bytesStringToInt(volStat.SnapshotOverheadOnDisk),
		"time":                                          volStat.Time,
		"totalDiskSpace":                                bytesStringToInt(volStat.TotalDiskSpace),
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", volStat.Time)
	if err != nil {
		return nil, time.Time{}, err
	}

	return fields, timestamp, nil
}

// bytesStringToInt converts a string in the form "13245 Bytes" to an int64 like 12346
func bytesStringToInt(byteString string) int64 {
	// strip " Byte" from the string
	str := strings.TrimRight(byteString, " Bytes")
	integer64, _ := strconv.ParseInt(str, 10, 64)

	//fmt.Printf("Converted %s to %d\n", byteString, integer64)
	return integer64
}

func init() {
	c := Dellstoragecenter{}
	inputs.Add("dellstoragecenter", func() telegraf.Input {
		return &c
	})
}
