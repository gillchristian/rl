// Package rl provides methods to handle reading list files.
package rl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	xdgOpen "github.com/skratchdot/open-golang/open"
)

// Add adds an item i at the end of the file.
func Add(file, i string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	i += "\n"

	_, err = f.WriteString(i)

	return err
}

// Done removes the last item in file.
func Done(file string) error {
	rl, err := read(file)

	if err != nil {
		return err
	}

	f := []byte(strings.Join(rl[1:], "\n") + "\n")

	return ioutil.WriteFile(file, f, os.ModeAppend)
}

// Open opens with xdg-open (or the Mac/Windows equivalent) the first item in file.
func Open(file string) error {
	i, err := first(file)

	if i != "" {
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

func first(file string) (string, error) {
	rl, err := read(file)

	if err != nil {
		return "", err
	}

	return rl[0], nil
}

func read(file string) ([]string, error) {
	b, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		fmt.Println("No items in your reading list!")
		return []string{}, err
	}

	if err != nil {
		return []string{}, err
	}

	rl := filterEmpty(strings.Split(string(b), "\n"))

	if len(rl) == 0 {
		fmt.Println("No items in your reading list!")
		return []string{}, nil
	}

	return rl, nil
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
