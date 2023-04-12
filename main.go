package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aattwwss/ssd-bot-go/reddit"
	"github.com/aattwwss/ssd-bot-go/sheets"
	"github.com/aattwwss/ssd-bot-go/ssd"
	"github.com/rs/zerolog/log"
)

const (
	SUBREDDIT   = "buildapcsales"
	LINK_PREFIX = "t3_"

	SPREADSHEET_ID = "1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4"
	SHEET_NAME     = "'Master List'" //take note of the single quote, which is needed for sheets with space in them
)

func main() {
	err := run()
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return
	}
}
func run() error {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	username := os.Getenv("BOT_USERNAME")
	password := os.Getenv("BOT_PASSWORD")
	token := os.Getenv("BOT_ACCESS_TOKEN")
	expireTimeMilli, _ := strconv.ParseInt(os.Getenv("BOT_TOKEN_EXPIRE_MILLI"), 10, 64)
	isDebug := strings.ToUpper(os.Getenv("IS_DEBUG")) == "TRUE"

	sheetValues, err := sheets.GetSheetsValues(SPREADSHEET_ID, SHEET_NAME)
	if err != nil {
		return err
	}

	var allSSDs []ssd.SSD
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
		allSSDs = append(allSSDs, ssd)
	}

	rc, err := reddit.NewRedditClient(clientId, clientSecret, username, password, token, expireTimeMilli, isDebug)
	if err != nil {
		return err
	}
	posts, err := rc.GetNewPosts(SUBREDDIT, 25)
	if err != nil {
		return err
	}

	botComments, err := rc.GetBotNewestComments(25)
	if err != nil {
		return err
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
		log.Info().Msgf("Found post about SSD: %s", post.Title)
		found := search(allSSDs, post.Title)
		if found == nil {
			log.Info().Msgf("SSD not found in database: %s", post.Title)
			continue
		}

		log.Info().Msgf("Found in database: %v", found)
		err = rc.SubmitComment("12ez9ws", found.ToMarkdown())
		if err != nil {
			return err
		}
	}
	return nil
}

func getStringAtIndexOrEmpty(arr []interface{}, i int) string {
	if i >= len(arr) {
		return ""
	}
	return fmt.Sprintf("%v", arr[i])
}

// first try to match the branch, then match the model
func search(allSSDs []ssd.SSD, title string) *ssd.SSD {
	for _, ssd := range allSSDs {
		title = strings.ToUpper(title)
		brand := strings.ToUpper(ssd.Brand)
		model := strings.ToUpper(ssd.Model)
		if strings.Contains(title, brand) && strings.Contains(title, model) {
			return &ssd
		}
	}
	return nil
}
