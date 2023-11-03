package user

import (
	"encoding/json"

	"smatyx.com/shared/server"
)

type UserLoginReq struct {
	Phone    string
	Email    string
	Password string
	OTP      string
}

func UserLogin(txt *server.Transaction) {
	body, err := txt.ReadEntireBody()
	if err != nil {
		panic(err)
	}
	
	req := UserLoginReq{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
}
