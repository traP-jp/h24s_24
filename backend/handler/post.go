package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, parentID uuid.UUID, rootID uuid.UUID) error
	GetPosts(ctx context.Context, after uuid.UUID, limit int) ([]*domain.Post, error)
}

type PostHandler struct {
	PostRepository     PostRepository
	ReactionRepository ReactionRepository
}

func (ph *PostHandler) PostPostsHandler(c echo.Context) error {
	// ctx := c.Request().Context()

	// ph.PostRepository.CreatePost()

	return nil
}

type GetPostsResponse struct {
	ID               uuid.UUID       `json:"id"`
	UserName         string          `json:"userName"`
	OriginalMessage  string          `json:"originalMessage"`
	ConvertedMessage string          `json:"convertedMessage"`
	RootID           uuid.UUID       `json:"rootID,omitempty"`
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

	afterStr := c.QueryParam("after")
	after, err := uuid.Parse(afterStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid param 'after'")
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

	posts, err := ph.PostRepository.GetPosts(ctx, after, limit)
	if err != nil {
		log.Printf("failed to get posts: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get posts")
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
				ID:    reaction.PostID,
				Count: reaction.Count,
			}
			reactions = append(reactions, r)
		}
		r.Reactions = reactions
		res = append(res, r)
	}

	return c.JSON(http.StatusOK, res)
}
