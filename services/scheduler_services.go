// CreateSchedule - Menambahkan jadwal dokter ke dalam database
package services

import (
	"MedSched/config"
	"MedSched/models"
	"database/sql"
	"fmt"
	"time"
)

func CreateSchedule(schedule *models.Schedule) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// Menyimpan data jadwal ke dalam database dengan available_slot
	sqlStatement := `INSERT INTO schedules (doctor_id, day_of_week, start_time, end_time, available_slot) 
                     VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, schedule.DoctorID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.AvailableSlot)
	if err != nil {
		return fmt.Errorf("Gagal menambahkan jadwal dokter ke dalam database: %v", err)
	}

	return nil
}

// GetAllSchedules - Mengambil semua jadwal dokter
func GetAllSchedules() ([]models.Schedule, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	var schedules []models.Schedule
	rows, err := db.Query(`SELECT id, doctor_id, day_of_week, start_time, end_time, available_slot, created_at FROM schedules`)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil data jadwal: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.DoctorID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.AvailableSlot, &schedule.CreatedAt); err != nil {
			return nil, fmt.Errorf("Gagal memindai data jadwal: %v", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// GetDoctorSchedules - Mengambil jadwal dokter (nama dokter, spesialisasi, hari, jumlah pasien yang bisa diterima, dan jam praktek)
func GetDoctorSchedules() ([]models.DoctorSchedule, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// Query untuk mengambil data jadwal dokter
	sqlStatement := `
        SELECT 
            d.name AS doctor_name, 
            d.specialization AS specialization, 
            s.day_of_week AS day, 
            s.available_slot AS available_slots, 
            s.start_time AS start_time, 
            s.end_time AS end_time
        FROM doctors d
        JOIN schedules s ON d.id = s.doctor_id
    `

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil data jadwal dokter: %v", err)
	}
	defer rows.Close()

	var doctorSchedules []models.DoctorSchedule

	for rows.Next() {
		var schedule models.DoctorSchedule
		err := rows.Scan(&schedule.DoctorName, &schedule.Specialization, &schedule.Day, &schedule.AvailableSlots, &schedule.StartTime, &schedule.EndTime)
		if err != nil {
			return nil, fmt.Errorf("Gagal memindai jadwal dokter: %v", err)
		}
		doctorSchedules = append(doctorSchedules, schedule)
	}

	return doctorSchedules, nil
}

// GetScheduleByID - Mencari jadwal dokter berdasarkan ID
func GetScheduleByID(id string) (*models.Schedule, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	var schedule models.Schedule
	sqlStatement := `SELECT id, doctor_id, day_of_week, start_time, end_time, available_slot, created_at FROM schedules WHERE id=$1`
	err = db.QueryRow(sqlStatement, id).Scan(&schedule.ID, &schedule.DoctorID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.AvailableSlot, &schedule.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Jadwal tidak ditemukan")
		}
		return nil, err
	}

	return &schedule, nil
}

// GetDoctorIDFromSchedule - Mencari doctor_id berdasarkan schedule_id
func GetDoctorIDFromSchedule(scheduleID int) (int, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var doctorID int
	sqlStatement := `SELECT doctor_id FROM schedules WHERE id=$1`
	err = db.QueryRow(sqlStatement, scheduleID).Scan(&doctorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("schedule not found")
		}
		return 0, err
	}

	return doctorID, nil
}

// UpdateSchedule - Memperbarui jadwal dokter berdasarkan ID
func UpdateSchedule(id string, updatedSchedule *models.Schedule) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `UPDATE schedules SET doctor_id=$1, day_of_week=$2, start_time=$3, end_time=$4, available_slot=$5 WHERE id=$6`
	_, err = db.Exec(sqlStatement, updatedSchedule.DoctorID, updatedSchedule.DayOfWeek, updatedSchedule.StartTime, updatedSchedule.EndTime, updatedSchedule.AvailableSlot, id)
	if err != nil {
		return fmt.Errorf("Gagal memperbarui jadwal: %v", err)
	}

	return nil
}

// DeleteSchedule - Menghapus jadwal dokter berdasarkan ID
func DeleteSchedule(id string) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `DELETE FROM schedules WHERE id=$1`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("Gagal menghapus jadwal: %v", err)
	}

	return nil
}

// CreateScheduleException - Menambahkan pengecualian untuk jadwal dokter
func CreateScheduleException(exception *models.ScheduleException) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// Menyimpan pengecualian jadwal ke dalam database
	sqlStatement := `INSERT INTO schedule_exceptions (schedule_id, exception_date) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, exception.ScheduleID, exception.ExceptionDate)
	if err != nil {
		return fmt.Errorf("Gagal menambahkan pengecualian jadwal: %v", err)
	}

	return nil
}

// GetAllScheduleExceptions - Mengambil semua pengecualian jadwal dokter
func GetAllScheduleExceptions() ([]models.ScheduleException, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exceptions []models.ScheduleException
	rows, err := db.Query(`SELECT id, schedule_id, exception_date FROM schedule_exceptions`)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil pengecualian jadwal: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exception models.ScheduleException
		if err := rows.Scan(&exception.ID, &exception.ScheduleID, &exception.ExceptionDate); err != nil {
			return nil, fmt.Errorf("Gagal memeriksa pengecualian jadwal: %v", err)
		}
		exceptions = append(exceptions, exception)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exceptions, nil
}

// GetScheduleExceptionByID - Mencari pengecualian jadwal berdasarkan ID
func GetScheduleExceptionByID(id string) (*models.ScheduleException, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exception models.ScheduleException
	sqlStatement := `SELECT id, schedule_id, exception_date FROM schedule_exceptions WHERE id=$1`
	err = db.QueryRow(sqlStatement, id).Scan(&exception.ID, &exception.ScheduleID, &exception.ExceptionDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pengecualian jadwal tidak ditemukan")
		}
		return nil, err
	}

	return &exception, nil
}

// CheckScheduleException - Mengecek apakah ada pengecualian untuk jadwal dokter
func CheckScheduleException(scheduleID int, appointmentTime time.Time) (*models.ScheduleException, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	var exception models.ScheduleException
	sqlStatement := `SELECT id, schedule_id, exception_date FROM schedule_exceptions WHERE schedule_id=$1 AND exception_date=$2`
	err = db.QueryRow(sqlStatement, scheduleID, appointmentTime.Format("2006-01-02")).Scan(&exception.ID, &exception.ScheduleID, &exception.ExceptionDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Tidak ada pengecualian
		}
		return nil, fmt.Errorf("Gagal memeriksa pengecualian jadwal: %v", err)
	}

	return &exception, nil
}

// DeleteScheduleException - Menghapus pengecualian jadwal dokter berdasarkan ID
func DeleteScheduleException(id string) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	sqlStatement := `DELETE FROM schedule_exceptions WHERE id=$1`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("Gagal menghapus pengecualian jadwal: %v", err)
	}

	return nil
}
