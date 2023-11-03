package server

import (
	"encoding/json"

	"smatyx.com/shared/cast"
)

type SingleRep struct {
	Data   any         `json:"data,omitempty"`
	Errors []*ErrorRep `json:"errors,omitempty"`
}

type PageRep struct {
	Total  int         `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Data   any         `json:"data,omitempty"`
	Errors []*ErrorRep `json:"errors,omitempty"`
}

type ErrorRep struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Message  string `json:"message"`
	Specific error  `json:"specific,omitempty"`
}



func SingleRepJson(data any, errs []*ErrorRep) []byte {
	result, err := json.Marshal(SingleRep{Data: data, Errors: errs})
	if err != nil {
		panic(err)
	}

	return result
}

func PageRepJson(total, offset, limit int, data any, errs []*ErrorRep) []byte {
	result, err := json.Marshal(
		PageRep{
			Total:  total,
			Offset: offset,
			Limit:  limit,
			Data:   data,
			Errors: errs,
		})
	if err != nil {
		panic(err)
	}

	return result
}

func (self ErrorTexts) Error() string {
	jsonBytes, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}

	return cast.BytesToString(jsonBytes)
}
