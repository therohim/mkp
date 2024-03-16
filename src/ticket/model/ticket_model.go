package model

type AddTicketRequest struct {
	Name        string             `json:"name"`
	Duration    int                `json:"duration"`
	Genre       []string           `json:"genre"`
	Cover       string             `json:"cover"`
	Description string             `json:"description"`
	Studio      []AddStudioRequest `json:"studio"`
}

type AddStudioRequest struct {
	ID       string               `json:"id"`
	Bioskop  string               `json:"bioskop"`
	Studio   string               `json:"studio"`
	Address  string               `json:"address"`
	Schedule []AddScheduleRequest `json:"schedule"`
}

type AddScheduleRequest struct {
	ID   string   `json:"id"`
	Date string   `json:"date"`
	Time string   `json:"time"`
	Seat []string `json:"seat"`
}

type EditTicketRequest struct {
	Name        string              `json:"name"`
	Duration    int                 `json:"duration"`
	Genre       []string            `json:"genre"`
	Cover       string              `json:"cover"`
	Description string              `json:"description"`
	Studio      []EditStudioRequest `json:"studio"`
}

type EditStudioRequest struct {
	ID       string                `json:"id"`
	Bioskop  string                `json:"bioskop"`
	Studio   string                `json:"studio"`
	Address  string                `json:"address"`
	Schedule []EditScheduleRequest `json:"schedule"`
}

type EditScheduleRequest struct {
	ID   string            `json:"id"`
	Date string            `json:"date"`
	Time string            `json:"time"`
	Seat []EditSeatRequest `json:"seat"`
}

type EditSeatRequest struct {
	ID   string `json:"id"`
	Seat string `json:"seat"`
}

type TicketResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Duration    int64            `json:"duration"`
	Genre       []string         `json:"genre"`
	Cover       string           `json:"cover"`
	Description string           `json:"description"`
	Studio      []StudioResponse `json:"studio"`
}

type StudioResponse struct {
	ID       string             `json:"id"`
	Bioskop  string             `json:"bioskop"`
	Studio   string             `json:"studio"`
	Address  string             `json:"address"`
	Schedule []ScheduleResponse `json:"schedule"`
}

type ScheduleResponse struct {
	Date string         `json:"date"`
	Time string         `json:"time"`
	Seat []SeatResponse `json:"seat"`
}

type SeatResponse struct {
	ID    string `json:"id"`
	Seat  string `json:"seat"`
	IsBuy bool   `json:"is_buy"`
}

type ListRequest struct {
	Page  string `query:"page"`
	Limit string `query:"limit"`
}
