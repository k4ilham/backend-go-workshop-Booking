package domain

import "time"

type Booking struct {
	ID            int64
	CustomerName  string
	CustomerPhone string
	ServiceID     int64
	BookingDate   time.Time
	BookingTime   string
	Status        string
	CreatedAt     time.Time
}
