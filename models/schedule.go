package models

import "time"

type Schedule struct {
	ID            int       `json:"id"`
	DoctorID      int       `json:"doctor_id"`
	DayOfWeek     string    `json:"day_of_week"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	AvailableSlot int       `json:"available_slot"`
	CreatedAt     time.Time `json:"created_at"`
}
