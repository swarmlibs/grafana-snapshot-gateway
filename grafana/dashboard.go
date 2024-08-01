package grafana

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
)

func (g *GrafanaClient) CreateDashboard(uid string, dashboard types.GrafanaDashboard) (*http.Response, error) {
	body, _ := json.Marshal(dashboard)
	req, _ := http.NewRequest("POST", g.Url+"/api/dashboards/db", bytes.NewBuffer(body))
	return g.Do(req)
}
