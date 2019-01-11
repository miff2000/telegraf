package dellstoragecenter

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type (
	// ScVolume This is a list of returned ScVolumes
	ScVolume struct {
		Name       string `json:"name"`
		InstanceID string `json:"instanceId"`
		SCName     string `json:"scName"`
		Active     bool   `json:"active"`
		Status     string `json:"status"`
	}

	// ScVolumeIoStat These are the IO stats returned for a ScVolume
	ScVolumeIoStat struct {
		Time             string `json:"time"`
		SCName           string `json:"scName"`
		InstanceID       string `json:"instancdId"`
		InstanceName     string `json:"instanceName"`
		ReadIOPS         int    `json:"readIops"`
		WriteIOPS        int    `json:"writeIops"`
		TotalIOPS        int    `json:"totalIops"`
		IOPending        int    `json:"ioPending"`
		ReadKbPerSecond  int    `json:"readKbPerSecond"`
		WriteKbPerSecond int    `json:"writeKbPerSecond"`
		TotalKbPerSecond int    `json:"totalKbPerSecond"`
		AverageKbPerIO   int    `json:"averageKbPerIo"`
		ReadLatency      int    `json:"readLatency"`
		XferLatency      int    `json:"xferLatency"`
		WriteLatency     int    `json:"writeLatency"`
	}

	apiConnection struct {
		BaseURL    string
		APIVersion string
		Username   string
		Password   string
		loggedIn   bool
		client     *http.Client
	}

	historicalIOUsageRequest struct {
		HistoricalFilter historicalFilter `json:"HistoricalFilter"`
	}

	historicalFilter struct {
		MaxCountReturn int    `json:"MaxCountReturn"`
		StartTime      string `json:"StartTime"`
		UseCurrent     bool   `json:"UseCurrent"`
		UseEndOfDay    bool   `json:"UseEndOfDay"`
		UseStartOfDay  bool   `json:"UseStartOfDay"`
	}
)

func newAPIConnection(baseURL string, apiVersion string, username string, password string) *apiConnection {
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			ResponseHeaderTimeout: time.Duration(30 * time.Second),
		},
		Timeout: time.Duration(30 * time.Second),
	}

	return &apiConnection{
		BaseURL:    baseURL,
		APIVersion: apiVersion,
		Username:   username,
		Password:   password,
		client:     client,
	}
}

func createHistoricalFilter() *historicalIOUsageRequest {
	return &historicalIOUsageRequest{
		HistoricalFilter: historicalFilter{
			MaxCountReturn: 1,
			StartTime:      time.Now().Add(-15 * time.Minute).Format("2006-01-02T15:04:05-07:00"),
			UseCurrent:     true,
			UseEndOfDay:    false,
			UseStartOfDay:  false,
		},
	}
}

func (a *apiConnection) Login() error {
	if a.Username == "" || a.Password == "" {
		return errors.New("A username and password and password must be provided in telegraf.conf")
	}

	headers := map[string]string{}
	headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(a.Username+":"+a.Password))

	response, err := a.invoke("POST", "/ApiConnection/Login", nil, headers)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Unexpected response code:[%d]", response.StatusCode)
	}

	return nil
}

func (a *apiConnection) DecodeResponseBody(body io.ReadCloser, out interface{}) error {
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("Failed to read response body:[%s]", err)
	}

	// fmt.Printf("Got response: %s", string(bodyBytes))

	if bodyBytes != nil && len(bodyBytes) > 0 {
		return json.Unmarshal(bodyBytes, out)
	}

	return nil

}

func (a *apiConnection) GetVolumeList() ([]ScVolume, error) {
	response, err := a.post("/StorageCenter/ScVolume/GetList", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	ScVolume := []ScVolume{}
	err = a.DecodeResponseBody(response.Body, &ScVolume)
	if err != nil {
		return nil, err
	}

	return ScVolume, nil
}

func (a *apiConnection) GetVolumeIoStats(scVolID string) ([]ScVolumeIoStat, error) {
	body := createHistoricalFilter()

	response, err := a.post("/StorageCenter/ScVolume/"+scVolID+"/GetHistoricalIoUsage", body)
	if err != nil {
		return nil, err
	}

	ScVolumeIoStat := []ScVolumeIoStat{}
	err = a.DecodeResponseBody(response.Body, &ScVolumeIoStat)
	if err != nil {
		return nil, err
	}

	return ScVolumeIoStat, nil
}

func (a *apiConnection) get(resource string) (*http.Response, error) {
	return a.invoke("GET", resource, nil, nil)
}

func (a *apiConnection) post(resource string, body interface{}) (*http.Response, error) {
	return a.invoke("POST", resource, body, nil)
}

func (a *apiConnection) invoke(method string, resource string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := a.BaseURL + "/api/rest/" + strings.Trim(resource, "/")

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for headerName, headerValue := range headers {
			req.Header.Set(headerName, headerValue)
		}
	}

	// Set the version of the Dell API to use, if set. If not set, use the default 4.1
	if a.APIVersion != "" {
		req.Header.Set("x-dell-api-version", a.APIVersion)
	} else {
		req.Header.Set("x-dell-api-version", "4.1")
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		fmt.Printf("%v", response)
		return nil, fmt.Errorf("Unexpected response from Storage Center API:[%s]", response.Status)
	}

	return response, nil
}
