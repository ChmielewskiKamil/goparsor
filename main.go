package main

import (
	"fmt"
	"goparsor/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s!\n", user.Username)
	fmt.Printf("Type in commands and create your new funky project, you monkey!\n")
	repl.Start(os.Stdin)
}
