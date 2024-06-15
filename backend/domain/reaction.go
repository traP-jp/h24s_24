package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	ReactionID int
	PostID     uuid.UUID
	Count      int
	Users      []string
}

type UserReaction struct {
	ReactionID int
	PostID     uuid.UUID
	CreatedAt  time.Time
}
