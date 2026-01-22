package ports

import (
	"time"

	"be-golang/internal/domain"
)

type UserRepository interface {
	GetByEmail(email string) (*domain.User, error)
	Create(u domain.User) (int64, error)
}

type BookingRepository interface {
	Create(b domain.Booking) (int64, error)
	ListLatest(limit int) ([]domain.Booking, error)
	CountOnDate(day time.Time) (int, error)
}

type ServiceRepository interface {
	Create(s domain.Service) (int64, error)
	Delete(id int64) error
	ListActive() ([]domain.Service, error)
}

type Logger interface {
	Log(action string, detail string, at time.Time) error
}

type Notifier interface {
	NotifyBookingCreated(b domain.Booking) error
}
