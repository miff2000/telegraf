package dellstoragecenter

import (
	"strconv"
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
  # interval = "60s"
	`
}

// Gather Collects Input data and pushes to Telegraf
func (d *Dellstoragecenter) Gather(acc telegraf.Accumulator) error {
	baseURL := "https://" + d.IPAddress + ":" + strconv.Itoa(d.Port)
	connection := newAPIConnection(baseURL, d.DellAPIVersion, d.Username, d.Password)
	err := connection.Login()
	if err != nil {
		return err
	}

	volumelist, err := connection.GetVolumeList()
	if err != nil {
		return err
	}
	for _, volume := range volumelist {

		volumestats, err := connection.GetVolumeIoStats(volume.InstanceID)
		if err != nil {
			return err
		}

		tags := map[string]string{
			"scName":     volume.SCName,
			"scVolume":   volume.Name,
			"instanceId": volume.InstanceID,
		}

		volumestat := volumestats[len(volumestats)-1]

		fields := map[string]interface{}{
			"readIops":         volumestat.ReadIOPS,
			"writeIops":        volumestat.WriteIOPS,
			"totalIops":        volumestat.TotalIOPS,
			"ioPending":        volumestat.IOPending,
			"readKbPerSecond":  volumestat.ReadKbPerSecond,
			"writeKbPerSecond": volumestat.WriteKbPerSecond,
			"totalKbPerSecond": volumestat.TotalKbPerSecond,
			"averageKbPerIo":   volumestat.AverageKbPerIO,
			"readLatency":      volumestat.ReadLatency,
			"writeLatency":     volumestat.WriteLatency,
			"xferLatency":      volumestat.XferLatency,
		}

		timestamp, err := time.Parse("2006-01-02T15:04:05Z", volumestat.Time)
		if err != nil {
			return err
		}

		acc.AddFields("dellstoragecenter", fields, tags, timestamp)
	}

	return nil
}

func init() {
	c := Dellstoragecenter{}
	inputs.Add("dellstoragecenter", func() telegraf.Input {
		return &c
	})
}
