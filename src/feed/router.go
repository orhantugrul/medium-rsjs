package feed

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/orhantugrul/medium-rsjs/src/util"
)

func BindRoutes(router *gin.RouterGroup) {
	router.Group("/feed").GET("/:username", getFeed)
}

func getFeed(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, util.Error{
			Code:    http.StatusBadRequest,
			Message: "Username required",
		})
		return
	}

	response, err := http.Get("https://medium.com/feed/" + username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, util.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch feed",
		})
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		context.JSON(http.StatusInternalServerError, util.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to read feed",
		})
		return
	}

	document, err := util.ParseDocument(body)
	if err != nil {
		context.JSON(http.StatusInternalServerError, util.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse feed",
		})
		return
	}

	feed := Feed{
		Title: document.Channel.Title,
		Link:  document.Channel.Link,
		Posts: make([]Post, 0, len(document.Channel.Items)),
	}

	for _, item := range document.Channel.Items {
		feed.Posts = append(feed.Posts, Post{
			Title:      item.Title,
			Link:       item.Link,
			Author:     item.Creator,
			Published:  item.PubDate,
			Content:    item.Content,
			Categories: item.Categories,
		})
	}

	context.JSON(http.StatusOK, feed)
}
