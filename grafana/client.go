package grafana

import (
	"bytes"
	"fmt"
	"net/http"
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

func (g *GrafanaClient) CreateFolder(uid string, title string) error {
	body := []byte(fmt.Sprintf(`{"uid": "%s","title": "%s"}`, uid, title))
	req, _ := http.NewRequest("POST", g.Url+"/api/folders", bytes.NewBuffer(body))
	if _, err := g.Do(req); err != nil {
		return err
	}
	return nil
}

func (g *GrafanaClient) CreateDashboard(uid, dashboard map[string]interface{}) error {
	return nil
}

func (g *GrafanaClient) CreateSnapshot(dashboard map[string]interface{}, name string, expires int, external bool, key string, deleteKey string) error {
	return nil
}

func (g *GrafanaClient) GetSnapshot(key string) (map[string]interface{}, error) {
	return nil, nil
}

func (g *GrafanaClient) DeleteSnapshot(key string) error {
	return nil
}

func (g *GrafanaClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.Username, g.Password)
	return g.http.Do(req)
}
