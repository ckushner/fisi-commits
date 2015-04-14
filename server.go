package main

import (
	"fmt"
	"net/http"
	"log"
	"regexp"
	"time"
	"github.com/google/go-github/github"
	"github.com/ckushner/fisi-commits/commits"	
)

func main() {

	/* Start tweetbot */

	tweet := make(chan *string)
	go commits.TweetService(tweet)

	/* Start github scraper */

	scrape := make(chan *string)
	score := make(chan *string)
	cycleTime := time.Now()
	commit_service := commits.GithubCommitService {
		ScrapeUser: scrape,
		UseCommit: func (user *string, commit *string) {
			fisi_regex := ".*(fuck\\s?it\\s?ship\\s?it)|(fisi).*"
			matched, err := regexp.MatchString(fisi_regex, *commit)
			if err != nil {
				fmt.Println(err)
				return
			}
			if matched {
				tweet <- commit
				score <- user
			}
		},
		LastCycle: &cycleTime,
	}
	go commit_service.Start()

	/* Send stored usernames to scraper service */
	
	ubuffhttp := make(chan *string)
	ubuffserv := make(chan *string)	
	ucheck := make(chan *string)	
	uregister := make(chan *github.User)
	
	var um commits.UserMap
	um = make(commits.UserMap)

	// register users before/after scrape so as not to interfere
	var serviceAddUser func()
	serviceAddUser = func() {
		select {
		case username := <-ubuffserv:
			if _, ok := um[*username]; ok == false {
				ucheck <- username
				if user := <-uregister; user != nil {
					um.RegisterUsers(user)
					// TODO: look at past 24 hours of user commits
				}
			}
			serviceAddUser()
		default:
			return
		}
	}

	var updateScore func()
	updateScore = func() {
		select {
		case username := <-score:
			if _, ok := um[*username]; ok == true {
				um[*username].Stats.NumCommitsCycle++
				fmt.Println("User: ", *username, " score: ", um[*username].Stats.NumCommitsCycle)
			}
			updateScore()
		default:
			return
		}
	}

	go func () {
		for {
			serviceAddUser()
			updateScore()
			// TODO: put in timed loop, update cycleTime
			for username, _ := range um {
			    scrape <- &username
			}
			time.Sleep(time.Duration(10)*time.Second)
		}
	}()

	// take off of ubuffhttp immediately so resp served
	go func () {
		for {
			buff := <-ubuffhttp
			ubuffserv <- buff
		}
	}()

	/* Start http server and github client for user tracking */
	
	user_service := commits.GithubUserService {
		CheckUser: ucheck,
		RegisterUser: uregister,
	}
	go user_service.Start()

	http.HandleFunc("/add", func (w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if len(query["user"]) == 0 || query["user"][0] == "" {
			fmt.Fprintf(w, "Error: did not send github username.\nFormat: /add?user=MYNAME")
			return;
		}
		user := query["user"][0]
		fmt.Println(user)
		if _, err := fmt.Fprintf(w, "Thanks, %s!\n", user); err != nil {
			fmt.Println(err)
			return
		}
		ubuffhttp <- &user
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
