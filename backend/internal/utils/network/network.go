package network

import (
	"fmt"
	"net"
	"strings"
)

// GetLocalIP retorna o IP local do host (ignora loopback e interfaces Docker/VM)
func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if ip := getIPv4FromInterface(iface); ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("não foi possível encontrar um IP local válido")
}

func getIPv4FromInterface(iface net.Interface) string {
	skipPrefixes := []string{"lo", "docker", "br-", "veth", "virbr"}
	if !isUsableInterface(iface, skipPrefixes) {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		ip := extractIP(addr)
		if ip != nil && ip.To4() != nil {
			return ip.String()
		}
	}

	return ""
}

func isUsableInterface(iface net.Interface, skipPrefixes []string) bool {
	if iface.Flags&net.FlagUp == 0 {
		return false
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(iface.Name, prefix) {
			return false
		}
	}

	return true
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
