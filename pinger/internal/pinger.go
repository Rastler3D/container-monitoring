package pinger

import (
	"context"
	"github.com/Rastler3D/container-monitoring/common/broker"
	"github.com/Rastler3D/container-monitoring/common/model"
	"github.com/Rastler3D/container-monitoring/pinger/internal/config"
	"github.com/Rastler3D/container-monitoring/pinger/internal/docker"
	"log"
	"sync"
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

func (pinger *Pinger) pingContainers(containers []string) []model.ContainerStatus {
	pingChannel := make(chan model.ContainerStatus, len(containers))
	var wg sync.WaitGroup

	for _, ip := range containers {
		wg.Add(1)
		go func() {
			log.Printf("Start pinging containers : %s", ip)
			defer wg.Done()
			ping, err := pinger.docker.PingContainer(ip)
			if err != nil {
				log.Printf("Failed to ping container %s: %s", ip, err)
				return
			}
			log.Printf("Container %s pinged successfully. Ping time:  %f", ping.IP, ping.PingTime)
			pingChannel <- ping
		}()
	}
	wg.Wait()
	close(pingChannel)

	pingResults := make([]model.ContainerStatus, 0, len(containers))
	for result := range pingChannel {
		pingResults = append(pingResults, result)
	}

	return pingResults
}
