package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID               uuid.UUID
	UserName         string
	OriginalMessage  string
	ConvertedMessage string
	ParentID         uuid.UUID
	RootID           uuid.UUID
	CreatedAt        time.Time
}
