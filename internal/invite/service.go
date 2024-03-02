package invite

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type InviteService struct {
	db *mongo.Database
}

type InviteModel struct {
	Email        string    `json:"email" bson:"email"`
	RegisteredAt time.Time `json:"reg_at" bson:"reg_at"`
}

func NewInviteService(cl *mongo.Client) *InviteService {
	db := cl.Database("invites")

	return &InviteService{
		db: db,
	}
}

func (s *InviteService) RegisterInvite(ctx context.Context, invite Invite) error {
	res, err := s.db.Collection(fmt.Sprintf("code-%s", invite.Code)).
		InsertOne(ctx, InviteModel{Email: invite.Email, RegisteredAt: invite.RegisteredAt})

	if err != nil {
		return err
	}
	_ = res
	return nil
}
