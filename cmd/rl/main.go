package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gillchristian/rl"
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
			Usage:     "Remove next item from the reading list and increase the count of read items.",
			UsageText: "$ rl done",
			Action:    func(c *cli.Context) error { return rl.Done(fileName) },
		},
		{
			Name:      "rm",
			Usage:     "Remove next item from the reading list (does not increase the count).",
			UsageText: "$ rl rm",
			Action:    func(c *cli.Context) error { return rl.Remove(fileName) },
		},
		{
			Name:      "count",
			Usage:     "Display the amount of read items.",
			UsageText: "$ rl count",
			Action:    func(c *cli.Context) error { return rl.Count(fileName) },
		},
		{
			Name:      "open",
			Usage:     "Open next item in the browser.",
			UsageText: "$ rl open",
			Action:    func(c *cli.Context) error { return rl.Open(fileName) },
		},
		{
			Name:      "show",
			Usage:     "Show next in the reading list.",
			UsageText: "$ rl show",
			Action:    func(c *cli.Context) error { return rl.Show(fileName) },
		},
	}

	_ = app.Run(os.Args)
}

func addOrshow(c *cli.Context) error {
	if c.NArg() == 0 {
		return rl.Show(fileName)
	}

	return rl.Add(fileName, c.Args()[0])
}

func add(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Println("Nothing to add!")
		return nil
	}

	return rl.Add(fileName, c.Args()[0])
}
