package grafana

import (
	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

func (g *GrafanaClient) CreateSnapshot(key string, snapshot types.GrafanaDashboardSnapshot) (*Response, error) {
	req, _ := g.NewRequest("POST", "/api/snapshots", snapshot)
	return g.Do(req)
}

func (g *GrafanaClient) GetSnapshot(key string) (*Response, error) {
	return nil, nil
}

func (g *GrafanaClient) DeleteSnapshot(key string) (*Response, error) {
	return nil, nil
}
