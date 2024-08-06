package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type GrafanaClient struct {
	Url      string
	Username string
	Password string
	http     *http.Client
	logger   log.Logger
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

func UnmarshalResponseBody(source io.ReadCloser, target any) error {
	defer source.Close()
	body, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func (g *GrafanaClient) SetLogger(logger log.Logger) {
	g.logger = logger
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
	req, err := http.NewRequest(method, g.Url+path, bytes.NewBuffer(bodybuf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (g *GrafanaClient) Do(req *http.Request) (*Response, error) {
	start := time.Now() // Start timer
	// Set basic auth if username and password are set
	if g.Username != "" && g.Password != "" {
		req.SetBasicAuth(g.Username, g.Password)
	}

	res, err := g.http.Do(req)

	stop := time.Now() // Stop timer
	latency := stop.Sub(start)

	msg := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	level.Info(g.logger).Log("msg", msg, "method", req.Method, "path", req.URL.Path, "status_code", res.StatusCode, "body_size", res.ContentLength, "latency", latency.String())

	return res, err
}
