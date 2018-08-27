package main

import (
	"fmt"
	"github.com/manohar9999/Twitlerts/twittertalk"
)

func main() {
	bearerToken := twittertalk.Oauth2Setup("config.json")
	var channel = make(chan twittertalk.OAuth2Response)
	go twittertalk.GetAllTweets("tim_cook", channel, bearerToken)
	for i := range channel {
		// for now, just printing them. I know.. I know..
		fmt.Println(i.FullText)
	}
}