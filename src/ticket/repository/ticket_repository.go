package repository

import (
	"context"
	"test-mkp/src/ticket/entity"

	"github.com/jmoiron/sqlx"
)

type TicketRepository interface {
	FindTicketByID(ticketID string) (*entity.Ticket, error)
	SaveTicket(u entity.Ticket) error
	DeleteTicketByID(ticketID string) error
	SaveStudio(studio ...entity.Studio) error
	InactiveStudioByTicketID(ticketID string) error
	SaveSchedule(schedule ...entity.Schedule) error
	InactiveScheduleByStudioID(studioID string) error
	SaveSeat(seats ...entity.Seat) error
	InactiveSeatByScheduleID(scheduleId string) error

	GetTicket(params entity.ListRequest) (resp []entity.Ticket, pagination entity.Pagination, err error)
	GetStudioByTicketID(ticketID string) (resp []entity.Studio, err error)
	GetScheduleByStudioID(studioID string) (resp []entity.Schedule, err error)
	GetSeatByScheduleID(scheduleID string) (resp []entity.Seat, err error)

	// use to control database transaction from service
	WithTx(ctx context.Context) (TicketRepository, error)
	Tx() *sqlx.Tx
}
