package pinger

import (
	"context"
	"github.com/Rastler3D/container-monitoring/common/broker"
	"github.com/Rastler3D/container-monitoring/common/model"
	"github.com/Rastler3D/container-monitoring/pinger/internal/config"
	"github.com/Rastler3D/container-monitoring/pinger/internal/docker"
	"log"
	"time"
)

type Pinger struct {
	broker broker.MessageBroker[[]model.ContainerStatus]
	config config.Config
	docker docker.Client
}

func NewPinger(config config.Config) (Pinger, error) {
	brokerClient, err := broker.NewMessageBroker[[]model.ContainerStatus](config.Broker.URL, config.Broker.Queue)
	if err != nil {
		return Pinger{}, err
	}
	dockerClient, err := docker.NewClient()
	if err != nil {
		return Pinger{}, err
	}

	return Pinger{
		broker: brokerClient,
		config: config,
		docker: dockerClient,
	}, nil
}

func (pinger *Pinger) Start() error {
	return pinger.StartWithContext(context.Background())
}

func (pinger *Pinger) StartWithContext(ctx context.Context) error {
	defer pinger.Close()
	log.Printf("Starting pinger")

	interval := time.Tick(pinger.config.Pinger.PingInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-interval:
			containers, err := pinger.docker.GetContainers()

			if err != nil {
				return err
			}
			log.Printf("Discovered containers: %v", containers)
			pings := pinger.pingContainers(containers)

			log.Printf("All containers pinged. Sending pings to message broker")
			err = pinger.broker.Publish(pings)
			if err != nil {
				return err
			}
		}

	}
}

func (pinger *Pinger) Close() {
	pinger.broker.Close()
	pinger.docker.Close()
}

func (pinger *Pinger) pingContainers(containers []docker.Container) []model.ContainerStatus {
	pingResults := make([]model.ContainerStatus, 0, len(containers))

	for _, container := range containers {
	ips:
		for _, ip := range container.Ips {
			log.Printf("Start pinging containers : %s", ip)

			ping, err := pinger.docker.PingContainer(container.Pid, ip)
			if err != nil {
				log.Printf("Failed to ping container %s: %s", ip, err)
				continue
			}
			log.Printf("Container %s pinged successfully. Ping time:  %f", ping.IP, ping.PingTime)

			pingResults = append(pingResults, ping)

			break ips

		}
	}

	return pingResults
}
