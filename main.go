package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aattwwss/ssd-bot-go/reddit"
)

const (
	SUBREDDIT   = "testingground4bots"
	LINK_PREFIX = "t3_"
)

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	username := os.Getenv("BOT_USERNAME")
	password := os.Getenv("BOT_PASSWORD")
	token := os.Getenv("BOT_ACCESS_TOKEN")
	expireTimeMilli, _ := strconv.ParseInt(os.Getenv("BOT_TOKEN_EXPIRE_MILLI"), 10, 64)
	rc, err := reddit.NewRedditClient(clientId, clientSecret, username, password, token, expireTimeMilli)
	if err != nil {
		log.Fatal(err)
	}
	posts, err := rc.GetNewPosts(SUBREDDIT, 10)
	if err != nil {
		log.Fatal(err)
	}

	botComments, err := rc.GetBotNewestComments(25)
	if err != nil {
		log.Fatal(err)
	}

	botCommentsMap := map[string]bool{}
	for _, comment := range botComments {
		linkId := strings.TrimPrefix(comment.LinkID, LINK_PREFIX)
		botCommentsMap[linkId] = true
	}

	for _, post := range posts {
		if !strings.Contains(strings.ToUpper(post.LinkFlairText), "SSD") {
			continue
		}
		_, ok := botCommentsMap[post.ID]
		if ok {
			continue
		}
		comments, err := rc.GetCommentsByPostId(post.ID, 10)
		if err != nil {
			log.Fatal(err)
		}
		for _, comment := range comments {
			fmt.Println(comment.Author)
		}
	}

	// err = rc.CreateComment("12ez9ws", content)

	if err != nil {
		log.Fatal(err)
	}

}
