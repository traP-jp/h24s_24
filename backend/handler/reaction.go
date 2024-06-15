package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type ReactionRepository interface {
	GetReactionsByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Reaction, error)
	PostReaction(ctx context.Context, postID uuid.UUID, reactionID int, userName string) error
}

type ReactionHandler struct {
	rr ReactionRepository
}

type PostReactionResponse struct {
	ID    int `json:"id"`
	Count int `json:"count"`
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

	userName := "" //TODO: getUserNameが来たら実装

	err = rh.rr.PostReaction(ctx, postID, reactionID, userName)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // duplicate reaction
				return c.JSON(http.StatusConflict, "already reacted")
			}
			if mysqlErr.Number == 1452 { // post not found
				return c.JSON(http.StatusNotFound, "post not found")
			}
		}
		log.Println("failed to post reaction: ", err)
		return c.JSON(500, "failed to post reaction")
	}

	reactions, err := rh.rr.GetReactionsByPostID(ctx, postID)
	if err != nil {
		log.Println("failed to get reactions: ", err)
		return c.JSON(500, "failed to get reactions")
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
