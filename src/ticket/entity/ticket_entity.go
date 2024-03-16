package entity

type Ticket struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Duration    int64   `db:"duration"`
	Genre       string  `db:"genre"`
	Cover       string  `db:"cover"`
	Description string  `db:"description"`
	Rating      float32 `db:"rating"`
	IsActive    bool    `db:"is_active"`
	Default
}

type Studio struct {
	ID       string `db:"id"`
	TicketID string `db:"ticket_id"`
	Studio   string `db:"studio"`
	Bioskop  string `db:"bioskop"`
	Address  string `db:"address"`
	IsActive bool   `db:"is_active"`
	Default
}

type Schedule struct {
	ID       string `db:"id"`
	StudioID string `db:"studio_id"`
	Date     string `db:"date"`
	Time     string `db:"time"`
	Default
}

type Seat struct {
	ID         string `db:"id"`
	ScheduleID string `db:"schedule_id"`
	Seat       string `db:"seat"`
	IsBuy      bool   `db:"is_buy"`
	Default
}
