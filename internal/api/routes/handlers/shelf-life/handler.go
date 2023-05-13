package shelflife

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
)

type ShelfLifeHandler struct {
	svc service.ShelfLifeServicer
	log *log.Logger
}

func New(svc service.ShelfLifeServicer, log *log.Logger) ShelfLifeHandler {
	return ShelfLifeHandler{
		svc: svc,
		log: log,
	}
}

func (h *ShelfLifeHandler) CreateShelfLife(ctx *fiber.Ctx) error {
	var payload *dto.CreateShelfLifeDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateShelfLife(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ShelfLifeHandler) FindShelfLifeByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.FindShelfLifeByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"shelf_life": result},
	))
}

func (h *ShelfLifeHandler) FindShelfLives(ctx *fiber.Ctx) error {
	filter := new(dto.ShelfLifeFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindShelfLifes(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"shelf_lives": result, "count": count},
	))
}

func (h *ShelfLifeHandler) UpdateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	payload := new(dto.UpdateShelfLifeDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateShelfLife(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ShelfLifeHandler) DeleteShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.DeleteShelfLife(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ShelfLifeHandler) RestoreShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.RestoreShelfLife(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ShelfLifeHandler) FindShelfLifeStatuses(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.FindShelfLifeStatuses(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"statuses": result},
	))
}
func (h *ShelfLifeHandler) CreateShelfLifeStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	statusID := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.CreateShelfLifeStatus(ctx.Context(), id, statusID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"status": result},
	))
}
func (h *ShelfLifeHandler) DeleteShelfLifeStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	statusID := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id, statusID); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}
