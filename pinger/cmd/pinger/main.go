package main

import (
	"context"
	"fmt"
	"github.com/Rastler3D/container-monitoring/pinger/internal"
	"github.com/Rastler3D/container-monitoring/pinger/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	pinger_, err := pinger.NewPinger(cfg)
	if err != nil {
		fmt.Println("Error creating pinger:", err)
		return
	}
	defer pinger_.Close()

	err = pinger_.StartWithContext(ctx)
	if err != nil {
		fmt.Println("Error creating pinger:", err)
		return
	}

}
