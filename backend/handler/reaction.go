package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type ReactionRepository interface {
	DeleteReaction(ctx context.Context, postID uuid.UUID, reactionID int, userName string) ([]*domain.Reaction, error)
	GetReactionsByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Reaction, error)
	GetReactionsByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]*domain.Reaction, error)
	GetReactionCountsGroupedByPostID(ctx context.Context, reactionID *int, since time.Time, until time.Time) ([]*domain.ReactionCount, error)
	PostReaction(ctx context.Context, postID uuid.UUID, reactionID int, userName string) error
	GetReactionsByUserName(ctx context.Context, postID uuid.UUID, userName string) ([]*domain.UserReaction, error)
}

type ReactionHandler struct {
	rr ReactionRepository
}

type PostReactionResponse struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type deleteReactionResponse struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

func (rh *ReactionHandler) DeleteReactionHandler(c echo.Context) error {
	ctx := c.Request().Context()

	postID, err := uuid.Parse(c.Param("postID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post id")
	}

	reactionID, err := strconv.Atoi(c.Param("reactionID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid reaction id")
	}

	username, err := getUserName(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get username")
	}

	reactions, err := rh.rr.DeleteReaction(ctx, postID, reactionID, username)
	if errors.Is(err, domain.ReactionNotFoundError) {
		return echo.NewHTTPError(http.StatusBadRequest, "reaction not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to delete reaction: %v", err))
	}

	response := make([]*deleteReactionResponse, len(reactions))
	for i, r := range reactions {
		response[i] = &deleteReactionResponse{r.ReactionID, r.Count}
	}
	return c.JSON(http.StatusOK, response)
}

func (rh *ReactionHandler) PostReactionHandler(c echo.Context) error {
	ctx := c.Request().Context()

	postID, err := uuid.Parse(c.Param("postID"))
	if err != nil {
		return c.JSON(400, "invalid post id")
	}

	reactionID, err := strconv.Atoi(c.Param("reactionID"))
	if err != nil {
		return c.JSON(400, "invalid reaction id")
	}

	userName, err := getUserName(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get username")
	}

	err = rh.rr.PostReaction(ctx, postID, reactionID, userName)
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		if me.Number == 1062 { // duplicate reaction
			return c.JSON(http.StatusConflict, "already reacted")
		}
		if me.Number == 1452 { // post not found
			return c.JSON(http.StatusNotFound, "post not found")
		}
	}
	if err != nil {
		log.Println("failed to post reaction: ", err)
		return c.JSON(http.StatusInternalServerError, "failed to post reaction")
	}

	reactions, err := rh.rr.GetReactionsByPostID(ctx, postID)
	if err != nil {
		log.Println("failed to get reactions: ", err)
		return c.JSON(http.StatusInternalServerError, "failed to get reactions")
	}

	res := make([]PostReactionResponse, 0, len(reactions))

	for _, reaction := range reactions {
		res = append(res, PostReactionResponse{
			ID:    reaction.ReactionID,
			Count: reaction.Count,
		})
	}

	return c.JSON(http.StatusOK, res)
}
