package system

import "net"

type Network struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

type NetworkInterface struct {
	Name      string   `json:"name"`
	MAC       string   `json:"mac"`
	MTU       int      `json:"mtu"`
	Up        bool     `json:"up"`
	Loopback  bool     `json:"loopback"`
	Multicast bool     `json:"multicast"`
	Addresses []string `json:"addresses"`
}

func GetNetwork() (Network, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return Network{}, err
	}

	network := Network{
		Interfaces: make([]NetworkInterface, 0, len(ifaces)),
	}

	for _, iface := range ifaces {

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		addresses := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}

		network.Interfaces = append(network.Interfaces, NetworkInterface{
			Name:      iface.Name,
			MAC:       iface.HardwareAddr.String(),
			MTU:       iface.MTU,
			Up:        iface.Flags&net.FlagUp != 0,
			Loopback:  iface.Flags&net.FlagLoopback != 0,
			Multicast: iface.Flags&net.FlagMulticast != 0,
			Addresses: addresses,
		})
	}

	return network, nil
}
