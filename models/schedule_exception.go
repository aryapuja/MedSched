package models

type ScheduleException struct {
	ID            int    `json:"id"`
	ScheduleID    int    `json:"schedule_id"`
	ExceptionDate string `json:"exception_date"`
}
