package system

import (
	"os"
	"strconv"
	"strings"
)

type Memory struct {
	Total     uint64 `json:"total_kb"`
	Free      uint64 `json:"free_kb"`
	Available uint64 `json:"available_kb"`
}

func GetMemory() (Memory, error) {

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return Memory{}, err
	}

	var memory Memory

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch fields[0] {

		case "MemTotal:":
			memory.Total = value

		case "MemFree:":
			memory.Free = value

		case "MemAvailable:":
			memory.Available = value
		}
	}

	return memory, nil
}
