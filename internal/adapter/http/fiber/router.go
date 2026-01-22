package fiber

import (
	"strconv"
	"time"

	"be-golang/internal/domain"
	"be-golang/internal/usecase"
	"be-golang/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	authLogin         *usecase.AuthLogin
	adminRegister     *usecase.AdminRegister
	bookingCreate     *usecase.BookingCreate
	bookingList       *usecase.BookingList
	dashboardStats    *usecase.DashboardStats
	serviceCreate     *usecase.ServiceCreate
	serviceDelete     *usecase.ServiceDelete
	serviceListActive *usecase.ServiceListActive
	jwt               *util.JWT
}

func NewHandlers(auth *usecase.AuthLogin, reg *usecase.AdminRegister, bc *usecase.BookingCreate, bl *usecase.BookingList, ds *usecase.DashboardStats, sc *usecase.ServiceCreate, sd *usecase.ServiceDelete, sla *usecase.ServiceListActive, jwt *util.JWT) *Handlers {
	return &Handlers{
		authLogin:         auth,
		adminRegister:     reg,
		bookingCreate:     bc,
		bookingList:       bl,
		dashboardStats:    ds,
		serviceCreate:     sc,
		serviceDelete:     sd,
		serviceListActive: sla,
		jwt:               jwt,
	}
}

func (h *Handlers) Register(app *fiber.App) {
	app.Post("/admin/login", h.login)
	app.Post("/admin/register", h.register)
	app.Post("/bookings", h.createBooking)
	app.Get("/bookings", h.listBookings)
	app.Get("/admin/dashboard", h.jwtMiddleware, h.dashboard)
	app.Post("/services", h.jwtMiddleware, h.createService)
	app.Delete("/services/:id", h.jwtMiddleware, h.deleteService)
	app.Get("/services", h.jwtMiddleware, h.listActiveServices)
}

func (h *Handlers) jwtMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if len(auth) < 8 || auth[:7] != "Bearer " {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	tokenStr := auth[7:]
	tok, err := h.jwt.Parse(tokenStr)
	if err != nil || !tok.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}

func (h *Handlers) login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	token, err := h.authLogin.Exec(body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_credentials"})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *Handlers) register(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	id, err := h.adminRegister.Exec(body.Email, body.Password)
	if err != nil {
		if err.Error() == "email_exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email_exists"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "register_failed"})
	}
	return c.JSON(fiber.Map{"id": id})
}

func (h *Handlers) createBooking(c *fiber.Ctx) error {
	var body struct {
		CustomerName  string `json:"customer_name"`
		CustomerPhone string `json:"customer_phone"`
		ServiceID     int64  `json:"service_id"`
		BookingDate   string `json:"booking_date"`
		BookingTime   string `json:"booking_time"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	date, err := time.Parse("2006-01-02", body.BookingDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_date"})
	}
	id, err := h.bookingCreate.Exec(domain.Booking{
		CustomerName:  body.CustomerName,
		CustomerPhone: body.CustomerPhone,
		ServiceID:     body.ServiceID,
		BookingDate:   date,
		BookingTime:   body.BookingTime,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "create_failed"})
	}
	return c.JSON(fiber.Map{"id": id})
}

func (h *Handlers) listBookings(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	items, err := h.bookingList.Exec(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "list_failed"})
	}
	return c.JSON(items)
}

func (h *Handlers) dashboard(c *fiber.Ctx) error {
	res, err := h.dashboardStats.Exec(time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "stats_failed"})
	}
	return c.JSON(res)
}

func (h *Handlers) createService(c *fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Price    int64  `json:"price"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	id, err := h.serviceCreate.Exec(domain.Service{
		Name:     body.Name,
		Price:    body.Price,
		IsActive: body.IsActive,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "create_failed"})
	}
	return c.JSON(fiber.Map{"id": id})
}

func (h *Handlers) deleteService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_id"})
	}
	err = h.serviceDelete.Exec(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "delete_failed"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) listActiveServices(c *fiber.Ctx) error {
	items, err := h.serviceListActive.Exec()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "list_failed"})
	}
	return c.JSON(items)
}
