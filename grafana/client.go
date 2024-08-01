package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

type GrafanaClient struct {
	Url      string
	Username string
	Password string
	http     *http.Client
}

func NewGrafanaClient(url, username, password string) *GrafanaClient {
	return &GrafanaClient{
		Url:      url,
		Username: username,
		Password: password,
		http:     &http.Client{},
	}
}

func (g *GrafanaClient) SetBasicAuth(username, password string) {
	g.Username = username
	g.Password = password
}

func (g *GrafanaClient) CreateFolder(uid string, title string) (*http.Response, error) {
	body := []byte(fmt.Sprintf(`{"uid": "%s","title": "%s"}`, uid, title))
	req, _ := http.NewRequest("POST", g.Url+"/api/folders", bytes.NewBuffer(body))
	return g.Do(req)
}

func (g *GrafanaClient) CreateDashboard(uid string, dashboard types.GrafanaDashboard) (*http.Response, error) {
	body, _ := json.Marshal(dashboard)
	req, _ := http.NewRequest("POST", g.Url+"/api/dashboards/db", bytes.NewBuffer(body))
	return g.Do(req)
}

func (g *GrafanaClient) CreateSnapshot(key string, snapshot types.GrafanaDashboardSnapshot) (*http.Response, error) {
	body, _ := json.Marshal(snapshot)
	req, _ := http.NewRequest("POST", g.Url+"/api/snapshots", bytes.NewBuffer(body))
	return g.Do(req)
}

func (g *GrafanaClient) GetSnapshot(key string) (*http.Response, error) {
	return nil, nil
}

func (g *GrafanaClient) DeleteSnapshot(key string) (*http.Response, error) {
	return nil, nil
}

func (g *GrafanaClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.Username, g.Password)
	return g.http.Do(req)
}

func (g *GrafanaClient) ShouldBindJSON(source io.Reader, target any) error {
	body, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
