package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Rastler3D/container-monitoring/backend/internal/config"
	"github.com/Rastler3D/container-monitoring/backend/internal/database"
	"github.com/Rastler3D/container-monitoring/backend/internal/server"
	"github.com/Rastler3D/container-monitoring/common/broker"
	"github.com/Rastler3D/container-monitoring/common/model"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	if err := database.InitDatabase(db); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}

	broker_, err := broker.NewMessageBroker[[]model.ContainerStatus](cfg.Broker.URL, cfg.Broker.Queue)
	if err != nil {
		fmt.Printf("Error creating message broker: %v\n", err)
		return
	}
	defer broker_.Close()

	server_ := server.NewServer(cfg.Server.Port, db, &broker_)

	err = server_.StartWithContext(ctx)
	if err != nil {
		fmt.Printf("Error starting server: %v\n\n", err)
		return
	}

}
