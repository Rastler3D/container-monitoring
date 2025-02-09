package model

import "time"

type ContainerStatus struct {
	IP       string    `json:"ip"`
	PingTime float64   `json:"ping_time"`
	LastPing time.Time `json:"last_ping"`
}
