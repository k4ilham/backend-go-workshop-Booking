package n8n

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"be-golang/internal/domain"
)

type Notifier struct {
	url   string
	httpc *http.Client
}

func New(url string) *Notifier {
	return &Notifier{url: url, httpc: &http.Client{Timeout: 5 * time.Second}}
}

func (n *Notifier) NotifyBookingCreated(b domain.Booking) error {
	if n.url == "" {
		return nil
	}
	payload := map[string]any{
		"event":          "booking_created",
		"id":             b.ID,
		"customer_name":  b.CustomerName,
		"customer_phone": b.CustomerPhone,
		"service_id":     b.ServiceID,
		"booking_date":   b.BookingDate.Format("2006-01-02"),
		"booking_time":   b.BookingTime,
		"status":         b.Status,
		"created_at":     b.CreatedAt.Format(time.RFC3339),
	}
	data, _ := json.Marshal(payload)
	_, err := n.httpc.Post(n.url, "application/json", bytes.NewReader(data))
	return err
}
