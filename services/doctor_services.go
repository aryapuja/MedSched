package services

import (
	"MedSched/config"
	"MedSched/models"
	"database/sql"
	"fmt"
)

// CreateDoctor - Menambahkan dokter baru ke database
func CreateDoctor(doctor *models.Doctor) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Menyimpan data dokter ke dalam database
	sqlStatement := `INSERT INTO doctors (name, specialization, contact_info) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, doctor.Name, doctor.Specialization, doctor.ContactInfo)
	if err != nil {
		return fmt.Errorf("Error inserting doctor into database: %v", err)
	}

	return nil
}

// GetAllDoctors - Mengambil semua data dokter
func GetAllDoctors() ([]models.Doctor, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}
	defer db.Close()

	var doctors []models.Doctor
	rows, err := db.Query(`SELECT id, name, specialization, contact_info FROM doctors`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching doctors: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var doctor models.Doctor
		if err := rows.Scan(&doctor.ID, &doctor.Name, &doctor.Specialization, &doctor.ContactInfo); err != nil {
			return nil, fmt.Errorf("Error scanning doctor: %v", err)
		}
		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over doctors: %v", err)
	}

	return doctors, nil
}

// GetDoctorByID - Get doctor data based on ID
func GetDoctorByID(id string) (*models.Doctor, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var doctor models.Doctor
	sqlStatement := `SELECT id, name, specialization, contact_info FROM doctors WHERE id=$1`
	err = db.QueryRow(sqlStatement, id).Scan(&doctor.ID, &doctor.Name, &doctor.Specialization, &doctor.ContactInfo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("doctor not found")
		}
		return nil, err
	}
	return &doctor, nil
}

// UpdateDoctor - Update doctor's data based on ID
func UpdateDoctor(id string, updatedDoctor *models.Doctor) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Update data dokter berdasarkan ID
	sqlStatement := `UPDATE doctors SET name=$1, specialization=$2, contact_info=$3 WHERE id=$4`
	_, err = db.Exec(sqlStatement, updatedDoctor.Name, updatedDoctor.Specialization, updatedDoctor.ContactInfo, id)
	if err != nil {
		return fmt.Errorf("Error updating doctor: %v", err)
	}

	return nil
}

// DeleteDoctor - Delete doctor's data based on ID
func DeleteDoctor(id string) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Hapus dokter berdasarkan ID
	sqlStatement := `DELETE FROM doctors WHERE id=$1`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("Error deleting doctor: %v", err)
	}

	return nil
}
