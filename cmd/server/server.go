package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/aattwwss/ssd-bot-go/internal/config"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/aattwwss/ssd-bot-go/pkg/reddit"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	LINK_PREFIX = "t3_"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}

	rc, err := reddit.NewRedditClient(cfg.ClientId, cfg.ClientSecret, cfg.Username, cfg.Password, cfg.Token, cfg.ExpireTimeMilli, cfg.OverrideOldBot)
	if err != nil {
		log.Fatal().Msgf("Init reddit client error: %v", err)
	}
	es, _ := elasticutil.NewElasticsearchClient(cfg.EsAddress)
	esRepo := ssd.NewEsRepository(es, "ssd-index")
	// doTest(esRepo)
	for {
		log.Info().Msg("Start searching...")
		err = run(context.Background(), cfg, rc, esRepo)
		if err != nil {
			log.Error().Msgf("Error during run: %v", err)
		}
		log.Info().Msg("End searching...")
		time.Sleep(15 * time.Minute)
	}
}

func run(ctx context.Context, cfg config.Config, rc *reddit.Client, esRepo *ssd.EsRepository) error {
	newSubmissions, err := rc.GetNewSubmissions(cfg.Subreddit, 25)
	if err != nil {
		return err
	}

	botComments, err := rc.GetUserNewestComments(25)
	if err != nil {
		return err
	}

	botCommentsMap := map[string]bool{}
	for _, comment := range botComments {
		linkId := strings.TrimPrefix(comment.LinkID, LINK_PREFIX)
		botCommentsMap[linkId] = true
	}

	for _, submission := range newSubmissions {
		if !strings.Contains(strings.ToUpper(submission.LinkFlairText), "SSD") {
			continue
		}
		if !cfg.OverrideOldBot {
			// do not comment if another bot already commented
			botToCheck := "SSDBot"
			botCommented := rc.IsCommentedByUser(submission.ID, botToCheck)
			if botCommented {
				log.Info().Msgf("%s already commented on this submission: %s", botToCheck, submission.Title)
				continue
			}
		}

		_, ok := botCommentsMap[submission.ID]
		if ok {
			log.Info().Msgf("This bot already commented on this submission: %s", submission.Title)
			continue
		}

		log.Info().Msgf("Found submission: %s", submission.Title)
		ssdList, err := esRepo.Search(ctx, cleanTitle(submission.Title))
		if err != nil {
			log.Error().Msgf("Error searching for ssd: %v", err)
			continue
		}
		if len(ssdList) == 0 {
			log.Info().Msgf("SSD not found in database: %s", submission.Title)
			continue
		}
		found := ssdList[0]
		err = rc.SubmitComment(submission.ID, found.ToMarkdown())
		if err != nil {
			return err
		}
		log.Info().Msgf("Post submitted for: %v", found)
		//rate limit submission of post to prevent getting rejected
		time.Sleep(1 * time.Second)
	}

	return nil
}

func cleanTitle(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`\[[^\]]+\]`).ReplaceAllString(s, "")
	stringsToRemove := []string{"ssd", "m2", "m.2", "nvme", "pcie", "gen"}
	for _, toReplace := range stringsToRemove {
		s = strings.ReplaceAll(s, toReplace, "")
	}

	stringsToReplace := map[string]string{
		" wd": " western digital",
	}
	for k, v := range stringsToReplace {
		if strings.Contains(s, k) {
			s = s + v
		}
	}
	return s
}

func doTest(esRepo *ssd.EsRepository) {
	// Open the input CSV file for reading
	inputFile, err := os.Open("test/input.csv")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// Open a temporary output file for writing
	outputFile, err := os.Create("test/output.csv")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Read all the records
	reader := csv.NewReader(inputFile)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(outputFile)
	writer.Comma = '|'

	// Print each record
	for _, record := range records {
		ssds, _ := esRepo.Search(context.Background(), cleanTitle(record[1]))
		var modelAndName string
		if len(ssds) != 0 {
			modelAndName = fmt.Sprintf("%s %s", ssds[0].Manufacturer, ssds[0].Name)
		}
		record = append(record, modelAndName)
		writer.Write(record)
		writer.Flush()
	}
}
