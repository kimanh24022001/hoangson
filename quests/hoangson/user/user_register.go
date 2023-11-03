package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/meta"
	"smatyx.com/shared/server"
	"smatyx.com/shared/types"
)

type UserRegisterReq struct {
	Email string
	Phone string
}

type UserRegisterRep struct {
	Ok bool `json:"ok"`
}

//api:doc user_register.md
func UserRegister(txt *server.Transaction) {
	body, err := txt.ReadEntireBody()
	if err != nil {
		panic(err)
	}

	req := UserRegisterReq{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	err = ValidateUserRegisterReq(&req)
	if err != nil {
		panic(err)
	}

	var now = time.Now()
	user := entities.User{
		Id:          types.Uuid(uuid.New()),
		State:       types.I16(0),
		CreatedTime: types.Time(now),
		UpdatedTime: types.Time(now),
	}
	if len(req.Email) != 0 {
		user.Email = types.String(req.Email)
	} else {
		user.Phone = types.String(req.Phone)
	}

	err = meta.User_Get(txt, &user)

	resultErrors := make([]*server.ErrorRep, 1)
	if err != nil {
		cErr := EmailAndPhoneConstraintError(err, req.Email, req.Phone)
		if cErr != nil {
			resultErrors[0] = cErr
		} else {
			fmt.Printf("ERROR: %v\n", err)
			resultErrors[0] = &server.Error_Unknown
		}
	}
	resultData := UserRegisterRep{Ok: (err == nil)}

	// TODO(duong): deliver a message to the user's mobile device or
	// email

	txt.WriteBytes(server.SingleRepJson(resultData, resultErrors))
}

func ValidateUserRegisterReq(req *UserRegisterReq) error {
	validateErr := server.NewValidateError(2)

	if len(req.Email) != 0 {
		err := ValidateEmail(req.Email)

		if err != nil {
			validateErr.Add("email", err.Error())

			return server.NewErrorRepCopy(
				server.Error_InputInvalid,
				"User input invalid email",
				validateErr)
		}

	} else if len(req.Phone) != 0 {
		err := ValidatePhone(req.Phone)
		if err != nil {
			validateErr.Add("phone", err.Error())

			return server.NewErrorRepCopy(
				server.Error_InputInvalid,
				"User entered an invalid phone number",
				validateErr)
		}
	} else {
		validateErr.Add("email", "Required at least 1 of 2 field: email or phone")
		validateErr.Add("phone", "Required at least 1 of 2 field: email or phone")

		return server.NewErrorRepCopy(
			server.Error_InputInvalid,
			"Required at least 1 of 2 field: email or phone",
			validateErr)
	}

	return nil
}
