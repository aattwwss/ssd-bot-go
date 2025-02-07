package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aattwwss/ssd-bot-go/internal/config"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"

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
		err = run(context.Background(), cfg, rc, esRepo)
		if err != nil {
			log.Error().Msgf("Error during run: %v", err)
		}
		time.Sleep(15 * time.Minute)
	}
}

func run(ctx context.Context, cfg config.Config, rc *reddit.Client, esRepo *ssd.EsRepository) error {
	log.Info().Msg("Start searching...")
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
		title := cleanTitle(submission.Title)
		ssdList, err := esRepo.Search(ctx, title)
		if err != nil {
			log.Error().Msgf("Error searching for ssd: %v", err)
			continue
		}
		ssdList = sanityCheck(title, ssdList)
		if len(ssdList) == 0 {
			log.Info().Msgf("SSD not found in database: %s", submission.Title)
			continue
		}
		sort.Slice(ssdList, func(i, j int) bool {
			iName := strings.ReplaceAll(ssdList[i].Name, "(w/ Heatsink)", "")
			jName := strings.ReplaceAll(ssdList[i].Name, "(w/ Heatsink)", "")
			if len(iName) == len(jName) {
				numI, _ := strconv.Atoi(ssdList[i].DriveID)
				numJ, _ := strconv.Atoi(ssdList[j].DriveID)
				return numI > numJ
			}
			return len(iName) > len(jName)
		})
		found := ssdList[0]
		err = rc.SubmitComment(submission.ID, found.ToMarkdown())
		if err != nil {
			return err
		}
		log.Info().Msgf("Post submitted for: %v", found)
		//rate limit submission of post to prevent getting rejected
		time.Sleep(1 * time.Second)
	}
	log.Info().Msg("End searching...")
	return nil
}

// rules to ensure no false positives
// 1. Manufacturer must be in the search query
// 2. Name must be in the search query (without the heatsink part)
func sanityCheck(searchQuery string, ssds []ssd.SSD) []ssd.SSD {
	var filtered []ssd.SSD
	for _, ssd := range ssds {
		if !strings.Contains(strings.ToLower(strings.ReplaceAll(searchQuery, " ", "")), strings.ToLower(strings.ReplaceAll(ssd.Manufacturer, " ", ""))) {
			continue
		}
		ssdName := strings.ReplaceAll(ssd.Name, "(w/ Heatsink)", "")
		words := strings.Split(ssdName, " ")
		hasMissingWord := false
		for _, word := range words {
			if !strings.Contains(strings.ToLower(strings.ReplaceAll(searchQuery, " ", "")), strings.ToLower(strings.ReplaceAll(word, " ", ""))) {
				hasMissingWord = true
				break
			}
		}
		if hasMissingWord {
			continue
		}
		filtered = append(filtered, ssd)
	}
	return filtered
}

func cleanTitle(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`\[[^\]]+\]`).ReplaceAllString(s, "")
	stringsToRemove := []string{"ssd", "m2", "m.2", "nvme", "pcie", "gen"}
	for _, toReplace := range stringsToRemove {
		s = strings.ReplaceAll(s, toReplace, "")
	}

	stringsToReplace := map[string]string{
		" wd":        " western digital",
		"team group": " teamgroup",
		"spatium": " msi spatium",
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
