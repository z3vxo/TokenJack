package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func InsertUser(c *sql.DB, user, pass string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	q := `INSERT INTO users(username, password) VALUES(?,?);`
	_, err = c.Exec(q, user, string(hash))

	return err

}

var ErrInvalidCredentials = errors.New("invalid creds")
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.DefaultCost)

func (d *DB) CheckAuth(user, pass string) error {
	var hash []byte
	err := d.conn.QueryRow(`SELECT password FROM users WHERE username = ?`, user).Scan(&hash)

	if errors.Is(err, sql.ErrNoRows) {
		bcrypt.CompareHashAndPassword(dummyHash, []byte(pass))
		return ErrInvalidCredentials
	}

	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword(hash, []byte(pass)) != nil {
		return ErrInvalidCredentials
	}
	return nil

}
