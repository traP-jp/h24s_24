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

type reactionSlim struct {
	PostID     uuid.UUID `db:"post_id"`
	ReactionID int       `db:"reaction_id"`
}

func (rr *ReactionRepository) GetReactionsCount(ctx context.Context, since time.Time, timeSpan time.Duration, reactionLimit int) (map[uuid.UUID]int, error) {
	var reactions []reactionSlim
	err := rr.DB.Select(&reactions, "SELECT post_id, reaction_id FROM posts_reactions WHERE created_at BETWEEN ? and ? ORDER BY created_at DESC LIMIT ?", since, since.Add(timeSpan), reactionLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}

	scores := make(map[uuid.UUID]int)
	for _, r := range reactions {
		scores[r.PostID]++
	}

	return scores, nil
}

func (rr *ReactionRepository) PostReaction(ctx context.Context, postID uuid.UUID, reactionID int, userName string) error {
	_, err := rr.DB.Exec("INSERT INTO posts_reactions (post_id, reaction_id, user_name) VALUES (?, ?, ?)", postID, reactionID, userName)
	if err != nil {
		return fmt.Errorf("failed to post reaction: %w", err)
	}

	return nil
}

func (rr *ReactionRepository) GetReactionsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*domain.Reaction, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(
		"SELECT post_id, reaction_id, COUNT(*) as count FROM posts_reactions WHERE post_id IN (?) GROUP BY post_id, reaction_id",
		postIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := rr.DB.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions by post ids: %w", err)
	}
	defer rows.Close()

	reactionsMap := make(map[uuid.UUID][]*domain.Reaction, len(postIDs))
	for rows.Next() {
		var postID uuid.UUID
		var reactionID int
		var count int
		if err := rows.Scan(&postID, &reactionID, &count); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		if _, ok := reactionsMap[postID]; !ok {
			reactionsMap[postID] = make([]*domain.Reaction, 0, 10)
		}
		reactionsMap[postID] = append(reactionsMap[postID], &domain.Reaction{
			PostID:     postID,
			ReactionID: reactionID,
			Count:      count,
		})
	}

	return reactionsMap, nil
}

func (rr *ReactionRepository) GetReactionsByUserName(ctx context.Context, postID uuid.UUID, userName string) ([]*domain.UserReaction, error) {
	var userReactions []Reaction
	err := rr.DB.Select(&userReactions, "SELECT * FROM posts_reactions WHERE post_id = ? AND user_name = ?", postID, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions by user name: %w", err)
	}

	userReactionsDomain := make([]*domain.UserReaction, 0, len(userReactions))
	for _, userReaction := range userReactions {
		userReactionsDomain = append(userReactionsDomain, &domain.UserReaction{
			PostID:     userReaction.PostID,
			ReactionID: userReaction.ReactionID,
			CreatedAt:  userReaction.CreatedAt,
		})
	}

	return userReactionsDomain, nil
}
