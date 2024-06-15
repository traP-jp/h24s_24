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

func (pr *PostRepository) CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, username string, parentID uuid.UUID) (uuid.UUID, error) {
	db := pr.db
	var rootID uuid.UUID

	if postID == parentID { // リプライじゃない
		rootID = postID
	} else {
		err := db.Get(&rootID, "SELECT root_id FROM posts WHERE id=?", parentID)
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("not found: %w", err)
		}
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to get %w", err)
		}
	}

	_, err := db.Exec("INSERT INTO posts (id, original_message, converted_message, user_name, parent_id, root_id) VALUES (?, ?, ?, ?, ?, ?)", postID, originalMessage, convertedMessage, username, parentID, rootID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert: %w", err)
	}
	return rootID, nil
}

func (pr *PostRepository) GetPostsAfter(ctx context.Context, after uuid.UUID, limit int) ([]*domain.Post, error) {
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

func (pr *PostRepository) GetLatestPosts(ctx context.Context, limit int) ([]*domain.Post, error) {
	var posts []post
	err := pr.db.Select(&posts, "SELECT * FROM posts ORDER BY created_at DESC LIMIT ?", limit)
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

func (pr *PostRepository) GetPostByID(c context.Context, postID uuid.UUID) (*domain.Post, error) {
	var p post
	err := pr.db.Select(&p, "SELECT * FROM posts WHERE id=?", postID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	post := &domain.Post{
		ID:               p.ID,
		UserName:         p.UserName,
		OriginalMessage:  p.OriginalMessage,
		ConvertedMessage: p.ConvertedMessage,
		ParentID:         p.ParentID,
		RootID:           p.RootID,
		CreatedAt:        p.CreatedAt,
	}
	return post, nil
}
