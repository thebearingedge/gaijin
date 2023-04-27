package main

import (
	"log"
	"net/http"
	"os"
	"repository/json"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

var userURL = "https://jsonplaceholder.typicode.com/users"
var postsURL = "https://jsonplaceholder.typicode.com/posts"

func main() {
	stdout := log.New(os.Stdout, "info - ", 0)
	stderr := log.New(os.Stderr, "error - ", 0)

	usersData := json.NewFetcher[User](*http.DefaultClient, userURL)

	if user, err := usersData.FetchById("1"); err != nil {
		stderr.Printf("user error is: %v\n", err)
	} else {
		stdout.Printf("user data is: %#v\n", user)
	}

	postsData := json.NewFetcher[Post](*http.DefaultClient, postsURL)

	if posts, err := postsData.FetchWhere("userId=1"); err != nil {
		stderr.Printf("posts error is %v\n", err)
	} else {
		stdout.Printf("posts data is: %#v\n", posts)
	}

}
