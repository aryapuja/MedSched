package handler

import (
	"MedSched/models"
	"MedSched/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Membuat dokter baru
// @Description Endpoint ini memungkinkan admin untuk membuat dokter baru.
// @Tags dokter
// @Accept json
// @Produce json
// @Param doctor body models.Doctor true "Data pendaftaran dokter"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors [post]
func CreateDoctor(c *gin.Context) {
	var doctor models.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := services.CreateDoctor(&doctor); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal membuat data dokter"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Data dokter berhasil dibuat"})
}

// @Summary Mengambil semua data dokter
// @Description Endpoint ini memungkinkan admin untuk mengambil semua dokter yang ada di sistem.
// @Tags dokter
// @Accept json
// @Produce json
// @Success 200 {array} models.Doctor
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors [get]
func GetAllDoctors(c *gin.Context) {
	doctors, err := services.GetAllDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal mengambil data dokter"})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

// @Summary Mengambil informasi dokter berdasarkan ID
// @Description Endpoint ini memungkinkan admin atau pengguna yang berwenang untuk melihat informasi dokter berdasarkan ID.
// @Tags dokter
// @Accept json
// @Produce json
// @Param id query string true "ID dokter"
// @Success 200 {object} models.Doctor
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors/by [get]
func GetDoctorByID(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	// Mengambil data dokter berdasarkan ID
	doctor, err := services.GetDoctorByID(id)
	if err != nil || doctor == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Data dokter tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

// @Summary Update data dokter
// @Description Endpoint ini memungkinkan admin untuk umengpdate data dokter berdasarkan ID.
// @Tags dokter
// @Accept json
// @Produce json
// @Param id query string true "ID dokter"
// @Param doctor body models.Doctor true "Data pembaruan dokter"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors/by [put]
func UpdateDoctor(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	// Periksa apakah dokter ada berdasarkan ID
	doctor, err := services.GetDoctorByID(id)
	if err != nil || doctor == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Data dokter tidak ditemukan"})
		return
	}

	// Bind data dokter untuk diupdate
	var updatedDoctor models.Doctor
	if err := c.ShouldBindJSON(&updatedDoctor); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Input tidak valid"})
		return
	}

	// Update dokter berdasarkan ID
	if err := services.UpdateDoctor(id, &updatedDoctor); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal memperbarui data dokter"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Data dokter berhasil diperbarui"})
}

// @Summary Menghapus dokter
// @Description Endpoint ini memungkinkan admin untuk menghapus dokter berdasarkan ID.
// @Tags dokter
// @Accept json
// @Produce json
// @Param id query string true "ID dokter"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /doctors/by [delete]
func DeleteDoctor(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"ID harus disediakan"})
		return
	}

	// Periksa apakah dokter ada berdasarkan ID
	doctor, err := services.GetDoctorByID(id)
	if err != nil || doctor == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"Data dokter tidak ditemukan"})
		return
	}

	// Hapus dokter berdasarkan ID
	if err := services.DeleteDoctor(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Gagal menghapus dokter"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Data dokter berhasil dihapus"})
}
