package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"bps-client-go/pkg/models"
	"bps-client-go/pkg/operations"
)

const (
	ClientVersion = "11.0"
	APIVersion    = "v2"
)

type BPS struct {
	Host          string
	User          string
	Password      string
	SessionID     string
	Client        *resty.Client
	ClientVersion []int

	ServerVersions map[string]interface{}
	CheckVersion   bool

	PrintRequests    bool
	ProfilingEnabled bool
	ProfilingData    map[string]map[string][]float64
	profilingMutex   sync.RWMutex

	Results        *models.DataModelProxy
	Capture        *models.DataModelProxy
	Administration *models.DataModelProxy
	Topology       *models.DataModelProxy
	LoadProfile    *models.DataModelProxy
	Network        *models.DataModelProxy
	EvasionProfile *models.DataModelProxy
	Remote         *models.DataModelProxy

	StrikeList        *operations.StrikeListOps
	Strikes           *operations.StrikesOps
	Reports           *operations.ReportsOps
	AppProfile        *operations.AppProfileOps
	Superflow         *operations.SuperflowOps
	TestModel         *operations.TestModelOps
	Statistics        *operations.StatisticsOps
	AdministrationOps *operations.AdministrationOps
	TopologyOps       *operations.TopologyOps
	LoadProfileOps    *operations.LoadProfileOps
	NetworkOps        *operations.NetworkOps
	EvasionProfileOps *operations.EvasionProfileOps
	RemoteOps         *operations.RemoteOps
}

func NewBPS(host, user, password string, checkVersion bool) *BPS {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(30 * time.Second)

	bps := &BPS{
		Host:            host,
		User:            user,
		Password:        password,
		Client:          client,
		ClientVersion:   parseVersion(ClientVersion),
		CheckVersion:    checkVersion,
		PrintRequests:   false,
		ProfilingEnabled: false,
		ProfilingData:   make(map[string]map[string][]float64),
	}

	bps.Results = models.NewDataModelProxy(bps, "results", "")
	bps.Capture = models.NewDataModelProxy(bps, "capture", "")
	bps.Administration = models.NewDataModelProxy(bps, "administration", "")
	bps.Topology = models.NewDataModelProxy(bps, "topology", "")
	bps.LoadProfile = models.NewDataModelProxy(bps, "loadProfile", "")
	bps.Network = models.NewDataModelProxy(bps, "network", "")
	bps.EvasionProfile = models.NewDataModelProxy(bps, "evasionProfile", "")
	bps.Remote = models.NewDataModelProxy(bps, "remote", "")

	bps.StrikeList        = &operations.StrikeListOps{Client: bps}
	bps.Strikes           = &operations.StrikesOps{Client: bps}
	bps.Reports           = &operations.ReportsOps{Client: bps}
	bps.AppProfile        = &operations.AppProfileOps{Client: bps}
	bps.Superflow         = &operations.SuperflowOps{Client: bps}
	bps.TestModel         = &operations.TestModelOps{Client: bps}
	bps.Statistics        = &operations.StatisticsOps{Client: bps}
	bps.AdministrationOps = &operations.AdministrationOps{Client: bps}
	bps.TopologyOps       = &operations.TopologyOps{Client: bps}
	bps.LoadProfileOps    = &operations.LoadProfileOps{Client: bps}
	bps.NetworkOps        = &operations.NetworkOps{Client: bps}
	bps.EvasionProfileOps = &operations.EvasionProfileOps{Client: bps}
	bps.RemoteOps         = &operations.RemoteOps{Client: bps}

	return bps
}

