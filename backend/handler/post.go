package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, username string, parentID uuid.UUID) (uuid.UUID, error)
	GetPostsAfter(ctx context.Context, after uuid.UUID, limit int) ([]*domain.Post, error)
	GetLatestPosts(ctx context.Context, limit int) ([]*domain.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error)
}

type PostHandler struct {
	PostRepository     PostRepository
	ReactionRepository ReactionRepository
}

type postPostsRequest struct {
	Message  string    `json:"message"`
	ParentID uuid.UUID `json:"parent_id"`
}

type postPostsResponse struct {
	OriginalMessage  string    `json:"original_message"`
	ConvertedMessage string    `json:"converted_message"`
	PostID           uuid.UUID `json:"post_id"`
	CreatedAt        time.Time `json:"created_at"`
	ParentID         uuid.UUID `json:"parent_id"`
	RootID           uuid.UUID `json:"root_id"`
}

func (ph *PostHandler) PostPostsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var post postPostsRequest
	err := c.Bind(&post)
	if err != nil {
		log.Printf("failed to bind: %v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind")
	}

	var username string
	username, err = getUsername(c)
	if err != nil {
		log.Printf("failed to get username: %v\n", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get username")
	}

	postID := uuid.New()
	parentID := post.ParentID
	if parentID == uuid.Nil {
		parentID = postID
	}

	convertedMessage := post.Message

	var rootID uuid.UUID
	rootID, err = ph.PostRepository.CreatePost(ctx, postID, post.Message, convertedMessage, username, parentID)

	if err != nil {
		log.Printf("failed to post: %v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "failed to post")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to post")
	}

	return c.JSON(http.StatusOK, postPostsResponse{
		OriginalMessage:  post.Message,
		ConvertedMessage: convertedMessage,
		PostID:           postID,
		CreatedAt:        time.Now(),
		ParentID:         parentID,
		RootID:           rootID,
	})
}

type GetPostsResponse struct {
	ID               uuid.UUID       `json:"id"`
	UserName         string          `json:"userName"`
	OriginalMessage  string          `json:"originalMessage"`
	ConvertedMessage string          `json:"convertedMessage"`
	RootID           uuid.UUID       `json:"rootID"`
	Reactions        []reactionCount `json:"reactions"`
	CreatedAt        time.Time       `json:"createdAt"`
}

type reactionCount struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

func (ph *PostHandler) GetPostsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	repostStr := c.QueryParam("repost")
	repost, err := strconv.ParseBool(repostStr)
	if err != nil {
		repost = false
	}

	if repost {
		return echo.NewHTTPError(http.StatusNotImplemented, "repost is not implemented")
	}

	useAfter := false
	afterStr := c.QueryParam("after")
	var after uuid.UUID
	if afterStr != "" {
		useAfter = true
		after, err = uuid.Parse(afterStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid param 'after'")
		}
	}

	var limit int
	limitStr := c.QueryParam("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid param 'limit'")
		}
	} else {
		limit = 30
	}

	var posts []*domain.Post
	if useAfter {
		posts, err = ph.PostRepository.GetPostsAfter(ctx, after, limit)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, "after post not found")
		}
		if err != nil {
			log.Printf("failed to get posts: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get posts")
		}
	} else {
		posts, err = ph.PostRepository.GetLatestPosts(ctx, limit)
		if err != nil {
			log.Printf("failed to get posts: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get posts")
		}
	}

	postReactionsMap := make(map[uuid.UUID][]*domain.Reaction)

	for _, post := range posts {
		reactions, err := ph.ReactionRepository.GetReactionsByPostID(ctx, post.ID)
		if err != nil {
			log.Printf("failed to get reactions: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions")
		}
		postReactionsMap[post.ID] = reactions
	}

	res := make([]GetPostsResponse, 0, len(posts))
	for _, post := range posts {
		r := GetPostsResponse{
			ID:               post.ID,
			UserName:         post.UserName,
			OriginalMessage:  post.OriginalMessage,
			ConvertedMessage: post.ConvertedMessage,
			RootID:           post.RootID,
			CreatedAt:        post.CreatedAt,
		}
		reactions := make([]reactionCount, 0, len(postReactionsMap[post.ID]))
		for _, reaction := range postReactionsMap[post.ID] {
			r := reactionCount{
				ID:    reaction.ReactionID,
				Count: reaction.Count,
			}
			reactions = append(reactions, r)
		}
		r.Reactions = reactions
		res = append(res, r)
	}

	return c.JSON(http.StatusOK, res)
}
