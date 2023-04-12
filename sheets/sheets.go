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
		allSSD = append(allSSD, ssd)
	}
	fmt.Println(allSSD)
}

func getStringAtIndexOrEmpty(arr []interface{}, i int) string {
	if i >= len(arr) {
		return ""
	}
	return fmt.Sprintf("%v", arr[i])

}
