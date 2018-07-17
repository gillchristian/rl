// Package rl provides methods to handle reading list files.
package rl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	xdgOpen "github.com/skratchdot/open-golang/open"
)

// ReadingList represents the content of a reading list file.
type ReadingList struct {
	Reads int      `json:"reads"`
	Added int      `json:"added"`
	Items []string `json:"items"`
}

// TODO: marshal & write function

// Add adds an item i at the end of the file.
func Add(file, i string) error {
	rl, err := Read(file)

	if err != nil {
		return err
	}

	rl.Items = append(rl.Items, i)
	rl.Added++

	return Write(file, rl)
}

// Done removes the first line in file and increases the count.
func Done(file string) error {
	rl, err := Read(file)

	if err != nil {
		return err
	}

	if len(rl.Items) > 0 {
		rl.Reads++
		rl.Items = rl.Items[1:]
	}

	return Write(file, rl)
}

// Remove removes the first line in file.
func Remove(file string) error {
	rl, err := Read(file)

	if err != nil {
		return err
	}

	rl.Items = rl.Items[1:]

	return Write(file, rl)
}

// Open opens with xdg-open (or the Mac/Windows equivalent) the first item in file.
func Open(file string) error {
	i, err := first(file)

	if i != "" {
		fmt.Println("Opening " + i)
		return xdgOpen.Start(i)
	}

	return err
}

// Show outputs the first item in file.
func Show(file string) error {
	i, err := first(file)

	if i != "" {
		fmt.Println(i)
	}

	return err
}

// Count outputs the count of read items in file.
func Count(file string) error {
	rl, err := Read(file)

	if err != nil {
		return err
	}

	// TODO: re-word these msgs
	fmt.Printf("Items in the reading list: %v\n", len(rl.Items))
	fmt.Printf("Items read:  %v\n", rl.Reads)
	fmt.Printf("Total items added: %v\n", rl.Added)

	return nil
}

// Read reads a ReadingList from file.
func Read(file string) (ReadingList, error) {
	b, err := ioutil.ReadFile(file) // TODO: create if not exist (?)
	if os.IsNotExist(err) {
		fmt.Println("No items in your reading list!")
		return ReadingList{}, err
	}

	if err != nil {
		return ReadingList{}, err
	}

	var rl ReadingList

	err = json.Unmarshal(b, &rl)

	if err != nil {
		return ReadingList{}, err
	}

	// TODO: read shouldn't care about this
	if len(rl.Items) == 0 {
		fmt.Println("No items in your reading list!")
		return ReadingList{}, nil
	}

	return rl, nil
}

// Write writes a ReadingList to file.
func Write(file string, rl ReadingList) error {
	b, err := json.Marshal(rl)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, b, os.ModeAppend)
}

func first(file string) (string, error) {
	rl, err := Read(file)

	if err != nil {
		return "", err
	}

	if len(rl.Items) == 0 {
		return "", nil
	}

	return rl.Items[0], nil
}

func filterEmpty(a []string) []string {
	b := a[:0]

	for _, s := range a {
		if s != "\n" && s != "" {
			b = append(b, s)
		}
	}

	return b
}
