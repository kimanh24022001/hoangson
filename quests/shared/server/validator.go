package server

import (
	"encoding/json"

	"smatyx.com/shared/cast"
)

type Validator struct {
	Keys  []string `json:"keys"`
	Texts []string `json:"texts"`
}

type ErrorTexts []string

// NOTE(duong): make a copy of Error and return that
func NewErrorRepCopy(err ErrorRep, message string, specific error) *ErrorRep {
	err.Message = message
	err.Specific = specific

	return &err
}

func (self *ErrorRep) Json() []byte {
	result, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return result
}

func (self *ErrorRep) Error() string {
	jsonBytes := self.Json()
	result := cast.BytesToString(jsonBytes)
	return result
}

func NewValidateError(capacity int) *Validator {
	result := &Validator{
		Keys:  make([]string, 0, capacity),
		Texts: make([]string, 0, capacity),
	}
	return result
}

func (self *Validator) Add(key, text string) {
	self.Keys = append(self.Keys, key)
	self.Texts = append(self.Texts, text)
}

func (self *Validator) Error() string {
	jsonBytes, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}

	return cast.BytesToString(jsonBytes)
}

func (self *Validator) Got() bool {
	return len(self.Keys) != 0
}

