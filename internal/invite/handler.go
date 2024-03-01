package invite

import "github.com/labstack/echo/v4"

// todo валидация емейл
// todo проверить код валидный

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) registerInvite(c echo.Context) error {
	return nil
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/invite", h.registerInvite)
}
