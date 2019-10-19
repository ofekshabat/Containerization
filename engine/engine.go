package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"github.com/docker/docker/pkg/reexec"
)

var duration = time.Duration(10)*time.Second
// ./engine container-name rootfs hostname address command args
// ./engine container1 /tmp/example/rootfs my-hostname 10.10.0.2/24 sh

func main() {
	netsetgoPath := "netsetgo"
	containerName := os.Args[1]
	rootfsPath := os.Args[2]
	hostname := os.Args[3]
	containerAddress := os.Args[4]
	command := os.Args[5:]

	args := append([]string{"nsInitialisation", containerName, rootfsPath, hostname}, command...)
	cmd := reexec.Command(args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | // Mount namespace
			syscall.CLONE_NEWUTS | // UTS IPC namespace
			syscall.CLONE_NEWIPC | // System V IPC namespace
			syscall.CLONE_NEWPID | // PID namespace
			syscall.CLONE_NEWNET | // Network namespace
			syscall.CLONE_NEWUSER, // User namespace

		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the reexec.Command - %s\n", err)
		time.Sleep(duration)
		os.Exit(1)
	}

	// run netsetgo using default args
	// note that netsetgo must be owned by root with the setuid bit set
	pid := fmt.Sprintf("%d", cmd.Process.Pid)	
	//netsetgoCmd := exec.Command(netsetgoPath, "-pid", pid,
	//	"-containerAddress", containerAddress)
	
	netsetgoCmd := exec.Command(netsetgoPath,
	 	"-pid", pid,
	 	"-bridgeAddress", "10.10.10.1/24",
	 	"-bridgeName", "brg0",
	 	"-containerAddress", containerAddress,
	 	"-vethNamePrefix", "veth")

	if err := netsetgoCmd.Run(); err != nil {
		fmt.Printf("Error running netsetgo - %s\n", err)
		time.Sleep(duration)		
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for the reexec.Command - %s\n", err)
		time.Sleep(duration)
		os.Exit(1)
	}
}

func init() {
	reexec.Register("nsInitialisation", nsInitialisation)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	// containerName := os.Args[1]
	newrootPath := os.Args[2]
	hostname := os.Args[3]

	if err := mountProc(newrootPath); err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		time.Sleep(duration)		
		os.Exit(1)
	}

	if err := pivotRoot(newrootPath); err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		time.Sleep(duration)
		os.Exit(1)
	}

	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
		time.Sleep(duration)
		os.Exit(1)
	}

	if err := waitForNetwork(); err != nil {
		fmt.Printf("Error waiting for network - %s\n", err)
		time.Sleep(duration)
		os.Exit(1)
	}

	nsRun()
}

func nsRun() {
	command := os.Args[4:]
	cmd := exec.Command(command[0],command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//fmt.Printf("%v",os.Environ())
	//cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		time.Sleep(duration)
		os.Exit(1)
	}
}
