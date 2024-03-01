package invite

import "time"

type Invite struct {
	Email        string `json:"email"`
	RegisteredAt time.Time
}
