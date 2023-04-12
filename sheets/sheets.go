package sheets

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GetSheetsValues(spreadSheetId, range_ string) ([][]interface{}, error) {
	// Hardcoding the credentials for now since we won't be calling the api frequently
	client, err := sheets.NewService(context.Background(), option.WithCredentialsFile("credential.json"))
	if err != nil {
		log.Error().Msgf("Failed to create Sheets API client: %v", err)
		return nil, err
	}

	// Fetch all the data from the sheet
	response, err := client.Spreadsheets.Values.Get(spreadSheetId, range_).Do()
	if err != nil {
		log.Error().Msgf("Failed to fetch data from sheet: %v", err)
		return nil, err
	}

	return response.Values, nil
}
