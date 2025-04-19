// CreateSchedule - Menambahkan jadwal dokter ke dalam database
package services

import (
	"MedSched/config"
	"MedSched/models"
	"database/sql"
	"fmt"
)

func CreateSchedule(schedule *models.Schedule) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Menyimpan data jadwal ke dalam database dengan available_slot
	sqlStatement := `INSERT INTO schedules (doctor_id, day_of_week, start_time, end_time, available_slot) 
                     VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, schedule.DoctorID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.AvailableSlot)
	if err != nil {
		return fmt.Errorf("Error inserting schedule into database: %v", err)
	}

	return nil
}

// GetAllSchedules - Mengambil semua jadwal dokter
func GetAllSchedules() ([]models.Schedule, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var schedules []models.Schedule
	rows, err := db.Query(`SELECT id, doctor_id, day_of_week, start_time, end_time, available_slot, created_at FROM schedules`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching schedules: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.DoctorID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.AvailableSlot, &schedule.CreatedAt); err != nil {
			return nil, fmt.Errorf("Error scanning schedule: %v", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// GetScheduleByID - Mencari jadwal dokter berdasarkan ID
func GetScheduleByID(id string) (*models.Schedule, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var schedule models.Schedule
	sqlStatement := `SELECT id, doctor_id, day_of_week, start_time, end_time, available_slot, created_at FROM schedules WHERE id=$1`
	err = db.QueryRow(sqlStatement, id).Scan(&schedule.ID, &schedule.DoctorID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.AvailableSlot, &schedule.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("schedule not found")
		}
		return nil, err
	}

	return &schedule, nil
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
		return fmt.Errorf("Error updating schedule: %v", err)
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
		return fmt.Errorf("Error deleting schedule: %v", err)
	}

	return nil
}

// CreateScheduleException - Menambahkan pengecualian untuk jadwal dokter
func CreateScheduleException(exception *models.ScheduleException) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Menyimpan pengecualian jadwal ke dalam database
	sqlStatement := `INSERT INTO schedule_exceptions (schedule_id, exception_date) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, exception.ScheduleID, exception.ExceptionDate)
	if err != nil {
		return fmt.Errorf("Error inserting schedule exception into database: %v", err)
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
		return nil, fmt.Errorf("Error fetching schedule exceptions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exception models.ScheduleException
		if err := rows.Scan(&exception.ID, &exception.ScheduleID, &exception.ExceptionDate); err != nil {
			return nil, fmt.Errorf("Error scanning schedule exception: %v", err)
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
			return nil, fmt.Errorf("schedule exception not found")
		}
		return nil, err
	}

	return &exception, nil
}

// DeleteScheduleException - Menghapus pengecualian jadwal dokter berdasarkan ID
func DeleteScheduleException(id string) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `DELETE FROM schedule_exceptions WHERE id=$1`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("Error deleting schedule exception: %v", err)
	}

	return nil
}
