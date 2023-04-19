package benchmark

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aattwwss/ssd-bot-go/config"
	"github.com/aattwwss/ssd-bot-go/reddit"
	"github.com/rs/zerolog/log"
)

func Benchmark(config config.Config) {
	// initPostsData(config)
	// read dataset.csv
	// parse into title and "brand + model"
	// for each title, run the search algo and see if match which the brand and model
	// calculate score e.g. if matched +1 else 0. Do some math
	// return score
	file, err := os.Open("benchmark/dataset.csv")
	if err != nil {
		log.Error().Msgf("open file error: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Error().Msgf("scanner error: %v", err)
		return
	}
}

func initPostsData(config config.Config) {
	rc, err := reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.IsDebug)
	if err != nil {
		log.Error().Msgf("Init reddit client error: %v", err)
		return
	}

	file, err := os.OpenFile("dataset.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Error().Msgf("Open file error: %v", err)
		return
	}
	defer file.Close()

	// Create a buffered writer to improve performance
	writer := bufio.NewWriter(file)

	posts := getLatestSSDPosts(rc)

	for _, post := range posts {
		ssd := getSSDFromPost(rc, post)
		data := fmt.Sprintf("%s\t%s\t%s\n", post.ID, post.Title, ssd)
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

func getSSDFromPost(rc *reddit.RedditClient, post reddit.Post) string {
	comments, err := rc.GetCommentsByPostId(post.ID, 100)
	if err != nil {
		log.Error().Msgf("GetCommentsByPostId error: %v", err)
		return ""
	}
	for _, comment := range comments {
		if strings.Contains(strings.ToUpper(comment.Author), "SSD") && strings.Contains(strings.ToUpper(comment.Author), "BOT") {
			res, err := substringBetween(comment.Body, "The ", " is a ")
			if err != nil {
				return ""
			}
			log.Info().Msgf("%s", res)
			return res
		}
	}
	return ""
}

func substringBetween(s, start, end string) (string, error) {
	startIndex := strings.Index(s, start)
	if startIndex == -1 {
		return "", fmt.Errorf("substring %q not found", start)
	}
	startIndex += len(start)
	endIndex := strings.Index(s[startIndex:], end)
	if endIndex == -1 {
		return "", fmt.Errorf("substring %q not found after index %d", end, startIndex)
	}
	endIndex += startIndex
	return s[startIndex:endIndex], nil
}
