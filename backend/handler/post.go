package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostRepository interface {
	CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, parentID uuid.UUID, rootID uuid.UUID) error
}

type PostHandler struct {
	PostRepository PostRepository
}

func (ph *PostHandler) PostPostsHandler(c echo.Context) error {
	// ctx := c.Request().Context()

	// ph.PostRepository.CreatePost()

	return nil
}
