package common

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
)

var DefaultIdKey = "id"

func GetIdByFiberCtx(ctx *fiber.Ctx, idKey ...string) (int, error) {
	key := DefaultIdKey
	if len(idKey) != 0 {
		key = idKey[0]
	}
	param := ctx.Params(key, "")
	id, err := strconv.Atoi(param)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetNameByFiberCtx(ctx *fiber.Ctx) (string, error) {
	param := ctx.Params("name")
	if param == "" {
		return "", fmt.Errorf("name is required")
	}
	return param, nil
}

func GetFilterByFiberCtx[T dto.Filter](ctx *fiber.Ctx, filter T) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.GetLimit() == 0 {
		filter.SetLimit(10)
	}
	if filter.GetOffset() < 0 {
		filter.SetOffset(0)
	}
	return nil
}

