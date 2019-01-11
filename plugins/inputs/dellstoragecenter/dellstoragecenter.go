package dellstoragecenter

import (
	"strconv"

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
	return "Return performance data from a Dell Storage Data Collector endpoint."
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
		scVolID := volume.InstanceID
		volumestats, err := connection.GetVolumeIoStats(scVolID)
		volumestat := volumestats[len(volumestats)-1]
		if err != nil {
			return err
		}

		tags := map[string]string{
			"scName":   volume.SCName,
			"scVolume": volume.Name,
		}

		fields := map[string]interface{}{
			"time":             volumestat.Time,
			"scName":           volumestat.SCName,
			"instanceId":       volumestat.InstanceID,
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

		acc.AddFields("dellstoragecenter", fields, tags)
	}

	return nil
}

func init() {
	c := Dellstoragecenter{}
	inputs.Add("dellstoragecenter", func() telegraf.Input {
		return &c
	})
}
