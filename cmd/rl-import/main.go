package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type item struct {
	Href string `json:"href"`
}

func main() {
	user := os.Args[1]
	token := os.Args[2]
	url := "https://api.pinboard.in/v1/posts/all?auth_token=" + user + ":" + token + "&format=json"

	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	items := []item{}

	err = json.Unmarshal(b, &items)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	links := make([]string, len(items))

	for i, item := range items {
		links[i] = item.Href
	}

	s := strings.Join(links, "\n")

	fmt.Println(s)
}
