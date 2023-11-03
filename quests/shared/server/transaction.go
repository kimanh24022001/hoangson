package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"smatyx.com/shared/cast"
	"time"
)

type TransactionHandler struct {
	Timeout            time.Duration
	MaxRequestBodySize int
	Function           func(*Transaction)
}

type Transaction struct {
	Handler     *TransactionHandler
	HttpWriter  http.ResponseWriter
	HttpRequest *http.Request
	Context     context.Context
	Mutex       *sync.Mutex
}

func NewTransaction(
	handler *TransactionHandler,
	writer http.ResponseWriter,
	request *http.Request) *Transaction {

	return &Transaction{
		Handler:     handler,
		HttpWriter:  writer,
		HttpRequest: request,
		Context:     request.Context(),
		Mutex:       &sync.Mutex{},
	}
}

func (txt *Transaction) Header() http.Header {
	return txt.HttpWriter.Header()
}

func (txt *Transaction) StatusCode(statusCode int) {
	txt.HttpWriter.WriteHeader(statusCode)
}

// IMPORTANT(duong): This function is not thread safe, please be use this synchronously
func (txt *Transaction) WriteBytes(bytes []byte) error {
	n, err := txt.HttpWriter.Write(bytes)
	if n != len(bytes) {
		return errors.New("Something went wrong while trying to write a response")
	}

	return err
}

// IMPORTANT(duong): This function is not thread safe, please be use this
// synchronously
func (txt *Transaction) WriteString(s string) error {
	bytes := cast.StringToBytes(s)

	n, err := txt.HttpWriter.Write(bytes)
	if n != len(bytes) {
		return errors.New("Something went wrong while trying to write a response")
	}

	return err
}

func (txt *Transaction) ReadEntireBody() ([]byte, error) {
	bytes := make([]byte, txt.Handler.MaxRequestBodySize+1)

	reader := txt.HttpRequest.Body
	n, err := reader.Read(bytes)

	if err != io.EOF {
		if err == nil {
			return nil, errors.New("Body size exceeds the limit")
		}
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("The request contains no body")
	}

	bytes = bytes[:n]
	return bytes, nil
}
