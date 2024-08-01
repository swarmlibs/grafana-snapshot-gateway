package grafana

import "encoding/json"

type GrafanaDashboard map[string]interface{}

type GrafanaDashboardSnapshot struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	Name      string                 `json:"name,omitempty"`
	Expires   int                    `json:"expires,omitempty"`
	External  bool                   `json:"external,omitempty"`
	Key       string                 `json:"key,omitempty"`
	DeleteKey string                 `json:"deleteKey,omitempty"`
}

func (g *GrafanaDashboardSnapshot) GetDashboardWithoutData() (GrafanaDashboard, error) {
	// Make a copy of the dashboard
	d, err := json.Marshal(g.Dashboard)
	if err != nil {
		return nil, err
	}

	// Unmarshal the dashboard to a new struct
	var dashboard GrafanaDashboard
	json.Unmarshal(d, &dashboard)

	// Remove snapshotData from dashboard
	if panels, ok := dashboard["panels"].([]interface{}); ok {
		for _, panel := range panels {
			if panelMap, ok := panel.(map[string]interface{}); ok {
				delete(panelMap, "snapshotData")
			}
		}
	}

	return dashboard, nil
}
