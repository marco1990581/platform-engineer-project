package system

import "os"

type Hostname struct {
	Name string `json:"hostname"`
}

func GetHostname() (Hostname, error) {

	name, err := os.Hostname()
	if err != nil {
		return Hostname{}, err
	}

	hostname := Hostname{
		Name: name,
	}

	return hostname, nil
}
