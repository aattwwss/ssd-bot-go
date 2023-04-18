package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aattwwss/ssd-bot-go/reddit"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}
	initPostsData(config)
}

func initPostsData(config Config) {
	rc, err := reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.IsDebug)
	if err != nil {
		log.Error().Msgf("Init reddit client error: %v", err)
		return
	}

	file, err := os.OpenFile("benchmark/posts.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Error().Msgf("Open file error: %v", err)
		return
	}
	defer file.Close()

	// Create a buffered writer to improve performance
	writer := bufio.NewWriter(file)

	posts := getLatestSSDPosts(rc)

	for _, post := range posts {
		data := fmt.Sprintf("%s|%s\n", post.ID, post.Title)
		_, err := writer.WriteString(data)
		if err != nil {
			log.Error().Msgf("Write to file error: %v", err)
			return
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Error().Msgf("Flush error: %v", err)
		return
	}
}
func getLatestSSDPosts(rc *reddit.RedditClient) []reddit.Post {
	total := 500
	pageSize := 100
	var res []reddit.Post
	q := "flair:SSD - Sata OR flair:SSD - M.2"

	opt := reddit.SearchPostOption{
		Subreddit:           "buildapcsales",
		Sort:                "new",
		Limit:               pageSize,
		Q:                   &q,
		RestrictToSubreddit: true,
		After:               nil,
	}

	titlesMap := make(map[string]bool)
	for i := 0; i < total; i = i + pageSize {
		log.Info().Msgf("querying: %v", pageSize)
		posts, err := rc.SearchPosts(opt)
		if err != nil {
			log.Error().Msgf("Search posts error: %v", err)
			return nil
		}
		if len(posts) == 0 {
			break
		}
		opt.After = &posts[len(posts)-1].Name
		for _, post := range posts {
			log.Info().Msgf("post: %v", post)
			_, ok := titlesMap[post.Title]
			if !ok {
				log.Info().Msgf("add new post: %v", post)
				res = append(res, post)
				titlesMap[post.Title] = true
			}
		}
	}
	return res
}
