package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

type Submission struct {
	ID            string `json:"id"`
	Subreddit     string `json:"subreddit"`
	Title         string `json:"title"`
	Name          string `json:"name"`
	LinkFlairText string `json:"link_flair_text"`
}

func (rc *RedditClient) GetNewSubmissions(subreddit string, limit int) ([]Submission, error) {
	redditUrl := fmt.Sprintf("https://oauth.reddit.com/r/%s/new?limit=%v", subreddit, limit)
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

	var listings Listing[Submission]
	err = json.NewDecoder(resp.Body).Decode(&listings)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	var posts []Submission
	for _, child := range listings.Data.Children {
		posts = append(posts, child.Data)
	}
	return posts, nil
}

type SubmissionComment struct {
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

func (rc *RedditClient) GetCommentsBySubmissionId(submissionId string, limit int) ([]SubmissionComment, error) {
	redditUrl := fmt.Sprintf("https://oauth.reddit.com/comments/%s?limit=%v&depth=1", submissionId, limit)
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

	var listings []Listing[SubmissionComment]
	err = json.NewDecoder(resp.Body).Decode(&listings)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return nil, err
	}
	var submissionComments []SubmissionComment
	// hardcoding to use the 2nd element as the first is the post information
	for _, child := range listings[1].Data.Children {
		submissionComments = append(submissionComments, child.Data)
	}
	return submissionComments, nil
}

func (rc *RedditClient) SubmitComment(postId, text string) error {
	data := url.Values{}
	data.Set("api_type", "json")
	data.Set("text", text)
	data.Set("thing_id", "t3_"+postId)

	req, err := rc.newRequest("POST", "https://oauth.reddit.com/api/comment", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Error().Msgf("Error creating request: %v", err)
		return err
	}
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		log.Error().Msgf("Error sending request: %v", err)
		return err
	}
	if resp.StatusCode/100 != 2 {
		log.Error().Msgf("Error request: %v", resp.Status)
		return errors.New("Received non OK status code: " + resp.Status)
	}
	defer resp.Body.Close()
	return nil
}

func (rc *RedditClient) IsCommentedByUser(submissionId string, author string) bool {
	comments, _ := rc.GetCommentsBySubmissionId(submissionId, 100)
	for _, comment := range comments {
		if comment.Author == author {
			return true
		}
	}
	return false
}
