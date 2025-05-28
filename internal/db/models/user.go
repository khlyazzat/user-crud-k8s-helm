package models

import (
	uuid "github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name          string
	Email         string
	Age           int
}