func (b *BPS) Login() (map[string]interface{}, error) {
	if err := b.connect(); err != nil {
		return nil, err
	}

	loginData := map[string]interface{}{
		"username":  b.User,
		"password":  b.Password,
		"sessionId": b.SessionID,
	}

	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(loginData).
		Post(fmt.Sprintf("https://%s/bps/api/v2/core/auth/login", b.Host))
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("login failed: status %d, %s", resp.StatusCode(), resp.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	b.ServerVersions = result

	if b.CheckVersion {
		if err := b.validateVersion(); err != nil {
			b.Logout()
			return nil, err
		}
	}
	return b.ServerVersions, nil
}

func (b *BPS) Logout() error {
	data := map[string]interface{}{
		"username":  b.User,
		"password":  b.Password,
		"sessionId": b.SessionID,
	}
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(fmt.Sprintf("https://%s/bps/api/v2/core/auth/logout", b.Host))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("logout failed: status %d, %s", resp.StatusCode(), resp.String())
	}
	b.disconnect()
	return nil
}

func (b *BPS) connect() error {
	auth := map[string]interface{}{"username": b.User, "password": b.Password}
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(auth).
		Post(fmt.Sprintf("https://%s/bps/api/v1/auth/session", b.Host))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("connection failed: %d, %s", resp.StatusCode(), resp.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}
	sess, ok := res["sessionId"].(string)
	if !ok {
		return fmt.Errorf("invalid sessionId")
	}
	key, ok := res["apiKey"].(string)
	if !ok {
		return fmt.Errorf("invalid apiKey")
	}
	b.SessionID = sess
	b.Client.SetHeaders(map[string]string{"sessionId": b.SessionID, "X-API-KEY": key})
	fmt.Printf("Successfully connected to %s.\n", b.Host)
	return nil
}

func (b *BPS) disconnect() {
	resp, err := b.Client.R().
		Delete(fmt.Sprintf("https://%s/bps/api/v1/auth/session", b.Host))
	if err == nil && resp.StatusCode() == 204 {
		b.SessionID = ""
		b.Client.SetHeaders(map[string]string{"sessionId": "", "X-API-KEY": ""})
	}
}

func (b *BPS) Get(path string, depth *int, params map[string]string) (interface{}, error) {
	if b.ProfilingEnabled {
		start := time.Now()
		defer b.recordTiming("Get", path, time.Since(start))
	}
	url := fmt.Sprintf("https://%s/bps/api/v2/core%s", b.Host, path)
	req := b.Client.R()
	if depth != nil {
		req.SetQueryParam("responseDepth", strconv.Itoa(*depth))
	}
	for k, v := range params {
		req.SetQueryParam(k, v)
	}
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 || resp.StatusCode() == 204 {
		return b.parseJSON(resp.Body())
	}
	return nil, fmt.Errorf("GET failed: %d, %s", resp.StatusCode(), resp.String())
}

func (b *BPS) Post(path string, data interface{}) (interface{}, error) {
	if b.ProfilingEnabled {
		start := time.Now()
		defer b.recordTiming("Post", path, time.Since(start))
	}
	url := fmt.Sprintf("https://%s/bps/api/v2/core/%s", b.Host, path)
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 || resp.StatusCode() == 202 || resp.StatusCode() == 204 {
		return b.parseJSON(resp.Body())
	}
	if resp.StatusCode() == 400 {
		return nil, fmt.Errorf("bad request: %s", resp.String())
	}
	return nil, fmt.Errorf("POST failed: %d, %s", resp.StatusCode(), resp.String())
}

func (b *BPS) Put(path string, value interface{}) error {
	url := fmt.Sprintf("https://%s/bps/api/v2/core/%s", b.Host, path)
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(value).
		Put(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("PUT failed: %s", resp.String())
	}
	return nil
}

func (b *BPS) Patch(path string, value interface{}) error {
	url := fmt.Sprintf("https://%s/bps/api/v2/core/%s", b.Host, path)
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(value).
		Patch(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("PATCH failed: %s", resp.String())
	}
	return nil
}

