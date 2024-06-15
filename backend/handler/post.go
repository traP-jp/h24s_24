package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, username string, parentID uuid.UUID) (uuid.UUID, error)
	GetPostsAfter(ctx context.Context, after uuid.UUID, limit int) ([]*domain.Post, error)
	GetLatestPosts(ctx context.Context, limit int) ([]*domain.Post, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*domain.Post, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.Post, error)
	GetAncestors(ctx context.Context, postID uuid.UUID) ([]*domain.Post, error)
	GetChildrenCountByParentIDs(ctx context.Context, parentIDs []uuid.UUID) (map[uuid.UUID]int, error)
}

type PostConverter interface {
	ConvertMessage(ctx context.Context, originalMessage string) (string, error)
}

type PostHandler struct {
	PostRepository     PostRepository
	ReactionRepository ReactionRepository

	pc PostConverter
}

type postPostsRequest struct {
	Message  string    `json:"message"`
	ParentID uuid.UUID `json:"parent_id"`
}

type postPostsResponse struct {
	OriginalMessage  string    `json:"original_message"`
	ConvertedMessage string    `json:"converted_message"`
	PostID           uuid.UUID `json:"id"`
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
	username, err = getUserName(c)
	if err != nil {
		log.Printf("failed to get username: %v\n", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get username")
	}

	postID := uuid.New()
	parentID := post.ParentID
	if parentID == uuid.Nil {
		parentID = postID
	}

	if post.Message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "message empty")
	}
	if utf8.RuneCountInString(post.Message) > 280 {
		return echo.NewHTTPError(http.StatusBadRequest, "message too long")
	}

	convertedMessage, err := ph.pc.ConvertMessage(ctx, post.Message)
	if err != nil {
		log.Printf("failed to convert message: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert message")
	}

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
	MyReactions      []int           `json:"myReactions"`
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

	loginUserName, err := getUserName(c)
	if err != nil {
		log.Printf("failed to get username: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get username")
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

		userReactions, err := ph.ReactionRepository.GetReactionsByUserName(ctx, post.ID, loginUserName)
		if err != nil {
			log.Printf("failed to get user reactions: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user reactions")
		}
		userReactionIDs := make([]int, 0, len(userReactions))
		for _, userReaction := range userReactions {
			userReactionIDs = append(userReactionIDs, userReaction.ReactionID)
		}
		r.MyReactions = userReactionIDs

		res = append(res, r)
	}

	return c.JSON(http.StatusOK, res)
}

type getPostResponse struct {
	ID               uuid.UUID       `json:"id"`
	UserName         string          `json:"user_name"`
	OriginalMessage  string          `json:"original_message"`
	ConvertedMessage string          `json:"converted_message"`
	MyReactions      []int           `json:"my_reactions"`
	Reactions        []reactionCount `json:"reactions"`
	Children         []postInfo      `json:"children"`
	Ancestors        []postInfo      `json:"ancestors"`
	CreatedAt        time.Time       `json:"created_at"`
}

type postInfo struct {
	Post struct {
		ID               uuid.UUID       `json:"id"`
		OriginalMessage  string          `json:"original_message"`
		ConvertedMessage string          `json:"converted_message"`
		UserName         string          `json:"user_name"`
		CreatedAt        time.Time       `json:"created_at"`
		Reactions        []reactionCount `json:"reactions"`
		MyReactions      []int           `json:"my_reactions"`
	} `json:"post"`
	ChildrenCount int `json:"children_count"`
}

