package usersetting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	usersetting "github.com/romankravchuk/muerta/internal/services/user-setting"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/setting"
)

func NewRouter(
	client postgres.Client,
	log *log.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := setting.New(client)
	svc := usersetting.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route(context.SettingID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.SettingID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Patch("/", jware.DeserializeUser, access.AdminOnly(log), handler.Restore)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
	})
	return router
}
