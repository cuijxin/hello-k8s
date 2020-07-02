package charts

import (
	"strings"

	"helm.sh/helm/v3/pkg/chart"
)

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

type file struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func findReadme(files []*chart.File) (file *chart.File) {
	for _, file := range files {
		for _, n := range readmeFileNames {
			if strings.EqualFold(file.Name, n) {
				return file
			}
		}
	}
	return nil
}
