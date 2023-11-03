package entities

import "smatyx.com/shared/types"

type UserState types.I16

const (
	UserState_Created UserState = iota
	UserState_Verified
)

//entity:table
type User struct {
	Id  types.Uuid

	Email          types.String `unique:"users_email"`
	Phone          types.String `unique:"users_phone"`
	FirstName      types.String
	LastName       types.String
	HashedPassword types.String
	PasswordSalt   types.String
	State          types.I16

	CreatedTime types.Time
	UpdatedTime types.Time
}
