package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aattwwss/ssd-bot-go/reddit"
	"github.com/aattwwss/ssd-bot-go/sheets"
	"github.com/aattwwss/ssd-bot-go/ssd"
)

const (
	SUBREDDIT   = "testingground4bots"
	LINK_PREFIX = "t3_"

	SPREADSHEET_ID = "1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4"
	SHEET_NAME     = "'Master List'" //take note of the single quote, which is needed for sheets with space in them
)

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	username := os.Getenv("BOT_USERNAME")
	password := os.Getenv("BOT_PASSWORD")
	token := os.Getenv("BOT_ACCESS_TOKEN")
	expireTimeMilli, _ := strconv.ParseInt(os.Getenv("BOT_TOKEN_EXPIRE_MILLI"), 10, 64)
	isDebug := strings.ToUpper(os.Getenv("IS_DEBUG")) == "TRUE"

	rc, err := reddit.NewRedditClient(clientId, clientSecret, username, password, token, expireTimeMilli, isDebug)
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
	sheetValues, err := sheets.GetSheetsValues(SPREADSHEET_ID, SHEET_NAME)
	if err != nil {
		log.Fatal(err)
	}
	var allSSD []ssd.SSD
	for i, row := range sheetValues {
		// skip the header
		if i == 0 {
			continue
		}
		// break at the end of the list of data
		if len(row) == 0 {
			break
		}

		ssd := ssd.SSD{
			Brand:         getStringAtIndexOrEmpty(row, 0),
			Model:         getStringAtIndexOrEmpty(row, 1),
			Interface:     getStringAtIndexOrEmpty(row, 2),
			FormFactor:    getStringAtIndexOrEmpty(row, 3),
			Capacity:      getStringAtIndexOrEmpty(row, 4),
			Controller:    getStringAtIndexOrEmpty(row, 5),
			Configuration: getStringAtIndexOrEmpty(row, 6),
			DRAM:          getStringAtIndexOrEmpty(row, 7),
			HMB:           getStringAtIndexOrEmpty(row, 8),
			NandBrand:     getStringAtIndexOrEmpty(row, 9),
			NandType:      getStringAtIndexOrEmpty(row, 10),
			Layers:        getStringAtIndexOrEmpty(row, 11),
			ReadWrite:     getStringAtIndexOrEmpty(row, 12),
			Category:      getStringAtIndexOrEmpty(row, 13),
			CellRow:       i + 1,
		}
		allSSD = append(allSSD, ssd)
	}
	fmt.Println(allSSD)

}

func getStringAtIndexOrEmpty(arr []interface{}, i int) string {
	if i >= len(arr) {
		return ""
	}
	return fmt.Sprintf("%v", arr[i])
}
