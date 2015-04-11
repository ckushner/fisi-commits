package commits

import (
	"fmt"
	"time"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

// create struct for the token source
type tokenSource struct {
  token *oauth2.Token
}

// add Token() method to satisfy oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error){
  return t.token, nil
}

func GetUserCommits(user string, useCommit func(*string)) {
	// Start client
	ts := &tokenSource{
		&oauth2.Token{AccessToken: "a1235512708e10706259ae110eece2cb73448d3b"},
	}
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	
	// Get repos for user
	r_opt := &github.RepositoryListOptions{Type: "all", Sort: "updated"}
	repos, _, err := client.Repositories.List(user, r_opt)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Get commits per repo
	baseTime := time.Now()
	lastRequest := baseTime.Add(-200*time.Hour)
	c_opt := &github.CommitsListOptions{Since: lastRequest}

	getRepoCommits := func(repo github.Repository) {
		// fmt.Println(*repo.Name)
		commits, _, err := client.Repositories.ListCommits(*repo.Owner.Login, *repo.Name, c_opt)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, commit := range commits {
			 go useCommit(commit.Commit.Message)
		}
	}
	for _, repo := range repos {
		if repo.UpdatedAt.Time.After(lastRequest) {
			go getRepoCommits(repo)
			fmt.Println(*repo.UpdatedAt)
		}
	}
}
