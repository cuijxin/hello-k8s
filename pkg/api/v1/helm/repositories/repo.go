package repositories

type RepoChartElement struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

type RepoChartList []RepoChartElement
