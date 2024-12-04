package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/logger"
	"github.com/razvanmarinn/chatroom/internal/services"
)

func AddToContext(sm *services.ServiceManager, logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "serviceManager", sm)
			ctx = context.WithValue(ctx, "logger", logger)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
