package handler

import (
	"MedSched/models"
	"MedSched/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Membuat jadwal untuk dokter
// @Description Endpoint ini memungkinkan admin untuk membuat jadwal untuk dokter.
// @Tags jadwal
// @Accept json
// @Produce json
// @Param schedule body models.Schedule true "Data jadwal dokter"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules [post]
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Periksa doctor_id valid
	if schedule.DoctorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID dokter tidak valid"})
		return
	}

	// Simpan jadwal ke database
	if err := services.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jadwal berhasil dibuat"})
}

// @Summary Mengambil semua jadwal dokter
// @Description Endpoint ini memungkinkan admin untuk mengambil semua jadwal dokter.
// @Tags jadwal
// @Accept json
// @Produce json
// @Success 200 {array} models.Schedule
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules [get]
func GetAllSchedules(c *gin.Context) {
	schedules, err := services.GetAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal mengambil jadwal"})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// @Summary Mengambil jadwal berdasarkan ID
// @Description Endpoint ini memungkinkan admin untuk mengambil jadwal dokter berdasarkan ID.
// @Tags jadwal
// @Accept json
// @Produce json
// @Param id query string true "ID jadwal"
// @Success 200 {object} models.Schedule
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/by [get]
func GetScheduleByID(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	schedule, err := services.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Jadwal tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// @Summary Memperbarui jadwal
// @Description Endpoint ini memungkinkan admin untuk memperbarui jadwal dokter berdasarkan ID.
// @Tags jadwal
// @Accept json
// @Produce json
// @Param id query string true "ID jadwal"
// @Param schedule body models.Schedule true "Data pembaruan jadwal"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/by [put]
func UpdateSchedule(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	schedule, err := services.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Jadwal tidak ditemukan"})
		return
	}

	var updatedSchedule models.Schedule
	if err := c.ShouldBindJSON(&updatedSchedule); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Input tidak valid"})
		return
	}

	err = services.UpdateSchedule(id, &updatedSchedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal memperbarui jadwal"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Jadwal berhasil diperbarui"})
}

// @Summary Menghapus jadwal
// @Description Endpoint ini memungkinkan admin untuk menghapus jadwal berdasarkan ID.
// @Tags jadwal
// @Accept json
// @Produce json
// @Param id query string true "ID jadwal"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/by [delete]
func DeleteSchedule(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	schedule, err := services.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Jadwal tidak ditemukan"})
		return
	}

	err = services.DeleteSchedule(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal menghapus jadwal"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Jadwal berhasil dihapus"})
}

// @Summary Membuat pengecualian jadwal
// @Description Endpoint ini memungkinkan admin untuk menambahkan tanggal pengecualian untuk jadwal.
// @Tags pengecualian jadwal
// @Accept json
// @Produce json
// @Param exception body models.ScheduleException true "Data pengecualian jadwal"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedule_exceptions [post]
func CreateScheduleException(c *gin.Context) {
	var exception models.ScheduleException
	if err := c.ShouldBindJSON(&exception); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Input tidak valid"})
		return
	}

	schedule, err := services.GetScheduleByID(fmt.Sprint(exception.ScheduleID))
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Jadwal tidak ditemukan"})
		return
	}

	if err := services.CreateScheduleException(&exception); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal membuat pengecualian jadwal"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Pengecualian jadwal berhasil dibuat"})
}

// @Summary Mengambil semua pengecualian jadwal
// @Description Endpoint ini memungkinkan admin untuk mengambil semua pengecualian jadwal yang ada.
// @Tags pengecualian jadwal
// @Accept json
// @Produce json
// @Success 200 {array} models.ScheduleException
// @Failure 500 {object} models.ErrorResponse
// @Router /schedule_exceptions [get]
func GetAllScheduleExceptions(c *gin.Context) {
	exceptions, err := services.GetAllScheduleExceptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal mengambil pengecualian jadwal"})
		return
	}

	c.JSON(http.StatusOK, exceptions) // Mengirim data pengecualian jadwal
}

// @Summary Menghapus pengecualian jadwal
// @Description Endpoint ini memungkinkan admin untuk menghapus pengecualian jadwal berdasarkan ID.
// @Tags pengecualian jadwal
// @Accept json
// @Produce json
// @Param id query string true "ID pengecualian jadwal"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedule_exceptions [delete]
func DeleteScheduleException(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	exception, err := services.GetScheduleExceptionByID(id)
	if err != nil || exception == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Pengecualian jadwal tidak ditemukan"})
		return
	}

	err = services.DeleteScheduleException(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal menghapus pengecualian jadwal"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Pengecualian jadwal berhasil dihapus"})
}

// @Summary Mengambil jadwal dokter
// @Description Endpoint ini memungkinkan semua orang melihat semua jadwal dokter (nama dokter, spesialisasi, slot yang tersedia, dan jam praktek).
// @Tags dokter
// @Accept json
// @Produce json
// @Success 200 {array} models.DoctorSchedule
// @Failure 500 {object} models.ErrorResponse
// @Router /doctor_schedule [get]
func GetDoctorSchedule(c *gin.Context) {
	doctorSchedules, err := services.GetDoctorSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jadwal dokter"})
		return
	}

	c.JSON(http.StatusOK, doctorSchedules)
}
