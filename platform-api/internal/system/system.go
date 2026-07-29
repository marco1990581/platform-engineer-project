package system

type SystemInfo struct {
	Hostname Hostname `json:"hostname"`
	Memory   Memory   `json:"memory"`
	Uptime   Uptime   `json:"uptime"`
	Network  Network  `json:"network"`
}

func GetSystemInfo() (SystemInfo, error) {
	hostname, err := GetHostname()
	if err != nil {
		return SystemInfo{}, err
	}

	memory, err := GetMemory()
	if err != nil {
		return SystemInfo{}, err
	}

	uptime, err := GetUptime()
	if err != nil {
		return SystemInfo{}, err
	}

	network, err := GetNetwork()
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		Hostname: hostname,
		Memory:   memory,
		Uptime:   uptime,
		Network:  network,
	}, nil
}
