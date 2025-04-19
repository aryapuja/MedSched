package services

import (
	"MedSched/config"
	"MedSched/models"
	"fmt"
)

// CreateAppointment - Menyimpan janji kontrol ke dalam database
func CreateAppointment(appointment *models.Appointment) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Menyimpan data janji kontrol ke dalam database
	sqlStatement := `INSERT INTO appointments (patient_id, doctor_id, schedule_id, appointment_time) 
                     VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, appointment.PatientID, appointment.DoctorID, appointment.ScheduleID, appointment.AppointmentTime)
	if err != nil {
		return fmt.Errorf("Error inserting appointment into database: %v", err)
	}

	return nil
}

// CountAppointments - Menghitung jumlah janji pada jadwal tertentu
func CountAppointments(scheduleID int) (int, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var count int
	// Menghitung jumlah janji yang terdaftar pada schedule_id yang sama
	sqlStatement := `SELECT COUNT(*) FROM appointments WHERE schedule_id=$1`
	err = db.QueryRow(sqlStatement, scheduleID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error counting appointments: %v", err)
	}

	return count, nil
}
