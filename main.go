package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aattwwss/ssd-bot-go/reddit"
	"github.com/aattwwss/ssd-bot-go/search"
	"github.com/aattwwss/ssd-bot-go/sheets"
	"github.com/aattwwss/ssd-bot-go/ssd"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	SUBREDDIT   = "buildapcsales"
	LINK_PREFIX = "t3_"

	SPREADSHEET_ID = "1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4"
	SHEET_NAME     = "'Master List'" //take note of the single quote, which is needed for sheets with space in them
)

type Config struct {
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Username     string `env:"BOT_USERNAME,notEmpty"`
	Password     string `env:"BOT_PASSWORD,notEmpty"`

	Token           string `env:"BOT_ACCESS_TOKEN"`
	ExpireTimeMilli int64  `env:"BOT_TOKEN_EXPIRE_MILLI"`
	IsDebug         bool   `env:"IS_DEBUG"`
}

func newConfig(clientId, clientSecret, username, password, token string, expireTimeMilli int64, isDebug bool) (*Config, error) {
	if clientId == "" || clientSecret == "" || username == "" || password == "" {
		return nil, errors.New("clientId, clientSecret, username and password cannot be empty")
	}

	config := Config{
		ClientId:        clientId,
		ClientSecret:    clientSecret,
		Username:        username,
		Password:        password,
		Token:           token,
		ExpireTimeMilli: expireTimeMilli,
		IsDebug:         isDebug,
	}

	if isDebug {
		config.IsDebug = isDebug
		config.Token = token
		config.ExpireTimeMilli = expireTimeMilli
	}

	return &config, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}

	rc, err := reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.IsDebug)
	if err != nil {
		log.Error().Msgf("Init reddit client error: %v", err)
		return
	}

	for {
		log.Info().Msgf("Scanning...")
		count, err := run(config, rc)
		if err != nil {
			log.Error().Msgf("Run error: %v", err)
		}
		log.Info().Msgf("Updated %v posts...", count)
		time.Sleep(10 * time.Minute)
	}
}

func run(config Config, rc *reddit.RedditClient) (int, error) {
	sheetValues, err := sheets.GetSheetsValues(SPREADSHEET_ID, SHEET_NAME)
	if err != nil {
		return 0, err
	}

	var allSSDs []ssd.SSD
	var searchDocuments []string
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
		searchDocuments = append(searchDocuments, search.ReplaceSpecialChar(strings.ToUpper(ssd.Brand)+" "+strings.ToUpper(ssd.Model), " "))
		allSSDs = append(allSSDs, ssd)
	}
	tfidf := search.NewTfIdf(searchDocuments)

	posts, err := rc.GetNewPosts(SUBREDDIT, 25)
	if err != nil {
		return 0, err
	}

	botComments, err := rc.GetBotNewestComments(25)
	if err != nil {
		return 0, err
	}

	botCommentsMap := map[string]bool{}
	for _, comment := range botComments {
		linkId := strings.TrimPrefix(comment.LinkID, LINK_PREFIX)
		botCommentsMap[linkId] = true
	}

	count := 0
	for _, post := range posts {
		if !strings.Contains(strings.ToUpper(post.LinkFlairText), "SSD") {
			continue
		}
		_, ok := botCommentsMap[post.ID]
		if ok {
			log.Info().Msgf("Already commented on this post: %s", post.Title)
			continue
		}
		log.Info().Msgf("Found post: %s", post.Title)
		// found := matchSsd(allSSDs, tfidf, post.Title)
		found := searchSsd(allSSDs, post.Title, tfidf)
		if found == nil {
			log.Info().Msgf("SSD not found in database: %s", post.Title)
			continue
		}

		err = rc.SubmitComment(post.ID, found.ToMarkdown())
		if err != nil {
			return 0, err
		}
		log.Info().Msgf("Post submitted for: %v", found)
		count++
		//rate limit submission of post to prevent getting rejected
		time.Sleep(10 * time.Second)
	}
	return count, nil
}

func getStringAtIndexOrEmpty(arr []interface{}, i int) string {
	if i >= len(arr) {
		return ""
	}
	return fmt.Sprintf("%v", arr[i])
}

// very naive searching algorithm for now
// first try to match the branch, then match the model
func matchSsd(allSSDs []ssd.SSD, tfidf *search.TfIdf, title string) *ssd.SSD {
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

// using tfidf to find the most relevant ssd
func searchSsd(allSSDs []ssd.SSD, postTitle string, tfidf *search.TfIdf) *ssd.SSD {
	postTitle = cleanTitle(postTitle)

	terms := strings.Fields(postTitle)
	scores := make([]float64, len(tfidf.Documents))
	for i := range tfidf.Documents {
		for _, term := range terms {
			scores[i] += tfidf.TfIdf(term, i)
		}
	}
	maxScore := 0.0
	maxIndex := 0
	for i, score := range scores {
		if score > maxScore {
			maxScore = score
			maxIndex = i
		}
	}

	// if not relevant at all, or if the post title does not contain the brand
	if maxScore == 0 || !strings.Contains(postTitle, strings.ToUpper(allSSDs[maxIndex].Brand)) || !strings.Contains(postTitle, strings.ToUpper(allSSDs[maxIndex].Model)) {
		log.Info().Msgf("Reject found ssd. score: %v, Title: %s, SSD: %v", maxScore, postTitle, allSSDs[maxIndex])
		return nil
	}

	return &allSSDs[maxIndex]
}

//temporary fix for some misalignment in post titles and the google sheets brands and models
func cleanTitle(title string) string {

	title = strings.ToUpper(title)
	replaceRules := map[string]string{
		"TEAMGROUP": "TEAMGROUP TEAM GROUP",
		"XPG":       "ADATA XPG",
	}

	for k := range replaceRules {
		if v, ok := replaceRules[k]; ok {
			title = strings.ReplaceAll(title, k, v)
		}
	}

	return title
}
