package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
)

func checkFirewall() (bool, error) {
	checkFirewallCmd := exec.Command("systemctl", "show", "--property", "MainPID", "--value", "firewalld")
	firewallCmdOutput, err := checkFirewallCmd.Output()
	if err != nil {
		return false, err
	}

	firewallPid, err := strconv.Atoi(strings.TrimSuffix(string(firewallCmdOutput), "\n"))
	if err != nil {
		return false, err
	}

	return firewallPid > 0, nil
}

func addPort(port string, root bool) error {
	if root {
		return exec.Command("firewall-cmd", "--add-port="+port+"/tcp").Run()
	}

	return exec.Command("sudo", "firewall-cmd", "--add-port="+port+"/tcp").Run()
}

func removePort(port string, root bool) error {
	if root {
		return exec.Command("firewall-cmd", "--remove-port="+port+"/tcp").Run()
	}

	return exec.Command("sudo", "firewall-cmd", "--remove-port="+port+"/tcp").Run()
}

func cleanup(c chan os.Signal, port string) {
	<-c
	removePort(port, os.Getuid() == 0)
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: serve <file> [port]")
		os.Exit(1)
	}

	file := os.Args[1]
	port := "8080"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	useFirewall, err := checkFirewall()
	if err != nil {
		fmt.Println("Error checking firewall status:", err)
		os.Exit(1)
	}

	if useFirewall {
		addPort(port, os.Getuid() == 0)
	} else {
		fmt.Println("firewalld is not running; continuing without changing firewall rules")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	if useFirewall {
		go cleanup(c, port)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	})

	fmt.Printf("Serving %s on port %s\n", file, port)
	err = http.ListenAndServe(":"+port, nil)
}
