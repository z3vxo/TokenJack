package setup

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/z3vxo/TokenJack/internal/db"

	"golang.org/x/term"
)

func GenUserAccount(conn *sql.DB) error {
	if user := os.Getenv("JACKER_USER"); user != "" {
		if pass := os.Getenv("JACKER_PASSWD"); pass != "" {
			fmt.Println("[*] Env vars detected")
			return db.InsertUser(conn, user, pass)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	var user string
	for user == "" {
		fmt.Print("[+] Enter Username: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		user = strings.TrimSpace(line)
	}

	fmt.Print("[*] Auto generate password? [Y/N]: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	ans := strings.ToLower(strings.TrimSpace(line))

	var password string
	if ans == "n" || ans == "no" {
		password, err = AskPassword()
		if err != nil {
			return err
		}
	} else {
		password, err := genPassword()
		if err != nil {
			return err
		}
		fmt.Printf("[*] Generated Password(save now, wont be shown again): %s\n", password)
	}

	return db.InsertUser(conn, user, password)
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func genPassword() (string, error) {
	b := make([]byte, 18)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}

func AskPassword() (string, error) {
	for {
		p1, err := readHidden("[+] Enter password: ")
		if err != nil {
			return "", err
		}
		p2, err := readHidden("[+] Confirm password: ")
		if err != nil {
			return "", err
		}

		if p1 != p2 {
			fmt.Println("[!] Passwords dont match")
			continue
		}
		return p1, nil
	}
}

func readHidden(prompt string) (string, error) {
	fmt.Print(prompt)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(b), err
}
