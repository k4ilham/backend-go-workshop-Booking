package postgres

import (
	"database/sql"
	"time"

	"be-golang/internal/domain"

	_ "github.com/lib/pq"
)

type Connection struct {
	DB *sql.DB
}

func New(dsn string) (*Connection, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Connection{DB: db}, nil
}

type UserRepo struct{ db *sql.DB }
type BookingRepo struct{ db *sql.DB }
type ServiceRepo struct{ db *sql.DB }

func (c *Connection) Users() *UserRepo       { return &UserRepo{db: c.DB} }
func (c *Connection) Bookings() *BookingRepo { return &BookingRepo{db: c.DB} }
func (c *Connection) Services() *ServiceRepo { return &ServiceRepo{db: c.DB} }

func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, created_at FROM users WHERE email=$1 LIMIT 1`, email)
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(u domain.User) (int64, error) {
	err := r.db.QueryRow(`INSERT INTO users (email, password_hash, created_at) VALUES ($1,$2,$3) RETURNING id`, u.Email, u.PasswordHash, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}
func (r *BookingRepo) Create(b domain.Booking) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO bookings (customer_name, customer_phone, service_id, booking_date, booking_time, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		b.CustomerName, b.CustomerPhone, b.ServiceID, b.BookingDate, b.BookingTime, b.Status, b.CreatedAt,
	).Scan(&b.ID)
	if err != nil {
		return 0, err
	}
	return b.ID, nil
}

func (r *BookingRepo) ListLatest(limit int) ([]domain.Booking, error) {
	rows, err := r.db.Query(
		`SELECT id, customer_name, customer_phone, service_id, booking_date, booking_time, status, created_at
		 FROM bookings ORDER BY created_at DESC LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Booking
	for rows.Next() {
		var b domain.Booking
		err = rows.Scan(&b.ID, &b.CustomerName, &b.CustomerPhone, &b.ServiceID, &b.BookingDate, &b.BookingTime, &b.Status, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BookingRepo) CountOnDate(day time.Time) (int, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	row := r.db.QueryRow(`SELECT COUNT(*) FROM bookings WHERE created_at >= $1 AND created_at < $2`, start, end)
	var c int
	err := row.Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (r *ServiceRepo) Create(s domain.Service) (int64, error) {
	err := r.db.QueryRow(`INSERT INTO services (name, price, is_active) VALUES ($1,$2,$3) RETURNING id`, s.Name, s.Price, s.IsActive).Scan(&s.ID)
	if err != nil {
		return 0, err
	}
	return s.ID, nil
}

func (r *ServiceRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM services WHERE id=$1`, id)
	return err
}

func (r *ServiceRepo) ListActive() ([]domain.Service, error) {
	rows, err := r.db.Query(`SELECT id, name, price, is_active FROM services WHERE is_active=TRUE ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Service
	for rows.Next() {
		var s domain.Service
		err = rows.Scan(&s.ID, &s.Name, &s.Price, &s.IsActive)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

var _ interface {
	GetByEmail(string) (*domain.User, error)
	Create(domain.User) (int64, error)
} = (*UserRepo)(nil)
var _ interface {
	Create(domain.Booking) (int64, error)
	ListLatest(int) ([]domain.Booking, error)
	CountOnDate(time.Time) (int, error)
} = (*BookingRepo)(nil)
var _ interface {
	Create(domain.Service) (int64, error)
	Delete(int64) error
	ListActive() ([]domain.Service, error)
} = (*ServiceRepo)(nil)
