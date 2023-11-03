package entities

import (
	"smatyx.com/shared/types"
)

type UUID types.Uuid

//entity:table
type Migration struct {
	Id          *types.String
	QueryText   *types.String
	AppliedTime *types.Time
}

//entity:table
type Contract struct {
	Id *UUID

	CreateUserId *UUID

	CreatedTime *types.Time
	UpdatedTime *types.Time
}
