package service

import (
	"errors"
	"strconv"
	"strings"
	"test-mkp/config"
	"test-mkp/src/ticket/entity"
	"test-mkp/src/ticket/model"
	"test-mkp/src/ticket/repository"
	"test-mkp/src/ticket/validation"
	"test-mkp/utils"
	"time"

	"github.com/google/uuid"
)

type ticketServiceImpl struct {
	repo     repository.TicketRepository
	validate validation.TicketValidation
}

func NewTicketService(
	repo repository.TicketRepository,
) TicketService {
	validate := new(validation.TicketValidation)
	return &ticketServiceImpl{
		repo:     repo,
		validate: *validate,
	}
}

func (s *ticketServiceImpl) GetTicket(request model.ListRequest) (response []model.TicketResponse, paginate utils.PaginationResponse, err error) {
	if request.Page == "" {
		request.Page = "1"
	}

	if request.Limit == "" {
		request.Limit = "10"
	}
	pageU64, err := strconv.ParseUint(request.Page, 10, 32)
	if err != nil {
		return
	}
	limitU64, err := strconv.ParseUint(request.Limit, 10, 32)
	if err != nil {
		return
	}
	params := entity.ListRequest{
		Page:  uint(pageU64),
		Limit: uint(limitU64),
	}

	tickets, page, err := s.repo.GetTicket(params)
	if err != nil {
		return
	}

	for _, v := range tickets {
		ticket := model.TicketResponse{
			ID:          v.ID,
			Name:        v.Name,
			Duration:    v.Duration,
			Genre:       strings.Split(v.Genre, ","),
			Cover:       v.Cover,
			Description: v.Description,
		}
		studio, err := s.repo.GetStudioByTicketID(v.ID)
		if err != nil {
			return response, paginate, err
		}
		ticket.Studio = []model.StudioResponse{}
		for _, v2 := range studio {
			std := model.StudioResponse{
				ID:       v2.ID,
				Bioskop:  v2.Bioskop,
				Studio:   v2.Studio,
				Address:  v2.Address,
				Schedule: []model.ScheduleResponse{},
			}

			schedule, err := s.repo.GetScheduleByStudioID(v2.ID)
			if err != nil {
				return response, paginate, err
			}
			for _, v3 := range schedule {
				scd := model.ScheduleResponse{
					Date: v3.Date,
					Time: v3.Time,
					Seat: []model.SeatResponse{},
				}
				seats, err := s.repo.GetSeatByScheduleID(v3.ID)
				if err != nil {
					return response, paginate, err
				}
				for _, v4 := range seats {
					scd.Seat = append(scd.Seat, model.SeatResponse{
						ID:    v4.ID,
						Seat:  v4.Seat,
						IsBuy: v4.IsBuy,
					})
				}
				std.Schedule = append(std.Schedule, scd)
			}
			ticket.Studio = append(ticket.Studio, std)
		}
		response = append(response, ticket)
	}

	paginate.Page = int16(pageU64)
	paginate.Limit = int16(limitU64)
	paginate.Total = uint32(page.Total)

	return
}

func (s *ticketServiceImpl) AddTicket(request model.AddTicketRequest) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	ticket := entity.Ticket{
		ID:          uuid.NewString(),
		Name:        request.Name,
		Duration:    int64(request.Duration),
		Description: request.Description,
		Cover:       request.Cover,
		IsActive:    true,
		Default: entity.Default{
			CreatedAt: utils.JakartaTime(time.Now()),
			UpdatedAt: utils.JakartaTime(time.Now()),
		},
	}
	for i, v := range request.Genre {
		ticket.Genre += v
		if i != len(request.Genre)-1 {
			ticket.Genre += ","
		}
	}

	studios := []entity.Studio{}
	schedules := []entity.Schedule{}
	seats := []entity.Seat{}
	for _, v := range request.Studio {
		studio := entity.Studio{
			ID:       uuid.NewString(),
			TicketID: ticket.ID,
			Bioskop:  v.Bioskop,
			Studio:   v.Studio,
			Address:  v.Address,
			IsActive: true,
			Default: entity.Default{
				CreatedAt: utils.JakartaTime(time.Now()),
				UpdatedAt: utils.JakartaTime(time.Now()),
				DeletedAt: nil,
			},
		}
		for _, v2 := range v.Schedule {
			schedule := entity.Schedule{
				ID:       uuid.NewString(),
				StudioID: studio.ID,
				Date:     v2.Date,
				Time:     v2.Time,
				Default: entity.Default{
					CreatedAt: utils.JakartaTime(time.Now()),
					UpdatedAt: utils.JakartaTime(time.Now()),
					DeletedAt: nil,
				},
			}

			for _, v3 := range v2.Seat {
				seat := entity.Seat{
					ID:         uuid.NewString(),
					ScheduleID: schedule.ID,
					Seat:       v3,
					IsBuy:      false,
					Default: entity.Default{
						CreatedAt: utils.JakartaTime(time.Now()),
						UpdatedAt: utils.JakartaTime(time.Now()),
						DeletedAt: nil,
					},
				}

				seats = append(seats, seat)
			}

			schedules = append(schedules, schedule)
		}
		studios = append(studios, studio)
	}
	tx, err := s.repo.WithTx(ctx)
	if err != nil {
		return err
	}

	if err := tx.SaveTicket(ticket); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.SaveStudio(studios...); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.SaveSchedule(schedules...); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.SaveSeat(seats...); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.Tx().Commit(); err != nil {
		tx.Tx().Rollback()
		return err
	}

	return nil
}

