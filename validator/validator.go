package validator

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var isValidEmail = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`).MatchString

func Init() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("custom validator initialization failed")
	}

	v.RegisterValidation("isEmail", isEmail)

	return nil
}

func isEmail(f validator.FieldLevel) bool {
	return isValidEmail(f.Field().String())
}

func GetValidationError(err error) *[]validationErrMsg {
	var validationErrMsgs []validationErrMsg
	for _, err := range err.(validator.ValidationErrors) {
		f := err.Field()
		t := err.Tag()
		p := err.Param()
		validationErrMsgs = append(validationErrMsgs, validationErrMsg{
			Field:   f,
			Message: getErrMessage(t, p),
		})
	}
	return &validationErrMsgs

}

type validationErrMsg struct {
	Field   string
	Message string
}

func getErrMessage(tag string, param string) string {
	errmessage := ""
	switch tag {
	case "isEmail":
		errmessage = "invalid email format"
	case "min":
		errmessage = "min required length " + param
	case "max":
		errmessage = "max supported length" + param
	default:
		errmessage = "should be " + param
	}
	return errmessage
}
