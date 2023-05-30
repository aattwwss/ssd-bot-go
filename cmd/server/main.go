package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/aattwwss/ssd-bot-go/pkg/reddit"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
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

	config := config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}

	rc, err := reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.OverrideOldBot, config.IsDebug)
	if err != nil {
		log.Fatal().Msgf("Init reddit client error: %v", err)
		log.Info().Msgf("Init reddit client error: %v", rc)
	}
	es, _ := elasticutil.NewElasticsearchClient(config.EsAddress)
	esRepo := ssd.NewEsSSDRepository(es, "ssd-index")
	doTest(esRepo)
	//	run(context.Background(), config, rc, esRepo)
	//
	// tpuRepo := ssd.NewTpuSSDRepository(config.TPUHost, config.TPUUsername, config.TPUSecret)
	// sync(tpuRepo, esRepo)
}

func sync(source, dest ssd.SSDRepository) {
	ssdSync := ssd.SSDSync{
		StartId:  1,
		EndId:    1500,
		Delay:    time.Duration(10),
		IdToSkip: []int{},
	}
	err := ssdSync.Sync(context.Background(), source, dest)
	if err != nil {
		log.Fatal().Msgf("sync error", err)
	}
}

func run(ctx context.Context, config config, rc *reddit.RedditClient, esRepo *ssd.EsSSDRepository) error {
	newSubmissions, err := rc.GetNewSubmissions(config.Subreddit, 25)
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
		var botCommented bool
		comments, _ := rc.GetCommentsBySubmissionId(submission.ID, 100)
		for _, comment := range comments {
			botCommented = comment.Author == "SSDBot"
		}
		if config.OverrideOldBot && botCommented {
			log.Info().Msgf("Another bot already commented on this submission: %s", submission.Title)
			continue
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

type config struct {
	// reddit config
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Username     string `env:"BOT_USERNAME,notEmpty"`
	Password     string `env:"BOT_PASSWORD,notEmpty"`
	Subreddit    string `env:"SUBREDDIT,notEmpty"`

	// techpowerup config
	TPUHost     string `env:"TPU_HOST,notEmpty"`
	TPUUsername string `env:"TPU_USERNAME,notEmpty"`
	TPUSecret   string `env:"TPU_SECRET,notEmpty"`

	// elasticsearch config
	EsAddress string `env:"ES_ADDRESS,notEmpty"`

	// application config
	OverrideOldBot bool `env:"OVERRIDE_OLD_BOT,notEmpty"`

	//debugging config
	Token           string `env:"BOT_ACCESS_TOKEN"`
	ExpireTimeMilli int64  `env:"BOT_TOKEN_EXPIRE_MILLI"`
	IsDebug         bool   `env:"IS_DEBUG"`
}

func doTest(esRepo *ssd.EsSSDRepository) {
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

	// for {
	// 	// Read a line from the input CSV file
	// 	if err != nil {
	// 		break // End of file
	// 	}
	//
	// 	// Process the line (example: convert to uppercase)
	//
	// 	// Write the processed line to the output CSV file
	// 	err = writer.Write(line)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	// Flush the writer to ensure the line is written immediately
	// 	writer.Flush()
	//
	// 	// Check for any writer error
	// 	if err := writer.Error(); err != nil {
	// 		panic(err)
	// 	}
	//
	// 	// Print the processed line
	// 	fmt.Println(line)
	// }
}