package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/razvanmarinn/chatroom/internal/services"
)

func AddServicesToContext(sm *services.ServiceManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "serviceManager", sm)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
