package runtime

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	. ".."
	containersfs "../fs"
)

const (
	terminalName = "gnome-terminal"
	engineName   = "./engine"
)

var runningContainers = make(map[string]*exec.Cmd)

// Starts a container's main process
func Start(containerInfo ContainerInfo) error {
	// CPU limit:    cgset -r cpu.shares=512 cpulimited
	// Memory limit: cgcreate -g memoty:/myGroup

	containerPath := containersfs.GetContainerPath(containerInfo.ContainerName)
	hostname := containerInfo.ContainerName
	cmd := exec.Command(terminalName, "--", engineName, containerInfo.ContainerName,
		containerPath, hostname, containerInfo.Address, containerInfo.CmdLine)
	err := cmd.Start()
	if err != nil {
		return err
	}

	runningContainers[containerInfo.ContainerName] = cmd

	return nil
}

// Stops all process in a container
func Stop(containerName string) error {
	cgroupsPath := fmt.Sprintf("/sys/fs/cgroup/pids/%s/cgroup.procs", containerName)
	data, err := ioutil.ReadFile(cgroupsPath)
	if err != nil {
		return err
	}

	pids := strings.Split(string(data), "\n")
	if len(pids) >= 2 {
		intPid, _ := strconv.Atoi(pids[1])
		syscall.Kill(intPid, syscall.SIGKILL)
	}
	return nil
}

// Stops a container and then restarts its main process
func Restart(containerInfo ContainerInfo) error {
	err := Stop(containerInfo.ContainerName)
	if err != nil {
		return err
	}
	err = Start(containerInfo)
	if err != nil {
		return err
	}
	return nil
}
