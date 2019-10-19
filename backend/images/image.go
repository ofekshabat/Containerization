package images

// Image contains a filesystem that can be imported
// into a container
type Image struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	PackageNames []string `json:"packageNames"`
}
