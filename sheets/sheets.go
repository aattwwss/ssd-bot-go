package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aattwwss/ssd-bot-go/ssd"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	SPREADSHEET_ID = "1B27_j9NDPU3cNlj2HKcrfpJKHkOf"
	SHEET_NAME     = "'Master List'" //take note of the single quote, which is needed for sheets with space in them
)

func main() {
	// Create a new Google Sheets API client with Application Default Credentials
	client, err := sheets.NewService(context.Background(), option.WithCredentialsFile("credential.json"))
	if err != nil {
		log.Fatalf("Failed to create Sheets API client: %v", err)
	}

	// Replace with your own spreadsheet ID
	spreadsheetID := "1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4"

	// Fetch all the data from the sheet
	response, err := client.Spreadsheets.Values.Get(spreadsheetID, "'Master List'").Do()
	if err != nil {
		log.Fatalf("Failed to fetch data from sheet: %v", err)
	}

	var allSSD []ssd.SSD
	for i, row := range response.Values {
		// skip the header
		if i == 0 {
			continue
		}
		// break at the end of the list of data
		if len(row) == 0 {
			break
		}

		ssd := ssd.SSD{
			Brand:         getFromSliceSafe(row, 0),
			Model:         getFromSliceSafe(row, 1),
			Interface:     getFromSliceSafe(row, 2),
			FormFactor:    getFromSliceSafe(row, 3),
			Capacity:      getFromSliceSafe(row, 4),
			Controller:    getFromSliceSafe(row, 5),
			Configuration: getFromSliceSafe(row, 6),
			DRAM:          getFromSliceSafe(row, 7),
			HMB:           getFromSliceSafe(row, 8),
			NandBrand:     getFromSliceSafe(row, 9),
			NandType:      getFromSliceSafe(row, 10),
			Layers:        getFromSliceSafe(row, 11),
			ReadWrite:     getFromSliceSafe(row, 12),
			Category:      getFromSliceSafe(row, 13),
			CellRow:       i + 1,
		}
		allSSD = append(allSSD, ssd)
	}
	fmt.Println(allSSD)
}

func getFromSliceSafe(row []interface{}, i int) string {
	if i >= len(row) {
		return ""
	}
	return row[i].(string)

}
