package users

import "github.com/joeychilson/lixy/pkg/context"

const ContextKey context.ContextKey = "user"

type User struct {
	ID       int32
	Email    string
	Username string
}
