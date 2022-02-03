package defined

import (
	"database/sql"
	"time"
)

type Role struct {
	Id       uint           `json:"id"`
	Name     string         `json:"name"`
	Desc     sql.NullString `json:"desc"`
	Created  time.Time      `json:"created"`
	Modified time.Time      `json:"updated"`
}
