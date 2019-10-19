package fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. ".."
)

const (
	containersPath = "containers"
)

// Returns the path of a container in the filesystem
func GetContainerPath(containerName string) string {
	return filepath.Join(containersPath, containerName)
}

// Returns the path of a container's config JSON file
func GetConfigPath(containerName string) string {
	return filepath.Join(GetContainerPath(containerName), "config.json")
}

// Creates a container in the filesystem by copying an image and
// writing a config JSON file
func Create(containerInfo ContainerInfo, imagePath string) error {
	containerPath := GetContainerPath(containerInfo.ContainerName)
	err := exec.Command("cp", "-r", imagePath, containerPath).Run()
	return err
}

// Renames a container filesystem
func Rename(oldName string, newName string) error {
	err := os.Rename(GetContainerPath(oldName), GetContainerPath(newName))
	return err
}

// Deletes a container filesystem
func Delete(containerName string) error {
	err := os.RemoveAll(GetContainerPath(containerName))
	return err
}

// Imports a container from a tarball (.tar.gz file)
func Import(containerName string, path string) (ContainerInfo, error) {
	var containerInfo ContainerInfo

	containerPath := GetContainerPath(containerName)
	os.Mkdir(containerPath, 0755)
	cmd := exec.Command("tar", "xzf", path)
	cmd.Dir = containerPath
	cmd.Run()

	configPath := GetConfigPath(containerName)
	configFile, err := os.Open(configPath)
	if err != nil {
		return containerInfo, err
	}
	configData, _ := ioutil.ReadAll(configFile)
	configFile.Close()

	err = json.Unmarshal([]byte(configData), &containerInfo)
	if err != nil {
		return containerInfo, err
	}
	os.Remove(configPath)

	containerInfo.ContainerName = containerName
	return containerInfo, nil
}

// Exports a container into a tarball (.tar.gz file)
func Export(containerInfo ContainerInfo, destinationPath string) error {
	// Write config file
	configPath := GetConfigPath(containerInfo.ContainerName)
	configJson, _ := json.Marshal(containerInfo)
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(configPath)

	workingDir, err := os.Getwd()

	command := fmt.Sprintf("tar czf %s *", destinationPath)
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = filepath.Join(workingDir, GetContainerPath(containerInfo.ContainerName))
	err = cmd.Run()

	return nil
}
