package validation

import "test-mkp/src/ticket/model"

type TicketValidation struct{}

func (v *TicketValidation) AddTicketValidate(request model.AddTicketRequest) error {
	// err := validation.ValidateStruct(
	// 	&request,
	// 	validation.Field(&request.ID, validation.Required),
	// 	validation.Field(&request.Price, validation.Required, validation.Min(1)),
	// 	validation.Field(&request.Qty, validation.Required, validation.Min(1)),
	// 	validation.Field(&request.Date, validation.Required),
	// 	validation.Field(&request.RefID, validation.Required),
	// 	validation.Field(&request.Signature, validation.Required),
	// )

	// return err
	return nil
}
