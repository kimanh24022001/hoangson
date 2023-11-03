// NOTE(auto): This file is auto-generated. Please don't modify.
package meta

import (
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/server"
)

func User_Insert(txt *server.Transaction, entity *entities.User) error {
	queryText := `INSERT INTO "users"("id","email","phone","first_name","last_name","hashed_password","password_salt","state","created_time","updated_time") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`
	query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		entity.Id, entity.Email, entity.Phone, entity.FirstName, entity.LastName, entity.HashedPassword, entity.PasswordSalt, entity.State, entity.CreatedTime, entity.UpdatedTime)

	err := query.Submit()

	return err
}

func User_Get(txt *server.Transaction, entity *entities.User) error {
	queryText := `SELECT * FROM  "users" WHERE "users"."id"=$1;`
	query := database.NewReadQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		"1eef0663-1da1-493e-9171-aeb44e6b32ff")

	err := query.Submit()
	print(1, quẻ)
	return err
}

/*
	func User_ResetPassword(txt *server.Transaction, entity *entities.User) error {
		queryText := `UPDATE  "users" SET "PassWord" = $1 WHERE "users"."id"=$2;`
		query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
			entity.PassWord, entity.Id)

		err := query.Submit()

		return err
	}
*/
func User_BatchInsert(txt *server.Transaction, list []entities.User) []error {
	l := len(list)

	queryText := `INSERT INTO "users"("id","email","phone","first_name","last_name","hashed_password","password_salt","state","created_time","updated_time") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`
	queries := database.NewQueries(txt.Context, l)
	for _, entity := range list {
		queries.AddWriteQuery(database.ValidateAffectedCountEqual(1), queryText,
			entity.Id, entity.Email, entity.Phone, entity.FirstName, entity.LastName, entity.HashedPassword, entity.PasswordSalt, entity.State, entity.CreatedTime, entity.UpdatedTime)
	}
	errs := queries.Submit()

	return errs
}

func User_Update(txt *server.Transaction, entity *entities.User) error {
	queryText := `UPDATE "users" SET "id"=$1,"email"=$2,"phone"=$3,"first_name"=$4,"last_name"=$5,"hashed_password"=$6,"password_salt"=$7,"state"=$8,"created_time"=$9,"updated_time"=$10 WHERE "users"."id"=$11;`
	query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		entity.Id, entity.Email, entity.Phone, entity.FirstName, entity.LastName, entity.HashedPassword, entity.PasswordSalt, entity.State, entity.CreatedTime, entity.UpdatedTime, entity.Id)

	err := query.Submit()

	return err
}

func User_BatchUpdate(txt *server.Transaction, list []entities.User) []error {
	l := len(list)

	queryText := `UPDATE "users" SET "id"=$1,"email"=$2,"phone"=$3,"first_name"=$4,"last_name"=$5,"hashed_password"=$6,"password_salt"=$7,"state"=$8,"created_time"=$9,"updated_time"=$10 WHERE "users"."id"=$11;`
	queries := database.NewQueries(txt.Context, l)
	for _, entity := range list {
		queries.AddWriteQuery(database.ValidateAffectedCountEqual(1), queryText,
			entity.Id, entity.Email, entity.Phone, entity.FirstName, entity.LastName, entity.HashedPassword, entity.PasswordSalt, entity.State, entity.CreatedTime, entity.UpdatedTime, entity.Id)
	}
	errs := queries.Submit()

	return errs
}

func Migration_Insert(txt *server.Transaction, entity *entities.Migration) error {
	queryText := `INSERT INTO "migrations"("id","query_text","applied_time") VALUES($1,$2,$3);`
	query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		entity.Id, entity.QueryText, entity.AppliedTime)

	err := query.Submit()

	return err
}

func Migration_BatchInsert(txt *server.Transaction, list []entities.Migration) []error {
	l := len(list)

	queryText := `INSERT INTO "migrations"("id","query_text","applied_time") VALUES($1,$2,$3);`
	queries := database.NewQueries(txt.Context, l)
	for _, entity := range list {
		queries.AddWriteQuery(database.ValidateAffectedCountEqual(1), queryText,
			entity.Id, entity.QueryText, entity.AppliedTime)
	}
	errs := queries.Submit()

	return errs
}

func Migration_Update(txt *server.Transaction, entity *entities.Migration) error {
	queryText := `UPDATE "migrations" SET "id"=$1,"query_text"=$2,"applied_time"=$3 WHERE "migrations"."id"=$4;`
	query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		entity.Id, entity.QueryText, entity.AppliedTime, entity.Id)

	err := query.Submit()

	return err
}

func Migration_BatchUpdate(txt *server.Transaction, list []entities.Migration) []error {
	l := len(list)

	queryText := `UPDATE "migrations" SET "id"=$1,"query_text"=$2,"applied_time"=$3 WHERE "migrations"."id"=$4;`
	queries := database.NewQueries(txt.Context, l)
	for _, entity := range list {
		queries.AddWriteQuery(database.ValidateAffectedCountEqual(1), queryText,
			entity.Id, entity.QueryText, entity.AppliedTime, entity.Id)
	}
	errs := queries.Submit()

	return errs
}

func Contract_Insert(txt *server.Transaction, entity *entities.Contract) error {
	queryText := `INSERT INTO "contracts"("created_time","updated_time") VALUES($1,$2);`
	query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		entity.CreatedTime, entity.UpdatedTime)

	err := query.Submit()

	return err
}

func Contract_BatchInsert(txt *server.Transaction, list []entities.Contract) []error {
	l := len(list)

	queryText := `INSERT INTO "contracts"("created_time","updated_time") VALUES($1,$2);`
	queries := database.NewQueries(txt.Context, l)
	for _, entity := range list {
		queries.AddWriteQuery(database.ValidateAffectedCountEqual(1), queryText,
			entity.CreatedTime, entity.UpdatedTime)
	}
	errs := queries.Submit()

	return errs
}

func Contract_Update(txt *server.Transaction, entity *entities.Contract) error {
	queryText := `UPDATE "contracts" SET "created_time"=$1,"updated_time"=$2 WHERE "contracts"."id"=$3;`
	query := database.NewWriteQuery(txt.Context, database.ValidateAffectedCountEqual(1), queryText,
		entity.CreatedTime, entity.UpdatedTime, entity.Id)

	err := query.Submit()

	return err
}

func Contract_BatchUpdate(txt *server.Transaction, list []entities.Contract) []error {
	l := len(list)

	queryText := `UPDATE "contracts" SET "created_time"=$1,"updated_time"=$2 WHERE "contracts"."id"=$3;`
	queries := database.NewQueries(txt.Context, l)
	for _, entity := range list {
		queries.AddWriteQuery(database.ValidateAffectedCountEqual(1), queryText,
			entity.CreatedTime, entity.UpdatedTime, entity.Id)
	}
	errs := queries.Submit()

	return errs
}
