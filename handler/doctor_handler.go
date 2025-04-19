package handler

import (
	"MedSched/models"
	"MedSched/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new doctor
// @Description This endpoint allows an admin to create a new doctor.
// @Accept json
// @Produce json
// @Param doctor body models.Doctor true "Doctor registration data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors [post]
func CreateDoctor(c *gin.Context) {
	var doctor models.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Menyimpan data dokter ke database
	if err := services.CreateDoctor(&doctor); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error creating doctor"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Doctor created successfully"})
}

// @Summary Get all doctors
// @Description This endpoint allows an admin to get all doctors in the system.
// @Accept json
// @Produce json
// @Success 200 {array} models.Doctor
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors [get]
func GetAllDoctors(c *gin.Context) {
	// Mengambil daftar dokter dari service
	doctors, err := services.GetAllDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error fetching doctors"})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

// @Summary Get a doctor's information by ID
// @Description This endpoint allows an admin or authorized user to view a doctor's information by ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the doctor"
// @Success 200 {object} models.Doctor
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors/by [get]
func GetDoctorByID(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	// Mengambil data dokter berdasarkan ID
	doctor, err := services.GetDoctorByID(id)
	if err != nil || doctor == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Doctor not found"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

// @Summary Update a doctor's data
// @Description This endpoint allows an admin to update a doctor's data based on ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the doctor to update"
// @Param doctor body models.Doctor true "Doctor update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors/by [put]
func UpdateDoctor(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	// Periksa apakah dokter ada berdasarkan ID
	doctor, err := services.GetDoctorByID(id)
	if err != nil || doctor == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Doctor not found"})
		return
	}

	// Bind data dokter untuk diupdate
	var updatedDoctor models.Doctor
	if err := c.ShouldBindJSON(&updatedDoctor); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Invalid input"})
		return
	}

	// Update dokter berdasarkan ID
	if err := services.UpdateDoctor(id, &updatedDoctor); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error updating doctor"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Doctor updated successfully"})
}

// @Summary Delete a doctor
// @Description This endpoint allows an admin to delete a doctor based on ID.
// @Accept json
// @Produce json
// @Param id query string true "ID of the doctor to delete"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors/by [delete]
func DeleteDoctor(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID must be provided"})
		return
	}

	// Periksa apakah dokter ada berdasarkan ID
	doctor, err := services.GetDoctorByID(id)
	if err != nil || doctor == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Doctor not found"})
		return
	}

	// Hapus dokter berdasarkan ID
	if err := services.DeleteDoctor(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error deleting doctor"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Doctor deleted successfully"})
}

