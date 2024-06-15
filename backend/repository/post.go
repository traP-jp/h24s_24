package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/traP-jp/h24s_24/domain"
)

type PostRepository struct {
	db *sqlx.DB
}

type post struct {
	ID               uuid.UUID `db:"id"`
	UserName         string    `db:"user_name"`
	OriginalMessage  string    `db:"original_message"`
	ConvertedMessage string    `db:"converted_message"`
	ParentID         uuid.UUID `db:"parent_id"`
	RootID           uuid.UUID `db:"root_id"`
	CreatedAt        time.Time `db:"created_at"`
}

func NewPostRepository(db *sqlx.DB) *PostRepository {

	return &PostRepository{db: db}
}

func (pr *PostRepository) CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, parentID uuid.UUID, rootID uuid.UUID) error {
	return nil
}

func (pr *PostRepository) GetPosts(ctx context.Context, after uuid.UUID, limit int) ([]*domain.Post, error) {
	var afterPost post
	err := pr.db.Get(&afterPost, "SELECT * FROM posts WHERE id = ?", after)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	var posts []post
	err = pr.db.Select(&posts, "SELECT * FROM posts WHERE created_at > ? ORDER BY created_at DESC LIMIT ?", afterPost.CreatedAt, limit)
	if err != nil {
		return nil, err
	}

	var domainPosts []*domain.Post
	for _, p := range posts {
		domainPosts = append(domainPosts, &domain.Post{
			ID:               p.ID,
			UserName:         p.UserName,
			OriginalMessage:  p.OriginalMessage,
			ConvertedMessage: p.ConvertedMessage,
			ParentID:         p.ParentID,
			RootID:           p.RootID,
			CreatedAt:        p.CreatedAt,
		})
	}

	return domainPosts, nil
}
