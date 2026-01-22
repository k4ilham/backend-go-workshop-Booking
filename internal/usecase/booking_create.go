package usecase

import (
	"time"

	"be-golang/internal/domain"
	"be-golang/internal/ports"
)

type BookingCreate struct {
	bookings ports.BookingRepository
	notifier ports.Notifier
	logger   ports.Logger
}

func NewBookingCreate(b ports.BookingRepository, n ports.Notifier, l ports.Logger) *BookingCreate {
	return &BookingCreate{bookings: b, notifier: n, logger: l}
}

func (u *BookingCreate) Exec(input domain.Booking) (int64, error) {
	now := time.Now().UTC()
	input.Status = "pending"
	input.CreatedAt = now
	id, err := u.bookings.Create(input)
	if err != nil {
		return 0, err
	}
	input.ID = id
	_ = u.notifier.NotifyBookingCreated(input)
	_ = u.logger.Log("booking_created", input.CustomerName, now)
	return id, nil
}
