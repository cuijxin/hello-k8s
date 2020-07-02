package releases

import (
	helmtime "helm.sh/helm/v3/pkg/time"
)

type ReleaseInfo struct {
	Revision    int           `json:"revision"`
	Updated     helmtime.Time `json:"updated"`
	Status      string        `json:"status"`
	Chart       string        `json:"chart"`
	AppVersion  string        `json:"app_version"`
	Description string        `json:"description"`
}

type ReleaseHistory []ReleaseInfo

type ReleaseElement struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Revision     string `json:"revision"`
	Updated      string `json:"updated"`
	Status       string `json:"status"`
	Chart        string `json:"chart"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`
	Notes        string `json:"notes,omitempty"`
}

type ReleaseList []ReleaseElement

type ReleaseOptions struct {
	Values          string   `json:"values"`
	SetValues       []string `json:"set"`
	SetStringValues []string `json:"set_string"`
}
