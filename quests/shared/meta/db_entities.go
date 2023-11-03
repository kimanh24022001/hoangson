// NOTE(auto): This file is auto-generated. Please don't modify.
package meta

import (
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
)


var User = database.NewPgEntity("users", entities.User{})
const (
	User_Id             = `"users"."id"`
	User_Email          = `"users"."email"`
	User_Phone          = `"users"."phone"`
	User_FirstName      = `"users"."first_name"`
	User_LastName       = `"users"."last_name"`
	User_HashedPassword = `"users"."hashed_password"`
	User_PasswordSalt   = `"users"."password_salt"`
	User_State          = `"users"."state"`
	User_CreatedTime    = `"users"."created_time"`
	User_UpdatedTime    = `"users"."updated_time"`
)

const (
	User_Id_Idx             = 0
	User_Email_Idx          = 1
	User_Phone_Idx          = 2
	User_FirstName_Idx      = 3
	User_LastName_Idx       = 4
	User_HashedPassword_Idx = 5
	User_PasswordSalt_Idx   = 6
	User_State_Idx          = 7
	User_CreatedTime_Idx    = 8
	User_UpdatedTime_Idx    = 9
)

var Migration = database.NewPgEntity("migrations", entities.Migration{})
const (
	Migration_Id          = `"migrations"."id"`
	Migration_QueryText   = `"migrations"."query_text"`
	Migration_AppliedTime = `"migrations"."applied_time"`
)

const (
	Migration_Id_Idx          = 0
	Migration_QueryText_Idx   = 1
	Migration_AppliedTime_Idx = 2
)

var Contract = database.NewPgEntity("contracts", entities.Contract{})
const (
	Contract_CreatedTime = `"contracts"."created_time"`
	Contract_UpdatedTime = `"contracts"."updated_time"`
)

const (
	Contract_CreatedTime_Idx = 0
	Contract_UpdatedTime_Idx = 1
)

