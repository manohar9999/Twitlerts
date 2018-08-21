package twittertalk

import (
	"net/http"
	"fmt"
	"net/url"
	"bytes"
	"log"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"strings"
)

const RateLimit = "https://api.twitter.com/1.1/application/rate_limit_status.json"
const UserTimeline = "https://api.twitter.com/1.1/statuses/user_timeline.json"

type TwitterConfig struct{
	ConsumerKey string
	ConsumerSecret string
}

type BearerToken struct {
	AccessToken string `json:"access_token"`
}

// Response object structure
type OAuth2Response struct {
	CreatedAt string `json:"created_at"`
	FullText string `json:"full_text"`
	Truncated bool `json:"truncated"`
	FavoriteCount int `json:"favorite_count"`
	RetweetCount int `json:"retweet_count"`
	Language string `json:"lang"`
	ID int64 `json:"id"`
	InReplyToStatusID string `json:"in_reply_to_status_id_str"`
	InReplyToScreenName string `json:"in_reply_to_screen_name"`
}

// Ref: https://developer.twitter.com/en/docs/basics/authentication/overview/application-only.html
func OAuth2Authenticate(config TwitterConfig) BearerToken {
	client := &http.Client{}
	// Encode consumer key and secret
	encodedKeySecret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		url.QueryEscape(config.ConsumerKey),
		url.QueryEscape(config.ConsumerSecret))))

	// Obtain a bearer token
	// The body of the request must be grant_type=client_credentials
	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	// The request must be a HTTP POST request
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	if err != nil {
		log.Fatal(err, client, req)
	}
	// The request must include an Authorization header formatted as
	// Basic <base64 encoded value from step 1>.
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedKeySecret))
	// The request must include a Content-Type header with
	// the value of application/x-www-form-urlencoded;charset=UTF-8.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))

	//Issue the request and get the bearer token from the JSON you get back
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, respBody)
	}

	var b BearerToken
	json.Unmarshal(respBody, &b)

	return b
}

// Returns a list of OAuth2Response response objects
func GetTweets(b BearerToken, screenname string, tweetcount int, includeentities bool) []OAuth2Response {
	client := &http.Client{}

	// Choose your API endpoint that supports application only auth context
	// and create a request object with that
	// Ref: https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline.html
	twitterEndPoint := fmt.Sprintf("%s?screen_name=%s&count=%d&include_entities=%t&tweet_mode=extended",UserTimeline,screenname,tweetcount,includeentities)
	req, err := http.NewRequest("GET", twitterEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Authenticate API requests with the bearer token
	// include an Authorization header formatted as
	// Bearer
	req.Header.Add("Authorization",
		fmt.Sprintf("Bearer %s", b.AccessToken))

	// Issue the request and get the JSON API response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	returnCode, _ := strconv.Atoi(strings.Split(resp.Status, " ")[0])
	if returnCode == http.StatusOK {
		oauth2resp := make([]OAuth2Response,0)
		json.Unmarshal(respBody, &oauth2resp)
		return oauth2resp
	}

	return nil
}