package model

import "github.com/google/uuid"

type Foo struct {
	UUID uuid.UUID `json:"foo_uuid"`
	Msg  string    `json:"msg"`
}
