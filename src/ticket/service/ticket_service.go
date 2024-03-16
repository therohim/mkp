package service

import (
	"test-mkp/src/ticket/model"
	"test-mkp/utils"
)

type TicketService interface {
	AddTicket(request model.AddTicketRequest) error
	DeleteTicket(ticketId string) error
	EditTicket(ticketId string, request model.EditTicketRequest) error
	GetTicket(request model.ListRequest) (response []model.TicketResponse, paginate utils.PaginationResponse, err error)
}
