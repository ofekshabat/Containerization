package fs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. ".."
	. "../../packages"
	packagesdb "../../packages/db"
)

const (
	imagesPath = "images"
)

// Returns the path of an image in the filesystem
func GetImagePath(imageName string) string {
	return filepath.Join(imagesPath, imageName)
}

// Creates an image from a base image and package list
func Create(image Image, baseImageName string) error {
	imagePath := GetImagePath(image.Name)
	baseImagePath := GetImagePath(baseImageName)

	exec.Command("cp", "-r", baseImagePath, imagePath).Run()

	for _, packageName := range image.PackageNames {
		_package, err := packagesdb.GetPackageInfo(packageName)
		if err != nil {
			return err
		}

		scriptPath, err := writePackageInstallScript(_package, imagePath)
		if err != nil {
			return err
		}

		exec.Command("/bin/sh", scriptPath).Run()
	}

	return nil
}

func writePackageInstallScript(_package Package, imagePath string) (string, error) {
	scriptPath := filepath.Join(imagePath, "install.sh")
	f, err := os.Create(scriptPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	chrootLine := fmt.Sprintf("chroot %s\n", imagePath)
	_, err = f.WriteString(chrootLine)
	_, err = f.WriteString(_package.InstallScript)

	return scriptPath, nil
}

// Deletes an image
func Delete(imageName string) error {
	err := os.RemoveAll(GetImagePath(imageName))
	return err
}
