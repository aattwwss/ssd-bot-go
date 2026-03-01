package reddit

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// UserComment represents a comment made by a user.
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

// GetUserNewestComments fetches the newest comments by the authenticated user.
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
		return nil, fmt.Errorf("received non OK status code: %s", resp.Status)
	}
	defer resp.Body.Close()

	var listing Listing[UserComment]
	err = json.NewDecoder(resp.Body).Decode(&listing)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding response body")
		return nil, err
	}
	var userComments []UserComment
	for _, child := range listing.Data.Children {
		userComments = append(userComments, child.Data)
	}
	return userComments, nil
}
