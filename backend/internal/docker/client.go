package docker

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Client encapsula a API do Docker
type Client struct {
	client *client.Client
}

// ContainerInfo representa informações de um container
type ContainerInfo struct {
	ID    string
	Name  string
	Image string
	Ports []Port
	State string
	IP    string
}

// Port representa uma porta exposta por um container
type Port struct {
	HostIP        string
	HostPort      string
	ContainerPort string
	Protocol      string
}

// NewClient cria um novo cliente Docker
func NewClient(host string) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{client: cli}, nil
}

// ListContainers lista todos os containers disponíveis
func (c *Client) ListContainers() ([]ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Lista todos os containers (running e stopped)
	containerList, err := c.client.ContainerList(ctx, container.ListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("status", "running"),
		),
	})
	if err != nil {
		return nil, err
	}

	var containers []ContainerInfo
	for _, container := range containerList {
		name := ""
		if len(container.Names) > 0 {
			// Remove o prefixo "/" do nome
			name = strings.TrimPrefix(container.Names[0], "/")
		}

		var ports []Port
		for _, port := range container.Ports {
			ports = append(ports, Port{
				HostIP:        port.IP,
				HostPort:      string(port.PublicPort),
				ContainerPort: string(port.PrivatePort),
				Protocol:      port.Type,
			})
		}

		// Obter detalhes do container para IP
		containerInfo, err := c.client.ContainerInspect(ctx, container.ID)
		if err != nil {
			continue
		}

		ip := ""
		// Tenta obter o IP do container
		if containerInfo.NetworkSettings != nil && len(containerInfo.NetworkSettings.Networks) > 0 {
			for _, network := range containerInfo.NetworkSettings.Networks {
				ip = network.IPAddress
				break
			}
		}

		containers = append(containers, ContainerInfo{
			ID:    container.ID,
			Name:  name,
			Image: container.Image,
			Ports: ports,
			State: container.State,
			IP:    ip,
		})
	}

	return containers, nil
}

// GetContainerStatus obtém o status atual de um container
func (c *Client) GetContainerStatus(containerID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containerInfo, err := c.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}

	if containerInfo.State.Running {
		return "running", nil
	}
	return "stopped", nil
}

// Close fecha o cliente Docker
func (c *Client) Close() error {
	return c.client.Close()
}
