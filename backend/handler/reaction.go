package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h24s_24/domain"
)

type ReactionRepository interface {
	GetReactionsByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Reaction, error)
}
