package v1

import (
	"fmt"
	"net"
)

const (
	NetMajorInterface = "eth0"
)

func GetMacAddr(interfaceName string) string {
	nets, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	macAddr := ""
	for _, net := range nets {
		if net.Name == interfaceName {
			macAddr = net.HardwareAddr.String()
			break
		}
	}
	if macAddr == "" {
		err := fmt.Errorf("mac address not found from interface: %s", interfaceName)
		panic(err)
	}

	return macAddr
}
