package models

type AppError string

func (e AppError) Error() string { return string(e) }

const (
	ErrNotFound AppError = "not found"
)

type SubJSON struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserUUID    string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type SubRespJSON struct {
	ID int `json:"id"`
}

type CostSumSubRespJSON struct {
	Sum int `json:"sum"`
}
