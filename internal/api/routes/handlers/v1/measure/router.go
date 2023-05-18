package measure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/measure"
	service "github.com/romankravchuk/muerta/internal/services/measure"
)

func NewRouter(
	client repositories.PostgresClient,
	log *log.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, handler.Create)
	router.Route(context.MeasureID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.MeasureID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
	})
	return router
}
