package commits

import (
	"fmt"
	// "github.com/ChimeraCoder/anaconda"
)

func TweetService (tweet chan *string) {
	// anaconda.SetConsumerKey("SetConsumerKey")
	// anaconda.SetConsumerSecret("SetConsumerSecret")
	// api := anaconda.NewTwitterApi(
	// 	"access-token",
	// 	"access-token-secret",
	// )
	
	// // Tweet commit messages as they are detected
	for {
		commit := <- tweet
		fmt.Println(*commit)
	// 	_, err := api.PostTweet(*commit, nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	}
}

// TODO: tweet weekly leaderboard (thursday nights)

// TODO: get tweets at account to register usernames
//       GET statuses/mentions_timeline with ""
