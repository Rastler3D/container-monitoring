package docker

import (
	"context"
	"fmt"
	"github.com/Rastler3D/container-monitoring/common/model"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-ping/ping"
	"github.com/vishvananda/netns"
	"log"
	"runtime"
	"time"
)

type Container struct {
	Pid int
	Ips []string
}

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

func (c *Client) GetContainers() ([]Container, error) {
	containers, err := c.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containersInfo = make([]Container, 0, len(containers))
	for _, cnt := range containers {
		dockerContainer, _ := c.cli.ContainerInspect(context.Background(), cnt.ID)
		ips := make([]string, 0, len(dockerContainer.NetworkSettings.Networks))

		for _, network := range dockerContainer.NetworkSettings.Networks {
			ips = append(ips, network.IPAddress)
		}
		containersInfo = append(containersInfo, Container{
			Pid: dockerContainer.State.Pid,
			Ips: ips,
		})
	}
	return containersInfo, nil
}

func (c *Client) PingContainer(pid int, ip string) (model.ContainerStatus, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origNS, err := netns.Get()
	if err != nil {
		log.Printf("Failed to get original netns when pinging container: %d", pid)
		return model.ContainerStatus{}, err
	}
	defer origNS.Close()

	containerNS, err := netns.GetFromPid(pid)
	if err != nil {
		log.Printf("Failed to get netns of container %d", pid)
		return model.ContainerStatus{}, err
	}
	defer containerNS.Close()

	if err := netns.Set(containerNS); err != nil {
		log.Printf("Failed to switch to netns of container %d", pid)
		return model.ContainerStatus{}, err
	}
	defer netns.Set(origNS)

	log.Printf("Switched to netns of container: %d", pid)

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
