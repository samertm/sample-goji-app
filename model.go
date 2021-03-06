package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/samertm/syncfbevents/db"
)

type User struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	AccessToken sql.NullString `db:"access_token"`
	ExpiresOn   *time.Time     `db:"expires_on"`
}

var userSchema = `
CREATE TABLE IF NOT EXISTS person (
  id SERIAL PRIMARY KEY,
  name TEXT,
  access_token TEXT,
  expires_on TIMESTAMP
)
`

func init() {
	db.DB.MustExec(userSchema)
}

type UserSpec struct {
	ID   int
	Name string
}

func GetCreateUser(name string) (User, error) {
	// Try to get the user once.
	u, err := GetUser(UserSpec{Name: name})
	if err == nil {
		// User exists, return them.
		return u, nil
	}
	// Create the user and then get them.
	if err := CreateUser(name); err != nil {
		return User{}, err
	}
	// Get the user one last time.
	return GetUser(UserSpec{Name: name})
}

func CreateUser(name string) error {
	query := "INSERT INTO person(name) VALUES ($1)"
	if _, err := db.DB.Exec(query, name); err != nil {
		return err
	}
	return nil
}

func GetUser(us UserSpec) (User, error) {
	u := User{}
	where := struct {
		col string
		val string
	}{}
	if us.ID != 0 {
		where.col = "id"
		where.val = strconv.Itoa(us.ID)
	} else if us.Name != "" {
		where.col = "name"
		where.val = us.Name
	} else {
		return User{}, errors.New("Empty user spec")
	}

	err := db.DB.Get(&u, fmt.Sprintf("SELECT * from person where %s=$1", where.col), where.val)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func SetAccessToken(u User, token string, expiresIn string) error {
	e, err := strconv.Atoi(expiresIn)
	if err != nil {
		return err
	}
	expiresOn := time.Now().Add(time.Duration(e) * time.Second)
	b := &db.Binder{}
	query := "UPDATE person SET access_token = " + b.Bind(token) + ", " +
		"expires_on = " + b.Bind(expiresOn) + " " +
		"WHERE id = " + b.Bind(u.ID)
	if _, err := db.DB.Exec(query, b.Items...); err != nil {
		return err
	}
	return nil
}
