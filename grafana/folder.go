package grafana

func (g *GrafanaClient) CreateFolder(uid string, title string) (*Response, error) {
	body := interface{}(map[string]interface{}{"uid": uid, "title": title})
	req, _ := g.NewRequest("POST", "/api/folders", body)
	return g.Do(req)
}

func (g *GrafanaClient) DeleteFolder(uid string) (*Response, error) {
	req, _ := g.NewRequest("DELETE", "/api/folders/"+uid, nil)
	return g.Do(req)
}
