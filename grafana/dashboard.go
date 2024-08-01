package grafana

import (
	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

func (g *GrafanaClient) CreateDashboard(uid string, dashboard types.GrafanaDashboard) (*Response, error) {
	req, _ := g.NewRequest("POST", "/api/dashboards/db", dashboard)
	return g.Do(req)
}
