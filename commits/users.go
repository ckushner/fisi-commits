package commits

import "github.com/google/go-github/github"

type User struct {
	Email *string
	Stats *UserStats
}

type UserStats struct {
	NumCommitsAll int
	NumCommitsCycle int
}

type UserMap map[string]*User // github username to User struct

func (um *UserMap) RegisterUsers(user *github.User) {
	stats := UserStats{NumCommitsAll:0,NumCommitsCycle:0}
	new_user := User{Email:user.Email,Stats:&stats}
	(*um)[*user.Login] = &new_user
}
