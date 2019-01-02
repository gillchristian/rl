package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gillchristian/rl"
)

const pinboardAPI = "https://api.pinboard.in/v1"

type item struct {
	Href string `json:"href"`
}

// TODO (refactor): this initialization is shared in both cmd/ packages
var fileName string

func init() {
	usr, _ := user.Current()
	dir := usr.HomeDir

	fileName = filepath.Join(dir, ".reading-list")
}

func main() {
	user := os.Args[1]
	token := os.Args[2]

	n, err := fetch(user, token)
	checkErr(err)

	err = write(n)
	checkErr(err)
}

func write(n rl.ReadingList) error {
	existing, err := rl.Read(fileName)

	if err != nil {
		return err
	}

	// TODO: use sync function
	existing.Items = append(existing.Items, n.Items...)
	existing.Added += len(n.Items)

	err = rl.Write(fileName, existing)

	if err != nil {
		return err
	}

	return nil
}

func fetch(user, token string) (rl.ReadingList, error) {
	url := pinboardAPI +
		fmt.Sprintf("/posts/all?auth_token=%s:%s&format=json", user, token)

	r, err := http.Get(url)
	if err != nil {
		return rl.ReadingList{}, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return rl.ReadingList{}, err
	}

	items := []item{}

	err = json.Unmarshal(b, &items)

	if err != nil {
		return rl.ReadingList{}, err
	}

	readingList := rl.ReadingList{
		Reads: 0,
		Added: 0,
		Items: make([]string, len(items)),
	}

	for i, item := range items {
		readingList.Items[i] = item.Href
	}

	return readingList, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
