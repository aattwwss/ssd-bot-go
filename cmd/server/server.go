package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
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
	LINK_PREFIX       = "t3_"
	ES_INDEX          = "ssd-index"
	POLL_INTERVAL     = 15 * time.Minute
	COMMENT_RATE_LIMIT = 1 * time.Second
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
	es, err := elasticutil.NewElasticsearchClient(cfg.EsAddress)
	if err != nil {
		log.Fatal().Msgf("Init elasticsearch client error: %v", err)
	}
	esRepo := ssd.NewEsRepository(es, ES_INDEX)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a goroutine to handle shutdown signals
	go func() {
		sig := <-sigChan
		log.Info().Msgf("Received signal %v, shutting down...", sig)
		cancel()
	}()

	// doTest(esRepo)
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Shutdown requested, exiting...")
			return
		default:
			err = run(ctx, cfg, rc, esRepo)
			if err != nil {
				log.Error().Msgf("Error during run: %v", err)
			}
			time.Sleep(POLL_INTERVAL)
		}
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
			jName := strings.ReplaceAll(ssdList[j].Name, "(w/ Heatsink)", "")
			if len(iName) == len(jName) {
				numI, errI := strconv.Atoi(ssdList[i].DriveID)
				numJ, errJ := strconv.Atoi(ssdList[j].DriveID)
				// If both are valid integers, compare numerically
				if errI == nil && errJ == nil {
					return numI > numJ
				}
				// If only one is valid, prefer the valid one
				if errI == nil {
					return true
				}
				if errJ == nil {
					return false
				}
				// If neither is valid, fall back to string comparison
				return ssdList[i].DriveID > ssdList[j].DriveID
			}
			return len(iName) > len(jName)
		})
		log.Info().Msgf("Final sorted filtered list %v", ssdList)
		found := ssdList[0]
		err = rc.SubmitComment(submission.ID, found.ToMarkdown())
		if err != nil {
			return err
		}
		log.Info().Msgf("Post submitted for: %v", found)
		//rate limit submission of post to prevent getting rejected
		time.Sleep(COMMENT_RATE_LIMIT)
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
		log.Debug().Msgf("checking %s %s", ssd.Manufacturer, ssd.Name)
		if !strings.Contains(strings.ToLower(strings.ReplaceAll(searchQuery, " ", "")), strings.ToLower(strings.ReplaceAll(ssd.Manufacturer, " ", ""))) {
			log.Debug().Msgf("skipping %s %s because manufacturer is missing from search query", ssd.Manufacturer, ssd.Name)
			continue
		}
		ssdName := strings.ReplaceAll(ssd.Name, "(w/ Heatsink)", "")
		words := strings.Split(ssdName, " ")
		hasMissingWord := false
		for _, word := range words {
			if !strings.Contains(strings.ToLower(strings.ReplaceAll(searchQuery, " ", "")), strings.ToLower(strings.ReplaceAll(word, " ", ""))) {
				hasMissingWord = true
				log.Debug().Msgf("skipping %s %s because %s is missing from search query", ssd.Manufacturer, ssd.Name, word)
				break
			}
		}
		if !hasMissingWord {
			log.Debug().Msgf("adding %s %s to filtered list", ssd.Manufacturer, ssd.Name)
			filtered = append(filtered, ssd)
		}
	}
	return filtered
}

func cleanTitle(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`\[[^\]]+\]`).ReplaceAllString(s, "")
	stringsToRemove := []string{"ssd", "m2", "m.2", "nvme", "pcie", "gen", "amazon"}
	for _, toReplace := range stringsToRemove {
		s = strings.ReplaceAll(s, toReplace, "")
	}

	stringsToReplace := map[string]string{
		" wd":        " western digital",
		"team group": " teamgroup",
		"spatium":    " msi spatium",
		"sn850x":     " western digital sn850x",
	}
	for k, v := range stringsToReplace {
		if strings.Contains(s, k) {
			s = s + v
		}
	}
	return s
}

func doTest(esRepo *ssd.EsRepository) error {
	// Open the input CSV file for reading
	inputFile, err := os.Open("test/input.csv")
	if err != nil {
		return fmt.Errorf("opening input file: %w", err)
	}
	defer inputFile.Close()

	// Open a temporary output file for writing
	outputFile, err := os.Create("test/output.csv")
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer outputFile.Close()

	// Read all the records
	reader := csv.NewReader(inputFile)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading CSV records: %w", err)
	}

	writer := csv.NewWriter(outputFile)
	writer.Comma = '|'

	// Print each record
	for _, record := range records {
		ssds, err := esRepo.Search(context.Background(), cleanTitle(record[1]))
		if err != nil {
			log.Error().Err(err).Msgf("Error searching for SSD: %s", record[1])
			continue
		}
		var modelAndName string
		if len(ssds) != 0 {
			modelAndName = fmt.Sprintf("%s %s", ssds[0].Manufacturer, ssds[0].Name)
		}
		record = append(record, modelAndName)
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("writing record: %w", err)
		}
		writer.Flush()
	}
	return nil
}