func (s *ticketServiceImpl) DeleteTicket(ticketId string) error {
	if ticketId == "" {
		return errors.New("ticket not found")
	}

	ticket, err := s.repo.FindTicketByID(ticketId)
	if err != nil {
		return err
	}

	if ticket == nil {
		return errors.New("ticket not found")
	}

	if !ticket.IsActive {
		return errors.New("ticket have been inactive")
	}

	if err := s.repo.DeleteTicketByID(ticketId); err != nil {
		return err
	}

	return nil
}

func (s *ticketServiceImpl) EditTicket(ticketId string, request model.EditTicketRequest) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	if ticketId == "" {
		return errors.New("ticket not found")
	}

	ticket, err := s.repo.FindTicketByID(ticketId)
	if err != nil {
		return err
	}

	if ticket == nil {
		return errors.New("ticket not found")
	}

	if !ticket.IsActive {
		return errors.New("ticket have been inactive")
	}

	ticketRequest := entity.Ticket{
		ID:          ticketId,
		Name:        request.Name,
		Duration:    int64(request.Duration),
		Description: request.Description,
		Cover:       request.Cover,
		IsActive:    true,
		Default: entity.Default{
			CreatedAt: utils.JakartaTime(time.Now()),
			UpdatedAt: utils.JakartaTime(time.Now()),
		},
	}
	for i, v := range request.Genre {
		ticket.Genre += v
		if i != len(request.Genre)-1 {
			ticket.Genre += ","
		}
	}

	studios := []entity.Studio{}
	schedules := []entity.Schedule{}
	seats := []entity.Seat{}
	for _, v := range request.Studio {
		studio := entity.Studio{
			ID:       v.ID,
			TicketID: ticket.ID,
			Bioskop:  v.Bioskop,
			Studio:   v.Studio,
			Address:  v.Address,
			IsActive: true,
			Default: entity.Default{
				CreatedAt: utils.JakartaTime(time.Now()),
				UpdatedAt: utils.JakartaTime(time.Now()),
				DeletedAt: nil,
			},
		}
		for _, v2 := range v.Schedule {
			schedule := entity.Schedule{
				ID:       v2.ID,
				StudioID: studio.ID,
				Date:     v2.Date,
				Time:     v2.Time,
				Default: entity.Default{
					CreatedAt: utils.JakartaTime(time.Now()),
					UpdatedAt: utils.JakartaTime(time.Now()),
					DeletedAt: nil,
				},
			}

			for _, v3 := range v2.Seat {
				seat := entity.Seat{
					ID:         v3.ID,
					ScheduleID: schedule.ID,
					Seat:       v3.Seat,
					IsBuy:      false,
					Default: entity.Default{
						CreatedAt: utils.JakartaTime(time.Now()),
						UpdatedAt: utils.JakartaTime(time.Now()),
						DeletedAt: nil,
					},
				}

				seats = append(seats, seat)
			}

			schedules = append(schedules, schedule)
		}
		studios = append(studios, studio)
	}
	tx, err := s.repo.WithTx(ctx)
	if err != nil {
		return err
	}

	if err := tx.SaveTicket(ticketRequest); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.SaveStudio(studios...); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.SaveSchedule(schedules...); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.SaveSeat(seats...); err != nil {
		tx.Tx().Rollback()
		return err
	}

	if err := tx.Tx().Commit(); err != nil {
		tx.Tx().Rollback()
		return err
	}

	return nil
}
