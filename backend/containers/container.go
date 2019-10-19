package containers

// ContainerInfo contains all the information needed to run
// a container in the system
type ContainerInfo struct {
	ContainerName string `json:"containerName"`
	BaseImageName string `json:"baseImageName"`
	CmdLine       string `json:"cmdLine"`
	State         string `json:"state"`     // "stopped" or "running"
	MaxCPU        int    `json:"maxCpu"`    // Percent
	MaxMemory     int    `json:"maxMemory"` // Mbs
	MaxPids       int    `json:"maxPids"`
	Address       string `json:"address"`
}
