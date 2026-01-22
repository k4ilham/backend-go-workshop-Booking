package app

import (
	"log"

	adapterfiber "be-golang/internal/adapter/http/fiber"
	"be-golang/internal/adapter/logger/turso"
	"be-golang/internal/adapter/notification/n8n"
	"be-golang/internal/adapter/repository/postgres"
	"be-golang/internal/config"
	"be-golang/internal/usecase"
	"be-golang/internal/util"

	fb "github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config) error {
	conn, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		return err
	}
	logAdapter := turso.New(cfg.TursoURL, cfg.TursoToken)
	notifier := n8n.New(cfg.N8NWebhookURL)
	j := util.NewJWT(cfg.JWTSecret)

	auth := usecase.NewAuthLogin(conn.Users(), logAdapter, j, cfg.TokenTTL)
	reg := usecase.NewAdminRegister(conn.Users(), logAdapter)
	bc := usecase.NewBookingCreate(conn.Bookings(), notifier, logAdapter)
	bl := usecase.NewBookingList(conn.Bookings())
	ds := usecase.NewDashboardStats(conn.Bookings())
	sc := usecase.NewServiceCreate(conn.Services())
	sd := usecase.NewServiceDelete(conn.Services())
	sla := usecase.NewServiceListActive(conn.Services())

	app := fb.New()
	handlers := adapterfiber.NewHandlers(auth, reg, bc, bl, ds, sc, sd, sla, j)
	handlers.Register(app)
	log.Println("server listening on", cfg.ServerAddr)
	return app.Listen(cfg.ServerAddr)
}
