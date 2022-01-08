package user

import "time"

type User struct {
	Id         int64 `datastore:"-"`
	Name       string
	Email      string
	Photo      string
	Timestamp  time.Time
	Timeupdate time.Time
}
