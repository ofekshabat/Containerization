package packages

// Package contains information about a software package
// that can be installed on an image
type Package struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	InstallScript string `json:"installScript"`
}
