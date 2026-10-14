package main

import (
	"fmt"
	"log"
	"os"

	"github.com/z3vxo/TokenJack/internal/core"
	"github.com/z3vxo/TokenJack/internal/setup"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		if err := setup.Run(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[*] Server setup!\n[+] run '%s run' to start\n", os.Args[0])
		return
	}

	server, err := core.New()
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
}
