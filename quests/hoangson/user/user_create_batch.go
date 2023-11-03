package user

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/meta"
	"smatyx.com/shared/security"
	"smatyx.com/shared/server"
	"smatyx.com/shared/types"
)

type UserCreateReq struct {
	Email     string
	Phone     string
	FirstName string
	LastName  string
	Password  string
}

type UserCreateBatchReq []UserCreateReq

type UserCreateRep struct {
	Id types.Uuid
}

type UserCreateBatchRep []UserCreateRep

func UserCreateBatch(txt *server.Transaction) {
	body, err := txt.ReadEntireBody()
	if err != nil {
		panic(err)
	}

	var batchReq UserCreateBatchReq
	err = json.Unmarshal(body, &batchReq)
	if err != nil {
		panic(err)
	}

	resultErrors := make([]*server.ErrorRep, len(batchReq))

	ok := true
	for i, req := range batchReq {
		validateErr := server.NewValidateError(5)
		{
			err := ValidateEmail(req.Email)
			if err != nil {
				validateErr.Add("email", err.Error())
			}
		}

		{
			err := ValidatePhone(req.Phone)
			if err != nil {
				validateErr.Add("phone", err.Error())
			}
		}

		// TODO: fill error messages
		if len(req.FirstName) == 0{
			validateErr.Add("firstName", "firstname")
		}

		if len(req.LastName) == 0{
			validateErr.Add("lastName", "lastname")
		}

		// TODO: Do stuff with password
		if len(req.Password) == 0 {
			validateErr.Add("password", "pw")
		}

		log.Printf("validate err: %v\n", validateErr)

		if validateErr.Got() {
			ok = false
			resultErrors[i] = server.NewErrorRepCopy(server.Error_InputInvalid, "", validateErr)			
		}
	}

	if !ok {
		txt.WriteBytes(server.SingleRepJson(nil, resultErrors))
		return
	}
	
	now := time.Now()
	users := make([]entities.User, 0, len(batchReq))

	for _, req := range batchReq {
		// NOTE(duong): might be slow
		salt := security.NewSalt()
		hash := security.HashPassword(req.Password, salt)

		user := entities.User{
			Id:             types.Uuid(uuid.New()),
			Email:          types.String(req.Email),
			Phone:          types.String(req.Phone),
			FirstName:      types.String(req.FirstName),
			LastName:       types.String(req.LastName),
			HashedPassword: types.String(hash),
			PasswordSalt:   types.String(salt),
			State:          types.I16(entities.UserState_Verified),
			CreatedTime:    types.Time(now),
			UpdatedTime:    types.Time(now),
		}

		users = append(users, user)
	}

	errs := meta.User_BatchInsert(txt, users)
	resultData := make([]any, len(users))

	for i, user := range users {
		err := errs[i]
		if err != nil {
			resultData[i] = nil

			cerr := EmailAndPhoneConstraintError(err, string(user.Email), string(user.Phone))
			if cerr != nil {
				resultErrors[i] = cerr
				continue
			}

			log.Printf("ERROR: %v\n", err)
			resultErrors[i] = &server.Error_Unknown
			continue
		}

		resultData[i] = UserCreateRep{user.Id}
		resultErrors[i] = nil
	}

	txt.WriteBytes(server.SingleRepJson(resultData, resultErrors))
}
