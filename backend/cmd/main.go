package main

import (
	"fmt"
	"log"
	"os"

	"github.com/z3vxo/TokenJack/internal/server"
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

	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}
