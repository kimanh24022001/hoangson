package server

var (
	Error_Unknown = ErrorRep{
		Code:    "0-0000",
		Name:    "Internal",
		Message: "This error is not specific; please contact the system administrator to determine what happened in the system."}
	Error_InputInvalid = ErrorRep{
		Code: "0-0001",
		Name: "InputInvalid"}
)
