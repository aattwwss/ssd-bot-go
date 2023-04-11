package reddit

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

// the actual payload returned from reddit api
type PostsDataResponse struct {
	Kind string                `json:"kind"`
	Data PostsDataResponseData `json:"data"`
}

type PostsDataResponseData struct {
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

// the actual payload returned from reddit api
type PostCommentsResponse []struct {
	Kind string                   `json:"kind"`
	Data PostCommentsResponseData `json:"data"`
}

type PostCommentsResponseData struct {
	After        any            `json:"after"`
	Dist         any            `json:"dist"`
	Modhash      any            `json:"modhash"`
	GeoFilter    string         `json:"geo_filter"`
	PostComments []PostComments `json:"children"`
	Before       any            `json:"before"`
}

type PostComments struct {
	Kind string           `json:"kind"`
	Data PostCommentsData `json:"data"`
}

type PostCommentsData struct {
	SubredditID    string `json:"subreddit_id"`
	Subreddit      string `json:"subreddit"`
	ID             string `json:"id"`
	Author         string `json:"author"`
	ParentID       string `json:"parent_id"`
	AuthorFullname string `json:"author_fullname"`
	Body           string `json:"body"`
	Name           string `json:"name"`
	IsSubmitter    bool   `json:"is_submitter"`
}

func (rc RedditClient) GetCommentsByPostId(postId string, limit int) ([]PostComments, error) {

	url := fmt.Sprintf("https://oauth.reddit.com/comments/%s?limit=%v&depth=1", postId, limit)
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

	var postCommentsResponse PostCommentsResponse
	err = json.NewDecoder(resp.Body).Decode(&postCommentsResponse)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	return postCommentsResponse[1].Data.PostComments, nil
}
