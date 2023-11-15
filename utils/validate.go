package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

func ValidateFields(s interface{}) []string {
	validate := validator.New()
	errMsgs := make([]string, 0)

	if errors := validate.Struct(s); errors != nil {
		for _, err := range errors.(validator.ValidationErrors) {
			validationErr := strings.Split(fmt.Sprint(err), "Error:")[1]
			errMsgs = append(errMsgs, fmt.Sprintf("Error: %s", validationErr))
		}

		return errMsgs
	}

	return nil
}
