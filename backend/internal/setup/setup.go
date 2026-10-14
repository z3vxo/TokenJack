package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/z3vxo/TokenJack/internal/config"
	"github.com/z3vxo/TokenJack/internal/db"
)

func Run() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	base := filepath.Join(home, ".local", "share", "token-jack")

	dir := []string{
		filepath.Join(base, "database"),
		filepath.Join(base, "logs"),
		filepath.Join(base, "config"),
		filepath.Join(base, "user-files")}

	files := []string{
		filepath.Join(base, "database", "jacker.db"),
		filepath.Join(base, "config", "jacker.conf"),
		filepath.Join(base, "logs", "access.log"),
		filepath.Join(base, "logs", "error.log"),
		filepath.Join(base, "logs", "genernal.log"),
	}

	for _, d := range dir {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	for _, f := range files {
		file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		file.Close()
	}

	dbPath := filepath.Join(base, "database", "jacker.db")
	cfgPath := filepath.Join(base, "config", "jacker.conf")

	fmt.Printf("[*] Directorys and Files created: %s\n", base)

	db, err := db.OpenAndMigrate(dbPath)
	if err != nil {
		return err
	}
	fmt.Print("[*] Database setup\n")

	if err := GenUserAccount(db); err != nil {
		return err
	}
	db.Close()

	if err := config.Generate(cfgPath); err != nil {
		return err
	}
	fmt.Printf("[*] Default config written to %s\n", cfgPath)

	return nil
}
