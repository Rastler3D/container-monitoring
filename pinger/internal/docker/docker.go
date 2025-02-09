package docker

import (
	"context"
	"fmt"
	"github.com/Rastler3D/container-monitoring/common/model"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-ping/ping"
	"time"
)

type Client struct {
	cli *client.Client
}

func NewClient() (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return Client{}, err
	}
	return Client{cli: cli}, nil
}

func (c *Client) GetContainers() ([]string, error) {
	containers, err := c.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, dockerContainer := range containers {
		for _, network := range dockerContainer.NetworkSettings.Networks {
			ips = append(ips, network.IPAddress)
		}
	}
	return ips, nil
}

func (c *Client) PingContainer(ip string) (model.ContainerStatus, error) {
	pinger, err := ping.NewPinger(ip)

	if err != nil {
		return model.ContainerStatus{}, err
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	err = pinger.Run()

	if err != nil {
		return model.ContainerStatus{}, err
	}
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return model.ContainerStatus{}, fmt.Errorf("received 0 packets")
	}

	result := model.ContainerStatus{
		IP:       ip,
		PingTime: stats.AvgRtt.Seconds() * 1000,
		LastPing: time.Now(),
	}

	return result, nil
}

func (c *Client) Close() {
	if err := c.cli.Close(); err != nil {
		panic(err)
	}
}
