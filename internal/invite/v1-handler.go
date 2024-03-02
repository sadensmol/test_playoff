package invite

import (
	"context"
	"net/mail"
	"strings"
	"time"

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
	RegisterInvite(ctx context.Context, invite Invite) error
}

func NewHandler(inviteService IInviteService) *Handler {
	return &Handler{
		inviteService: inviteService,
	}
}

func (h *Handler) registerInvite(c echo.Context) error {
	req := new(InviteRequest)
	err := c.Bind(&req)
	if err != nil {
		log.Error().Err(err).Msg("failed to bind request")
		return nil
	}

	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		log.Error().Err(err).Msgf("invalid email passed %v", req)
		return nil
	}

	if strings.TrimSpace(req.Code) == "" {
		log.Error().Msgf("code is blank %v", req)
		return nil
	}

	return h.inviteService.RegisterInvite(context.TODO(), Invite{Code: req.Code, Email: req.Email, RegisteredAt: time.Now()})
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/api/v1/invite", h.registerInvite)
}
