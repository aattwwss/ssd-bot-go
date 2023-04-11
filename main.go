package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aattwwss/ssd-bot-go/reddit"
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
	posts, err := rc.GetNewPosts("buildapcsales", 100)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range posts {
		if !strings.Contains(strings.ToUpper(post.Data.LinkFlairText), "SSD") {
			continue
		}
		fmt.Println(post.Data.ID)
	}
}
