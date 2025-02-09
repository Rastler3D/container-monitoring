package api

import (
	"context"
	"github.com/Rastler3D/container-monitoring/backend/internal/service"
	"github.com/Rastler3D/container-monitoring/common/broker"
	"github.com/Rastler3D/container-monitoring/common/model"
	"log"
)

type Consumer struct {
	pingService *service.PingService
	broker      *broker.MessageBroker[[]model.ContainerStatus]
	logger      *log.Logger
}

func NewConsumer(pingService *service.PingService, broker *broker.MessageBroker[[]model.ContainerStatus], logger *log.Logger) Consumer {
	return Consumer{pingService: pingService, logger: logger, broker: broker}
}

func (consumer *Consumer) Start() error {
	return consumer.StartWithContext(context.Background())
}

func (consumer *Consumer) StartWithContext(ctx context.Context) error {

	channel, err := consumer.broker.Consume()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				consumer.logger.Printf("%v", ctx.Err())
				return
			case statuses := <-channel:
				if err := consumer.pingService.AddContainerStatuses(statuses); err != nil {
					consumer.logger.Printf("Error saving container status: %v", err)
				}
			}
		}
	}()

	return nil
}
