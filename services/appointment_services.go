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

	sqlStatement := `INSERT INTO appointments (patient_id, doctor_id, schedule_id, appointment_time) 
                     VALUES ($1, $2, $3, $4) RETURNING id`

	err = db.QueryRow(sqlStatement, appointment.PatientID, appointment.DoctorID, appointment.ScheduleID, appointment.AppointmentTime).Scan(&appointment.ID)
	if err != nil {
		return fmt.Errorf("Gagal memasukkan data janji temu ke database: %v", err)
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
		return 0, fmt.Errorf("Gagal menghitung janji temu: %v", err)
	}

	return count, nil
}

// CheckIfAppointmentExists - Mengecek apakah pasien sudah memiliki janji pada schedule yang sama
func CheckIfAppointmentExists(patientID int, scheduleID int) (bool, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	sqlStatement := `SELECT COUNT(*) FROM appointments WHERE patient_id = $1 AND schedule_id = $2`
	err = db.QueryRow(sqlStatement, patientID, scheduleID).Scan(&count)
	if err != nil {
		return false, err
	}

	// Jika ada 1 atau lebih janji yang ditemukan, return true
	return count > 0, nil
}

// GetPatientEmailByID - Mengambil email pasien berdasarkan PatientID
func GetPatientEmailByID(patientID int) (string, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var email string
	sqlStatement := `SELECT email FROM users WHERE id = $1`
	err = db.QueryRow(sqlStatement, patientID).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("Gagal mengambil email pasien: %v", err)
	}

	return email, nil
}

// GetDoctorNameByID - Mengambil nama dokter berdasarkan doctor_id
func GetDoctorNameByID(doctorID int) (string, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var doctorName string
	sqlStatement := `SELECT name FROM doctors WHERE id = $1`
	err = db.QueryRow(sqlStatement, doctorID).Scan(&doctorName)
	if err != nil {
		return "", fmt.Errorf("Gagal mengambil nama dokter: %v", err)
	}

	return doctorName, nil
}

// GetAllAppointments - Mengambil semua data appointment
func GetAllAppointments() ([]models.AppointmentWithDetails, error) {
    db, err := config.SetupDatabase()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    var appointments []models.AppointmentWithDetails

    // Query untuk mengambil appointment dengan informasi pasien dan dokter
    query := `
        SELECT a.id, a.patient_id, u.username AS patient_name, a.doctor_id, d.name AS doctor_name, a.schedule_id, a.appointment_time
        FROM appointments a
        JOIN users u ON a.patient_id = u.id
        JOIN doctors d ON a.doctor_id = d.id
    `
    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("Gagal mengambil data appointment: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var appointment models.AppointmentWithDetails
        if err := rows.Scan(&appointment.ID, &appointment.PatientID, &appointment.PatientName, &appointment.DoctorID, &appointment.DoctorName, &appointment.ScheduleID, &appointment.AppointmentTime); err != nil {
            return nil, fmt.Errorf("Gagal memindai data appointment: %v", err)
        }
        appointments = append(appointments, appointment)
    }

    return appointments, nil
}

// GetAppointmentsWithDetail - Mengambil data appointment berdasarkan tanggal
func GetAppointmentsWithDetail(date string) ([]models.AppointmentWithDetails, error) {
    db, err := config.SetupDatabase()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    var appointments []models.AppointmentWithDetails

    // Query untuk mengambil appointment berdasarkan tanggal
    query := `
        SELECT a.id, a.patient_id, u.username AS patient_name, a.doctor_id, d.name AS doctor_name, a.schedule_id, a.appointment_time
        FROM appointments a
        JOIN users u ON a.patient_id = u.id
        JOIN doctors d ON a.doctor_id = d.id
        WHERE DATE(a.appointment_time) = $1
    `
    rows, err := db.Query(query, date)
    if err != nil {
        return nil, fmt.Errorf("Gagal mengambil data appointment berdasarkan tanggal: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var appointment models.AppointmentWithDetails
        if err := rows.Scan(&appointment.ID, &appointment.PatientID, &appointment.PatientName, &appointment.DoctorID, &appointment.DoctorName, &appointment.ScheduleID, &appointment.AppointmentTime); err != nil {
            return nil, fmt.Errorf("Gagal memindai data appointment: %v", err)
        }
        appointments = append(appointments, appointment)
    }

    return appointments, nil
}
