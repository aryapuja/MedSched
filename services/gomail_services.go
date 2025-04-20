package services

import (
	"MedSched/config"
	"MedSched/models"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

func SendEmail(recipientEmail, subject, body string) error {
	m := gomail.NewMessage()

	// Set email pengirim dan penerima
	m.SetHeader("From", "frozenforce4321@gmail.com") // email pengirim
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Set konfigurasi server SMTP
	d := gomail.NewDialer("smtp.gmail.com", 587, "frozenforce4321@gmail.com", "wrlw wwvg rgav yrnr")

	// Kirim email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gagal mengirim email: %v", err)
	}

	return nil
}

func LogNotification(userID int, appointmentID int, notificationType string, status string) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `
        INSERT INTO Notifications (user_id, appointment_id, notification_type, status, sent_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = db.Exec(sqlStatement, userID, appointmentID, notificationType, status, time.Now())
	if err != nil {
		return fmt.Errorf("gagal mencatat pengiriman notifikasi: %v", err)
	}

	return nil
}

// SendReminderEmail - Mengirimkan email pengingat H-1
func SendReminderEmail(patientEmail, appointmentTime string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "frozenforce4321@gmail.com")
	m.SetHeader("To", patientEmail)
	m.SetHeader("Subject", "Pengingat Jadwal Kontrol")
	m.SetBody("text/html", fmt.Sprintf("<h1>Reminder</h1><p>Jadwal kontrol Anda dilakukan pada: %s</p>", appointmentTime))

	d := gomail.NewDialer("smtp.gmail.com", 587, "frozenforce4321@gmail.com", "wrlw wwvg rgav yrnr")
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gagal mengirimkan email pengingat: %v", err)
	}
	return nil
}

// GetAppointmentsByDate - Mendapatkan janji berdasarkan tanggal
func GetAppointmentsByDate(date string) ([]models.Appointment, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, patient_id, schedule_id, appointment_time FROM appointments WHERE DATE(appointment_time) = $1`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(&appointment.ID, &appointment.PatientID, &appointment.ScheduleID, &appointment.AppointmentTime); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

// SendReminderEmails - Mengirim pengingat H-1 untuk semua janji
func SendReminderEmails() error {
	tomorrow := time.Now().Add(24 * time.Hour)
	appointmentDate := tomorrow.Format("2006-01-02")

	appointments, err := GetAppointmentsByDate(appointmentDate)
	if err != nil {
		return err
	}

	for _, appointment := range appointments {
		patientEmail, err := GetPatientEmailByID(appointment.PatientID)
		if err != nil {
			fmt.Println("Gagal mengambil email pasien:", err)
			continue
		}

		err = SendReminderEmail(patientEmail, appointment.AppointmentTime)
		if err != nil {
			fmt.Println("Gagal mengirimkan pengingat:", err)
		}
	}

	return nil
}

// StartReminderScheduler - Menjadwalkan pengiriman pengingat H-1
func StartReminderScheduler() {
	c := cron.New()

	_, err := c.AddFunc("0 9 * * *", func() { // Setiap jam 9 pagi
		if err := SendReminderEmails(); err != nil {
			fmt.Println("Gagal mengirimkan email pengingat:", err)
		}
	})

	if err != nil {
		fmt.Println("Gagal menjadwalkan cron job pengingat:", err)
		return
	}

	c.Start()
}
