package role

import "time"

type Role struct {
	Id       uint      `json:"id"`
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"updated"`
}
