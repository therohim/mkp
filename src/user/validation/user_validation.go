package validation

import (
	"net/mail"
	"test-mkp/src/user/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserValidation struct{}

func (userValidation *UserValidation) LoginValidate(request model.LoginRequest) error {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.Identity, validation.Required),
		validation.Field(&request.Password, validation.Required, validation.Length(6, 999).Error("min 6 character")),
	)

	return err
}

func (userValidation *UserValidation) RegisterValidate(request model.RegisterRequest) error {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Email, validation.Required, validation.When(userValidation.validationEmail(request.Email), validation.Required.Error("Email not valid"))),
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.Password, validation.Required, validation.Length(6, 999).Error("min 6 character")),
	)

	return err
}

func (userValidation *UserValidation) validationEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
