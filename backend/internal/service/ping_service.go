package service

import (
	"github.com/Rastler3D/container-monitoring/backend/internal/repository"
	"github.com/Rastler3D/container-monitoring/common/model"
)

type PingService struct {
	repo *repository.PingRepository
}

func NewPingService(repo *repository.PingRepository) PingService {
	return PingService{repo: repo}
}

func (s *PingService) AddContainerStatuses(status []model.ContainerStatus) error {
	return s.repo.AddContainerStatuses(status)
}

func (s *PingService) GetAllContainerStatuses() ([]model.ContainerStatus, error) {
	return s.repo.GetAllContainerStatuses()
}
