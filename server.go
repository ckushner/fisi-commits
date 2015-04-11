package main

import (
	"fmt"
	"regexp"
	"github.com/ckushner/fisi-commits/commits"	
)

func main() {
	// TODO: start endpoint for accepted usernames
	
	sendCommits := func(commit *string) {
		fisi_regex := ".*(fuck\\s?it\\s?ship\\s?it)|(fisi).*"
		matched, err := regexp.MatchString(fisi_regex, *commit)
		if err != nil {
			fmt.Println(err)
			return
		}
		if matched {
			commits.Tweet(*commit)
		}
	}

	// TODO: put in timed loop, and time as param
	// TODO: cache repos checked and don't recheck in same cycle
	commits.GetUserCommits("ckushner", sendCommits)
}
