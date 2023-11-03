// NOTE(auto): This file is auto-generated. Please don't modify.
package main

import (
	"os"
	"path/filepath"
	"runtime"

	"smatyx.com/shared/meta"
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/migrate"
)

var basePath = func () string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}()
func realPath(pathName string) string{
	if len(pathName) == 0 && pathName[0] == '/' {
		return pathName
	}

	return basePath + "/" + pathName
}

func main() {
	callsFile := realPath("../../shared/meta/db_calls.go")
	entitiesFile := realPath("../../shared/meta/db_entities.go")

	os.WriteFile(callsFile, []byte{}, 0666)
	os.WriteFile(entitiesFile, []byte{}, 0666)

	callsFd, err := os.OpenFile(callsFile, os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer callsFd.Close()

	entitiesFd, err := os.OpenFile(entitiesFile, os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer entitiesFd.Close()

	callsFd.WriteString(meta.DbCallsHeader)
	entitiesFd.WriteString(meta.DbEntitiesHeader)
	// NOTE(auto): Migrate entities.User
	{
		pgEntity := database.NewPgEntity("users", entities.User{})
		migrate.AddMigrate("users", pgEntity)
		callsFd.WriteString(meta.MakeDbCalls(pgEntity))
		entitiesFd.WriteString(meta.MakeDbEntities(pgEntity))
	}

	// NOTE(auto): Migrate entities.Migration
	{
		pgEntity := database.NewPgEntity("migrations", entities.Migration{})
		migrate.AddMigrate("migrations", pgEntity)
		callsFd.WriteString(meta.MakeDbCalls(pgEntity))
		entitiesFd.WriteString(meta.MakeDbEntities(pgEntity))
	}

	// NOTE(auto): Migrate entities.Contract
	{
		pgEntity := database.NewPgEntity("contracts", entities.Contract{})
		migrate.AddMigrate("contracts", pgEntity)
		callsFd.WriteString(meta.MakeDbCalls(pgEntity))
		entitiesFd.WriteString(meta.MakeDbEntities(pgEntity))
	}

	migrate.DoMigrate()
}