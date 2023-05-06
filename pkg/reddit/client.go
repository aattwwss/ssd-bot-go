package reddit

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type tokenRes struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type RedditClient struct {
	httpClient     *http.Client
	clientId       string
	clientSecret   string
	username       string
	password       string
	overrideOldBot bool

	accessToken          string
	tokenExpireTimeMilli int64

	isDebug bool
}

const (
	userAgent = "SSD bot v2.0 by /u/_SSD_BOT_ github.com/aattwwss/ssd-bot-go" //need to set user agent to prevent getting blocked by reddit
)

func NewRedditClient(clientId, clientSecret, username, password, accessToken string, expireTimeMilli int64, overrideOldBot, isDebug bool) (*RedditClient, error) {
	if clientId == "" || clientSecret == "" || username == "" || password == "" {
		return nil, errors.New("clientId, clientSecret, username, password cannot be empty")
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	rc := RedditClient{
		httpClient:           httpClient,
		clientId:             clientId,
		clientSecret:         clientSecret,
		username:             username,
		password:             password,
		overrideOldBot:       overrideOldBot,
		accessToken:          accessToken,
		tokenExpireTimeMilli: expireTimeMilli,
		isDebug:              isDebug,
	}

	err := rc.RefreshToken()
	if err != nil {
		return nil, err
	}

	return &rc, nil
}

func (rc *RedditClient) RefreshToken() error {
	now := time.Now()
	durationFromExpire := time.UnixMilli(rc.tokenExpireTimeMilli).Sub(now).Minutes()
	if durationFromExpire > 30 {
		log.Info().Msgf("Token is still valid for %v minutes. Refresh only with 30 minutes left.", int(durationFromExpire))
		return nil
	}
	// Set the form data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", rc.username)
	data.Set("password", rc.password)

	// Create a new POST request
	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	req.Header.Add("User-Agent", userAgent)
	if err != nil {
		log.Error().Msgf("Error creating request: %v", err)
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.SetBasicAuth(rc.clientId, rc.clientSecret)

	// Send the request
	resp, err := retryHttpRequest(rc.httpClient, req, 5, time.Minute)
	if err != nil {
		log.Error().Msgf("Error sending request: %v", err)
		return err
	}
	if resp.StatusCode/100 != 2 {
		log.Error().Msgf("Error request: %v", resp.Status)
		return errors.New("Received non OK status code: " + resp.Status)
	}

	defer resp.Body.Close()

	var tokenRes tokenRes
	err = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if err != nil {
		log.Error().Msgf("Error decoding response body:", err)
		return err
	}

	rc.accessToken = tokenRes.AccessToken
	rc.tokenExpireTimeMilli = now.Add(time.Duration(tokenRes.ExpiresIn) * time.Second).UnixMilli()

	if rc.isDebug {
		log.Info().Msgf("token res: %v", tokenRes)
		log.Info().Msgf("access token: %v", rc.accessToken)
		log.Info().Msgf("expire time milli: %v", rc.tokenExpireTimeMilli)
	}
	return nil
}

func (rc *RedditClient) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer "+rc.accessToken)
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

func retryHttpRequest(client *http.Client, req *http.Request, attempts int, sleep time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < attempts; i++ {
		resp, err = client.Do(req)
		if err == nil || resp.StatusCode/100 == 2 {
			return resp, nil
		}

		log.Error().Msgf("Error sending request: %v", err)
		time.Sleep(sleep)
		sleep *= 2 // increase delay exponentially
	}

	return nil, errors.New("http request exceeded retry attempts")
}
