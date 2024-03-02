package invite

import "time"

type Invite struct {
	Code         string
	Email        string
	RegisteredAt time.Time
}
