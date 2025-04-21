package models

import "time"

type AppointmentWithDetails struct {
	ID              int       `json:"id"`
	PatientID       int       `json:"patient_id"`
	PatientName     string    `json:"patient_name"`
	DoctorID        int       `json:"doctor_id"`
	DoctorName      string    `json:"doctor_name"`
	ScheduleID      int       `json:"schedule_id"`
	AppointmentTime time.Time `json:"appointment_time"`
}
