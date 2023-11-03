package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"smatyx.com/shared/cast"
	"smatyx.com/shared/database"
	"smatyx.com/shared/server"

	"smatyx.com/hoangson/user"
)

const DbUrl = "postgres://zvrujuoc:RiIBrsrws4EFvbnNFbjFcCs4qu2RPp4P@tiny.db.elephantsql.com/zvrujuoc?sslmode=disable&pool_max_conns=16"
// const DbUrl = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&pool_max_conns=16"

func main() {
	time.Local = time.UTC
	database.InitPg(DbUrl)

	multiplexer := server.Multiplexer{}
	multiplexer.Map(server.MethodGet, "/", &server.TransactionHandler{
		Timeout:            5 * time.Second,
		MaxRequestBodySize: cast.Kilobyte(64),
		Function: func(txt *server.Transaction) {
			txt.WriteString("Hello from smatyx!\n")
		},
	})

	user.InitUser(&multiplexer)

	serv := http.Server{
		Addr:                         ":8080",
		Handler:                      &multiplexer,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  1 * time.Second,
		MaxHeaderBytes:               cast.Kilobyte(64),
	}

	fmt.Println(server.CoolSmatyxLogo)
	log.Printf("Server is running... Look at %v\n", serv.Addr)
	err := serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
