package docker_helpers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

type machineCommand struct {
}

func (m *machineCommand) Create(driver, name string, opts ...string) error {
	args := []string{
		"create",
		"--driver", driver,
	}
	for _, opt := range opts {
		args = append(args, "--"+opt)
	}
	args = append(args, name)

	cmd := exec.Command("docker-machine", args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Debugln("Executing", cmd.Path, cmd.Args)
	return cmd.Run()
}

func (m *machineCommand) Provision(name string) error {
	cmd := exec.Command("docker-machine", "provision", name)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *machineCommand) Remove(name string) error {
	cmd := exec.Command("docker-machine", "rm", "-y", name)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *machineCommand) List(nodeFilter string) (machines []string, err error) {
	cmd := exec.Command("docker-machine", "ls", "-q")
	cmd.Env = os.Environ()
	data, err := cmd.Output()
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewReader(data))
	for {
		var line string

		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var query string
		if n, _ := fmt.Sscanf(line, nodeFilter, &query); n != 1 {
			continue
		}

		machines = append(machines, line)
	}
}

func (m *machineCommand) get(args ...string) (out string, err error) {
	// Execute docker-machine to fetch IP
	cmd := exec.Command("docker-machine", args...)
	cmd.Env = os.Environ()
	data, err := cmd.Output()
	if err != nil {
		return
	}

	// Save the IP
	out = strings.TrimSpace(string(data))
	if out == "" {
		err = fmt.Errorf("failed to get %v", args)
	}
	return
}

func (m *machineCommand) IP(name string) (string, error) {
	return m.get("ip", name)
}

func (m *machineCommand) URL(name string) (string, error) {
	return m.get("url", name)
}

func (m *machineCommand) CertPath(name string) (string, error) {
	return m.get("inspect", name, "-f", "{{.HostOptions.AuthOptions.StorePath}}")
}

func (m *machineCommand) Status(name string) (string, error) {
	return m.get("status", name)
}

func (m *machineCommand) CanConnect(name string) bool {
	// Execute docker-machine config which actively ask the machine if it is up and online
	cmd := exec.Command("docker-machine", "config", name)
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err == nil {
		return true
	}
	return false
}

func (m *machineCommand) Credentials(name string) (dc DockerCredentials, err error) {
	if !m.CanConnect(name) {
		err = errors.New("Can't connect")
		return
	}

	dc.TLSVerify = true
	dc.Host, err = m.URL(name)
	if err == nil {
		dc.CertPath, err = m.CertPath(name)
	}
	return
}

func NewMachineCommand() Machine {
	return &machineCommand{}
}
