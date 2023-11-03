package user

import (
	"time"

	"smatyx.com/shared/cast"
	"smatyx.com/shared/server"
)

func InitUser(multiplexer *server.Multiplexer) {
	multiplexer.Map(server.MethodGet,
		"/users",
		&server.TransactionHandler{
			Timeout:            5 * time.Second,
			MaxRequestBodySize: cast.Kilobyte(64),
			Function:           UserSearch,
		})

	multiplexer.Map(server.MethodPut,
		"/users",
		&server.TransactionHandler{
			Timeout:            5 * time.Second,
			MaxRequestBodySize: cast.Kilobyte(64),
			Function:           UserCreateBatch,
		})

	multiplexer.Map(server.MethodPut,
		"/users/register",
		&server.TransactionHandler{
			Timeout:            5 * time.Second,
			MaxRequestBodySize: cast.Kilobyte(64),
			Function:           UserRegister,
		})
}