func (ph *PostHandler) GetPostHandler(c echo.Context) error {
	ctx := c.Request().Context()
	loginUserName, err := getUserName(c)
	if err != nil {
		log.Printf("failed to get username: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get username")
	}

	postID, err := uuid.Parse(c.Param("postID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid post id")
	}

	post, err := ph.PostRepository.GetPost(ctx, postID)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, "post not found")
	}
	if err != nil {
		log.Printf("failed to get post: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	reactions, err := ph.ReactionRepository.GetReactionsByPostID(ctx, post.ID)
	if err != nil {
		log.Printf("failed to get reactions: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions")
	}

	resReactions := make([]reactionCount, 0, len(reactions))
	for _, re := range reactions {
		resReactions = append(resReactions, reactionCount{
			ID:    re.ReactionID,
			Count: re.Count,
		})
	}

	myReactions, err := ph.ReactionRepository.GetReactionsByUserName(ctx, post.ID, loginUserName)
	if err != nil {
		log.Printf("failed to get reactions by user name: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions by user name")
	}
	myReactionIDs := make([]int, 0, len(myReactions))
	for _, myReaction := range myReactions {
		myReactionIDs = append(myReactionIDs, myReaction.ReactionID)
	}

	res := getPostResponse{
		ID:               post.ID,
		UserName:         post.UserName,
		OriginalMessage:  post.OriginalMessage,
		ConvertedMessage: post.ConvertedMessage,
		Reactions:        resReactions,
		MyReactions:      myReactionIDs,
		CreatedAt:        post.CreatedAt,
	}

	children, err := ph.PostRepository.GetChildren(ctx, post.ID)
	if err != nil {
		log.Printf("failed to get children: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get children")
	}

	ancestors, err := ph.PostRepository.GetAncestors(ctx, post.ID)
	if err != nil {
		log.Printf("failed to get ancestors: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get ancestors")
	}

	result := make([][]postInfo, 0, 2)
	for _, posts := range [][]*domain.Post{children, ancestors} {
		postIDs := make([]uuid.UUID, 0, len(posts))
		for _, post := range posts {
			postIDs = append(postIDs, post.ID)
		}

		postIDReactionMap, err := ph.ReactionRepository.GetReactionsByPostIDs(ctx, postIDs)
		if err != nil {
			log.Printf("failed to get reactions by post ids: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions by post ids")
		}

		postIDChildrenCountMap, err := ph.PostRepository.GetChildrenCountByParentIDs(ctx, postIDs)
		if err != nil {
			log.Printf("failed to get posts count by parent ids: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get posts count by parent ids")
		}

		postInfos := make([]postInfo, 0, len(posts))
		for _, post := range posts {
			myReactions, err := ph.ReactionRepository.GetReactionsByUserName(ctx, post.ID, loginUserName)
			if err != nil {
				log.Printf("failed to get reactions by user name: %v\n", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions by user name")
			}
			myReactionIDs := make([]int, 0, len(myReactions))
			for _, myReaction := range myReactions {
				myReactionIDs = append(myReactionIDs, myReaction.ReactionID)
			}

			pInfo := postInfo{
				Post: struct {
					ID               uuid.UUID       `json:"id"`
					OriginalMessage  string          `json:"original_message"`
					ConvertedMessage string          `json:"converted_message"`
					UserName         string          `json:"user_name"`
					CreatedAt        time.Time       `json:"created_at"`
					Reactions        []reactionCount `json:"reactions"`
					MyReactions      []int           `json:"my_reactions"`
				}{
					ID:               post.ID,
					OriginalMessage:  post.OriginalMessage,
					ConvertedMessage: post.ConvertedMessage,
					UserName:         post.UserName,
					CreatedAt:        post.CreatedAt,
					MyReactions:      myReactionIDs,
				},
				ChildrenCount: postIDChildrenCountMap[post.ID],
			}

			reactions, ok := postIDReactionMap[post.ID]
			if !ok {
				pInfo.Post.Reactions = []reactionCount{}
				postInfos = append(postInfos, pInfo)
				continue
			}

			reCount := make([]reactionCount, 0, len(reactions))
			for _, re := range reactions {
				reCount = append(reCount, reactionCount{
					ID:    re.ReactionID,
					Count: re.Count,
				})
			}
			pInfo.Post.Reactions = reCount

			postInfos = append(postInfos, pInfo)
		}

		result = append(result, postInfos)
	}

	res.Children = result[0]
	res.Ancestors = result[1]

	return c.JSON(http.StatusOK, res)
}
