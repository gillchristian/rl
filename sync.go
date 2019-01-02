package rl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TODO: use as defult but receive from flags
const gistFileName = "reading-list"
const description = "reading list (https://github.com/gillchristian/rl)"
const gitHubGistAPI = "https://api.github.com/gists"

// GistFile is the gist file content.
type GistFile struct {
	// TODO: Content should be of type ReadingList to Un/Marshalled together
	Content string `json:"content"`
}

// GithubGist represents the response/payload for getting or updating a GitHub
// gist, respectively. For mor information check:
// https://developer.github.com/v3/gists/#create-a-gist
// https://developer.github.com/v3/gists/#get-a-single-gist
type GithubGist struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

var client = &http.Client{}

// TODO: merge SyncWithGist & CreateGist

// SyncWithGist syncs the local reading list with one in a GitHub Gist.
func SyncWithGist(file, token, gistID string) error {
	remote, err := fetch(gistID, token)
	if err != nil {
		return err
	}

	local, err := Read(file)
	if err != nil {
		return err
	}

	merged := merge(remote.Items, local.Items)
	delta := len(merged) - len(remote.Items)

	remote.Items = merged
	remote.Added += delta

	return update(gistID, token, remote)
}

// CreateGist creates a GitHub Gist and saves the local reading list to it.
func CreateGist(file, token string) (string, error) {
	local, err := Read(file)
	if err != nil {
		return "", err
	}

	return create(file, local)
}

// TODO: duplicated code in fetch, create & update

func fetch(gistID, token string) (ReadingList, error) {
	// TODO: if gistID is empty create gist
	req, err := http.NewRequest(http.MethodGet, gistURL(gistID), nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return ReadingList{}, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ReadingList{}, err
	}

	var remote GithubGist
	err = json.Unmarshal(b, &remote)
	if err != nil {
		return ReadingList{}, err
	}

	s, ok := remote.Files[gistFileName]
	var remoteContent ReadingList
	if !ok {
		remoteContent = ReadingList{0, 0, []string{}}
	}

	err = json.Unmarshal([]byte(s.Content), &remoteContent)
	if err != nil {
		remoteContent = ReadingList{0, 0, []string{}}
	}

	return remoteContent, nil
}

func create(token string, content ReadingList) (string, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return "", err
	}

	files := map[string]GistFile{gistFileName: GistFile{string(b)}}
	gistFile := GithubGist{Description: description, Files: files}

	b, err = json.Marshal(gistFile)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, gitHubGistAPI, bytes.NewReader(b))
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var newGist GithubGist
	err = json.Unmarshal(b, &newGist)
	if err != nil {
		return "", err
	}

	return newGist.ID, nil
}

func update(gistID, token string, content ReadingList) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	files := map[string]GistFile{gistFileName: GistFile{string(b)}}
	gistFile := GithubGist{Description: description, Files: files}

	b, err = json.Marshal(gistFile)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, gistURL(gistID), bytes.NewReader(b))
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	return nil
}

func gistURL(gistID string) string {
	return fmt.Sprintf("%s/%s", gitHubGistAPI, gistID)
}

// merge merges slice b into slice a.
// All duplicated values from b are excluded.
func merge(a, b []string) []string {
	set := make(map[string]bool)

	for _, item := range a {
		set[item] = true
	}

	result := a

	for _, item := range b {
		if !set[item] {
			result = append(result, item)
		}
	}

	return result
}
