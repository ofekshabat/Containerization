package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const cgroupsPath = "/sys/fs/cgroup"

func writeControlGroup(containerName string) {
	// ioPath := getCgroupPath("io", containerName)
	// memoryPath := getCgroupPath("memory", containerName)
	pidsPath := getCgroupPath("pids", containerName)
	// cpuPath := getCgroupPath("cpu", containerName)

	os.Mkdir(pidsPath, 0755)

	// Limit the number of process inside the cgroups
	// To check run ":() { : | : & }; :" and we can look the tree at the host with ps fax
	ioutil.WriteFile(filepath.Join(pidsPath, "pids.max"), []byte("20"), 0700)

	// Removes the new cgroup in place after the container exits
	ioutil.WriteFile(filepath.Join(pidsPath, "notify_on_release"), []byte("1"), 0700)
	ioutil.WriteFile(filepath.Join(pidsPath, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
}

func getCgroupPath(subsystem string, containerName string) string {
	return filepath.Join(cgroupsPath, subsystem, containerName)
}
