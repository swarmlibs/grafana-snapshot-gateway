package grafana

import (
	"bytes"
	"fmt"
	"net/http"
)

func (g *GrafanaClient) CreateFolder(uid string, title string) (*http.Response, error) {
	body := []byte(fmt.Sprintf(`{"uid": "%s","title": "%s"}`, uid, title))
	req, _ := http.NewRequest("POST", g.Url+"/api/folders", bytes.NewBuffer(body))
	return g.Do(req)
}
