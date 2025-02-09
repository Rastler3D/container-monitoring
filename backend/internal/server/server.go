package server

import (
	"context"
	"database/sql"
	"github.com/Rastler3D/container-monitoring/backend/internal/api"
	"github.com/Rastler3D/container-monitoring/backend/internal/repository"
	"github.com/Rastler3D/container-monitoring/backend/internal/service"
	"github.com/Rastler3D/container-monitoring/common/broker"
	"github.com/Rastler3D/container-monitoring/common/model"
	"log"
	"net/http"
)

type Server struct {
	Port     string
	Consumer api.Consumer
	Handler  api.Handler
}

func NewServer(port string, db *sql.DB, broker *broker.MessageBroker[[]model.ContainerStatus]) *Server {

	pingRepository := repository.NewPingRepository(db)
	pingService := service.NewPingService(&pingRepository)
	consumer := api.NewConsumer(&pingService, broker, log.Default())
	handler := api.NewHandler(&pingService, log.Default())

	return &Server{
		Port:     port,
		Consumer: consumer,
		Handler:  handler,
	}
}

func (s *Server) StartWithContext(ctx context.Context) error {
	log.Printf("Starting server on port %s\n", s.Port)
	err := s.Consumer.StartWithContext(ctx)
	if err != nil {
		return err
	}

	server := http.Server{Addr: ":" + s.Port, Handler: s.Handler.Router()}

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}
func (s *Server) Start() error {
	return s.StartWithContext(context.Background())
}
