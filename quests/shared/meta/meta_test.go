package meta

import (
	"testing"

	"smatyx.com/internal/database"
	"smatyx.com/internal/entities"
)

func TestMakeDbCalls(t *testing.T) {
	entity := database.NewPgEntity("users", entities.User{})
	MakeDbCalls(entity)
}
