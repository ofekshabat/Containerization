package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

func waitForNetwork() error {
	maxWait := time.Second * 3
	checkInterval := time.Second
	timeStarted := time.Now()

	for {
		interfaces, err := net.Interfaces()
		if err != nil {
			return err
		}

		// pretty basic check ...
		// > 1 as a lo device will already exist
		if len(interfaces) > 1 {
			return nil
		}

		if time.Since(timeStarted) > maxWait {
			return fmt.Errorf("Timeout after %s waiting for network", maxWait)
		}

		time.Sleep(checkInterval)
	}
}

func exposePort(externalPort int, internalAddress string, internalPort int) {
	address1 := fmt.Sprintf("tcp-listen:%d,reuseaddr,fork", externalPort)
	address2 := fmt.Sprintf("tcp-connect:%s:%d", internalAddress, internalPort)
	exec.Command("socat", address1, address2).Start()
}
