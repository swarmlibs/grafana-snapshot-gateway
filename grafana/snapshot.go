package grafana

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

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
