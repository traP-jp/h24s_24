package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/traP-jp/h24s_24/domain"
)

type ReactionRepository struct {
	DB *sqlx.DB
}

func NewReactionRepository(db *sqlx.DB) *ReactionRepository {
	return &ReactionRepository{DB: db}
}

type Reaction struct {
	ID         int       `db:"id"`
	UserName   string    `db:"user_name"`
	ReactionID int       `db:"reaction_id"`
	PostID     uuid.UUID `db:"post_id"`
	CreatedAt  time.Time `db:"created_at"`
}

func (rr *ReactionRepository) GetReactionsByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Reaction, error) {
	var postReactions []Reaction
	err := rr.DB.Select(&postReactions, "SELECT * FROM posts_reactions WHERE post_id = ? ORDER BY created_at DESC", postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions by post id: %w", err)
	}

	reactionUsersMap := make(map[int][]string, len(postReactions))
	for _, postReaction := range postReactions {
		if _, ok := reactionUsersMap[postReaction.ReactionID]; !ok {
			reactionUsersMap[postReaction.ReactionID] = make([]string, 0, 10)
		}
		reactionUsersMap[postReaction.ReactionID] = append(reactionUsersMap[postReaction.ReactionID], postReaction.UserName)
	}

	reactions := make([]*domain.Reaction, 0, len(reactionUsersMap))
	for reactionID, users := range reactionUsersMap {
		reactions = append(reactions, &domain.Reaction{
			PostID:     postID,
			ReactionID: reactionID,
			Users:      users,
			Count:      len(users),
		})
	}

	return reactions, nil
}