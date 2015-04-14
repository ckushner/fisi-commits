package commits

import (
	"fmt"
	"time"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

type GithubCommitService struct {
	ScrapeUser chan *string	
	UseCommit  func(*string, *string)
	LastCycle *time.Time
}

type GithubUserService struct {
	CheckUser chan *string
	RegisterUser chan *github.User
}

func (service *GithubCommitService) Start() {
	
	client := client()

	getUserCommits := func(user *string) {
		// Get repos for user
		r_opt := &github.RepositoryListOptions{Type: "all", Sort: "updated"}
		repos, _, err := client.Repositories.List(*user, r_opt)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		// Get commits per repo
		lastRequest := service.LastCycle.Add(-200*time.Hour)
		c_opt := &github.CommitsListOptions{Since: lastRequest}

		// TODO: cache repos checked and don't recheck in same cycle
		getRepoCommits := func(repo github.Repository) {
			commits, _, err := client.Repositories.ListCommits(*repo.Owner.Login, *repo.Name, c_opt)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, commit := range commits {
				 go service.UseCommit(user, commit.Commit.Message)
			}
		}

		for _, repo := range repos {
			if repo.UpdatedAt.Time.After(lastRequest) {
				go getRepoCommits(repo)
				fmt.Println(*repo.Name)
			}
		}
	}

	// Service requests when needed by scraper
	for {
		user := <- service.ScrapeUser
		go getUserCommits(user)
	}
}

func (service *GithubUserService) Start() {
	
	client := client()

	checkUserName := func(username string) {
		user, _, err := client.Users.Get(username)
		if err != nil {
			fmt.Println(err)
			service.RegisterUser <- nil
		}
		// if user was not valid github user -> err == 1
		// only gets valid github users pass
		service.RegisterUser <- user
	}

	// Service requests when needed by add user endpoint
	for {
		user := <- service.CheckUser
		go checkUserName(*user)
	}
}

type tokenSource struct {
  token *oauth2.Token
}

// add Token() method to satisfy oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error){
  return t.token, nil
}

func client() *github.Client {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: "AccessToken"},
	}
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}
