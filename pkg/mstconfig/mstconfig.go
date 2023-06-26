package mstconfig

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type SrIovConfig struct {
	Enabled  bool
	VfsCount int
}

type Config struct {
	SrIov SrIovConfig
}

func Query(pciAddress string) (*Config, error) {
	cmd := exec.Command("mstconfig", "-d", pciAddress, "query")
	stdout, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("mstconfig: %w", err)
	}

	return parseQuery(bytes.NewReader(stdout))
}

func parseQuery(output io.Reader) (*Config, error) {
	config := &Config{}

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)

	parsingConfig := false
	for scanner.Scan() {
		// If we see the "Configurations:" line, then we know we're about to start
		// parsing the config.
		if strings.Contains(scanner.Text(), "Configurations:") {
			parsingConfig = true
			continue
		}

		// So long as we're not parsing the config, we can skip this line.
		if !parsingConfig {
			continue
		}

		key, val := parseConfig(scanner.Text())

		switch key {
		case "SRIOV_EN":
			enabled, err := parseConfigValueBool(val)
			if err != nil {
				return nil, fmt.Errorf("mstconfig: failed to parse SRIOV_EN %w", err)
			}
			config.SrIov.Enabled = enabled
		case "NUM_OF_VFS":
			vfsCount, err := parseConfigValueInt(val)
			if err != nil {
				return nil, fmt.Errorf("mstconfig: failed to parse NUM_OF_VFS %w", err)
			}
			config.SrIov.VfsCount = vfsCount
		}
	}

	return config, nil
}

func parseConfig(line string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
	return parts[0], strings.TrimSpace(parts[1])
}

func parseConfigValueBool(value string) (bool, error) {
	switch value {
	case "True(1)":
		return true, nil
	case "False(0)":
		return false, nil
	}

	return false, fmt.Errorf("mstconfig: failed to parse bool %q", value)
}

func parseConfigValueInt(value string) (int, error) {
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("mstconfig: failed to parse int %w", err)
	}

	return num, nil
}
