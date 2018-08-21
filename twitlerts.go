package main

import (
	"github.com/manohar9999/Twitlerts/twittertalk"
	"os"
	"encoding/json"
	"log"
	"fmt"
)

func main() {

	file, err  := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	twitterconfig := twittertalk.TwitterConfig{}
	// Reads the json and stores it in TwitterConfig.
	decoder.Decode(&twitterconfig)

	bearerToken := twittertalk.OAuth2Authenticate(twittertalk.TwitterConfig{twitterconfig.ConsumerKey,twitterconfig.ConsumerSecret})
	tweets := twittertalk.GetTweets(bearerToken,"tim_cook", 10,true)
	if tweets != nil {
		for _, tweet := range tweets {
			fmt.Println(tweet.ID, tweet.FullText)
		}
	}
}