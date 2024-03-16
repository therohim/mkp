package repository

import (
	"context"
	"database/sql"
	"strings"
	"test-mkp/config"
	"test-mkp/src/ticket/entity"
	"test-mkp/utils"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type ticketRepositoryImpl struct {
	db     *sqlx.DB
	tx     *sqlx.Tx
	withTx bool
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &ticketRepositoryImpl{
		db:     db,
		tx:     nil,
		withTx: false,
	}
}

func (repo *ticketRepositoryImpl) WithTx(ctx context.Context) (TicketRepository, error) {
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return repo, err
	}

	return &ticketRepositoryImpl{
		db:     repo.db,
		tx:     tx,
		withTx: true,
	}, nil
}

func (repo *ticketRepositoryImpl) Tx() *sqlx.Tx {
	return repo.tx
}

func (t *ticketRepositoryImpl) GetTicket(params entity.ListRequest) (resp []entity.Ticket, pagination entity.Pagination, err error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	resp = []entity.Ticket{}
	pagination = entity.Pagination{}

	statement := sq.Select().
		From("tickets a").
		Where(sq.Eq{"a.is_active": true, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	//prepare for pagination
	pageStatement := statement.Columns("COUNT(1) as pagination_total")
	pageSql, pageArgs, err := pageStatement.ToSql()
	if err != nil {
		return
	}

	err = t.db.GetContext(ctx, &pagination, pageSql, pageArgs...)
	if err != nil {
		return
	}

	pagination.Page = params.Page
	pagination.Limit = params.Limit

	page := (params.Page - 1) * params.Limit
	statement = statement.Limit(uint64(params.Limit)).Offset(uint64(page))

	statement = statement.Columns("id, name, duration, genre, cover, description, rating, created_at, updated_at, deleted_at")
	statement = statement.OrderBy("a.created_at desc")
	sql, args, err := statement.ToSql()
	if err != nil {
		return
	}

	rows, err := t.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.Ticket
		err := rows.StructScan(&item)
		if err != nil {
			return resp, pagination, err
		}

		resp = append(resp, item)
	}

	return resp, pagination, nil
}

func (t *ticketRepositoryImpl) GetStudioByTicketID(ticketID string) (resp []entity.Studio, err error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	resp = []entity.Studio{}

	statement := sq.Select().
		From("ticket_studio a").
		Where(sq.Eq{"a.is_active": true, "a.ticket_id": ticketID, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	statement = statement.Columns("id, bioskop, address, studio, created_at, updated_at, deleted_at")
	statement = statement.OrderBy("a.created_at desc")
	sql, args, err := statement.ToSql()
	if err != nil {
		return
	}

	rows, err := t.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.Studio
		err := rows.StructScan(&item)
		if err != nil {
			return resp, err
		}

		resp = append(resp, item)
	}

	return resp, nil
}

func (t *ticketRepositoryImpl) GetScheduleByStudioID(studioID string) (resp []entity.Schedule, err error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	resp = []entity.Schedule{}

	statement := sq.Select().
		From("ticket_schedule a").
		Where(sq.Eq{"a.studio_id": studioID, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	statement = statement.Columns("id, date, time, created_at, updated_at, deleted_at")
	statement = statement.OrderBy("a.created_at desc")
	sql, args, err := statement.ToSql()
	if err != nil {
		return
	}

	rows, err := t.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.Schedule
		err := rows.StructScan(&item)
		if err != nil {
			return resp, err
		}

		resp = append(resp, item)
	}

	return resp, nil
}

func (t *ticketRepositoryImpl) GetSeatByScheduleID(scheduleID string) (resp []entity.Seat, err error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	resp = []entity.Seat{}

	statement := sq.Select().
		From("ticket_seats a").
		Where(sq.Eq{"a.schedule_id": scheduleID, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	statement = statement.Columns("id, seat, is_buy, created_at, updated_at, deleted_at")
	statement = statement.OrderBy("a.created_at desc")
	sql, args, err := statement.ToSql()
	if err != nil {
		return
	}

	rows, err := t.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.Seat
		err := rows.StructScan(&item)
		if err != nil {
			return resp, err
		}

		resp = append(resp, item)
	}

	return resp, nil
}

func (r *ticketRepositoryImpl) FindTicketByID(ticketID string) (*entity.Ticket, error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	q := sq.Select().
		Columns("a.id", "a.name", "a.duration", "genre", "cover", "description", "rating", "is_active", "created_at", "updated_at", "deleted_at").
		From("tickets a").
		Where(sq.Eq{"a.id": ticketID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	response := entity.Ticket{}
	if err := r.db.GetContext(ctx, &response, sql, args...); err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return nil, nil
		}

		return &response, err
	}

	return &response, nil
}

func (r *ticketRepositoryImpl) SaveTicket(u entity.Ticket) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Insert("tickets").
		Columns("id", "name", "duration", "genre", "cover", "description", "is_active", "created_at", "updated_at").
		Values(u.ID, u.Name, u.Duration, u.Genre, u.Cover, u.Description, u.IsActive, u.CreatedAt, u.UpdatedAt).
		Suffix(`ON CONFLICT(id) DO UPDATE SET 
			name = EXCLUDED.name, 
			duration = EXCLUDED.duration, 
			genre = EXCLUDED.genre, 
			description = EXCLUDED.description,
			updated_at = now()
		`).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := statement.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) DeleteTicketByID(ticketID string) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Update("tickets").
		Set("is_active", false).
		Set("deleted_at", utils.JakartaTime(time.Now())).
		Set("updated_at", utils.JakartaTime(time.Now())).
		Where(sq.Eq{
			"id": ticketID,
		}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) SaveStudio(studio ...entity.Studio) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Insert("ticket_studio").
		Columns("id", "ticket_id", "bioskop", "studio", "address", "is_active", "created_at", "updated_at", "deleted_at").
		PlaceholderFormat(sq.Dollar)

	for _, v := range studio {
		statement = statement.Values(v.ID, v.TicketID, v.Bioskop, v.Studio, v.Address, v.IsActive, v.CreatedAt, v.UpdatedAt, v.DeletedAt)
	}
	statement = statement.Suffix(`ON CONFLICT(id) DO UPDATE SET 
		bioskop = EXCLUDED.bioskop, 
		address = EXCLUDED.address, 
		is_active = EXCLUDED.is_active,
		studio = EXCLUDED.studio,
		deleted_at = EXCLUDED.deleted_at,
		updated_at = now()
	`)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) InactiveStudioByTicketID(ticketID string) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Update("ticket_studio").
		Set("is_active", false).
		Set("deleted_at", utils.JakartaTime(time.Now())).
		Set("updated_at", utils.JakartaTime(time.Now())).
		Where(sq.Eq{
			"ticket_id": ticketID,
		}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) SaveSchedule(schedule ...entity.Schedule) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Insert("ticket_schedule").
		Columns("id", "studio_id", "date", "time", "created_at", "updated_at", "deleted_at").
		PlaceholderFormat(sq.Dollar)

	for _, v := range schedule {
		statement = statement.Values(v.ID, v.StudioID, v.Date, v.Time, v.CreatedAt, v.UpdatedAt, v.DeletedAt)
	}
	statement = statement.Suffix(`ON CONFLICT(id) DO UPDATE SET 
		date = EXCLUDED.date, 
		time = EXCLUDED.time,
		deleted_at = EXCLUDED.deleted_at,
		updated_at = now()
	`)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) InactiveScheduleByStudioID(studioID string) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Update("ticket_schedule").
		Set("deleted_at", utils.JakartaTime(time.Now())).
		Set("updated_at", utils.JakartaTime(time.Now())).
		Where(sq.Eq{
			"studio_id": studioID,
		}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) SaveSeat(seats ...entity.Seat) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Insert("ticket_seats").
		Columns("id", "schedule_id", "seat", "is_buy", "created_at", "updated_at", "deleted_at").
		PlaceholderFormat(sq.Dollar)

	for _, v := range seats {
		statement = statement.Values(v.ID, v.ScheduleID, v.Seat, v.IsBuy, v.CreatedAt, v.UpdatedAt, v.DeletedAt)
	}
	statement = statement.Suffix(`ON CONFLICT(id) DO UPDATE SET 
		seat = EXCLUDED.seat, 
		deleted_at = EXCLUDED.deleted_at,
		updated_at = now()
	`)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (r *ticketRepositoryImpl) InactiveSeatByScheduleID(scheduleId string) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Update("ticket_seats").
		Set("deleted_at", utils.JakartaTime(time.Now())).
		Set("updated_at", utils.JakartaTime(time.Now())).
		Where(sq.Eq{
			"schedule_id": scheduleId,
		}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if r.withTx {
		_, err = r.tx.ExecContext(ctx, sql, args...)
	} else {
		_, err = r.db.ExecContext(ctx, sql, args...)
	}

	return err
}
