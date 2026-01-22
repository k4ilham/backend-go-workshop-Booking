package usecase

import (
	"be-golang/internal/domain"
	"be-golang/internal/ports"
)

type ServiceCreate struct {
	services ports.ServiceRepository
}

func NewServiceCreate(s ports.ServiceRepository) *ServiceCreate {
	return &ServiceCreate{services: s}
}

func (u *ServiceCreate) Exec(s domain.Service) (int64, error) {
	return u.services.Create(s)
}

type ServiceDelete struct {
	services ports.ServiceRepository
}

func NewServiceDelete(s ports.ServiceRepository) *ServiceDelete {
	return &ServiceDelete{services: s}
}

func (u *ServiceDelete) Exec(id int64) error {
	return u.services.Delete(id)
}

type ServiceListActive struct {
	services ports.ServiceRepository
}

func NewServiceListActive(s ports.ServiceRepository) *ServiceListActive {
	return &ServiceListActive{services: s}
}

func (u *ServiceListActive) Exec() ([]domain.Service, error) {
	return u.services.ListActive()
}
