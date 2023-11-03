package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"smatyx.com/config"
	"smatyx.com/shared/cast"
)

const (
	MethodGet = iota
	MethodPut
	MethodPost
	MethodDelete
	MethodCount
)

func MethodResolve(method string) (int, error) {
	if method == "GET" {
		return MethodGet, nil
	} else if method == "PUT" {
		return MethodPut, nil
	} else if method == "POST" {
		return MethodPost, nil
	} else if method == "DELETE" {
		return MethodDelete, nil
	}
	return -1, errors.New("Invalid method")
}

type Multiplexer struct {
	Entries [MethodCount]map[string]*TransactionHandler
}

func (multiplexer *Multiplexer) Map(method int, path string, handler *TransactionHandler) error {
	path, err := CleanPath(path)
	if err != nil {
		return err
	}

	for i := 0; i < MethodCount; i++ {
		if multiplexer.Entries[i] == nil {
			multiplexer.Entries[i] = make(map[string]*TransactionHandler)
		}
	}

	multiplexer.Entries[method][path] = handler

	return nil
}

func (multiplexer *Multiplexer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path, err := CleanPath(request.URL.Path)
	defer func() {
		// TODO: error log
	}()

	if err != nil {
		panic(err)
	}

	method, err := MethodResolve(request.Method)
	if err != nil {
		panic(err)
	}

	handler, ok := multiplexer.Entries[method][path]
	if !ok {
		// TODO: return not found
		panic("404")
	}

	ctx, cancelCtx := context.WithTimeout(request.Context(), handler.Timeout)
	defer cancelCtx()
	done := make(chan struct{})
	panicChannel := make(chan any, 1)
	request = request.WithContext(ctx)

	txt := NewTransaction(handler, writer, request)

	go func() {
		defer close(done)

		defer func() {
			if p := recover(); p != nil {
				// switch t := p.(type) {
				// case *internal.Error:
				// 	txt.StatusCode(http.StatusBadRequest)
				// 	txt.WriteBytes(t.Json())

				// default:
				// 	// NOTE(duong): Only turn this off for the production
				// 	// environment. In the development environment, we
				// 	// should acknowledge where the error is.
				// 	if config.Debug {
				// 		panicChannel <- p
				// 		panic(p)
				// 	}
				// }

				if config.Debug {
					panicChannel <- p
					panic(p)
				}

				// NOTE: print out the stack to avoid losing it.
				debug.PrintStack()
			}
		}()
		handler.Function(txt)
	}()

	select {
	case p := <-panicChannel:
		panic(p)
	case <-done:
		txt.Mutex.Lock()
		defer txt.Mutex.Unlock()
	case <-ctx.Done():
		txt.Mutex.Lock()
		defer txt.Mutex.Unlock()
		err := ctx.Err()
		if err == context.DeadlineExceeded {
			writer.WriteHeader(http.StatusServiceUnavailable)
			// TODO: make better error body
			writer.Write([]byte("Timeout"))
		} else {
			writer.WriteHeader(http.StatusServiceUnavailable)
			// TODO: write error body
		}
		log.Printf("%v\n", err)
	}
}

var PathContainsInvalidCharacter = errors.New("Path contains invalid character")

// To avoid mistakes and to consider the combination of "method" and
// "path" as similar to calling a "function", our path will not
// contain "/.." or "/.".
//
// When we encounter "/.." or "." in the path, return an error.
//
// Then do something like:
// - Replace multiple slashes with a single slash.
// - Remove the "/" at the begining and the end of path
func CleanPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if strings.Contains(path, "/.") {
		return path, PathContainsInvalidCharacter
	}

	src := cast.StringToBytes(path)
	dst := make([]byte, len(path))
	srcRead, dstWrite := 0, 0
	justSlash := false

	// NOTE(duong): Replace multiple slashes with a single slash.
	{
		for i := 0; i < len(path); i++ {
			if path[i] == '/' {
				if !justSlash && srcRead != i {
					copy(dst[dstWrite:], src[srcRead:i])
					dstWrite += i - srcRead
					justSlash = true
				}
				continue
			}

			if justSlash {
				justSlash = false
				srcRead = i
				dst[dstWrite] = '/'
				dstWrite++
			}
		}
		if !justSlash {
			copy(dst[dstWrite:], src[srcRead:])
			dstWrite += len(path) - srcRead
		}
		dst = dst[:dstWrite]
		path = cast.BytesToString(dst)
	}

	if len(path) != 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if len(path) != 0 && path[0] == '/' {
		path = path[1:]
	}

	return path, nil
}

func StringToBytes(path string) {
	panic("unimplemented")
}
