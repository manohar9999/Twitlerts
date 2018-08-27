package twittertalk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type TwitterConfig struct{
	ConsumerKey string
	ConsumerSecret string
}

const RateLimit = "https://api.twitter.com/1.1/application/rate_limit_status.json"
const UserTimeline = "https://api.twitter.com/1.1/statuses/user_timeline.json"

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
// param in:
//	config TwitterConfig - config object which contains key and secret
// param out:
//	BearerToken - returns the token required to authenticate the requests
func oauth2Authenticate(config TwitterConfig) BearerToken {
	client := &http.Client{}
	// encode consumer key and secret
	encodedKeySecret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		url.QueryEscape(config.ConsumerKey),
		url.QueryEscape(config.ConsumerSecret))))

	// obtain a bearer token
	// the body of the request must be grant_type=client_credentials
	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	// The request must be a HTTP POST request
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	if err != nil {
		log.Fatal(err, client, req)
	}
	// the request must include an Authorization header formatted as
	// basic <base64 encoded value from step 1>.
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedKeySecret))
	// The request must include a Content-Type header with
	// the value of application/x-www-form-urlencoded;charset=UTF-8.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))

	// issue the request and get the bearer token from the JSON you get back
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

// Exposed method that handles the auth flow
//	param in: config file name which contains key and secret
//		in json format
//	param out: returns the token required to authenticate the requests
func Oauth2Setup(filename string) BearerToken {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	twitterconfig := TwitterConfig{}
	// reads the json and stores it in TwitterConfig.
	decoder.Decode(&twitterconfig)
	bearerToken := oauth2Authenticate(TwitterConfig{twitterconfig.ConsumerKey, twitterconfig.ConsumerSecret})

	return bearerToken
}
