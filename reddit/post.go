package reddit

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Listing[T any] struct {
	Kind string         `json:"kind"`
	Data ListingData[T] `json:"data"`
}

type ListingData[T any] struct {
	After     string        `json:"after"`
	Dist      int           `json:"dist"`
	Modhash   string        `json:"modhash"`
	GeoFilter string        `json:"geo_filter"`
	Children  []Children[T] `json:"children"`
	Before    string        `json:"before"`
}

type Children[T any] struct {
	Kind string `json:"kind"`
	Data T      `json:"data"`
}

type Post struct {
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

	var postDataListing Listing[Post]
	err = json.NewDecoder(resp.Body).Decode(&postDataListing)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	var posts []Post
	for _, child := range postDataListing.Data.Children {
		posts = append(posts, child.Data)
	}
	return posts, nil
}

type PostComment struct {
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

func (rc RedditClient) GetCommentsByPostId(postId string, limit int) ([]PostComment, error) {

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

	var postCommentsListing []Listing[PostComment]
	err = json.NewDecoder(resp.Body).Decode(&postCommentsListing)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	var postComments []PostComment
	for _, child := range postCommentsListing[1].Data.Children {
		postComments = append(postComments, child.Data)
	}
	return postComments, nil
}
