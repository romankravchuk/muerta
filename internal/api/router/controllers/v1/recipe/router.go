package recipe

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	svc "github.com/romankravchuk/muerta/internal/services/recipe"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repo "github.com/romankravchuk/muerta/internal/storage/postgres/recipe"
)

func NewRouter(
	client postgres.Client,
	log logger.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route(context.RecipeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.RecipeID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
		router.Patch("/", jware.DeserializeUser, access.AdminOnly(log), handler.Restore)
		router.Route("/ingredients", func(router fiber.Router) {
			router.Get("/", handler.FindRecipeIngredients)
			router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.AddIngredient)
			router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.UpdateIngredient)
			router.Delete(
				"/",
				jware.DeserializeUser,
				access.AdminOnly(log),
				handler.RemoveIngredient,
			)
		})
		router.Route("/steps", func(router fiber.Router) {
			router.Get("/", handler.FindSteps)
			router.Route(context.StepID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.StepID))
				router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.AddStep)
				router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.RemoveStep)
			})
		})
	})
	return router
}
