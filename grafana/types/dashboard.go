package types

func NewGrafanaDashboard() *GrafanaDashboard {
	return &GrafanaDashboard{}
}

type GrafanaDashboard struct {
	Dashboard GrafanaDashboardModel `json:"dashboard"`
	FolderUid string                `json:"folderUid,omitempty"`
	Message   string                `json:"message,omitempty"`
	Overwrite bool                  `json:"overwrite,omitempty"`
}

func (g *GrafanaDashboard) SetFolderUid(folderUid string) {
	g.FolderUid = folderUid
}

func (g *GrafanaDashboard) SetMessage(message string) {
	g.Message = message
}

func (g *GrafanaDashboard) SetOverwrite(overwrite bool) {
	g.Overwrite = overwrite
}

func (g *GrafanaDashboard) SetDashboardModel(dashboard GrafanaDashboardModel) {
	g.Dashboard = dashboard
}

type GrafanaDashboardModel map[string]interface{}

func (g *GrafanaDashboardModel) Set(key string, value string) {
	(*g)[key] = value
}

func (g *GrafanaDashboardModel) Get(key string) string {
	if value, ok := (*g)[key].(string); ok {
		return value
	}
	return ""
}
