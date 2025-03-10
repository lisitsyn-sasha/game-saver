package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func OK() Response {
	return Response{
		Status: "StatusOk",
	}
}

func Error(msg string) Response {
	return Response{
		Status: "StatusError",
		Error:  msg,
	}
}

// ValidationError предназначена для обработки ошибок валидации
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: "StatusError",
		Error:  strings.Join(errMsgs, ", "),
	}
}
