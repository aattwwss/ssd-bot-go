package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/aattwwss/ssd-bot-go/pkg/reddit"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
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

	_, err = reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.OverrideOldBot, config.IsDebug)
	if err != nil {
		log.Fatal().Msgf("Init reddit client error: %v", err)
	}
	es, _ := elasticutil.NewElasticsearchClient(config.EsAddress)
	esRepo := ssd.NewEsSSDRepository(es, "ssd-index")
	doTest(esRepo)
	// tpuRepo := ssd.NewTpuSSDRepository(config.TPUHost, config.TPUUsername, config.TPUSecret)
	// ssdSync := ssd.SSDSync{
	// 	StartId:  1,
	// 	Delay:    time.Duration(100),
	// 	IdToSkip: []int{},
	// }
	// err = ssdSync.Sync(context.Background(), tpuRepo, esRepo)
	// if err != nil {
	// 	log.Fatal().Msgf("sync error", err)
	// }
	// ssd, _ := esRepo.FindById(context.Background(), "123")
	// ssds, _ := esRepo.SearchBasic(context.Background(), "corsair")
	// title := "[SSD] Sabrent Rocket 2230 NVMe 4.0 1TB - $102.99 ($109.99-$7 with code SSCSA536)"
	// sss, _ := esRepo.Search(context.Background(), cleanTitle(title))
	// log.Info().Msgf("%v", ssd)
	// log.Info().Msgf("%v", ssds)
	// log.Info().Msgf("cleaned title: %v", cleanTitle(title))
	// log.Info().Msgf("%v", sss[0])
	// ssd.DriveID = ssd.DriveID + "_new"
	// ssd.Capacity = "some capacity"
	// esRepo.Insert(context.Background(), *ssd)
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
