package grafana

import (
	"bytes"
	"fmt"
)

func (g *GrafanaClient) CreateFolder(uid string, title string) (*Response, error) {
	body := []byte(fmt.Sprintf(`{"uid": "%s","title": "%s"}`, uid, title))
	req, _ := g.NewRequest("POST", g.Url+"/api/folders", bytes.NewBuffer(body))
	return g.Do(req)
}
