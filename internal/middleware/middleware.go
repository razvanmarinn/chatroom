package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	_db "github.com/razvanmarinn/chatroom/internal/db"
)

func AddRepositoriesToContext(userRepo *_db.UserRepository, roomRepo *_db.RoomRepository, messagesRepo *_db.MessageRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "userRepo", userRepo)
			ctx = context.WithValue(ctx, "roomRepo", roomRepo)
			ctx = context.WithValue(ctx, "messageRepo", roomRepo)
			

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
