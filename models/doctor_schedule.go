package models

import "time"

type DoctorSchedule struct {
	DoctorName     string    `json:"doctor_name"`
	Specialization string    `json:"specialization"`
	Day            string    `json:"day"`
	AvailableSlots int       `json:"available_slots"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
}
