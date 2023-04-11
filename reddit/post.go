package reddit

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type PostsDataResponse struct {
	Kind string `json:"kind"`
	Data Data   `json:"data"`
}

type Data struct {
	After     string `json:"after"`
	Dist      int    `json:"dist"`
	Modhash   any    `json:"modhash"`
	GeoFilter string `json:"geo_filter"`
	Posts     []Post `json:"children"`
	Before    any    `json:"before"`
}

type Post struct {
	Kind string   `json:"kind"`
	Data PostData `json:"data"`
}

type PostData struct {
	ID            string `json:"id"`
	Subreddit     string `json:"subreddit"`
	Title         string `json:"title"`
	Name          string `json:"name"`
	LinkFlairText string `json:"link_flair_text"`
}

func (rc RedditClient) GetNewPosts(subreddit string, limit int) ([]Post, error) {
	url := fmt.Sprintf("https://oauth.reddit.com/r/%s/new?limit=%v", subreddit, limit)
	req, err := rc.newRequest("GET", url, nil)
	if err != nil {
		log.Error().Msgf("Error creating request: %v", err)
		return nil, err
	}
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		log.Error().Msgf("Error sending request: %v", err)
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		log.Error().Msgf("Error request: %v", resp.Status)
		return nil, errors.New("Received non OK status code: " + resp.Status)
	}
	defer resp.Body.Close()

	var postDataResponse PostsDataResponse
	err = json.NewDecoder(resp.Body).Decode(&postDataResponse)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	return postDataResponse.Data.Posts, nil
}
