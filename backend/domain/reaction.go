package domain

import "github.com/google/uuid"

type Reaction struct {
	ReactionID int
	PostID     uuid.UUID
	Count      int
	Users      []string
}
