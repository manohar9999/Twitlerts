package twittertalk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Used to retrieve all the tweets of the user.
// since_id and max_id are used to keep track
// of the cursor
// param in:
//	screenname string - Twitter user screen name
//  ch chan OAuth2Response - channel of OAuth2Response objects
// param out:
//	None
func GetAllTweets(screenname string, ch chan OAuth2Response, token BearerToken) {
	defer close(ch)

	// maximum number of tweets you can retrieve at a time are 200.
	tweets := GetTweets(token, screenname, 200, true, 0,0)
	// firstid is later used to extract all the
	// new tweets using since_id once the older
	// tweets are extracted
	var firstid int64
	// lastid if constantly used to retrieve
	// the next 200 tweets until all the tweets
	// are retrieved.
	var lastid int64
	if len(tweets) > 0 {
		// Save it in the first pass
		firstid = tweets[0].ID
		lastid = tweets[len(tweets)-1].ID
	}

	// Process them
	ProcessTweets(tweets, ch)

	var firstpass bool = true
	for oldertweets := GetTweets(token, screenname, 200, true, 0,lastid); oldertweets != nil ;  {
		if !firstpass {
			oldertweets = GetTweets(token, screenname, 200, true, 0,lastid)
		}
		firstpass = false
		lastid = oldertweets[len(oldertweets)-1].ID
		// length of oldertweets has to be greater than 1
		// since it's inclusive of tweet of lastid
		if len(oldertweets) > 1 {
			// process older tweets
			ProcessTweets(oldertweets, ch)
			//for _, tweet := range oldertweets {
			//	ch <- tweet
			//}
		} else {
			// no older tweets
			// break out of for loop
			break
		}
	}
	newertweets := GetTweets(token, screenname, 200, true, firstid,0)
	if len(newertweets) > 0 {
		// process newer tweets
		ProcessTweets(tweets, ch)
	}
}

// pushing the response objects onto the channel
func ProcessTweets(tweets []OAuth2Response, ch chan OAuth2Response)  {
	for _, tweet := range tweets {
		ch <- tweet
	}
}

// param in:
//	btoken BearerToken - This is required to authenticate
//	screenname string - Twitter user screen name
//	tweetcount int - Number of tweets to retrieve
//  includeentities bool - includes hashtags, user mentions, links,
// 		stock tickers (symbols), Twitter polls, and attached media
//		This is helpful preprocessed information
//	sinceid int64 - Returns newer results with an ID greater than sinceid
//		This is ignored if it is 0
//	maxid int64 - Returns older results less than or equal to this maxid
//		This is ignored if it is 0
// param out:
// 	Returns a list of OAuth2Response response objects
func GetTweets(token BearerToken, screenname string, tweetcount int, includeentities bool, sinceid int64, maxid int64) []OAuth2Response {
	client := &http.Client{}

	// Choose your API endpoint that supports application only auth context
	// and create a request object with that
	// Ref: https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline.html
	var twitterEndPoint string
	// Do not use since_id or max_id if they are 0
	if sinceid != 0 {
		twitterEndPoint = fmt.Sprintf("%s?screen_name=%s&count=%d&include_entities=%t&tweet_mode=extended&include_rts=true&since_id=%d", UserTimeline, screenname, tweetcount, includeentities, sinceid)
	} else if maxid!= 0 {
		twitterEndPoint = fmt.Sprintf("%s?screen_name=%s&count=%d&include_entities=%t&tweet_mode=extended&include_rts=true&max_id=%d", UserTimeline, screenname, tweetcount, includeentities, maxid)
	} else {
		twitterEndPoint = fmt.Sprintf("%s?screen_name=%s&count=%d&include_entities=%t&tweet_mode=extended&include_rts=true", UserTimeline, screenname, tweetcount, includeentities)

	}
	req, err := http.NewRequest("GET", twitterEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Authenticate API requests with the bearer token
	// include an Authorization header formatted as
	// Bearer
	req.Header.Add("Authorization",
		fmt.Sprintf("Bearer %s", token.AccessToken))

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