func (b *BPS) Delete(path string) (interface{}, error) {
	url := fmt.Sprintf("https://%s/bps/api/v2/core/%s", b.Host, path)
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		Delete(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 || resp.StatusCode() == 204 {
		return b.parseJSON(resp.Body())
	}
	return nil, fmt.Errorf("DELETE failed: %d, %s", resp.StatusCode(), resp.String())
}

func (b *BPS) Export(path, file string, params map[string]interface{}) error {
	url := fmt.Sprintf("https://%s/bps/api/v2/core/%s", b.Host, path)
	params["filepath"] = file
	resp, err := b.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() == 200 || resp.StatusCode() == 204 {
		return b.downloadFile(fmt.Sprintf("https://%s%s", b.Host, resp.String()), file)
	}
	return fmt.Errorf("EXPORT failed: %d, %s", resp.StatusCode(), resp.String())
}

func (b *BPS) Import(path, filename string, params map[string]interface{}) (interface{}, error) {
	url := fmt.Sprintf("https://%s/bps/api/v2/core/%s", b.Host, path)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	resp, err := b.Client.R().
		SetFileReader("file", filename, file).
		SetFormData(map[string]string{"fileInfo": fmt.Sprintf("%v", params)}).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 || resp.StatusCode() == 204 {
		return b.parseJSON(resp.Body())
	}
	return nil, fmt.Errorf("IMPORT failed: %d, %s", resp.StatusCode(), resp.String())
}

func (b *BPS) downloadFile(url, file string) error {
	resp, err := b.Client.R().SetHeader("Content-Type", "application/json").Get(url)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, strings.NewReader(resp.String()))
	return err
}

func (b *BPS) parseJSON(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return string(data), nil
	}
	return res, nil
}

func parseVersion(v string) []int {
	parts := strings.Split(v, ".")
	out := make([]int, 2)
	for i, p := range parts {
		if i >= 2 {
			break
		}
		if val, err := strconv.Atoi(p); err == nil {
			out[i] = val
		}
	}
	return out
}

func compareVersions(a, b []int) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] > b[i] {
			return 1
		}
		if a[i] < b[i] {
			return -1
		}
	}
	return 0
}

func (b *BPS) validateVersion() error {
	if b.ServerVersions == nil {
		return nil
	}
	verStr, ok := b.ServerVersions["apiServer"].(string)
	if !ok {
		return nil
	}
	serverVer := parseVersion(verStr)
	cmp := compareVersions(serverVer, b.ClientVersion)
	if cmp > 0 {
		return fmt.Errorf("client version is older than server version")
	}
	if cmp < 0 {
		fmt.Println("Warning: client version is newer than server version")
	}
	return nil
}

func (b *BPS) EnableProfiling(e bool) {
	b.profilingMutex.Lock()
	defer b.profilingMutex.Unlock()
	if e {
		b.ProfilingData = make(map[string]map[string][]float64)
	}
	b.ProfilingEnabled = e
}

func (b *BPS) recordTiming(method, args string, dur time.Duration) {
	if !b.ProfilingEnabled {
		return
	}
	b.profilingMutex.Lock()
	defer b.profilingMutex.Unlock()
	if b.ProfilingData[method] == nil {
		b.ProfilingData[method] = make(map[string][]float64)
	}
	b.ProfilingData[method][args] = append(b.ProfilingData[method][args], dur.Seconds())
}

func (b *BPS) PrintVersions() {
	server := "N/A"
	if b.ServerVersions != nil {
		if v, ok := b.ServerVersions["apiServer"].(string); ok {
			server = v
		}
	}
	fmt.Printf("Client Version: %s\nServer Version: %s\n", ClientVersion, server)
}

func (b *BPS) PrintProfilingData() {
	b.profilingMutex.RLock()
	defer b.profilingMutex.RUnlock()
	for method, calls := range b.ProfilingData {
		for args, times := range calls {
			count := len(times)
			if count == 0 {
				continue
			}
			var sum, min, max float64
			min, max = times[0], times[0]
			for _, t := range times {
				sum += t
				if t < min {
					min = t
				}
				if t > max {
					max = t
				}
			}
			fmt.Printf("%s %s: n=%d, avg=%.6f, min=%.6f, max=%.6f\n", method, args, count, sum/float64(count), min, max)
		}
	}
}