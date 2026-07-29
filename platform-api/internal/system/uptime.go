package system

import (
	"os"
	"strconv"
	"strings"
)

type Uptime struct {
	Seconds float64 `json:"seconds"`
}

func GetUptime() (Uptime, error) {

	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return Uptime{}, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return Uptime{}, nil
	}

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return Uptime{}, err
	}

	uptime := Uptime{
		Seconds: seconds,
	}

	return uptime, nil
}
