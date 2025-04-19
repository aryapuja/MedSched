package models

import "time"

type Appointment struct {
	ID              int       `json:"id"`
	PatientID       int       `json:"patient_id"`
	DoctorID        int       `json:"doctor_id"`
	ScheduleID      int       `json:"schedule_id"`
	AppointmentTime string    `json:"appointment_time"`
	CreatedAt       time.Time `json:"created_at"`
}
