package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	xdgOpen "github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

var fileName string

func init() {
	usr, _ := user.Current()
	dir := usr.HomeDir

	fileName = filepath.Join(dir, ".reading-list")
}

func main() {
	app := cli.NewApp()

	app.Name = "rl"
	app.Version = "0.0.1"
	app.Author = "Christian Gill (gillchristiang@gmail.com)"
	app.Usage = "FIFO reading list"
	app.UsageText = "$ rl        # show next item\n   $ rl [item] # add item"

	app.Action = addOrshow

	app.Commands = []cli.Command{
		{
			Name:      "add",
			Usage:     "Add item to the reading list.",
			UsageText: "$ rl add item",
			Action:    add,
		},
		{
			Name:      "done",
			Usage:     "Remove next item from the reading list.",
			UsageText: "$ rl done",
			Action:    done,
		},
		{
			Name:      "open",
			Usage:     "Open next item in the browser.",
			UsageText: "$ rl open",
			Action:    open,
		},
		{
			Name:      "show",
			Usage:     "Show next in the reading list.",
			UsageText: "$ rl show",
			Action:    show,
		},
	}

	_ = app.Run(os.Args)
}

func addOrshow(c *cli.Context) error {
	if c.NArg() == 0 {
		return show(c)
	}

	return add(c)
}

func add(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Println("Nothing to add!")
		return nil
	}

	i := c.Args()[0] + "\n"

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	_, err = f.WriteString(i)

	return err
}

func done(c *cli.Context) error {
	b, err := ioutil.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("No items in your reading list!")
		return err
	}

	if err != nil {
		return err
	}

	rl := filterEmpty(strings.Split(string(b), "\n"))

	if len(rl) == 0 {
		fmt.Println("No items in your reading list!")
		return nil
	}

	f := []byte(strings.Join(rl[1:], "\n") + "\n")

	return ioutil.WriteFile(fileName, f, os.ModeAppend)
}

func open(c *cli.Context) error {
	b, err := ioutil.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("No items in your reading list!")
		return err
	}

	if err != nil {
		return err
	}

	rl := filterEmpty(strings.Split(string(b), "\n"))

	if len(rl) == 0 {
		fmt.Println("No items in your reading list!")
		return nil
	}

	return xdgOpen.Start(rl[0])
}

func show(c *cli.Context) error {
	b, err := ioutil.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("No items in your reading list!")
		return err
	}

	if err != nil {
		return err
	}

	rl := filterEmpty(strings.Split(string(b), "\n"))

	if len(rl) == 0 {
		fmt.Println("No items in your reading list!")
		return nil
	}

	fmt.Println(rl[0])

	return nil
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
