<h1 align="center">reading-list</h1>

A FIFO reading list CLI.

## Install

```
$ go get github.com/gillchristian/rl/cmd/rl
```

## Use

```
NAME:
   rl - FIFO reading list

USAGE:
   $ rl        # show next item
   $ rl [item] # add item

VERSION:
   0.0.1

AUTHOR:
   Christian Gill (gillchristiang@gmail.com)

COMMANDS:
     add      Add item to the reading list.
     done     Remove next item from the reading list and increase the count of read items.
     rm       Remove next item from the reading list (does not increase the count).
     count    Display the amount of read items.
     open     Open next item in the browser.
     show     Show next in the reading list.
     sync     sync current file with a remote one (GitHub Gist).
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
