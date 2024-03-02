package invite

import (
	"context"
	"net/mail"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// todo валидация емейл
// todo проверить код валидный
type InviteRequest struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

type Handler struct {
	inviteService IInviteService
}

type IInviteService interface {
	RegisterInvite(ctx context.Context, email string) error
}

func NewHandler(inviteService IInviteService) *Handler {
	return &Handler{
		inviteService: inviteService,
	}
}

func (h *Handler) registerInvite(c echo.Context) error {
	req := new(InviteRequest)
	c.Bind(req)

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		log.Error().Err(err).Msg("invalid email passed")
		return nil
	}

	//todo add invite code

	return h.inviteService.RegisterInvite(context.TODO(), Invite{Code:})
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/invite", h.registerInvite)
}
