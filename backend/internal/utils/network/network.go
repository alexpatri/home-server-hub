package network

import (
	"fmt"
	"net"
)

// GetLocalIP retorna o IP local do host (ignora interfaces loopback)
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ip := extractIP(addr)
		if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("não foi possível encontrar o IP local")
}

func extractIP(addr net.Addr) net.IP {
	if ipNet, ok := addr.(*net.IPNet); ok {
		return ipNet.IP
	}

	if ipAddr, ok := addr.(*net.IPAddr); ok {
		return ipAddr.IP
	}

	return nil
}
