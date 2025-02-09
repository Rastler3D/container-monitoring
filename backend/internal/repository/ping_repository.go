package repository

import (
	"database/sql"
	"fmt"
	"github.com/Rastler3D/container-monitoring/common/model"
	"strings"
)

type PingRepository struct {
	db *sql.DB
}

func NewPingRepository(db *sql.DB) PingRepository {
	return PingRepository{db: db}
}

func (r *PingRepository) AddContainerStatuses(statuses []model.ContainerStatus) error {

	valueStrings := make([]string, 0, len(statuses))
	valueArgs := make([]interface{}, 0, len(statuses)*3)
	for i, status := range statuses {
		// PostgreSQL uses placeholders like $1, $2, etc.
		// For MySQL, use '?' placeholders instead.
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, status.IP, status.PingTime, status.LastPing)
	}
	query := fmt.Sprintf("INSERT INTO containermonitor (ip_address, ping_time, last_ping) VALUES %s ON CONFLICT (ip_address) DO UPDATE SET ping_time = EXCLUDED.ping_time, last_ping = EXCLUDED.last_ping",
		strings.Join(valueStrings, ","))

	_, err := r.db.Exec(query, valueArgs...)
	return err
}

func (r *PingRepository) GetAllContainerStatuses() ([]model.ContainerStatus, error) {
	rows, err := r.db.Query("SELECT ip_address, ping_time, last_ping FROM containermonitor ORDER BY last_ping DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make([]model.ContainerStatus, 0)
	for rows.Next() {
		var s model.ContainerStatus
		if err := rows.Scan(&s.IP, &s.PingTime, &s.LastPing); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}
