package grafana

import (
	"bytes"
	"encoding/json"

	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

func (g *GrafanaClient) CreateDashboard(uid string, dashboard types.GrafanaDashboard) (*Response, error) {
	body, _ := json.Marshal(dashboard)
	req, _ := g.NewRequest("POST", g.Url+"/api/dashboards/db", bytes.NewBuffer(body))
	return g.Do(req)
}
