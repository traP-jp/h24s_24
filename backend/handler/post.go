package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostRepository interface {
	CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, parentID uuid.UUID) error
}

type PostHandler struct {
	PostRepository PostRepository
}

type post struct {
	Message  string `json:"message"`
	ParentId string `json:"parent_id"`
}

func (ph *PostHandler) PostPostsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var post post
	err := c.Bind(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to bind: %v\n", err))
	}

	postID := uuid.New()
	var parentID uuid.UUID
	if len(post.ParentId) == 0 {
		parentID = postID
	} else {
		parentID, err = uuid.Parse(post.ParentId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse parent_id: %v\n", err))
		}
	}

	convertedMessage := post.Message
	return ph.PostRepository.CreatePost(ctx, postID, post.Message, convertedMessage, parentID)
}
