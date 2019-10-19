package db

import (
	"fmt"
	"os/exec"
	"strings"
)

const serviceName = "mongod"

func isServiceRunning() bool {
	out, err := exec.Command("systemctl", "show", "-p", "SubState", "--value", serviceName).Output()
	if err != nil {
		panic(err)
	}
	state := strings.TrimSpace(string(out[:]))
	return state == "running"
}

// Starts the MongoDB service
func startService() {
	fmt.Printf("Starting DB service... ")
	serviceCommand("start").Run()
	fmt.Println()
}

// Stops the MongoDB service
func stopService() {
	fmt.Printf("Stopping DB service... ")
	serviceCommand("stop").Run()
	fmt.Println()
}

func serviceCommand(command string) *exec.Cmd {
	return exec.Command("service", serviceName, command)
}
