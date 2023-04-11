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

	arr := []string{
		fmt.Sprintf("The Crucial P1 is a *QLC* **Budget NVMe** SSD."),
		fmt.Sprintf("* Interface: **x4 PCIe 3.0/NVMe**"),
		fmt.Sprintf("* Form Factor: **M.2**"),
		fmt.Sprintf("* Controller: **SMI SM2263**"),
		fmt.Sprintf("* Configuration: **Dual-core, 4-ch, 4-CE/ch**"),
		fmt.Sprintf("* DRAM: **Yes**"),
		fmt.Sprintf("* HMB: **No**"),
		fmt.Sprintf("* NAND Brand: **Micron**"),
		fmt.Sprintf("* NAND Type: **QLC**"),
		fmt.Sprintf("* 2D/3D NAND: **3D**"),
		fmt.Sprintf("* Layers: **64**"),
		fmt.Sprintf("* R/W: **2000/1700**"),
		fmt.Sprintf("[Click here to view this SSD in the tier list](https://docs.google.com/spreadsheets/d/1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4/edit#gid=0&amp;range=A63:V63),"),
		fmt.Sprintf("[Click here to view camelcamelcamel product search page](https://camelcamelcamel.com/search?sq=),."),
		fmt.Sprintf("---\n^(Suggestions, concerns, errors? Message us directly or submit an issue on), [^(Github!)](https://github.com/aattwwss/ssd-bot-go)"),
	}
	content := strings.Join(arr, "\n\n")

	err = rc.CreateComment("12ez9ws", content)

	if err != nil {
		log.Fatal(err)
	}

}
