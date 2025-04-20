package handler

import (
	"MedSched/models"
	"MedSched/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Membuat janji kontrol untuk pasien
// @Description Endpoint ini memungkinkan pasien atau admin untuk membuat janji kontrol dengan dokter.
// @Tags appointment
// @Accept json
// @Produce json
// @Param appointment body models.Appointment true "Data janji kontrol"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /appointments [post]
func CreateAppointment(c *gin.Context) {
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Jika request oleh pasien (user), set patient_id otomatis
	if c.GetString("role") == "pasien" {
		appointment.PatientID = c.GetInt("userID")
	} else if c.GetString("role") == "admin" {
		// Jika role admin, pastikan patient_id di body JSON dimasukkan
		if appointment.PatientID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "patient_id harus disediakan"})
			return
		}
	}

	// Cari doctor_id berdasarkan schedule_id
	doctorID, err := services.GetDoctorIDFromSchedule(appointment.ScheduleID)
	if err != nil || doctorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schedule_id tidak valid atau dokter tidak ditemukan"})
		return
	}
	appointment.DoctorID = doctorID

	// Memastikan appointment_time berisi tanggal yang valid
	parsedTime, err := time.Parse("2006-01-02", appointment.AppointmentTime) // Format YYYY-MM-DD
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid"})
		return
	}

	appointment.AppointmentTime = parsedTime.Format("2006-01-02") // Mengonversi kembali ke string format YYYY-MM-DD

	// Cek apakah tanggal yang dipilih lebih kecil dari tanggal hari ini
	if parsedTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Tanggal yang dipilih sudah lewat"})
		return
	}

	// Cek apakah ada pengecualian pada jadwal ini
	exception, err := services.CheckScheduleException(appointment.ScheduleID, parsedTime)
	if err != nil || exception != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Dokter tidak tersedia pada tanggal yang dipilih"})
		return
	}

	// Cek apakah pasien sudah mendaftar ke jadwal yang sama
	exists, err := services.CheckIfAppointmentExists(appointment.PatientID, appointment.ScheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa janji yang ada"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pasien sudah terdaftar di jadwal ini"})
		return
	}

	// Cek apakah jadwal dokter yang dipilih masih tersedia
	schedule, err := services.GetScheduleByID(fmt.Sprint(appointment.ScheduleID))
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Jadwal tidak ditemukan"})
		return
	}

	// Hitung jumlah janji yang sudah ada untuk schedule_id yang sama
	currentCount, err := services.CountAppointments(appointment.ScheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah janji"})
		return
	}

	// Periksa apakah jumlah janji yang ada sudah melebihi available_slot
	if currentCount >= schedule.AvailableSlot {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada slot yang tersedia untuk jadwal ini"})
		return
	}

	// Cek apakah slot tersedia
	if schedule.AvailableSlot <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Tidak ada slot yang tersedia untuk jadwal ini"})
		return
	}

	// Simpan janji kontrol ke database
	if err := services.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal membuat janji kontrol"})
		return
	}

	// Mengambil nama dokter berdasarkan doctor_id
	doctorName, err := services.GetDoctorNameByID(appointment.DoctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil nama dokter"})
		return
	}

	// Ambil email pasien berdasarkan PatientID
	patientEmail, err := services.GetPatientEmailByID(appointment.PatientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil email pasien"})
		return
	}
	subject := "Konfirmasi Jadwal Kontrol"
	body := fmt.Sprintf("<h1>Jadwal Kontrol</h1><p>Jadwal pertemuan Anda dengan %s dijadwalkan pada: %s</p>", appointment.AppointmentTime, doctorName)

	if err := services.SendEmail(patientEmail, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirimkan email konfirmasi"})
		return
	}

	// Log pengiriman email ke tabel Notifications
	notificationType := "Konfirmasi Jadwal Kontrol"
	status := "Terkirim"
	if err := services.LogNotification(appointment.PatientID, appointment.ID, notificationType, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat pengiriman notifikasi"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Janji kontrol berhasil dibuat"})
}
