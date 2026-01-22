package usecase

import (
	"time"

	"be-golang/internal/domain"
	"be-golang/internal/ports"
)

type DashboardStats struct {
	bookings ports.BookingRepository
}

type DashboardResult struct {
	TotalToday int
	Latest     []domain.Booking
}

func NewDashboardStats(b ports.BookingRepository) *DashboardStats {
	return &DashboardStats{bookings: b}
}

func (u *DashboardStats) Exec(today time.Time) (DashboardResult, error) {
	count, err := u.bookings.CountOnDate(today)
	if err != nil {
		return DashboardResult{}, err
	}
	latest, err := u.bookings.ListLatest(10)
	if err != nil {
		return DashboardResult{}, err
	}
	return DashboardResult{TotalToday: count, Latest: latest}, nil
}
