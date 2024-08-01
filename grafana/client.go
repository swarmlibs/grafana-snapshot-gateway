package grafana

import (
	"encoding/json"
	"io"
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

func UnmarshalResponseBody(source io.Reader, target any) error {
	body, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func (g *GrafanaClient) SetBasicAuth(username, password string) {
	g.Username = username
	g.Password = password
}

func (g *GrafanaClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.Username, g.Password)
	return g.http.Do(req)
}
