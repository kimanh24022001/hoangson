package user

import (
	"smatyx.com/shared/server"
	"smatyx.com/shared/database"
)

var (
	Error_ExistEmail = server.ErrorRep{
		Code: "U-0001",
		Name: "ExistEmail",
	}

	Error_ExistPhone = server.ErrorRep{
		Code: "U-0002",
		Name: "ExistPhone",
	}
)

func EmailAndPhoneConstraintError(err error, email, phone string) *server.ErrorRep {
	pgErr := database.PgError(err)

	if pgErr == nil {
		return nil
	}

	if !(pgErr.Code == database.ErrCode_UniqueViolation) {
		return nil
	}

	if pgErr.ConstraintName == "users_email" {
		return server.NewErrorRepCopy(Error_ExistEmail,
			"The user with this email already exists.",
			server.ErrorTexts([]string{email}))
	}

	if pgErr.ConstraintName == "users_phone" {
		return server.NewErrorRepCopy(Error_ExistPhone,
			"The user with this phone already exists.",
			server.ErrorTexts([]string{phone}))
	}

	return nil
}
