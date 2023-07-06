package model

import "github.com/google/uuid"

type Foo struct {
	UUID  uuid.UUID `json:"foo_uuid" mapstructure:"foo_uuid"`
	Msg   string    `json:"msg" mapstructure:"msg"`
	Email string    `json:"email" mapstructure:"email" validate:"email"`
	Phone string    `json:"phone" mapstructure:"phone" validate:"e164"`
}

// PostFoo is the request body of the POST foo endpoint
type PostFoo struct {
	Msg   string `json:"msg"  mapstructure:"msg" validate:"lowercase,required"`
	Email string `json:"email" mapstructure:"email" validate:"email,required"`
	Phone string `json:"phone" mapstructure:"phone" validate:"e164,required"`
}

// Update foo is the request body of the PUT foo endpoint
type UpdateFoo struct {
	Msg   string `json:"msg" mapstructure:"msg" validate:"lowercase"`
	Phone string `json:"phone" mapstructure:"phone" validate:"e164"`
}
