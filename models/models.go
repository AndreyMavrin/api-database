package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jackc/pgtype"
)

var DB *sql.DB

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

type UserUpdate struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type Forum struct {
	Posts   int    `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Post struct {
	Author   string           `json:"author"`
	Created  time.Time        `json:"created"`
	Forum    string           `json:"forum"`
	ID       int              `json:"id"`
	IsEdited bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   JsonNullInt64    `json:"parent"`
	Thread   int              `json:"thread,"`
	Path     pgtype.Int8Array `json:"-"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"-"`
}

type PostFull struct {
	Author *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread"`
}

type Status struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}

type Error struct {
	Message string `json:"message"`
}
