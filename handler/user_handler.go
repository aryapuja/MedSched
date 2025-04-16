package handler

import (
	"MedSched/models"
	"MedSched/services"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// @Summary Login a user and generate a JWT token
// @Description This endpoint allows a registered user to log in with their username and password, and receive a JWT token for authentication.
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "User login data"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := services.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Ambil data pengguna dari database
	user, err := services.GetUserByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Mengembalikan token JWT
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"full_name": user.FullName,
		"role":      user.Role,
	})
}

// @Summary Register a new user
// @Description This endpoint allows a new user to register with their username, email, password, full name, gender, and birth date.
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if user.Role == "" {
		user.Role = "pasien"
	}

	// Validasi kolom diisi
	if user.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Full name is required"})
		return
	}
	if user.Gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gender is required"})
		return
	}
	if user.BirthDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Birth date is required"})
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	// Validasi email format menggunakan regex
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validasi panjang password
	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	// Cek apakah username atau email sudah terdaftar
	if err := services.RegisterUser(&user); err != nil {
		// Jika ada error saat pendaftaran, periksa jenis error
		if err.Error() == "username already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already registered"})
			return
		}
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		// Log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
