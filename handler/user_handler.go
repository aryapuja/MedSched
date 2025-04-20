package handler

import (
	"MedSched/models"
	"MedSched/services"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// @Summary Masuk sebagai user dan menghasilkan token JWT
// @Description Endpoint ini memungkinkan user yang terdaftar untuk login dengan username dan password, serta menerima token JWT untuk autentikasi.
// @Tags autentikasi
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Data login user"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON request body
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	token, err := services.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Ambil data user dari database
	user, err := services.GetUserByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Mengembalikan token JWT
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"full_name": user.FullName,
		"role":      user.Role,
	})
}

// @Summary Daftar user baru
// @Description Endpoint ini memungkinkan user baru untuk mendaftar dengan username, email, password, nama lengkap, jenis kelamin, dan tanggal lahir.
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "Data pendaftaran user"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if user.Role == "" {
		user.Role = "pasien"
	}

	// Validasi kolom diisi
	if user.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Full name harus diisi"})
		return
	}
	if user.Gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gender harus diisi"})
		return
	}
	if user.BirthDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Birth date harus diisi"})
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username harus diisi"})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email harus diisi"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password harus diisi"})
		return
	}

	// Validasi email format menggunakan regex
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email tidak valid"})
		return
	}

	// Validasi panjang password
	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password harus minimal 6 karakte"})
		return
	}

	// Cek apakah username atau email sudah terdaftar
	if err := services.RegisterUser(&user); err != nil {
		// Jika ada error saat pendaftaran, periksa jenis error
		if err.Error() == "username already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username sudah terdaftar"})
			return
		}
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
			return
		}
		// Log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error membuat user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dibuat"})
}

// @Summary Ambil data semua Users
// @Description Endpoint ini memungkinkan admin untuk melihat user.
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	// Panggil service untuk mengambil semua user
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Ambil data user berdasarkan parameter query
// @Description Endpoint ini memungkinkan admin untuk melihat user berdasarkan username atau role.
// @Accept json
// @Produce json
// @Param username query string false "Username user"
// @Param role query string false "Role user"
// @Success 200 {array} models.User "Daftar user"
// @Failure 400 {object} models.ErrorResponse "Permintaan tidak valid"
// @Failure 404 {object} models.ErrorResponse "user tidak ditemukan"
// @Failure 500 {object} models.ErrorResponse "Error server internal"
// @Router /users/by [get]
func GetUsersByQuery(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	role := c.DefaultQuery("role", "")

	// ambil data user berdasarkan query
	if username != "" {
		// Ambil user berdasarkan username
		users, err := services.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, users)
	} else if role != "" {
		// Ambil user berdasarkan role
		users, err := services.GetUserByRole(role)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User dengan role yang ditentukan tidak ditemukan"})
			return
		}
		if len(users) == 0 {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"Tidak ada user ditemukan dengan role: " + role})
			return
		}
		c.JSON(http.StatusOK, users)
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Username atau role harus disediakan"})
	}
}

// @Summary Update data user
// @Description Endpoint ini memungkinkan admin untuk mengupdate data user berdasarkan username atau ID user.
// @Accept json
// @Produce json
// @Param username query string false "Username user"
// @Param id query string false "ID User"
// @Param user body models.User true "Data pembaruan user"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/by [put]
func UpdateUser(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	id := c.DefaultQuery("id", "")

	// Cek apakah salah satu parameter query ada
	if username == "" && id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Username or ID must be provided"})
		return
	}

	// Periksa ID/Username
	var user *models.User
	var err error
	if username != "" {
		user, err = services.GetUserByUsername(username)
	} else if id != "" {
		user, err = services.GetUserByID(id)
	}

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"User tidak ditemukan"})
		return
	}

	// Bind data user untuk diupdate
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Invalid input"})
		return
	}

	// Update data user
	var updateErr error
	if username != "" {
		updateErr = services.UpdateUser(username, &updatedUser)
	} else if id != "" {
		updateErr = services.UpdateUser(id, &updatedUser)
	}

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error updating user"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User updated successfully"})
}

// @Summary Mengambil data user itu sendiri
// @Description Endpoint ini memungkinkan user melihat informasi user itu sendiri.
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users/me [get]
func GetMyProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user belum diverifikasi"})
		return
	}

	user, err := services.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update data user
// @Description Endpoint ini memungkinkan user mengubah informasi user itu sendiri.
// @Accept json
// @Produce json
// @Param user body models.User true "User update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/me [put]
func UpdateMyProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user belum diverifikasi"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Input tidak valid"})
		return
	}

	// Update data user
	if err := services.UpdateUser(username.(string), &user); err != nil {
		if err.Error() == "user tidak ditemukan" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error updating user"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User updated successfully"})
}

// @Summary Menghapus user
// @Description Endpoint ini memungkinkan admin untuk menghapus informasi user berdasaarkan Username/ID.
// @Accept json
// @Produce json
// @Param username query string false "Username userr"
// @Param id query string false "ID User"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/by [delete]
func DeleteUser(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	id := c.DefaultQuery("id", "")

	// Cek apakah salah satu parameter query ada
	if username == "" && id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"username / ID harus disediakan"})
		return
	}

	// Periksa ID/Username
	var user *models.User
	var err error
	if username != "" {
		user, err = services.GetUserByUsername(username)
	} else if id != "" {
		user, err = services.GetUserByID(id)
	}

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"User tidak ditemukan"})
		return
	}

	if username != "" {
		// Hapus berdasarkan username
		err = services.DeleteUser(username)
	} else if id != "" {
		// Hapus berdasarkan ID
		err = services.DeleteUser(id)
	}

	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error menghapus data user"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Data user berhasil dihapus"})
}
