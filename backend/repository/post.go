package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
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
	err = pr.db.Select(&posts, "SELECT * FROM posts WHERE created_at > ? ORDER BY created_at ASC LIMIT ?", afterPost.CreatedAt, limit)
	if err != nil {
		return nil, err
	}

	slices.Reverse(posts)

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

func (pr *PostRepository) GetPostsBefore(ctx context.Context, before uuid.UUID, limit int) ([]*domain.Post, error) {
	var beforePost post
	err := pr.db.Get(&beforePost, "SELECT * FROM posts WHERE id = ?", before)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	var posts []post
	err = pr.db.Select(&posts, "SELECT * FROM posts WHERE created_at < ? ORDER BY created_at DESC LIMIT ?", beforePost.CreatedAt, limit)
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

func (pr *PostRepository) GetPost(ctx context.Context, postID uuid.UUID) (*domain.Post, error) {
	var p post
	err := pr.db.Get(&p, "SELECT * FROM posts WHERE id = ?", postID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &domain.Post{
		ID:               p.ID,
		UserName:         p.UserName,
		OriginalMessage:  p.OriginalMessage,
		ConvertedMessage: p.ConvertedMessage,
		ParentID:         p.ParentID,
		RootID:           p.RootID,
		CreatedAt:        p.CreatedAt,
	}, nil
}

func (pr *PostRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.Post, error) {
	var posts []post
	err := pr.db.Select(&posts, "SELECT * FROM posts WHERE parent_id = ? AND id != ? ORDER BY created_at DESC", parentID, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
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

func (pr *PostRepository) GetAncestors(ctx context.Context, postID uuid.UUID) ([]*domain.Post, error) {
	var p post
	err := pr.db.Get(&p, "SELECT * FROM posts WHERE id = ?", postID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	if p.ID == p.ParentID {
		return []*domain.Post{}, nil
	}

	var posts []post
	err = pr.db.Select(&posts,
		"WITH RECURSIVE parent_cte AS ("+
			"SELECT id, user_name, original_message, converted_message, parent_id, root_id, created_at "+
			"FROM posts "+
			"WHERE id = ? "+
			"UNION ALL "+
			"SELECT p.id AS `id`, p.user_name AS `user_name`, p.original_message AS `original_message`, p.converted_message AS `converted_message`, p.parent_id AS `parent_id`, p.root_id AS `root_id`, p.created_at AS `created_at` "+
			"FROM posts p JOIN parent_cte c ON p.id = c.parent_id "+
			"WHERE p.id <> p.parent_id"+
			") "+
			"SELECT id, user_name, original_message, converted_message, parent_id, root_id, created_at "+
			"FROM posts WHERE id = ? AND id <> ? AND id <> ? "+
			"UNION ALL SELECT * FROM parent_cte ORDER BY created_at ASC;", p.ParentID, p.RootID, p.ID, p.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ancestors: %w", err)
	}

	domainPosts := make([]*domain.Post, 0, len(posts))
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

func (pr *PostRepository) GetChildrenCountByParentIDs(ctx context.Context, parentIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(parentIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In("SELECT parent_id, COUNT(*) FROM posts WHERE parent_id IN (?) GROUP BY parent_id", parentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	rows, err := pr.db.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	counts := make(map[uuid.UUID]int, len(parentIDs))
	for rows.Next() {
		var parentID uuid.UUID
		var count int
		err := rows.Scan(&parentID, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		counts[parentID] = count
	}

	return counts, nil
}

func (pr *PostRepository) GetRandomPost(ctx context.Context) (*domain.Post, error) {
	var p post
	err := pr.db.Get(&p, "SELECT * FROM posts ORDER BY RAND() LIMIT 1")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &domain.Post{
		ID:               p.ID,
		UserName:         p.UserName,
		OriginalMessage:  p.OriginalMessage,
		ConvertedMessage: p.ConvertedMessage,
		ParentID:         p.ParentID,
		RootID:           p.RootID,
		CreatedAt:        p.CreatedAt,
	}, nil
}
