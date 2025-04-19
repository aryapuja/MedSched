package handler

import (
	"MedSched/models"
	"MedSched/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Create an appointment for a patient
// @Description This endpoint allows a patient or admin to create an appointment with a doctor.
// @Accept json
// @Produce json
// @Param appointment body models.Appointment true "Appointment data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /appointments [post]
func CreateAppointment(c *gin.Context) {
	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Log nilai user_id dari context
	fmt.Printf("User ID from context: %d\n", c.GetInt("userID")) // Debugging log

	// Jika request oleh pasien (user), set patient_id otomatis
	if c.GetString("role") == "pasien" {
		// Ambil user_id dari JWT token (misalnya "id" adalah claims di token)
		appointment.PatientID = c.GetInt("userID")
	} else if c.GetString("role") == "admin" {
		// Jika role admin, pastikan patient_id di body JSON dimasukkan
		if appointment.PatientID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "patient_id must be provided"})
			return
		}
	}

	// Cari doctor_id berdasarkan schedule_id
	doctorID, err := services.GetDoctorIDFromSchedule(appointment.ScheduleID)
	if err != nil || doctorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id or doctor not found"})
		return
	}
	appointment.DoctorID = doctorID

	// Log nilai appointment.PatientID untuk memeriksa apakah sudah terisi
	fmt.Printf("Received Appointment: PatientID=%d, ScheduleID=%d, AppointmentTime=%s\n",
		appointment.PatientID, appointment.ScheduleID, appointment.AppointmentTime)

	// Memastikan appointment_time berisi tanggal yang valid dan menambah waktu default
	parsedTime, err := time.Parse("2006-01-02", appointment.AppointmentTime) // Format YYYY-MM-DD
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Konversi parsedTime ke string setelah parsing
	appointment.AppointmentTime = parsedTime.Format("2006-01-02") // Mengonversi kembali ke string format YYYY-MM-DD

	// Cek apakah jadwal dokter yang dipilih masih tersedia
	schedule, err := services.GetScheduleByID(fmt.Sprint(appointment.ScheduleID))
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Schedule not found"})
		return
	}

	// Cek apakah slot tersedia
	if schedule.AvailableSlot <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"No available slots for this schedule"})
		return
	}

	// Cek apakah tanggal yang dipilih lebih kecil dari tanggal hari ini
	if parsedTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"The selected date is in the past"})
		return
	}

	// Cek apakah ada pengecualian pada jadwal ini (apakah tanggal kontrol dokter tersedia)
	exception, err := services.CheckScheduleException(appointment.ScheduleID, parsedTime)
	if err != nil || exception != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Doctor is not available on the selected date"})
		return
	}

	// Simpan janji kontrol ke database
	if err := services.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error creating appointment"})
		fmt.Printf("Received Appointment: %+v\n", err) // Log data appointment
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Appointment created successfully"})
}
