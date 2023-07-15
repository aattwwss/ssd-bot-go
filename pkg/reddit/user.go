package reddit

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type UserComment struct {
	SubredditID    string `json:"subreddit_id"`
	Subreddit      string `json:"subreddit"`
	ID             string `json:"id"`
	Author         string `json:"author"`
	ParentID       string `json:"parent_id"`
	AuthorFullname string `json:"author_fullname"`
	Body           string `json:"body"`
	LinkID         string `json:"link_id"`
	Name           string `json:"name"`
}

func (rc *Client) GetUserNewestComments(limit int) ([]UserComment, error) {
	redditUrl := fmt.Sprintf("https://oauth.reddit.com/user/%s/comments?limit=%v", rc.username, limit)
	req, err := rc.newRequest("GET", redditUrl, nil)
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

	var listing Listing[UserComment]
	err = json.NewDecoder(resp.Body).Decode(&listing)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	var userComments []UserComment
	for _, child := range listing.Data.Children {
		userComments = append(userComments, child.Data)
	}
	return userComments, nil
}
