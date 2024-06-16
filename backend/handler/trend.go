package handler

import (
	"cmp"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type TrendHandler struct {
	pr PostRepository
	rr ReactionRepository
}

type reaction struct {
	ReactionID int `json:"id"`
	Count      int `json:"count"`
}

type getTrendResponse struct {
	PostID           uuid.UUID   `json:"id"`
	Username         string      `json:"user_name"`
	OriginalMessage  string      `json:"original_message"`
	ConvertedMessage string      `json:"converted_message"`
	Reactions        []*reaction `json:"reactions"`
	RootID           uuid.UUID   `json:"root_id"`
	CreatedAt        time.Time   `json:"created_at"`
	MyReactions      []int       `json:"my_reactions"`
}

func (tr *TrendHandler) GetTrendHandler(c echo.Context) error {
	until := time.Now()
	since := until.Add(-time.Hour * 30)
	postLimit := 30

	ctx := c.Request().Context()

	reactionIDRef := (*int)(nil)
	if reactionIDString := c.QueryParams().Get("reaction_id"); reactionIDString != "" {
		reactionID, err := strconv.Atoi(reactionIDString)
		if err != nil {
			log.Printf("failed to parse reaction_id: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "failed to parse reaction_id")
		}
		reactionIDRef = &reactionID
	}

	counts, err := tr.rr.GetReactionCountsGroupedByPostID(ctx, reactionIDRef, since, until)
	if err != nil {
		log.Printf("failed to get reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions")
	}
	slices.SortFunc(counts, func(a, b *domain.ReactionCount) int {
		n := b.Count - a.Count
		if n != 0 {
			return n
		}
		return cmp.Compare(a.PostID.ID(), b.PostID.ID())
	})

	postIDs := make([]uuid.UUID, len(counts))
	for i, v := range counts {
		postIDs[i] = v.PostID
	}

	reactionsMap, err := tr.rr.GetReactionsByPostIDs(ctx, postIDs)
	if err != nil {
		log.Printf("failed to get reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions")
	}

	posts := make([]*getTrendResponse, min(postLimit, len(counts)))
	for i, v := range counts {
		post, err := tr.pr.GetPost(ctx, v.PostID)
		if err != nil {
			log.Printf("failed to get post: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
		}

		reactions := reactionsMap[v.PostID]
		reactionsSlice := make([]*reaction, len(reactions))
		for i, v := range reactions {
			reactionsSlice[i] = &reaction{v.ReactionID, v.Count}
		}

		myReactions, err := tr.rr.GetReactionsByUserName(ctx, post.ID, post.UserName)
		if err != nil {
			log.Printf("failed to get my reactions: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get my reaction")
		}
		myReactionsSlice := make([]int, len(myReactions))
		for i, userReaction := range myReactions {
			myReactionsSlice[i] = userReaction.ReactionID
		}

		posts[i] = &getTrendResponse{
			PostID:           post.ID,
			Username:         post.UserName,
			OriginalMessage:  post.OriginalMessage,
			ConvertedMessage: post.ConvertedMessage,
			Reactions:        reactionsSlice,
			RootID:           post.RootID,
			CreatedAt:        post.CreatedAt.Local(),
			MyReactions:      myReactionsSlice,
		}
	}
	return c.JSON(http.StatusOK, posts)
}
