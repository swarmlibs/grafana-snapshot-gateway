package types

import (
	"encoding/json"
	"net/url"
)

type GrafanaDashboardSnapshot struct {
	Dashboard GrafanaDashboardModel `json:"dashboard"`
	Name      string                `json:"name,omitempty"`
	Expires   int                   `json:"expires,omitempty"`
	External  bool                  `json:"external,omitempty"`
	Key       string                `json:"key,omitempty"`
	DeleteKey string                `json:"deleteKey,omitempty"`
}

func (g *GrafanaDashboardSnapshot) SetKey(key string) {
	g.Key = key
	g.Dashboard.Set("uid", key)
}

func (g *GrafanaDashboardSnapshot) GetDashboardModel() (GrafanaDashboardModel, error) {
	// Make a copy of the dashboard
	d, err := json.Marshal(g.Dashboard)
	if err != nil {
		return nil, err
	}

	// Unmarshal the dashboard to a new struct
	var dashboard GrafanaDashboardModel
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

type GrafanaDashboardSnapshotResponse struct {
	Key       string `json:"key"`
	Url       string `json:"url"`
	DeleteKey string `json:"deleteKey"`
	DeleteUrl string `json:"deleteUrl"`
}

func (g *GrafanaDashboardSnapshotResponse) OverrideDeleteUrl(host string) (*GrafanaDashboardSnapshotResponse, error) {
	deleteUrl, err := url.Parse(g.DeleteUrl)
	if err != nil {
		return nil, err
	}
	deleteUrl.Host = host
	g.DeleteUrl = deleteUrl.String()
	return g, nil
}
