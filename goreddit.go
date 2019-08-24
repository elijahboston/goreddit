package main

import (
	json "encoding/json"
	"flag"
	"fmt"
	ioutil "io/ioutil"
	http "net/http"
)

type APIResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Id        string `json:"id"`
				Subreddit string `json:"subreddit"`
				Title     string `json:"title"`
				Ups       int    `json:"ups"`
				Downs     int    `json:"downs"`
				Score     int    `json:"score"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func getThreads(body []byte) (*APIResponse, error) {
	var t = new(APIResponse)
	err := json.Unmarshal(body, &t)

	if err != nil {
		panic(err.Error())
	}

	return t, err
}

func main() {
	var subreddit string
	flag.StringVar(&subreddit, "s", "all", "name of the subreddit")

	flag.Parse()

	// Format URL string for API request
	url := fmt.Sprintf("https://www.reddit.com/r/%s/.json", subreddit)

	resp, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()

	// Parse the response body from JSON
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
	}

	if resp.StatusCode != 200 {
		err := fmt.Sprintf("Bad Status: %s", resp.Status)
		panic(err)
	}

	t, err := getThreads([]byte(body))

	fmt.Printf("/%s - %d threads\n", subreddit, len(t.Data.Children))

	for _, c := range t.Data.Children {
		fmt.Printf("%s [%d/%d]\n\n", c.Data.Title, c.Data.Ups, c.Data.Downs)
	}
}
