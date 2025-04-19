package handler

import (
	"MedSched/models"
	"MedSched/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create a schedule for a doctor
// @Description This endpoint allows an admin to create a schedule for a doctor.
// @Accept json
// @Produce json
// @Param schedule body models.Schedule true "Schedule data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules [post]
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Periksa doctor_id valid
	if schedule.DoctorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	// Simpan jadwal ke database
	if err := services.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully"})
}

// @Summary Get all schedules
// @Description This endpoint allows an admin to get all schedules for doctors.
// @Accept json
// @Produce json
// @Success 200 {array} models.Schedule
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules [get]
func GetAllSchedules(c *gin.Context) {
	schedules, err := services.GetAllSchedules() // Mendapatkan semua jadwal dari service
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error fetching schedules"})
		return
	}

	c.JSON(http.StatusOK, schedules) // Mengirim data jadwal dokter
}

// @Summary Get a schedule by ID
// @Description This endpoint allows an admin to get a schedule for a doctor by its ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the schedule"
// @Success 200 {object} models.Schedule
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/by [get]
func GetScheduleByID(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	schedule, err := services.GetScheduleByID(id) // Mengambil jadwal berdasarkan ID
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, schedule) // Mengirim data jadwal berdasarkan ID
}

// @Summary Update a schedule
// @Description This endpoint allows an admin to update a schedule for a doctor by ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the schedule"
// @Param schedule body models.Schedule true "Schedule update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/by [put]
func UpdateSchedule(c *gin.Context) {
	// Mendapatkan ID dari query parameter
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	// Validasi apakah jadwal dengan ID tersebut ada
	schedule, err := services.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Schedule not found"})
		return
	}

	// Binding data yang akan diperbarui
	var updatedSchedule models.Schedule
	if err := c.ShouldBindJSON(&updatedSchedule); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Invalid input"})
		return
	}

	// Memperbarui jadwal berdasarkan ID
	err = services.UpdateSchedule(id, &updatedSchedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error updating schedule"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Schedule updated successfully"})
}

// @Summary Delete a schedule
// @Description This endpoint allows an admin to delete a schedule by its ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the schedule"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/by [delete]
func DeleteSchedule(c *gin.Context) {
	// Mendapatkan ID dari query parameter
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	// Validasi apakah jadwal dengan ID tersebut ada
	schedule, err := services.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Schedule not found"})
		return
	}

	// Menghapus jadwal berdasarkan ID
	err = services.DeleteSchedule(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error deleting schedule"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Schedule deleted successfully"})
}

// @Summary Create an exception for a schedule
// @Description This endpoint allows an admin to add an exception date for a schedule.
// @Accept json
// @Produce json
// @Param exception body models.ScheduleException true "Schedule exception data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedule_exceptions [post]
func CreateScheduleException(c *gin.Context) {
	var exception models.ScheduleException
	if err := c.ShouldBindJSON(&exception); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Invalid input"})
		return
	}

	// Periksa apakah schedule_id valid (ada dalam tabel schedules)
	schedule, err := services.GetScheduleByID(fmt.Sprint(exception.ScheduleID))
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Schedule not found"})
		return
	}

	// Menyimpan pengecualian jadwal ke database
	if err := services.CreateScheduleException(&exception); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error creating schedule exception"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Schedule exception created successfully"})
}

// @Summary Get all schedule exceptions
// @Description This endpoint allows an admin to get all schedule exceptions.
// @Accept json
// @Produce json
// @Success 200 {array} models.ScheduleException
// @Failure 500 {object} models.ErrorResponse
// @Router /schedule_exceptions [get]
func GetAllScheduleExceptions(c *gin.Context) {
	exceptions, err := services.GetAllScheduleExceptions() // Mendapatkan semua pengecualian jadwal dari service
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error fetching schedule exceptions"})
		return
	}

	c.JSON(http.StatusOK, exceptions) // Mengirim data pengecualian jadwal
}

// @Summary Delete a schedule exception
// @Description This endpoint allows an admin to delete a schedule exception by its ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the schedule exception"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedule_exceptions [delete]
func DeleteScheduleException(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	// Periksa apakah pengecualian dengan ID tersebut ada di tabel schedule_exceptions
	exception, err := services.GetScheduleExceptionByID(id)
	if err != nil || exception == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Schedule exception not found"})
		return
	}

	// Menghapus pengecualian jadwal berdasarkan ID
	err = services.DeleteScheduleException(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error deleting schedule exception"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Schedule exception deleted successfully"})
}
