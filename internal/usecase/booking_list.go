package usecase

import (
	"be-golang/internal/domain"
	"be-golang/internal/ports"
)

type BookingList struct {
	bookings ports.BookingRepository
}

func NewBookingList(b ports.BookingRepository) *BookingList {
	return &BookingList{bookings: b}
}

func (u *BookingList) Exec(limit int) ([]domain.Booking, error) {
	return u.bookings.ListLatest(limit)
}
