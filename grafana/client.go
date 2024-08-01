package grafana

import (
	"bytes"
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

type Request = http.Request
type Response = http.Response

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

func (g *GrafanaClient) NewRequest(method string, path string, body any) (*Request, error) {
	bodybuf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, path, bytes.NewBuffer(bodybuf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (g *GrafanaClient) Do(req *http.Request) (*Response, error) {
	req.SetBasicAuth(g.Username, g.Password)
	return g.http.Do(req)
}
