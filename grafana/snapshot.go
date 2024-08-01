package grafana

import (
	"bytes"
	"encoding/json"

	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

func (g *GrafanaClient) CreateSnapshot(key string, snapshot types.GrafanaDashboardSnapshot) (*Response, error) {
	body, _ := json.Marshal(snapshot)
	req, _ := g.NewRequest("POST", g.Url+"/api/snapshots", bytes.NewBuffer(body))
	return g.Do(req)
}

func (g *GrafanaClient) GetSnapshot(key string) (*Response, error) {
	return nil, nil
}

func (g *GrafanaClient) DeleteSnapshot(key string) (*Response, error) {
	return nil, nil
}
